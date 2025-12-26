package milestone

import (
	"github.com/faisal/yo/internal/state"
)

// Milestone represents a framework milestone
type Milestone struct {
	ID       int
	Name     string
	Criteria []string
}

// All returns all defined milestones
func All() []Milestone {
	return []Milestone{
		{
			ID:   0,
			Name: "Clear Launch Blockers",
			Criteria: []string{
				"Can sign up with fresh account",
				"Load at least 1 template",
				"Deploy successfully",
				"See success/failure status",
				"Repeatable 3 times",
			},
		},
		{
			ID:   1,
			Name: "Core Workflow Complete",
			Criteria: []string{
				"Full RED/YELLOW/GREEN cycle works",
				"Timer tracks accurately",
				"Tasks archive correctly",
				"State persists across sessions",
			},
		},
		{
			ID:   2,
			Name: "Activity Tracking",
			Criteria: []string{
				"File watcher logs changes",
				"Focus score calculates correctly",
				"Activity queries work",
			},
		},
		{
			ID:   3,
			Name: "Full Feature Set",
			Criteria: []string{
				"Backlog management works",
				"Stats calculate weekly",
				"Emergency bypass tracked",
			},
		},
		{
			ID:   4,
			Name: "Production Ready",
			Criteria: []string{
				"All tests passing",
				"Cross-platform builds work",
				"Documentation complete",
				"Dogfooded for 1 week",
			},
		},
	}
}

// Current returns the current milestone
func Current(s *state.State) *Milestone {
	milestones := All()
	if s.Milestone.Current >= len(milestones) {
		return nil
	}
	m := milestones[s.Milestone.Current]
	return &m
}

// Progress calculates completion percentage for the current milestone
func Progress(s *state.State, completed int) float64 {
	m := Current(s)
	if m == nil || len(m.Criteria) == 0 {
		return 100
	}
	return (float64(completed) / float64(len(m.Criteria))) * 100
}

// Advance moves to the next milestone
func Advance(s *state.State) bool {
	milestones := All()
	if s.Milestone.Current >= len(milestones)-1 {
		return false
	}

	s.Milestone.Current++
	if s.Milestone.Current < len(milestones) {
		s.Milestone.Name = milestones[s.Milestone.Current].Name
	} else {
		s.Milestone.Name = "All Complete"
	}

	return true
}

// IsComplete checks if all milestones are complete
func IsComplete(s *state.State) bool {
	return s.Milestone.Current >= len(All())
}

// GetByID returns a milestone by ID
func GetByID(id int) *Milestone {
	milestones := All()
	if id < 0 || id >= len(milestones) {
		return nil
	}
	m := milestones[id]
	return &m
}

// Status represents the status of all milestones
type Status struct {
	Current   int
	Total     int
	Completed []bool
	Names     []string
}

// GetStatus returns the overall milestone status
func GetStatus(s *state.State) *Status {
	milestones := All()
	status := &Status{
		Current:   s.Milestone.Current,
		Total:     len(milestones),
		Completed: make([]bool, len(milestones)),
		Names:     make([]string, len(milestones)),
	}

	for i, m := range milestones {
		status.Names[i] = m.Name
		status.Completed[i] = i < s.Milestone.Current
	}

	return status
}
