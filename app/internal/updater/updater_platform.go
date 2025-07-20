package updater

// PlatformUpdater defines platform-specific update operations
type PlatformUpdater interface {
	// CanReplaceRunningExecutable returns true if the platform allows
	// replacing a running executable
	CanReplaceRunningExecutable() bool

	// ReplaceExecutable performs platform-specific executable replacement
	ReplaceExecutable(src, dst string) error

	// MakeExecutable sets executable permissions on the file
	MakeExecutable(path string) error
	
	// RestartApplication restarts the application after update
	RestartApplication() error
}

// platformUpdaterFactory creates the appropriate updater for the current platform
var platformUpdaterFactory func() PlatformUpdater

// GetPlatformUpdater returns the platform-specific updater
func GetPlatformUpdater() PlatformUpdater {
	if platformUpdaterFactory == nil {
		panic("platform updater factory not initialized")
	}
	return platformUpdaterFactory()
}
