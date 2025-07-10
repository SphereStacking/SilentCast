package config

import (
	"gopkg.in/yaml.v3"
)

// MapCustomKeys converts custom key names to standard keys before unmarshaling
func MapCustomKeys(data []byte) (mappedData []byte, hasPrefix bool, err error) {
	// If using standard keys (same as Config struct YAML tags), return as-is
	if KeyDaemon == "daemon" && KeyHotkeys == "hotkeys" &&
		KeyShortcuts == "spells" && KeyActions == "grimoire" {
		// Check if prefix is set even with standard keys
		var raw map[string]interface{}
		if err := yaml.Unmarshal(data, &raw); err != nil {
			return data, false, nil // Return original on error
		}
		if hotkeys, ok := raw["hotkeys"].(map[string]interface{}); ok {
			_, hasPrefix = hotkeys["prefix"]
		}
		return data, hasPrefix, nil
	}
	
	// Parse YAML into a generic map
	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, false, err
	}
	
	// Check if prefix was explicitly set
	prefixSet := false
	if hotkeys, ok := raw[KeyHotkeys].(map[string]interface{}); ok {
		_, prefixSet = hotkeys["prefix"]
	}
	
	// Create new map with standard keys
	mapped := make(map[string]interface{})
	
	// Map custom keys to standard keys (Config struct YAML tags)
	keyMap := map[string]string{
		KeyDaemon:    "daemon",
		KeyHotkeys:   "hotkeys",
		KeyShortcuts: "spells",   // Config struct field: Spells
		KeyActions:   "grimoire", // Config struct field: Grimoire
		KeyLogger:    "logger",
		KeyUpdater:   "updater",
	}
	
	for customKey, standardKey := range keyMap {
		if value, exists := raw[customKey]; exists {
			mapped[standardKey] = value
		}
	}
	
	// Convert back to YAML
	result, err := yaml.Marshal(mapped)
	return result, prefixSet, err
}