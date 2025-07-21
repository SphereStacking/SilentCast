//go:build windows && !darwin && !linux

package commands

import (
	"fmt"
	"os"

	"github.com/SphereStacking/silentcast/internal/service"
)

// ServiceCommand handles Windows service management
type ServiceCommand struct {
	action string
	onRun  func() error
}

// NewServiceInstallCommand creates a new service install command
func NewServiceInstallCommand(onRun func() error) Command {
	return &ServiceCommand{
		action: "install",
		onRun:  onRun,
	}
}

// NewServiceUninstallCommand creates a new service uninstall command
func NewServiceUninstallCommand() Command {
	return &ServiceCommand{
		action: "uninstall",
	}
}

// NewServiceStartCommand creates a new service start command
func NewServiceStartCommand() Command {
	return &ServiceCommand{
		action: "start",
	}
}

// NewServiceStopCommand creates a new service stop command
func NewServiceStopCommand() Command {
	return &ServiceCommand{
		action: "stop",
	}
}

// NewServiceStatusCommand creates a new service status command
func NewServiceStatusCommand() Command {
	return &ServiceCommand{
		action: "status",
	}
}

// Name returns the command name
func (c *ServiceCommand) Name() string {
	return fmt.Sprintf("Service%s", c.action)
}

// Description returns the command description
func (c *ServiceCommand) Description() string {
	switch c.action {
	case "install":
		return "Install SilentCast as Windows service"
	case "uninstall":
		return "Uninstall SilentCast Windows service"
	case "start":
		return "Start SilentCast Windows service"
	case "stop":
		return "Stop SilentCast Windows service"
	case "status":
		return "Show SilentCast Windows service status"
	default:
		return "Unknown service command"
	}
}

// FlagName returns the flag name
func (c *ServiceCommand) FlagName() string {
	return fmt.Sprintf("service-%s", c.action)
}

// IsActive checks if the command should run
func (c *ServiceCommand) IsActive(flags interface{}) bool {
	f, ok := flags.(*Flags)
	if !ok {
		return false
	}

	switch c.action {
	case "install":
		return f.ServiceInstall
	case "uninstall":
		return f.ServiceUninstall
	case "start":
		return f.ServiceStart
	case "stop":
		return f.ServiceStop
	case "status":
		return f.ServiceStatus
	default:
		return false
	}
}

// Execute runs the command
func (c *ServiceCommand) Execute(flags interface{}) error {
	// Check if running as administrator
	if !isAdmin() {
		return fmt.Errorf("service management requires administrator privileges")
	}

	mgr := service.NewManager(c.onRun)

	switch c.action {
	case "install":
		fmt.Println("🔧 Installing SilentCast service...")
		if err := mgr.Install(); err != nil {
			return fmt.Errorf("failed to install service: %w", err)
		}
		fmt.Println("✅ Service installed successfully")
		fmt.Println("   • Service name: SilentCast")
		fmt.Println("   • Start type: Automatic")
		fmt.Println("   • To start the service: silentcast --service-start")

	case "uninstall":
		fmt.Println("🔧 Uninstalling SilentCast service...")
		if err := mgr.Uninstall(); err != nil {
			return fmt.Errorf("failed to uninstall service: %w", err)
		}
		fmt.Println("✅ Service uninstalled successfully")

	case "start":
		fmt.Println("🚀 Starting SilentCast service...")
		if err := mgr.Start(); err != nil {
			return fmt.Errorf("failed to start service: %w", err)
		}
		fmt.Println("✅ Service started successfully")

	case "stop":
		fmt.Println("🛑 Stopping SilentCast service...")
		if err := mgr.Stop(); err != nil {
			return fmt.Errorf("failed to stop service: %w", err)
		}
		fmt.Println("✅ Service stopped successfully")

	case "status":
		status, err := mgr.Status()
		if err != nil {
			return fmt.Errorf("failed to get service status: %w", err)
		}

		fmt.Println("📊 SilentCast Service Status")
		fmt.Println("============================")
		fmt.Printf("Installed: %v\n", status.Installed)
		if status.Installed {
			fmt.Printf("Running:   %v\n", status.Running)
			fmt.Printf("Start type: %s\n", status.StartType)
		}
		fmt.Printf("Status:    %s\n", status.Message)
	}

	return nil
}

// Group returns the command group
func (c *ServiceCommand) Group() string {
	return "service"
}

// HasOptions returns if this command has additional options
func (c *ServiceCommand) HasOptions() bool {
	return false
}

// isAdmin checks if running with administrator privileges
func isAdmin() bool {
	// Try to open the service manager - this requires admin rights
	file, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	file.Close()
	return true
}
