package config

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"

	appErrors "github.com/SphereStacking/silentcast/internal/errors"
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

// Validate checks if the configuration files are valid
func (l *Loader) Validate() ([]string, error) {
	// Load configuration without defaults for validation
	cfg := &Config{
		Shortcuts: make(map[string]string),
		Actions:   make(map[string]ActionConfig),
	}

	hasConfig := false

	// Load each config file and merge for validation (preserves all values, including invalid ones)
	for _, path := range l.configPaths {
		if err := l.loadFileForValidation(path, cfg); err != nil {
			if !os.IsNotExist(err) {
				return nil, appErrors.Wrap(appErrors.ErrorTypeConfig, "failed to load configuration file", err).
					WithContext("path", path).
					WithContext("operation", "validation")
			}
			// File doesn't exist, skip it
			continue
		}
		hasConfig = true
	}

	if !hasConfig {
		return nil, appErrors.ErrConfigNotFound.WithContext("searched_paths", l.configPaths).
			WithContext("operation", "validation")
	}

	// Use the new comprehensive validator
	validator := NewValidator()
	validationErrors := validator.Validate(cfg)

	// Convert to string slice for backward compatibility
	var errors []string
	for _, e := range validationErrors {
		errors = append(errors, e.Error())
	}

	return errors, nil
}

// Load reads and merges configuration files
func (l *Loader) Load() (*Config, error) {
	// Start with empty config, defaults will be applied after loading
	cfg := &Config{
		Shortcuts: make(map[string]string),
		Actions:   make(map[string]ActionConfig),
	}

	hasConfig := false

	// Load each config file and merge
	for _, path := range l.configPaths {
		if err := l.loadFile(path, cfg); err != nil {
			if !os.IsNotExist(err) {
				return nil, appErrors.Wrap(appErrors.ErrorTypeConfig, "failed to load configuration file", err).
					WithContext("path", path).
					WithContext("operation", "load")
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
		return nil, appErrors.Wrap(appErrors.ErrorTypeValidation, "configuration validation failed", err).
			WithContext("config_paths", l.configPaths)
	}

	return cfg, nil
}

// LoadRaw loads configuration without validation or defaults (for export)
func (l *Loader) LoadRaw() (*Config, error) {
	// Start with empty config
	cfg := &Config{
		Shortcuts: make(map[string]string),
		Actions:   make(map[string]ActionConfig),
	}

	hasConfig := false

	// Load each config file and merge
	for _, path := range l.configPaths {
		if err := l.loadFile(path, cfg); err != nil {
			if !os.IsNotExist(err) {
				return nil, appErrors.Wrap(appErrors.ErrorTypeConfig, "failed to load configuration file", err).
					WithContext("path", path).
					WithContext("operation", "load_raw")
			}
			// File doesn't exist, skip it
			continue
		}
		hasConfig = true
	}

	if !hasConfig {
		return nil, appErrors.ErrConfigNotFound.WithContext("searched_paths", l.configPaths).
			WithContext("operation", "load_raw")
	}

	// Don't apply defaults or validate - return raw config
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
		return appErrors.Wrap(appErrors.ErrorTypeConfig, "failed to map custom keys", err).
			WithContext("path", path)
	}

	// Create a temporary config to load into
	temp := &Config{
		Shortcuts:           make(map[string]string),
		Actions:             make(map[string]ActionConfig),
		prefixExplicitlySet: hasPrefix,
	}

	if err := yaml.Unmarshal(mappedData, temp); err != nil {
		return appErrors.Wrap(appErrors.ErrorTypeConfig, "failed to parse YAML", err).
			WithContext("path", path)
	}

	// Merge the configurations
	l.merge(cfg, temp)

	return nil
}

// loadFileForValidation loads a config file for validation, preserving all values including invalid ones
func (l *Loader) loadFileForValidation(path string, cfg *Config) error {
	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Apply key mapping
	mappedData, hasPrefix, err := MapCustomKeys(data)
	if err != nil {
		return appErrors.Wrap(appErrors.ErrorTypeConfig, "failed to map custom keys", err).
			WithContext("path", path)
	}

	// Parse into temporary config
	temp := &Config{
		Shortcuts:           make(map[string]string),
		Actions:             make(map[string]ActionConfig),
		prefixExplicitlySet: hasPrefix,
	}

	if err := yaml.Unmarshal(mappedData, temp); err != nil {
		return appErrors.Wrap(appErrors.ErrorTypeConfig, "failed to parse YAML", err).
			WithContext("path", path)
	}

	// Merge the configurations using validation-specific merge
	l.mergeForValidation(cfg, temp)

	return nil
}

// mergeForValidation combines two configurations for validation, preserving all values including invalid ones
func (l *Loader) mergeForValidation(dst, src *Config) {
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

	// Merge hotkey config - preserve ALL values for validation, including negative ones
	if src.prefixExplicitlySet {
		dst.Hotkeys.Prefix = src.Hotkeys.Prefix
		dst.prefixExplicitlySet = true
	} else if src.Hotkeys.Prefix != "" {
		dst.Hotkeys.Prefix = src.Hotkeys.Prefix
	}
	// For validation, preserve ALL timeout values (even negative/zero)
	if src.Hotkeys.Timeout != 0 {
		dst.Hotkeys.Timeout = src.Hotkeys.Timeout
	}
	if src.Hotkeys.SequenceTimeout != 0 {
		dst.Hotkeys.SequenceTimeout = src.Hotkeys.SequenceTimeout
	}

	// Merge logger config
	if src.Logger.Level != "" {
		dst.Logger.Level = src.Logger.Level
	}
	if src.Logger.File != "" {
		dst.Logger.File = src.Logger.File
	}
	// For validation, preserve ALL values including negative ones
	if src.Logger.MaxSize != 0 {
		dst.Logger.MaxSize = src.Logger.MaxSize
	}
	if src.Logger.MaxBackups != 0 {
		dst.Logger.MaxBackups = src.Logger.MaxBackups
	}
	if src.Logger.MaxAge != 0 {
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
	for k := range src.Actions {
		dst.Actions[k] = src.Actions[k]
	}

	// Merge updater config
	if src.Updater.Enabled {
		dst.Updater.Enabled = src.Updater.Enabled
	}
	if src.Updater.CheckInterval != "" {
		dst.Updater.CheckInterval = src.Updater.CheckInterval
	}
	if src.Updater.Prerelease {
		dst.Updater.Prerelease = src.Updater.Prerelease
	}
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
	for k := range src.Actions {
		dst.Actions[k] = src.Actions[k]
	}
}

// validate checks if the configuration is valid
func (l *Loader) validate(cfg *Config) error {
	// Use the new comprehensive validator
	validator := NewValidator()
	errors := validator.Validate(cfg)

	if len(errors) > 0 {
		// Return the first error for backward compatibility
		// In the future, we might want to return all errors
		return appErrors.NewValidationError("config", errors[0].Error()).
			WithContext("total_errors", len(errors))
	}

	return nil
}
