//go:build darwin

package commands

import (
	"fmt"
	"os"

	"github.com/SphereStacking/silentcast/internal/service"
)

// ServiceCommand handles macOS service management
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
		return "Install SilentCast as macOS LaunchAgent"
	case "uninstall":
		return "Uninstall SilentCast LaunchAgent"
	case "start":
		return "Start SilentCast LaunchAgent"
	case "stop":
		return "Stop SilentCast LaunchAgent"
	case "status":
		return "Show SilentCast LaunchAgent status"
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
	mgr := service.NewManager(c.onRun)

	switch c.action {
	case "install":
		fmt.Println("🔧 Installing SilentCast LaunchAgent...")
		
		// Check if already installed
		status, _ := mgr.Status()
		if status.Installed {
			return fmt.Errorf("service already installed")
		}
		
		if err := mgr.Install(); err != nil {
			return fmt.Errorf("failed to install service: %w", err)
		}
		
		fmt.Println("✅ LaunchAgent installed successfully")
		fmt.Println("   • Service name: com.spherestacking.silentcast")
		fmt.Println("   • Start type: Automatic (login)")
		fmt.Println("   • Logs: ~/Library/Logs/silentcast.log")
		fmt.Println("   • To start now: silentcast --service-start")
		
	case "uninstall":
		fmt.Println("🔧 Uninstalling SilentCast LaunchAgent...")
		
		// Check if sudo is needed for system-level service
		if c.needsSudo() {
			return fmt.Errorf("system-level service management requires sudo")
		}
		
		if err := mgr.Uninstall(); err != nil {
			return fmt.Errorf("failed to uninstall service: %w", err)
		}
		fmt.Println("✅ LaunchAgent uninstalled successfully")
		
	case "start":
		fmt.Println("🚀 Starting SilentCast LaunchAgent...")
		if err := mgr.Start(); err != nil {
			return fmt.Errorf("failed to start service: %w", err)
		}
		fmt.Println("✅ LaunchAgent started successfully")
		
	case "stop":
		fmt.Println("🛑 Stopping SilentCast LaunchAgent...")
		if err := mgr.Stop(); err != nil {
			return fmt.Errorf("failed to stop service: %w", err)
		}
		fmt.Println("✅ LaunchAgent stopped successfully")
		
	case "status":
		status, err := mgr.Status()
		if err != nil {
			return fmt.Errorf("failed to get service status: %w", err)
		}
		
		fmt.Println("📊 SilentCast LaunchAgent Status")
		fmt.Println("================================")
		fmt.Printf("Installed: %v\n", status.Installed)
		if status.Installed {
			fmt.Printf("Running:   %v\n", status.Running)
			fmt.Printf("Start type: %s\n", status.StartType)
			fmt.Printf("Plist: ~/Library/LaunchAgents/com.spherestacking.silentcast.plist\n")
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

// needsSudo checks if sudo is required
func (c *ServiceCommand) needsSudo() bool {
	// Check if we need elevated privileges
	// For user-level LaunchAgent, we don't need sudo
	// For system-level LaunchDaemon, we would need sudo
	return false
}

// isRoot checks if running as root
func isRoot() bool {
	return os.Geteuid() == 0
}