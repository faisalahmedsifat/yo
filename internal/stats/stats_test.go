package stats

import (
	"testing"
	"time"

	"github.com/faisal/yo/internal/activity"
)

func TestCalculate(t *testing.T) {
	entries := []activity.Entry{
		{Type: activity.TypeTaskComplete, ActualHours: 4.0, EstimatedHours: 4.0},
		{Type: activity.TypeTaskComplete, ActualHours: 3.0, EstimatedHours: 2.0},
		{Type: activity.TypeEmergencyBypass},
		{Type: activity.TypeFileChange, Untracked: false},
		{Type: activity.TypeFileChange, Untracked: false},
		{Type: activity.TypeFileChange, Untracked: true},
	}

	stats := Calculate(entries)

	if stats.TasksCompleted != 2 {
		t.Errorf("Expected 2 tasks completed, got %d", stats.TasksCompleted)
	}

	if stats.Bypasses != 1 {
		t.Errorf("Expected 1 bypass, got %d", stats.Bypasses)
	}

	if stats.TotalChanges != 3 {
		t.Errorf("Expected 3 total changes, got %d", stats.TotalChanges)
	}

	if stats.OnTaskChanges != 2 {
		t.Errorf("Expected 2 on-task changes, got %d", stats.OnTaskChanges)
	}

	// Focus score should be 2/3 = 66.67%
	if stats.FocusScore < 66 || stats.FocusScore > 67 {
		t.Errorf("Expected focus score ~66.67%%, got %f", stats.FocusScore)
	}
}

func TestCalculateEmpty(t *testing.T) {
	stats := Calculate([]activity.Entry{})

	if stats.TasksCompleted != 0 {
		t.Error("Expected 0 tasks completed")
	}

	if stats.FocusScore != 100 {
		t.Errorf("Expected 100%% focus score for empty, got %f", stats.FocusScore)
	}
}

func TestGetWeekRange(t *testing.T) {
	// Test with a known date (Wednesday Dec 27, 2024)
	date := time.Date(2024, 12, 27, 12, 0, 0, 0, time.UTC)
	start, end := GetWeekRange(date)

	// Week should start on Monday Dec 23
	expectedStart := time.Date(2024, 12, 23, 0, 0, 0, 0, time.UTC)
	if !start.Equal(expectedStart) {
		t.Errorf("Expected week start %v, got %v", expectedStart, start)
	}

	// Week should end on Monday Dec 30
	expectedEnd := time.Date(2024, 12, 30, 0, 0, 0, 0, time.UTC)
	if !end.Equal(expectedEnd) {
		t.Errorf("Expected week end %v, got %v", expectedEnd, end)
	}
}

func TestGenerateInsights(t *testing.T) {
	// Good stats
	good := &WeekStats{
		AvgAccuracy:    95,
		FocusScore:     85,
		Bypasses:       1,
		TasksCompleted: 5,
	}

	insights := GenerateInsights(good)
	if len(insights.Messages) == 0 {
		t.Error("Expected some insights")
	}

	// Check for positive messages
	found := false
	for _, msg := range insights.Messages {
		if contains(msg, "Great") || contains(msg, "Excellent") {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected positive insight for good stats")
	}

	// Bad stats
	bad := &WeekStats{
		AvgAccuracy:    200,
		FocusScore:     40,
		Bypasses:       10,
		TasksCompleted: 0,
	}

	insights = GenerateInsights(bad)
	if len(insights.Messages) < 3 {
		t.Error("Expected multiple warnings for bad stats")
	}
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
