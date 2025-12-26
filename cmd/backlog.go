package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/faisalahmedsifat/yo/internal/backlog"
	"github.com/faisalahmedsifat/yo/internal/state"
	"github.com/faisalahmedsifat/yo/internal/workspace"
	"github.com/spf13/cobra"
)

var (
	listP0Only bool
	listP1Only bool
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List backlog items",
	Long: `Display items from the backlog organized by priority.

Examples:
  yo list       - Show all items
  yo list --p0  - Show only P0 items
  yo list --p1  - Show only P1 items`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !workspace.IsInitialized() {
			return fmt.Errorf("workspace not initialized. Run 'yo init' first")
		}

		b, err := backlog.Load()
		if err != nil {
			return err
		}

		if listP0Only {
			printPriorityItems("P0 - Launch Blockers", b.Items[backlog.P0])
			return nil
		}

		if listP1Only {
			printPriorityItems("P1 - Paying User Blockers", b.Items[backlog.P1])
			return nil
		}

		fmt.Println()
		fmt.Println("📋 Backlog")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()

		priorities := []struct {
			key   string
			title string
		}{
			{backlog.P0, "P0 - Launch Blockers"},
			{backlog.P1, "P1 - Paying User Blockers"},
			{backlog.P2, "P2 - Nice to Have"},
			{backlog.P3, "P3 - Future Improvements"},
		}

		totalItems := b.Total()
		for _, p := range priorities {
			if len(b.Items[p.key]) > 0 {
				printPriorityItems(p.title, b.Items[p.key])
			}
		}

		if totalItems == 0 {
			fmt.Println("  (no items in backlog)")
			fmt.Println()
			fmt.Println("  Add items with: yo add \"description\"")
		}

		return nil
	},
}

func printPriorityItems(title string, items []backlog.Item) {
	fmt.Printf("  ## %s (%d)\n", title, len(items))
	if len(items) == 0 {
		fmt.Println("    (none)")
	}
	for i, item := range items {
		checkbox := "[ ]"
		if item.Checked {
			checkbox = "[x]"
		}
		fmt.Printf("    %d. %s %s\n", i+1, checkbox, item.Text)
	}
	fmt.Println()
}

var addCmd = &cobra.Command{
	Use:   "add [description]",
	Short: "Add item to backlog",
	Long: `Add a new item to the backlog.

Examples:
  yo add "OAuth login broken"
  yo add -i  # Interactive mode`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !workspace.IsInitialized() {
			return fmt.Errorf("workspace not initialized. Run 'yo init' first")
		}

		var description string
		var priority string

		interactive, _ := cmd.Flags().GetBool("interactive")

		if interactive || len(args) == 0 {
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Description: ")
			description, _ = reader.ReadString('\n')
			description = strings.TrimSpace(description)

			if description == "" {
				return fmt.Errorf("description cannot be empty")
			}

			fmt.Println()
			fmt.Println("Priority:")
			fmt.Println("  0. P0 - Launch blocker")
			fmt.Println("  1. P1 - Paying user blocker")
			fmt.Println("  2. P2 - Nice to have")
			fmt.Println("  3. P3 - Future improvement")
			fmt.Print("> ")
			prioStr, _ := reader.ReadString('\n')
			prioStr = strings.TrimSpace(prioStr)

			switch prioStr {
			case "0":
				priority = backlog.P0
			case "1":
				priority = backlog.P1
			case "2":
				priority = backlog.P2
			case "3":
				priority = backlog.P3
			default:
				priority = backlog.P2 // Default
			}
		} else {
			description = strings.Join(args, " ")
			priority = backlog.P2 // Default priority for non-interactive
		}

		// Add to backlog using backlog package
		b, err := backlog.Load()
		if err != nil {
			return err
		}

		if err := b.Add(description, priority); err != nil {
			return err
		}

		fmt.Println()
		fmt.Printf("✅ Added to %s: %s\n", priority, description)
		fmt.Println()
		fmt.Print("Work on it now? (y/n): ")

		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		if strings.TrimSpace(strings.ToLower(response)) == "y" {
			fmt.Println("  Start with: yo red")
		}

		return nil
	},
}

var nextCmd = &cobra.Command{
	Use:   "next",
	Short: "Pick next task from backlog and start RED LIGHT",
	Long:  `Select the next task from your backlog (prioritizing P0 items) and start RED LIGHT.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !workspace.IsInitialized() {
			return fmt.Errorf("workspace not initialized. Run 'yo init' first")
		}

		s, err := state.Load()
		if err != nil {
			return err
		}

		// Check if already in a task
		if s.CurrentStage != "none" && s.CurrentStage != "" {
			fmt.Printf("⚠️  Already in %s stage with task: %s\n", strings.ToUpper(s.CurrentStage), s.CurrentTaskID)
			fmt.Println("   Complete current task first with 'yo done' or 'yo off'")
			return nil
		}

		bl, err := backlog.Load()
		if err != nil {
			return err
		}

		// Get unchecked items using backlog package
		available := bl.GetUnchecked()

		if len(available) == 0 {
			fmt.Println("🎉 Backlog is empty! Add items with: yo add \"description\"")
			return nil
		}

		fmt.Println()
		fmt.Println("📋 Pick next task:")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()

		for i, item := range available {
			fmt.Printf("  [%d] %s (%s)\n", i+1, item.Text, item.Priority)
		}
		fmt.Println("  [q] Quit")
		fmt.Println()

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(response)

		if response == "q" || response == "" {
			return nil
		}

		var choice int
		if _, err := fmt.Sscanf(response, "%d", &choice); err != nil || choice < 1 || choice > len(available) {
			return fmt.Errorf("invalid choice")
		}

		selected := available[choice-1]

		// Generate task ID from description
		taskID := sanitizeForTaskID(selected.Text)

		// Update current_task.md with the problem from backlog
		taskPath, err := workspace.GetCurrentTaskPath()
		if err != nil {
			return err
		}
		content, err := os.ReadFile(taskPath)
		if err != nil {
			return err
		}
		text := string(content)
		text = strings.Replace(text, "<!-- Describe the problem clearly -->", selected.Text, 1)
		os.WriteFile(taskPath, []byte(text), 0644)

		// Update state to RED LIGHT
		s.SetStage("red")
		s.CurrentTaskID = taskID
		if err := s.Save(); err != nil {
			return err
		}

		fmt.Println()
		fmt.Println("🔴 RED LIGHT Started!")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("  Task: %s\n", selected.Text)
		fmt.Printf("  ID:   %s\n", taskID)
		fmt.Printf("  From: %s\n", selected.Priority)
		fmt.Println()
		fmt.Println("  The problem has been added to current_task.md")
		fmt.Println()
		fmt.Println("  Next steps:")
		fmt.Println("    1. Edit .yo/current_task.md to complete RED LIGHT")
		fmt.Println("    2. Or run: yo red  (interactive mode)")
		fmt.Println("    3. Then: yo yellow  (plan your solution)")

		return nil
	},
}

func sanitizeForTaskID(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "_")
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

func init() {
	listCmd.Flags().BoolVar(&listP0Only, "p0", false, "Show only P0 items")
	listCmd.Flags().BoolVar(&listP1Only, "p1", false, "Show only P1 items")
	addCmd.Flags().BoolP("interactive", "i", false, "Interactive mode")

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(nextCmd)
}
