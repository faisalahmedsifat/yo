package cmd

import (
	"fmt"

	"github.com/faisal/yo/internal/stats"
	"github.com/faisal/yo/internal/timer"
	"github.com/faisal/yo/internal/workspace"
	"github.com/spf13/cobra"
)

var statsWeek string

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show weekly statistics",
	Long: `Display weekly statistics including:
  - Tasks completed
  - Time estimate accuracy
  - Emergency bypasses
  - Focus score

Examples:
  yo stats                    - Current week
  yo stats --week 2024-12-20  - Specific week`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !workspace.IsInitialized() {
			return fmt.Errorf("workspace not initialized. Run 'yo init' first")
		}

		var weekStart, weekEnd = stats.GetCurrentWeekRange()

		if statsWeek != "" {
			parsed, err := timer.ParseDuration(statsWeek)
			if err != nil {
				// Try parsing as date
				return fmt.Errorf("invalid date format: %s (use YYYY-MM-DD)", statsWeek)
			}
			_ = parsed // Unused in this case
		}

		weekStats, err := stats.ForWeek(weekStart)
		if err != nil {
			return err
		}

		// Display
		fmt.Println()
		fmt.Printf("📊 Week of %s - %s\n", weekStart.Format("Jan 2"), weekEnd.Format("Jan 2"))
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()

		fmt.Printf("  Tasks completed:    %d\n", weekStats.TasksCompleted)
		fmt.Printf("  Emergency bypasses: %d", weekStats.Bypasses)
		if weekStats.Bypasses > 5 {
			fmt.Print(" 🚨")
		}
		fmt.Println()
		fmt.Println()

		if weekStats.TasksCompleted > 0 {
			fmt.Println("  Time estimate accuracy:")
			fmt.Printf("    Average: %.0f%%\n", weekStats.AvgAccuracy)
			if weekStats.AvgAccuracy >= 80 && weekStats.AvgAccuracy <= 120 {
				fmt.Println("    🎯 Great estimations!")
			} else if weekStats.AvgAccuracy > 120 {
				fmt.Println("    📝 Tasks taking longer than estimated")
			} else {
				fmt.Println("    ⚡ Tasks completing faster than estimated")
			}
			fmt.Println()
		}

		fmt.Printf("  Focus score: %.0f%%\n", weekStats.FocusScore)
		if weekStats.FocusScore >= 80 {
			fmt.Println("  🌟 Excellent focus!")
		}
		fmt.Println()

		// Generate and display insights
		insights := stats.GenerateInsights(weekStats)
		if len(insights.Messages) > 0 {
			fmt.Println("  Insights:")
			for _, msg := range insights.Messages {
				fmt.Printf("    %s\n", msg)
			}
			fmt.Println()
		}

		// Cache stats
		weekStats.Save()

		return nil
	},
}

func init() {
	statsCmd.Flags().StringVar(&statsWeek, "week", "", "Week start date (YYYY-MM-DD)")
	rootCmd.AddCommand(statsCmd)
}
