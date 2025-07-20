package config

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// ValidationError represents a single validation error with context
type ValidationError struct {
	Field      string // Field path (e.g., "grimoire.editor.command")
	Value      interface{}
	Message    string
	Suggestion string
	Line       int // Line number in YAML file
	Column     int // Column number in YAML file
}

func (e ValidationError) Error() string {
	if e.Line > 0 {
		return fmt.Sprintf("line %d: %s: %s", e.Line, e.Field, e.Message)
	}
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

// Validator performs comprehensive configuration validation
type Validator struct {
	config     *Config
	errors     []ValidationError
	yamlNode   *yaml.Node
	lineMapper map[string]int // Maps field paths to line numbers
}

// NewValidator creates a new configuration validator
func NewValidator() *Validator {
	return &Validator{
		errors:     make([]ValidationError, 0),
		lineMapper: make(map[string]int),
	}
}

// Validate performs comprehensive validation on the configuration
func (v *Validator) Validate(cfg *Config) []ValidationError {
	v.config = cfg
	v.errors = make([]ValidationError, 0)

	// Validate all sections
	v.validateHotkeys()
	v.validateDaemon()
	v.validateLogger()
	v.validateSpells()
	v.validateGrimoire()
	v.validateUpdater()

	return v.errors
}

// ValidateWithYAML performs validation with YAML node information for line numbers
func (v *Validator) ValidateWithYAML(cfg *Config, yamlContent []byte) []ValidationError {
	// Parse YAML to get node information
	var node yaml.Node
	if err := yaml.Unmarshal(yamlContent, &node); err == nil {
		v.yamlNode = &node
		v.buildLineMapper(&node, "")
	}

	return v.Validate(cfg)
}

// buildLineMapper recursively builds a map of field paths to line numbers
func (v *Validator) buildLineMapper(node *yaml.Node, prefix string) {
	if node == nil {
		return
	}

	switch node.Kind {
	case yaml.DocumentNode:
		if len(node.Content) > 0 {
			v.buildLineMapper(node.Content[0], prefix)
		}
	case yaml.MappingNode:
		for i := 0; i < len(node.Content); i += 2 {
			if i+1 < len(node.Content) {
				key := node.Content[i].Value
				fieldPath := key
				if prefix != "" {
					fieldPath = prefix + "." + key
				}
				v.lineMapper[fieldPath] = node.Content[i].Line
				v.buildLineMapper(node.Content[i+1], fieldPath)
			}
		}
	case yaml.SequenceNode:
		for i, item := range node.Content {
			fieldPath := fmt.Sprintf("%s[%d]", prefix, i)
			v.lineMapper[fieldPath] = item.Line
			v.buildLineMapper(item, fieldPath)
		}
	}
}

// addError adds a validation error
func (v *Validator) addError(field string, value interface{}, message string, suggestion string) {
	line := v.lineMapper[field]
	v.errors = append(v.errors, ValidationError{
		Field:      field,
		Value:      value,
		Message:    message,
		Suggestion: suggestion,
		Line:       line,
	})
}

// validateHotkeys validates hotkey configuration
func (v *Validator) validateHotkeys() {
	// Validate prefix key
	if v.config.Hotkeys.Prefix == "" && !v.config.prefixExplicitlySet {
		v.addError("hotkeys.prefix", "", "prefix key is required",
			"Add a prefix key like 'alt+space' or 'ctrl+shift+x'")
	}

	// Validate timeout values
	if v.config.Hotkeys.Timeout < 0 {
		v.addError("hotkeys.timeout", v.config.Hotkeys.Timeout,
			"timeout must be non-negative",
			"Use 0 for no timeout or a positive millisecond value")
	} else if v.config.Hotkeys.Timeout > 0 && v.config.Hotkeys.Timeout < 100 {
		v.addError("hotkeys.timeout", v.config.Hotkeys.Timeout,
			"timeout is very short (less than 100ms)",
			"Consider using a value between 500-2000ms for better usability")
	}

	// Validate sequence timeout
	if v.config.Hotkeys.SequenceTimeout < 0 {
		v.addError("hotkeys.sequence_timeout", v.config.Hotkeys.SequenceTimeout,
			"sequence timeout must be non-negative",
			"Use 0 for no timeout or a positive millisecond value")
	} else if v.config.Hotkeys.SequenceTimeout > 0 && v.config.Hotkeys.SequenceTimeout.ToDuration() < v.config.Hotkeys.Timeout.ToDuration() {
		v.addError("hotkeys.sequence_timeout", v.config.Hotkeys.SequenceTimeout,
			"sequence timeout should be greater than or equal to timeout",
			"Sequence timeout is the total time for multi-key sequences")
	}
}

// validateDaemon validates daemon configuration
func (v *Validator) validateDaemon() {
	// Validate log level
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if v.config.Daemon.LogLevel != "" && !validLogLevels[v.config.Daemon.LogLevel] {
		v.addError("daemon.log_level", v.config.Daemon.LogLevel,
			fmt.Sprintf("invalid log level '%s'", v.config.Daemon.LogLevel),
			"Use one of: debug, info, warn, error")
	}
}

// validateLogger validates logger configuration
func (v *Validator) validateLogger() {
	// Validate log file path
	if v.config.Logger.File != "" {
		expandedPath := os.ExpandEnv(v.config.Logger.File)
		dir := filepath.Dir(expandedPath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			v.addError("logger.file", v.config.Logger.File,
				"log file directory does not exist",
				fmt.Sprintf("Create directory '%s' or use a different path", dir))
		}
	}

	// Validate log sizes
	if v.config.Logger.MaxSize < 0 {
		v.addError("logger.max_size", v.config.Logger.MaxSize,
			"max size must be non-negative",
			"Use 0 for unlimited or a positive value in megabytes")
	}

	if v.config.Logger.MaxBackups < 0 {
		v.addError("logger.max_backups", v.config.Logger.MaxBackups,
			"max backups must be non-negative",
			"Use 0 for no backups or a positive value")
	}

	if v.config.Logger.MaxAge < 0 {
		v.addError("logger.max_age", v.config.Logger.MaxAge,
			"max age must be non-negative",
			"Use 0 for no age limit or a positive value in days")
	}
}

// validateSpells validates spell definitions
func (v *Validator) validateSpells() {
	// Check for prefix key conflicts
	if v.config.Hotkeys.Prefix != "" {
		prefixParts := strings.Split(v.config.Hotkeys.Prefix, "+")
		lastPart := prefixParts[len(prefixParts)-1]

		for spell := range v.config.Shortcuts {
			if spell == lastPart {
				v.addError(fmt.Sprintf("spells.%s", spell), spell,
					"spell conflicts with prefix key",
					"Use a different key or change the prefix key")
			}
		}
	}

	// Check for references to non-existent actions
	for spell, action := range v.config.Shortcuts {
		if _, exists := v.config.Actions[action]; !exists {
			v.addError(fmt.Sprintf("spells.%s", spell), action,
				fmt.Sprintf("references non-existent grimoire action '%s'", action),
				"Create the action in the grimoire section or fix the reference")
		}
	}

	// Validate spell key formats
	for spell := range v.config.Shortcuts {
		if err := validateSpellKey(spell); err != nil {
			v.addError(fmt.Sprintf("spells.%s", spell), spell,
				err.Error(),
				"Use single keys (e.g., 'e') or sequences (e.g., 'g,s')")
		}
	}
}

// validateGrimoire validates grimoire action definitions
func (v *Validator) validateGrimoire() {
	for name, action := range v.config.Actions {
		fieldPrefix := fmt.Sprintf("grimoire.%s", name)

		// Validate action type
		if action.Type == "" {
			v.addError(fieldPrefix+".type", "", "type is required",
				"Specify 'app', 'script', or 'url'")
			continue
		}

		validTypes := map[string]bool{"app": true, "script": true, "url": true}
		if !validTypes[action.Type] {
			v.addError(fieldPrefix+".type", action.Type,
				fmt.Sprintf("invalid type '%s'", action.Type),
				"Use 'app', 'script', or 'url'")
			continue
		}

		// Validate command
		if action.Command == "" {
			v.addError(fieldPrefix+".command", "", "command is required",
				"Specify the command, application path, or URL")
			continue
		}

		// Type-specific validation
		switch action.Type {
		case "app":
			v.validateAppAction(fieldPrefix, &action)
		case "script":
			v.validateScriptAction(fieldPrefix, &action)
		case "url":
			v.validateURLAction(fieldPrefix, &action)
		}

		// Validate common fields
		v.validateCommonActionFields(fieldPrefix, &action)
	}
}

// validateAppAction validates app-specific action fields
func (v *Validator) validateAppAction(fieldPrefix string, action *ActionConfig) {
	expandedCmd := os.ExpandEnv(action.Command)

	// Check if it's an absolute path
	if filepath.IsAbs(expandedCmd) {
		if _, err := os.Stat(expandedCmd); os.IsNotExist(err) {
			v.addError(fieldPrefix+".command", action.Command,
				"application file does not exist",
				"Check the file path or use a command in PATH")
		}
	} else {
		// Check if command exists in PATH
		if _, err := exec.LookPath(expandedCmd); err != nil {
			// Don't error for common commands that might be aliases
			commonCommands := map[string]bool{
				"code": true, "vim": true, "emacs": true, "nano": true,
				"chrome": true, "firefox": true, "safari": true,
			}
			if !commonCommands[expandedCmd] {
				v.addError(fieldPrefix+".command", action.Command,
					"application not found in PATH",
					"Use full path or ensure the application is in PATH")
			}
		}
	}
}

// validateScriptAction validates script-specific action fields
func (v *Validator) validateScriptAction(fieldPrefix string, action *ActionConfig) {
	// Validate shell if specified
	if action.Shell != "" {
		shellPath := action.Shell
		if !filepath.IsAbs(shellPath) {
			if path, err := exec.LookPath(shellPath); err == nil {
				shellPath = path
			} else {
				v.addError(fieldPrefix+".shell", action.Shell,
					"specified shell not found",
					fmt.Sprintf("Install %s or use a different shell", action.Shell))
			}
		} else if _, err := os.Stat(shellPath); os.IsNotExist(err) {
			v.addError(fieldPrefix+".shell", action.Shell,
				"specified shell does not exist",
				"Check the shell path")
		}
	}

	// Validate timeout
	if action.Timeout < 0 {
		v.addError(fieldPrefix+".timeout", action.Timeout,
			"timeout must be non-negative",
			"Use 0 for no timeout or a positive value in seconds")
	} else if action.Timeout > 3600 {
		v.addError(fieldPrefix+".timeout", action.Timeout,
			"timeout is very long (over 1 hour)",
			"Consider if such a long timeout is necessary")
	}

	// Validate working directory
	if action.WorkingDir != "" {
		expandedDir := os.ExpandEnv(action.WorkingDir)
		if !strings.Contains(expandedDir, "$") { // Skip if contains unresolved variables
			if _, err := os.Stat(expandedDir); os.IsNotExist(err) {
				v.addError(fieldPrefix+".working_dir", action.WorkingDir,
					"working directory does not exist",
					"Create the directory or use a different path")
			}
		}
	}
}

// validateURLAction validates URL-specific action fields
func (v *Validator) validateURLAction(fieldPrefix string, action *ActionConfig) {
	urlStr := strings.TrimSpace(action.Command)

	// Add scheme if missing
	if !strings.Contains(urlStr, "://") {
		urlStr = "https://" + urlStr
	}

	// Parse and validate URL
	u, err := url.Parse(urlStr)
	if err != nil {
		v.addError(fieldPrefix+".command", action.Command,
			"invalid URL format",
			"Use a valid URL like 'https://example.com'")
		return
	}

	// Validate scheme
	validSchemes := map[string]bool{
		"http": true, "https": true, "file": true,
		"ftp": true, "mailto": true,
	}

	if !validSchemes[u.Scheme] {
		v.addError(fieldPrefix+".command", action.Command,
			fmt.Sprintf("unsupported URL scheme '%s'", u.Scheme),
			"Use http, https, file, ftp, or mailto")
	}

	// Warn about localhost URLs in production
	if u.Hostname() == "localhost" || u.Hostname() == "127.0.0.1" {
		v.addError(fieldPrefix+".command", action.Command,
			"URL points to localhost",
			"This will only work on your local machine")
	}
}

// validateCommonActionFields validates fields common to all action types
func (v *Validator) validateCommonActionFields(fieldPrefix string, action *ActionConfig) {
	// Validate mutually exclusive options
	if action.ShowOutput && action.Terminal {
		v.addError(fieldPrefix, nil,
			"show_output and terminal are mutually exclusive",
			"Use either show_output for notifications or terminal for interactive display")
	}

	if action.Terminal && action.Type == "url" {
		v.addError(fieldPrefix+".terminal", true,
			"terminal option is not applicable for URL actions",
			"Remove the terminal option")
	}

	if action.KeepOpen && !action.Terminal {
		v.addError(fieldPrefix+".keep_open", true,
			"keep_open requires terminal to be true",
			"Add 'terminal: true' or remove keep_open")
	}

	if action.Admin && runtime.GOOS == "linux" {
		// Check for elevation tools on Linux
		tools := []string{"pkexec", "gksudo", "kdesudo", "sudo"}
		found := false
		for _, tool := range tools {
			if _, err := exec.LookPath(tool); err == nil {
				found = true
				break
			}
		}
		if !found {
			v.addError(fieldPrefix+".admin", true,
				"no elevation tool found for admin execution",
				"Install pkexec or gksudo for GUI elevation prompts")
		}
	}
}

// validateUpdater validates updater configuration
func (v *Validator) validateUpdater() {
	if !v.config.Updater.Enabled {
		return
	}

	// Validate check interval
	if v.config.Updater.CheckInterval != "" {
		duration, err := time.ParseDuration(v.config.Updater.CheckInterval)
		if err != nil {
			v.addError("updater.check_interval", v.config.Updater.CheckInterval,
				"invalid duration format",
				"Use format like '24h', '1h30m', or '90m'")
		} else if duration < time.Hour {
			v.addError("updater.check_interval", v.config.Updater.CheckInterval,
				"check interval is very short (less than 1 hour)",
				"Consider using a longer interval to reduce server load")
		}
	}
}

// validateSpellKey validates a spell key format
func validateSpellKey(key string) error {
	if key == "" {
		return fmt.Errorf("empty key")
	}

	// Check for sequence
	if strings.Contains(key, ",") {
		parts := strings.Split(key, ",")
		if len(parts) < 2 {
			return fmt.Errorf("invalid sequence format")
		}
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part == "" {
				return fmt.Errorf("empty key in sequence")
			}
			if strings.Contains(part, "+") {
				return fmt.Errorf("modifiers not allowed in sequences")
			}
		}
		return nil
	}

	// Single key validation
	parts := strings.Split(key, "+")
	if len(parts) > 1 {
		// Has modifiers
		validModifiers := map[string]bool{
			"ctrl": true, "alt": true, "shift": true, "cmd": true,
			"control": true, "option": true, "command": true, "super": true,
		}
		for i := 0; i < len(parts)-1; i++ {
			if !validModifiers[strings.ToLower(parts[i])] {
				return fmt.Errorf("invalid modifier '%s'", parts[i])
			}
		}
	}

	return nil
}
