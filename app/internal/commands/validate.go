package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/SphereStacking/silentcast/internal/config"
	"github.com/SphereStacking/silentcast/internal/notify"
	"github.com/SphereStacking/silentcast/internal/permission"
	"gopkg.in/yaml.v3"
)

// ValidateConfigCommand validates the configuration file
type ValidateConfigCommand struct {
	getConfigPath func() string
}

// NewValidateConfigCommand creates a new validate config command
func NewValidateConfigCommand(getConfigPath func() string) Command {
	return &ValidateConfigCommand{
		getConfigPath: getConfigPath,
	}
}

// Name returns the command name
func (c *ValidateConfigCommand) Name() string {
	return "Validate Config"
}

// Description returns the command description
func (c *ValidateConfigCommand) Description() string {
	return "Validate configuration and exit"
}

// FlagName returns the flag name
func (c *ValidateConfigCommand) FlagName() string {
	return "validate-config"
}

// IsActive checks if the command should run
func (c *ValidateConfigCommand) IsActive(flags interface{}) bool {
	f, ok := flags.(*Flags)
	if !ok {
		return false
	}
	return f.ValidateConfig
}

// Execute runs the command
func (c *ValidateConfigCommand) Execute(flags interface{}) error {
	// Type check flags
	f, ok := flags.(*Flags)
	if !ok {
		return fmt.Errorf("invalid flags type")
	}
	_ = f // We know it's the right type now

	configPath := c.getConfigPath()
	fmt.Printf("🔍 Validating configuration from: %s\n\n", configPath)

	// Create loader
	loader := config.NewLoader(configPath)

	// Validate configuration with enhanced error reporting
	validationErrors, err := c.validateWithLineNumbers(loader)
	if err != nil {
		fmt.Printf("❌ Configuration loading failed: %v\n", err)
		return err
	}

	if len(validationErrors) > 0 {
		fmt.Println("❌ Configuration validation failed:")
		for _, ve := range validationErrors {
			fmt.Printf("   • %s\n", ve)
		}
		return fmt.Errorf("configuration validation failed: %d errors found", len(validationErrors))
	}

	// Load configuration to check deeper issues
	cfg, err := loader.Load()
	if err != nil {
		fmt.Printf("❌ Configuration loading failed: %v\n", err)
		return err
	}

	// Check spell references
	fmt.Println("🔍 Checking spell references...")
	invalidSpells := 0
	for sequence, spellName := range cfg.Shortcuts {
		if _, exists := cfg.Actions[spellName]; !exists {
			fmt.Printf("   ❌ Spell '%s' references non-existent action '%s'\n", sequence, spellName)
			invalidSpells++
		} else {
			fmt.Printf("   ✅ %s → %s\n", sequence, spellName)
		}
	}

	if invalidSpells > 0 {
		return fmt.Errorf("%d invalid spell references found", invalidSpells)
	}

	// Check action configurations
	fmt.Println("\n🔍 Checking action configurations...")
	for name, action := range cfg.Actions {
		issues := validateAction(name, &action)
		if len(issues) > 0 {
			fmt.Printf("   ❌ Action '%s' has issues:\n", name)
			for _, issue := range issues {
				fmt.Printf("      • %s\n", issue)
			}
		} else {
			fmt.Printf("   ✅ %s (%s)\n", name, action.Type)
		}
	}

	// Check permissions
	fmt.Println("\n🔍 Checking permissions...")
	ctx := context.Background()
	permManager, err := permission.NewManager()
	if err != nil {
		fmt.Printf("   ⚠️  Could not check permissions: %v\n", err)
	} else {
		permissions, err := permManager.Check(ctx)
		if err != nil {
			fmt.Printf("   ⚠️  Permission check error: %v\n", err)
		} else {
			for _, perm := range permissions {
				if perm.Required {
					if perm.Status == permission.StatusGranted {
						fmt.Printf("   ✅ %s: Granted\n", perm.Type)
					} else {
						fmt.Printf("   ❌ %s: %s - %s\n", perm.Type, perm.Status, perm.Description)
					}
				}
			}
		}
	}

	// Check notification system
	fmt.Println("\n🔍 Checking notification system...")
	notifier := notify.NewManager()
	available := notifier.GetAvailableNotifiers()
	if len(available) > 0 {
		fmt.Printf("   ✅ Available notifiers: %s\n", strings.Join(available, ", "))
	} else {
		fmt.Printf("   ⚠️  No system notifiers available (will use console)\n")
	}

	// Summary
	fmt.Println("\n✅ Configuration is valid!")
	fmt.Printf("   • %d spells configured\n", len(cfg.Shortcuts))
	fmt.Printf("   • %d actions defined\n", len(cfg.Actions))
	fmt.Printf("   • Prefix key: %s\n", cfg.Hotkeys.Prefix)
	fmt.Printf("   • Config file: %s\n", filepath.Join(configPath, config.ConfigName+".yml"))

	return nil
}

