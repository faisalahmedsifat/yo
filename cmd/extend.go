package cmd

import (
	"fmt"
	"strings"

	"github.com/faisal/yo/internal/state"
	"github.com/faisal/yo/internal/timer"
	"github.com/faisal/yo/internal/workspace"
	"github.com/spf13/cobra"
)

var extendCmd = &cobra.Command{
	Use:   "extend <duration> [reason]",
	Short: "Extend the timer threshold",
	Long: `Add more time to the current threshold.

This does NOT create tech debt - it's just adjusting your estimate.
Use 'yo defer' to log conscious decisions to defer work.

Examples:
  yo extend 1h "API more complex than expected"
  yo extend 30m
  yo extend 2h`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !workspace.IsInitialized() {
			return fmt.Errorf("workspace not initialized. Run 'yo init' first")
		}

		s, err := state.Load()
		if err != nil {
			return err
		}

		if s.CurrentStage != "green" {
			return fmt.Errorf("can only extend timer in GREEN LIGHT. Current stage: %s", s.CurrentStage)
		}

		// Parse duration using timer package
		duration, err := timer.ParseDuration(args[0])
		if err != nil || duration == 0 {
			return fmt.Errorf("invalid duration: %s (use format like 1h, 30m, 1.5h)", args[0])
		}

		// Get reason
		reason := ""
		if len(args) > 1 {
			reason = strings.Join(args[1:], " ")
		}

		oldThreshold := s.Timer.ThresholdHours
		s.ExtendTimer(duration, reason)

		if err := s.Save(); err != nil {
			return err
		}

		fmt.Println()
		fmt.Println("⏱️  Timer Extended!")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("  Added:         +%s\n", timer.FormatHours(duration))
		fmt.Printf("  Old threshold: %s\n", timer.FormatHours(oldThreshold))
		fmt.Printf("  New threshold: %s\n", timer.FormatHours(s.Timer.ThresholdHours))
		if reason != "" {
			fmt.Printf("  Reason:        %s\n", reason)
		}
		fmt.Println()
		fmt.Printf("  Total extensions: %d\n", len(s.Timer.Extensions))
		fmt.Println()

		// Tip about tech debt
		fmt.Println("  💡 Tip: If you're skipping work to ship faster, log it:")
		fmt.Println("     yo defer -i")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(extendCmd)
}
