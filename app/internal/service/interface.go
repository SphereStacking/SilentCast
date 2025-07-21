package service

// Manager defines the interface for system service management
type Manager interface {
	// Install installs the service
	Install() error

	// Uninstall removes the service
	Uninstall() error

	// Start starts the service
	Start() error

	// Stop stops the service
	Stop() error

	// Status returns the current service status
	Status() (ServiceStatus, error)

	// Run executes the service (called by service manager)
	Run() error
}

// ServiceStatus represents the current state of the service
type ServiceStatus struct {
	Installed bool
	Running   bool
	StartType string // "auto", "manual", "disabled"
	Message   string
}
