//go:build darwin

package hotkey

func init() {
	keyMapperFactory = func() KeyMapper {
		return &darwinKeyMapper{}
	}
}

type darwinKeyMapper struct{}

func (d *darwinKeyMapper) GetModifierMap() map[string]string {
	return map[string]string{
		"cmd":     "cmd",
		"command": "cmd",
		"ctrl":    "ctrl",
		"control": "ctrl",
		"alt":     "alt",
		"option":  "alt",
		"opt":     "alt",
		"shift":   "shift",
	}
}

func (d *darwinKeyMapper) GetSpecialKeys() map[string]uint16 {
	// macOS-specific key codes
	return map[string]uint16{
		"fn":       0x3F, // Function key
		"capslock": 0x39,
		"clear":    0x47,
		"help":     0x72,
	}
}

func (d *darwinKeyMapper) GetKeyNameFromRawcode(rawcode uint16) string {
	// macOS virtual key code mapping
	keyMap := map[uint16]string{
		// Letters
		0: "a", 11: "b", 8: "c", 2: "d", 14: "e", 3: "f", 5: "g", 4: "h",
		34: "i", 38: "j", 40: "k", 37: "l", 46: "m", 45: "n", 31: "o", 35: "p",
		12: "q", 15: "r", 1: "s", 17: "t", 32: "u", 9: "v", 13: "w", 7: "x",
		16: "y", 6: "z",
		
		// Numbers
		29: "0", 18: "1", 19: "2", 20: "3", 21: "4", 23: "5", 22: "6", 26: "7", 28: "8", 25: "9",
		
		// Function keys
		122: "f1", 120: "f2", 99: "f3", 118: "f4", 96: "f5", 97: "f6",
		98: "f7", 100: "f8", 101: "f9", 109: "f10", 103: "f11", 111: "f12",
		
		// Special keys
		49: "space", 36: "enter", 48: "tab", 53: "esc", 51: "backspace",
		117: "delete", 114: "insert", 115: "home", 119: "end", 116: "pageup", 121: "pagedown",
		126: "up", 125: "down", 123: "left", 124: "right",
		
		// Numpad
		71: "numlock", 75: "divide", 67: "multiply", 78: "subtract", 69: "add",
		82: "numpad0", 83: "numpad1", 84: "numpad2", 85: "numpad3",
		86: "numpad4", 87: "numpad5", 88: "numpad6", 89: "numpad7",
		91: "numpad8", 92: "numpad9",
		
		// Others
		57: "capslock", 107: "scrolllock", 113: "pause",
		27: "minus", 24: "equal", 33: "leftbracket", 30: "rightbracket",
		41: "semicolon", 39: "apostrophe", 50: "grave", 42: "backslash",
		43: "comma", 47: "period", 44: "slash",
	}
	
	if name, ok := keyMap[rawcode]; ok {
		return name
	}
	
	return ""
}

func (d *darwinKeyMapper) IsModifierKey(rawcode uint16) (string, bool) {
	modifierMap := map[uint16]string{
		59:  "ctrl",  // Left Ctrl
		62:  "ctrl",  // Right Ctrl
		56:  "shift", // Left Shift
		60:  "shift", // Right Shift
		58:  "alt",   // Left Alt (Option)
		61:  "alt",   // Right Alt (Option)
		55:  "cmd",   // Left Command
		54:  "cmd",   // Right Command
	}
	
	if name, ok := modifierMap[rawcode]; ok {
		return name, true
	}
	
	return "", false
}
