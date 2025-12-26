package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/faisalahmedsifat/yo/internal/activity"
	"github.com/faisalahmedsifat/yo/internal/state"
	"github.com/faisalahmedsifat/yo/internal/task"
	"github.com/faisalahmedsifat/yo/internal/timer"
	"github.com/faisalahmedsifat/yo/internal/workspace"
	"github.com/spf13/cobra"
)

var (
	goTimeOverride string
)

var goCmd = &cobra.Command{
	Use:   "go",
	Short: "Start GREEN LIGHT - execute with timer",
	Long: `GREEN LIGHT is the execution phase.
This command:
  - Verifies RED and YELLOW are complete
  - Starts a count-up timer
  - Sets the threshold from your time estimate
  - Begins tracking your work

Use --time to override the estimated time.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !workspace.IsInitialized() {
			return fmt.Errorf("workspace not initialized. Run 'yo init' first")
		}

		s, err := state.Load()
		if err != nil {
			return err
		}

		// Check stage
		if s.CurrentStage == "green" {
			elapsed := s.GetElapsed()
			fmt.Printf("⚠️  Already in GREEN LIGHT (timer: %s)\n", timer.FormatDuration(elapsed))
			fmt.Println("   Use 'yo done' to complete or 'yo extend' to add time.")
			return nil
		}

		if s.CurrentStage == "none" || s.CurrentStage == "" {
			return fmt.Errorf("complete RED and YELLOW first. Run 'yo red'")
		}

		taskPath, err := workspace.GetCurrentTaskPath()
		if err != nil {
			return err
		}

		// Validate RED
		redResult, err := task.ValidateRed(taskPath)
		if err != nil {
			return err
		}
		if !redResult.Valid {
			fmt.Println("❌ RED LIGHT incomplete:")
			for _, e := range redResult.Errors {
				fmt.Printf("   - %s\n", e.Message)
			}
			return fmt.Errorf("complete RED LIGHT first")
		}

		// Validate YELLOW
		yellowResult, err := task.ValidateYellow(taskPath)
		if err != nil {
			return err
		}
		if !yellowResult.Valid {
			fmt.Println("❌ YELLOW LIGHT incomplete:")
			for _, e := range yellowResult.Errors {
				fmt.Printf("   - %s\n", e.Message)
			}
			return fmt.Errorf("complete YELLOW LIGHT first")
		}

		// Get time estimate using timer package
		var estimatedHours float64
		if goTimeOverride != "" {
			estimatedHours, err = timer.ParseDuration(goTimeOverride)
			if err != nil || estimatedHours == 0 {
				return fmt.Errorf("invalid time format: %s (use format like 2h, 1.5h, 30m)", goTimeOverride)
			}
		} else {
			estimatedHours, err = task.GetTimeEstimate(taskPath)
			if err != nil {
				fmt.Println("⚠️  No time estimate found in task.")
				fmt.Print("   Enter time estimate (e.g., 2h, 1.5h): ")
				var timeStr string
				fmt.Scanln(&timeStr)
				estimatedHours, err = timer.ParseDuration(timeStr)
				if err != nil || estimatedHours == 0 {
					return fmt.Errorf("invalid time format")
				}
			}
		}

		// Update task file with green light info
		content, err := os.ReadFile(taskPath)
		if err != nil {
			return err
		}
		text := string(content)
		now := time.Now()
		text = strings.Replace(text, "### Timer Started:", fmt.Sprintf("### Timer Started: %s", now.Format("2006-01-02 15:04")), 1)
		text = strings.Replace(text, "### Estimated Time:", fmt.Sprintf("### Estimated Time: %.1fh", estimatedHours), 1)
		if err := os.WriteFile(taskPath, []byte(text), 0644); err != nil {
			return err
		}

		// Update state
		oldStage := s.CurrentStage
		s.SetStage("green")
		s.StartTimer(estimatedHours)
		s.StartSession()

		// Get current repo
		cwd, _ := os.Getwd()
		s.CurrentTaskRepo = cwd

		if err := s.Save(); err != nil {
			return err
		}

		// Log stage change
		activity.LogStageChange(oldStage, "green", s.CurrentTaskID)

		fmt.Println()
		fmt.Println("🟢 GREEN LIGHT - Execution Started!")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("   Task:      %s\n", s.CurrentTaskID)
		fmt.Printf("   Threshold: %s\n", timer.FormatHours(estimatedHours))
		fmt.Printf("   Started:   %s\n", now.Format("15:04"))
		fmt.Println()
		fmt.Println("   Timer is running! You've got this. 💪")
		fmt.Println()
		fmt.Println("   Commands:")
		fmt.Println("     yo status  - Check progress")
		fmt.Println("     yo timer   - Show timer only")
		fmt.Println("     yo extend  - Add more time")
		fmt.Println("     yo done    - Mark complete")
		return nil
	},
}

func init() {
	goCmd.Flags().StringVar(&goTimeOverride, "time", "", "Override time estimate (e.g., 2h, 1.5h, 30m)")
	rootCmd.AddCommand(goCmd)
}
