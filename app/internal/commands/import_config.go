package commands

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/SphereStacking/silentcast/internal/config"
	"github.com/SphereStacking/silentcast/pkg/logger"
	"gopkg.in/yaml.v3"
)

// ImportConfigCommand imports configuration from backup
type ImportConfigCommand struct {
	configPathFunc       func() string
	configSearchPathFunc func() []string
}

// NewImportConfigCommand creates a new import config command
func NewImportConfigCommand(configPath func() string, configSearchPaths func() []string) Command {
	return &ImportConfigCommand{
		configPathFunc:       configPath,
		configSearchPathFunc: configSearchPaths,
	}
}

// Name returns the command name
func (c *ImportConfigCommand) Name() string {
	return "ImportConfig"
}

// Description returns the command description
func (c *ImportConfigCommand) Description() string {
	return "Import configuration from backup file"
}

// FlagName returns the flag name
func (c *ImportConfigCommand) FlagName() string {
	return "import-config"
}

// IsActive checks if the command should run
func (c *ImportConfigCommand) IsActive(flags interface{}) bool {
	f, ok := flags.(*Flags)
	if !ok {
		return false
	}
	return f.ImportConfig != ""
}

// Execute runs the command
func (c *ImportConfigCommand) Execute(flags interface{}) error {
	f, ok := flags.(*Flags)
	if !ok {
		return fmt.Errorf("invalid flags type")
	}

	// Determine input reader
	var reader io.Reader
	var closeFunc func() error

	if f.ImportConfig == "-" {
		// Import from stdin
		reader = os.Stdin
		fmt.Fprintln(os.Stderr, "📥 Reading configuration from stdin...")
	} else {
		// Import from file
		file, err := os.Open(f.ImportConfig)
		if err != nil {
			return fmt.Errorf("failed to open import file: %w", err)
		}
		reader = file
		closeFunc = file.Close
		fmt.Fprintf(os.Stderr, "📥 Importing configuration from: %s\n", f.ImportConfig)
	}

	// Detect format and import
	var err error
	if strings.HasSuffix(strings.ToLower(f.ImportConfig), ".tar.gz") ||
		strings.HasSuffix(strings.ToLower(f.ImportConfig), ".tgz") {
		err = c.importTarGz(reader)
	} else {
		// Default to YAML
		err = c.importYAML(reader)
	}

	if closeFunc != nil {
		if closeErr := closeFunc(); closeErr != nil {
			logger.Warn("Failed to close input file: %v", closeErr)
		}
	}

	if err != nil {
		return fmt.Errorf("import failed: %w", err)
	}

	fmt.Fprintln(os.Stderr, "✅ Configuration imported successfully")
	fmt.Fprintln(os.Stderr, "   • Run 'silentcast --validate-config' to verify")
	fmt.Fprintln(os.Stderr, "   • Restart SilentCast to apply changes")

	return nil
}

// importYAML imports configuration from YAML
func (c *ImportConfigCommand) importYAML(reader io.Reader) error {
	// Read and parse the YAML
	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	// Validate it's valid YAML
	var testConfig config.Config
	if err := yaml.Unmarshal(data, &testConfig); err != nil {
		return fmt.Errorf("invalid YAML configuration: %w", err)
	}

	// Get target path
	configDir := c.configPathFunc()
	configPath := filepath.Join(configDir, "spellbook.yml")

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Backup existing config if it exists
	if err := c.backupExistingConfig(configPath); err != nil {
		return fmt.Errorf("failed to backup existing config: %w", err)
	}

	// Write the new configuration
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	fmt.Fprintf(os.Stderr, "   • Configuration written to: %s\n", configPath)

	// Validate the imported configuration
	loader := config.NewLoader(configDir)
	if _, err := loader.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Warning: Imported configuration has validation errors:\n")
		fmt.Fprintf(os.Stderr, "   %v\n", err)
		fmt.Fprintf(os.Stderr, "   • Original config backed up as: %s.backup\n", configPath)
	}

	return nil
}

// importTarGz imports configuration from tar.gz archive
func (c *ImportConfigCommand) importTarGz(reader io.Reader) error {
	// Create gzip reader
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	// Create tar reader
	tarReader := tar.NewReader(gzReader)

	// Get target directory
	configDir := c.configPathFunc()

	// Process each file in the archive
	filesImported := 0
	for {
		header, err := tarReader.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		// Skip directories and non-config files
		if header.Typeflag != tar.TypeReg {
			continue
		}

		// Only import .yml files and skip metadata
		if !strings.HasSuffix(header.Name, ".yml") || header.Name == "EXPORT_INFO.txt" {
			continue
		}

		// Determine target path
		targetPath := filepath.Join(configDir, header.Name)

		// Ensure the file goes into the config directory
		if !strings.HasPrefix(filepath.Clean(targetPath), filepath.Clean(configDir)) {
			return fmt.Errorf("invalid file path in archive: %s", header.Name)
		}

		// Read file content
		content, err := io.ReadAll(tarReader)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", header.Name, err)
		}

		// Validate YAML
		var testConfig config.Config
		if err := yaml.Unmarshal(content, &testConfig); err != nil {
			fmt.Fprintf(os.Stderr, "   ⚠️  Skipping invalid YAML file: %s\n", header.Name)
			continue
		}

		// Backup existing file if it exists
		if err := c.backupExistingConfig(targetPath); err != nil {
			return fmt.Errorf("failed to backup %s: %w", targetPath, err)
		}

		// Write the file
		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", targetPath, err)
		}

		fmt.Fprintf(os.Stderr, "   • Imported: %s\n", header.Name)
		filesImported++
	}

	if filesImported == 0 {
		return fmt.Errorf("no valid configuration files found in archive")
	}

	// Validate the imported configuration
	loader := config.NewLoader(configDir)
	if _, err := loader.Load(); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Warning: Imported configuration has validation errors:\n")
		fmt.Fprintf(os.Stderr, "   %v\n", err)
	}

	return nil
}

// backupExistingConfig creates a backup of existing config file
func (c *ImportConfigCommand) backupExistingConfig(configPath string) error {
	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// No existing file to backup
		return nil
	}

	// Create backup with timestamp
	backupPath := fmt.Sprintf("%s.backup.%s", configPath, time.Now().Format("20060102-150405"))

	// Read existing file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read existing config: %w", err)
	}

	// Write backup
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write backup: %w", err)
	}

	fmt.Fprintf(os.Stderr, "   • Backed up existing config to: %s\n", filepath.Base(backupPath))
	return nil
}

// Group returns the command group
func (c *ImportConfigCommand) Group() string {
	return "config"
}

// HasOptions returns if this command has additional options
func (c *ImportConfigCommand) HasOptions() bool {
	return true
}
