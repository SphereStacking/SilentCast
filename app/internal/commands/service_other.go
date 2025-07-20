//go:build !windows && !darwin && !linux

package commands

import (
	"fmt"
	"runtime"
)

// Service commands are not available on non-Windows platforms
// These functions return nil to avoid compilation errors

func NewServiceInstallCommand(onRun func() error) Command {
	return &UnsupportedCommand{
		name: "service-install",
		msg:  fmt.Sprintf("Service management is only available on Windows (current: %s)", runtime.GOOS),
	}
}

func NewServiceUninstallCommand() Command {
	return &UnsupportedCommand{
		name: "service-uninstall",
		msg:  fmt.Sprintf("Service management is only available on Windows (current: %s)", runtime.GOOS),
	}
}

func NewServiceStartCommand() Command {
	return &UnsupportedCommand{
		name: "service-start",
		msg:  fmt.Sprintf("Service management is only available on Windows (current: %s)", runtime.GOOS),
	}
}

func NewServiceStopCommand() Command {
	return &UnsupportedCommand{
		name: "service-stop",
		msg:  fmt.Sprintf("Service management is only available on Windows (current: %s)", runtime.GOOS),
	}
}

func NewServiceStatusCommand() Command {
	return &UnsupportedCommand{
		name: "service-status",
		msg:  fmt.Sprintf("Service management is only available on Windows (current: %s)", runtime.GOOS),
	}
}

// UnsupportedCommand represents a command not available on this platform
type UnsupportedCommand struct {
	name string
	msg  string
}

func (c *UnsupportedCommand) Name() string        { return c.name }
func (c *UnsupportedCommand) Description() string { return c.msg }
func (c *UnsupportedCommand) FlagName() string    { return c.name }
func (c *UnsupportedCommand) Group() string       { return "service" }
func (c *UnsupportedCommand) HasOptions() bool    { return false }

func (c *UnsupportedCommand) IsActive(flags interface{}) bool {
	f, ok := flags.(*Flags)
	if !ok {
		return false
	}
	
	switch c.name {
	case "service-install":
		return f.ServiceInstall
	case "service-uninstall":
		return f.ServiceUninstall
	case "service-start":
		return f.ServiceStart
	case "service-stop":
		return f.ServiceStop
	case "service-status":
		return f.ServiceStatus
	}
	return false
}

func (c *UnsupportedCommand) Execute(flags interface{}) error {
	return fmt.Errorf(c.msg)
}