#!/bin/bash
# ticket.sh - Unified ticket management interface

set -e

# Resolve symlinks to get the actual script location
SCRIPT_PATH="$0"
if [ -L "$SCRIPT_PATH" ]; then
    SCRIPT_PATH="$(readlink -f "$SCRIPT_PATH")"
fi
SCRIPT_DIR="$(cd "$(dirname "$SCRIPT_PATH")" && pwd)"
TICKET_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Show usage
show_usage() {
    echo "Ticket Management System"
    echo ""
    echo "Usage: $0 <command> [options]"
    echo ""
    echo "Commands:"
    echo "  new, n      Create a new ticket"
    echo "  list, ls    List tickets with filters"
    echo "  status, s   Update ticket status"
    echo "  report, r   Generate status report"
    echo "  show        Show a specific ticket"
    echo "  edit        Edit a ticket (opens in \$EDITOR)"
    echo "  help, h     Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 new --title \"Add new feature\" --type feature"
    echo "  $0 list --status todo"
    echo "  $0 status T001 in_progress"
    echo "  $0 report"
    echo "  $0 show T001"
    echo ""
    echo "For command-specific help:"
    echo "  $0 <command> --help"
}

# Show a specific ticket
show_ticket() {
    if [ -z "$1" ]; then
        echo "Error: Ticket ID required"
        echo "Usage: $0 show TICKET_ID"
        exit 1
    fi
    
    TICKET_ID="$1"
    TICKET_FILE=$(find "$TICKET_DIR/tickets" -name "$TICKET_ID-*.yml" -type f 2>/dev/null | head -1)
    
    if [ -z "$TICKET_FILE" ]; then
        echo -e "${RED}Error: Ticket $TICKET_ID not found${NC}"
        exit 1
    fi
    
    echo -e "${BLUE}=== Ticket $TICKET_ID ===${NC}"
    echo ""
    cat "$TICKET_FILE"
}

# Edit a ticket
edit_ticket() {
    if [ -z "$1" ]; then
        echo "Error: Ticket ID required"
        echo "Usage: $0 edit TICKET_ID"
        exit 1
    fi
    
    TICKET_ID="$1"
    TICKET_FILE=$(find "$TICKET_DIR/tickets" -name "$TICKET_ID-*.yml" -type f 2>/dev/null | head -1)
    
    if [ -z "$TICKET_FILE" ]; then
        echo -e "${RED}Error: Ticket $TICKET_ID not found${NC}"
        exit 1
    fi
    
    # Use default editor or vi
    EDITOR="${EDITOR:-vi}"
    
    echo -e "${YELLOW}Opening $TICKET_ID in $EDITOR...${NC}"
    $EDITOR "$TICKET_FILE"
    
    # Update the updated date
    DATE=$(date +%Y-%m-%d)
    sed -i.bak "s/^updated: .*/updated: \"$DATE\"/" "$TICKET_FILE"
    rm -f "$TICKET_FILE.bak"
    
    echo -e "${GREEN}Ticket $TICKET_ID updated${NC}"
}

# Quick status display
quick_status() {
    echo -e "${BLUE}=== Ticket Summary ===${NC}"
    echo ""
    
    TODO_COUNT=$(find "$TICKET_DIR/tickets" -name "*.yml" -type f -not -path "*/archive/*" -exec grep -l "^status: todo" {} + 2>/dev/null | wc -l)
    IN_PROGRESS_COUNT=$(find "$TICKET_DIR/tickets" -name "*.yml" -type f -not -path "*/archive/*" -exec grep -l "^status: in_progress" {} + 2>/dev/null | wc -l)
    COMPLETED_COUNT=$(find "$TICKET_DIR/tickets" -name "*.yml" -type f -not -path "*/archive/*" -exec grep -l "^status: completed" {} + 2>/dev/null | wc -l)
    ARCHIVED_COUNT=$(find "$TICKET_DIR/tickets/archive" -name "*.yml" -type f 2>/dev/null | wc -l)
    TOTAL_COMPLETED=$((COMPLETED_COUNT + ARCHIVED_COUNT))
    
    echo -e "Todo:        ${YELLOW}$TODO_COUNT${NC}"
    echo -e "In Progress: ${BLUE}$IN_PROGRESS_COUNT${NC}"
    echo -e "Completed:   ${GREEN}$TOTAL_COMPLETED${NC}"
    echo ""
    echo "Use '$0 list' for detailed view"
}

# Main command router
case "${1:-}" in
    new|n)
        shift
        "$SCRIPT_DIR/ticket-new.sh" "$@"
        ;;
    list|ls)
        shift
        "$SCRIPT_DIR/ticket-list.sh" "$@"
        ;;
    status|s)
        shift
        "$SCRIPT_DIR/ticket-status.sh" "$@"
        ;;
    report|r)
        shift
        "$SCRIPT_DIR/ticket-report.sh" "$@"
        ;;
    show)
        shift
        show_ticket "$@"
        ;;
    edit)
        shift
        edit_ticket "$@"
        ;;
    help|h|--help|-h)
        show_usage
        ;;
    "")
        quick_status
        ;;
    *)
        echo -e "${RED}Error: Unknown command '$1'${NC}"
        echo ""
        show_usage
        exit 1
        ;;
esac