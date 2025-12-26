package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/faisalahmedsifat/yo/internal/state"
	"github.com/faisalahmedsifat/yo/internal/timer"
	"github.com/faisalahmedsifat/yo/internal/workspace"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current yo status",
	Long:  `Display the current stage, task, timer, and session information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !workspace.IsInitialized() {
			return fmt.Errorf("workspace not initialized. Run 'yo init' first")
		}

		s, err := state.Load()
		if err != nil {
			return err
		}

		printStatus(s)
		return nil
	},
}

func printStatus(s *state.State) {
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("  yo status\n")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// Stage with emoji
	stageEmoji := map[string]string{
		"none":   "⚪",
		"red":    "🔴",
		"yellow": "🟡",
		"green":  "🟢",
	}
	emoji := stageEmoji[s.CurrentStage]
	if emoji == "" {
		emoji = "⚪"
	}
	fmt.Printf("  Stage: %s %s\n", emoji, strings.ToUpper(s.CurrentStage))

	// Task
	if s.CurrentTaskID != "" {
		fmt.Printf("  Task:  %s\n", s.CurrentTaskID)
		if s.CurrentTaskRepo != "" {
			fmt.Printf("  Repo:  %s\n", s.CurrentTaskRepo)
		}
	} else {
		fmt.Println("  Task:  (none)")
	}
	fmt.Println()

	// Timer - using timer package
	if s.CurrentStage == "green" && !s.Timer.StartedAt.IsZero() {
		status := timer.GetStatus(s)

		fmt.Println("  Timer:")
		fmt.Printf("    Elapsed:   %s\n", timer.FormatDuration(status.Elapsed))
		fmt.Printf("    Threshold: %s\n", timer.FormatHours(status.ThresholdHours))
		fmt.Printf("    Progress:  %.0f%%\n", status.Progress)

		if status.Extensions > 0 {
			fmt.Printf("    Extensions: %d\n", status.Extensions)
		}
	}

	// Session
	if s.Session.Active {
		sessionDuration := time.Since(s.Session.StartedAt)
		fmt.Println()
		fmt.Println("  Session:")
		fmt.Printf("    Started:  %s\n", s.Session.StartedAt.Format("15:04"))
		fmt.Printf("    Duration: %s\n", timer.FormatDuration(sessionDuration))
	}

	// Bypasses
	if s.EmergencyBypasses.Today > 0 || s.EmergencyBypasses.ThisWeek > 0 {
		fmt.Println()
		fmt.Println("  Emergency Bypasses:")
		fmt.Printf("    Today: %d/1, This Week: %d/5\n",
			s.EmergencyBypasses.Today, s.EmergencyBypasses.ThisWeek)
	}

	// Milestone
	fmt.Println()
	fmt.Printf("  Milestone: %d - %s\n", s.Milestone.Current, s.Milestone.Name)

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Next actions
	switch s.CurrentStage {
	case "none":
		fmt.Println("  Next: yo red  (define a problem)")
	case "red":
		fmt.Println("  Next: yo verify red  (validate problem definition)")
		fmt.Println("        yo yellow   (analyze and plan)")
	case "yellow":
		fmt.Println("  Next: yo verify yellow  (validate plan)")
		fmt.Println("        yo go            (start execution)")
	case "green":
		fmt.Println("  Next: yo done    (mark complete)")
		fmt.Println("        yo extend  (extend timer)")
	}
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
