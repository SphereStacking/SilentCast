#!/bin/bash
# ticket-status.sh - Update ticket status

set -e

TICKET_DIR="$(cd "$(dirname "$0")/../.." && pwd)/.ticket"

# Parse arguments
if [ $# -lt 2 ]; then
    echo "Usage: $0 TICKET_ID NEW_STATUS"
    echo "Valid statuses: todo, in_progress, review, testing, completed, blocked"
    echo "Example: $0 T001 in_progress"
    exit 1
fi

TICKET_ID="$1"
NEW_STATUS="$2"
DATE=$(date +%Y-%m-%d)

# Validate status
VALID_STATUSES="todo in_progress review testing completed blocked"
if ! echo "$VALID_STATUSES" | grep -q "\b$NEW_STATUS\b"; then
    echo "Error: Invalid status '$NEW_STATUS'"
    echo "Valid statuses: $VALID_STATUSES"
    exit 1
fi

# Find ticket file
TICKET_FILE=$(find "$TICKET_DIR" -name "$TICKET_ID-*.yml" -type f 2>/dev/null | head -1)

if [ -z "$TICKET_FILE" ]; then
    echo "Error: Ticket $TICKET_ID not found"
    exit 1
fi

echo "Found ticket: $TICKET_FILE"

# Get current status
CURRENT_STATUS=$(grep "^status:" "$TICKET_FILE" | awk '{print $2}')
echo "Current status: $CURRENT_STATUS"

# Update status and updated date
sed -i.bak "s/^status: .*/status: $NEW_STATUS/" "$TICKET_FILE"
sed -i.bak "s/^updated: .*/updated: \"$DATE\"/" "$TICKET_FILE"
rm -f "$TICKET_FILE.bak"

echo "Updated status: $CURRENT_STATUS -> $NEW_STATUS"

# Handle special cases
if [ "$NEW_STATUS" = "completed" ]; then
    # Get current year and month
    YEAR=$(date +%Y)
    MONTH=$(date +%m)
    
    # Create archive directory if needed
    ARCHIVE_DIR="$TICKET_DIR/tickets/archive/$YEAR/$MONTH"
    mkdir -p "$ARCHIVE_DIR"
    
    # Move to archive
    BASENAME=$(basename "$TICKET_FILE")
    mv "$TICKET_FILE" "$ARCHIVE_DIR/$BASENAME"
    echo "Moved to: $ARCHIVE_DIR/$BASENAME"
    
elif [ "$CURRENT_STATUS" = "completed" ] && [ "$NEW_STATUS" != "completed" ]; then
    # Moving back from completed - move to tickets directory
    echo "Warning: Moving ticket out of completed status"
    BASENAME=$(basename "$TICKET_FILE")
    mv "$TICKET_FILE" "$TICKET_DIR/tickets/$BASENAME"
    echo "Moved back to: $TICKET_DIR/tickets/$BASENAME"
fi

echo "Status update complete!"