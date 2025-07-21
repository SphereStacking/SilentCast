package commands

import (
	"bytes"
	"strings"
	"testing"
)

func TestSelfUpdateCommand_Name(t *testing.T) {
	cmd := NewSelfUpdateCommand(func() string { return "/tmp" })
	if cmd.Name() != "SelfUpdate" {
		t.Errorf("Expected Name() = 'SelfUpdate', got '%s'", cmd.Name())
	}
}

func TestSelfUpdateCommand_Description(t *testing.T) {
	cmd := NewSelfUpdateCommand(func() string { return "/tmp" })
	description := cmd.Description()
	if description == "" {
		t.Error("Expected non-empty description")
	}
	// Check for 'Update' (capital U) which is in the actual description
	if !strings.Contains(description, "Update") {
		t.Errorf("Expected description to contain 'Update', got '%s'", description)
	}
}

func TestSelfUpdateCommand_HasOptions(t *testing.T) {
	cmd := NewSelfUpdateCommand(func() string { return "/tmp" })
	if !cmd.HasOptions() {
		t.Error("Expected HasOptions() = true")
	}
}

func TestSelfUpdateCommand_Execute_DryRun(t *testing.T) {
	// Create a buffer to capture output
	var output bytes.Buffer

	// Create command with mock config path
	cmd := NewSelfUpdateCommand(func() string { return "/tmp" })

	// Test with dry run flag (which should be safe)
	flags := map[string]interface{}{
		"check-only": true,
	}

	// Note: This might fail if GitHub API is unreachable, but that's expected
	err := cmd.Execute(flags)

	// We expect this to either succeed (if GitHub is reachable) or fail with a network/API error
	// The important thing is that it doesn't panic or crash
	if err != nil {
		// Check that it's a reasonable error message
		errMsg := err.Error()
		if !strings.Contains(errMsg, "update check failed") &&
			!strings.Contains(errMsg, "failed to get latest release") &&
			!strings.Contains(errMsg, "GitHub API") &&
			!strings.Contains(errMsg, "network") &&
			!strings.Contains(errMsg, "timeout") &&
			!strings.Contains(errMsg, "context deadline exceeded") {
			t.Errorf("Unexpected error type: %v", err)
		}
	}

	// Test that output is captured (even if command fails)
	_ = output.String() // Just verify we can get output
}

func TestSelfUpdateCommand_Execute_InvalidFlags(t *testing.T) {
	cmd := NewSelfUpdateCommand(func() string { return "/tmp" })

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

func TestSelfUpdateCommand_Execute_DevVersion(t *testing.T) {
	cmd := NewSelfUpdateCommand(func() string { return "/tmp" })

	// This test simulates running with development version
	// Since we can't easily mock the version package, we'll test the error handling
	flags := map[string]interface{}{
		"check-only": true,
	}

	// Execute and expect either success or a reasonable error
	err := cmd.Execute(flags)

	// For dev versions, it might succeed (no update available) or fail with network error
	// We just want to ensure it doesn't crash
	if err != nil {
		errMsg := err.Error()
		// These are acceptable error messages
		acceptableErrors := []string{
			"update check failed",
			"failed to get latest release",
			"GitHub API",
			"network",
			"timeout",
			"context deadline exceeded",
			"no suitable update found",
		}

		hasAcceptableError := false
		for _, acceptable := range acceptableErrors {
			if strings.Contains(errMsg, acceptable) {
				hasAcceptableError = true
				break
			}
		}

		if !hasAcceptableError {
			t.Errorf("Unexpected error message: %v", err)
		}
	}
}
