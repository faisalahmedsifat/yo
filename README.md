# yo CLI

A productivity CLI that enforces the **RED/YELLOW/GREEN** framework through observation, tracking, and helpful nudges.

```
🔴 RED LIGHT    → Define the problem before coding
🟡 YELLOW LIGHT → Analyze and plan your solution  
🟢 GREEN LIGHT  → Execute with a timer
```

## Quick Start

```bash
# Build and install
make build
./yo init

# Start a task
yo red -i      # Define the problem (interactive)
yo yellow -i   # Analyze and plan
yo go --time 2h   # Start execution

# Track progress
yo status      # Current state
yo timer       # Timer only

# Complete
yo done        # Mark task complete
yo off         # End session
```

---

## System Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                           yo CLI                                     │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────┐  ┌───────────┐  │
│  │   Commands   │  │  File Watch  │  │   Timer   │  │  Activity │  │
│  │   (Cobra)    │  │  (fsnotify)  │  │           │  │    Log    │  │
│  └──────┬───────┘  └──────┬───────┘  └─────┬─────┘  └─────┬─────┘  │
│         │                 │                │              │         │
│         └─────────────────┴────────────────┴──────────────┘         │
│                                    │                                 │
│                        ┌───────────▼───────────┐                    │
│                        │    State Manager      │                    │
│                        │   (JSON + JSONL)      │                    │
│                        └───────────┬───────────┘                    │
│                                    │                                 │
│         ┌──────────────────────────┼──────────────────────────┐     │
│         │                          │                          │     │
│  ┌──────▼──────┐  ┌───────────────▼───────────────┐  ┌───────▼───┐ │
│  │  Project    │  │       Global Daemon           │  │  Stats &  │ │
│  │  .yo/       │  │       ~/.yo/                  │  │  Reports  │ │
│  │             │  │                               │  │           │ │
│  │ state.json  │  │  watcher_config.json          │  │ stats/    │ │
│  │ activity.   │  │  watcher.pid                  │  │ sessions/ │ │
│  │   jsonl     │  │                               │  │           │ │
│  └─────────────┘  └───────────────────────────────┘  └───────────┘ │
└─────────────────────────────────────────────────────────────────────┘
```

### Directory Structure

```
Project Directory            Global (~/.yo/)
├── .yo/                     ├── watcher_config.json  # Watch dirs, current project
│   ├── current_task.md      ├── watcher.pid          # Daemon PID
│   ├── backlog.md           
│   ├── tech_debt_log.md     
│   ├── state.json           
│   ├── config.json          
│   ├── activity.jsonl       
│   ├── done/                
│   │   └── 2024-12-27_task.md
│   ├── sessions/            
│   │   └── 2024-12-27_session.json
│   └── stats/               
│       └── 2024-week-52.json
```

### File Watcher Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                  Global Watcher Daemon                       │
│                  (single process for all projects)           │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│   watch_dirs: [~/Dev, ~/work]                               │
│                     │                                        │
│   ┌─────────────────▼─────────────────┐                     │
│   │        Git Repo Discovery          │                     │
│   │   (finds all .git directories)     │                     │
│   └─────────────────┬─────────────────┘                     │
│                     │                                        │
│   ┌─────────────────▼─────────────────┐                     │
│   │      fsnotify Watchers             │                     │
│   │  (watches all repos recursively)   │                     │
│   └─────────────────┬─────────────────┘                     │
│                     │                                        │
│   ┌─────────────────▼─────────────────┐                     │
│   │      Event Handler                 │                     │
│   │  - Debounce (1 sec)                │                     │
│   │  - Find repo for file              │                     │
│   │  - Check if current task repo      │                     │
│   │  - Log to activity.jsonl           │                     │
│   └─────────────────┬─────────────────┘                     │
│                     │                                        │
│   Active Project: ~/Dev/myapp/.yo/activity.jsonl            │
│   Logs: {"repo": "~/Dev/other", "untracked": true}          │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## Commands

### Setup
```bash
yo init              # Create .yo/ workspace
yo version           # Show version info
yo status            # Current state, timer, session
```

### Workflow (RED → YELLOW → GREEN)
```bash
yo red [-i]          # Define problem (-i for interactive)
yo verify red        # Validate RED complete
yo yellow [-i]       # Analyze and plan
yo verify yellow     # Validate YELLOW complete
yo go [--time 4h]    # Start GREEN + timer
```

### Timer
```bash
yo timer             # Show elapsed / threshold
yo extend 1h "reason" # Add time to threshold
```

### Completion
```bash
yo done              # Complete task (verify criteria)
yo off               # End work session
yo next              # Pick next task from backlog
```

### File Watcher
```bash
yo watch             # Start in foreground
yo watch --bg        # Start in background
yo watch stop        # Stop daemon
yo watch status      # Check if running
```

### Activity & Focus
```bash
yo activity          # Today's activity
yo activity --week   # This week
yo focus             # Focus score
```

### Backlog
```bash
yo list              # Show backlog
yo list --p0         # P0 items only
yo add "description" # Add item
yo add -i            # Interactive add
```

### Emergency
```bash
yo bypass "reason"   # Skip framework (tracked, limited)
```

### Analytics
```bash
yo stats             # Weekly stats
yo stats --week 2024-12-20
yo milestone         # Milestone progress
yo milestone complete
```

### Configuration
```bash
yo config list       # Show config
yo config set editor vim
yo config get editor
```

---

## The RED/YELLOW/GREEN Framework

### 🔴 RED LIGHT - Problem Definition
**Before writing code, define:**
- What is the problem?
- What is the impact?
- What is the severity (P0-P3)?

### 🟡 YELLOW LIGHT - Analysis & Planning
**Before implementing, plan:**
- Root cause analysis (3 levels deep)
- 3+ solution options with time estimates
- Decision with rationale
- Success criteria (minimum 2)

### 🟢 GREEN LIGHT - Execution
**Now code with:**
- Timer running (count-up from 0)
- Focus tracking (via file watcher)
- Threshold notifications (100%, 150%, 200%)

---

## Installation

### From Source
```bash
git clone https://github.com/faisal/yo.git
cd yo
make build
make install  # Installs to $GOPATH/bin
```

### Cross-Platform Builds
```bash
make build-all  # Builds for linux, darwin (amd64, arm64)
```

---

## Testing

```bash
# All tests
make test

