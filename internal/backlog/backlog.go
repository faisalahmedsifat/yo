package backlog

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/faisal/yo/internal/workspace"
)

// Priority levels
const (
	P0 = "P0"
	P1 = "P1"
	P2 = "P2"
	P3 = "P3"
)

// Item represents a backlog item
type Item struct {
	Text     string
	Priority string
	Checked  bool
	Line     int
}

// Backlog holds all backlog items
type Backlog struct {
	Items map[string][]Item
	Path  string
}

// Load loads the backlog from disk
func Load() (*Backlog, error) {
	path, err := workspace.GetBacklogPath()
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read backlog: %w", err)
	}

	return Parse(string(content), path), nil
}

// Parse parses backlog markdown content
func Parse(content, path string) *Backlog {
	b := &Backlog{
		Items: make(map[string][]Item),
		Path:  path,
	}
	b.Items[P0] = []Item{}
	b.Items[P1] = []Item{}
	b.Items[P2] = []Item{}
	b.Items[P3] = []Item{}

	lines := strings.Split(content, "\n")
	currentPriority := ""
	itemRe := regexp.MustCompile(`^\s*-\s*\[([ xX])\]\s*(.+)$`)

	for lineNum, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Detect priority section
		if strings.Contains(trimmed, "P0 -") || strings.Contains(trimmed, "## P0") {
			currentPriority = P0
			continue
		}
		if strings.Contains(trimmed, "P1 -") || strings.Contains(trimmed, "## P1") {
			currentPriority = P1
			continue
		}
		if strings.Contains(trimmed, "P2 -") || strings.Contains(trimmed, "## P2") {
			currentPriority = P2
			continue
		}
		if strings.Contains(trimmed, "P3 -") || strings.Contains(trimmed, "## P3") {
			currentPriority = P3
			continue
		}

		// Parse items
		if currentPriority != "" {
			matches := itemRe.FindStringSubmatch(line)
			if len(matches) > 2 {
				item := Item{
					Text:     strings.TrimSpace(matches[2]),
					Priority: currentPriority,
					Checked:  matches[1] == "x" || matches[1] == "X",
					Line:     lineNum,
				}
				b.Items[currentPriority] = append(b.Items[currentPriority], item)
			}
		}
	}

	return b
}

// Add adds an item to the backlog
func (b *Backlog) Add(text, priority string) error {
	content, err := os.ReadFile(b.Path)
	if err != nil {
		return err
	}

	markdown := string(content)
	item := fmt.Sprintf("- [ ] %s\n", text)

	sectionMarkers := map[string]string{
		P0: "## P0 - Launch Blockers",
		P1: "## P1 - Paying User Blockers",
		P2: "## P2 - Nice to Have",
		P3: "## P3 - Future Improvements",
	}

	marker := sectionMarkers[priority]
	if idx := strings.Index(markdown, marker); idx != -1 {
		endOfLine := idx + len(marker)
		for endOfLine < len(markdown) && markdown[endOfLine] != '\n' {
			endOfLine++
		}
		if endOfLine < len(markdown) {
			endOfLine++
		}
		markdown = markdown[:endOfLine] + item + markdown[endOfLine:]
	} else {
		markdown += fmt.Sprintf("\n%s\n%s", marker, item)
	}

	return os.WriteFile(b.Path, []byte(markdown), 0644)
}

// GetUnchecked returns all unchecked items ordered by priority
func (b *Backlog) GetUnchecked() []Item {
	var items []Item
	for _, p := range []string{P0, P1, P2, P3} {
		for _, item := range b.Items[p] {
			if !item.Checked {
				items = append(items, item)
			}
		}
	}
	return items
}

// Count returns the count of items by priority
func (b *Backlog) Count() map[string]int {
	counts := make(map[string]int)
	for p, items := range b.Items {
		counts[p] = len(items)
	}
	return counts
}

// CountUnchecked returns the count of unchecked items by priority
func (b *Backlog) CountUnchecked() map[string]int {
	counts := make(map[string]int)
	for p, items := range b.Items {
		for _, item := range items {
			if !item.Checked {
				counts[p]++
			}
		}
	}
	return counts
}

// Total returns total item count
func (b *Backlog) Total() int {
	total := 0
	for _, items := range b.Items {
		total += len(items)
	}
	return total
}

// TotalUnchecked returns total unchecked item count
func (b *Backlog) TotalUnchecked() int {
	total := 0
	for _, items := range b.Items {
		for _, item := range items {
			if !item.Checked {
				total++
			}
		}
	}
	return total
}
