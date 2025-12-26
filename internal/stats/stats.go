package stats

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/faisal/yo/internal/activity"
	"github.com/faisal/yo/internal/state"
)

// WeekStats holds weekly statistics
type WeekStats struct {
	WeekStart      time.Time `json:"week_start"`
	WeekEnd        time.Time `json:"week_end"`
	TasksCompleted int       `json:"tasks_completed"`
	Bypasses       int       `json:"bypasses"`
	AvgAccuracy    float64   `json:"avg_accuracy"`
	FocusScore     float64   `json:"focus_score"`
	OnTaskChanges  int       `json:"on_task_changes"`
	TotalChanges   int       `json:"total_changes"`
	TotalHours     float64   `json:"total_hours"`
}

// Calculate generates stats from activity entries
func Calculate(entries []activity.Entry) *WeekStats {
	stats := &WeekStats{}

	var totalAccuracy float64
	accuracyCount := 0

	for _, e := range entries {
		switch e.Type {
		case activity.TypeTaskComplete:
			stats.TasksCompleted++
			if e.EstimatedHours > 0 && e.ActualHours > 0 {
				accuracy := (e.EstimatedHours / e.ActualHours) * 100
				totalAccuracy += accuracy
				accuracyCount++
			}
			stats.TotalHours += e.ActualHours

		case activity.TypeEmergencyBypass:
			stats.Bypasses++

		case activity.TypeFileChange:
			stats.TotalChanges++
			if !e.Untracked {
				stats.OnTaskChanges++
			}
		}
	}

	if accuracyCount > 0 {
		stats.AvgAccuracy = totalAccuracy / float64(accuracyCount)
	}

	if stats.TotalChanges > 0 {
		stats.FocusScore = (float64(stats.OnTaskChanges) / float64(stats.TotalChanges)) * 100
	} else {
		stats.FocusScore = 100
	}

	return stats
}

// GetWeekRange returns start and end of a week containing the given date
func GetWeekRange(date time.Time) (start, end time.Time) {
	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	start = time.Date(date.Year(), date.Month(), date.Day()-weekday+1, 0, 0, 0, 0, date.Location())
	end = start.Add(7 * 24 * time.Hour)
	return
}

// GetCurrentWeekRange returns the current week's range
func GetCurrentWeekRange() (start, end time.Time) {
	return GetWeekRange(time.Now())
}

// ForWeek calculates stats for a specific week
func ForWeek(weekStart time.Time) (*WeekStats, error) {
	start, end := GetWeekRange(weekStart)

	entries, err := activity.Query(start, end)
	if err != nil {
		return nil, err
	}

	stats := Calculate(entries)
	stats.WeekStart = start
	stats.WeekEnd = end

	return stats, nil
}

// Save caches stats to disk
func (s *WeekStats) Save() error {
	yoDir, err := state.GetYoDir()
	if err != nil {
		return err
	}

	year, week := s.WeekStart.ISOWeek()
	filename := fmt.Sprintf("%d-week-%02d.json", year, week)
	statsPath := filepath.Join(yoDir, "stats", filename)

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(statsPath, data, 0644)
}

// Load loads cached stats from disk
func Load(weekStart time.Time) (*WeekStats, error) {
	yoDir, err := state.GetYoDir()
	if err != nil {
		return nil, err
	}

	year, week := weekStart.ISOWeek()
	filename := fmt.Sprintf("%d-week-%02d.json", year, week)
	statsPath := filepath.Join(yoDir, "stats", filename)

	data, err := os.ReadFile(statsPath)
	if err != nil {
		return nil, err
	}

	var stats WeekStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

// Insights generates insights from stats
type Insights struct {
	Messages []string
}

// GenerateInsights analyzes stats and provides recommendations
func GenerateInsights(s *WeekStats) *Insights {
	insights := &Insights{}

	// Estimation accuracy insights
	if s.AvgAccuracy > 0 {
		if s.AvgAccuracy >= 80 && s.AvgAccuracy <= 120 {
			insights.Messages = append(insights.Messages, "🎯 Great estimation accuracy!")
		} else if s.AvgAccuracy < 80 {
			insights.Messages = append(insights.Messages, "📝 Tasks completing faster than estimated - tighten estimates")
		} else {
			insights.Messages = append(insights.Messages, "📝 Tasks taking longer - add 20% buffer to estimates")
		}
	}

	// Focus score insights
	if s.FocusScore < 60 {
		insights.Messages = append(insights.Messages, "⚠️ Low focus score - reduce context switching")
	} else if s.FocusScore >= 80 {
		insights.Messages = append(insights.Messages, "🌟 Excellent focus!")
	}

	// Bypass insights
	if s.Bypasses > 5 {
		insights.Messages = append(insights.Messages, "🚨 Too many emergency bypasses - improve planning")
	}

	// Productivity insights
	if s.TasksCompleted == 0 {
		insights.Messages = append(insights.Messages, "💡 No tasks completed - break work into smaller chunks")
	} else if s.TasksCompleted >= 5 {
		insights.Messages = append(insights.Messages, "🏆 Great productivity!")
	}

	return insights
}
