//go:build darwin

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
	serviceName        = "com.spherestacking.silentcast"
	serviceLabel       = "SilentCast"
	serviceDescription = "Silent hotkey-driven task runner"
	
	// LaunchAgent paths
	userAgentPath   = "~/Library/LaunchAgents"
	systemAgentPath = "/Library/LaunchAgents"
	daemonPath      = "/Library/LaunchDaemons"
)

// plistTemplate is the macOS launchd plist template
const plistTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>{{.Label}}</string>
    
    <key>ProgramArguments</key>
    <array>
        <string>{{.Executable}}</string>
        <string>--no-tray</string>
    </array>
    
    <key>RunAtLoad</key>
    <{{.RunAtLoad}}/>
    
    <key>KeepAlive</key>
    <dict>
        <key>SuccessfulExit</key>
        <false/>
        <key>Crashed</key>
        <true/>
    </dict>
    
    <key>ProcessType</key>
    <string>Interactive</string>
    
    <key>Nice</key>
    <integer>0</integer>
    
    <key>StandardOutPath</key>
    <string>{{.LogPath}}/silentcast.log</string>
    
    <key>StandardErrorPath</key>
    <string>{{.LogPath}}/silentcast.error.log</string>
    
    <key>WorkingDirectory</key>
    <string>{{.WorkingDirectory}}</string>
    
    {{if .IsSystemLevel}}
    <key>UserName</key>
    <string>{{.UserName}}</string>
    {{end}}
    
    <key>EnvironmentVariables</key>
    <dict>
        <key>PATH</key>
        <string>/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin</string>
    </dict>
    
    <key>LimitLoadToSessionType</key>
    <array>
        <string>Aqua</string>
        {{if .IsSystemLevel}}
        <string>LoginWindow</string>
        {{end}}
    </array>
