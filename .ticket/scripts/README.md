# Ticket Management Scripts

This directory contains scripts for managing SilentCast development tickets.

## 🚀 Quick Start

Use the unified `ticket.sh` script for all ticket operations:

```bash
# Show quick summary
./.ticket/scripts/ticket.sh

# Create a new ticket
./.ticket/scripts/ticket.sh new --title "My new feature" --type feature

# List all tickets
./.ticket/scripts/ticket.sh list

# Update ticket status
./.ticket/scripts/ticket.sh status T001 in_progress

# Show specific ticket
./.ticket/scripts/ticket.sh show T001

# Edit a ticket
./.ticket/scripts/ticket.sh edit T001

# Generate report
./.ticket/scripts/ticket.sh report
```

## 📝 Available Commands

### Main Script: `ticket.sh`

The unified interface for all ticket operations:

- `new` (or `n`) - Create a new ticket
- `list` (or `ls`) - List tickets with filters
- `status` (or `s`) - Update ticket status
- `report` (or `r`) - Generate status report
- `show` - Display a specific ticket
- `edit` - Edit a ticket in your $EDITOR
- `help` (or `h`) - Show help

### Individual Scripts

You can also use the individual scripts directly:

- `ticket-new.sh` - Create new tickets
- `ticket-list.sh` - List and filter tickets
- `ticket-status.sh` - Update ticket status
- `ticket-report.sh` - Generate reports

## 💡 Tips

1. **Quick status check**: Just run `ticket.sh` with no arguments
2. **Use shortcuts**: `n` for new, `ls` for list, `s` for status, `r` for report
3. **Filter tickets**: `ticket.sh ls --status todo --priority high`
4. **Export reports**: `ticket.sh report > STATUS.md`

## 🔧 Configuration

- Tickets are stored in `.ticket/tickets/`
- Completed tickets move to `.ticket/tickets/archive/{year}/{month}/`
- Templates are in `.ticket/meta/templates/`
- Configuration is in `.ticket/meta/config.yml`