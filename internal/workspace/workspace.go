package workspace

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/faisalahmedsifat/yo/internal/config"
	"github.com/faisalahmedsifat/yo/internal/state"
	"github.com/faisalahmedsifat/yo/internal/templates"
)

// Init creates the .yo workspace structure
func Init() error {
	yoDir, err := state.GetYoDir()
	if err != nil {
		return err
	}

	// Check if already initialized
	if _, err := os.Stat(yoDir); err == nil {
		return fmt.Errorf("workspace already initialized at %s", yoDir)
	}

	// Create directory structure
	dirs := []string{
		yoDir,
		filepath.Join(yoDir, "done"),
		filepath.Join(yoDir, "sessions"),
		filepath.Join(yoDir, "stats"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create template files
	files := map[string]string{
		"current_task.md":  templates.CurrentTask,
		"backlog.md":       templates.Backlog,
		"tech_debt_log.md": templates.TechDebtLog,
	}

	for name, content := range files {
		path := filepath.Join(yoDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("failed to create %s: %w", name, err)
		}
	}

	// Create state.json
	s := state.NewState()
	if err := s.Save(); err != nil {
		return fmt.Errorf("failed to create state.json: %w", err)
	}

	// Create config.json
	cfg := config.Default()
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to create config.json: %w", err)
	}

	// Create empty activity.jsonl
	activityPath := filepath.Join(yoDir, "activity.jsonl")
	if err := os.WriteFile(activityPath, []byte{}, 0644); err != nil {
		return fmt.Errorf("failed to create activity.jsonl: %w", err)
	}

	return nil
}

// IsInitialized checks if workspace is initialized
func IsInitialized() bool {
	yoDir, err := state.GetYoDir()
	if err != nil {
		return false
	}
	_, err = os.Stat(yoDir)
	return err == nil
}

// GetCurrentTaskPath returns path to current_task.md
func GetCurrentTaskPath() (string, error) {
	yoDir, err := state.GetYoDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(yoDir, "current_task.md"), nil
}

// GetBacklogPath returns path to backlog.md
func GetBacklogPath() (string, error) {
	yoDir, err := state.GetYoDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(yoDir, "backlog.md"), nil
}

// GetTechDebtPath returns path to tech_debt_log.md
func GetTechDebtPath() (string, error) {
	yoDir, err := state.GetYoDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(yoDir, "tech_debt_log.md"), nil
}
