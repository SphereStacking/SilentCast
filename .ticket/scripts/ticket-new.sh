#!/bin/bash
# ticket-new.sh - Create a new ticket

set -e

TICKET_DIR="$(cd "$(dirname "$0")/../.." && pwd)/.ticket"
CONFIG_FILE="$TICKET_DIR/meta/config.yml"
TEMPLATE_DIR="$TICKET_DIR/meta/templates"

# Default values
TYPE="feature"
PRIORITY="medium"
TITLE=""
TEMPLATE="default"

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --type)
            TYPE="$2"
            shift 2
            ;;
        --priority)
            PRIORITY="$2"
            shift 2
            ;;
        --title)
            TITLE="$2"
            shift 2
            ;;
        --template)
            TEMPLATE="$2"
            shift 2
            ;;
        --help|-h)
            echo "Usage: $0 [options]"
            echo "Options:"
            echo "  --type TYPE        Ticket type (feature, bug, refactor, docs, test, chore)"
            echo "  --priority PRIO    Priority (critical, high, medium, low)"
            echo "  --title TITLE      Ticket title (required)"
            echo "  --template TMPL    Template to use (default, feature, bug)"
            echo "  --help             Show this help message"
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            exit 1
            ;;
    esac
done

# Validate required fields
if [ -z "$TITLE" ]; then
    echo "Error: --title is required"
    exit 1
fi

# Get next ticket ID from config
NEXT_ID=$(grep "next_id:" "$CONFIG_FILE" | awk '{print $2}')
if [ -z "$NEXT_ID" ]; then
    echo "Error: Could not read next_id from config"
    exit 1
fi

# Format ticket ID
TICKET_ID=$(printf "T%03d" "$NEXT_ID")

# Create slug from title
SLUG=$(echo "$TITLE" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9]/-/g' | sed 's/--*/-/g' | sed 's/^-//' | sed 's/-$//' | cut -c1-50)

# Determine destination directory
DEST_DIR="$TICKET_DIR/tickets"

# Check if directory exists
if [ ! -d "$DEST_DIR" ]; then
    echo "Error: Directory $DEST_DIR does not exist"
    exit 1
fi

# Template file
TEMPLATE_FILE="$TEMPLATE_DIR/$TEMPLATE.yml"
if [ ! -f "$TEMPLATE_FILE" ]; then
    echo "Error: Template $TEMPLATE not found"
    echo "Available templates:"
    ls "$TEMPLATE_DIR" | sed 's/\.yml$//'
    exit 1
fi

# Create ticket file
TICKET_FILE="$DEST_DIR/$TICKET_ID-$SLUG.yml"
DATE=$(date +%Y-%m-%d)

# Copy template and update fields
cp "$TEMPLATE_FILE" "$TICKET_FILE"

# Update ticket fields using sed (portable)
sed -i.bak "s/^id: .*/id: $TICKET_ID/" "$TICKET_FILE"
sed -i.bak "s/^title: .*/title: \"$TITLE\"/" "$TICKET_FILE"
sed -i.bak "s/^type: .*/type: $TYPE/" "$TICKET_FILE"
sed -i.bak "s/^priority: .*/priority: $PRIORITY/" "$TICKET_FILE"
sed -i.bak "s/^created: .*/created: \"$DATE\"/" "$TICKET_FILE"
sed -i.bak "s/^updated: .*/updated: \"$DATE\"/" "$TICKET_FILE"


# Remove backup files
rm -f "$TICKET_FILE.bak"

# Update next_id in config
NEW_ID=$((NEXT_ID + 1))
sed -i.bak "s/^next_id: .*/next_id: $NEW_ID/" "$CONFIG_FILE"
rm -f "$CONFIG_FILE.bak"

echo "Created ticket: $TICKET_FILE"
echo "Ticket ID: $TICKET_ID"
echo "Next ID updated to: $NEW_ID"
echo ""
echo "Don't forget to:"
echo "1. Fill in the description and other fields"
echo "2. Add appropriate labels"
echo "3. Set dependencies if needed"