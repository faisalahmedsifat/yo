package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/faisal/yo/internal/activity"
	"github.com/faisal/yo/internal/state"
	"github.com/faisal/yo/internal/workspace"
	"github.com/spf13/cobra"
)

var (
	yellowEdit bool
)

var yellowCmd = &cobra.Command{
	Use:   "yellow",
	Short: "Start YELLOW LIGHT - analyze and plan",
	Long: `YELLOW LIGHT is the analysis and planning phase.
Before executing, you must:
  - Analyze root cause (3 levels deep)
  - List at least 3 solution options
  - Choose one with rationale
  - Define success criteria

By default, runs in interactive mode with guided prompts.
Use -e to open the editor directly.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !workspace.IsInitialized() {
			return fmt.Errorf("workspace not initialized. Run 'yo init' first")
		}

		s, err := state.Load()
		if err != nil {
			return err
		}

		// Check stage
		if s.CurrentStage == "none" || s.CurrentStage == "" {
			return fmt.Errorf("complete RED LIGHT first. Run 'yo red'")
		}

		if s.CurrentStage == "green" {
			return fmt.Errorf("already in GREEN LIGHT. Complete current task first with 'yo done'")
		}

		taskPath, err := workspace.GetCurrentTaskPath()
		if err != nil {
			return err
		}

		if yellowEdit {
			// Open editor
			fmt.Println("🟡 Opening current_task.md for YELLOW LIGHT...")
			fmt.Println("   Fill in the analysis and planning section.")
			fmt.Println()

			if err := openEditor(taskPath); err != nil {
				return fmt.Errorf("failed to open editor: %w", err)
			}

			// Update state
			oldStage := s.CurrentStage
			s.SetStage("yellow")
			if err := s.Save(); err != nil {
				return err
			}

			// Log stage change
			activity.LogStageChange(oldStage, "yellow", s.CurrentTaskID)

			fmt.Println()
			fmt.Println("✅ YELLOW LIGHT started!")
			fmt.Println("   Next: yo go  (start execution)")
			return nil
		}

		// Default: Interactive mode
		return runYellowInteractive(taskPath, s)
	},
}

func runYellowInteractive(taskPath string, s *state.State) error {
	fmt.Println("🟡 YELLOW LIGHT - Interactive Mode")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	// Root cause analysis
	fmt.Println("Root Cause Analysis (dig 3 levels deep)")
	fmt.Println()

	fmt.Print("Immediate cause (what directly causes this):\n> ")
	immediateCause, _ := reader.ReadString('\n')
	immediateCause = strings.TrimSpace(immediateCause)

	fmt.Print("\nUnderlying cause (why does the immediate cause exist):\n> ")
	underlyingCause, _ := reader.ReadString('\n')
	underlyingCause = strings.TrimSpace(underlyingCause)

	fmt.Print("\nSystem cause (what systemic issue allows this):\n> ")
	systemCause, _ := reader.ReadString('\n')
	systemCause = strings.TrimSpace(systemCause)

	// Solution options (minimum 3)
	fmt.Println()
	fmt.Println("Solution Options (minimum 3)")
	fmt.Println()

	type option struct {
		desc string
		time string
		pros string
		cons string
	}

	options := make([]option, 3)
	labels := []string{"A", "B", "C"}

	for i := 0; i < 3; i++ {
		fmt.Printf("Option %s:\n", labels[i])

		fmt.Print("  Description: ")
		options[i].desc, _ = reader.ReadString('\n')
		options[i].desc = strings.TrimSpace(options[i].desc)

		fmt.Print("  Time estimate (e.g., 2h, 4h): ")
		options[i].time, _ = reader.ReadString('\n')
		options[i].time = strings.TrimSpace(options[i].time)

		fmt.Print("  Pros: ")
		options[i].pros, _ = reader.ReadString('\n')
		options[i].pros = strings.TrimSpace(options[i].pros)

		fmt.Print("  Cons: ")
		options[i].cons, _ = reader.ReadString('\n')
		options[i].cons = strings.TrimSpace(options[i].cons)

		fmt.Println()
	}

	// Decision
	fmt.Print("Which option do you choose? (A/B/C): ")
	chosenStr, _ := reader.ReadString('\n')
	chosenStr = strings.TrimSpace(strings.ToUpper(chosenStr))
	chosenIdx := 0
	switch chosenStr {
	case "B":
		chosenIdx = 1
	case "C":
		chosenIdx = 2
	}

	fmt.Print("\nWhy this option? ")
	reason, _ := reader.ReadString('\n')
	reason = strings.TrimSpace(reason)

	// Implementation steps
	fmt.Println()
	fmt.Println("Implementation Steps (enter each step, empty line to finish):")
	var steps []string
	for i := 1; ; i++ {
		fmt.Printf("  %d. ", i)
		step, _ := reader.ReadString('\n')
		step = strings.TrimSpace(step)
		if step == "" {
			break
		}
		steps = append(steps, step)
	}

	// Success criteria
	fmt.Println()
	fmt.Println("Success Criteria (enter each criterion, empty line to finish):")
	var criteria []string
	for i := 1; ; i++ {
		fmt.Printf("  %d. ", i)
		criterion, _ := reader.ReadString('\n')
		criterion = strings.TrimSpace(criterion)
		if criterion == "" {
			break
		}
		criteria = append(criteria, criterion)
	}

	if len(criteria) < 2 {
		return fmt.Errorf("need at least 2 success criteria")
	}

	// Read and update template
	content, err := os.ReadFile(taskPath)
	if err != nil {
		return err
	}

	text := string(content)

	// Replace root cause placeholders
	text = strings.Replace(text, "<!-- What directly causes this? -->", immediateCause, 1)
	text = strings.Replace(text, "<!-- Why does the immediate cause exist? -->", underlyingCause, 1)
	text = strings.Replace(text, "<!-- What systemic issue allows this? -->", systemCause, 1)

	// Update options
	for i, opt := range options {
		label := labels[i]
		optSection := fmt.Sprintf(`#### Option %s:
- Description: %s
- Time estimate: %s
- Pros: %s
- Cons: %s`, label, opt.desc, opt.time, opt.pros, opt.cons)

		placeholder := fmt.Sprintf(`#### Option %s:
- Description:
- Time estimate:
- Pros:
- Cons:`, label)

		text = strings.Replace(text, placeholder, optSection, 1)
	}

	// Update decision
	text = strings.Replace(text, "**Chosen option:** \n**Reason:**",
		fmt.Sprintf("**Chosen option:** %s\n**Reason:** %s", chosenStr, reason), 1)

	// Update implementation steps
	var stepLines strings.Builder
	for i, step := range steps {
		stepLines.WriteString(fmt.Sprintf("%d. %s\n", i+1, step))
	}
	text = strings.Replace(text, "### Implementation Steps\n1. \n2. \n3. ",
		"### Implementation Steps\n"+stepLines.String(), 1)

	// Update success criteria
	var criteriaLines strings.Builder
	for _, c := range criteria {
		criteriaLines.WriteString(fmt.Sprintf("- [ ] %s\n", c))
	}
	text = strings.Replace(text, "### Success Criteria\n- [ ] \n- [ ] \n- [ ] ",
		"### Success Criteria\n"+criteriaLines.String(), 1)

	// Write updated content
	if err := os.WriteFile(taskPath, []byte(text), 0644); err != nil {
		return err
	}

	// Update state
	oldStage := s.CurrentStage
	s.SetStage("yellow")
	if err := s.Save(); err != nil {
		return err
	}

	// Log stage change
	activity.LogStageChange(oldStage, "yellow", s.CurrentTaskID)

	fmt.Println()
	fmt.Println("✅ YELLOW LIGHT complete!")
	fmt.Printf("   Time estimate: %s\n", options[chosenIdx].time)
	fmt.Println("   Next: yo go  (start execution with timer)")
	return nil
}

func init() {
	yellowCmd.Flags().BoolVarP(&yellowEdit, "edit", "e", false, "Open in editor instead of interactive mode")
	rootCmd.AddCommand(yellowCmd)
}
