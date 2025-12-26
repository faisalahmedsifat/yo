package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version information (set during build)
	Version   = "0.1.0"
	BuildDate = "development"
	GitCommit = "unknown"
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "yo",
	Short: "A CLI to enforce the RED/YELLOW/GREEN productivity framework",
	Long: `yo is a productivity CLI that helps you follow the RED/YELLOW/GREEN framework:

🔴 RED LIGHT: Define the problem clearly before coding
🟡 YELLOW LIGHT: Analyze and plan your solution
🟢 GREEN LIGHT: Execute with a timer and focus tracking

yo tracks your progress, monitors file changes, and provides helpful
nudges to keep you on track.`,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags can be added here
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
}
