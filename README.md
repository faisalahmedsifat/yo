# yo

**Stop coding before you think.** `yo` is a CLI that enforces the RED/YELLOW/GREEN framework to help you solve problems methodically.

```
🔴 RED    → What's the problem? (Don't touch code yet)
🟡 YELLOW → What's the plan? (Think before you type)
🟢 GREEN  → Execute with a timer
```

## Why yo?

**The problem:** You see a bug, jump into code, realize 2 hours later you solved the wrong thing.

**The solution:** Force yourself to define the problem and plan before coding. Track your time. Log what you're deferring.

**yo enforces this by:**
- Blocking you from "starting work" until you've defined the problem
- Making you consider 3 options before picking one
- Timing your execution phase
- Tracking tech debt (shortcuts you consciously choose)

---

## Installation

```bash
git clone https://github.com/faisalahmedsifat/yo.git
cd yo
go install
```

---

## Complete Workflow

### 1. Initialize a workspace

```bash
cd your-project
yo init
```

Creates `.yo/` with:
- `current_task.md` - Your RED/YELLOW/GREEN task
- `backlog.md` - Prioritized task list
- `tech_debt_log.md` - Conscious shortcuts you're taking

### 2. Manage your backlog

```bash
# Add tasks
yo add "Fix login button on mobile"
yo add -i                    # Interactive (choose priority)

# View backlog
yo list                      # All items
yo list --p0                 # P0 (launch blockers) only

# Pick next task → starts RED LIGHT
yo next
```

### 3. RED LIGHT - Define the problem

```bash
yo red -i     # Interactive mode (recommended)
```

**You must answer:**
- What's the problem?
- What's the impact?
- What's the severity (P0-P3)?

```bash
yo verify red   # Check if complete
```

### 4. YELLOW LIGHT - Plan your solution

```bash
yo yellow -i   # Interactive mode
```

**You must define:**
- Root cause (3 levels deep - immediate, underlying, system)
- 3 solution options with time estimates
- Your decision and why
- Success criteria

**Log tech debt here:**
```bash
yo defer "No OAuth - using password only for MVP"
yo defer -i   # Interactive with more details
```

```bash
yo verify yellow   # Check if complete
```

### 5. GREEN LIGHT - Execute

```bash
yo go                # Uses estimated time from YELLOW
yo go --time 2h      # Override time estimate
```

**Now you can code.** Timer is running.

```bash
yo status   # Current state + timer
yo timer    # Timer only

# Need more time?
yo extend 30m "API was more complex"
```

### 6. Complete the task

```bash
yo done
```

Prompts you to verify success criteria were met. Archives the task.

### 7. End your session

```bash
yo off   # Ends session, shows focus score
```

---

## Tech Debt Tracking

Tech debt is **not** timer extensions. Tech debt is **conscious decisions to defer work**.

Use `yo defer` during YELLOW LIGHT when you're choosing to skip something:

```bash
yo defer "No retry button - users can click again"
yo defer -i   # Interactive mode
```

**Creates entries like:**
```markdown
## Deferred on 2024-12-27
**Task:** deploy_feature

**What:** No retry button
**Why skipped:** Users can click deploy again, not critical
**Come back when:** When user complains
**Estimated fix time:** 2h
```

View with: `cat .yo/tech_debt_log.md`

---

## Emergency Bypass

Sometimes you need to skip the framework:

```bash
yo bypass "production is down"
```

**Limited to:** 1/day, 5/week. Tracked for accountability.

---

## Activity & Stats

```bash
yo activity          # Today's changes
yo activity --week   # This week

yo focus             # Focus score (on-task vs off-task)
yo stats             # Weekly statistics
yo milestone         # Track your progress
```

---

## File Watcher (Optional)

Track file changes across repos:

```bash
yo watch             # Foreground
yo watch --bg        # Background daemon
yo watch stop        # Stop daemon
yo watch status      # Check if running
```

Logs when you're working on the wrong repo (off-task activity).

---

## All Commands

| Command | Description |
|---------|-------------|
| `yo init` | Initialize workspace |
| `yo status` | Show current state |
| `yo red -i` | Define problem (interactive) |
| `yo yellow -i` | Plan solution (interactive) |
| `yo go` | Start GREEN LIGHT with timer |
| `yo timer` | Show timer |
| `yo extend 1h` | Add time to threshold |
| `yo done` | Complete task |
| `yo off` | End session |
| `yo list` | Show backlog |
| `yo add "task"` | Add to backlog |
| `yo next` | Pick next task |
| `yo defer "what"` | Log tech debt |
| `yo bypass "why"` | Emergency skip |
| `yo activity` | Show activity |
| `yo focus` | Show focus score |
| `yo stats` | Weekly stats |
| `yo milestone` | Progress milestones |
| `yo config list` | Show config |
| `yo watch` | Start file watcher |

---

## Configuration

```bash
yo config list              # Show all settings
yo config set editor vim    # Set editor
yo config get editor        # Get value
```

Settings in `.yo/config.json`:
- `notifications` - Desktop notifications (true/false)
- `editor` - Editor for opening files
- `max_bypass_day` - Bypass limit per day (default: 1)
- `max_bypass_week` - Bypass limit per week (default: 5)

---

## Directory Structure

```
your-project/
└── .yo/
    ├── current_task.md    # Current RED/YELLOW/GREEN task
    ├── backlog.md         # Prioritized backlog
    ├── tech_debt_log.md   # Conscious shortcuts
    ├── state.json         # Timer, stage, session
    ├── config.json        # Settings
    ├── activity.jsonl     # Activity log
    ├── done/              # Archived completed tasks
    ├── sessions/          # Session summaries
    └── stats/             # Weekly statistics
```

---

## The Framework Mindset

### Before yo:
1. See bug
2. Jump into code
3. Realize 2 hours later you fixed the wrong thing
4. Feel bad

### After yo:
1. See bug → `yo red -i`
2. Define exact problem, impact, severity
3. Analyze root cause → `yo yellow -i`  
4. Consider 3 options, pick one, define success criteria
5. Log what you're skipping → `yo defer`
6. Start timer → `yo go`
7. Code with focus
8. Verify success criteria → `yo done`

**The few minutes spent planning save hours of wasted work.**

---

## License

MIT
