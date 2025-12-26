package notify

import (
	"testing"
)

func TestNew(t *testing.T) {
	n := New(true)
	if !n.Enabled {
		t.Error("Expected notifier to be enabled")
	}

	n = New(false)
	if n.Enabled {
		t.Error("Expected notifier to be disabled")
	}
}

func TestSendDisabled(t *testing.T) {
	n := New(false)

	// Should not error when disabled
	err := n.Send("Test", "Message")
	if err != nil {
		t.Errorf("Expected no error when disabled, got: %v", err)
	}
}

func TestTimerMilestones(t *testing.T) {
	n := New(false) // Disabled to avoid actual notifications

	// These should not error
	if err := n.TimerMilestone100("test_task"); err != nil {
		t.Errorf("TimerMilestone100 error: %v", err)
	}

	if err := n.TimerMilestone150("test_task"); err != nil {
		t.Errorf("TimerMilestone150 error: %v", err)
	}

	if err := n.TimerMilestone200("test_task"); err != nil {
		t.Errorf("TimerMilestone200 error: %v", err)
	}
}

func TestBypassNotifications(t *testing.T) {
	n := New(false)

	if err := n.BypassStarted(30); err != nil {
		t.Errorf("BypassStarted error: %v", err)
	}

	if err := n.BypassEnded(); err != nil {
		t.Errorf("BypassEnded error: %v", err)
	}
}

func TestTaskNotifications(t *testing.T) {
	n := New(false)

	if err := n.TaskComplete("test_task", 95.5); err != nil {
		t.Errorf("TaskComplete error: %v", err)
	}

	if err := n.SessionEnded(85.0); err != nil {
		t.Errorf("SessionEnded error: %v", err)
	}

	if err := n.GreenLightStarted("test_task", 2.0); err != nil {
		t.Errorf("GreenLightStarted error: %v", err)
	}
}
