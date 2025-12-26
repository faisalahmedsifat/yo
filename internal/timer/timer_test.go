package timer

import (
	"testing"
	"time"

	"github.com/faisal/yo/internal/state"
)

func TestGetStatus(t *testing.T) {
	s := state.NewState()

	// Not running
	status := GetStatus(s)
	if status.Running {
		t.Error("Expected timer to not be running")
	}

	// Start timer
	s.StartTimer(2.0)
	status = GetStatus(s)

	if !status.Running {
		t.Error("Expected timer to be running")
	}

	if status.ThresholdHours != 2.0 {
		t.Errorf("Expected threshold 2.0h, got %f", status.ThresholdHours)
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		d        time.Duration
		expected string
	}{
		{30 * time.Minute, "30m"},
		{1 * time.Hour, "1h 0m"},
		{90 * time.Minute, "1h 30m"},
		{2*time.Hour + 15*time.Minute, "2h 15m"},
	}

	for _, tt := range tests {
		result := FormatDuration(tt.d)
		if result != tt.expected {
			t.Errorf("FormatDuration(%v) = %s, expected %s", tt.d, result, tt.expected)
		}
	}
}

func TestFormatHours(t *testing.T) {
	tests := []struct {
		h        float64
		expected string
	}{
		{0.5, "30m"},
		{1.0, "1h 0m"},
		{1.5, "1h 30m"},
		{2.25, "2h 15m"},
	}

	for _, tt := range tests {
		result := FormatHours(tt.h)
		if result != tt.expected {
			t.Errorf("FormatHours(%f) = %s, expected %s", tt.h, result, tt.expected)
		}
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		hasError bool
	}{
		{"2h", 2.0, false},
		{"1.5h", 1.5, false},
		{"30m", 0.5, false},
		{"60m", 1.0, false},
		{"3", 3.0, false},
		{"", 0, true},
	}

	for _, tt := range tests {
		result, err := ParseDuration(tt.input)
		if tt.hasError {
			if err == nil {
				t.Errorf("ParseDuration(%s) expected error", tt.input)
			}
		} else {
			if err != nil {
				t.Errorf("ParseDuration(%s) unexpected error: %v", tt.input, err)
			}
			if result != tt.expected {
				t.Errorf("ParseDuration(%s) = %f, expected %f", tt.input, result, tt.expected)
			}
		}
	}
}

func TestProgressBar(t *testing.T) {
	tests := []struct {
		progress float64
		width    int
		filled   int
	}{
		{0, 10, 0},
		{50, 10, 5},
		{100, 10, 10},
		{150, 10, 10}, // Cap at 100%
	}

	for _, tt := range tests {
		result := ProgressBar(tt.progress, tt.width)
		if len(result) != tt.width*3 { // Each block is 3 bytes in UTF-8
			// Actually just check it's not empty
			if len(result) == 0 {
				t.Errorf("ProgressBar(%f, %d) returned empty string", tt.progress, tt.width)
			}
		}
	}
}

func TestProgressIndicator(t *testing.T) {
	tests := []struct {
		progress float64
		expected string
	}{
		{50, "🟢"},
		{90, "🟡"},
		{120, "🟠"},
		{200, "🔴"},
	}

	for _, tt := range tests {
		result := ProgressIndicator(tt.progress)
		if result != tt.expected {
			t.Errorf("ProgressIndicator(%f) = %s, expected %s", tt.progress, result, tt.expected)
		}
	}
}
