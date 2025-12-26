package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/faisal/yo/internal/activity"
	"github.com/faisal/yo/internal/state"
	"github.com/faisal/yo/internal/task"
	"github.com/faisal/yo/internal/templates"
	"github.com/faisal/yo/internal/timer"
	"github.com/faisal/yo/internal/workspace"
	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done",
	Short: "Complete the current task",
	Long: `Mark the current task as complete.
This will:
  - Verify success criteria
  - Stop the timer
  - Calculate accuracy (actual vs estimated)
  - Archive the task to done/
  - Clear current_task.md`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !workspace.IsInitialized() {
			return fmt.Errorf("workspace not initialized. Run 'yo init' first")
		}

		s, err := state.Load()
		if err != nil {
			return err
		}

		if s.CurrentStage != "green" {
			return fmt.Errorf("can only complete tasks in GREEN LIGHT. Current stage: %s", s.CurrentStage)
		}

		taskPath, err := workspace.GetCurrentTaskPath()
		if err != nil {
			return err
		}

		// Get success criteria
		criteria, err := task.GetSuccessCriteria(taskPath)
		if err != nil {
			fmt.Println("⚠️  Could not load success criteria")
			criteria = []string{}
		}

		// Verify success criteria
		if len(criteria) > 0 {
			fmt.Println("📋 Verify Success Criteria")
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Println()

			reader := bufio.NewReader(os.Stdin)
			allMet := true

			for _, c := range criteria {
				fmt.Printf("  [?] %s (y/n): ", c)
				response, _ := reader.ReadString('\n')
				response = strings.TrimSpace(strings.ToLower(response))
				if response != "y" && response != "yes" {
					allMet = false
				}
			}

			if !allMet {
				fmt.Println()
				fmt.Println("❌ Not all success criteria met.")
				fmt.Print("   Continue working? (y/n): ")
				response, _ := reader.ReadString('\n')
				if strings.TrimSpace(strings.ToLower(response)) != "n" {
					fmt.Println("   Keep working! You've got this. 💪")
					return nil
				}
			}
		}

		// Calculate time
		elapsed := s.GetElapsed()
		actualHours := elapsed.Hours()
		estimatedHours := s.Timer.EstimatedHours
		accuracy := (estimatedHours / actualHours) * 100
		if actualHours == 0 {
			accuracy = 100
		}

		// Archive task
		if err := archiveTask(s, taskPath); err != nil {
			fmt.Printf("⚠️  Failed to archive task: %v\n", err)
		}

		// Log completion
		activity.LogTaskComplete(s.CurrentTaskID, actualHours, estimatedHours)

		// Reset state
		taskID := s.CurrentTaskID
		s.SetStage("none")
		s.StopTimer()
		s.CurrentTaskID = ""
		s.CurrentTaskRepo = ""
		if err := s.Save(); err != nil {
			return err
		}

		// Reset current_task.md using templates package
		resetCurrentTask()

		fmt.Println()
		fmt.Println("✅ Task Complete!")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("  Task:      %s\n", taskID)
		fmt.Printf("  Actual:    %s\n", timer.FormatDuration(elapsed))
		fmt.Printf("  Estimated: %s\n", timer.FormatHours(estimatedHours))
		fmt.Printf("  Accuracy:  %.0f%%\n", accuracy)
		fmt.Println()

		if accuracy >= 80 && accuracy <= 120 {
			fmt.Println("  🎯 Great estimation!")
		} else if actualHours < estimatedHours {
			fmt.Println("  ⚡ Faster than expected!")
		} else {
			fmt.Println("  📝 Consider adding buffer to future estimates")
		}

		fmt.Println()
		fmt.Println("  Next: yo next  (pick another task)")
		fmt.Println("        yo off   (end session)")
		return nil
	},
}

func archiveTask(s *state.State, taskPath string) error {
	yoDir, err := state.GetYoDir()
	if err != nil {
		return err
	}

	// Create archive filename
	date := time.Now().Format("2006-01-02")
	archiveName := fmt.Sprintf("%s_%s.md", date, s.CurrentTaskID)
	archivePath := filepath.Join(yoDir, "done", archiveName)

	// Copy file
	content, err := os.ReadFile(taskPath)
	if err != nil {
		return err
	}

	// Add completion metadata
	metadata := fmt.Sprintf("\n---\n\n## Completion Metadata\n- Completed: %s\n- Actual Time: %s\n- Estimated Time: %s\n",
		time.Now().Format("2006-01-02 15:04"),
		timer.FormatDuration(s.GetElapsed()),
		timer.FormatHours(s.Timer.EstimatedHours))

	content = append(content, []byte(metadata)...)

	return os.WriteFile(archivePath, content, 0644)
}

func resetCurrentTask() {
	taskPath, err := workspace.GetCurrentTaskPath()
	if err != nil {
		return
	}

	// Use templates package instead of embedded template
	os.WriteFile(taskPath, []byte(templates.CurrentTask), 0644)
}

func init() {
	rootCmd.AddCommand(doneCmd)
}
