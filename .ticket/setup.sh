#!/bin/bash
# setup.sh - Initialize SilentCast ticket management system

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
TICKET_DIR="$SCRIPT_DIR"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo -e "${BLUE}Ticket Management System Setup${NC}"
echo "========================================="
echo ""

# Check if already set up
if [ -L "$PROJECT_ROOT/ticket" ]; then
    echo -e "${YELLOW}Warning: Symlink 'ticket' already exists in project root${NC}"
    read -p "Do you want to recreate it? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Skipping symlink creation..."
    else
        rm -f "$PROJECT_ROOT/ticket"
        ln -s .ticket/scripts/ticket.sh "$PROJECT_ROOT/ticket"
        echo -e "${GREEN}✓ Recreated symlink${NC}"
    fi
else
    # Create symlink
    ln -s .ticket/scripts/ticket.sh "$PROJECT_ROOT/ticket"
    echo -e "${GREEN}✓ Created symlink 'ticket' in project root${NC}"
fi

# Check directory structure
echo ""
echo "Checking directory structure..."

# Create directories if they don't exist
DIRS=(
    "$TICKET_DIR/meta"
    "$TICKET_DIR/meta/templates"
    "$TICKET_DIR/tickets"
    "$TICKET_DIR/tickets/archive/$(date +%Y)/$(date +%m)"
    "$TICKET_DIR/scripts"
)

for dir in "${DIRS[@]}"; do
    if [ ! -d "$dir" ]; then
        mkdir -p "$dir"
        echo -e "${GREEN}✓ Created $dir${NC}"
    else
        echo -e "  Directory exists: $dir"
    fi
done

# Check essential files
echo ""
echo "Checking essential files..."

FILES=(
    "$TICKET_DIR/meta/config.yml"
    "$TICKET_DIR/meta/labels.yml"
    "$TICKET_DIR/meta/templates/default.yml"
    "$TICKET_DIR/scripts/ticket.sh"
    "$TICKET_DIR/scripts/ticket-new.sh"
    "$TICKET_DIR/scripts/ticket-list.sh"
    "$TICKET_DIR/scripts/ticket-status.sh"
    "$TICKET_DIR/scripts/ticket-report.sh"
)

MISSING_FILES=0
for file in "${FILES[@]}"; do
    if [ ! -f "$file" ]; then
        echo -e "${RED}✗ Missing: $file${NC}"
        MISSING_FILES=$((MISSING_FILES + 1))
    else
        echo -e "  File exists: $file"
    fi
done

if [ $MISSING_FILES -gt 0 ]; then
    echo ""
    echo -e "${RED}Error: Some essential files are missing!${NC}"
    echo "Please ensure all ticket system files are properly installed."
    exit 1
fi

# Make scripts executable
echo ""
echo "Setting script permissions..."
chmod +x "$TICKET_DIR/scripts/"*.sh
echo -e "${GREEN}✓ All scripts are executable${NC}"

# Initialize config if needed
if ! grep -q "next_id:" "$TICKET_DIR/meta/config.yml" 2>/dev/null; then
    echo ""
    echo -e "${YELLOW}Warning: config.yml might need initialization${NC}"
fi

# Add to .gitignore if needed
if [ -f "$PROJECT_ROOT/.gitignore" ]; then
    if ! grep -q "TICKET_STATUS.md" "$PROJECT_ROOT/.gitignore" 2>/dev/null; then
        echo ""
        echo "Adding TICKET_STATUS.md to .gitignore..."
        echo -e "\n# Ticket management\nTICKET_STATUS.md" >> "$PROJECT_ROOT/.gitignore"
        echo -e "${GREEN}✓ Updated .gitignore${NC}"
    fi
fi

# Show usage
echo ""
echo -e "${GREEN}Setup complete!${NC}"
echo ""
echo "You can now use the ticket management system:"
echo ""
echo "  ${BLUE}./ticket${NC}              - Show quick status"
echo "  ${BLUE}./ticket new${NC}          - Create a new ticket"
echo "  ${BLUE}./ticket list${NC}         - List all tickets"
echo "  ${BLUE}./ticket show T001${NC}    - Show ticket details"
echo "  ${BLUE}./ticket help${NC}         - Show all commands"
echo ""
echo "For more information, see:"
echo "  - .ticket/README.md"
echo "  - .ticket/scripts/README.md"
echo ""

# Quick status
echo "Current ticket status:"
"$PROJECT_ROOT/ticket" 2>/dev/null || echo "No tickets found yet."