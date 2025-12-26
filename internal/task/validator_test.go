package task

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateRedComplete(t *testing.T) {
	// Create temp file with complete RED LIGHT
	content := `# Current Task

## 🔴 RED LIGHT - Problem Definition

### What's the Problem?
The deploy button is broken

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
`

	tmpFile := createTempFile(t, content)
	defer os.Remove(tmpFile)

	result, err := ValidateRed(tmpFile)
	if err != nil {
		t.Fatalf("ValidateRed failed: %v", err)
	}

	if !result.Valid {
		t.Errorf("Expected RED to be valid, got errors: %v", result.Errors)
	}
}

func TestValidateRedIncomplete(t *testing.T) {
	// Create temp file with incomplete RED LIGHT (no impact selected)
	content := `# Current Task

## 🔴 RED LIGHT - Problem Definition

### What's the Problem?
<!-- Describe the problem clearly -->

### Impact
- [ ] Blocks launch
- [ ] Blocks paying users
- [ ] Causes user frustration
- [ ] Tech debt accumulation
- [ ] Other: ___

### Severity
- [ ] P0 - Launch blocker
- [ ] P1 - Paying user blocker
- [ ] P2 - Nice to have
- [ ] P3 - Future improvement

---

## 🟡 YELLOW LIGHT - Analysis & Planning
`

	tmpFile := createTempFile(t, content)
	defer os.Remove(tmpFile)

	result, err := ValidateRed(tmpFile)
	if err != nil {
		t.Fatalf("ValidateRed failed: %v", err)
	}

	if result.Valid {
		t.Error("Expected RED to be invalid")
	}

	if len(result.Errors) == 0 {
		t.Error("Expected validation errors")
	}
}

func TestValidateYellowComplete(t *testing.T) {
	content := `# Current Task

## 🔴 RED LIGHT - Problem Definition
(complete)

---

## 🟡 YELLOW LIGHT - Analysis & Planning

### Root Cause Analysis
**Immediate cause:**
API endpoint missing

**Underlying cause:**
No route defined

**System cause:**
Missing tests

### Solution Options

#### Option A:
- Description: Quick fix
- Time estimate: 2h
- Pros: Fast
- Cons: Hacky

#### Option B:
- Description: Proper fix
- Time estimate: 4h
- Pros: Clean
- Cons: Slower

#### Option C:
- Description: Refactor
- Time estimate: 8h
- Pros: Best
- Cons: Longest

### Decision
**Chosen option:** B
**Reason:** Balance of speed and quality

### Implementation Steps
1. Add route
2. Add tests
3. Deploy

### Success Criteria
- [ ] Route works
- [ ] Tests pass
- [ ] No errors in prod

---

## 🟢 GREEN LIGHT - Execution
`

	tmpFile := createTempFile(t, content)
	defer os.Remove(tmpFile)

	result, err := ValidateYellow(tmpFile)
	if err != nil {
		t.Fatalf("ValidateYellow failed: %v", err)
	}

	if !result.Valid {
		t.Errorf("Expected YELLOW to be valid, got errors: %v", result.Errors)
	}
}

func TestValidateYellowIncompleteOptions(t *testing.T) {
	content := `# Current Task

## 🟡 YELLOW LIGHT - Analysis & Planning

### Solution Options

#### Option A:
- Description: Quick fix
- Time estimate: 2h

### Decision
**Chosen option:** A

### Success Criteria
- [ ] Works
- [ ] Tests pass

---

## 🟢 GREEN LIGHT - Execution
`

	tmpFile := createTempFile(t, content)
	defer os.Remove(tmpFile)

	result, err := ValidateYellow(tmpFile)
	if err != nil {
		t.Fatalf("ValidateYellow failed: %v", err)
	}

	if result.Valid {
		t.Error("Expected YELLOW to be invalid (only 1 option)")
	}

	// Should have error about needing 3 options
	found := false
	for _, e := range result.Errors {
		if e.Field == "Solution Options" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected error about Solution Options")
	}
}

func TestGetSuccessCriteria(t *testing.T) {
	content := `# Current Task

## 🟡 YELLOW LIGHT - Analysis & Planning

### Success Criteria
- [ ] Route works
- [ ] Tests pass
- [ ] No errors in prod

---

## 🟢 GREEN LIGHT - Execution
`

	tmpFile := createTempFile(t, content)
	defer os.Remove(tmpFile)

	criteria, err := GetSuccessCriteria(tmpFile)
	if err != nil {
		t.Fatalf("GetSuccessCriteria failed: %v", err)
	}

	if len(criteria) != 3 {
		t.Errorf("Expected 3 criteria, got %d", len(criteria))
	}

	expected := []string{"Route works", "Tests pass", "No errors in prod"}
	for i, c := range expected {
		if criteria[i] != c {
			t.Errorf("Expected criteria[%d] = '%s', got '%s'", i, c, criteria[i])
		}
	}
}

func TestGetTimeEstimate(t *testing.T) {
	content := `# Task

## YELLOW

### Decision
**Chosen option:** A
Time estimate: 4h
`

	tmpFile := createTempFile(t, content)
	defer os.Remove(tmpFile)

	hours, err := GetTimeEstimate(tmpFile)
	if err != nil {
		t.Fatalf("GetTimeEstimate failed: %v", err)
	}

	if hours != 4.0 {
		t.Errorf("Expected 4.0 hours, got %f", hours)
	}
}

func createTempFile(t *testing.T, content string) string {
	t.Helper()

	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_task.md")

	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	return tmpFile
}
