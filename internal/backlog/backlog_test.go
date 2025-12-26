package backlog

import (
	"testing"
)

func TestParse(t *testing.T) {
	content := `# Backlog

## P0 - Launch Blockers
- [x] Fix login bug
- [ ] Deploy button missing

## P1 - Paying User Blockers
- [ ] Payment flow broken

## P2 - Nice to Have
- [ ] Dark mode

## P3 - Future Improvements
`

	b := Parse(content, "/test/backlog.md")

	if len(b.Items[P0]) != 2 {
		t.Errorf("Expected 2 P0 items, got %d", len(b.Items[P0]))
	}

	if len(b.Items[P1]) != 1 {
		t.Errorf("Expected 1 P1 item, got %d", len(b.Items[P1]))
	}

	if len(b.Items[P2]) != 1 {
		t.Errorf("Expected 1 P2 item, got %d", len(b.Items[P2]))
	}

	if len(b.Items[P3]) != 0 {
		t.Errorf("Expected 0 P3 items, got %d", len(b.Items[P3]))
	}

	// Check first P0 item is checked
	if !b.Items[P0][0].Checked {
		t.Error("Expected first P0 item to be checked")
	}

	// Check second P0 item is not checked
	if b.Items[P0][1].Checked {
		t.Error("Expected second P0 item to be unchecked")
	}
}

func TestGetUnchecked(t *testing.T) {
	content := `# Backlog

## P0 - Launch Blockers
- [x] Done item
- [ ] P0 todo

## P1 - Paying User Blockers
- [ ] P1 todo
`

	b := Parse(content, "/test/backlog.md")
	unchecked := b.GetUnchecked()

	if len(unchecked) != 2 {
		t.Errorf("Expected 2 unchecked items, got %d", len(unchecked))
	}

	// Should be ordered by priority
	if unchecked[0].Priority != P0 {
		t.Error("Expected first unchecked to be P0")
	}

	if unchecked[1].Priority != P1 {
		t.Error("Expected second unchecked to be P1")
	}
}

func TestCount(t *testing.T) {
	content := `# Backlog

## P0 - Launch Blockers
- [ ] Item 1
- [ ] Item 2
- [ ] Item 3

## P1 - Paying User Blockers
- [ ] Item 4
`

	b := Parse(content, "/test/backlog.md")
	counts := b.Count()

	if counts[P0] != 3 {
		t.Errorf("Expected P0 count 3, got %d", counts[P0])
	}

	if counts[P1] != 1 {
		t.Errorf("Expected P1 count 1, got %d", counts[P1])
	}
}

func TestTotal(t *testing.T) {
	content := `# Backlog

## P0 - Launch Blockers
- [ ] Item 1
- [x] Item 2

## P1 - Paying User Blockers
- [ ] Item 3
`

	b := Parse(content, "/test/backlog.md")

	if b.Total() != 3 {
		t.Errorf("Expected total 3, got %d", b.Total())
	}

	if b.TotalUnchecked() != 2 {
		t.Errorf("Expected total unchecked 2, got %d", b.TotalUnchecked())
	}
}

func TestCountUnchecked(t *testing.T) {
	content := `# Backlog

## P0 - Launch Blockers
- [ ] Unchecked 1
- [x] Checked
- [ ] Unchecked 2

## P1 - Paying User Blockers
- [x] All done
`

	b := Parse(content, "/test/backlog.md")
	counts := b.CountUnchecked()

	if counts[P0] != 2 {
		t.Errorf("Expected P0 unchecked count 2, got %d", counts[P0])
	}

	if counts[P1] != 0 {
		t.Errorf("Expected P1 unchecked count 0, got %d", counts[P1])
	}
}

func TestParseEmpty(t *testing.T) {
	content := `# Backlog

## P0 - Launch Blockers

## P1 - Paying User Blockers

## P2 - Nice to Have

## P3 - Future Improvements
`

	b := Parse(content, "/test/backlog.md")

	if b.Total() != 0 {
		t.Errorf("Expected empty backlog, got %d items", b.Total())
	}
}
