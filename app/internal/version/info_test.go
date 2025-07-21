package version

import (
	"encoding/json"
	"runtime"
	"strings"
	"testing"
)

func TestGetBuildInfo(t *testing.T) {
	info := GetBuildInfo()

	// Test basic fields are populated
	if info == nil {
		t.Fatal("GetBuildInfo() returned nil")
	}

	if info.Version == "" {
		t.Error("Version should not be empty")
	}

	if info.GoVersion == "" {
		t.Error("GoVersion should not be empty")
	}

	if info.GoArch == "" {
		t.Error("GoArch should not be empty")
	}

	if info.GoOS == "" {
		t.Error("GoOS should not be empty")
	}

	if info.Compiler == "" {
		t.Error("Compiler should not be empty")
	}

	// Test runtime information
	if info.NumCPU <= 0 {
		t.Error("NumCPU should be positive")
	}

	if info.NumGoroutine <= 0 {
		t.Error("NumGoroutine should be positive")
	}

	// Test CGO status
	if info.CGOEnabled != "enabled" && info.CGOEnabled != "disabled" {
		t.Errorf("CGOEnabled should be 'enabled' or 'disabled', got: %s", info.CGOEnabled)
	}

	// Test build tags (should be slice, even if empty)
	if info.BuildTags == nil {
		t.Error("BuildTags should not be nil")
	}
}

func TestBuildInfoWithBuildTags(t *testing.T) {
	// Save original value
	originalBuildTags := BuildTags
	defer func() {
		BuildTags = originalBuildTags
	}()

	// Test with build tags
	BuildTags = "nogohook,notray,test"
	info := GetBuildInfo()

	expectedTags := []string{"nogohook", "notray", "test"}
	if len(info.BuildTags) != len(expectedTags) {
		t.Errorf("Expected %d build tags, got %d", len(expectedTags), len(info.BuildTags))
	}

	for i, tag := range expectedTags {
		if i >= len(info.BuildTags) || info.BuildTags[i] != tag {
			t.Errorf("Expected build tag %s at index %d, got %s", tag, i, info.BuildTags[i])
		}
	}
}

func TestBuildInfoWithEmptyBuildTags(t *testing.T) {
	// Save original value
	originalBuildTags := BuildTags
	defer func() {
		BuildTags = originalBuildTags
	}()

	// Test with empty build tags
	BuildTags = ""
	info := GetBuildInfo()

	if len(info.BuildTags) != 0 {
		t.Errorf("Expected empty build tags, got %v", info.BuildTags)
	}
}

func TestFormatHuman(t *testing.T) {
	info := GetBuildInfo()
	human := info.FormatHuman()

	// Test that output contains expected sections
	expectedSections := []string{
		"SilentCast",
		"Version Information:",
		"Go Information:",
		"Build Configuration:",
		"Runtime Information:",
		"Version:",
		"Go Version:",
		"Platform:",
		"CPU Cores:",
	}

	for _, section := range expectedSections {
		if !strings.Contains(human, section) {
			t.Errorf("FormatHuman() should contain '%s'", section)
		}
	}

	// Test that version info appears in output
	if !strings.Contains(human, info.Version) {
		t.Error("FormatHuman() should contain version number")
	}

	if !strings.Contains(human, info.GoVersion) {
		t.Error("FormatHuman() should contain Go version")
	}
}

func TestFormatJSON(t *testing.T) {
	info := GetBuildInfo()
	jsonStr, err := info.FormatJSON()

	if err != nil {
		t.Fatalf("FormatJSON() returned error: %v", err)
	}

	if jsonStr == "" {
		t.Error("FormatJSON() returned empty string")
	}

	// Test that it's valid JSON by unmarshaling
	var unmarshaled BuildInfo
	err = json.Unmarshal([]byte(jsonStr), &unmarshaled)
	if err != nil {
		t.Errorf("FormatJSON() produced invalid JSON: %v", err)
	}

	// Test that key fields are preserved
	if unmarshaled.Version != info.Version {
		t.Error("JSON marshaling/unmarshaling changed Version")
	}

	if unmarshaled.GoVersion != info.GoVersion {
		t.Error("JSON marshaling/unmarshaling changed GoVersion")
	}

	if unmarshaled.NumCPU != info.NumCPU {
		t.Error("JSON marshaling/unmarshaling changed NumCPU")
	}
}

