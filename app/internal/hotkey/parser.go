package hotkey

import (
	"fmt"
	"strings"
)

// Parser parses key sequence strings into KeySequence objects
type Parser struct {
	// Platform-specific key mappings
	keyMap      map[string]uint16
	modifierMap map[string]string
}

// NewParser creates a new key parser for the current platform
func NewParser() *Parser {
	p := &Parser{
		keyMap:      make(map[string]uint16),
		modifierMap: make(map[string]string),
	}

	// Initialize common key mappings
	p.initCommonKeys()

	// Initialize platform-specific mappings
	keyMapper := GetKeyMapper()
	p.modifierMap = keyMapper.GetModifierMap()

	// Merge in any platform-specific special keys
	for key, code := range keyMapper.GetSpecialKeys() {
		p.keyMap[key] = code
	}

	return p
}

// Parse parses a key sequence string into a KeySequence
func (p *Parser) Parse(sequence string) (KeySequence, error) {
	if sequence == "" {
		return KeySequence{}, ParseError{Input: sequence, Message: "empty sequence"}
	}

	// Split by comma for sequential keys
	parts := strings.Split(sequence, ",")
	keys := make([]Key, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			return KeySequence{}, ParseError{Input: sequence, Message: "empty key in sequence"}
		}

		key, err := p.parseKey(part)
		if err != nil {
			return KeySequence{}, err
		}

		keys = append(keys, key)
	}

	return KeySequence{Keys: keys}, nil
}

// parseKey parses a single key combination (e.g., "ctrl+a" or "a")
func (p *Parser) parseKey(keyStr string) (Key, error) {
	parts := strings.Split(strings.ToLower(keyStr), "+")

	if len(parts) == 0 {
		return Key{}, ParseError{Input: keyStr, Message: "empty key"}
	}

	// Last part is the main key
	mainKey := parts[len(parts)-1]
	modifiers := parts[:len(parts)-1]

	// Normalize modifiers
	normalizedMods := make([]string, 0, len(modifiers))
	for _, mod := range modifiers {
		normalized, ok := p.modifierMap[mod]
		if !ok {
			return Key{}, ParseError{Input: keyStr, Message: "unknown modifier: " + mod}
		}
		normalizedMods = append(normalizedMods, normalized)
	}

	// Get key code
	keyCode, ok := p.keyMap[mainKey]
	if !ok {
		// Check if it's a single letter or number
		if len(mainKey) == 1 {
			char := mainKey[0]
			if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') {
				keyCode = uint16(char)
			} else {
				return Key{}, ParseError{Input: keyStr, Message: "unknown key: " + mainKey}
			}
		} else {
			return Key{}, ParseError{Input: keyStr, Message: "unknown key: " + mainKey}
		}
	}

	return Key{
		Code:      keyCode,
		Modifiers: normalizedMods,
		Name:      mainKey,
	}, nil
}

// initCommonKeys initializes common key mappings
func (p *Parser) initCommonKeys() {
	// Letters (lowercase)
	for i := 'a'; i <= 'z'; i++ {
		p.keyMap[string(i)] = uint16(i)
	}

	// Numbers
	for i := '0'; i <= '9'; i++ {
		p.keyMap[string(i)] = uint16(i)
	}

	// Function keys
	for i := 1; i <= 12; i++ {
		key := fmt.Sprintf("f%d", i)
		p.keyMap[key] = uint16(0x70 + i - 1) // F1 = 0x70
	}

	// Special keys
	p.keyMap["space"] = 0x20
	p.keyMap["enter"] = 0x0D
	p.keyMap["return"] = 0x0D
	p.keyMap["tab"] = 0x09
	p.keyMap["esc"] = 0x1B
	p.keyMap["escape"] = 0x1B
	p.keyMap["backspace"] = 0x08
	p.keyMap["delete"] = 0x2E
	p.keyMap["insert"] = 0x2D
	p.keyMap["home"] = 0x24
	p.keyMap["end"] = 0x23
	p.keyMap["pageup"] = 0x21
	p.keyMap["pagedown"] = 0x22
	p.keyMap["up"] = 0x26
	p.keyMap["down"] = 0x28
	p.keyMap["left"] = 0x25
	p.keyMap["right"] = 0x27
}
