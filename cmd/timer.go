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
	Long:  `Display the elapsed time, threshold, and progress percentage.`,
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
	fmt.Printf("  Threshold: %s\n", timer.FormatHours(status.ThresholdHours))

	if status.Overtime > 0 {
		fmt.Printf("  Overtime:  %s\n", timer.FormatDuration(status.Overtime))
	}

	if status.Extensions > 0 {
		fmt.Println()
		fmt.Printf("  Extensions: %d (added %s total)\n",
			status.Extensions, formatExtensions(s.Timer.Extensions))
	}

	fmt.Println()
	fmt.Println("  Commands:")
	fmt.Println("    yo extend 1h \"reason\"  - Add time")
	fmt.Println("    yo done                - Complete task")
	fmt.Println()
}

func formatExtensions(extensions []state.Extension) string {
	var total float64
	for _, e := range extensions {
		total += e.Hours
	}
	return timer.FormatHours(total)
}

func init() {
	rootCmd.AddCommand(timerCmd)
}
