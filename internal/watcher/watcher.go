package watcher

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// GlobalConfig holds the global watcher configuration
type GlobalConfig struct {
	WatchDirs  []string `json:"watch_dirs"`
	PidFile    string   `json:"pid_file"`
	CurrentDir string   `json:"current_dir"` // The active project directory
}

// Watcher manages file watching across all projects
type Watcher struct {
	fsWatcher *fsnotify.Watcher
	config    *GlobalConfig
	repos     map[string]bool // Detected git repositories
	debouncer map[string]time.Time
	mu        sync.Mutex
	stopChan  chan struct{}
	logFile   *os.File
}

// GetGlobalYoDir returns the global .yo directory in user's home
func GetGlobalYoDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".yo"), nil
}

// GetGlobalConfigPath returns path to global config
func GetGlobalConfigPath() (string, error) {
	globalDir, err := GetGlobalYoDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(globalDir, "watcher_config.json"), nil
}

// GetPidFilePath returns path to the daemon PID file
func GetPidFilePath() (string, error) {
	globalDir, err := GetGlobalYoDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(globalDir, "watcher.pid"), nil
}

// LoadGlobalConfig loads or creates the global watcher config
func LoadGlobalConfig() (*GlobalConfig, error) {
	configPath, err := GetGlobalConfigPath()
	if err != nil {
		return nil, err
	}

	// Create global dir if needed
	globalDir, _ := GetGlobalYoDir()
	os.MkdirAll(globalDir, 0755)

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create default config
			home, _ := os.UserHomeDir()
			cfg := &GlobalConfig{
				WatchDirs: []string{filepath.Join(home, "Dev")},
			}
			cfg.Save()
			return cfg, nil
		}
		return nil, err
	}

	var cfg GlobalConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Save saves the global config
func (c *GlobalConfig) Save() error {
	configPath, err := GetGlobalConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// New creates a new watcher
func New() (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	config, err := LoadGlobalConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &Watcher{
		fsWatcher: fsWatcher,
		config:    config,
		repos:     make(map[string]bool),
		debouncer: make(map[string]time.Time),
		stopChan:  make(chan struct{}),
	}, nil
}

// Start begins watching all configured directories
func (w *Watcher) Start() error {
	// Discover git repos in watch dirs
	if err := w.discoverRepos(); err != nil {
		return err
	}

	// Add all repos to watcher
	for repo := range w.repos {
		if err := w.addDirRecursively(repo); err != nil {
			fmt.Printf("Warning: failed to watch %s: %v\n", repo, err)
		}
	}

	// Write PID file
	if err := w.writePidFile(); err != nil {
		return err
	}

	fmt.Printf("🔍 Watching %d repositories\n", len(w.repos))
	for repo := range w.repos {
		fmt.Printf("   %s\n", repo)
	}

	go w.eventLoop()
	return nil
}

// Stop stops the watcher
func (w *Watcher) Stop() {
	close(w.stopChan)
	w.fsWatcher.Close()
	w.removePidFile()
	if w.logFile != nil {
		w.logFile.Close()
	}
}

// discoverRepos finds all git repositories in watch directories
func (w *Watcher) discoverRepos() error {
	for _, watchDir := range w.config.WatchDirs {
		watchDir = expandPath(watchDir)

		err := filepath.Walk(watchDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Skip inaccessible paths
			}

			if info.IsDir() && info.Name() == ".git" {
				repoDir := filepath.Dir(path)
				w.repos[repoDir] = true
				return filepath.SkipDir
			}

			// Skip common non-project directories
			if info.IsDir() {
				name := info.Name()
				if name == "node_modules" || name == "vendor" || name == ".git" ||
					name == "__pycache__" || name == ".cache" || name == "dist" ||
					name == "build" || name == ".next" {
					return filepath.SkipDir
				}
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

// addDirRecursively adds a directory and subdirectories to the watcher
func (w *Watcher) addDirRecursively(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			name := info.Name()
			// Skip common non-project directories
			if name == "node_modules" || name == "vendor" || name == ".git" ||
				name == "__pycache__" || name == ".cache" || name == "dist" ||
				name == "build" || name == ".next" || name == ".yo" {
				return filepath.SkipDir
			}

			return w.fsWatcher.Add(path)
		}

		return nil
	})
}

// eventLoop handles file system events
func (w *Watcher) eventLoop() {
	for {
		select {
		case <-w.stopChan:
			return

		case event, ok := <-w.fsWatcher.Events:
			if !ok {
				return
			}

			if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
				w.handleFileChange(event.Name)
			}

		case err, ok := <-w.fsWatcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("Watcher error: %v\n", err)
		}
	}
}

