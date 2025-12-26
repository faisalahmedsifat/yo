package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/faisalahmedsifat/yo/internal/activity"
	"github.com/faisalahmedsifat/yo/internal/state"
	"github.com/faisalahmedsifat/yo/internal/timer"
	"github.com/faisalahmedsifat/yo/internal/workspace"
	"github.com/spf13/cobra"
)

var offCmd = &cobra.Command{
	Use:   "off",
	Short: "End work session",
	Long: `End your work session.
This will:
  - Calculate session duration
  - Show activity summary
  - Calculate focus score
  - Save session to sessions/
  - Optionally pause or complete current task`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !workspace.IsInitialized() {
			return fmt.Errorf("workspace not initialized. Run 'yo init' first")
		}

		s, err := state.Load()
		if err != nil {
			return err
		}

		if !s.Session.Active {
			return fmt.Errorf("no active session. Start with 'yo go'")
		}

		sessionDuration := time.Since(s.Session.StartedAt)
		reader := bufio.NewReader(os.Stdin)

		fmt.Println()
		fmt.Println("📊 Session Summary")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()
		fmt.Printf("  Duration:  %s\n", timer.FormatDuration(sessionDuration))
		fmt.Printf("  Started:   %s\n", s.Session.StartedAt.Format("15:04"))
		fmt.Printf("  Ended:     %s\n", time.Now().Format("15:04"))
		fmt.Println()

		// Get activity summary
		entries, _ := activity.QueryToday()
		focusScore := calculateFocusScore(entries, s.CurrentTaskRepo)
		fmt.Printf("  Focus:     %.0f%%\n", focusScore)

		// Handle current task
		if s.CurrentStage == "green" {
			fmt.Println()
			fmt.Printf("  Current task: %s (IN PROGRESS)\n", s.CurrentTaskID)
			elapsed := s.GetElapsed()
			fmt.Printf("  Time on task: %s\n", timer.FormatDuration(elapsed))
			fmt.Println()

			fmt.Print("  Task is incomplete. Roll to tomorrow? (y/n): ")
			response, _ := reader.ReadString('\n')
			if strings.TrimSpace(strings.ToLower(response)) == "n" {
				fmt.Println("  Task left in progress. Continue tomorrow with 'yo status'")
			}
		}

		// Save session
		if err := saveSession(s, sessionDuration, focusScore); err != nil {
			fmt.Printf("⚠️  Failed to save session: %v\n", err)
		}

		// Log session end
		activity.LogSessionEnd(int(sessionDuration.Minutes()), s.CurrentTaskRepo, focusScore)

		// End session but keep task state
		s.EndSession()
		if err := s.Save(); err != nil {
			return err
		}

		fmt.Println()
		fmt.Println("👋 Session ended. Great work today!")
		fmt.Println()

		// Ask about backlog
		fmt.Print("  Any new issues for backlog? (y/n): ")
		response, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(response)) == "y" {
			fmt.Println("  Add them with: yo add \"issue description\"")
		}

		return nil
	},
}

func calculateFocusScore(entries []activity.Entry, currentRepo string) float64 {
	if len(entries) == 0 || currentRepo == "" {
		return 100 // Default to 100% if no data
	}

	var onTask, total int
	for _, e := range entries {
		if e.Type == activity.TypeFileChange {
			total++
			if e.Repo == currentRepo || !e.Untracked {
				onTask++
			}
		}
	}

	if total == 0 {
		return 100
	}

	return (float64(onTask) / float64(total)) * 100
}

// SessionData represents saved session info
type SessionData struct {
	Date          string  `json:"date"`
	StartedAt     string  `json:"started_at"`
	EndedAt       string  `json:"ended_at"`
	DurationMins  int     `json:"duration_minutes"`
	FocusScore    float64 `json:"focus_score"`
	TaskID        string  `json:"task_id,omitempty"`
	TaskCompleted bool    `json:"task_completed"`
}

func saveSession(s *state.State, duration time.Duration, focusScore float64) error {
	yoDir, err := state.GetYoDir()
	if err != nil {
		return err
	}

	sessionData := SessionData{
		Date:          time.Now().Format("2006-01-02"),
		StartedAt:     s.Session.StartedAt.Format(time.RFC3339),
		EndedAt:       time.Now().Format(time.RFC3339),
		DurationMins:  int(duration.Minutes()),
		FocusScore:    focusScore,
		TaskID:        s.CurrentTaskID,
		TaskCompleted: s.CurrentStage == "none",
	}

	data, err := json.MarshalIndent(sessionData, "", "  ")
	if err != nil {
		return err
	}

	// Create filename
	filename := fmt.Sprintf("%s_session_%d.json",
		time.Now().Format("2006-01-02"),
		time.Now().Unix()%100000)
	sessionPath := filepath.Join(yoDir, "sessions", filename)

	return os.WriteFile(sessionPath, data, 0644)
}

func init() {
	rootCmd.AddCommand(offCmd)
}
