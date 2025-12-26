package cmd

import (
	"fmt"

	"github.com/faisalahmedsifat/yo/internal/config"
	"github.com/faisalahmedsifat/yo/internal/workspace"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `View and modify yo configuration settings.`,
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configuration values",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !workspace.IsInitialized() {
			return fmt.Errorf("workspace not initialized. Run 'yo init' first")
		}

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		fmt.Println()
		fmt.Println("⚙️  Configuration")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()
		fmt.Printf("  notifications:  %v\n", cfg.Notifications)
		fmt.Printf("  editor:         %s\n", cfg.Editor)
		fmt.Printf("  watch_dirs:     %v\n", cfg.WatchDirs)
		fmt.Printf("  max_bypass_day: %d\n", cfg.MaxBypassDay)
		fmt.Printf("  max_bypass_week: %d\n", cfg.MaxBypassWeek)
		fmt.Println()

		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value.

Available keys:
  notifications  - on/off
  editor         - path to editor`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !workspace.IsInitialized() {
			return fmt.Errorf("workspace not initialized. Run 'yo init' first")
		}

		key := args[0]
		value := args[1]

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if err := cfg.Set(key, value); err != nil {
			return err
		}

		if err := cfg.Save(); err != nil {
			return err
		}

		fmt.Printf("✅ Set %s = %s\n", key, value)
		return nil
	},
}

var configGetCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "Get a configuration value",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if !workspace.IsInitialized() {
			return fmt.Errorf("workspace not initialized. Run 'yo init' first")
		}

		key := args[0]

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		value, err := cfg.Get(key)
		if err != nil {
			return err
		}

		fmt.Printf("%s = %s\n", key, value)
		return nil
	},
}

func init() {
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	rootCmd.AddCommand(configCmd)
}
