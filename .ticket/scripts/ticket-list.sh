#!/bin/bash
# ticket-list.sh - List tickets with filters

set -e

TICKET_DIR="$(cd "$(dirname "$0")/../.." && pwd)/.ticket"

# Default values
STATUS=""
PRIORITY=""
TYPE=""
LABEL=""
FORMAT="simple"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --status)
            STATUS="$2"
            shift 2
            ;;
        --priority)
            PRIORITY="$2"
            shift 2
            ;;
        --type)
            TYPE="$2"
            shift 2
            ;;
        --label)
            LABEL="$2"
            shift 2
            ;;
        --format)
            FORMAT="$2"
            shift 2
            ;;
        --help|-h)
            echo "Usage: $0 [options]"
            echo "Options:"
            echo "  --status STATUS    Filter by status (todo, in_progress, etc.)"
            echo "  --priority PRIO    Filter by priority (critical, high, medium, low)"
            echo "  --type TYPE        Filter by type (feature, bug, etc.)"
            echo "  --label LABEL      Filter by label"
            echo "  --format FORMAT    Output format (simple, detailed, csv)"
            echo "  --help             Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Find tickets in tickets directory (including archive)
TICKETS=$(find "$TICKET_DIR/tickets" -name '*.yml' -type f 2>/dev/null | sort || true)

if [ -z "$TICKETS" ]; then
    echo "No tickets found"
    exit 0
fi

# Header based on format
case $FORMAT in
    simple)
        printf "%-6s %-12s %-8s %-50s\n" "ID" "STATUS" "PRIORITY" "TITLE"
        printf "%-6s %-12s %-8s %-50s\n" "------" "------------" "--------" "--------------------------------------------------"
        ;;
    detailed)
        echo "=== TICKETS ==="
        ;;
    csv)
        echo "ID,Status,Priority,Type,Title"
        ;;
esac

# Process each ticket
for ticket in $TICKETS; do
    # Skip if file doesn't exist (race condition protection)
    [ -f "$ticket" ] || continue
    
    # Extract fields
    ID=$(grep "^id:" "$ticket" 2>/dev/null | awk '{print $2}')
    TSTATUS=$(grep "^status:" "$ticket" 2>/dev/null | awk '{print $2}')
    TPRIORITY=$(grep "^priority:" "$ticket" 2>/dev/null | awk '{print $2}')
    TTYPE=$(grep "^type:" "$ticket" 2>/dev/null | awk '{print $2}')
    TITLE=$(grep "^title:" "$ticket" 2>/dev/null | cut -d'"' -f2)
    
    # Apply filters
    if [ -n "$STATUS" ] && [ "$TSTATUS" != "$STATUS" ]; then
        continue
    fi
    
    if [ -n "$PRIORITY" ] && [ "$TPRIORITY" != "$PRIORITY" ]; then
        continue
    fi
    
    if [ -n "$TYPE" ] && [ "$TTYPE" != "$TYPE" ]; then
        continue
    fi
    
    if [ -n "$LABEL" ]; then
        if ! grep -q "- $LABEL" "$ticket" 2>/dev/null; then
            continue
        fi
    fi
    
    # Output based on format
    case $FORMAT in
        simple)
            printf "%-6s %-12s %-8s %-50s\n" "$ID" "$TSTATUS" "$TPRIORITY" "${TITLE:0:50}"
            ;;
        detailed)
            echo ""
            echo "Ticket: $ID"
            echo "Title: $TITLE"
            echo "Status: $TSTATUS"
            echo "Priority: $TPRIORITY"
            echo "Type: $TTYPE"
            echo "File: $ticket"
            echo "---"
            ;;
        csv)
            echo "$ID,$TSTATUS,$TPRIORITY,$TTYPE,\"$TITLE\""
            ;;
    esac
done

# Summary for simple format
if [ "$FORMAT" = "simple" ]; then
    echo ""
    TOTAL=$(echo "$TICKETS" | wc -l)
    echo "Total tickets shown: $TOTAL"
fi