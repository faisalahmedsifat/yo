package activity

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupTestWorkspace(t *testing.T) (string, func()) {
	t.Helper()

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "yo-activity-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create .yo directory
	yoDir := filepath.Join(tmpDir, ".yo")
	if err := os.MkdirAll(yoDir, 0755); err != nil {
		t.Fatalf("Failed to create .yo dir: %v", err)
	}

	// Create empty activity.jsonl
	activityPath := filepath.Join(yoDir, "activity.jsonl")
	if err := os.WriteFile(activityPath, []byte{}, 0644); err != nil {
		t.Fatalf("Failed to create activity.jsonl: %v", err)
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

func TestAppendEntry(t *testing.T) {
	_, cleanup := setupTestWorkspace(t)
	defer cleanup()

	entry := Entry{
		Type:  TypeFileChange,
		Repo:  "~/Dev/myapp",
		File:  "src/main.go",
		Task:  "test_task",
		Stage: "green",
	}

	if err := Append(entry); err != nil {
		t.Fatalf("Failed to append entry: %v", err)
	}

	// Read back
	entries, err := QueryToday()
	if err != nil {
		t.Fatalf("Failed to query today: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(entries))
	}

	if entries[0].Type != TypeFileChange {
		t.Errorf("Expected type file_change, got %s", entries[0].Type)
	}

	if entries[0].File != "src/main.go" {
		t.Errorf("Expected file 'src/main.go', got %s", entries[0].File)
	}
}

func TestLogStageChange(t *testing.T) {
	_, cleanup := setupTestWorkspace(t)
	defer cleanup()

	if err := LogStageChange("red", "yellow", "test_task"); err != nil {
		t.Fatalf("Failed to log stage change: %v", err)
	}

	entries, err := QueryToday()
	if err != nil {
		t.Fatalf("Failed to query: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(entries))
	}

	if entries[0].Type != TypeStageChange {
		t.Errorf("Expected type stage_change, got %s", entries[0].Type)
	}

	if entries[0].From != "red" {
		t.Errorf("Expected from 'red', got %s", entries[0].From)
	}

	if entries[0].To != "yellow" {
		t.Errorf("Expected to 'yellow', got %s", entries[0].To)
	}
}

func TestLogTaskComplete(t *testing.T) {
	_, cleanup := setupTestWorkspace(t)
	defer cleanup()

	if err := LogTaskComplete("test_task", 4.5, 4.0); err != nil {
		t.Fatalf("Failed to log task complete: %v", err)
	}

	entries, err := QueryToday()
	if err != nil {
		t.Fatalf("Failed to query: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(entries))
	}

	if entries[0].Type != TypeTaskComplete {
		t.Errorf("Expected type task_complete, got %s", entries[0].Type)
	}

	if entries[0].ActualHours != 4.5 {
		t.Errorf("Expected actual hours 4.5, got %f", entries[0].ActualHours)
	}

	if entries[0].EstimatedHours != 4.0 {
		t.Errorf("Expected estimated hours 4.0, got %f", entries[0].EstimatedHours)
	}
}

func TestLogEmergencyBypass(t *testing.T) {
	_, cleanup := setupTestWorkspace(t)
	defer cleanup()

	if err := LogEmergencyBypass("production down", 1, 3); err != nil {
		t.Fatalf("Failed to log bypass: %v", err)
	}

	entries, err := QueryToday()
	if err != nil {
		t.Fatalf("Failed to query: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(entries))
	}

	if entries[0].Reason != "production down" {
		t.Errorf("Expected reason 'production down', got %s", entries[0].Reason)
	}

	if entries[0].CountToday != 1 {
		t.Errorf("Expected count today 1, got %d", entries[0].CountToday)
	}
}

func TestQueryTimeRange(t *testing.T) {
	_, cleanup := setupTestWorkspace(t)
	defer cleanup()

	// Add multiple entries
	for i := 0; i < 5; i++ {
		Append(Entry{
			Type: TypeFileChange,
			File: "test.go",
		})
	}

	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24 * time.Hour)

	entries, err := Query(start, end)
	if err != nil {
		t.Fatalf("Failed to query: %v", err)
	}

	if len(entries) != 5 {
		t.Errorf("Expected 5 entries, got %d", len(entries))
	}
}

func TestQueryEmptyLog(t *testing.T) {
	_, cleanup := setupTestWorkspace(t)
	defer cleanup()

	entries, err := QueryToday()
	if err != nil {
		t.Fatalf("Failed to query: %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(entries))
	}
}
