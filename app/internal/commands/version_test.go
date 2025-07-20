package commands

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	tests := []struct {
		name         string
		version      string
		flags        *Flags
		wantActive   bool
		wantContains string
	}{
		{
			name:         "version flag set",
			version:      "1.2.3",
			flags:        &Flags{Version: true},
			wantActive:   true,
			wantContains: "SilentCast v1.2.3",
		},
		{
			name:         "version flag not set",
			version:      "1.2.3",
			flags:        &Flags{Version: false},
			wantActive:   false,
			wantContains: "",
		},
		{
			name:         "dev version",
			version:      "0.1.0-dev",
			flags:        &Flags{Version: true},
			wantActive:   true,
			wantContains: "SilentCast v0.1.0-dev",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewVersionCommand(tt.version)

			// Test metadata
			if cmd.Name() != "Version" {
				t.Errorf("Name() = %v, want %v", cmd.Name(), "Version")
			}

			if cmd.Description() != "Print version and exit" {
				t.Errorf("Description() = %v, want %v", cmd.Description(), "Print version and exit")
			}

			if cmd.FlagName() != "version" {
				t.Errorf("FlagName() = %v, want %v", cmd.FlagName(), "version")
			}

			if cmd.Group() != "core" {
				t.Errorf("Group() = %v, want %v", cmd.Group(), "core")
			}

			if cmd.HasOptions() != false {
				t.Errorf("HasOptions() = %v, want %v", cmd.HasOptions(), false)
			}

			// Test IsActive
			if got := cmd.IsActive(tt.flags); got != tt.wantActive {
				t.Errorf("IsActive() = %v, want %v", got, tt.wantActive)
			}

			// Test Execute if active
			if tt.wantActive {
				// Capture stdout
				old := os.Stdout
				r, w, _ := os.Pipe()
				os.Stdout = w

				err := cmd.Execute(tt.flags)

				w.Close()
				os.Stdout = old

				if err != nil {
					t.Errorf("Execute() error = %v, want nil", err)
				}

				// Read captured output
				var buf bytes.Buffer
				buf.ReadFrom(r)
				output := buf.String()

				if !strings.Contains(output, tt.wantContains) {
					t.Errorf("Execute() output = %v, want to contain %v", output, tt.wantContains)
				}
			}
		})
	}
}

func TestVersionCommand_InvalidFlags(t *testing.T) {
	cmd := NewVersionCommand("1.0.0")

	// Test with invalid flags type
	if cmd.IsActive("invalid") {
		t.Error("IsActive() should return false for invalid flags type")
	}

	// Test Execute with invalid flags type
	err := cmd.Execute("invalid")
	if err == nil || !strings.Contains(err.Error(), "invalid flags type") {
		t.Errorf("Execute() should return error for invalid flags type, got %v", err)
	}
}
