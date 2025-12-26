package task

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ValidationError represents a validation failure
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// ValidationResult holds the result of validating a task file
type ValidationResult struct {
	Valid  bool
	Errors []ValidationError
}

// ValidateRed validates the RED LIGHT section
func ValidateRed(filepath string) (*ValidationResult, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read task file: %w", err)
	}

	result := &ValidationResult{Valid: true}
	text := string(content)

	// Check for problem description
	if !strings.Contains(text, "## 🔴 RED LIGHT") {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "RED LIGHT section",
			Message: "missing RED LIGHT section",
		})
		return result, nil
	}

	// Extract RED LIGHT section
	redSection := extractSection(text, "## 🔴 RED LIGHT", "## 🟡 YELLOW LIGHT")

	// Check for problem statement (not just template)
	if strings.Contains(redSection, "<!-- Describe the problem clearly -->") &&
		!hasContentAfter(redSection, "### What's the Problem?", "<!-- Describe the problem clearly -->") {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "Problem",
			Message: "problem description is empty",
		})
	}

	// Check for at least one impact checkbox
	impactChecked := countCheckedBoxes(redSection, "### Impact", "### Severity")
	if impactChecked == 0 {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "Impact",
			Message: "at least one impact must be selected",
		})
	}

	// Check for severity selection
	severityChecked := countCheckedBoxes(redSection, "### Severity", "---")
	if severityChecked == 0 {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "Severity",
			Message: "severity must be selected (P0-P3)",
		})
	}

	return result, nil
}

// ValidateYellow validates the YELLOW LIGHT section
func ValidateYellow(filepath string) (*ValidationResult, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read task file: %w", err)
	}

	result := &ValidationResult{Valid: true}
	text := string(content)

	// Check for YELLOW LIGHT section
	if !strings.Contains(text, "## 🟡 YELLOW LIGHT") {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "YELLOW LIGHT section",
			Message: "missing YELLOW LIGHT section",
		})
		return result, nil
	}

	// Extract YELLOW LIGHT section
	yellowSection := extractSection(text, "## 🟡 YELLOW LIGHT", "## 🟢 GREEN LIGHT")

	// Check for at least 3 solution options
	optionCount := strings.Count(yellowSection, "#### Option")
	if optionCount < 3 {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "Solution Options",
			Message: fmt.Sprintf("need at least 3 solution options, found %d", optionCount),
		})
	}

	// Check for decision
	if !strings.Contains(yellowSection, "**Chosen option:**") ||
		strings.Contains(yellowSection, "**Chosen option:** \n") {
		hasDecision := false
		lines := strings.Split(yellowSection, "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "**Chosen option:**") && len(strings.TrimPrefix(line, "**Chosen option:**")) > 1 {
				hasDecision = true
				break
			}
		}
		if !hasDecision {
			result.Valid = false
			result.Errors = append(result.Errors, ValidationError{
				Field:   "Decision",
				Message: "chosen option must be specified",
			})
		}
	}

	// Check for success criteria (at least 1)
	successChecked := countAllBoxes(yellowSection, "### Success Criteria", "---")
	if successChecked < 1 {
		result.Valid = false
		result.Errors = append(result.Errors, ValidationError{
			Field:   "Success Criteria",
			Message: fmt.Sprintf("need at least 1 success criterion, found %d", successChecked),
		})
	}

	return result, nil
}

// GetSuccessCriteria extracts success criteria from the task file
func GetSuccessCriteria(filepath string) ([]string, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	text := string(content)
	yellowSection := extractSection(text, "## 🟡 YELLOW LIGHT", "## 🟢 GREEN LIGHT")
	criteriaSection := extractSection(yellowSection, "### Success Criteria", "---")

	var criteria []string
	lines := strings.Split(criteriaSection, "\n")
	checkboxRe := regexp.MustCompile(`^\s*-\s*\[[ xX]\]\s*(.+)$`)

	for _, line := range lines {
		matches := checkboxRe.FindStringSubmatch(line)
		if len(matches) > 1 {
			criteria = append(criteria, strings.TrimSpace(matches[1]))
		}
	}

	return criteria, nil
}

// GetTimeEstimate extracts the time estimate from YELLOW LIGHT
func GetTimeEstimate(filepath string) (float64, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return 0, err
	}

	text := string(content)

	// Look for "Time estimate:" patterns with hours
	timeReH := regexp.MustCompile(`[Tt]ime\s*estimate[:\s]+(\d+(?:\.\d+)?)\s*h`)
	matches := timeReH.FindStringSubmatch(text)
	if len(matches) > 1 {
		var hours float64
		fmt.Sscanf(matches[1], "%f", &hours)
		return hours, nil
	}

	// Look for "Time estimate:" patterns with minutes
	timeReM := regexp.MustCompile(`[Tt]ime\s*estimate[:\s]+(\d+(?:\.\d+)?)\s*m`)
	matches = timeReM.FindStringSubmatch(text)
	if len(matches) > 1 {
		var minutes float64
		fmt.Sscanf(matches[1], "%f", &minutes)
		return minutes / 60, nil
	}

	// Try standalone "Xh" format
	hoursRe := regexp.MustCompile(`(\d+(?:\.\d+)?)\s*(?:hours?|h)\b`)
	matches = hoursRe.FindStringSubmatch(text)
	if len(matches) > 1 {
		var hours float64
		fmt.Sscanf(matches[1], "%f", &hours)
		return hours, nil
	}

	// Try standalone "Xm" format
	minsRe := regexp.MustCompile(`(\d+(?:\.\d+)?)\s*m\b`)
	matches = minsRe.FindStringSubmatch(text)
	if len(matches) > 1 {
		var minutes float64
		fmt.Sscanf(matches[1], "%f", &minutes)
		return minutes / 60, nil
	}

	return 0, fmt.Errorf("no time estimate found")
}

// Helper functions

func extractSection(text, start, end string) string {
	startIdx := strings.Index(text, start)
	if startIdx == -1 {
		return ""
	}
	endIdx := strings.Index(text[startIdx+len(start):], end)
	if endIdx == -1 {
		return text[startIdx:]
	}
	return text[startIdx : startIdx+len(start)+endIdx]
}

func countCheckedBoxes(text, start, end string) int {
	section := extractSection(text, start, end)
	return strings.Count(section, "[x]") + strings.Count(section, "[X]")
}

func countAllBoxes(text, start, end string) int {
	section := extractSection(text, start, end)
	return strings.Count(section, "[ ]") + strings.Count(section, "[x]") + strings.Count(section, "[X]")
}

func hasContentAfter(text, marker, template string) bool {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		if strings.Contains(line, marker) && i+1 < len(lines) {
			nextLine := strings.TrimSpace(lines[i+1])
			if nextLine != "" && !strings.Contains(nextLine, template) && !strings.HasPrefix(nextLine, "<!--") {
				return true
			}
		}
	}
	return false
}

// OpenInEditor opens the file in the user's editor
func OpenInEditor(filepath string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = "nano" // Default fallback
	}

	// Use exec to run the editor
	cmd := fmt.Sprintf("%s %s", editor, filepath)
	return runShellCommand(cmd)
}

func runShellCommand(cmd string) error {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/sh"
	}

	proc := &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}

	argv := []string{shell, "-c", cmd}
	process, err := os.StartProcess(shell, argv, proc)
	if err != nil {
		return err
	}

	_, err = process.Wait()
	return err
}

// PromptConfirm asks for yes/no confirmation
func PromptConfirm(question string) bool {
	fmt.Printf("%s (y/n): ", question)
	reader := bufio.NewReader(os.Stdin)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}
