//go:build windows && !darwin && !linux

package service

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

const (
	serviceName        = "SilentCast"
	serviceDisplayName = "SilentCast Hotkey Service"
	serviceDescription = "Silent hotkey-driven task runner service"
)

// WindowsManager implements service management for Windows
type WindowsManager struct {
	executable string
	onRun      func() error
}

// NewManager creates a new Windows service manager
func NewManager(onRun func() error) Manager {
	exe, err := os.Executable()
	if err != nil {
		exe = os.Args[0]
	}
	
	return &WindowsManager{
		executable: exe,
		onRun:      onRun,
	}
}

// Install installs the Windows service
func (m *WindowsManager) Install() error {
	exepath, err := filepath.Abs(m.executable)
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	mgr, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to service manager: %w", err)
	}
	defer mgr.Disconnect()

	// Check if service already exists
	s, err := mgr.OpenService(serviceName)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", serviceName)
	}

	// Create the service
	config := mgr.Config{
		ServiceType:  windows.SERVICE_WIN32_OWN_PROCESS,
		StartType:    uint32(windows.SERVICE_AUTO_START),
		ErrorControl: uint32(windows.SERVICE_ERROR_NORMAL),
		DisplayName:  serviceDisplayName,
		Description:  serviceDescription,
	}

	s, err = mgr.CreateService(serviceName, exepath, config)
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}
	defer s.Close()

	// Recovery actions configuration can be added later if needed
	// The service will still work without automatic recovery actions

	// Create event log source
	err = eventlog.InstallAsEventCreate(serviceName, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		// Non-fatal error
		fmt.Printf("Warning: failed to create event log source: %v\n", err)
	}

	return nil
}

// Uninstall removes the Windows service
func (m *WindowsManager) Uninstall() error {
	mgr, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to service manager: %w", err)
	}
	defer mgr.Disconnect()

	s, err := mgr.OpenService(serviceName)
	if err != nil {
		return fmt.Errorf("service %s not found", serviceName)
	}
	defer s.Close()

	// Stop the service if running
	status, err := s.Query()
	if err == nil && status.State == svc.Running {
		_, err = s.Control(svc.Stop)
		if err != nil {
			fmt.Printf("Warning: failed to stop service: %v\n", err)
		}
		// Wait for service to stop
		time.Sleep(2 * time.Second)
	}

	// Delete the service
	err = s.Delete()
	if err != nil {
		return fmt.Errorf("failed to delete service: %w", err)
	}

	// Remove event log source
	eventlog.Remove(serviceName)

	return nil
}

// Start starts the Windows service
func (m *WindowsManager) Start() error {
	mgr, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to service manager: %w", err)
	}
	defer mgr.Disconnect()

	s, err := mgr.OpenService(serviceName)
	if err != nil {
		return fmt.Errorf("service %s not found", serviceName)
	}
	defer s.Close()

	err = s.Start()
	if err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	return nil
}

// Stop stops the Windows service
func (m *WindowsManager) Stop() error {
	mgr, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to service manager: %w", err)
	}
	defer mgr.Disconnect()

	s, err := mgr.OpenService(serviceName)
	if err != nil {
		return fmt.Errorf("service %s not found", serviceName)
	}
	defer s.Close()

	status, err := s.Control(svc.Stop)
	if err != nil {
		return fmt.Errorf("failed to stop service: %w", err)
	}

	if status.State != svc.Stopped {
		// Wait a bit for service to stop
		timeout := time.Now().Add(30 * time.Second)
		for time.Now().Before(timeout) {
			status, err = s.Query()
			if err != nil {
				return fmt.Errorf("failed to query service status: %w", err)
			}
			if status.State == svc.Stopped {
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
		
		if status.State != svc.Stopped {
			return fmt.Errorf("service did not stop within timeout")
		}
	}

	return nil
}

// Status returns the current service status
func (m *WindowsManager) Status() (ServiceStatus, error) {
	result := ServiceStatus{
		Installed: false,
		Running:   false,
	}

	mgr, err := mgr.Connect()
	if err != nil {
		return result, fmt.Errorf("failed to connect to service manager: %w", err)
	}
	defer mgr.Disconnect()

	s, err := mgr.OpenService(serviceName)
	if err != nil {
		result.Message = "Service not installed"
		return result, nil
	}
	defer s.Close()

	result.Installed = true

	config, err := s.Config()
	if err == nil {
		switch config.StartType {
		case uint32(windows.SERVICE_AUTO_START):
			result.StartType = "auto"
		case uint32(windows.SERVICE_DEMAND_START):
			result.StartType = "manual"
		case uint32(windows.SERVICE_DISABLED):
			result.StartType = "disabled"
		}
	}

	status, err := s.Query()
	if err != nil {
		return result, fmt.Errorf("failed to query service status: %w", err)
	}

	switch status.State {
	case svc.Running:
		result.Running = true
		result.Message = "Service is running"
	case svc.Stopped:
		result.Message = "Service is stopped"
	case svc.StartPending:
		result.Message = "Service is starting"
	case svc.StopPending:
		result.Message = "Service is stopping"
	case svc.Paused:
		result.Message = "Service is paused"
	default:
		result.Message = fmt.Sprintf("Service state: %v", status.State)
	}

	return result, nil
}

// Run executes the service
func (m *WindowsManager) Run() error {
	// Check if running interactively
	isIntSess, err := svc.IsAnInteractiveSession()
	if err != nil {
		return fmt.Errorf("failed to determine if we are running in an interactive session: %w", err)
	}

	if isIntSess {
		// Running interactively, just run the main function
		return m.onRun()
	}

	// Running as a service
	err = svc.Run(serviceName, &windowsService{
		onRun: m.onRun,
	})
	if err != nil {
		return fmt.Errorf("failed to run service: %w", err)
	}

	return nil
}

// windowsService implements svc.Handler
type windowsService struct {
	onRun func() error
}

// Execute is called by the service manager
func (ws *windowsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown

	changes <- svc.Status{State: svc.StartPending}

	// Start the main application in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- ws.onRun()
	}()

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	for {
		select {
		case err := <-errChan:
			if err != nil {
				// Log error to event log
				elog, _ := eventlog.Open(serviceName)
				if elog != nil {
					elog.Error(1, fmt.Sprintf("Service error: %v", err))
					elog.Close()
				}
				changes <- svc.Status{State: svc.StopPending}
				return false, 1
			}
			changes <- svc.Status{State: svc.StopPending}
			return false, 0

		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				changes <- svc.Status{State: svc.StopPending}
				// TODO: Signal the main application to stop gracefully
				return false, 0
			default:
				// Ignore unexpected commands
			}
		}
	}
}