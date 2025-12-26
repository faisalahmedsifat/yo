package watcher

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestHome(t *testing.T) (string, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "yo-watcher-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create global .yo directory
	globalYoDir := filepath.Join(tmpDir, ".yo")
	os.MkdirAll(globalYoDir, 0755)

	// Mock home directory by setting it
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)

	cleanup := func() {
		os.Setenv("HOME", oldHome)
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

func TestGetGlobalYoDir(t *testing.T) {
	tmpDir, cleanup := setupTestHome(t)
	defer cleanup()

	dir, err := GetGlobalYoDir()
	if err != nil {
		t.Fatalf("GetGlobalYoDir failed: %v", err)
	}

	expected := filepath.Join(tmpDir, ".yo")
	if dir != expected {
		t.Errorf("Expected %s, got %s", expected, dir)
	}
}

func TestGetGlobalConfigPath(t *testing.T) {
	tmpDir, cleanup := setupTestHome(t)
	defer cleanup()

	path, err := GetGlobalConfigPath()
	if err != nil {
		t.Fatalf("GetGlobalConfigPath failed: %v", err)
	}

	expected := filepath.Join(tmpDir, ".yo", "watcher_config.json")
	if path != expected {
		t.Errorf("Expected %s, got %s", expected, path)
	}
}

func TestGetPidFilePath(t *testing.T) {
	tmpDir, cleanup := setupTestHome(t)
	defer cleanup()

	path, err := GetPidFilePath()
	if err != nil {
		t.Fatalf("GetPidFilePath failed: %v", err)
	}

	expected := filepath.Join(tmpDir, ".yo", "watcher.pid")
	if path != expected {
		t.Errorf("Expected %s, got %s", expected, path)
	}
}

func TestLoadGlobalConfig(t *testing.T) {
	_, cleanup := setupTestHome(t)
	defer cleanup()

	// Load should create default config if not exists
	cfg, err := LoadGlobalConfig()
	if err != nil {
		t.Fatalf("LoadGlobalConfig failed: %v", err)
	}

	if cfg == nil {
		t.Fatal("Expected non-nil config")
	}

	if len(cfg.WatchDirs) == 0 {
		t.Error("Expected at least one watch dir")
	}
}

func TestGlobalConfigSaveLoad(t *testing.T) {
	_, cleanup := setupTestHome(t)
	defer cleanup()

	// Create and save config
	cfg := &GlobalConfig{
		WatchDirs:  []string{"/home/test/Dev", "/home/test/work"},
		CurrentDir: "/home/test/Dev/myapp",
	}

	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load config
	loaded, err := LoadGlobalConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if len(loaded.WatchDirs) != 2 {
		t.Errorf("Expected 2 watch dirs, got %d", len(loaded.WatchDirs))
	}

	if loaded.CurrentDir != "/home/test/Dev/myapp" {
		t.Errorf("Expected CurrentDir '/home/test/Dev/myapp', got '%s'", loaded.CurrentDir)
	}
}

func TestSetCurrentProject(t *testing.T) {
	_, cleanup := setupTestHome(t)
	defer cleanup()

	// Create initial config
	cfg := &GlobalConfig{WatchDirs: []string{"/home/test"}}
	cfg.Save()

	// Set current project
	if err := SetCurrentProject("/home/test/myproject"); err != nil {
		t.Fatalf("SetCurrentProject failed: %v", err)
	}

	// Load and verify
	loaded, err := LoadGlobalConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if loaded.CurrentDir != "/home/test/myproject" {
		t.Errorf("Expected CurrentDir '/home/test/myproject', got '%s'", loaded.CurrentDir)
	}
}

func TestExpandPath(t *testing.T) {
	home, _ := os.UserHomeDir()

	tests := []struct {
		input    string
		expected string
	}{
		{"~/Dev", filepath.Join(home, "Dev")},
		{"/absolute/path", "/absolute/path"},
		{"relative/path", "relative/path"},
	}

	for _, tt := range tests {
		result := expandPath(tt.input)
		if result != tt.expected {
			t.Errorf("expandPath(%s) = %s, expected %s", tt.input, result, tt.expected)
		}
	}
}

func TestIsRunningNotRunning(t *testing.T) {
	_, cleanup := setupTestHome(t)
	defer cleanup()

	// Should return false when no pid file
	running, pid := IsRunning()
	if running {
		t.Error("Expected IsRunning to be false when no pid file")
	}
	if pid != 0 {
		t.Errorf("Expected pid 0, got %d", pid)
	}
}

func TestIsRunningStalePid(t *testing.T) {
	tmpDir, cleanup := setupTestHome(t)
	defer cleanup()

	// Write a pid file with a non-existent process
	pidPath := filepath.Join(tmpDir, ".yo", "watcher.pid")
	os.WriteFile(pidPath, []byte("999999999"), 0644)

	// Should return false and clean up the stale pid file
	running, _ := IsRunning()
	if running {
		t.Error("Expected IsRunning to be false for stale pid")
	}
}

func TestNewWatcher(t *testing.T) {
	_, cleanup := setupTestHome(t)
	defer cleanup()

	// Create config
	cfg := &GlobalConfig{WatchDirs: []string{"/tmp"}}
	cfg.Save()

	w, err := New()
	if err != nil {
		t.Fatalf("Failed to create watcher: %v", err)
	}

	if w == nil {
		t.Fatal("Expected non-nil watcher")
	}

	if w.fsWatcher == nil {
		t.Error("Expected non-nil fsWatcher")
	}

	if w.repos == nil {
		t.Error("Expected non-nil repos map")
	}

	// Clean up
	w.Stop()
}

func TestWatcherFindRepo(t *testing.T) {
	_, cleanup := setupTestHome(t)
	defer cleanup()

	w := &Watcher{
		repos: map[string]bool{
			"/home/user/Dev/project1": true,
			"/home/user/Dev/project2": true,
		},
	}

	tests := []struct {
		path     string
		expected string
	}{
		{"/home/user/Dev/project1/src/main.go", "/home/user/Dev/project1"},
		{"/home/user/Dev/project2/lib/util.go", "/home/user/Dev/project2"},
		{"/home/user/other/file.go", ""},
	}

	for _, tt := range tests {
		result := w.findRepo(tt.path)
		if result != tt.expected {
			t.Errorf("findRepo(%s) = %s, expected %s", tt.path, result, tt.expected)
		}
	}
}
