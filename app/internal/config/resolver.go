package config

// PlatformResolver defines the interface for platform-specific configuration
type PlatformResolver interface {
	// GetPlatformConfigFile returns the platform-specific config file name
	GetPlatformConfigFile() string
	
	// GetDefaultConfigPath returns the default configuration directory
	GetDefaultConfigPath() string
}

// platformFactory creates the appropriate platform resolver
var platformFactory func() PlatformResolver

// GetPlatformResolver returns the platform-specific resolver
func GetPlatformResolver() PlatformResolver {
	if platformFactory == nil {
		panic("platform factory not initialized")
	}
	return platformFactory()
}