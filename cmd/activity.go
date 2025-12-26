package cmd

import (
	"fmt"

	"github.com/faisalahmedsifat/yo/internal/activity"
	"github.com/faisalahmedsifat/yo/internal/timer"
	"github.com/faisalahmedsifat/yo/internal/workspace"
	"github.com/spf13/cobra"
)

var (
	activityYesterday bool
	activityWeek      bool
	activityRepo      string
)

var activityCmd = &cobra.Command{
	Use:   "activity",
	Short: "Show activity log",
	Long: `Display activity log entries.

Examples:
  yo activity             - Today's activity
  yo activity --yesterday - Yesterday's activity
  yo activity --week      - This week's activity`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !workspace.IsInitialized() {
			return fmt.Errorf("workspace not initialized. Run 'yo init' first")
		}

		var entries []activity.Entry
		var err error

		if activityYesterday {
			entries, err = activity.QueryYesterday()
		} else if activityWeek {
			entries, err = activity.QueryThisWeek()
		} else {
			entries, err = activity.QueryToday()
		}

		if err != nil {
			return err
		}

		// Filter by repo if specified
		if activityRepo != "" {
			var filtered []activity.Entry
			for _, e := range entries {
				if e.Repo == activityRepo {
					filtered = append(filtered, e)
				}
			}
			entries = filtered
		}

		// Display
		fmt.Println()
		fmt.Println("📊 Activity")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()

		if len(entries) == 0 {
			fmt.Println("  No activity recorded.")
			return nil
		}

		// Group by type
		stageChanges := 0
		fileChanges := 0
		bypasses := 0

		for _, e := range entries {
			switch e.Type {
			case activity.TypeStageChange:
				stageChanges++
			case activity.TypeFileChange:
				fileChanges++
			case activity.TypeEmergencyBypass:
				bypasses++
			}
		}

		fmt.Printf("  Stage changes: %d\n", stageChanges)
		fmt.Printf("  File changes:  %d\n", fileChanges)
		fmt.Printf("  Bypasses:      %d\n", bypasses)
		fmt.Println()

		// Show recent entries
		fmt.Println("  Recent events:")
		limit := 10
		if len(entries) < limit {
			limit = len(entries)
		}
		for i := len(entries) - 1; i >= len(entries)-limit && i >= 0; i-- {
			e := entries[i]
			ts := e.Timestamp.Format("15:04")
			switch e.Type {
			case activity.TypeStageChange:
				fmt.Printf("    %s  %s → %s\n", ts, e.From, e.To)
			case activity.TypeFileChange:
				fmt.Printf("    %s  📝 %s\n", ts, e.File)
			case activity.TypeEmergencyBypass:
				fmt.Printf("    %s  🚨 Bypass: %s\n", ts, e.Reason)
			case activity.TypeTaskComplete:
				fmt.Printf("    %s  ✅ Completed: %s\n", ts, e.Task)
			}
		}
		fmt.Println()

		return nil
	},
}

var focusCmd = &cobra.Command{
	Use:   "focus",
	Short: "Show focus score",
	Long:  `Calculate and display your focus score based on today's activity.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !workspace.IsInitialized() {
			return fmt.Errorf("workspace not initialized. Run 'yo init' first")
		}

		entries, err := activity.QueryToday()
		if err != nil {
			return err
		}

		// Calculate focus
		var onTask, untracked int
		repoTime := make(map[string]int)

		for _, e := range entries {
			if e.Type == activity.TypeFileChange {
				repoTime[e.Repo]++
				if e.Untracked {
					untracked++
				} else {
					onTask++
				}
			}
		}

		total := onTask + untracked
		var focusScore float64
		if total > 0 {
			focusScore = (float64(onTask) / float64(total)) * 100
		} else {
			focusScore = 100
		}

		fmt.Println()
		fmt.Println("🎯 Focus Score")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()
		fmt.Printf("  Score: %.0f%%\n", focusScore)
		fmt.Println()

		// Display progress bar using timer package
		bar := timer.ProgressBar(focusScore, 20)
		indicator := timer.ProgressIndicator(100 - focusScore) // Invert for focus
		fmt.Printf("  %s [%s]\n", indicator, bar)
		fmt.Println()

		if len(repoTime) > 0 {
			fmt.Println("  Time by repo:")
			for repo, count := range repoTime {
				fmt.Printf("    %s: %d changes\n", repo, count)
			}
		}
		fmt.Println()

		if focusScore >= 80 {
			fmt.Println("  🌟 Great focus today!")
		} else if focusScore >= 60 {
			fmt.Println("  👍 Good focus, but room for improvement")
		} else {
			fmt.Println("  ⚠️  Consider reducing context switching")
		}

		return nil
	},
}

func init() {
	activityCmd.Flags().BoolVar(&activityYesterday, "yesterday", false, "Show yesterday's activity")
	activityCmd.Flags().BoolVar(&activityWeek, "week", false, "Show this week's activity")
	activityCmd.Flags().StringVar(&activityRepo, "repo", "", "Filter by repo")

	rootCmd.AddCommand(activityCmd)
	rootCmd.AddCommand(focusCmd)
}
