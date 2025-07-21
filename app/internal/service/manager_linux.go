//go:build linux

package service

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

const (
	serviceName        = "silentcast"
	serviceDisplayName = "SilentCast"
	serviceDescription = "Silent hotkey-driven task runner"

	// Systemd paths
	systemdUserPath   = ".config/systemd/user"
	systemdSystemPath = "/etc/systemd/system"

	// XDG autostart paths
	xdgAutostartPath = ".config/autostart"

	// Service file names
	systemdServiceFile = serviceName + ".service"
	desktopFile        = serviceName + ".desktop"
)

// systemdTemplate is the systemd service unit template
const systemdTemplate = `[Unit]
Description={{.Description}}
Documentation=https://github.com/SphereStacking/silentcast
After=graphical-session.target

[Service]
Type=simple
ExecStart={{.ExecStart}}
Restart=on-failure
RestartSec=5
Environment="PATH=/usr/local/bin:/usr/bin:/bin"
{{if .WorkingDirectory}}WorkingDirectory={{.WorkingDirectory}}{{end}}

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=read-only
ReadWritePaths={{.ConfigDir}} {{.LogDir}}

[Install]
WantedBy=default.target
`

// desktopTemplate is the XDG desktop entry template
const desktopTemplate = `[Desktop Entry]
Type=Application
Name={{.Name}}
Comment={{.Comment}}
Exec={{.Exec}}
Icon={{.Icon}}
Terminal=false
Categories=Utility;System;
StartupNotify=false
X-GNOME-Autostart-enabled=true
Hidden=false
`

// commandExecutor is the function type for executing commands
type commandExecutor func(name string, args ...string) ([]byte, error)

// LinuxManager implements service management for Linux
type LinuxManager struct {
	executable   string
	onRun        func() error
	useSystemd   bool
	isSystemWide bool
	execCommand  commandExecutor
}

// defaultCommandExecutor uses os/exec to run commands
func defaultCommandExecutor(name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	return cmd.Output()
}

// getCurrentUser wraps user.Current for testing
var getCurrentUser = user.Current

// NewManager creates a new Linux service manager
func NewManager(onRun func() error) Manager {
	exe, err := os.Executable()
	if err != nil {
		exe = os.Args[0]
	}

	// Resolve symlinks to get actual executable
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		exe = resolved
	}

	return &LinuxManager{
		executable:   exe,
		onRun:        onRun,
		useSystemd:   hasSystemd(),
		isSystemWide: false, // Default to user installation
		execCommand:  defaultCommandExecutor,
	}
}

// hasSystemd checks if systemd is available
func hasSystemd() bool {
	// Check if systemctl exists
	_, err := exec.LookPath("systemctl")
	if err != nil {
		return false
	}

	// Check if systemd is running
	cmd := exec.Command("systemctl", "--version")
	return cmd.Run() == nil
}

// Install installs the Linux service
func (m *LinuxManager) Install() error {
	var errors []error

	// Install systemd service if available
	if m.useSystemd {
		if err := m.installSystemdService(); err != nil {
			errors = append(errors, fmt.Errorf("systemd: %w", err))
		}
	}

	// Always install XDG autostart for GUI environments
	if err := m.installXDGAutostart(); err != nil {
		errors = append(errors, fmt.Errorf("XDG autostart: %w", err))
	}

	if len(errors) > 0 {
		return fmt.Errorf("installation errors: %v", errors)
	}

	return nil
}

// installSystemdService installs the systemd service
func (m *LinuxManager) installSystemdService() error {
	// Get current user
	currentUser, err := getCurrentUser()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	// Determine service path
	var servicePath string
	if m.isSystemWide {
		servicePath = filepath.Join(systemdSystemPath, systemdServiceFile)
	} else {
		userServiceDir := filepath.Join(currentUser.HomeDir, systemdUserPath)
		if mkdirErr := os.MkdirAll(userServiceDir, 0o755); mkdirErr != nil {
			return fmt.Errorf("failed to create systemd user directory: %w", mkdirErr)
		}
		servicePath = filepath.Join(userServiceDir, systemdServiceFile)
	}

	// Check if already installed
	if _, statErr := os.Stat(servicePath); statErr == nil {
		return fmt.Errorf("systemd service already installed at %s", servicePath)
	}

	// Prepare template data
	configDir := filepath.Join(currentUser.HomeDir, ".config", "silentcast")
	logDir := filepath.Join(currentUser.HomeDir, ".local", "share", "silentcast")

	data := struct {
		Description      string
		ExecStart        string
		WorkingDirectory string
		ConfigDir        string
		LogDir           string
	}{
		Description:      serviceDescription,
		ExecStart:        m.executable + " --no-tray",
		WorkingDirectory: filepath.Dir(m.executable),
		ConfigDir:        configDir,
		LogDir:           logDir,
	}

	// Generate service file from template
	tmpl, err := template.New("systemd").Parse(systemdTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse systemd template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute systemd template: %w", err)
	}

	// Write service file
	if err := os.WriteFile(servicePath, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("failed to write service file: %w", err)
	}

	// Reload systemd
	scope := "--user"
	if m.isSystemWide {
		scope = "--system"
	}

	if output, err := m.execCombinedOutput("systemctl", scope, "daemon-reload"); err != nil {
		os.Remove(servicePath)
		return fmt.Errorf("failed to reload systemd: %w\nOutput: %s", err, output)
	}

	// Enable service
	if output, err := m.execCombinedOutput("systemctl", scope, "enable", serviceName); err != nil {
		os.Remove(servicePath)
		return fmt.Errorf("failed to enable service: %w\nOutput: %s", err, output)
	}

	return nil
}

