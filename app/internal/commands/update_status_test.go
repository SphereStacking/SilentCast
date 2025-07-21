package commands

import (
	"strings"
	"testing"
)

func TestUpdateStatusCommand_Name(t *testing.T) {
	cmd := NewUpdateStatusCommand(func() string { return "/tmp" })
	if cmd.Name() != "UpdateStatus" {
		t.Errorf("Expected Name() = 'UpdateStatus', got '%s'", cmd.Name())
	}
}

func TestUpdateStatusCommand_Description(t *testing.T) {
	cmd := NewUpdateStatusCommand(func() string { return "/tmp" })
	description := cmd.Description()
	if description == "" {
		t.Error("Expected non-empty description")
	}
	if !strings.Contains(description, "update") {
		t.Errorf("Expected description to contain 'update', got '%s'", description)
	}
}

func TestUpdateStatusCommand_FlagName(t *testing.T) {
	cmd := NewUpdateStatusCommand(func() string { return "/tmp" })
	if cmd.FlagName() != "update-status" {
		t.Errorf("Expected FlagName() = 'update-status', got '%s'", cmd.FlagName())
	}
}

func TestUpdateStatusCommand_Group(t *testing.T) {
	cmd := NewUpdateStatusCommand(func() string { return "/tmp" })
	if cmd.Group() != "utility" {
		t.Errorf("Expected Group() = 'utility', got '%s'", cmd.Group())
	}
}

func TestUpdateStatusCommand_HasOptions(t *testing.T) {
	cmd := NewUpdateStatusCommand(func() string { return "/tmp" })
	if !cmd.HasOptions() {
		t.Error("Expected HasOptions() = true")
	}
}

func TestUpdateStatusCommand_IsActive(t *testing.T) {
	cmd := NewUpdateStatusCommand(func() string { return "/tmp" })

	// Test with Flags struct
	flags := &Flags{UpdateStatus: true}
	if !cmd.IsActive(flags) {
		t.Error("Expected IsActive() = true with UpdateStatus flag set")
	}

	flags.UpdateStatus = false
	if cmd.IsActive(flags) {
		t.Error("Expected IsActive() = false with UpdateStatus flag unset")
	}

	// Test with map for unit tests
	flagsMap := map[string]interface{}{"update-status": true}
	if !cmd.IsActive(flagsMap) {
		t.Error("Expected IsActive() = true with update-status in map")
	}

	flagsMap["update-status"] = false
	if cmd.IsActive(flagsMap) {
		t.Error("Expected IsActive() = false with update-status false in map")
	}
}

func TestUpdateStatusCommand_Execute_InvalidFlags(t *testing.T) {
	cmd := NewUpdateStatusCommand(func() string { return "/tmp" })

	// Test with invalid flags type
	err := cmd.Execute("invalid")
	if err == nil {
		t.Error("Expected error for invalid flags type")
	}

	expectedMsg := "invalid flags type"
	if !strings.Contains(err.Error(), expectedMsg) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestUpdateStatusCommand_Execute_NetworkFailure(t *testing.T) {
	cmd := NewUpdateStatusCommand(func() string { return "/tmp" })

	// Test with valid flags that should result in network failure
	flags := map[string]interface{}{
		"update-status": true,
	}

	// This should not fail the command itself, just show the error
	err := cmd.Execute(flags)
	if err != nil {
		t.Errorf("Command should not fail for network errors, got: %v", err)
	}
}

func TestFormatUpdateSize(t *testing.T) {
	tests := []struct {
		size     int64
		expected string
	}{
		{0, "0 B"},
		{500, "500 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := formatUpdateSize(tt.size)
			if result != tt.expected {
				t.Errorf("formatUpdateSize(%d) = %s, want %s", tt.size, result, tt.expected)
			}
		})
	}
}

func TestIndentText(t *testing.T) {
	input := "line1\nline2\nline3"
	expected := "  line1\n  line2\n  line3"
	result := indentText(input, "  ")

	if result != expected {
		t.Errorf("indentText() = %q, want %q", result, expected)
	}

	// Test with empty lines
	input = "line1\n\nline3"
	expected = "  line1\n\n  line3"
	result = indentText(input, "  ")

	if result != expected {
		t.Errorf("indentText() with empty lines = %q, want %q", result, expected)
	}
}
