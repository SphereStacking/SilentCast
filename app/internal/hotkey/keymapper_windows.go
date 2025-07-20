//go:build windows

package hotkey

func init() {
	keyMapperFactory = func() KeyMapper {
		return &windowsKeyMapper{}
	}
}

type windowsKeyMapper struct{}

func (w *windowsKeyMapper) GetModifierMap() map[string]string {
	return map[string]string{
		"win":     "win",
		"windows": "win",
		"ctrl":    "ctrl",
		"control": "ctrl",
		"alt":     "alt",
		"shift":   "shift",
	}
}

func (w *windowsKeyMapper) GetSpecialKeys() map[string]uint16 {
	// Windows uses standard virtual key codes
	return map[string]uint16{
		"printscreen": 0x2C,
		"scrolllock":  0x91,
		"pause":       0x13,
		"numlock":     0x90,
	}
}

func (w *windowsKeyMapper) GetKeyNameFromRawcode(rawcode uint16) string {
	// Windows virtual key code mapping
	keyMap := map[uint16]string{
		// Letters
		30: "a", 48: "b", 46: "c", 32: "d", 18: "e", 33: "f", 34: "g", 35: "h",
		23: "i", 36: "j", 37: "k", 38: "l", 50: "m", 49: "n", 24: "o", 25: "p",
		16: "q", 19: "r", 31: "s", 20: "t", 22: "u", 47: "v", 17: "w", 45: "x",
		21: "y", 44: "z",

		// Numbers
		11: "0", 2: "1", 3: "2", 4: "3", 5: "4", 6: "5", 7: "6", 8: "7", 9: "8", 10: "9",

		// Function keys
		59: "f1", 60: "f2", 61: "f3", 62: "f4", 63: "f5", 64: "f6",
		65: "f7", 66: "f8", 67: "f9", 68: "f10", 87: "f11", 88: "f12",

		// Special keys
		57: "space", 28: "enter", 15: "tab", 1: "esc", 14: "backspace",
		83: "delete", 82: "insert", 71: "home", 79: "end", 73: "pageup", 81: "pagedown",
		72: "up", 80: "down", 75: "left", 77: "right",

		// Numpad
		69: "numlock", 181: "divide", 55: "multiply", 74: "subtract", 78: "add",
		96: "numpad0", 97: "numpad1", 98: "numpad2", 99: "numpad3",
		100: "numpad4", 101: "numpad5", 102: "numpad6", 103: "numpad7",
		104: "numpad8", 105: "numpad9",

		// Others
		58: "capslock", 70: "scrolllock", 119: "pause",
		12: "minus", 13: "equal", 26: "leftbracket", 27: "rightbracket",
		39: "semicolon", 40: "apostrophe", 41: "grave", 43: "backslash",
		51: "comma", 52: "period", 191: "slash",
	}

	if name, ok := keyMap[rawcode]; ok {
		return name
	}

	return ""
}

func (w *windowsKeyMapper) IsModifierKey(rawcode uint16) (string, bool) {
	modifierMap := map[uint16]string{
		29:  "ctrl",  // Left Ctrl
		157: "ctrl",  // Right Ctrl
		42:  "shift", // Left Shift
		54:  "shift", // Right Shift
		56:  "alt",   // Left Alt
		184: "alt",   // Right Alt
		91:  "win",   // Left Windows
		92:  "win",   // Right Windows
	}

	if name, ok := modifierMap[rawcode]; ok {
		return name, true
	}

	return "", false
}
