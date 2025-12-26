package tests

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestFullWorkflow tests the complete RED -> YELLOW -> GREEN -> DONE flow
func TestFullWorkflow(t *testing.T) {
	// Create temp directory for test workspace
	tmpDir, err := os.MkdirTemp("", "yo-integration-*")
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

	// Test: yo init
	t.Run("yo init", func(t *testing.T) {
		output := runYo(t, tmpDir, yoBinary, "init")
		if !strings.Contains(output, "Workspace initialized") {
			t.Errorf("Expected 'Workspace initialized', got: %s", output)
		}

		// Verify .yo directory was created
		yoDir := filepath.Join(tmpDir, ".yo")
		if _, err := os.Stat(yoDir); os.IsNotExist(err) {
			t.Error(".yo directory was not created")
		}

		// Verify files were created
		expectedFiles := []string{
			"current_task.md",
			"backlog.md",
			"tech_debt_log.md",
			"state.json",
			"config.json",
			"activity.jsonl",
		}
		for _, f := range expectedFiles {
			path := filepath.Join(yoDir, f)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf("Expected file %s was not created", f)
			}
		}
	})

	// Test: yo status (after init)
	t.Run("yo status after init", func(t *testing.T) {
		output := runYo(t, tmpDir, yoBinary, "status")
		if !strings.Contains(output, "NONE") {
			t.Errorf("Expected stage 'NONE', got: %s", output)
		}
	})

	// Test: yo version
	t.Run("yo version", func(t *testing.T) {
		output := runYo(t, tmpDir, yoBinary, "version")
		if !strings.Contains(output, "yo version") {
			t.Errorf("Expected version info, got: %s", output)
		}
	})

	// Test: yo list (empty backlog)
	t.Run("yo list empty", func(t *testing.T) {
		output := runYo(t, tmpDir, yoBinary, "list")
		if !strings.Contains(output, "Backlog") {
			t.Errorf("Expected backlog output, got: %s", output)
		}
	})

	// Test: yo config list
	t.Run("yo config list", func(t *testing.T) {
		output := runYo(t, tmpDir, yoBinary, "config", "list")
		if !strings.Contains(output, "Configuration") {
			t.Errorf("Expected configuration output, got: %s", output)
		}
		if !strings.Contains(output, "notifications") {
			t.Errorf("Expected 'notifications' in config, got: %s", output)
		}
	})

	// Test: Simulate RED LIGHT by modifying the task file
	t.Run("simulate RED LIGHT", func(t *testing.T) {
		taskPath := filepath.Join(tmpDir, ".yo", "current_task.md")
		content := `# Current Task

## 🔴 RED LIGHT - Problem Definition

### What's the Problem?
Integration test problem

### Impact
- [x] Blocks launch
- [ ] Blocks paying users
- [ ] Causes user frustration
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
Test cause

**Underlying cause:**
Test underlying

**System cause:**
Test system

### Solution Options

#### Option A:
- Description: Quick fix
- Time estimate: 1h
- Pros: Fast
- Cons: Hacky

#### Option B:
- Description: Proper fix
- Time estimate: 2h
- Pros: Clean
- Cons: Slower

#### Option C:
- Description: Refactor
- Time estimate: 4h
- Pros: Best
- Cons: Longest

### Decision
**Chosen option:** A
**Reason:** Testing

### Implementation Steps
1. Step 1
2. Step 2

### Success Criteria
- [ ] Criterion 1
- [ ] Criterion 2

---

## 🟢 GREEN LIGHT - Execution

### Timer Started:
### Estimated Time:

### Notes:
### Blockers:

---

## ✅ Completion
`
		if err := os.WriteFile(taskPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write task file: %v", err)
		}

		// Update state to red
		statePath := filepath.Join(tmpDir, ".yo", "state.json")
		stateContent := `{
  "version": "1.0.0",
  "current_stage": "red",
  "current_task_id": "test_task",
  "timer": { "estimated_hours": 0, "threshold_hours": 0 },
  "session": { "active": false },
  "emergency_bypasses": { "today": 0, "this_week": 0, "last_reset": "2024-12-27" }
}`
		if err := os.WriteFile(statePath, []byte(stateContent), 0644); err != nil {
			t.Fatalf("Failed to write state file: %v", err)
		}
	})

	// Test: yo verify red
	t.Run("yo verify red", func(t *testing.T) {
		output := runYo(t, tmpDir, yoBinary, "verify", "red")
		if !strings.Contains(output, "RED LIGHT is complete") {
			t.Errorf("Expected RED to be complete, got: %s", output)
		}
	})

	// Test: yo verify yellow
	t.Run("yo verify yellow", func(t *testing.T) {
		output := runYo(t, tmpDir, yoBinary, "verify", "yellow")
		if !strings.Contains(output, "YELLOW LIGHT is complete") {
			t.Errorf("Expected YELLOW to be complete, got: %s", output)
		}
	})

	// Test: yo go --time 1h
	t.Run("yo go", func(t *testing.T) {
		// First update state to yellow
		statePath := filepath.Join(tmpDir, ".yo", "state.json")
		stateContent := `{
  "version": "1.0.0",
  "current_stage": "yellow",
  "current_task_id": "test_task",
  "timer": { "estimated_hours": 0, "threshold_hours": 0 },
  "session": { "active": false },
  "emergency_bypasses": { "today": 0, "this_week": 0, "last_reset": "2024-12-27" }
}`
		if err := os.WriteFile(statePath, []byte(stateContent), 0644); err != nil {
			t.Fatalf("Failed to write state file: %v", err)
		}

		output := runYo(t, tmpDir, yoBinary, "go", "--time", "1h")
		if !strings.Contains(output, "GREEN LIGHT") {
			t.Errorf("Expected GREEN LIGHT start, got: %s", output)
		}
	})

	// Test: yo timer
	t.Run("yo timer", func(t *testing.T) {
		output := runYo(t, tmpDir, yoBinary, "timer")
		if !strings.Contains(output, "Timer") {
			t.Errorf("Expected timer output, got: %s", output)
		}
	})

	// Test: yo status (should show green)
	t.Run("yo status green", func(t *testing.T) {
		output := runYo(t, tmpDir, yoBinary, "status")
		if !strings.Contains(output, "GREEN") {
			t.Errorf("Expected stage 'GREEN', got: %s", output)
		}
	})

	// Test: yo activity
	t.Run("yo activity", func(t *testing.T) {
		output := runYo(t, tmpDir, yoBinary, "activity")
		if !strings.Contains(output, "Activity") {
			t.Errorf("Expected activity output, got: %s", output)
		}
	})

	// Test: yo focus
	t.Run("yo focus", func(t *testing.T) {
		output := runYo(t, tmpDir, yoBinary, "focus")
		if !strings.Contains(output, "Focus") {
			t.Errorf("Expected focus output, got: %s", output)
		}
	})

	// Test: yo stats
	t.Run("yo stats", func(t *testing.T) {
		output := runYo(t, tmpDir, yoBinary, "stats")
		if !strings.Contains(output, "Week of") {
			t.Errorf("Expected stats output, got: %s", output)
		}
	})
}

