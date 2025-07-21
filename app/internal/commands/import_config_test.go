package commands

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestImportConfigCommand(t *testing.T) {
	tests := []struct {
		name         string
		importFile   string
		inputContent string
		wantErr      bool
		checkResult  func(t *testing.T, tmpDir string)
	}{
		{
			name:       "import valid yaml",
			importFile: "test.yml",
			inputContent: `spells:
  e: editor
grimoire:
  editor:
    type: app
    command: vi`,
			wantErr: false,
			checkResult: func(t *testing.T, tmpDir string) {
				// Check that spellbook.yml was created
				configPath := filepath.Join(tmpDir, "spellbook.yml")
				if _, err := os.Stat(configPath); err != nil {
					t.Errorf("Config file not created: %v", err)
				}

				// Read and verify content
				content, err := os.ReadFile(configPath)
				if err != nil {
					t.Errorf("Failed to read config: %v", err)
				}

				if !strings.Contains(string(content), "editor") {
					t.Error("Imported config missing expected content")
				}
			},
		},
		{
			name:         "import invalid yaml",
			importFile:   "invalid.yml",
			inputContent: "invalid: yaml: content:",
			wantErr:      true,
		},
		{
			name:       "import with backup",
			importFile: "new.yml",
			inputContent: `spells:
  t: terminal`,
			wantErr: false,
			checkResult: func(t *testing.T, tmpDir string) {
				// Should have created a backup
				files, _ := os.ReadDir(tmpDir)
				backupFound := false
				for _, f := range files {
					if strings.Contains(f.Name(), ".backup.") {
						backupFound = true
						break
					}
				}
				if !backupFound {
					t.Error("No backup file created")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tmpDir := t.TempDir()

			// Create existing config for backup test
			if tt.name == "import with backup" {
				existingConfig := filepath.Join(tmpDir, "spellbook.yml")
				os.WriteFile(existingConfig, []byte("existing: config"), 0o644)
			}

			// Create import file
			importPath := filepath.Join(tmpDir, tt.importFile)
			if err := os.WriteFile(importPath, []byte(tt.inputContent), 0o644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Create command
			cmd := NewImportConfigCommand(
				func() string { return tmpDir },
				func() []string { return []string{tmpDir} },
			)

			// Create flags
			flags := &Flags{
				ImportConfig: importPath,
			}

			// Check IsActive
			if !cmd.IsActive(flags) {
				t.Error("Command should be active with import config set")
			}

			// Execute
			err := cmd.Execute(flags)

			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && tt.checkResult != nil {
				tt.checkResult(t, tmpDir)
			}
		})
	}
}

func TestImportConfigCommand_TarGz(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()

	// Create a tar.gz with config
	var buf bytes.Buffer

	// Create gzip writer
	gzWriter := gzip.NewWriter(&buf)

	// Create tar writer
	tarWriter := tar.NewWriter(gzWriter)

	// Add spellbook.yml
	configContent := `spells:
  e: editor
grimoire:
  editor:
    type: app
    command: vi`

	header := &tar.Header{
		Name: "spellbook.yml",
		Mode: 0o644,
		Size: int64(len(configContent)),
	}

	if err := tarWriter.WriteHeader(header); err != nil {
		t.Fatalf("Failed to write tar header: %v", err)
	}

	if _, err := tarWriter.Write([]byte(configContent)); err != nil {
		t.Fatalf("Failed to write tar content: %v", err)
	}

	// Close writers
	tarWriter.Close()
	gzWriter.Close()

	// Write tar.gz file
	archivePath := filepath.Join(tmpDir, "backup.tar.gz")
	if err := os.WriteFile(archivePath, buf.Bytes(), 0o644); err != nil {
		t.Fatalf("Failed to write archive: %v", err)
	}

	// Create command
	cmd := NewImportConfigCommand(
		func() string { return tmpDir },
		func() []string { return []string{tmpDir} },
	)

	// Create flags
	flags := &Flags{
		ImportConfig: archivePath,
	}

	// Execute
	err := cmd.Execute(flags)
	if err != nil {
		t.Errorf("Execute() error = %v, want nil", err)
	}

	// Check that config was imported
	configPath := filepath.Join(tmpDir, "spellbook.yml")
	if _, err := os.Stat(configPath); err != nil {
		t.Errorf("Config file not created: %v", err)
	}
}

func TestImportConfigCommand_Properties(t *testing.T) {
	cmd := NewImportConfigCommand(
		func() string { return "." },
		func() []string { return []string{"."} },
	)

	if cmd.Name() != "ImportConfig" {
		t.Errorf("Name() = %v, want %v", cmd.Name(), "ImportConfig")
	}

	if cmd.FlagName() != "import-config" {
		t.Errorf("FlagName() = %v, want %v", cmd.FlagName(), "import-config")
	}

	if !strings.Contains(cmd.Description(), "Import") {
		t.Errorf("Description() = %v, should contain 'Import'", cmd.Description())
	}

	if cmd.Group() != "config" {
		t.Errorf("Group() = %v, want %v", cmd.Group(), "config")
	}

	if !cmd.HasOptions() {
		t.Error("HasOptions() = false, want true")
	}
}

func TestImportConfigCommand_NotActive(t *testing.T) {
	cmd := NewImportConfigCommand(
		func() string { return "." },
		func() []string { return []string{"."} },
	)

	flags := &Flags{
		ImportConfig: "", // Empty means not active
	}

	if cmd.IsActive(flags) {
		t.Error("Command should not be active with empty import config")
	}
}
