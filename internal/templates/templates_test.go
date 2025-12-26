package templates

import (
	"strings"
	"testing"
)

func TestCurrentTaskTemplate(t *testing.T) {
	// Check required sections
	sections := []string{
		"🔴 RED LIGHT",
		"🟡 YELLOW LIGHT",
		"🟢 GREEN LIGHT",
		"What's the Problem?",
		"Impact",
		"Severity",
		"Root Cause Analysis",
		"Solution Options",
		"Option A",
		"Option B",
		"Option C",
		"Success Criteria",
		"Completion",
	}

	for _, section := range sections {
		if !strings.Contains(CurrentTask, section) {
			t.Errorf("CurrentTask template missing section: %s", section)
		}
	}
}

func TestBacklogTemplate(t *testing.T) {
	priorities := []string{"P0", "P1", "P2", "P3"}

	for _, p := range priorities {
		if !strings.Contains(Backlog, p) {
			t.Errorf("Backlog template missing priority: %s", p)
		}
	}
}

func TestTechDebtLogTemplate(t *testing.T) {
	if !strings.Contains(TechDebtLog, "Tech Debt") {
		t.Error("TechDebtLog template missing title")
	}
}

func TestSessionSummaryTemplate(t *testing.T) {
	fields := []string{
		"{{.Date}}",
		"{{.StartedAt}}",
		"{{.EndedAt}}",
		"{{.Duration}}",
		"{{.FocusScore}}",
	}

	for _, field := range fields {
		if !strings.Contains(SessionSummary, field) {
			t.Errorf("SessionSummary template missing field: %s", field)
		}
	}
}

func TestTemplatesNotEmpty(t *testing.T) {
	if len(CurrentTask) < 100 {
		t.Error("CurrentTask template is too short")
	}

	if len(Backlog) < 50 {
		t.Error("Backlog template is too short")
	}

	if len(TechDebtLog) < 20 {
		t.Error("TechDebtLog template is too short")
	}
}
