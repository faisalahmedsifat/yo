package activity

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/faisal/yo/internal/state"
)

// EntryType defines the type of activity entry
type EntryType string

const (
	TypeFileChange      EntryType = "file_change"
	TypeStageChange     EntryType = "stage_change"
	TypeSessionEnd      EntryType = "session_end"
	TypeTimerMilestone  EntryType = "timer_milestone"
	TypeEmergencyBypass EntryType = "emergency_bypass"
	TypeTaskComplete    EntryType = "task_complete"
)

// Entry represents a single activity log entry
type Entry struct {
	Timestamp time.Time `json:"ts"`
	Type      EntryType `json:"type"`

	// For file_change
	Repo      string `json:"repo,omitempty"`
	File      string `json:"file,omitempty"`
	Task      string `json:"task,omitempty"`
	Stage     string `json:"stage,omitempty"`
	Untracked bool   `json:"untracked,omitempty"`

	// For stage_change
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`

	// For session_end
	DurationMinutes int     `json:"duration_minutes,omitempty"`
	PrimaryRepo     string  `json:"primary_repo,omitempty"`
	FocusPercent    float64 `json:"focus_percent,omitempty"`

	// For timer_milestone
	Milestone      string  `json:"milestone,omitempty"`
	ActualHours    float64 `json:"actual_hours,omitempty"`
	EstimatedHours float64 `json:"estimated_hours,omitempty"`

	// For emergency_bypass
	Reason     string `json:"reason,omitempty"`
	CountToday int    `json:"count_today,omitempty"`
	CountWeek  int    `json:"count_week,omitempty"`
}

// getActivityPath returns the path to activity.jsonl
func getActivityPath() (string, error) {
	yoDir, err := state.GetYoDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(yoDir, "activity.jsonl"), nil
}

// Append appends an activity entry to the log
func Append(entry Entry) error {
	activityPath, err := getActivityPath()
	if err != nil {
		return err
	}

	entry.Timestamp = time.Now()

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal entry: %w", err)
	}

	f, err := os.OpenFile(activityPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open activity log: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(string(data) + "\n"); err != nil {
		return fmt.Errorf("failed to write entry: %w", err)
	}

	return nil
}

// LogStageChange logs a stage transition
func LogStageChange(from, to, task string) error {
	return Append(Entry{
		Type: TypeStageChange,
		From: from,
		To:   to,
		Task: task,
	})
}

// LogTaskComplete logs task completion
func LogTaskComplete(taskID string, actualHours, estimatedHours float64) error {
	return Append(Entry{
		Type:           TypeTaskComplete,
		Task:           taskID,
		ActualHours:    actualHours,
		EstimatedHours: estimatedHours,
	})
}

// LogSessionEnd logs the end of a work session
func LogSessionEnd(durationMinutes int, primaryRepo string, focusPercent float64) error {
	return Append(Entry{
		Type:            TypeSessionEnd,
		DurationMinutes: durationMinutes,
		PrimaryRepo:     primaryRepo,
		FocusPercent:    focusPercent,
	})
}

// LogTimerMilestone logs a timer milestone (100%, 150%, 200%)
func LogTimerMilestone(milestone string, actualHours, estimatedHours float64) error {
	return Append(Entry{
		Type:           TypeTimerMilestone,
		Milestone:      milestone,
		ActualHours:    actualHours,
		EstimatedHours: estimatedHours,
	})
}

// LogEmergencyBypass logs an emergency bypass
func LogEmergencyBypass(reason string, countToday, countWeek int) error {
	return Append(Entry{
		Type:       TypeEmergencyBypass,
		Reason:     reason,
		CountToday: countToday,
		CountWeek:  countWeek,
	})
}

// Query returns entries matching the given time range
func Query(start, end time.Time) ([]Entry, error) {
	activityPath, err := getActivityPath()
	if err != nil {
		return nil, err
	}

	f, err := os.Open(activityPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []Entry{}, nil
		}
		return nil, fmt.Errorf("failed to open activity log: %w", err)
	}
	defer f.Close()

	var entries []Entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var entry Entry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue // Skip malformed entries
		}
		if (entry.Timestamp.After(start) || entry.Timestamp.Equal(start)) &&
			(entry.Timestamp.Before(end) || entry.Timestamp.Equal(end)) {
			entries = append(entries, entry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read activity log: %w", err)
	}

	return entries, nil
}

// QueryToday returns today's entries
func QueryToday() ([]Entry, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	end := start.Add(24 * time.Hour)
	return Query(start, end)
}

// QueryYesterday returns yesterday's entries
func QueryYesterday() ([]Entry, error) {
	now := time.Now()
	start := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, now.Location())
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return Query(start, end)
}

// QueryThisWeek returns entries from the current week (Monday to now)
func QueryThisWeek() ([]Entry, error) {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	start := time.Date(now.Year(), now.Month(), now.Day()-weekday+1, 0, 0, 0, 0, now.Location())
	end := now.Add(24 * time.Hour)
	return Query(start, end)
}
