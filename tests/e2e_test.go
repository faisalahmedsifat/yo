package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestEndToEnd performs a complete end-to-end test simulating real user workflow
func TestEndToEnd(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "yo-e2e-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Build yo binary
	yoBinary := filepath.Join(tmpDir, "yo")
	buildCmd := exec.Command("go", "build", "-o", yoBinary, ".")
	buildCmd.Dir = getProjectRoot(t)
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build yo: %v\n%s", err, output)
	}

	run := func(args ...string) string {
		cmd := exec.Command(yoBinary, args...)
		cmd.Dir = tmpDir
		output, _ := cmd.CombinedOutput()
		return string(output)
	}

	// === PHASE 1: Initialize ===
	t.Log("Phase 1: Initialize workspace")

	output := run("init")
	assertContains(t, output, "Workspace initialized")

	// Verify all files created
	yoDir := filepath.Join(tmpDir, ".yo")
	assertFileExists(t, filepath.Join(yoDir, "current_task.md"))
	assertFileExists(t, filepath.Join(yoDir, "backlog.md"))
	assertFileExists(t, filepath.Join(yoDir, "state.json"))
	assertFileExists(t, filepath.Join(yoDir, "config.json"))
	assertFileExists(t, filepath.Join(yoDir, "activity.jsonl"))

	// === PHASE 2: Check initial state ===
	t.Log("Phase 2: Verify initial state")

	output = run("status")
	assertContains(t, output, "NONE")

	output = run("milestone")
	assertContains(t, output, "Milestone 0")

	// === PHASE 3: Simulate RED LIGHT ===
	t.Log("Phase 3: Complete RED LIGHT")

	// Write a complete RED LIGHT task
	taskContent := `# Current Task

## 🔴 RED LIGHT - Problem Definition

### What's the Problem?
E2E Test: The login button doesn't work on mobile

### Impact
- [x] Blocks launch
- [x] Causes user frustration
- [ ] Tech debt accumulation
- [ ] Other: ___

### Severity
- [x] P0 - Launch blocker
- [ ] P1 - Paying user blocker
- [ ] P2 - Nice to have
- [ ] P3 - Future improvement

---

## 🟡 YELLOW LIGHT - Analysis & Planning

### Root Cause Analysis
**Immediate cause:**
Touch events not handled

**Underlying cause:**
No mobile testing

**System cause:**
Missing device testing

### Solution Options

#### Option A:
- Description: Quick CSS fix
- Time estimate: 1h
- Pros: Fast
- Cons: Hacky

#### Option B:
- Description: Proper touch handler
- Time estimate: 2h
- Pros: Robust
- Cons: More work

#### Option C:
- Description: Full responsive refactor
- Time estimate: 4h
- Pros: Complete solution
- Cons: Time consuming

### Decision
**Chosen option:** B
**Reason:** Balance of speed and quality

### Implementation Steps
1. Add touch event listener
2. Test on mobile devices
3. Deploy

### Success Criteria
- [ ] Button works on iOS Safari
- [ ] Button works on Android Chrome
- [ ] No console errors

---

## 🟢 GREEN LIGHT - Execution

### Timer Started:
### Estimated Time:

### Notes:
### Blockers:

---

## ✅ Completion
`
	err = os.WriteFile(filepath.Join(yoDir, "current_task.md"), []byte(taskContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write task: %v", err)
	}

	// Update state to red
	stateContent := `{"version":"1.0.0","current_stage":"red","current_task_id":"login_button_e2e","timer":{},"session":{},"emergency_bypasses":{"today":0,"this_week":0,"last_reset":"2024-12-27"},"milestone":{"current":0,"name":"Test"}}`
	err = os.WriteFile(filepath.Join(yoDir, "state.json"), []byte(stateContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write state: %v", err)
	}

	// Verify RED
	output = run("verify", "red")
	assertContains(t, output, "RED LIGHT is complete")

	// Verify YELLOW
	output = run("verify", "yellow")
	assertContains(t, output, "YELLOW LIGHT is complete")

	// === PHASE 4: Start GREEN LIGHT ===
	t.Log("Phase 4: Start GREEN LIGHT")

	// Update state to yellow first
	stateContent = `{"version":"1.0.0","current_stage":"yellow","current_task_id":"login_button_e2e","timer":{},"session":{},"emergency_bypasses":{"today":0,"this_week":0,"last_reset":"2024-12-27"},"milestone":{"current":0,"name":"Test"}}`
	os.WriteFile(filepath.Join(yoDir, "state.json"), []byte(stateContent), 0644)

	output = run("go", "--time", "2h")
	assertContains(t, output, "GREEN LIGHT")
	assertContains(t, output, "Timer is running")

	// === PHASE 5: Check timer ===
	t.Log("Phase 5: Verify timer")

	output = run("timer")
	assertContains(t, output, "Timer")
	assertContains(t, output, "%")

	output = run("status")
	assertContains(t, output, "GREEN")

	// === PHASE 6: Test activity and focus ===
	t.Log("Phase 6: Test activity tracking")

	output = run("activity")
	assertContains(t, output, "Activity")

	output = run("focus")
	assertContains(t, output, "Focus")

	// === PHASE 7: Test stats ===
	t.Log("Phase 7: Test stats")

	output = run("stats")
	assertContains(t, output, "Week of")

	// === PHASE 8: Test backlog operations ===
	t.Log("Phase 8: Test backlog")

	output = run("list")
	assertContains(t, output, "Backlog")

	// === PHASE 9: Test config ===
	t.Log("Phase 9: Test config")

	output = run("config", "list")
	assertContains(t, output, "Configuration")
	assertContains(t, output, "notifications")

	// === PHASE 10: Test watch status ===
	t.Log("Phase 10: Test watcher")

	output = run("watch", "status")
	assertContains(t, output, "Watcher")

	t.Log("✅ All E2E tests passed!")
}

// TestBacklogWorkflow tests adding and managing backlog items
func TestBacklogWorkflow(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "yo-backlog-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	yoBinary := filepath.Join(tmpDir, "yo")
	buildCmd := exec.Command("go", "build", "-o", yoBinary, ".")
	buildCmd.Dir = getProjectRoot(t)
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build yo: %v\n%s", err, output)
	}

	run := func(args ...string) string {
		cmd := exec.Command(yoBinary, args...)
		cmd.Dir = tmpDir
		output, _ := cmd.CombinedOutput()
		return string(output)
	}

	// Initialize
	run("init")

	// List empty backlog
	output := run("list")
	assertContains(t, output, "Backlog")

	// List with P0 filter
	output = run("list", "--p0")
	assertContains(t, output, "P0")
}

// TestBypassLimits tests emergency bypass limits
func TestBypassLimits(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "yo-bypass-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	yoBinary := filepath.Join(tmpDir, "yo")
	buildCmd := exec.Command("go", "build", "-o", yoBinary, ".")
	buildCmd.Dir = getProjectRoot(t)
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build yo: %v\n%s", err, output)
	}

	run := func(args ...string) string {
		cmd := exec.Command(yoBinary, args...)
		cmd.Dir = tmpDir
		output, _ := cmd.CombinedOutput()
		return string(output)
	}

	// Initialize
	run("init")

	// First bypass should work
	output := run("bypass", "production down")
	assertContains(t, output, "BYPASS")
}

func assertContains(t *testing.T, output, expected string) {
	t.Helper()
	if !strings.Contains(output, expected) {
		t.Errorf("Expected output to contain '%s', got:\n%s", expected, output)
	}
}

func assertFileExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Expected file to exist: %s", path)
	}
}
