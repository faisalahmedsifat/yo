package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestDir(t *testing.T) (string, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "yo-workspace-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	oldWd, _ := os.Getwd()
	os.Chdir(tmpDir)

	cleanup := func() {
		os.Chdir(oldWd)
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

func TestInit(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	// Init should create .yo directory
	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	yoDir := filepath.Join(tmpDir, ".yo")

	// Verify directory was created
	if _, err := os.Stat(yoDir); os.IsNotExist(err) {
		t.Error(".yo directory was not created")
	}

	// Verify files were created
	expectedFiles := []string{
		"current_task.md",
		"backlog.md",
		"tech_debt_log.md",
		"state.json",
		"config.json",
		"activity.jsonl",
	}

	for _, f := range expectedFiles {
		path := filepath.Join(yoDir, f)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected file %s was not created", f)
		}
	}

	// Verify subdirectories were created
	expectedDirs := []string{"done", "sessions", "stats"}
	for _, d := range expectedDirs {
		path := filepath.Join(yoDir, d)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected directory %s was not created", d)
		}
	}
}

func TestInitAlreadyExists(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	// First init should succeed
	if err := Init(); err != nil {
		t.Fatalf("First init failed: %v", err)
	}

	// Second init should fail
	if err := Init(); err == nil {
		t.Error("Expected error for second init")
	}
}

func TestIsInitialized(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	// Should be false before init
	if IsInitialized() {
		t.Error("Expected IsInitialized to be false before Init")
	}

	// Init
	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	// Should be true after init
	if !IsInitialized() {
		t.Error("Expected IsInitialized to be true after Init")
	}
}

func TestGetCurrentTaskPath(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	path, err := GetCurrentTaskPath()
	if err != nil {
		t.Fatalf("GetCurrentTaskPath failed: %v", err)
	}

	expected := filepath.Join(tmpDir, ".yo", "current_task.md")
	if path != expected {
		t.Errorf("Expected path %s, got %s", expected, path)
	}
}

func TestGetBacklogPath(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	path, err := GetBacklogPath()
	if err != nil {
		t.Fatalf("GetBacklogPath failed: %v", err)
	}

	expected := filepath.Join(tmpDir, ".yo", "backlog.md")
	if path != expected {
		t.Errorf("Expected path %s, got %s", expected, path)
	}
}

func TestGetTechDebtPath(t *testing.T) {
	tmpDir, cleanup := setupTestDir(t)
	defer cleanup()

	path, err := GetTechDebtPath()
	if err != nil {
		t.Fatalf("GetTechDebtPath failed: %v", err)
	}

	expected := filepath.Join(tmpDir, ".yo", "tech_debt_log.md")
	if path != expected {
		t.Errorf("Expected path %s, got %s", expected, path)
	}
}

func TestCurrentTaskTemplate(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	taskPath, _ := GetCurrentTaskPath()
	content, err := os.ReadFile(taskPath)
	if err != nil {
		t.Fatalf("Failed to read task file: %v", err)
	}

	text := string(content)

	// Check for required sections
	requiredSections := []string{
		"🔴 RED LIGHT",
		"🟡 YELLOW LIGHT",
		"🟢 GREEN LIGHT",
		"What's the Problem?",
		"Impact",
		"Severity",
		"Solution Options",
		"Success Criteria",
	}

	for _, section := range requiredSections {
		if !contains(text, section) {
			t.Errorf("Expected template to contain '%s'", section)
		}
	}
}

func TestBacklogTemplate(t *testing.T) {
	_, cleanup := setupTestDir(t)
	defer cleanup()

	if err := Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	backlogPath, _ := GetBacklogPath()
	content, err := os.ReadFile(backlogPath)
	if err != nil {
		t.Fatalf("Failed to read backlog file: %v", err)
	}

	text := string(content)

	// Check for priority sections
	priorities := []string{"P0", "P1", "P2", "P3"}
	for _, p := range priorities {
		if !contains(text, p) {
			t.Errorf("Expected backlog to contain '%s' section", p)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