// Group returns the command group
func (c *ValidateConfigCommand) Group() string {
	return "config"
}

// HasOptions returns if this command has additional options
func (c *ValidateConfigCommand) HasOptions() bool {
	return false
}

// validateAction checks an action configuration for issues
func validateAction(name string, action *config.ActionConfig) []string {
	var issues []string

	// Check required fields
	if action.Type == "" {
		issues = append(issues, "missing type")
	}
	if action.Command == "" {
		issues = append(issues, "missing command")
	}

	// Type-specific validation
	switch action.Type {
	case "":
		// Already handled above
	case "app":
		// App type validation
		if action.ShowOutput {
			issues = append(issues, "show_output is not applicable for app type")
		}
	case "script":
		// Script type validation
		if action.Shell != "" {
			// Validate shell exists
			validShells := []string{"bash", "sh", "zsh", "fish", "powershell", "pwsh", "cmd"}
			found := false
			for _, valid := range validShells {
				if action.Shell == valid {
					found = true
					break
				}
			}
			if !found {
				issues = append(issues, fmt.Sprintf("unknown shell: %s", action.Shell))
			}
		}
	case "url":
		// URL type validation
		if action.ShowOutput {
			issues = append(issues, "show_output is not applicable for url type")
		}
		if action.KeepOpen {
			issues = append(issues, "keep_open is not applicable for url type")
		}
	default:
		issues = append(issues, fmt.Sprintf("unknown type: %s", action.Type))
	}

	// Validate timeout
	if action.Timeout < 0 {
		issues = append(issues, "timeout cannot be negative")
	}

	// Validate conflicting options
	if action.ShowOutput && action.Terminal {
		issues = append(issues, "show_output and terminal options may conflict")
	}

	return issues
}

// validateWithLineNumbers validates configuration files with YAML line number reporting
func (c *ValidateConfigCommand) validateWithLineNumbers(loader *config.Loader) ([]string, error) {
	// Get the config paths from the loader
	configPaths := []string{
		filepath.Join(c.getConfigPath(), "spellbook.yml"),
		filepath.Join(c.getConfigPath(), config.GetPlatformResolver().GetPlatformConfigFile()),
	}

	var allErrors []string

	// Load configuration without defaults for validation  
	cfg := &config.Config{
		Shortcuts: make(map[string]string),
		Actions:   make(map[string]config.ActionConfig),
	}

	hasConfig := false

	// Validate each config file separately for better error reporting
	for _, path := range configPaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue // Skip non-existent files
		}

		// Read file content for line number validation
		yamlContent, err := os.ReadFile(path)
		if err != nil {
			allErrors = append(allErrors, fmt.Sprintf("%s: failed to read file: %v", filepath.Base(path), err))
			continue
		}

		// Create a temporary config for this file
		tempCfg := &config.Config{
			Shortcuts: make(map[string]string),
			Actions:   make(map[string]config.ActionConfig),
		}

		// Try to parse YAML syntax
		if err := yaml.Unmarshal(yamlContent, tempCfg); err != nil {
			allErrors = append(allErrors, fmt.Sprintf("%s: YAML syntax error: %v", filepath.Base(path), err))
			continue
		}

		// Use the validator with line number information
		validator := config.NewValidator()
		validationErrors := validator.ValidateWithYAML(tempCfg, yamlContent)

		// Add file context to each error
		for _, ve := range validationErrors {
			filename := filepath.Base(path)
			allErrors = append(allErrors, fmt.Sprintf("%s: %s", filename, ve.Error()))
		}

		// Merge this config into the main config
		mergeConfigs(cfg, tempCfg)
		hasConfig = true
	}

	if !hasConfig {
		return nil, fmt.Errorf("no configuration files found")
	}

	return allErrors, nil
}

// mergeConfigs merges source configuration into target
func mergeConfigs(target, source *config.Config) {
	// Merge shortcuts
	for key, value := range source.Shortcuts {
		target.Shortcuts[key] = value
	}

	// Merge actions
	for key, value := range source.Actions {
		target.Actions[key] = value
	}

	// Merge other fields if they're not zero values
	if source.Hotkeys.Prefix != "" {
		target.Hotkeys.Prefix = source.Hotkeys.Prefix
	}
	if source.Hotkeys.Timeout != 0 {
		target.Hotkeys.Timeout = source.Hotkeys.Timeout
	}
	if source.Hotkeys.SequenceTimeout != 0 {
		target.Hotkeys.SequenceTimeout = source.Hotkeys.SequenceTimeout
	}
}
