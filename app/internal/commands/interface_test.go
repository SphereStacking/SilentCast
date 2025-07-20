package commands

import (
	"errors"
	"testing"
)

// mockCommand is a test implementation of Command interface
type mockCommand struct {
	name        string
	description string
	flagName    string
	isActive    bool
	executeErr  error
	group       string
	hasOptions  bool
	executed    bool
}

func (m *mockCommand) Name() string        { return m.name }
func (m *mockCommand) Description() string { return m.description }
func (m *mockCommand) FlagName() string    { return m.flagName }
func (m *mockCommand) IsActive(flags interface{}) bool {
	return m.isActive
}
func (m *mockCommand) Execute(flags interface{}) error {
	m.executed = true
	return m.executeErr
}
func (m *mockCommand) Group() string    { return m.group }
func (m *mockCommand) HasOptions() bool { return m.hasOptions }

func TestRegistry_Register(t *testing.T) {
	tests := []struct {
		name     string
		commands []Command
		wantLen  int
	}{
		{
			name: "register single command",
			commands: []Command{
				&mockCommand{name: "test1", flagName: "test1"},
			},
			wantLen: 1,
		},
		{
			name: "register multiple commands",
			commands: []Command{
				&mockCommand{name: "test1", flagName: "test1"},
				&mockCommand{name: "test2", flagName: "test2"},
				&mockCommand{name: "test3", flagName: "test3"},
			},
			wantLen: 3,
		},
		{
			name:     "register no commands",
			commands: []Command{},
			wantLen:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRegistry()
			for _, cmd := range tt.commands {
				r.Register(cmd)
			}
			if len(r.commands) != tt.wantLen {
				t.Errorf("Registry has %d commands, want %d", len(r.commands), tt.wantLen)
			}
		})
	}
}

func TestRegistry_RegisterAll(t *testing.T) {
	r := NewRegistry()
	commands := []Command{
		&mockCommand{name: "test1", flagName: "test1"},
		&mockCommand{name: "test2", flagName: "test2"},
		&mockCommand{name: "test3", flagName: "test3"},
	}

	r.RegisterAll(commands...)

	if len(r.commands) != 3 {
		t.Errorf("Registry has %d commands, want 3", len(r.commands))
	}
}

func TestRegistry_Execute(t *testing.T) {
	tests := []struct {
		name         string
		commands     []Command
		flags        interface{}
		wantExecuted bool
		wantErr      bool
		errMsg       string
	}{
		{
			name: "execute active command",
			commands: []Command{
				&mockCommand{name: "test1", isActive: true},
			},
			flags:        &Flags{},
			wantExecuted: true,
			wantErr:      false,
		},
		{
			name: "execute command with error",
			commands: []Command{
				&mockCommand{name: "test1", isActive: true, executeErr: errors.New("execution failed")},
			},
			flags:        &Flags{},
			wantExecuted: true,
			wantErr:      true,
			errMsg:       "execution failed",
		},
		{
			name: "no active command",
			commands: []Command{
				&mockCommand{name: "test1", isActive: false},
				&mockCommand{name: "test2", isActive: false},
			},
			flags:        &Flags{},
			wantExecuted: false,
			wantErr:      false,
		},
		{
			name: "multiple commands but only one active",
			commands: []Command{
				&mockCommand{name: "test1", isActive: false},
				&mockCommand{name: "test2", isActive: true},
				&mockCommand{name: "test3", isActive: false},
			},
			flags:        &Flags{},
			wantExecuted: true,
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRegistry()
			r.RegisterAll(tt.commands...)

			executed, err := r.Execute(tt.flags)

			if executed != tt.wantExecuted {
				t.Errorf("Execute() executed = %v, want %v", executed, tt.wantExecuted)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && err != nil && err.Error() != tt.errMsg {
				t.Errorf("Execute() error = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestRegistry_GetGroups(t *testing.T) {
	tests := []struct {
		name       string
		commands   []Command
		wantGroups []string
	}{
		{
			name: "single group",
			commands: []Command{
				&mockCommand{name: "test1", group: "core"},
				&mockCommand{name: "test2", group: "core"},
			},
			wantGroups: []string{"core"},
		},
		{
			name: "multiple groups",
			commands: []Command{
				&mockCommand{name: "test1", group: "core"},
				&mockCommand{name: "test2", group: "config"},
				&mockCommand{name: "test3", group: "debug"},
			},
			wantGroups: []string{"core", "config", "debug"},
		},
		{
			name: "commands with same group",
			commands: []Command{
				&mockCommand{name: "test1", group: "config"},
				&mockCommand{name: "test2", group: "config"},
				&mockCommand{name: "test3", group: "config"},
			},
			wantGroups: []string{"config"},
		},
		{
			name:       "no commands",
			commands:   []Command{},
			wantGroups: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRegistry()
			r.RegisterAll(tt.commands...)

			groups := r.GetGroups()

			// Extract group names
			groupNames := make([]string, len(groups))
			for i, g := range groups {
				groupNames[i] = g.Name
			}

			// Check if all expected groups are present
			if len(groupNames) != len(tt.wantGroups) {
				t.Errorf("GetGroups() returned %d groups, want %d", len(groupNames), len(tt.wantGroups))
			}

			// Verify each expected group exists
			for _, want := range tt.wantGroups {
				found := false
				for _, got := range groupNames {
					if got == want {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("GetGroups() missing group %s", want)
				}
			}
		})
	}
}

func TestRegistry_GetGroupCommands(t *testing.T) {
	r := NewRegistry()
	cmd1 := &mockCommand{name: "test1", group: "config", description: "Test command 1"}
	cmd2 := &mockCommand{name: "test2", group: "config", description: "Test command 2"}
	cmd3 := &mockCommand{name: "test3", group: "debug", description: "Test command 3"}

	r.RegisterAll(cmd1, cmd2, cmd3)

	groups := r.GetGroups()

	// Find config group
	var configGroup *CommandGroup
	for _, g := range groups {
		if g.Name == "config" {
			configGroup = &g
			break
		}
	}

	if configGroup == nil {
		t.Fatal("config group not found")
	}

	if len(configGroup.Commands) != 2 {
		t.Errorf("config group has %d commands, want 2", len(configGroup.Commands))
	}

	// Verify commands in group
	foundCmd1 := false
	foundCmd2 := false
	for _, cmd := range configGroup.Commands {
		if cmd.Name() == "test1" {
			foundCmd1 = true
		}
		if cmd.Name() == "test2" {
			foundCmd2 = true
		}
	}

	if !foundCmd1 || !foundCmd2 {
		t.Error("config group missing expected commands")
	}
}
