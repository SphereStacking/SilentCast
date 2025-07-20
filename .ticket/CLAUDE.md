# Claude Guide for Ticket System

This guide helps Claude understand and effectively use the ticket management system.

## 🎯 System Overview

The ticket system is a Git-friendly task management solution that uses YAML files for tracking development tasks. Each ticket is a separate file, making it easy to track changes and avoid merge conflicts.

## 📁 Structure

```
.ticket/
├── meta/               # Configuration and templates
│   ├── config.yml      # System configuration (next_id, etc.)
│   ├── labels.yml      # Label definitions
│   └── templates/      # Ticket templates (default, feature, bug)
├── tickets/            # All active tickets
│   └── archive/        # Completed tickets (organized by year/month)
├── scripts/            # Management scripts
└── CLAUDE.md           # This file
```

## 🛠️ Primary Interface

Always use the unified `ticket` command from the project root:

```bash
./ticket              # Quick status summary
./ticket new          # Create new ticket
./ticket list         # List all tickets
./ticket show T001    # Show ticket details
./ticket edit T001    # Edit ticket in editor
./ticket status T001 in_progress  # Update status
./ticket report       # Generate full report
```

## 📝 Creating Tickets

### Using the script (recommended):
```bash
./ticket new --title "Implement new feature" --type feature --priority high
```

### Ticket fields:
- **id**: Auto-assigned (T001, T002, etc.)
- **title**: Brief description
- **type**: feature, bug, refactor, docs, test, chore
- **priority**: critical, high, medium, low
- **status**: todo, in_progress, review, testing, completed, blocked
- **labels**: Use for categorization (e.g., foundation, core-features, cli)

## 🔄 Workflow

### 1. Start work on a ticket:
```bash
./ticket status T001 in_progress
./ticket show T001  # Review details
```

### 2. Complete a ticket:
```bash
./ticket status T001 completed
# Automatically moves to tickets/archive/YYYY/MM/
```

### 3. Check progress:
```bash
./ticket list --status in_progress
./ticket report > TICKET_STATUS.md
```

## 🏷️ Label System

Labels replace the old phase system for flexible categorization:
- **foundation**: Core infrastructure
- **core-features**: Main functionality
- **cli**: Command-line interface
- **testing**: Test-related
- **docs**: Documentation
- **performance**: Performance improvements
- **security**: Security-related

## 💡 Best Practices for Claude

### 1. Always check current status first:
```bash
./ticket
./ticket list --status todo --priority high
```

### 2. When implementing features:
- Find related tickets: `./ticket list --label foundation`
- Show ticket details: `./ticket show T003`
- Update status when starting: `./ticket status T003 in_progress`
- Mark completed when done: `./ticket status T003 completed`

### 3. Creating comprehensive tickets:
```bash
./ticket new --title "Add streaming output support" \
            --type feature \
            --priority high
            
# Then edit to add details:
./ticket edit T005
```

### 4. Tracking work:
- Use time_spent field in tickets
- Update notes section with findings
- Link related tickets in the related field

### 5. Reporting:
```bash
# Generate report for user
./ticket report

# Quick status for context
./ticket list --status todo
```

## 🔍 Common Queries

### Find high-priority TODOs:
```bash
./ticket list --status todo --priority high
```

### Find tickets by label:
```bash
./ticket list --label core-features
```

### Show all in-progress work:
```bash
./ticket list --status in_progress
```

### Find specific ticket:
```bash
./ticket show T001
```

## ⚠️ Important Notes

1. **File locations**: All tickets are in `.ticket/tickets/`
2. **Completed tickets**: Automatically moved to `archive/YYYY/MM/`
3. **Never edit**: `meta/config.yml` directly (especially next_id)
4. **Always use scripts**: Don't manually move ticket files
5. **Status transitions**: Use `./ticket status` command

## 🚀 Quick Reference

```bash
# Essential commands
./ticket                          # Status overview
./ticket new --title "..." --type feature --priority high
./ticket list --status todo
./ticket show T001
./ticket status T001 in_progress
./ticket status T001 completed
./ticket report

# Shortcuts
./ticket n     # new
./ticket ls    # list  
./ticket s     # status
./ticket r     # report
```

## 📊 Understanding Ticket Files

Example ticket structure:
```yaml
id: T001
title: "Create OutputManager Interface"
type: feature
priority: high
status: completed
created: "2025-01-17"
updated: "2025-01-17"
assignee: null
labels:
  - core
  - output-management
  - foundation

description: |
  Detailed description...

tasks:
  - [x] Completed task
  - [ ] Pending task

acceptance_criteria:
  - Criteria 1
  - Criteria 2

dependencies: [T002, T003]
related: [T004]
files:
  - app/internal/output/interface.go

time_estimate: 4h
time_spent: 3h

notes: |
  Implementation notes...
```

## 🤝 Integration with Development

When working on project features:

1. **Before implementing**: Check for existing tickets
2. **During implementation**: Update ticket status and notes
3. **After implementation**: Mark completed, create follow-up tickets if needed
4. **For documentation**: Create docs type tickets
5. **For bugs**: Use bug template with reproduction steps

Remember: The ticket system is the source of truth for development tasks!