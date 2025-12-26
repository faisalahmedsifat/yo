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
