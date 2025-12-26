package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/faisal/yo/internal/activity"
	"github.com/faisal/yo/internal/config"
	"github.com/faisal/yo/internal/state"
	"github.com/faisal/yo/internal/workspace"
	"github.com/spf13/cobra"
)

var bypassCmd = &cobra.Command{
	Use:   "bypass [reason]",
	Short: "Emergency bypass - skip the framework",
	Long: `Skip the RED/YELLOW/GREEN framework for emergencies.
Tracked and limited to:
  - 1 per day
  - 5 per week

Use sparingly for genuine emergencies only.`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !workspace.IsInitialized() {
			return fmt.Errorf("workspace not initialized. Run 'yo init' first")
		}

		reason := strings.Join(args, " ")

		s, err := state.Load()
		if err != nil {
			return err
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		// Check limits
		if s.EmergencyBypasses.Today >= cfg.MaxBypassDay {
			fmt.Printf("⚠️  Already used %d bypass today (max %d/day)\n",
				s.EmergencyBypasses.Today, cfg.MaxBypassDay)
		}

		if s.EmergencyBypasses.ThisWeek >= cfg.MaxBypassWeek {
			fmt.Printf("⚠️  Already used %d bypasses this week (max %d/week)\n",
				s.EmergencyBypasses.ThisWeek, cfg.MaxBypassWeek)
		}

		if s.EmergencyBypasses.Today >= cfg.MaxBypassDay ||
			s.EmergencyBypasses.ThisWeek >= cfg.MaxBypassWeek {
			fmt.Println()
			fmt.Print("Proceed anyway? (y/n): ")
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			if strings.TrimSpace(strings.ToLower(response)) != "y" {
				fmt.Println("Bypass cancelled.")
				return nil
			}
		}

		// Increment counters
		s.EmergencyBypasses.Today++
		s.EmergencyBypasses.ThisWeek++
		if err := s.Save(); err != nil {
			return err
		}

		// Log bypass
		activity.LogEmergencyBypass(reason, s.EmergencyBypasses.Today, s.EmergencyBypasses.ThisWeek)

		fmt.Println()
		fmt.Println("🚨 BYPASS Active!")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("   Reason: %s\n", reason)
		fmt.Printf("   Time limit: 30 minutes\n")
		fmt.Printf("   Bypasses today: %d/%d\n", s.EmergencyBypasses.Today, cfg.MaxBypassDay)
		fmt.Printf("   Bypasses this week: %d/%d\n", s.EmergencyBypasses.ThisWeek, cfg.MaxBypassWeek)
		fmt.Println()
		fmt.Println("   ⏰ Remember to document what you fixed!")

		// Start a 30-minute timer (in background)
		go func() {
			time.Sleep(30 * time.Minute)
			// In a full implementation, this would trigger a notification
		}()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(bypassCmd)
}
