package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/faisalahmedsifat/yo/internal/state"
	"github.com/faisalahmedsifat/yo/internal/timer"
	"github.com/faisalahmedsifat/yo/internal/workspace"
	"github.com/spf13/cobra"
)

var timerWatch bool

var timerCmd = &cobra.Command{
	Use:   "timer",
	Short: "Show the current timer",
	Long: `Display the elapsed time, estimate, and progress percentage.

Use -w or --watch for a live updating timer.`,
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

		if timerWatch {
			return runLiveTimer(s)
		}

		printTimer(s)
		return nil
	},
}

func runLiveTimer(s *state.State) error {
	// Handle Ctrl+C gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// Clear screen and hide cursor
	fmt.Print("\033[?25l")       // Hide cursor
	defer fmt.Print("\033[?25h") // Show cursor on exit

	// Initial draw
	drawLiveTimer(s)

	for {
		select {
		case <-sigChan:
			fmt.Println("\n\n  Timer stopped. Task still running.")
			fmt.Println("  Use 'yo done' when complete.")
			return nil
		case <-ticker.C:
			drawLiveTimer(s)
		}
	}
}

func drawLiveTimer(s *state.State) {
	// Move cursor to top and clear
	fmt.Print("\033[H\033[2J")

	status := timer.GetStatus(s)

	fmt.Println()
	fmt.Println("⏱️  Live Timer (Ctrl+C to exit)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// Progress bar
	bar := timer.ProgressBar(status.Progress, 30)
	indicator := timer.ProgressIndicator(status.Progress)

	fmt.Printf("  %s [%s] %.1f%%\n", indicator, bar, status.Progress)
	fmt.Println()
	fmt.Printf("  Elapsed:   %s\n", timer.FormatDurationWithSeconds(status.Elapsed))
	fmt.Printf("  Estimate:  %s\n", timer.FormatHours(status.ThresholdHours))

	if status.Overtime > 0 {
		fmt.Println()
		fmt.Printf("  ⚠️  Over estimate by %s\n", timer.FormatDuration(status.Overtime))
	}

	fmt.Println()
	fmt.Printf("  Task: %s\n", s.CurrentTaskID)
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
	fmt.Println("    yo timer -w  - Live timer")
	fmt.Println("    yo done      - Complete task")
	fmt.Println()
}

func init() {
	timerCmd.Flags().BoolVarP(&timerWatch, "watch", "w", false, "Live updating timer")
	rootCmd.AddCommand(timerCmd)
}