// TestInitIdempotency verifies that init fails if already initialized
func TestInitIdempotency(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "yo-init-test-*")
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

	// First init should succeed
	runYo(t, tmpDir, yoBinary, "init")

	// Second init should fail
	cmd := exec.Command(yoBinary, "init")
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("Expected second init to fail")
	}
	if !strings.Contains(string(output), "already initialized") {
		t.Errorf("Expected 'already initialized' error, got: %s", output)
	}
}

// TestStageEnforcement verifies stage transition rules
func TestStageEnforcement(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "yo-stage-test-*")
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

	runYo(t, tmpDir, yoBinary, "init")

	// Try to go directly to GREEN without RED/YELLOW
	cmd := exec.Command(yoBinary, "go")
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Error("Expected 'go' to fail without RED/YELLOW")
	}
	if !strings.Contains(string(output), "complete RED") || !strings.Contains(string(output), "first") {
		t.Errorf("Expected error about completing RED first, got: %s", output)
	}
}

func runYo(t *testing.T, workDir, binary string, args ...string) string {
	t.Helper()

	cmd := exec.Command(binary, args...)
	cmd.Dir = workDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		// Some commands are expected to exit with an error
		// Return combined output for assertion
		return stdout.String() + stderr.String()
	}

	return stdout.String()
}

func getProjectRoot(t *testing.T) string {
	t.Helper()

	// Get the directory containing this test file
	_, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Navigate up to project root
	root, err := filepath.Abs("../")
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	return root
}

// TestInteractiveInputWithSpaces verifies that interactive commands handle multi-word input
func TestInteractiveInputWithSpaces(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "yo-input-test-*")
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

	// Initialize
	runYo(t, tmpDir, yoBinary, "init")

	// Test yo defer with multi-word input
	// yo defer accepts a description with spaces
	t.Run("defer with spaces", func(t *testing.T) {
		cmd := exec.Command(yoBinary, "defer", "No OAuth support - using password only for MVP")
		cmd.Dir = tmpDir
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("defer command failed: %v\n%s", err, output)
		}
		if !strings.Contains(string(output), "Tech debt logged") {
			t.Errorf("Expected 'Tech debt logged', got: %s", output)
		}

		// Verify tech debt log contains the full text
		techDebt, err := os.ReadFile(filepath.Join(tmpDir, ".yo", "tech_debt_log.md"))
		if err != nil {
			t.Fatalf("Failed to read tech_debt_log.md: %v", err)
		}
		if !strings.Contains(string(techDebt), "No OAuth support - using password only for MVP") {
			t.Errorf("Tech debt log missing full description: %s", techDebt)
		}
	})

	// Test yo add with multi-word description
	t.Run("add with spaces", func(t *testing.T) {
		cmd := exec.Command(yoBinary, "add", "Fix the login button on mobile Safari")
		cmd.Dir = tmpDir
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("add command failed: %v\n%s", err, output)
		}
		if !strings.Contains(string(output), "Added to") {
			t.Errorf("Expected 'Added to', got: %s", output)
		}

		// Verify backlog contains the full text
		backlog, err := os.ReadFile(filepath.Join(tmpDir, ".yo", "backlog.md"))
		if err != nil {
			t.Fatalf("Failed to read backlog.md: %v", err)
		}
		if !strings.Contains(string(backlog), "Fix the login button on mobile Safari") {
			t.Errorf("Backlog missing full description: %s", backlog)
		}
	})

	// Test yo bypass with multi-word reason
	t.Run("bypass with spaces", func(t *testing.T) {
		cmd := exec.Command(yoBinary, "bypass", "Production is on fire and users are complaining")
		cmd.Dir = tmpDir
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("bypass command failed: %v\n%s", err, output)
		}
		if !strings.Contains(string(output), "BYPASS") {
			t.Errorf("Expected 'BYPASS', got: %s", output)
		}
	})
}
