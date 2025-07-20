package commands

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/SphereStacking/silentcast/internal/config"
)

// ListSpellsCommand lists all configured spells
type ListSpellsCommand struct {
	getConfigPath func() string
}

// NewListSpellsCommand creates a new list spells command
func NewListSpellsCommand(getConfigPath func() string) Command {
	return &ListSpellsCommand{
		getConfigPath: getConfigPath,
	}
}

// Name returns the command name
func (c *ListSpellsCommand) Name() string {
	return "List Spells"
}

// Description returns the command description
func (c *ListSpellsCommand) Description() string {
	return "List all configured spells"
}

// FlagName returns the flag name
func (c *ListSpellsCommand) FlagName() string {
	return "list-spells"
}

// IsActive checks if the command should run
func (c *ListSpellsCommand) IsActive(flags interface{}) bool {
	f, ok := flags.(*Flags)
	if !ok {
		return false
	}
	return f.ListSpells
}

// Execute runs the command
func (c *ListSpellsCommand) Execute(flags interface{}) error {
	f, ok := flags.(*Flags)
	if !ok {
		return fmt.Errorf("invalid flags type")
	}

	// Load configuration
	loader := config.NewLoader(c.getConfigPath())
	cfg, err := loader.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Collect spell information
	type spellInfo struct {
		sequence    string
		name        string
		actionType  string
		command     string
		description string
	}

	spells := make([]spellInfo, 0, len(cfg.Shortcuts))

	for sequence, spellName := range cfg.Shortcuts {
		// Apply filter if provided
		if f.ListFilter != "" {
			filterLower := strings.ToLower(f.ListFilter)
			if !strings.Contains(strings.ToLower(sequence), filterLower) &&
				!strings.Contains(strings.ToLower(spellName), filterLower) &&
				!strings.Contains(strings.ToLower(cfg.Actions[spellName].Description), filterLower) {
				continue
			}
		}

		action, exists := cfg.Actions[spellName]
		if !exists {
			// Spell references non-existent action
			spells = append(spells, spellInfo{
				sequence:    sequence,
				name:        spellName,
				actionType:  "unknown",
				command:     "N/A",
				description: "⚠️  Action not found",
			})
			continue
		}

		// Format command based on type
		command := action.Command
		if action.Type == "url" && !strings.HasPrefix(command, "http") {
			command = "https://" + command
		}

		// Truncate long commands
		if len(command) > 40 {
			command = command[:37] + "..."
		}

		spells = append(spells, spellInfo{
			sequence:    sequence,
			name:        spellName,
			actionType:  action.Type,
			command:     command,
			description: action.Description,
		})
	}

	// Sort by sequence for consistent output
	sort.Slice(spells, func(i, j int) bool {
		// Sort single-key spells before multi-key sequences
		iHasComma := strings.Contains(spells[i].sequence, ",")
		jHasComma := strings.Contains(spells[j].sequence, ",")

		if iHasComma != jHasComma {
			return !iHasComma // Single keys first
		}

		return spells[i].sequence < spells[j].sequence
	})

	// Display results
	if len(spells) == 0 {
		if f.ListFilter != "" {
			fmt.Printf("No spells found matching filter: %s\n", f.ListFilter)
		} else {
			fmt.Println("No spells configured")
		}
		return nil
	}

	// Print header
	fmt.Printf("\n🪄 SilentCast Spells")
	if f.ListFilter != "" {
		fmt.Printf(" (filtered: %s)", f.ListFilter)
	}
	fmt.Printf("\n")
	fmt.Printf("Prefix: %s\n\n", cfg.Hotkeys.Prefix)

	// Use tabwriter for aligned output
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	// Header
	fmt.Fprintln(w, "SEQUENCE\tNAME\tTYPE\tCOMMAND\tDESCRIPTION")
	fmt.Fprintln(w, "--------\t----\t----\t-------\t-----------")

	// Spells
	for _, spell := range spells {
		// Add prefix to sequence for clarity
		fullSequence := cfg.Hotkeys.Prefix + " → " + spell.sequence

		// Add type emoji
		typeIcon := ""
		switch spell.actionType {
		case "app":
			typeIcon = "📱 "
		case "script":
			typeIcon = "📜 "
		case "url":
			typeIcon = "🌐 "
		default:
			typeIcon = "❓ "
		}

		fmt.Fprintf(w, "%s\t%s\t%s%s\t%s\t%s\n",
			fullSequence,
			spell.name,
			typeIcon,
			spell.actionType,
			spell.command,
			spell.description,
		)
	}

	w.Flush()

	// Summary
	fmt.Printf("\n📊 Total: %d spells\n", len(spells))

	// Count by type
	typeCounts := make(map[string]int)
	for _, spell := range spells {
		typeCounts[spell.actionType]++
	}

	if len(typeCounts) > 1 {
		fmt.Print("   By type: ")
		first := true
		for actionType, count := range typeCounts {
			if !first {
				fmt.Print(", ")
			}
			fmt.Printf("%s (%d)", actionType, count)
			first = false
		}
		fmt.Println()
	}

	// Show warnings if any
	warningCount := 0
	for _, spell := range spells {
		if spell.actionType == "unknown" {
			warningCount++
		}
	}
	if warningCount > 0 {
		fmt.Printf("\n⚠️  Warning: %d spell(s) reference non-existent actions\n", warningCount)
	}

	return nil
}

// Group returns the command group
func (c *ListSpellsCommand) Group() string {
	return "spell"
}

// HasOptions returns if this command has additional options
func (c *ListSpellsCommand) HasOptions() bool {
	return true
}
