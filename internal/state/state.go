package state

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// State represents the current state of the yo workspace
type State struct {
	Version           string            `json:"version"`
	CurrentStage      string            `json:"current_stage"` // none, red, yellow, green
	CurrentTaskID     string            `json:"current_task_id"`
	CurrentTaskRepo   string            `json:"current_task_repo"`
	Timer             Timer             `json:"timer"`
	Session           Session           `json:"session"`
	EmergencyBypasses EmergencyBypasses `json:"emergency_bypasses"`
}

// Timer tracks the current task timer
type Timer struct {
	StartedAt      time.Time   `json:"started_at,omitempty"`
	EstimatedHours float64     `json:"estimated_hours"`
	ThresholdHours float64     `json:"threshold_hours"`
	Paused         bool        `json:"paused"`
	Extensions     []Extension `json:"extensions,omitempty"`
	Notifications  []string    `json:"notifications,omitempty"` // track which milestones were notified
}

// Extension represents a timer extension
type Extension struct {
	At     time.Time `json:"at"`
	Hours  float64   `json:"hours"`
	Reason string    `json:"reason"`
}

// Session tracks the current work session
type Session struct {
	StartedAt time.Time `json:"started_at,omitempty"`
	Active    bool      `json:"active"`
}

// EmergencyBypasses tracks bypass usage
type EmergencyBypasses struct {
	Today     int    `json:"today"`
	ThisWeek  int    `json:"this_week"`
	LastReset string `json:"last_reset"`
}

// NewState creates a new default state
func NewState() *State {
	return &State{
		Version:      "1.0.0",
		CurrentStage: "none",
		Timer: Timer{
			EstimatedHours: 0,
			ThresholdHours: 0,
			Paused:         false,
		},
		Session: Session{
			Active: false,
		},
		EmergencyBypasses: EmergencyBypasses{
			Today:     0,
			ThisWeek:  0,
			LastReset: time.Now().Format("2006-01-02"),
		},
	}
}

// GetYoDir returns the path to the .yo directory
func GetYoDir() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}
	return filepath.Join(cwd, ".yo"), nil
}

// GetStatePath returns the path to state.json
func GetStatePath() (string, error) {
	yoDir, err := GetYoDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(yoDir, "state.json"), nil
}

// Load loads state from disk
func Load() (*State, error) {
	statePath, err := GetStatePath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(statePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("workspace not initialized. Run 'yo init' first")
		}
		return nil, fmt.Errorf("failed to read state: %w", err)
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state: %w", err)
	}

	// Reset daily/weekly counters if needed
	state.resetBypassCounters()

	return &state, nil
}

// Save saves state to disk
func (s *State) Save() error {
	statePath, err := GetStatePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(statePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write state: %w", err)
	}

	return nil
}

// resetBypassCounters resets daily/weekly bypass counters if needed
func (s *State) resetBypassCounters() {
	today := time.Now().Format("2006-01-02")
	lastReset, err := time.Parse("2006-01-02", s.EmergencyBypasses.LastReset)
	if err != nil {
		s.EmergencyBypasses.LastReset = today
		s.EmergencyBypasses.Today = 0
		s.EmergencyBypasses.ThisWeek = 0
		return
	}

	now := time.Now()
	// Reset daily counter
	if now.Format("2006-01-02") != lastReset.Format("2006-01-02") {
		s.EmergencyBypasses.Today = 0
	}

	// Reset weekly counter (if different week)
	_, lastWeek := lastReset.ISOWeek()
	_, currentWeek := now.ISOWeek()
	if currentWeek != lastWeek || now.Year() != lastReset.Year() {
		s.EmergencyBypasses.ThisWeek = 0
	}

	s.EmergencyBypasses.LastReset = today
}

// StartTimer starts the timer for the current task
func (s *State) StartTimer(estimatedHours float64) {
	s.Timer = Timer{
		StartedAt:      time.Now(),
		EstimatedHours: estimatedHours,
		ThresholdHours: estimatedHours,
		Paused:         false,
		Extensions:     []Extension{},
		Notifications:  []string{},
	}
}

// GetElapsed returns the elapsed time since timer started
func (s *State) GetElapsed() time.Duration {
	if s.Timer.StartedAt.IsZero() {
		return 0
	}
	return time.Since(s.Timer.StartedAt)
}

// GetProgress returns the progress percentage (elapsed / threshold)
func (s *State) GetProgress() float64 {
	if s.Timer.ThresholdHours == 0 {
		return 0
	}
	elapsedHours := s.GetElapsed().Hours()
	return (elapsedHours / s.Timer.ThresholdHours) * 100
}

// ExtendTimer extends the timer threshold
func (s *State) ExtendTimer(hours float64, reason string) {
	s.Timer.Extensions = append(s.Timer.Extensions, Extension{
		At:     time.Now(),
		Hours:  hours,
		Reason: reason,
	})
	s.Timer.ThresholdHours += hours
}

// StopTimer stops the timer
func (s *State) StopTimer() {
	s.Timer.StartedAt = time.Time{}
	s.Timer.Paused = true
}

// SetStage sets the current stage
func (s *State) SetStage(stage string) {
	s.CurrentStage = stage
}

// StartSession starts a new work session
func (s *State) StartSession() {
	s.Session = Session{
		StartedAt: time.Now(),
		Active:    true,
	}
}

// EndSession ends the current work session
func (s *State) EndSession() {
	s.Session.Active = false
}
