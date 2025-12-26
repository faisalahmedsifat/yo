package templates

// CurrentTask is the template for .yo/current_task.md
const CurrentTask = `# Current Task

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

### Root Cause Analysis
**Immediate cause:**
<!-- What directly causes this? -->

**Underlying cause:**
<!-- Why does the immediate cause exist? -->

**System cause:**
<!-- What systemic issue allows this? -->

### Solution Options

#### Option A:
- Description:
- Time estimate:
- Pros:
- Cons:

#### Option B:
- Description:
- Time estimate:
- Pros:
- Cons:

#### Option C:
- Description:
- Time estimate:
- Pros:
- Cons:

### Decision
**Chosen option:** 
**Reason:**

### What I'm Deferring (Tech Debt)
<!-- Things you know should be done but are choosing to skip for now -->
<!-- Use 'yo defer -i' to log these properly, or list them here: -->
- 

### Implementation Steps
1. 
2. 
3. 

### Success Criteria
- [ ] 
- [ ] 
- [ ] 

---

## 🟢 GREEN LIGHT - Execution

### Timer Started:
### Estimated Time:

### Notes:
<!-- Add notes as you work -->

### Blockers:
<!-- Document any blockers -->

---

## ✅ Completion

### Actual Time:
### Accuracy:
### Lessons Learned:
`

// Backlog is the template for .yo/backlog.md
const Backlog = `# Backlog

## P0 - Launch Blockers


## P1 - Paying User Blockers


## P2 - Nice to Have


## P3 - Future Improvements

`

// TechDebtLog is the template for .yo/tech_debt_log.md
const TechDebtLog = `# Tech Debt Log

Tech debt = shortcuts you CHOSE to take to ship faster.
Log decisions you're consciously deferring during YELLOW LIGHT.

Use: yo defer -i

---

`

// SessionSummary is the template for session summaries
const SessionSummary = `# Session Summary

Date: {{.Date}}
Started: {{.StartedAt}}
Ended: {{.EndedAt}}
Duration: {{.Duration}}

## Focus Score
{{.FocusScore}}%

## Activity
{{.Activity}}

## Tasks
{{.Tasks}}
`

// Agents is the template for .yo/AGENTS.md
const Agents = `# AGENTS.md - AI Agent Instructions

This file provides instructions for AI agents working on this project.

## Overview

This project uses the **yo** CLI to enforce the RED/YELLOW/GREEN methodology.
As an AI agent, you MUST follow this workflow when making changes.

## The Workflow

` + "```" + `
🔴 RED    → Define the problem (no code changes yet)
🟡 YELLOW → Plan the solution (analyze, consider options)
🟢 GREEN  → Execute with a timer (now you can write code)
` + "```" + `

## File Reference

| File | Purpose | When to Update |
|------|---------|----------------|
| ` + "`" + `.yo/current_task.md` + "`" + ` | Active task definition | During RED/YELLOW phases |
| ` + "`" + `.yo/backlog.md` + "`" + ` | Prioritized task list | When adding/completing tasks |
| ` + "`" + `.yo/tech_debt_log.md` + "`" + ` | Conscious shortcuts | When deferring work |
| ` + "`" + `.yo/state.json` + "`" + ` | Current stage/timer | Managed by yo CLI |
| ` + "`" + `.yo/activity.jsonl` + "`" + ` | Activity log | Managed by yo CLI |

## Commands Reference

### Stage Transitions
- ` + "`" + `yo red` + "`" + ` - Start/update problem definition
- ` + "`" + `yo yellow` + "`" + ` - Start/update planning phase
- ` + "`" + `yo go --time 2h` + "`" + ` - Start execution with timer
- ` + "`" + `yo done` + "`" + ` - Complete current task

### Verification
- ` + "`" + `yo verify red` + "`" + ` - Check if RED phase is complete
- ` + "`" + `yo verify yellow` + "`" + ` - Check if YELLOW phase is complete
- ` + "`" + `yo status` + "`" + ` - Show current stage and timer

### Backlog Management
- ` + "`" + `yo list` + "`" + ` - View all backlog items
- ` + "`" + `yo list --p0` + "`" + ` - View P0 (launch blockers) only
- ` + "`" + `yo add "description"` + "`" + ` - Add item to backlog
- ` + "`" + `yo next` + "`" + ` - Pick next task from backlog

### Tech Debt
- ` + "`" + `yo defer "what I'm skipping"` + "`" + ` - Log a conscious shortcut

## Agent Instructions

### Before Making Changes
1. Run ` + "`" + `yo status` + "`" + ` to check current state
2. If in GREEN stage, proceed with code changes
3. If not in GREEN stage, complete RED and YELLOW first

### Starting New Work
1. Check ` + "`" + `yo list` + "`" + ` for existing backlog items
2. Use ` + "`" + `yo next` + "`" + ` to pick a task, OR
3. Start with ` + "`" + `yo red` + "`" + ` for a new problem

### During RED Phase
Fill in ` + "`" + `.yo/current_task.md` + "`" + `:
- **What's the Problem?** - Clear problem statement
- **Impact** - Check applicable boxes
- **Severity** - P0/P1/P2/P3

### During YELLOW Phase
Complete the planning section:
- **Root Cause Analysis** - 3 levels deep
- **Solution Options** - At least 2-3 options with estimates
- **Decision** - Chosen option with reasoning
- **Success Criteria** - Testable criteria

### During GREEN Phase
- Execute the plan
- Run ` + "`" + `yo timer` + "`" + ` to check remaining time
- If blocked, document in the Blockers section
- If deferring work, use ` + "`" + `yo defer` + "`" + `

### Completing Work
1. Verify success criteria are met
2. Run ` + "`" + `yo done` + "`" + `
3. Task is archived to ` + "`" + `.yo/done/` + "`" + `

## Priority Levels

| Level | Description | SLA |
|-------|-------------|-----|
| P0 | Launch blocker | Fix immediately |
| P1 | Paying user blocker | Fix this sprint |
| P2 | Nice to have | Backlog |
| P3 | Future improvement | Long-term |

## Best Practices

1. **Never skip RED** - Define problems before coding
2. **Consider options** - Don't jump to first solution
3. **Track tech debt** - Use ` + "`" + `yo defer` + "`" + ` for shortcuts
4. **Stay focused** - Work on current task only
5. **Time-box** - Respect timer estimates
`
