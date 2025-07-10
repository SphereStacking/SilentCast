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