// installXDGAutostart installs XDG autostart entry
func (m *LinuxManager) installXDGAutostart() error {
	// Use HOME environment variable (can be overridden in tests)
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		currentUser, err := getCurrentUser()
		if err != nil {
			return fmt.Errorf("failed to get current user: %w", err)
		}
		homeDir = currentUser.HomeDir
	}

	// Create autostart directory
	autostartDir := filepath.Join(homeDir, xdgAutostartPath)
	if err := os.MkdirAll(autostartDir, 0o755); err != nil {
		return fmt.Errorf("failed to create autostart directory: %w", err)
	}

	desktopPath := filepath.Join(autostartDir, desktopFile)

	// Check if already installed
	if _, err := os.Stat(desktopPath); err == nil {
		return fmt.Errorf("XDG autostart already installed at %s", desktopPath)
	}

	// Prepare template data
	data := struct {
		Name    string
		Comment string
		Exec    string
		Icon    string
	}{
		Name:    serviceDisplayName,
		Comment: serviceDescription,
		Exec:    m.executable + " --no-tray",
		Icon:    serviceName,
	}

	// Generate desktop file from template
	tmpl, err := template.New("desktop").Parse(desktopTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse desktop template: %w", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute desktop template: %w", err)
	}

	// Write desktop file
	if err := os.WriteFile(desktopPath, buf.Bytes(), 0o644); err != nil {
		return fmt.Errorf("failed to write desktop file: %w", err)
	}

	return nil
}

// Uninstall removes the Linux service
func (m *LinuxManager) Uninstall() error {
	var errors []error

	// Uninstall systemd service
	if m.useSystemd {
		if err := m.uninstallSystemdService(); err != nil {
			errors = append(errors, fmt.Errorf("systemd: %w", err))
		}
	}

	// Uninstall XDG autostart
	if err := m.uninstallXDGAutostart(); err != nil {
		errors = append(errors, fmt.Errorf("XDG autostart: %w", err))
	}

	if len(errors) > 0 {
		return fmt.Errorf("uninstallation errors: %v", errors)
	}

	return nil
}

// uninstallSystemdService removes the systemd service
func (m *LinuxManager) uninstallSystemdService() error {
	currentUser, err := getCurrentUser()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	// Determine service path
	var servicePath string
	scope := "--user"
	if m.isSystemWide {
		servicePath = filepath.Join(systemdSystemPath, systemdServiceFile)
		scope = "--system"
	} else {
		servicePath = filepath.Join(currentUser.HomeDir, systemdUserPath, systemdServiceFile)
	}

	// Check if installed
	if _, err := os.Stat(servicePath); os.IsNotExist(err) {
		return fmt.Errorf("systemd service not installed")
	}

	// Stop service if running
	_ = m.execRun("systemctl", scope, "stop", serviceName) // Ignore error, service might not be running

	// Disable service
	if output, err := m.execCombinedOutput("systemctl", scope, "disable", serviceName); err != nil {
		fmt.Printf("Warning: failed to disable service: %v\nOutput: %s\n", err, output)
	}

	// Remove service file
	if err := os.Remove(servicePath); err != nil {
		return fmt.Errorf("failed to remove service file: %w", err)
	}

	// Reload systemd
	_ = m.execRun("systemctl", scope, "daemon-reload") // Ignore error, best effort cleanup

	return nil
}