</dict>
</plist>`

// DarwinManager implements service management for macOS
type DarwinManager struct {
	executable    string
	onRun         func() error
	isSystemLevel bool
}

// NewManager creates a new macOS service manager
func NewManager(onRun func() error) Manager {
	exe, err := os.Executable()
	if err != nil {
		exe = os.Args[0]
	}
	
	return &DarwinManager{
		executable:    exe,
		onRun:         onRun,
		isSystemLevel: false, // Default to user level
	}
}

// getPlistPath returns the path to the plist file
func (m *DarwinManager) getPlistPath() (string, error) {
	plistName := serviceName + ".plist"
	
	if m.isSystemLevel {
		return filepath.Join(systemAgentPath, plistName), nil
	}
	
	// User-level LaunchAgent
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}
	
	agentPath := strings.Replace(userAgentPath, "~", currentUser.HomeDir, 1)
	return filepath.Join(agentPath, plistName), nil
}

// Install installs the macOS service
func (m *DarwinManager) Install() error {
	// Ensure LaunchAgents directory exists
	plistPath, err := m.getPlistPath()
	if err != nil {
		return err
	}
	
	plistDir := filepath.Dir(plistPath)
	if err := os.MkdirAll(plistDir, 0755); err != nil {
		return fmt.Errorf("failed to create LaunchAgents directory: %w", err)
	}
	
	// Check if already installed
	if _, err := os.Stat(plistPath); err == nil {
		return fmt.Errorf("service %s already installed at %s", serviceName, plistPath)
	}
	
	// Get current user info
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}
	
	// Prepare template data
	logPath := filepath.Join(currentUser.HomeDir, "Library", "Logs")
	workingDir := filepath.Dir(m.executable)
	
	data := struct {
		Label             string
		Executable        string
		RunAtLoad         string
		LogPath           string
		WorkingDirectory  string
		IsSystemLevel     bool
		UserName          string
	}{
		Label:            serviceName,
		Executable:       m.executable,
		RunAtLoad:        "true",
		LogPath:          logPath,
		WorkingDirectory: workingDir,
		IsSystemLevel:    m.isSystemLevel,
		UserName:         currentUser.Username,
	}
	
	// Generate plist from template
	tmpl, err := template.New("plist").Parse(plistTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse plist template: %w", err)
	}
	
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("failed to execute plist template: %w", err)
	}
	
	// Write plist file
	if err := os.WriteFile(plistPath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("failed to write plist file: %w", err)
	}
	
	// Load the service
	cmd := exec.Command("launchctl", "load", plistPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		// Try to clean up the plist file
		os.Remove(plistPath)
		return fmt.Errorf("failed to load service: %w\nOutput: %s", err, output)
	}
	
	return nil
}

// Uninstall removes the macOS service
func (m *DarwinManager) Uninstall() error {
	plistPath, err := m.getPlistPath()
	if err != nil {
		return err
	}
	
	// Check if installed
	if _, err := os.Stat(plistPath); os.IsNotExist(err) {
		return fmt.Errorf("service %s not installed", serviceName)
	}
	
	// Unload the service
	cmd := exec.Command("launchctl", "unload", plistPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		// Log but continue - service might already be unloaded
		fmt.Printf("Warning: failed to unload service: %v\nOutput: %s\n", err, output)
	}
	
	// Remove from launchctl (for newer macOS versions)
	cmd = exec.Command("launchctl", "remove", serviceName)
	cmd.Run() // Ignore errors as this might fail on older macOS
	
	// Remove plist file
	if err := os.Remove(plistPath); err != nil {
		return fmt.Errorf("failed to remove plist file: %w", err)
	}
	
	return nil
}

// Start starts the macOS service
func (m *DarwinManager) Start() error {
	// Check if installed
	plistPath, err := m.getPlistPath()
	if err != nil {
		return err
	}
	
	if _, err := os.Stat(plistPath); os.IsNotExist(err) {
		return fmt.Errorf("service %s not installed", serviceName)
	}
	
	// Start the service
	cmd := exec.Command("launchctl", "start", serviceName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to start service: %w\nOutput: %s", err, output)
	}
	
	return nil
}

// Stop stops the macOS service
func (m *DarwinManager) Stop() error {
	// Stop the service
	cmd := exec.Command("launchctl", "stop", serviceName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to stop service: %w\nOutput: %s", err, output)
	}
	
	// Wait a bit for service to stop
	time.Sleep(2 * time.Second)
	
	return nil
}

// Status returns the current service status
func (m *DarwinManager) Status() (ServiceStatus, error) {
	result := ServiceStatus{
		Installed: false,
		Running:   false,
		StartType: "manual",
	}
	
	// Check if plist exists
	plistPath, err := m.getPlistPath()
	if err != nil {
		return result, err
	}
	
	if _, err := os.Stat(plistPath); os.IsNotExist(err) {
		result.Message = "Service not installed"
		return result, nil
	}
	
	result.Installed = true
	
	// Check if service is loaded
	cmd := exec.Command("launchctl", "list")
	output, err := cmd.Output()
	if err != nil {
		return result, fmt.Errorf("failed to list services: %w", err)
	}
	
	if strings.Contains(string(output), serviceName) {
		// Service is loaded, check if running
		cmd = exec.Command("launchctl", "print", "gui/"+fmt.Sprintf("%d", os.Getuid())+"/"+serviceName)
		output, err = cmd.Output()
		if err == nil {
			// Parse output to check state
			outputStr := string(output)
			if strings.Contains(outputStr, "state = running") {
				result.Running = true
				result.Message = "Service is running"
			} else if strings.Contains(outputStr, "state = waiting") {
				result.Message = "Service is loaded but not running"
			} else {
				result.Message = "Service is loaded with unknown state"
			}
			
			// Check if RunAtLoad is true
			if strings.Contains(outputStr, "RunAtLoad = true") {
				result.StartType = "auto"
			}
		} else {
			// Fallback: assume loaded but state unknown
			result.Message = "Service is loaded"
		}
	} else {
		result.Message = "Service is installed but not loaded"
	}
	
	return result, nil
}

// Run executes the service
func (m *DarwinManager) Run() error {
	// On macOS, when running as a LaunchAgent/Daemon, we just run normally
	// The service management is handled by launchd
	return m.onRun()
}