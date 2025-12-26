package milestone

import (
	"testing"

	"github.com/faisal/yo/internal/state"
)

func TestAll(t *testing.T) {
	milestones := All()

	if len(milestones) != 5 {
		t.Errorf("Expected 5 milestones, got %d", len(milestones))
	}

	// Check first milestone
	if milestones[0].Name != "Clear Launch Blockers" {
		t.Errorf("Expected first milestone 'Clear Launch Blockers', got '%s'", milestones[0].Name)
	}

	// Check last milestone
	if milestones[4].Name != "Production Ready" {
		t.Errorf("Expected last milestone 'Production Ready', got '%s'", milestones[4].Name)
	}
}

func TestCurrent(t *testing.T) {
	s := state.NewState()
	s.Milestone.Current = 0

	m := Current(s)
	if m == nil {
		t.Fatal("Expected non-nil milestone")
	}

	if m.ID != 0 {
		t.Errorf("Expected milestone ID 0, got %d", m.ID)
	}
}

func TestAdvance(t *testing.T) {
	s := state.NewState()
	s.Milestone.Current = 0

	// Advance should succeed
	if !Advance(s) {
		t.Error("Expected advance to succeed")
	}

	if s.Milestone.Current != 1 {
		t.Errorf("Expected current milestone 1, got %d", s.Milestone.Current)
	}

	// Advance to end
	s.Milestone.Current = 4
	if Advance(s) {
		t.Error("Expected advance to fail at end")
	}
}

func TestIsComplete(t *testing.T) {
	s := state.NewState()

	s.Milestone.Current = 0
	if IsComplete(s) {
		t.Error("Expected not complete at milestone 0")
	}

	s.Milestone.Current = 5
	if !IsComplete(s) {
		t.Error("Expected complete at milestone 5")
	}
}

func TestProgress(t *testing.T) {
	s := state.NewState()
	s.Milestone.Current = 0

	// 0 of 5 criteria complete
	progress := Progress(s, 0)
	if progress != 0 {
		t.Errorf("Expected 0%% progress, got %f", progress)
	}

	// 2 of 5 criteria complete
	progress = Progress(s, 2)
	if progress != 40 {
		t.Errorf("Expected 40%% progress, got %f", progress)
	}

	// All complete
	progress = Progress(s, 5)
	if progress != 100 {
		t.Errorf("Expected 100%% progress, got %f", progress)
	}
}

func TestGetByID(t *testing.T) {
	m := GetByID(0)
	if m == nil {
		t.Fatal("Expected non-nil milestone for ID 0")
	}
	if m.Name != "Clear Launch Blockers" {
		t.Errorf("Expected 'Clear Launch Blockers', got '%s'", m.Name)
	}

	m = GetByID(100)
	if m != nil {
		t.Error("Expected nil for invalid ID")
	}

	m = GetByID(-1)
	if m != nil {
		t.Error("Expected nil for negative ID")
	}
}

func TestGetStatus(t *testing.T) {
	s := state.NewState()
	s.Milestone.Current = 2

	status := GetStatus(s)

	if status.Current != 2 {
		t.Errorf("Expected current 2, got %d", status.Current)
	}

	if status.Total != 5 {
		t.Errorf("Expected total 5, got %d", status.Total)
	}

	// First two should be completed
	if !status.Completed[0] || !status.Completed[1] {
		t.Error("Expected first two milestones to be completed")
	}

	// Current and beyond should not be completed
	if status.Completed[2] {
		t.Error("Expected milestone 2 to not be completed")
	}
}
