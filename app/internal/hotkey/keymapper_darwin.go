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