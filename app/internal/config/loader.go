package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// Loader handles configuration loading and merging
type Loader struct {
	configPaths []string
}

// NewLoader creates a new configuration loader
func NewLoader(basePath string) *Loader {
	paths := []string{
		filepath.Join(basePath, ConfigName+".yml"),
	}
	
	// Add OS-specific config file using platform resolver
	platform := GetPlatformResolver()
	paths = append(paths, filepath.Join(basePath, platform.GetPlatformConfigFile()))
	
	return &Loader{
		configPaths: paths,
	}
}

// Load reads and merges configuration files
func (l *Loader) Load() (*Config, error) {
	// Start with empty config, defaults will be applied after loading
	cfg := &Config{
		Shortcuts:   make(map[string]string),
		Actions: make(map[string]ActionConfig),
	}
	
	hasConfig := false
	
	// Load each config file and merge
	for _, path := range l.configPaths {
		if err := l.loadFile(path, cfg); err != nil {
			if !os.IsNotExist(err) {
				return nil, fmt.Errorf("failed to load %s: %w", path, err)
			}
			// File doesn't exist, skip it
			continue
		}
		hasConfig = true
	}
	
	// Apply defaults only if config was loaded and values are not set
	if hasConfig {
		l.applyDefaults(cfg)
	}
	
	// Validate the configuration
	if err := l.validate(cfg); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}
	
	return cfg, nil
}

// applyDefaults sets default values for unset fields
func (l *Loader) applyDefaults(cfg *Config) {
	// Daemon defaults
	if cfg.Daemon.LogLevel == "" {
		cfg.Daemon.LogLevel = "info"
	}
	
	// Hotkey defaults
	// Only apply default prefix if it wasn't explicitly set (even to empty)
	if cfg.Hotkeys.Prefix == "" && !cfg.prefixExplicitlySet {
		cfg.Hotkeys.Prefix = "alt+space"
	}
	if cfg.Hotkeys.Timeout == 0 {
		cfg.Hotkeys.Timeout = Duration(1000 * time.Millisecond)
	}
	if cfg.Hotkeys.SequenceTimeout == 0 {
		cfg.Hotkeys.SequenceTimeout = Duration(2000 * time.Millisecond)
	}
	
	// Logger defaults
	if cfg.Logger.Level == "" {
		cfg.Logger.Level = "info"
	}
	if cfg.Logger.MaxSize == 0 {
		cfg.Logger.MaxSize = 10 // 10MB
	}
	if cfg.Logger.MaxBackups == 0 {
		cfg.Logger.MaxBackups = 3
	}
	if cfg.Logger.MaxAge == 0 {
		cfg.Logger.MaxAge = 7 // 7 days
	}
}

// loadFile reads a single configuration file and merges it into the config
func (l *Loader) loadFile(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	
	// Map custom keys to standard keys if needed
	mappedData, hasPrefix, err := MapCustomKeys(data)
	if err != nil {
		return fmt.Errorf("failed to map custom keys: %w", err)
	}
	
	// Create a temporary config to load into
	temp := &Config{
		Shortcuts:   make(map[string]string),
		Actions: make(map[string]ActionConfig),
		prefixExplicitlySet: hasPrefix,
	}
	
	if err := yaml.Unmarshal(mappedData, temp); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}
	
	// Merge the configurations
	l.merge(cfg, temp)
	
	return nil
}

// merge combines two configurations, with 'src' overriding 'dst'
func (l *Loader) merge(dst, src *Config) {
	// Merge daemon config
	if src.Daemon.LogLevel != "" {
		dst.Daemon.LogLevel = src.Daemon.LogLevel
	}
	if src.Daemon.AutoStart {
		dst.Daemon.AutoStart = src.Daemon.AutoStart
	}
	if !src.Daemon.ConfigWatch {
		dst.Daemon.ConfigWatch = src.Daemon.ConfigWatch
	}
	
	// Merge hotkey config
	// For prefix, we need to check if it was explicitly set
	// If a config has hotkeys section, it means prefix was configured (even if empty)
	if src.prefixExplicitlySet {
		dst.Hotkeys.Prefix = src.Hotkeys.Prefix
		dst.prefixExplicitlySet = true
	} else if src.Hotkeys.Prefix != "" {
		// Only update prefix if not empty and not explicitly set
		dst.Hotkeys.Prefix = src.Hotkeys.Prefix
	}
	if src.Hotkeys.Timeout > 0 {
		dst.Hotkeys.Timeout = src.Hotkeys.Timeout
	}
	if src.Hotkeys.SequenceTimeout > 0 {
		dst.Hotkeys.SequenceTimeout = src.Hotkeys.SequenceTimeout
	}
	
	// Merge logger config
	if src.Logger.Level != "" {
		dst.Logger.Level = src.Logger.Level
	}
	if src.Logger.File != "" {
		dst.Logger.File = src.Logger.File
	}
	if src.Logger.MaxSize > 0 {
		dst.Logger.MaxSize = src.Logger.MaxSize
	}
	if src.Logger.MaxBackups > 0 {
		dst.Logger.MaxBackups = src.Logger.MaxBackups
	}
	if src.Logger.MaxAge > 0 {
		dst.Logger.MaxAge = src.Logger.MaxAge
	}
	if src.Logger.Compress {
		dst.Logger.Compress = src.Logger.Compress
	}
	
	// Merge spells (overwrite)
	for k, v := range src.Shortcuts {
		dst.Shortcuts[k] = v
	}
	
	// Merge grimoire (overwrite)
	for k, v := range src.Actions {
		dst.Actions[k] = v
	}
}

// validate checks if the configuration is valid
func (l *Loader) validate(cfg *Config) error {
	// Check if any configuration was loaded
	if len(cfg.Shortcuts) == 0 && len(cfg.Actions) == 0 {
		// Allow empty config for now
	}
	
	// Validate prefix key - empty is allowed if explicitly set
	if cfg.Hotkeys.Prefix == "" && !cfg.prefixExplicitlySet {
		return fmt.Errorf("hotkeys.prefix is required")
	}
	
	// Validate spells reference existing grimoire entries
	for spell, action := range cfg.Shortcuts {
		if _, exists := cfg.Actions[action]; !exists {
			return fmt.Errorf("spell '%s' references non-existent grimoire action '%s'", spell, action)
		}
	}
	
	// Validate grimoire entries
	for name, action := range cfg.Actions {
		if action.Type == "" {
			return fmt.Errorf("grimoire action '%s' missing type", name)
		}
		if action.Type != "app" && action.Type != "script" {
			return fmt.Errorf("grimoire action '%s' has invalid type '%s' (must be 'app' or 'script')", name, action.Type)
		}
		if action.Command == "" {
			return fmt.Errorf("grimoire action '%s' missing command", name)
		}
	}
	
	// Validate log level
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLogLevels[cfg.Daemon.LogLevel] {
		return fmt.Errorf("invalid log level '%s'", cfg.Daemon.LogLevel)
	}
	
	return nil
}