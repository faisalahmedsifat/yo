package cmd

import (
	"fmt"

	"github.com/faisalahmedsifat/yo/internal/workspace"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a yo workspace in the current directory",
	Long: `Creates a .yo/ directory structure with:
  - current_task.md: The RED/YELLOW/GREEN task template
  - backlog.md: Your prioritized task backlog
  - tech_debt_log.md: Deferred decisions and shortcuts
  - AGENTS.md: Instructions for AI agents
  - state.json: Current state (stage, timer, etc.)
  - config.json: User configuration
  - activity.jsonl: Activity log (file changes, events)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := workspace.Init(); err != nil {
			return err
		}
		fmt.Println("✅ Workspace initialized!")
		fmt.Println("")
		fmt.Println("Created .yo/ with:")
		fmt.Println("  📄 current_task.md  - Your RED/YELLOW/GREEN task")
		fmt.Println("  📋 backlog.md       - Prioritized task list")
		fmt.Println("  📝 tech_debt_log.md - Deferred decisions")
		fmt.Println("  🤖 AGENTS.md        - AI agent instructions")
		fmt.Println("  ⚙️  state.json       - Current state")
		fmt.Println("  🔧 config.json      - Configuration")
		fmt.Println("  📊 activity.jsonl   - Activity log")
		fmt.Println("")
		fmt.Println("Get started with:")
		fmt.Println("  yo red    # Define your first problem")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
