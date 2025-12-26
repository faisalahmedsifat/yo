package cmd

import (
	"fmt"

	"github.com/faisal/yo/internal/task"
	"github.com/faisal/yo/internal/workspace"
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify [red|yellow]",
	Short: "Verify that a phase is complete",
	Long: `Verify validates that a phase (RED or YELLOW) is properly completed.

Examples:
  yo verify red     - Validate RED LIGHT is complete
  yo verify yellow  - Validate YELLOW LIGHT is complete`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !workspace.IsInitialized() {
			return fmt.Errorf("workspace not initialized. Run 'yo init' first")
		}

		phase := args[0]
		taskPath, err := workspace.GetCurrentTaskPath()
		if err != nil {
			return err
		}

		switch phase {
		case "red":
			return verifyRed(taskPath)
		case "yellow":
			return verifyYellow(taskPath)
		default:
			return fmt.Errorf("unknown phase: %s (use 'red' or 'yellow')", phase)
		}
	},
}

func verifyRed(taskPath string) error {
	result, err := task.ValidateRed(taskPath)
	if err != nil {
		return err
	}

	if result.Valid {
		fmt.Println("✅ RED LIGHT is complete!")
		fmt.Println("   Next: yo yellow  (analyze and plan)")
		return nil
	}

	fmt.Println("❌ RED LIGHT is incomplete:")
	for _, e := range result.Errors {
		fmt.Printf("   - %s: %s\n", e.Field, e.Message)
	}
	fmt.Println()
	fmt.Println("   Fix these issues with: yo red")
	return fmt.Errorf("validation failed")
}

func verifyYellow(taskPath string) error {
	result, err := task.ValidateYellow(taskPath)
	if err != nil {
		return err
	}

	if result.Valid {
		fmt.Println("✅ YELLOW LIGHT is complete!")
		fmt.Println("   Next: yo go  (start execution)")
		return nil
	}

	fmt.Println("❌ YELLOW LIGHT is incomplete:")
	for _, e := range result.Errors {
		fmt.Printf("   - %s: %s\n", e.Field, e.Message)
	}
	fmt.Println()
	fmt.Println("   Fix these issues with: yo yellow")
	return fmt.Errorf("validation failed")
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
