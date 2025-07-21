package commands

import (
	"archive/tar"
	"compress/gzip"
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

// ExportConfigCommand exports configuration for backup/sharing
type ExportConfigCommand struct {
	configPathFunc       func() string
	configSearchPathFunc func() []string
}

// NewExportConfigCommand creates a new export config command
func NewExportConfigCommand(configPath func() string, configSearchPaths func() []string) Command {
	return &ExportConfigCommand{
		configPathFunc:       configPath,
		configSearchPathFunc: configSearchPaths,
	}
}

// Name returns the command name
func (c *ExportConfigCommand) Name() string {
	return "ExportConfig"
}

// Description returns the command description
func (c *ExportConfigCommand) Description() string {
	return "Export configuration for backup or sharing"
}

// FlagName returns the flag name
func (c *ExportConfigCommand) FlagName() string {
	return "export-config"
}

// IsActive checks if the command should run
func (c *ExportConfigCommand) IsActive(flags interface{}) bool {
	f, ok := flags.(*Flags)
	if !ok {
		return false
	}
	return f.ExportConfig != ""
}

// Execute runs the command
func (c *ExportConfigCommand) Execute(flags interface{}) error {
	f, ok := flags.(*Flags)
	if !ok {
		return fmt.Errorf("invalid flags type")
	}

	// Load configuration without validation (for export)
	configDir := c.configPathFunc()
	loader := config.NewLoader(configDir)
	cfg, err := loader.LoadRaw()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Determine output writer
	var writer io.Writer
	var closeFunc func() error

	if f.ExportConfig == "-" || f.ExportConfig == "" {
		// Export to stdout
		writer = os.Stdout
	} else {
		// Export to file
		file, fileErr := os.Create(f.ExportConfig)
		if fileErr != nil {
			return fmt.Errorf("failed to create export file: %w", fileErr)
		}
		writer = file
		closeFunc = file.Close
	}

	// Export based on format
	switch strings.ToLower(f.ExportFormat) {
	case "yaml":
		err = c.exportYAML(cfg, writer)
	case "tar.gz", "tgz", "targz":
		err = c.exportTarGz(cfg, configDir, writer)
	default:
		if closeFunc != nil {
			if closeErr := closeFunc(); closeErr != nil {
				logger.Warn("Failed to close export file: %v", closeErr)
			}
		}
		return fmt.Errorf("unsupported export format: %s (supported: yaml, tar.gz)", f.ExportFormat)
	}

	if err != nil {
		if closeFunc != nil {
			if closeErr := closeFunc(); closeErr != nil {
				logger.Warn("Failed to close export file: %v", closeErr)
			}
		}
		return fmt.Errorf("export failed: %w", err)
	}

	if closeFunc != nil {
		if err := closeFunc(); err != nil {
			return fmt.Errorf("failed to close export file: %w", err)
		}
	}

	// Print success message (only if not exporting to stdout)
	if f.ExportConfig != "-" && f.ExportConfig != "" {
		fmt.Fprintf(os.Stderr, "✅ Configuration exported to: %s\n", f.ExportConfig)
	}

	return nil
}

// exportYAML exports configuration as YAML
func (c *ExportConfigCommand) exportYAML(cfg *config.Config, writer io.Writer) error {
	encoder := yaml.NewEncoder(writer)
	encoder.SetIndent(2)

	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("failed to encode configuration: %w", err)
	}

	return encoder.Close()
}

// exportTarGz exports configuration as tar.gz archive
func (c *ExportConfigCommand) exportTarGz(_ *config.Config, configDir string, writer io.Writer) error {
	// Create gzip writer
	gzWriter := gzip.NewWriter(writer)
	defer gzWriter.Close()

	// Create tar writer
	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	// Add main config file
	mainConfigPath := filepath.Join(configDir, "spellbook.yml")
	if err := c.addFileToTar(tarWriter, mainConfigPath, "spellbook.yml"); err != nil {
		return err
	}

	// Check for OS-specific config files
	osConfigs := []string{
		filepath.Join(configDir, "spellbook.windows.yml"),
		filepath.Join(configDir, "spellbook.darwin.yml"),
		filepath.Join(configDir, "spellbook.linux.yml"),
	}

	for _, osConfig := range osConfigs {
		if _, err := os.Stat(osConfig); err == nil {
			baseName := filepath.Base(osConfig)
			if err := c.addFileToTar(tarWriter, osConfig, baseName); err != nil {
				return err
			}
		}
	}

	// Add metadata file with export info
	metadata := fmt.Sprintf("# SilentCast Configuration Export\n# Exported at: %s\n",
		time.Now().Format(time.RFC3339))

	metadataHeader := &tar.Header{
		Name: "EXPORT_INFO.txt",
		Mode: 0o644,
		Size: int64(len(metadata)),
	}

	if err := tarWriter.WriteHeader(metadataHeader); err != nil {
		return fmt.Errorf("failed to write metadata header: %w", err)
	}

	if _, err := tarWriter.Write([]byte(metadata)); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	return nil
}

// addFileToTar adds a file to the tar archive
func (c *ExportConfigCommand) addFileToTar(tw *tar.Writer, filePath, archiveName string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file %s: %w", filePath, err)
	}

	header := &tar.Header{
		Name:    archiveName,
		Mode:    int64(stat.Mode()),
		Size:    stat.Size(),
		ModTime: stat.ModTime(),
	}

	if err := tw.WriteHeader(header); err != nil {
		return fmt.Errorf("failed to write tar header for %s: %w", archiveName, err)
	}

	if _, err := io.Copy(tw, file); err != nil {
		return fmt.Errorf("failed to write file content for %s: %w", archiveName, err)
	}

	return nil
}

// Group returns the command group
func (c *ExportConfigCommand) Group() string {
	return "config"
}

// HasOptions returns if this command has additional options
func (c *ExportConfigCommand) HasOptions() bool {
	return true
}
