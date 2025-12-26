package state

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupTestWorkspace(t *testing.T) (string, func()) {
	t.Helper()

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "yo-test-*")
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

func TestNewState(t *testing.T) {
	s := NewState()

	if s.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", s.Version)
	}

	if s.CurrentStage != "none" {
		t.Errorf("Expected stage 'none', got %s", s.CurrentStage)
	}

	if s.Session.Active {
		t.Error("Expected session to be inactive")
	}
}

func TestStateSaveLoad(t *testing.T) {
	_, cleanup := setupTestWorkspace(t)
	defer cleanup()

	// Create and save state
	s := NewState()
	s.CurrentStage = "green"
	s.CurrentTaskID = "test_task"
	s.Timer.EstimatedHours = 4.0

	if err := s.Save(); err != nil {
		t.Fatalf("Failed to save state: %v", err)
	}

	// Load state
	loaded, err := Load()
	if err != nil {
		t.Fatalf("Failed to load state: %v", err)
	}

	if loaded.CurrentStage != "green" {
		t.Errorf("Expected stage 'green', got %s", loaded.CurrentStage)
	}

	if loaded.CurrentTaskID != "test_task" {
		t.Errorf("Expected task 'test_task', got %s", loaded.CurrentTaskID)
	}

	if loaded.Timer.EstimatedHours != 4.0 {
		t.Errorf("Expected estimated hours 4.0, got %f", loaded.Timer.EstimatedHours)
	}
}

func TestTimerOperations(t *testing.T) {
	s := NewState()

	// Start timer
	s.StartTimer(2.0)

	if s.Timer.EstimatedHours != 2.0 {
		t.Errorf("Expected estimated hours 2.0, got %f", s.Timer.EstimatedHours)
	}

	if s.Timer.ThresholdHours != 2.0 {
		t.Errorf("Expected threshold hours 2.0, got %f", s.Timer.ThresholdHours)
	}

	if s.Timer.StartedAt.IsZero() {
		t.Error("Expected timer to be started")
	}

	// Extend timer
	s.ExtendTimer(1.0, "test reason")

	if s.Timer.ThresholdHours != 3.0 {
		t.Errorf("Expected threshold hours 3.0, got %f", s.Timer.ThresholdHours)
	}

	if len(s.Timer.Extensions) != 1 {
		t.Errorf("Expected 1 extension, got %d", len(s.Timer.Extensions))
	}

	if s.Timer.Extensions[0].Reason != "test reason" {
		t.Errorf("Expected reason 'test reason', got %s", s.Timer.Extensions[0].Reason)
	}
}

func TestTimerProgress(t *testing.T) {
	s := NewState()
	s.Timer.StartedAt = time.Now().Add(-1 * time.Hour)
	s.Timer.ThresholdHours = 2.0

	elapsed := s.GetElapsed()
	if elapsed < 59*time.Minute || elapsed > 61*time.Minute {
		t.Errorf("Expected elapsed ~1h, got %v", elapsed)
	}

	progress := s.GetProgress()
	if progress < 45 || progress > 55 {
		t.Errorf("Expected progress ~50%%, got %f", progress)
	}
}

func TestSessionOperations(t *testing.T) {
	s := NewState()

	if s.Session.Active {
		t.Error("Expected session to be inactive initially")
	}

	s.StartSession()

	if !s.Session.Active {
		t.Error("Expected session to be active after start")
	}

	if s.Session.StartedAt.IsZero() {
		t.Error("Expected session start time to be set")
	}

	s.EndSession()

	if s.Session.Active {
		t.Error("Expected session to be inactive after end")
	}
}

func TestStageTransitions(t *testing.T) {
	s := NewState()

	stages := []string{"none", "red", "yellow", "green"}
	for _, stage := range stages {
		s.SetStage(stage)
		if s.CurrentStage != stage {
			t.Errorf("Expected stage %s, got %s", stage, s.CurrentStage)
		}
	}
}

func TestBypassCounterReset(t *testing.T) {
	s := NewState()

	// Set bypass counts
	s.EmergencyBypasses.Today = 2
	s.EmergencyBypasses.ThisWeek = 5
	s.EmergencyBypasses.LastReset = time.Now().Add(-48 * time.Hour).Format("2006-01-02")

	// Reset should happen when loading
	s.resetBypassCounters()

	if s.EmergencyBypasses.Today != 0 {
		t.Errorf("Expected today count to reset, got %d", s.EmergencyBypasses.Today)
	}

	// LastReset should be updated to today
	if s.EmergencyBypasses.LastReset != time.Now().Format("2006-01-02") {
		t.Errorf("Expected last reset to be today, got %s", s.EmergencyBypasses.LastReset)
	}
}
