//go:build linux

package hotkey

func init() {
	keyMapperFactory = func() KeyMapper {
		return &linuxKeyMapper{}
	}
}

type linuxKeyMapper struct{}

func (l *linuxKeyMapper) GetModifierMap() map[string]string {
	return map[string]string{
		// Control key
		"ctrl":    "ctrl",
		"control": "ctrl",

		// Alt key
		"alt":    "alt",
		"opt":    "alt",
		"option": "alt",

		// Shift key
		"shift": "shift",

		// Super/Windows/Command key
		"cmd":     "super",
		"command": "super",
		"win":     "super",
		"windows": "super",
		"meta":    "super",
		"super":   "super",
	}
}

func (l *linuxKeyMapper) GetSpecialKeys() map[string]uint16 {
	// Linux doesn't need special key codes in the same way
	// as the key names are handled differently by gohook
	return map[string]uint16{}
}

func (l *linuxKeyMapper) GetKeyNameFromRawcode(rawcode uint16) string {
	// Linux X11 keycodes (subtract 8 from X11 keycode to get rawcode)
	keyMap := map[uint16]string{
		// Letters
		38: "a", 56: "b", 54: "c", 40: "d", 26: "e", 41: "f", 42: "g", 43: "h",
		31: "i", 44: "j", 45: "k", 46: "l", 58: "m", 57: "n", 32: "o", 33: "p",
		24: "q", 27: "r", 39: "s", 28: "t", 30: "u", 55: "v", 25: "w", 53: "x",
		29: "y", 52: "z",

		// Numbers
		19: "0", 10: "1", 11: "2", 12: "3", 13: "4", 14: "5", 15: "6", 16: "7", 17: "8", 18: "9",

		// Function keys
		67: "f1", 68: "f2", 69: "f3", 70: "f4", 71: "f5", 72: "f6",
		73: "f7", 74: "f8", 75: "f9", 76: "f10", 95: "f11", 96: "f12",

		// Special keys
		65: "space", 36: "enter", 23: "tab", 9: "esc", 22: "backspace",
		119: "delete", 118: "insert", 110: "home", 115: "end", 112: "pageup", 117: "pagedown",
		111: "up", 116: "down", 113: "left", 114: "right",

		// Numpad
		77: "numlock", 106: "divide", 63: "multiply", 82: "subtract", 86: "add",
		90: "numpad0", 87: "numpad1", 88: "numpad2", 89: "numpad3",
		83: "numpad4", 84: "numpad5", 85: "numpad6", 79: "numpad7",
		80: "numpad8", 81: "numpad9",

		// Others
		66: "capslock", 78: "scrolllock", 127: "pause",
		20: "minus", 21: "equal", 34: "leftbracket", 35: "rightbracket",
		47: "semicolon", 48: "apostrophe", 49: "grave", 51: "backslash",
		59: "comma", 60: "period", 61: "slash",
	}

	if name, ok := keyMap[rawcode]; ok {
		return name
	}

	return ""
}

func (l *linuxKeyMapper) IsModifierKey(rawcode uint16) (string, bool) {
	modifierMap := map[uint16]string{
		37:  "ctrl",  // Left Ctrl
		105: "ctrl",  // Right Ctrl
		50:  "shift", // Left Shift
		62:  "shift", // Right Shift
		64:  "alt",   // Left Alt
		108: "alt",   // Right Alt
		133: "super", // Left Super (Windows key)
		134: "super", // Right Super
	}

	if name, ok := modifierMap[rawcode]; ok {
		return name, true
	}

	return "", false
}
