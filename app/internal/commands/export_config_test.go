package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestExportConfigCommand(t *testing.T) {
	tests := []struct {
		name        string
		exportFile  string
		exportFormat string
		wantErr     bool
		checkOutput func(t *testing.T, output []byte)
	}{
		{
			name:        "export to stdout yaml",
			exportFile:  "-",
			exportFormat: "yaml",
			wantErr:     false,
			checkOutput: func(t *testing.T, output []byte) {
				// Check that it's valid YAML
				var data map[string]interface{}
				if err := yaml.Unmarshal(output, &data); err != nil {
					t.Errorf("Output is not valid YAML: %v", err)
				}
				
				// Check for expected sections
				if _, ok := data["spells"]; !ok {
					t.Error("Output missing 'spells' section")
				}
				if _, ok := data["grimoire"]; !ok {
					t.Error("Output missing 'grimoire' section")
				}
			},
		},
		{
			name:        "export to file yaml",
			exportFile:  "test_export.yml",
			exportFormat: "yaml",
			wantErr:     false,
			checkOutput: func(t *testing.T, output []byte) {
				// Check that file was created
				if _, err := os.Stat("test_export.yml"); err != nil {
					t.Errorf("Export file not created: %v", err)
				}
				
				// Clean up
				os.Remove("test_export.yml")
			},
		},
		{
			name:        "unsupported format",
			exportFile:  "-",
			exportFormat: "xml",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory with config
			tmpDir := t.TempDir()
			configFile := filepath.Join(tmpDir, "spellbook.yml")
			configContent := `spells:
  e: editor
grimoire:
  editor:
    type: app
    command: vi`
			
			if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
				t.Fatalf("Failed to create test config: %v", err)
			}

			// Create command
			cmd := NewExportConfigCommand(
				func() string { return tmpDir },
				func() []string { return []string{tmpDir} },
			)

			// Create flags
			flags := &Flags{
				ExportConfig: tt.exportFile,
				ExportFormat: tt.exportFormat,
			}

			// Check IsActive
			if !cmd.IsActive(flags) {
				t.Error("Command should be active with export config set")
			}

			// Capture output if exporting to stdout
			var output bytes.Buffer
			if tt.exportFile == "-" {
				oldStdout := os.Stdout
				r, w, _ := os.Pipe()
				os.Stdout = w
				
				err := cmd.Execute(flags)
				
				w.Close()
				os.Stdout = oldStdout
				
				buf := make([]byte, 1024)
				n, _ := r.Read(buf)
				output.Write(buf[:n])
				
				if (err != nil) != tt.wantErr {
					t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				}
				
				if !tt.wantErr && tt.checkOutput != nil {
					tt.checkOutput(t, output.Bytes())
				}
			} else {
				// Change to temp directory for file output
				oldDir, _ := os.Getwd()
				os.Chdir(tmpDir)
				defer os.Chdir(oldDir)
				
				err := cmd.Execute(flags)
				
				if (err != nil) != tt.wantErr {
					t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				}
				
				if !tt.wantErr && tt.checkOutput != nil {
					tt.checkOutput(t, nil)
				}
			}
		})
	}
}

func TestExportConfigCommand_Properties(t *testing.T) {
	cmd := NewExportConfigCommand(
		func() string { return "." },
		func() []string { return []string{"."} },
	)

	if cmd.Name() != "ExportConfig" {
		t.Errorf("Name() = %v, want %v", cmd.Name(), "ExportConfig")
	}

	if cmd.FlagName() != "export-config" {
		t.Errorf("FlagName() = %v, want %v", cmd.FlagName(), "export-config")
	}

	if !strings.Contains(cmd.Description(), "Export") {
		t.Errorf("Description() = %v, should contain 'Export'", cmd.Description())
	}

	if cmd.Group() != "config" {
		t.Errorf("Group() = %v, want %v", cmd.Group(), "config")
	}

	if !cmd.HasOptions() {
		t.Error("HasOptions() = false, want true")
	}
}

func TestExportConfigCommand_NotActive(t *testing.T) {
	cmd := NewExportConfigCommand(
		func() string { return "." },
		func() []string { return []string{"."} },
	)

	flags := &Flags{
		ExportConfig: "", // Empty means not active
	}

	if cmd.IsActive(flags) {
		t.Error("Command should not be active with empty export config")
	}
}