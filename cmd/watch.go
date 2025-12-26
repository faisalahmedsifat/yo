package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/faisal/yo/internal/watcher"
	"github.com/faisal/yo/internal/workspace"
	"github.com/spf13/cobra"
)

var watchBackground bool

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Start the file watcher daemon",
	Long: `Start a global file watcher that monitors all projects in your watch_dirs.

The watcher:
  - Runs as a single daemon for all projects
  - Detects file changes in git repositories
  - Logs activity to the current active project's .yo/activity.jsonl
  - Marks changes as "untracked" if outside the current task's repo

Use --bg to run in background.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if already running
		if running, pid := watcher.IsRunning(); running {
			fmt.Printf("⚠️  Watcher already running (PID: %d)\n", pid)
			fmt.Println("   Use 'yo watch stop' to stop it first.")
			return nil
		}

		// Set current project directory
		if workspace.IsInitialized() {
			cwd, _ := os.Getwd()
			watcher.SetCurrentProject(cwd)
		}

		if watchBackground {
			return startBackground()
		}

		return startForeground()
	},
}

var watchStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the file watcher daemon",
	RunE: func(cmd *cobra.Command, args []string) error {
		running, pid := watcher.IsRunning()
		if !running {
			fmt.Println("⚠️  Watcher is not running.")
			return nil
		}

		// Send SIGTERM to the process
		process, err := os.FindProcess(pid)
		if err != nil {
			return fmt.Errorf("failed to find process: %w", err)
		}

		if err := process.Signal(syscall.SIGTERM); err != nil {
			return fmt.Errorf("failed to stop watcher: %w", err)
		}

		fmt.Printf("✅ Watcher stopped (PID: %d)\n", pid)
		return nil
	},
}

var watchStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check if watcher daemon is running",
	RunE: func(cmd *cobra.Command, args []string) error {
		running, pid := watcher.IsRunning()
		if running {
			fmt.Printf("🟢 Watcher is running (PID: %d)\n", pid)
		} else {
			fmt.Println("⚪ Watcher is not running")
			fmt.Println("   Start with: yo watch")
		}
		return nil
	},
}

func startForeground() error {
	w, err := watcher.New()
	if err != nil {
		return err
	}

	if err := w.Start(); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("🔍 File watcher started (press Ctrl+C to stop)")
	fmt.Println()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println()
	fmt.Println("👋 Stopping watcher...")
	w.Stop()

	return nil
}

func startBackground() error {
	// Get the path to the yo binary
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Start a new process without --bg flag
	cmd := exec.Command(executable, "watch")
	cmd.Stdout = nil
	cmd.Stderr = nil
	cmd.Stdin = nil

	// Detach from parent
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setsid: true,
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start background watcher: %w", err)
	}

	fmt.Printf("🔍 Watcher started in background (PID: %d)\n", cmd.Process.Pid)
	fmt.Println("   Stop with: yo watch stop")
	fmt.Println("   Status:    yo watch status")

	return nil
}

func init() {
	watchCmd.Flags().BoolVar(&watchBackground, "bg", false, "Run in background")
	watchCmd.AddCommand(watchStopCmd)
	watchCmd.AddCommand(watchStatusCmd)
	rootCmd.AddCommand(watchCmd)
}
