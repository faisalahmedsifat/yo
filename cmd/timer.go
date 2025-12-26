package cmd

import (
	"fmt"

	"github.com/faisal/yo/internal/state"
	"github.com/faisal/yo/internal/timer"
	"github.com/faisal/yo/internal/workspace"
	"github.com/spf13/cobra"
)

var timerCmd = &cobra.Command{
	Use:   "timer",
	Short: "Show the current timer",
	Long:  `Display the elapsed time, estimate, and progress percentage.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !workspace.IsInitialized() {
			return fmt.Errorf("workspace not initialized. Run 'yo init' first")
		}

		s, err := state.Load()
		if err != nil {
			return err
		}

		if s.CurrentStage != "green" {
			return fmt.Errorf("timer only runs in GREEN LIGHT. Current stage: %s", s.CurrentStage)
		}

		if s.Timer.StartedAt.IsZero() {
			return fmt.Errorf("timer not started. Run 'yo go' first")
		}

		printTimer(s)
		return nil
	},
}

func printTimer(s *state.State) {
	status := timer.GetStatus(s)

	fmt.Println()
	fmt.Println("⏱️  Timer")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// Progress bar using timer package
	bar := timer.ProgressBar(status.Progress, 30)
	indicator := timer.ProgressIndicator(status.Progress)

	fmt.Printf("  %s [%s] %.0f%%\n", indicator, bar, status.Progress)
	fmt.Println()
	fmt.Printf("  Elapsed:   %s\n", timer.FormatDuration(status.Elapsed))
	fmt.Printf("  Estimate:  %s\n", timer.FormatHours(status.ThresholdHours))

	if status.Overtime > 0 {
		fmt.Println()
		fmt.Printf("  ⚠️  Over estimate by %s\n", timer.FormatDuration(status.Overtime))
	}

	fmt.Println()
	fmt.Println("  Commands:")
	fmt.Println("    yo done   - Complete task")
	fmt.Println("    yo status - Full status")
	fmt.Println()
}

func init() {
	rootCmd.AddCommand(timerCmd)
}
