package notify

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Notifier sends desktop notifications
type Notifier struct {
	Enabled bool
}

// New creates a new notifier
func New(enabled bool) *Notifier {
	return &Notifier{Enabled: enabled}
}

// Send sends a desktop notification
func (n *Notifier) Send(title, message string) error {
	if !n.Enabled {
		return nil
	}

	switch runtime.GOOS {
	case "darwin":
		return n.sendMacOS(title, message)
	case "linux":
		return n.sendLinux(title, message)
	default:
		// Fallback: just print to console
		fmt.Printf("🔔 %s: %s\n", title, message)
		return nil
	}
}

// sendMacOS sends notification via osascript
func (n *Notifier) sendMacOS(title, message string) error {
	script := fmt.Sprintf(`display notification "%s" with title "%s"`, message, title)
	return exec.Command("osascript", "-e", script).Run()
}

// sendLinux sends notification via notify-send
func (n *Notifier) sendLinux(title, message string) error {
	return exec.Command("notify-send", title, message).Run()
}

// Timer notifications
func (n *Notifier) TimerMilestone100(taskName string) error {
	return n.Send("⏱️ Timer", fmt.Sprintf("Hit your estimate for %s. Still on track?", taskName))
}

func (n *Notifier) TimerMilestone150(taskName string) error {
	return n.Send("⚠️ Timer", fmt.Sprintf("50%% over estimate for %s", taskName))
}

func (n *Notifier) TimerMilestone200(taskName string) error {
	return n.Send("🚨 Timer", fmt.Sprintf("Doubled your estimate for %s!", taskName))
}

// Bypass notifications
func (n *Notifier) BypassStarted(minutes int) error {
	return n.Send("🚨 BYPASS Active", fmt.Sprintf("Fix it in %d minutes", minutes))
}

func (n *Notifier) BypassEnded() error {
	return n.Send("⏰ Bypass Expired", "Time to document what you fixed")
}

// Task notifications
func (n *Notifier) TaskComplete(taskName string, accuracy float64) error {
	return n.Send("✅ Task Complete", fmt.Sprintf("%s (%.0f%% accuracy)", taskName, accuracy))
}

func (n *Notifier) SessionEnded(focusScore float64) error {
	return n.Send("👋 Session Ended", fmt.Sprintf("Focus score: %.0f%%", focusScore))
}

// GreenLightStarted notifies that execution has begun
func (n *Notifier) GreenLightStarted(taskName string, hours float64) error {
	return n.Send("🟢 GREEN LIGHT", fmt.Sprintf("Timer started: %.1fh for %s", hours, taskName))
}
