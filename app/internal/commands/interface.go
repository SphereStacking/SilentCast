package commands

// Command defines the interface for all CLI commands
type Command interface {
	// Name returns the command name for display
	Name() string

	// Description returns a short description of the command
	Description() string

	// FlagName returns the flag name used to invoke this command
	FlagName() string

	// IsActive checks if this command should be executed based on flags
	IsActive(flags interface{}) bool

	// Execute runs the command with the given flags
	Execute(flags interface{}) error

	// Group returns the command group for help organization
	Group() string

	// HasOptions returns true if this command has additional options
	HasOptions() bool
}

// CommandGroup represents a group of related commands for help display
type CommandGroup struct {
	Name        string
	Description string
	Commands    []Command
}

// Registry holds all available commands
type Registry struct {
	commands []Command
}

// NewRegistry creates a new command registry
func NewRegistry() *Registry {
	return &Registry{
		commands: make([]Command, 0),
	}
}

// Register adds a command to the registry
func (r *Registry) Register(cmd Command) {
	r.commands = append(r.commands, cmd)
}

// RegisterAll registers multiple commands at once
func (r *Registry) RegisterAll(cmds ...Command) {
	for _, cmd := range cmds {
		r.Register(cmd)
	}
}

// GetCommands returns all registered commands
func (r *Registry) GetCommands() []Command {
	return r.commands
}

// GetGroups returns commands organized by groups
func (r *Registry) GetGroups() []CommandGroup {
	groupMap := make(map[string][]Command)

	// Group commands by their group name
	for _, cmd := range r.commands {
		group := cmd.Group()
		groupMap[group] = append(groupMap[group], cmd)
	}

	// Define group order and descriptions
	groupOrder := []struct {
		name        string
		description string
	}{
		{"core", "Basic application controls"},
		{"config", "Manage and inspect configuration"},
		{"spell", "Work with spells (keyboard shortcuts)"},
		{"debug", "Debug and test functionality"},
		{"execution", "Run commands and spells"},
		{"utility", "Utility commands"},
		{"service", "Service management"},
	}

	// Build result in defined order
	var groups []CommandGroup
	for _, g := range groupOrder {
		if cmds, exists := groupMap[g.name]; exists {
			groups = append(groups, CommandGroup{
				Name:        g.name,
				Description: g.description,
				Commands:    cmds,
			})
		}
	}

	return groups
}

// Execute runs the first active command based on flags
func (r *Registry) Execute(flags interface{}) (bool, error) {
	for _, cmd := range r.commands {
		if cmd.IsActive(flags) {
			err := cmd.Execute(flags)
			return true, err
		}
	}
	return false, nil
}
