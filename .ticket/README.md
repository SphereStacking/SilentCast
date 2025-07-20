# Ticket Management System

This directory contains a Git-friendly ticket management system for tracking development tasks.

## 🚀 Installation

### Quick Install

Run the setup script to initialize the ticket system:

```bash
./.ticket/setup.sh
```

This will:
- Create a `ticket` symlink in the project root
- Verify directory structure
- Check all required files
- Set correct permissions
- Update .gitignore if needed

### Manual Installation

If you prefer to set up manually:

1. **Create symlink**
   ```bash
   ln -s .ticket/scripts/ticket.sh ticket
   ```

2. **Set permissions**
   ```bash
   chmod +x .ticket/scripts/*.sh
   chmod +x .ticket/setup.sh
   ```

3. **Update .gitignore**
   Add the following line:
   ```
   TICKET_STATUS.md
   ```

4. **Verify structure**
   Ensure these directories exist:
   - `.ticket/meta/`
   - `.ticket/meta/templates/`
   - `.ticket/tickets/`
   - `.ticket/tickets/archive/`
   - `.ticket/scripts/`

## 📁 Directory Structure

```
.ticket/
├── meta/               # System configuration and metadata
│   ├── config.yml      # Main configuration
│   ├── labels.yml      # Label definitions
│   └── templates/      # Ticket templates
├── tickets/           # All tickets
│   └── archive/       # Completed tickets
│       └── 2025/01/   # Organized by year/month
└── scripts/           # Management scripts
```

## 🎫 Ticket Format

Each ticket is a YAML file with the following structure:

```yaml
id: T001                    # Unique ticket ID
title: "Ticket Title"       # Brief description
type: feature              # feature|bug|refactor|docs|test|chore
priority: high             # critical|high|medium|low
status: todo               # todo|in_progress|review|testing|completed|blocked
created: "2025-01-17"      # Creation date
updated: "2025-01-17"      # Last update date
assignee: null             # Assigned developer
labels: []                 # List of labels (e.g., foundation, core-features, etc.)

description: |
  Detailed description of the ticket

tasks:
  - [ ] Subtask 1
  - [ ] Subtask 2

acceptance_criteria:
  - Criteria 1
  - Criteria 2

dependencies: []           # Ticket IDs this depends on
related: []               # Related ticket IDs
files: []                 # Files to modify/create
time_estimate: 2h         # Estimated time
time_spent: 0h            # Actual time spent
```

## 📋 Workflow

### Ticket Status Flow

```
todo → in_progress → review → testing → completed
         ↓             ↓         ↓
      blocked      blocked   blocked
```

### Creating a New Ticket

1. Use the script: `.ticket/scripts/ticket-new.sh --title "My Task"`
2. Or manually:
   - Copy template from `meta/templates/`
   - Save as `tickets/T{number}-{slug}.yml`
   - Update next_id in `meta/config.yml`

### Moving Tickets

- All active tickets stay in `tickets/`
- Completed tickets: Automatically moved to `tickets/archive/{year}/{month}/`
- Status changes: Use `.ticket/scripts/ticket-status.sh T001 completed`

## 🏷️ Labels

Labels are defined in `meta/labels.yml` and grouped by category:

- **Components**: core, output-management, notification, config, etc.
- **Platforms**: windows, macos, linux, cross-platform
- **Special**: breaking-change, security, performance, ux
- **Testing**: needs-tests, has-tests, test-failing
- **Documentation**: needs-docs, has-docs

## 📊 Core Fields

- **id**: Unique identifier (T001, T002, etc.)
- **status**: todo, in_progress, review, testing, completed, blocked
- **priority**: critical, high, medium, low
- **type**: feature, bug, refactor, docs, test, chore
- **labels**: Flexible categorization (foundation, core-features, cli, etc.)

## 🛠️ Utility Scripts

Ticket management scripts are located in `.ticket/scripts/`:

### Main Interface
- **`ticket.sh`** - Unified interface for all operations
  ```bash
  ./.ticket/scripts/ticket.sh         # Quick status
  ./.ticket/scripts/ticket.sh new     # Create ticket
  ./.ticket/scripts/ticket.sh list    # List tickets
  ./.ticket/scripts/ticket.sh show T001  # Show ticket
  ```

### Individual Scripts
- `ticket-new.sh` - Create a new ticket
- `ticket-status.sh` - Update ticket status
- `ticket-list.sh` - List tickets with filters
- `ticket-report.sh` - Generate progress reports

See `.ticket/scripts/README.md` for detailed documentation.

## 📚 For AI Assistants (Claude)

If you're an AI assistant working with this codebase, please read `.ticket/CLAUDE.md` for specific guidance on using the ticket system effectively.

## 📝 Guidelines

1. **One ticket = One feature/bug/task**
   - Keep tickets focused and atomic
   - Break large features into multiple tickets

2. **Clear descriptions**
   - Include context and background
   - Define clear acceptance criteria

3. **Regular updates**
   - Update status when starting/completing work
   - Add time spent for tracking

4. **Use dependencies**
   - Link related tickets
   - Mark blockers clearly

5. **File naming**
   - Automatic: `T{number}-{slug}.yml`
   - Keep slugs under 50 characters

6. **Label-based organization**
   - Use labels for categorization instead of phases
   - Examples: foundation, core-features, cli, testing, docs
   - Multiple labels per ticket for flexible organization

## 🔍 Common Queries

### Find all TODO tickets
```bash
grep -l "status: todo" tickets/*.yml
```

### Find high priority tickets
```bash
grep -l "priority: high" */*.yml
```

### Find tickets by label
```bash
grep -l "output-management" */*.yml
```

### Count tickets by status
```bash
grep "status:" */*.yml | cut -d: -f3 | sort | uniq -c
```

## 🚀 Getting Started

### First Steps

1. **Create your first ticket**:
   ```bash
   ./ticket new --title "My first task" --type feature
   ```

2. **List tickets**:
   ```bash
   ./ticket list
   ```

3. **Update status**:
   ```bash
   ./ticket status T001 in_progress
   ```

## 📈 Reporting

Generate a status report:
```bash
./.ticket/scripts/ticket-report.sh > TICKET_STATUS.md
```

This will create a markdown report with:
- Overall summary
- Status distribution
- Priority breakdown
- Recent updates

## 🔧 Troubleshooting

### "Permission denied" error
```bash
chmod +x .ticket/scripts/*.sh
```

### "ticket: command not found"
```bash
# Make sure you're in the project root
cd /path/to/project
./ticket help
```

### Missing directories
```bash
# Run setup to create all directories
./.ticket/setup.sh
```