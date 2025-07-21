//go:build linux

package updater

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLinuxUpdater_CanReplaceRunningExecutable(t *testing.T) {
	updater := &linuxUpdater{}
	if !updater.CanReplaceRunningExecutable() {
		t.Error("Linux should be able to replace running executables")
	}
}

func TestLinuxUpdater_MakeExecutable(t *testing.T) {
	// Create a temporary file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_executable")

	if err := os.WriteFile(tmpFile, []byte("#!/bin/sh\necho test"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	updater := &linuxUpdater{}
	if err := updater.MakeExecutable(tmpFile); err != nil {
		t.Errorf("MakeExecutable failed: %v", err)
	}

	// Check permissions
	info, err := os.Stat(tmpFile)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	// Check if file is executable (0o755 = -rwxr-xr-x)
	mode := info.Mode()
	if mode.Perm() != 0o755 {
		t.Errorf("Expected permissions 0o755, got %v", mode.Perm())
	}
}

func TestLinuxUpdater_RestartApplication(t *testing.T) {
	t.Skip("Skipping RestartApplication test as it would replace the test process")

	// The actual test would be:
	// updater := &linuxUpdater{}
	// err := updater.RestartApplication()
	// This would replace the current process, so we can't test it directly
}
