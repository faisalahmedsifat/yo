package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/faisalahmedsifat/yo/internal/activity"
	"github.com/faisalahmedsifat/yo/internal/state"
	"github.com/faisalahmedsifat/yo/internal/task"
	"github.com/faisalahmedsifat/yo/internal/workspace"
	"github.com/spf13/cobra"
)

var (
	redEdit bool
)

var redCmd = &cobra.Command{
	Use:   "red",
	Short: "Start RED LIGHT - define the problem",
	Long: `RED LIGHT is the first phase of the framework.
Before writing any code, you must clearly define:
  - What is the problem?
  - What is the impact?
  - What is the severity (P0-P3)?

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

		taskPath, err := workspace.GetCurrentTaskPath()
		if err != nil {
			return err
		}

		// Check if already in RED stage (e.g., from yo next)
		if s.CurrentStage == "red" {
			fmt.Printf("⚠️  Already in RED stage for: %s\n", s.CurrentTaskID)
			fmt.Println()
			fmt.Println("  [1] Continue - Fill in missing impact/severity")
			fmt.Println("  [2] Start new - Begin fresh RED LIGHT")
			fmt.Println("  [q] Cancel")
			fmt.Print("> ")

			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(response)

			switch response {
			case "1":
				// Get existing problem from task file
				content, _ := os.ReadFile(taskPath)
				problem := extractProblem(string(content))
				if problem == "" {
					problem = s.CurrentTaskID
				}
				return runRedContinue(taskPath, s, problem)
			case "2":
				// Start fresh
				return runRedInteractive(taskPath, s)
			default:
				return nil
			}
		}

		// Check if in other stages
		if s.CurrentStage != "none" && s.CurrentStage != "" {
			fmt.Printf("⚠️  Already in %s stage.\n", strings.ToUpper(s.CurrentStage))
			if !promptConfirm("Start a new RED LIGHT anyway?") {
				return nil
			}
		}

		if redEdit {
			// Open editor
			fmt.Println("🔴 Opening current_task.md for RED LIGHT...")
			fmt.Println("   Fill in the problem definition, then save and close.")
			fmt.Println()

			if err := openEditor(taskPath); err != nil {
				return fmt.Errorf("failed to open editor: %w", err)
			}

			// Update state
			oldStage := s.CurrentStage
			s.SetStage("red")
			if err := s.Save(); err != nil {
				return err
			}

			// Log stage change
			activity.LogStageChange(oldStage, "red", s.CurrentTaskID)

			fmt.Println()
			fmt.Println("✅ RED LIGHT started!")
			fmt.Println("   Next: yo yellow  (analyze and plan)")
			return nil
		}

		// Default: Interactive mode
		return runRedInteractive(taskPath, s)
	},
}

func runRedInteractive(taskPath string, s *state.State) error {
	fmt.Println("🔴 RED LIGHT - Interactive Mode")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()

	// Get problem description
	fmt.Print("What's the problem?\n> ")
	problem := readLine()
	if problem == "" {
		return fmt.Errorf("problem description cannot be empty")
	}

	// Get impact (multi-select simulation)
	fmt.Println()
	fmt.Println("What's the impact? (enter numbers, comma-separated)")
	fmt.Println("  1. Blocks launch")
	fmt.Println("  2. Blocks paying users")
	fmt.Println("  3. Causes user frustration")
	fmt.Println("  4. Tech debt accumulation")
	fmt.Println("  5. Other")
	fmt.Print("> ")
	impactStr := readLine()
	impacts := parseNumbers(impactStr)
	if len(impacts) == 0 {
		return fmt.Errorf("at least one impact must be selected")
	}

	// Get severity
	fmt.Println()
	fmt.Println("Severity?")
	fmt.Println("  0. P0 - Launch blocker")
	fmt.Println("  1. P1 - Paying user blocker")
	fmt.Println("  2. P2 - Nice to have")
	fmt.Println("  3. P3 - Future improvement")
	fmt.Print("> ")
	severityStr := readLine()
	severity := 0
	fmt.Sscanf(severityStr, "%d", &severity)
	if severity < 0 || severity > 3 {
		return fmt.Errorf("invalid severity (0-3)")
	}

	// Read current template
	content, err := os.ReadFile(taskPath)
	if err != nil {
		return err
	}

	// Update template with answers
	text := string(content)

	// Replace problem placeholder
	text = strings.Replace(text, "<!-- Describe the problem clearly -->", problem, 1)

	// Check impact boxes
	impactMap := map[int]string{
		1: "Blocks launch",
		2: "Blocks paying users",
		3: "Causes user frustration",
		4: "Tech debt accumulation",
		5: "Other",
	}
	for _, i := range impacts {
		if label, ok := impactMap[i]; ok {
			text = strings.Replace(text, "- [ ] "+label, "- [x] "+label, 1)
		}
	}

	// Check severity box
	severityLabels := []string{
		"P0 - Launch blocker",
		"P1 - Paying user blocker",
		"P2 - Nice to have",
		"P3 - Future improvement",
	}
	text = strings.Replace(text, "- [ ] "+severityLabels[severity], "- [x] "+severityLabels[severity], 1)

	// Write updated content
	if err := os.WriteFile(taskPath, []byte(text), 0644); err != nil {
		return err
	}

	// Update state
	oldStage := s.CurrentStage
	s.SetStage("red")
	taskID := sanitizeTaskID(problem)
	s.CurrentTaskID = taskID
	if err := s.Save(); err != nil {
		return err
	}

	// Log stage change
	activity.LogStageChange(oldStage, "red", taskID)

	fmt.Println()
	fmt.Println("✅ RED LIGHT complete!")
	fmt.Printf("   Task ID: %s\n", taskID)
	fmt.Println("   Next: yo yellow  (analyze and plan)")
	return nil
}

func openEditor(filepath string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = "nano"
	}

	cmd := exec.Command(editor, filepath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func parseNumbers(s string) []int {
	var nums []int
	parts := strings.Split(s, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		var n int
		if _, err := fmt.Sscanf(p, "%d", &n); err == nil {
			nums = append(nums, n)
		}
	}
	return nums
}

func sanitizeTaskID(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "_")
	// Keep only alphanumeric and underscores
	var result strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			result.WriteRune(r)
		}
	}
	id := result.String()
	if len(id) > 30 {
		id = id[:30]
	}
	return id
}

func promptConfirm(question string) bool {
	return task.PromptConfirm(question)
}

// extractProblem extracts the problem description from task file
func extractProblem(text string) string {
	if idx := strings.Index(text, "### What's the Problem?"); idx != -1 {
		endIdx := strings.Index(text[idx:], "### Impact")
		if endIdx != -1 {
			problem := strings.TrimSpace(text[idx+len("### What's the Problem?") : idx+endIdx])
			problem = strings.TrimPrefix(problem, "<!-- Describe the problem clearly -->")
			return strings.TrimSpace(problem)
		}
	}
	return ""
}

// runRedContinue fills in impact/severity for an existing problem (from yo next)
func runRedContinue(taskPath string, s *state.State, problem string) error {
	fmt.Println()
	fmt.Println("🔴 Continue RED LIGHT")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("  Problem: %s\n", problem)
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)

	// Impact
	fmt.Println("What's the impact? (enter numbers, comma-separated)")
	fmt.Println("  1. Blocks launch")
	fmt.Println("  2. Blocks paying users")
	fmt.Println("  3. Causes user frustration")
	fmt.Println("  4. Tech debt accumulation")
	fmt.Println("  5. Other")
	fmt.Print("> ")
	impactStr, _ := reader.ReadString('\n')
	impacts := parseNumbers(strings.TrimSpace(impactStr))

	// Severity
	fmt.Println()
	fmt.Println("Severity?")
	fmt.Println("  0. P0 - Launch blocker")
	fmt.Println("  1. P1 - Paying user blocker")
	fmt.Println("  2. P2 - Nice to have")
	fmt.Println("  3. P3 - Future improvement")
	fmt.Print("> ")
	severityStr, _ := reader.ReadString('\n')
	var severity int
	fmt.Sscanf(strings.TrimSpace(severityStr), "%d", &severity)

	// Read and update current task file
	content, err := os.ReadFile(taskPath)
	if err != nil {
		return err
	}
	text := string(content)

	// Update impact checkboxes
	impactMap := map[int]string{
		1: "Blocks launch",
		2: "Blocks paying users",
		3: "Causes user frustration",
		4: "Tech debt accumulation",
		5: "Other",
	}
	for _, i := range impacts {
		if label, ok := impactMap[i]; ok {
			text = strings.Replace(text, "- [ ] "+label, "- [x] "+label, 1)
		}
	}

	// Update severity
	severities := []string{
		"P0 - Launch blocker",
		"P1 - Paying user blocker",
		"P2 - Nice to have",
		"P3 - Future improvement",
	}
	if severity >= 0 && severity < len(severities) {
		text = strings.Replace(text, "- [ ] "+severities[severity], "- [x] "+severities[severity], 1)
	}

	if err := os.WriteFile(taskPath, []byte(text), 0644); err != nil {
		return err
	}

	// Log activity
	activity.LogStageChange("red", "red", s.CurrentTaskID)

	fmt.Println()
	fmt.Println("✅ RED LIGHT complete!")
	fmt.Printf("   Task ID: %s\n", s.CurrentTaskID)
	fmt.Println("   Next: yo yellow  (analyze and plan)")

	return nil
}

func init() {
	redCmd.Flags().BoolVarP(&redEdit, "edit", "e", false, "Open in editor instead of interactive mode")
	rootCmd.AddCommand(redCmd)
}
