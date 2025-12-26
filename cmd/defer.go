package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/faisalahmedsifat/yo/internal/state"
	"github.com/faisalahmedsifat/yo/internal/workspace"
	"github.com/spf13/cobra"
)

var deferInteractive bool

var deferCmd = &cobra.Command{
	Use:   "defer [description]",
	Short: "Log tech debt - a conscious decision to defer work",
	Long: `Log tech debt during YELLOW LIGHT planning.

Tech debt = shortcuts you CHOSE to take to ship faster.
It's when you say: "I know the proper way, but I'm doing the quick way."

Examples:
  yo defer "No retry button - users can click deploy again"
  yo defer -i                    # Interactive mode with guided prompts

This logs to .yo/tech_debt_log.md for future reference.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !workspace.IsInitialized() {
			return fmt.Errorf("workspace not initialized. Run 'yo init' first")
		}

		s, err := state.Load()
		if err != nil {
			return err
		}

		// Encourage use during YELLOW LIGHT
		if s.CurrentStage != "yellow" && s.CurrentStage != "green" {
			fmt.Println("💡 Tip: Log tech debt during YELLOW LIGHT when deciding what to defer.")
			fmt.Println()
		}

		if deferInteractive || len(args) == 0 {
			return runDeferInteractive(s)
		}

		// Quick mode
		what := strings.Join(args, " ")
		if err := logTechDebt(s.CurrentTaskID, what, "Deferred for faster shipping", "When needed", "TBD"); err != nil {
			return err
		}
		fmt.Println()
		fmt.Println("✅ Tech debt logged!")
		fmt.Printf("   What: %s\n", what)
		fmt.Println("   View with: cat .yo/tech_debt_log.md")
		return nil
	},
}

func runDeferInteractive(s *state.State) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("📝 Log Tech Debt")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("Tech debt = conscious shortcuts you're choosing to take.")
	fmt.Println()

	// What
	fmt.Print("What are you deferring?\n> ")
	what, _ := reader.ReadString('\n')
	what = strings.TrimSpace(what)
	if what == "" {
		return fmt.Errorf("description cannot be empty")
	}

	// Why
	fmt.Print("\nWhy are you skipping it? (what's the tradeoff)\n> ")
	why, _ := reader.ReadString('\n')
	why = strings.TrimSpace(why)
	if why == "" {
		why = "Ship faster for MVP"
	}

	// When to fix
	fmt.Print("\nWhen should you come back to fix it?\n")
	fmt.Println("  Examples: 'After 50 users', 'Before public launch', 'When user complains'")
	fmt.Print("> ")
	when, _ := reader.ReadString('\n')
	when = strings.TrimSpace(when)
	if when == "" {
		when = "When needed"
	}

	// Time estimate
	fmt.Print("\nEstimated time to fix later? (e.g., 2h, 4h)\n> ")
	estimate, _ := reader.ReadString('\n')
	estimate = strings.TrimSpace(estimate)
	if estimate == "" {
		estimate = "TBD"
	}

	taskID := s.CurrentTaskID
	if taskID == "" {
		taskID = "General"
	}

	if err := logTechDebt(taskID, what, why, when, estimate); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("✅ Tech debt logged!")
	fmt.Println("   View with: cat .yo/tech_debt_log.md")
	fmt.Println()
	fmt.Println("   Remember: This is a CONSCIOUS choice, not bad code.")
	fmt.Println("   You'll fix it when the time is right. 👍")

	return nil
}

func logTechDebt(taskID, what, why, when, estimate string) error {
	techDebtPath, err := workspace.GetTechDebtPath()
	if err != nil {
		return err
	}

	entry := fmt.Sprintf(`
## Deferred on %s
**Task:** %s

**What:** %s
**Why skipped:** %s
**Come back when:** %s
**Estimated fix time:** %s

---
`,
		time.Now().Format("2006-01-02"),
		taskID,
		what,
		why,
		when,
		estimate,
	)

	f, err := os.OpenFile(techDebtPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(entry)
	return err
}

func init() {
	deferCmd.Flags().BoolVarP(&deferInteractive, "interactive", "i", false, "Interactive mode with guided prompts")
	rootCmd.AddCommand(deferCmd)
}
