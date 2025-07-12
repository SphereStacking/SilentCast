package hotkey

// KeyMapper defines the interface for platform-specific key mappings
type KeyMapper interface {
	// GetModifierMap returns the platform-specific modifier key mappings
	GetModifierMap() map[string]string

	// GetSpecialKeys returns any platform-specific special key codes
	GetSpecialKeys() map[string]uint16
}

// keyMapperFactory creates the appropriate key mapper for the current platform
var keyMapperFactory func() KeyMapper

// GetKeyMapper returns the platform-specific key mapper
func GetKeyMapper() KeyMapper {
	if keyMapperFactory == nil {
		panic("key mapper factory not initialized")
	}
	return keyMapperFactory()
}