func TestFormatCompact(t *testing.T) {
	info := GetBuildInfo()
	compact := info.FormatCompact()

	if compact == "" {
		t.Error("FormatCompact() returned empty string")
	}

	// Test that compact format contains key information
	if !strings.Contains(compact, "SilentCast") {
		t.Error("FormatCompact() should contain 'SilentCast'")
	}

	if !strings.Contains(compact, info.Version) {
		t.Error("FormatCompact() should contain version")
	}

	if !strings.Contains(compact, info.GoOS) {
		t.Error("FormatCompact() should contain OS")
	}

	if !strings.Contains(compact, info.GoArch) {
		t.Error("FormatCompact() should contain architecture")
	}
}

func TestGetVersionString(t *testing.T) {
	version := GetVersionString()

	if version == "" {
		t.Error("GetVersionString() returned empty string")
	}

	// Should match the Version variable
	if version != Version {
		t.Errorf("GetVersionString() = %s, want %s", version, Version)
	}
}

func TestGetShortVersion(t *testing.T) {
	// Save original values
	originalVersion := Version
	originalGitCommit := GitCommit
	defer func() {
		Version = originalVersion
		GitCommit = originalGitCommit
	}()

	// Test with known commit
	Version = "1.0.0"
	GitCommit = "abc123def456"

	short := GetShortVersion()
	expected := "1.0.0-abc123de"

	if short != expected {
		t.Errorf("GetShortVersion() = %s, want %s", short, expected)
	}

	// Test with unknown commit
	GitCommit = "unknown"
	short = GetShortVersion()

	if short != Version {
		t.Errorf("GetShortVersion() with unknown commit = %s, want %s", short, Version)
	}

	// Test with short commit
	GitCommit = "abc123"
	short = GetShortVersion()

	if short != Version {
		t.Errorf("GetShortVersion() with short commit = %s, want %s", short, Version)
	}
}

func TestBuildInfoFields(t *testing.T) {
	info := GetBuildInfo()

	// Test that runtime information matches runtime package
	if info.GoVersion != runtime.Version() {
		t.Errorf("GoVersion mismatch: got %s, want %s", info.GoVersion, runtime.Version())
	}

	if info.GoArch != runtime.GOARCH {
		t.Errorf("GoArch mismatch: got %s, want %s", info.GoArch, runtime.GOARCH)
	}

	if info.GoOS != runtime.GOOS {
		t.Errorf("GoOS mismatch: got %s, want %s", info.GoOS, runtime.GOOS)
	}

	if info.Compiler != runtime.Compiler {
		t.Errorf("Compiler mismatch: got %s, want %s", info.Compiler, runtime.Compiler)
	}

	if info.NumCPU != runtime.NumCPU() {
		t.Errorf("NumCPU mismatch: got %d, want %d", info.NumCPU, runtime.NumCPU())
	}
}

func TestMinFunction(t *testing.T) {
	tests := []struct {
		a, b, want int
	}{
		{1, 2, 1},
		{2, 1, 1},
		{5, 5, 5},
		{0, 10, 0},
		{-1, 5, -1},
		{-5, -3, -5},
	}

	for _, tt := range tests {
		got := min(tt.a, tt.b)
		if got != tt.want {
			t.Errorf("min(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestBuildInfoStructure(t *testing.T) {
	info := BuildInfo{
		Version:      "test-version",
		GitCommit:    "test-commit",
		BuildTime:    "test-time",
		GoVersion:    "test-go",
		GoArch:       "test-arch",
		GoOS:         "test-os",
		CGOEnabled:   "test-cgo",
		BuildTags:    []string{"tag1", "tag2"},
		Compiler:     "test-compiler",
		BuildMode:    "test-mode",
		NumCPU:       4,
		NumGoroutine: 10,
	}

	// Test JSON marshaling
	data, err := json.Marshal(info)
	if err != nil {
		t.Fatalf("Failed to marshal BuildInfo: %v", err)
	}

	var unmarshaled BuildInfo
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal BuildInfo: %v", err)
	}

	// Test all fields are preserved
	if unmarshaled.Version != info.Version {
		t.Error("Version not preserved in JSON")
	}
	if unmarshaled.GitCommit != info.GitCommit {
		t.Error("GitCommit not preserved in JSON")
	}
	if len(unmarshaled.BuildTags) != len(info.BuildTags) {
		t.Error("BuildTags not preserved in JSON")
	}
}

func TestIsCGOEnabled(t *testing.T) {
	// This test just ensures the function doesn't panic
	enabled := isCGOEnabled()

	// The result should be deterministic for a given build
	enabled2 := isCGOEnabled()
	if enabled != enabled2 {
		t.Error("isCGOEnabled() should return consistent results")
	}

	// Should be bool
	if enabled != true && enabled != false {
		t.Error("isCGOEnabled() should return a boolean")
	}
}
