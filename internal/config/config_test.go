package config

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestWorkspace(t *testing.T) (string, func()) {
	t.Helper()

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "yo-config-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create .yo directory
	yoDir := filepath.Join(tmpDir, ".yo")
	if err := os.MkdirAll(yoDir, 0755); err != nil {
		t.Fatalf("Failed to create .yo dir: %v", err)
	}

	// Change to temp directory
	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)

	cleanup := func() {
		os.Chdir(oldWd)
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.Notifications != true {
		t.Error("Expected notifications to be true by default")
	}

	if cfg.MaxBypassDay != 1 {
		t.Errorf("Expected MaxBypassDay=1, got %d", cfg.MaxBypassDay)
	}

	if cfg.MaxBypassWeek != 5 {
		t.Errorf("Expected MaxBypassWeek=5, got %d", cfg.MaxBypassWeek)
	}

	if len(cfg.WatchDirs) == 0 {
		t.Error("Expected at least one watch dir")
	}
}

func TestConfigSaveLoad(t *testing.T) {
	_, cleanup := setupTestWorkspace(t)
	defer cleanup()

	// Create and save config
	cfg := Default()
	cfg.Notifications = false
	cfg.Editor = "vim"

	if err := cfg.Save(); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Load config
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if loaded.Notifications != false {
		t.Error("Expected notifications to be false")
	}

	if loaded.Editor != "vim" {
		t.Errorf("Expected editor 'vim', got '%s'", loaded.Editor)
	}
}

func TestConfigSet(t *testing.T) {
	cfg := Default()

	// Test setting notifications
	if err := cfg.Set("notifications", "off"); err != nil {
		t.Fatalf("Failed to set notifications: %v", err)
	}
	if cfg.Notifications != false {
		t.Error("Expected notifications to be false")
	}

	if err := cfg.Set("notifications", "on"); err != nil {
		t.Fatalf("Failed to set notifications: %v", err)
	}
	if cfg.Notifications != true {
		t.Error("Expected notifications to be true")
	}

	// Test setting editor
	if err := cfg.Set("editor", "nano"); err != nil {
		t.Fatalf("Failed to set editor: %v", err)
	}
	if cfg.Editor != "nano" {
		t.Errorf("Expected editor 'nano', got '%s'", cfg.Editor)
	}

	// Test unknown key
	if err := cfg.Set("unknown_key", "value"); err == nil {
		t.Error("Expected error for unknown key")
	}
}

func TestConfigGet(t *testing.T) {
	cfg := Default()
	cfg.Notifications = true
	cfg.Editor = "code"

	// Test getting notifications
	val, err := cfg.Get("notifications")
	if err != nil {
		t.Fatalf("Failed to get notifications: %v", err)
	}
	if val != "on" {
		t.Errorf("Expected 'on', got '%s'", val)
	}

	// Test getting editor
	val, err = cfg.Get("editor")
	if err != nil {
		t.Fatalf("Failed to get editor: %v", err)
	}
	if val != "code" {
		t.Errorf("Expected 'code', got '%s'", val)
	}

	// Test getting watch_dirs
	val, err = cfg.Get("watch_dirs")
	if err != nil {
		t.Fatalf("Failed to get watch_dirs: %v", err)
	}
	if val == "" {
		t.Error("Expected non-empty watch_dirs")
	}

	// Test unknown key
	_, err = cfg.Get("unknown")
	if err == nil {
		t.Error("Expected error for unknown key")
	}
}

func TestLoadNonExistent(t *testing.T) {
	_, cleanup := setupTestWorkspace(t)
	defer cleanup()

	// Load should return default config if file doesn't exist
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	if cfg == nil {
		t.Error("Expected non-nil config")
	}

	// Should have default values
	if cfg.MaxBypassDay != 1 {
		t.Errorf("Expected MaxBypassDay=1, got %d", cfg.MaxBypassDay)
	}
}
