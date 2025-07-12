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
