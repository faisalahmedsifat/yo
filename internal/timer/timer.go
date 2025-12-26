package timer

import (
	"fmt"
	"time"

	"github.com/faisalahmedsifat/yo/internal/activity"
	"github.com/faisalahmedsifat/yo/internal/state"
)

// Milestones for timer notifications
const (
	Milestone100 = "100%"
	Milestone150 = "150%"
	Milestone200 = "200%"
)

// Status represents the current timer status
type Status struct {
	Running        bool
	Elapsed        time.Duration
	ElapsedHours   float64
	Threshold      time.Duration
	ThresholdHours float64
	Progress       float64
	Extensions     int
	Overtime       time.Duration
}

// GetStatus returns the current timer status
func GetStatus(s *state.State) *Status {
	if s.Timer.StartedAt.IsZero() {
		return &Status{Running: false}
	}

	elapsed := s.GetElapsed()
	threshold := time.Duration(s.Timer.ThresholdHours * float64(time.Hour))
	progress := s.GetProgress()

	status := &Status{
		Running:        true,
		Elapsed:        elapsed,
		ElapsedHours:   elapsed.Hours(),
		Threshold:      threshold,
		ThresholdHours: s.Timer.ThresholdHours,
		Progress:       progress,
		Extensions:     len(s.Timer.Extensions),
	}

	if elapsed > threshold {
		status.Overtime = elapsed - threshold
	}

	return status
}

// CheckMilestones checks if any notification milestones have been reached
// Returns list of newly reached milestones
func CheckMilestones(s *state.State) []string {
	var newMilestones []string
	progress := s.GetProgress()

	milestones := []struct {
		threshold float64
		name      string
	}{
		{100, Milestone100},
		{150, Milestone150},
		{200, Milestone200},
	}

	for _, m := range milestones {
		if progress >= m.threshold && !hasNotified(s, m.name) {
			newMilestones = append(newMilestones, m.name)
			s.Timer.Notifications = append(s.Timer.Notifications, m.name)

			// Log to activity
			activity.LogTimerMilestone(m.name, s.GetElapsed().Hours(), s.Timer.EstimatedHours)
		}
	}

	return newMilestones
}

// hasNotified checks if a milestone notification was already sent
func hasNotified(s *state.State, milestone string) bool {
	for _, n := range s.Timer.Notifications {
		if n == milestone {
			return true
		}
	}
	return false
}

// FormatDuration formats a duration as "Xh Ym"
func FormatDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}

// FormatHours formats hours as "Xh Ym"
func FormatHours(h float64) string {
	hours := int(h)
	minutes := int((h - float64(hours)) * 60)
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

// ParseDuration parses a duration string like "2h", "30m", "1.5h"
func ParseDuration(s string) (float64, error) {
	s = trimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty duration")
	}

	// Check for 'm' suffix (minutes)
	if len(s) > 1 && s[len(s)-1] == 'm' {
		var minutes float64
		if n, _ := fmt.Sscanf(s[:len(s)-1], "%f", &minutes); n == 1 {
			return minutes / 60, nil
		}
	}

	// Check for 'h' suffix (hours)
	if len(s) > 1 && s[len(s)-1] == 'h' {
		var hours float64
		if n, _ := fmt.Sscanf(s[:len(s)-1], "%f", &hours); n == 1 {
			return hours, nil
		}
	}

	// Try just a number (assume hours)
	var hours float64
	if n, _ := fmt.Sscanf(s, "%f", &hours); n == 1 {
		return hours, nil
	}

	return 0, fmt.Errorf("invalid duration format: %s", s)
}

func trimSpace(s string) string {
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

// ProgressBar returns a visual progress bar
func ProgressBar(progress float64, width int) string {
	filled := int((progress / 100) * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}

	bar := ""
	for i := 0; i < filled; i++ {
		bar += "█"
	}
	for i := filled; i < width; i++ {
		bar += "░"
	}
	return bar
}

// ProgressIndicator returns an emoji indicator based on progress
func ProgressIndicator(progress float64) string {
	switch {
	case progress < 80:
		return "🟢"
	case progress < 100:
		return "🟡"
	case progress < 150:
		return "🟠"
	default:
		return "🔴"
	}
}