# Unit tests only
make test-unit

# Integration tests
make test-integration

# With coverage
go test ./... -cover
```

### Test Coverage
| Package | Coverage |
|---------|----------|
| internal/activity | 72.3% |
| internal/state | 68.9% |
| internal/task | 56.1% |

---

## Configuration

### Project Config (`.yo/config.json`)
```json
{
  "watch_dirs": ["~/Dev"],
  "notifications": true,
  "editor": "vim",
  "max_bypass_day": 1,
  "max_bypass_week": 5
}
```

### Global Watcher Config (`~/.yo/watcher_config.json`)
```json
{
  "watch_dirs": ["~/Dev", "~/work"],
  "current_dir": "/home/user/Dev/myapp"
}
```

---

## State Files

### `state.json`
```json
{
  "version": "1.0.0",
  "current_stage": "green",
  "current_task_id": "deploy_button",
  "timer": {
    "started_at": "2024-12-27T14:30:00Z",
    "estimated_hours": 4.0,
    "threshold_hours": 5.0,
    "extensions": [{"hours": 1.0, "reason": "API complex"}]
  },
  "session": {"active": true, "started_at": "..."},
  "emergency_bypasses": {"today": 0, "this_week": 1}
}
```

### `activity.jsonl`
```jsonl
{"ts":"...","type":"stage_change","from":"yellow","to":"green","task":"deploy"}
{"ts":"...","type":"file_change","repo":"~/Dev/app","file":"main.go"}
{"ts":"...","type":"task_complete","task":"deploy","actual_hours":4.5}
```

---

## License

MIT
