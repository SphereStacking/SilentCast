package action

import (
	"context"
	"testing"

	"github.com/SphereStacking/silentcast/internal/config"
)

func TestManager_Execute(t *testing.T) {
	grimoire := map[string]config.ActionConfig{
		"echo_test": {
			Type:        "script",
			Command: "echo 'test'",
			Description: "Test echo command",
		},
		"invalid_app": {
			Type:        "app",
			Command: "/non/existent/app",
			Description: "Non-existent app",
		},
		"unknown_type": {
			Type:        "unknown",
			Command: "something",
		},
	}
	
	manager := NewManager(grimoire)
	ctx := context.Background()
	
	tests := []struct {
		name      string
		spellName string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "Valid script spell",
			spellName: "echo_test",
			wantErr:   false,
		},
		{
			name:      "Non-existent spell",
			spellName: "non_existent",
			wantErr:   true,
			errMsg:    "not found in grimoire",
		},
		{
			name:      "Unknown action type",
			spellName: "unknown_type",
			wantErr:   true,
			errMsg:    "unknown action type",
		},
		{
			name:      "Invalid app path",
			spellName: "invalid_app",
			wantErr:   true,
			errMsg:    "application not found",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.Execute(ctx, tt.spellName)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errMsg, err.Error())
				}
			}
		})
	}
}

func TestAppExecutor_Execute(t *testing.T) {
	tests := []struct {
		name    string
		config  config.ActionConfig
		wantErr bool
	}{
		{
			name: "Valid executable in PATH",
			config: config.ActionConfig{
				Type:        "app",
				Command: "echo", // Should exist in PATH
			},
			wantErr: false,
		},
		{
			name: "Non-existent application",
			config: config.ActionConfig{
				Type:        "app",
				Command: "/non/existent/application",
			},
			wantErr: true,
		},
		{
			name: "With arguments",
			config: config.ActionConfig{
				Type:        "app",
				Command: "echo",
				Args:        []string{"hello", "world"},
			},
			wantErr: false,
		},
		{
			name: "With environment variables",
			config: config.ActionConfig{
				Type:        "app",
				Command: "echo",
				Env: map[string]string{
					"TEST_VAR": "test_value",
				},
			},
			wantErr: false,
		},
		{
			name: "Environment variable expansion",
			config: config.ActionConfig{
				Type:        "app",
				Command: "$HOME/non_existent", // Will expand but won't exist
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewAppExecutor(tt.config)
			ctx := context.Background()
			
			err := executor.Execute(ctx)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestScriptExecutor_Execute(t *testing.T) {
	tests := []struct {
		name    string
		config  config.ActionConfig
		wantErr bool
	}{
		{
			name: "Simple echo command",
			config: config.ActionConfig{
				Type:        "script",
				Command: "echo test",
			},
			wantErr: false,
		},
		{
			name: "Command with exit code 0",
			config: config.ActionConfig{
				Type:        "script",
				Command: "true",
			},
			wantErr: false,
		},
		{
			name: "Command with non-zero exit code",
			config: config.ActionConfig{
				Type:        "script",
				Command: "false",
			},
			wantErr: true,
		},
		{
			name: "With working directory",
			config: config.ActionConfig{
				Type:        "script",
				Command: "pwd",
				WorkingDir:  "/tmp",
			},
			wantErr: false,
		},
		{
			name: "Environment variable in command",
			config: config.ActionConfig{
				Type:        "script",
				Command: "echo $HOME",
			},
			wantErr: false,
		},
		{
			name: "Empty command",
			config: config.ActionConfig{
				Type:        "script",
				Command: "",
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := NewScriptExecutor(tt.config)
			ctx := context.Background()
			
			err := executor.Execute(ctx)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExecutor_String(t *testing.T) {
	tests := []struct {
		name     string
		executor Executor
		want     string
	}{
		{
			name: "App with description",
			executor: NewAppExecutor(config.ActionConfig{
				Type:        "app",
				Command: "/usr/bin/app",
				Description: "My cool app",
			}),
			want: "My cool app",
		},
		{
			name: "App without description",
			executor: NewAppExecutor(config.ActionConfig{
				Type:        "app",
				Command: "/usr/bin/app",
			}),
			want: "Launch /usr/bin/app",
		},
		{
			name: "Script with description",
			executor: NewScriptExecutor(config.ActionConfig{
				Type:        "script",
				Command: "git status",
				Description: "Show git status",
			}),
			want: "Show git status",
		},
		{
			name: "Script without description",
			executor: NewScriptExecutor(config.ActionConfig{
				Type:        "script",
				Command: "git status",
			}),
			want: "Run script: git status",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.executor.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && (s[0:len(substr)] == substr || contains(s[1:], substr)))
}