// uninstallXDGAutostart removes XDG autostart entry
func (m *LinuxManager) uninstallXDGAutostart() error {
	// Use HOME environment variable (can be overridden in tests)
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		currentUser, err := getCurrentUser()
		if err != nil {
			return fmt.Errorf("failed to get current user: %w", err)
		}
		homeDir = currentUser.HomeDir
	}

	desktopPath := filepath.Join(homeDir, xdgAutostartPath, desktopFile)

	// Check if installed
	if _, err := os.Stat(desktopPath); os.IsNotExist(err) {
		return fmt.Errorf("XDG autostart not installed")
	}

	// Remove desktop file
	if err := os.Remove(desktopPath); err != nil {
		return fmt.Errorf("failed to remove desktop file: %w", err)
	}

	return nil
}

// Start starts the Linux service
func (m *LinuxManager) Start() error {
	if !m.useSystemd {
		return fmt.Errorf("systemd not available, cannot manage service")
	}

	scope := "--user"
	if m.isSystemWide {
		scope = "--system"
	}

	if output, err := m.execCombinedOutput("systemctl", scope, "start", serviceName); err != nil {
		return fmt.Errorf("failed to start service: %w\nOutput: %s", err, output)
	}

	return nil
}

// Stop stops the Linux service
func (m *LinuxManager) Stop() error {
	if !m.useSystemd {
		return fmt.Errorf("systemd not available, cannot manage service")
	}

	scope := "--user"
	if m.isSystemWide {
		scope = "--system"
	}

	if output, err := m.execCombinedOutput("systemctl", scope, "stop", serviceName); err != nil {
		return fmt.Errorf("failed to stop service: %w\nOutput: %s", err, output)
	}

	// Wait a bit for service to stop
	time.Sleep(2 * time.Second)

	return nil
}

// Status returns the current service status
func (m *LinuxManager) Status() (ServiceStatus, error) {
	result := ServiceStatus{
		Installed: false,
		Running:   false,
		StartType: "manual",
	}

	// Check systemd status
	if m.useSystemd {
		if status, err := m.getSystemdStatus(); err == nil {
			result = status
		}
	}

	// Check XDG autostart
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		if currentUser, err := getCurrentUser(); err == nil {
			homeDir = currentUser.HomeDir
		}
	}

	if homeDir != "" {
		desktopPath := filepath.Join(homeDir, xdgAutostartPath, desktopFile)
		if _, err := os.Stat(desktopPath); err == nil {
			result.Installed = true
			if !m.useSystemd {
				result.Message = "XDG autostart installed (no systemd)"
			}
		}
	}

	if !result.Installed {
		result.Message = "Service not installed"
	}

	return result, nil
}

// getSystemdStatus gets status from systemd
func (m *LinuxManager) getSystemdStatus() (ServiceStatus, error) {
	result := ServiceStatus{
		Installed: false,
		Running:   false,
		StartType: "manual",
	}

	scope := "--user"
	if m.isSystemWide {
		scope = "--system"
	}

	// Check if service exists
	output, err := m.execCommand("systemctl", scope, "list-unit-files", serviceName+".service")
	if err != nil || !strings.Contains(string(output), serviceName) {
		return result, fmt.Errorf("service not found")
	}

	result.Installed = true

	// Check if enabled
	if output, err := m.execCommand("systemctl", scope, "is-enabled", serviceName); err == nil {
		enabled := strings.TrimSpace(string(output))
		if enabled == "enabled" {
			result.StartType = "auto"
		}
	}

	// Check if active
	if output, err := m.execCommand("systemctl", scope, "is-active", serviceName); err == nil {
		active := strings.TrimSpace(string(output))
		if active == "active" {
			result.Running = true
			result.Message = "Service is running"
		} else {
			result.Message = fmt.Sprintf("Service is %s", active)
		}
	}

	return result, nil
}

// Run executes the service
func (m *LinuxManager) Run() error {
	// On Linux, just run the main function
	// Service management is handled by systemd
	return m.onRun()
}

// execCombinedOutput runs a command and returns combined output
func (m *LinuxManager) execCombinedOutput(name string, args ...string) ([]byte, error) {
	if m.execCommand != nil {
		// For testing, just use regular output
		return m.execCommand(name, args...)
	}
	cmd := exec.Command(name, args...)
	return cmd.CombinedOutput()
}

// execRun runs a command and returns error
func (m *LinuxManager) execRun(name string, args ...string) error {
	if m.execCommand != nil {
		// For testing, just check if command would succeed
		_, err := m.execCommand(name, args...)
		return err
	}
	cmd := exec.Command(name, args...)
	return cmd.Run()
}
