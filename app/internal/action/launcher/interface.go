package launcher

import "context"

// Launcher defines the interface for launching applications
type Launcher interface {
	// Launch starts an application with the given arguments
	Launch(ctx context.Context, appPath string, args []string) error
	
	// IsAvailable checks if the launcher can run on the current platform
	IsAvailable() bool
	
	// String returns a description of the launcher
	String() string
}

// AppResolver resolves application paths based on platform conventions
type AppResolver interface {
	// ResolvePath resolves an application name to its full path
	ResolvePath(appName string) (string, error)
	
	// FindExecutable searches for an executable in standard locations
	FindExecutable(name string) (string, error)
}