// handleFileChange processes a file change event
func (w *Watcher) handleFileChange(path string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Debounce: ignore if we saw this file in the last second
	if lastSeen, ok := w.debouncer[path]; ok {
		if time.Since(lastSeen) < time.Second {
			return
		}
	}
	w.debouncer[path] = time.Now()

	// Skip hidden files and directories
	if strings.Contains(path, "/.") {
		return
	}

	// Find which repo this file belongs to
	repo := w.findRepo(path)
	if repo == "" {
		return
	}

	// Check if this is the current active project
	currentDir := w.config.CurrentDir
	isUntracked := currentDir != "" && !strings.HasPrefix(path, currentDir)

	// Log to the current project's activity.jsonl if there is one
	w.logActivity(path, repo, isUntracked)
}

// findRepo finds which repo a file belongs to
func (w *Watcher) findRepo(path string) string {
	for repo := range w.repos {
		if strings.HasPrefix(path, repo) {
			return repo
		}
	}
	return ""
}

// logActivity logs a file change to the activity log
func (w *Watcher) logActivity(path, repo string, untracked bool) {
	// Find the .yo directory for the current project
	currentDir := w.config.CurrentDir
	if currentDir == "" {
		return // No active project
	}

	activityPath := filepath.Join(currentDir, ".yo", "activity.jsonl")
	if _, err := os.Stat(filepath.Dir(activityPath)); os.IsNotExist(err) {
		return // No .yo directory in current project
	}

	// Create log entry
	entry := map[string]interface{}{
		"ts":        time.Now().Format(time.RFC3339Nano),
		"type":      "file_change",
		"repo":      repo,
		"file":      strings.TrimPrefix(path, repo+"/"),
		"untracked": untracked,
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return
	}

	f, err := os.OpenFile(activityPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	f.WriteString(string(data) + "\n")
}

// writePidFile writes the current process ID to the pid file
func (w *Watcher) writePidFile() error {
	pidPath, err := GetPidFilePath()
	if err != nil {
		return err
	}

	pid := fmt.Sprintf("%d", os.Getpid())
	return os.WriteFile(pidPath, []byte(pid), 0644)
}

// removePidFile removes the pid file
func (w *Watcher) removePidFile() {
	pidPath, _ := GetPidFilePath()
	os.Remove(pidPath)
}

// IsRunning checks if the watcher daemon is already running
func IsRunning() (bool, int) {
	pidPath, err := GetPidFilePath()
	if err != nil {
		return false, 0
	}

	data, err := os.ReadFile(pidPath)
	if err != nil {
		return false, 0
	}

	var pid int
	fmt.Sscanf(string(data), "%d", &pid)

	// Check if process exists
	process, err := os.FindProcess(pid)
	if err != nil {
		return false, 0
	}

	// On Unix, FindProcess always succeeds, so we need to send signal 0
	// to check if process actually exists
	err = process.Signal(os.Signal(nil))
	if err != nil {
		os.Remove(pidPath) // Clean up stale pid file
		return false, 0
	}

	return true, pid
}

// SetCurrentProject updates the current active project directory
func SetCurrentProject(dir string) error {
	config, err := LoadGlobalConfig()
	if err != nil {
		return err
	}

	config.CurrentDir = dir
	return config.Save()
}

// expandPath expands ~ to home directory
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}
