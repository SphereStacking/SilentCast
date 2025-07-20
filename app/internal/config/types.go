package config

import (
	"time"

	"github.com/SphereStacking/silentcast/internal/terminal"
	"gopkg.in/yaml.v3"
)

// Config represents the main configuration structure
type Config struct {
	Daemon       DaemonConfig            `yaml:"daemon"`
	Hotkeys      HotkeyConfig            `yaml:"hotkeys"`
	Shortcuts    map[string]string       `yaml:"spells"`   // YAMLでは"spells"だがコードではShortcuts
	Actions      map[string]ActionConfig `yaml:"grimoire"` // YAMLでは"grimoire"だがコードではActions
	Logger       LoggerConfig            `yaml:"logger"`
	Updater      UpdaterConfig           `yaml:"updater"`
	Notification NotificationConfig      `yaml:"notification"`
	Performance  PerformanceConfig       `yaml:"performance"`

	// Internal fields (not from YAML)
	prefixExplicitlySet bool `yaml:"-"`
}

// UnmarshalYAML implements yaml.Unmarshaler for Config
func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Use an alias to avoid infinite recursion
	type configAlias Config
	alias := (*configAlias)(c)

	// First unmarshal into a map to check what fields are present
	var raw map[string]interface{}
	if err := unmarshal(&raw); err != nil {
		return err
	}

	// Check if hotkeys.prefix was explicitly set
	if hotkeys, ok := raw["hotkeys"].(map[string]interface{}); ok {
		if _, hasPrefix := hotkeys["prefix"]; hasPrefix {
			c.prefixExplicitlySet = true
		}
	}

	// Now unmarshal the full structure
	// Re-encode the map to YAML and decode into our struct
	data, err := yaml.Marshal(raw)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, alias)
}

// DaemonConfig contains daemon-related settings
type DaemonConfig struct {
	AutoStart   bool   `yaml:"auto_start"`
	LogLevel    string `yaml:"log_level"`
	ConfigWatch bool   `yaml:"config_watch"`
}

// LoggerConfig contains logger-related settings
type LoggerConfig struct {
	Level      string `yaml:"level"`       // debug, info, warn, error
	File       string `yaml:"file"`        // log file path
	MaxSize    int    `yaml:"max_size"`    // megabytes
	MaxBackups int    `yaml:"max_backups"` // number of old files to keep
	MaxAge     int    `yaml:"max_age"`     // days
	Compress   bool   `yaml:"compress"`    // compress old files
}

// UpdaterConfig contains auto-update settings
type UpdaterConfig struct {
	Enabled       bool   `yaml:"enabled"`        // Enable auto-update checks
	CheckInterval string `yaml:"check_interval"` // How often to check (e.g., "24h")
	AutoInstall   bool   `yaml:"auto_install"`   // Automatically install updates
	Prerelease    bool   `yaml:"prerelease"`     // Include pre-release versions
}

// NotificationConfig contains notification-related settings
type NotificationConfig struct {
	EnableTimeout   bool `yaml:"enable_timeout,omitempty"`    // Enable timeout notifications (default: true)
	EnableWarning   bool `yaml:"enable_warning,omitempty"`    // Enable warning before timeout (default: true)
	Sound           bool `yaml:"sound,omitempty"`             // Play sound for notifications (default: true)
	MaxOutputLength int  `yaml:"max_output_length,omitempty"` // Max output length in notifications (default: 1024)
}

// HotkeyConfig contains hotkey-related settings
type HotkeyConfig struct {
	Prefix          string   `yaml:"prefix"`
	Timeout         Duration `yaml:"timeout"`
	SequenceTimeout Duration `yaml:"sequence_timeout"`
}

// Duration is a wrapper around time.Duration for YAML unmarshaling
type Duration time.Duration

// UnmarshalYAML implements yaml.Unmarshaler for Duration
func (d *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var ms int
	if err := unmarshal(&ms); err != nil {
		return err
	}
	*d = Duration(time.Duration(ms) * time.Millisecond)
	return nil
}

// ToDuration converts Duration to time.Duration
func (d Duration) ToDuration() time.Duration {
	return time.Duration(d)
}

// ActionConfig represents an action that can be executed
type ActionConfig struct {
	Type        string            `yaml:"type"`    // "app", "script", or "url"
	Command     string            `yaml:"command"` // Path or command
	Args        []string          `yaml:"args,omitempty"`
	Env         map[string]string `yaml:"env,omitempty"`
	WorkingDir  string            `yaml:"working_dir,omitempty"`
	Description string            `yaml:"description,omitempty"`

	// Output control
	ShowOutput bool `yaml:"show_output,omitempty"` // Show command output as notification
	KeepOpen   bool `yaml:"keep_open,omitempty"`   // Keep terminal open after execution

	// Execution control
	Timeout        int    `yaml:"timeout,omitempty"`         // Timeout in seconds (0 = no timeout)
	GracePeriod    int    `yaml:"grace_period,omitempty"`    // Grace period before SIGKILL in seconds (default: 5)
	TimeoutWarning int    `yaml:"timeout_warning,omitempty"` // Warning before timeout in seconds (0 = no warning)
	Shell          string `yaml:"shell,omitempty"`           // Custom shell to use
	Interpreter    bool   `yaml:"interpreter,omitempty"`     // Run directly in interpreter mode (bypass shell)
	Admin          bool   `yaml:"admin,omitempty"`           // Run with elevated privileges
	Terminal       bool   `yaml:"terminal,omitempty"`        // Force run in terminal
	ForceTerminal  bool   `yaml:"force_terminal,omitempty"`  // Force terminal even in GUI/tray mode
	
	// Terminal customization
	TerminalCustomization *terminal.Customization `yaml:"terminal_customization,omitempty"` // Visual customization for terminal window
}
