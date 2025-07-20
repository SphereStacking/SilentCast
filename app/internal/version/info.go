package version

import (
	"encoding/json"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"
)

// BuildInfo contains comprehensive build and runtime information
type BuildInfo struct {
	// Version information
	Version   string `json:"version"`
	GitCommit string `json:"git_commit"`
	BuildTime string `json:"build_time"`
	
	// Go information
	GoVersion   string `json:"go_version"`
	GoArch      string `json:"go_arch"`
	GoOS        string `json:"go_os"`
	CGOEnabled  string `json:"cgo_enabled"`
	
	// Build configuration
	BuildTags   []string `json:"build_tags"`
	Compiler    string   `json:"compiler"`
	BuildMode   string   `json:"build_mode"`
	
	// Runtime information
	NumCPU    int `json:"num_cpu"`
	NumGoroutine int `json:"num_goroutine"`
}

// Build-time variables injected via ldflags
var (
	// Version is the semantic version
	Version = "dev"
	// GitCommit is the git commit hash
	GitCommit = "unknown"
	// BuildTime is the build timestamp
	BuildTime = "unknown"
	// BuildTags are the build tags used
	BuildTags = ""
)

// GetBuildInfo returns comprehensive build and runtime information
func GetBuildInfo() *BuildInfo {
	info := &BuildInfo{
		Version:   Version,
		GitCommit: GitCommit,
		BuildTime: BuildTime,
		
		GoVersion: runtime.Version(),
		GoArch:    runtime.GOARCH,
		GoOS:      runtime.GOOS,
		Compiler:  runtime.Compiler,
		
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
	}
	
	// Determine CGO status
	if isCGOEnabled() {
		info.CGOEnabled = "enabled"
	} else {
		info.CGOEnabled = "disabled"
	}
	
	// Parse build tags
	if BuildTags != "" {
		info.BuildTags = strings.Split(BuildTags, ",")
	} else {
		info.BuildTags = []string{}
	}
	
	// Try to get additional build info from debug.BuildInfo
	if buildInfo, ok := debug.ReadBuildInfo(); ok {
		// Check build settings for additional info
		for _, setting := range buildInfo.Settings {
			switch setting.Key {
			case "-tags":
				if BuildTags == "" {
					info.BuildTags = strings.Split(setting.Value, ",")
				}
			case "CGO_ENABLED":
				if setting.Value == "1" {
					info.CGOEnabled = "enabled"
				} else {
					info.CGOEnabled = "disabled"
				}
			case "-buildmode":
				info.BuildMode = setting.Value
			}
		}
	}
	
	return info
}

// isCGOEnabled checks if CGO is enabled at build time
func isCGOEnabled() bool {
	// This function will be optimized away if CGO is disabled
	// If CGO is enabled, the function exists and returns true
	// If CGO is disabled, the function doesn't exist and this returns false
	return _cgoEnabled()
}

// CGO detection is handled by build-tag specific files (cgo.go/nocgo.go)

// FormatHuman returns human-readable version information
func (info *BuildInfo) FormatHuman() string {
	var sb strings.Builder
	
	sb.WriteString("🪄 SilentCast - Silent hotkey-driven task runner\n\n")
	
	// Version information
	sb.WriteString("Version Information:\n")
	sb.WriteString(fmt.Sprintf("  Version:     %s\n", info.Version))
	sb.WriteString(fmt.Sprintf("  Git Commit:  %s\n", info.GitCommit))
	sb.WriteString(fmt.Sprintf("  Build Time:  %s\n", info.BuildTime))
	sb.WriteString("\n")
	
	// Go information
	sb.WriteString("Go Information:\n")
	sb.WriteString(fmt.Sprintf("  Go Version:  %s\n", info.GoVersion))
	sb.WriteString(fmt.Sprintf("  Platform:    %s/%s\n", info.GoOS, info.GoArch))
	sb.WriteString(fmt.Sprintf("  Compiler:    %s\n", info.Compiler))
	sb.WriteString(fmt.Sprintf("  CGO:         %s\n", info.CGOEnabled))
	if info.BuildMode != "" {
		sb.WriteString(fmt.Sprintf("  Build Mode:  %s\n", info.BuildMode))
	}
	sb.WriteString("\n")
	
	// Build configuration
	sb.WriteString("Build Configuration:\n")
	if len(info.BuildTags) > 0 {
		sb.WriteString(fmt.Sprintf("  Build Tags:  %s\n", strings.Join(info.BuildTags, ", ")))
	} else {
		sb.WriteString("  Build Tags:  none\n")
	}
	sb.WriteString("\n")
	
	// Runtime information
	sb.WriteString("Runtime Information:\n")
	sb.WriteString(fmt.Sprintf("  CPU Cores:   %d\n", info.NumCPU))
	sb.WriteString(fmt.Sprintf("  Goroutines:  %d\n", info.NumGoroutine))
	
	return sb.String()
}

// FormatJSON returns JSON-formatted version information
func (info *BuildInfo) FormatJSON() (string, error) {
	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal build info to JSON: %w", err)
	}
	return string(data), nil
}

// FormatCompact returns a compact version string
func (info *BuildInfo) FormatCompact() string {
	return fmt.Sprintf("SilentCast v%s (%s, %s/%s, %s)", 
		info.Version, info.GitCommit[:min(len(info.GitCommit), 8)], 
		info.GoOS, info.GoArch, info.CGOEnabled)
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetVersionString returns a simple version string for backward compatibility
func GetVersionString() string {
	return Version
}

// GetShortVersion returns a short version with commit hash
func GetShortVersion() string {
	if GitCommit != "unknown" && len(GitCommit) >= 8 {
		return fmt.Sprintf("%s-%s", Version, GitCommit[:8])
	}
	return Version
}