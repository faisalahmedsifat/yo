package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/faisalahmedsifat/yo/internal/milestone"
	"github.com/faisalahmedsifat/yo/internal/state"
	"github.com/faisalahmedsifat/yo/internal/workspace"
	"github.com/spf13/cobra"
)

var milestoneCmd = &cobra.Command{
	Use:   "milestone",
	Short: "Track milestone progress",
	Long:  `View and manage progress through framework milestones.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return showMilestoneStatus()
	},
}

var milestoneStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current milestone status",
	RunE: func(cmd *cobra.Command, args []string) error {
		return showMilestoneStatus()
	},
}

var milestoneCompleteCmd = &cobra.Command{
	Use:   "complete",
	Short: "Mark current milestone as complete",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !workspace.IsInitialized() {
			return fmt.Errorf("workspace not initialized. Run 'yo init' first")
		}

		s, err := state.Load()
		if err != nil {
			return err
		}

		if milestone.IsComplete(s) {
			fmt.Println("🎉 All milestones complete!")
			return nil
		}

		current := milestone.Current(s)
		if current == nil {
			return fmt.Errorf("no current milestone")
		}

		reader := bufio.NewReader(os.Stdin)

		fmt.Println()
		fmt.Printf("📍 Completing Milestone %d: %s\n", current.ID, current.Name)
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()
		fmt.Println("Verify each criterion:")
		fmt.Println()

		allMet := true
		for _, c := range current.Criteria {
			fmt.Printf("  [?] %s (y/n): ", c)
			response, _ := reader.ReadString('\n')
			if strings.TrimSpace(strings.ToLower(response)) != "y" {
				allMet = false
			}
		}

		if !allMet {
			fmt.Println()
			fmt.Println("❌ Not all criteria met. Keep working!")
			return nil
		}

		// Advance milestone using package
		if milestone.Advance(s) {
			if err := s.Save(); err != nil {
				return err
			}
		}

		fmt.Println()
		fmt.Println("🎉 Milestone complete!")
		if !milestone.IsComplete(s) {
			next := milestone.Current(s)
			fmt.Printf("   Next: Milestone %d - %s\n", next.ID, next.Name)
		} else {
			fmt.Println("   All milestones complete! 🏆")
		}

		return nil
	},
}

func showMilestoneStatus() error {
	if !workspace.IsInitialized() {
		return fmt.Errorf("workspace not initialized. Run 'yo init' first")
	}

	s, err := state.Load()
	if err != nil {
		return err
	}

	status := milestone.GetStatus(s)
	milestones := milestone.All()

	fmt.Println()
	fmt.Println("📍 Milestone Progress")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	for i, m := range milestones {
		var icon string
		if i < status.Current {
			icon = "✅"
		} else if i == status.Current {
			icon = "🔄"
		} else {
			icon = "⬜"
		}
		fmt.Printf("  %s Milestone %d: %s\n", icon, i, m.Name)
	}

	fmt.Println()

	if !milestone.IsComplete(s) {
		current := milestone.Current(s)
		fmt.Printf("  Current: Milestone %d - %s\n", current.ID, current.Name)
		fmt.Println()
		fmt.Println("  Exit criteria:")
		for _, c := range current.Criteria {
			fmt.Printf("    [ ] %s\n", c)
		}
		fmt.Println()
		fmt.Println("  Complete with: yo milestone complete")
	} else {
		fmt.Println("  🏆 All milestones complete!")
	}

	return nil
}

func init() {
	milestoneCmd.AddCommand(milestoneStatusCmd)
	milestoneCmd.AddCommand(milestoneCompleteCmd)
	rootCmd.AddCommand(milestoneCmd)
}
