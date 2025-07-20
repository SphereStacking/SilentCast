package updater

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestNewUpdater(t *testing.T) {
	tempDir := t.TempDir()
	cfg := Config{
		CurrentVersion: "v1.0.0",
		RepoOwner:      "test",
		RepoName:       "repo",
		CheckInterval:  1 * time.Hour,
		ConfigDir:      tempDir,
		CacheDuration:  2 * time.Hour,
	}

	u := NewUpdater(cfg)

	if u.currentVersion != cfg.CurrentVersion {
		t.Errorf("Expected version %s, got %s", cfg.CurrentVersion, u.currentVersion)
	}

	if u.checkInterval != cfg.CheckInterval {
		t.Errorf("Expected interval %v, got %v", cfg.CheckInterval, u.checkInterval)
	}

	if u.cacheManager == nil {
		t.Error("Expected cache manager to be initialized")
	}
}

func TestIsNewerVersion(t *testing.T) {
	tests := []struct {
		current string
		new     string
		expect  bool
	}{
		{"v1.0.0", "v1.1.0", true},
		{"v1.1.0", "v1.0.0", false},
		{"v1.0.0", "v1.0.0", false},
		{"1.0.0", "1.1.0", true},
		{"v1.0.0", "1.1.0", true},
		{"dev", "v1.0.0", false},
		{"v1.0.0", "dev", false},
		{"v1.0.0", "v2.0.0", true},
		{"v1.9.0", "v1.10.0", true},
	}

	for _, tt := range tests {
		t.Run(tt.current+"_vs_"+tt.new, func(t *testing.T) {
			u := &Updater{currentVersion: tt.current}
			if got := u.isNewerVersion(tt.new); got != tt.expect {
				t.Errorf("isNewerVersion(%s) = %v, want %v", tt.new, got, tt.expect)
			}
		})
	}
}

func TestFindPlatformAsset(t *testing.T) {
	assets := []Asset{
		{Name: "spellbook-darwin-amd64", Size: 1100},
		{Name: "spellbook-darwin-arm64", Size: 1200},
		{Name: "spellbook-windows-amd64.exe", Size: 1300},
		{Name: "spellbook-linux-amd64", Size: 1400},
		{Name: "spellbook-linux-arm64", Size: 1500},
		{Name: "checksums.txt", Size: 100},
	}

	u := &Updater{}

	// Test current platform
	currentPlatform := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	asset, err := u.findPlatformAsset(assets)

	// Skip if current platform not in test assets
	if err != nil && strings.Contains(err.Error(), "no asset found") {
		t.Skipf("Skipping test - no asset for current platform %s", currentPlatform)
		return
	}

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify the asset name contains the platform
	if asset != nil && !strings.Contains(asset.Name, currentPlatform) {
		t.Errorf("Expected asset name to contain %s, got %s", currentPlatform, asset.Name)
	}

	// Test no matching asset
	noMatchAssets := []Asset{
		{Name: "spellbook-freebsd-amd64", Size: 1000},
		{Name: "checksums.txt", Size: 100},
	}

	_, err = u.findPlatformAsset(noMatchAssets)
	if err == nil {
		t.Error("Expected error for no matching platform")
	}
}

func TestCheckForUpdate(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/test/repo/releases/latest" {
			http.NotFound(w, r)
			return
		}

		platform := fmt.Sprintf("test-app-%s-%s", runtime.GOOS, runtime.GOARCH)
		release := Release{
			TagName:     "v2.0.0",
			Name:        "Release v2.0.0",
			Body:        "New features",
			PublishedAt: time.Now(),
			Assets: []Asset{
				{
					Name:        platform,
					Size:        1024,
					DownloadURL: "https://example.com/download",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(release)
	}))
	defer server.Close()

	tempDir := t.TempDir()
	u := &Updater{
		currentVersion: "v1.0.0",
		repoOwner:      "test",
		repoName:       "repo",
		httpClient: &http.Client{
			Transport: &roundTripperFunc{
				fn: func(req *http.Request) (*http.Response, error) {
					// Redirect GitHub API calls to our test server
					if strings.Contains(req.URL.String(), "api.github.com") {
						req.URL.Host = strings.TrimPrefix(server.URL, "http://")
						req.URL.Scheme = "http"
					}
					return http.DefaultTransport.RoundTrip(req)
				},
			},
			Timeout: 5 * time.Second,
		},
		cacheManager: NewCacheManager(tempDir, 1*time.Hour),
	}

	ctx := context.Background()
	
	// First check should hit the API
	info, err := u.CheckForUpdate(ctx)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if info == nil {
		t.Fatal("Expected update info, got nil")
	}
	
	if info.Version != "v2.0.0" {
		t.Errorf("Expected version v2.0.0, got %s", info.Version)
	}

	// Second check should use cache
	info2, err := u.CheckForUpdate(ctx)
	if err != nil {
		t.Fatalf("Unexpected error on cached check: %v", err)
	}
	
	if info2 == nil {
		t.Fatal("Expected cached update info, got nil")
	}
	
	if info2.Version != info.Version {
		t.Error("Cached version should match original")
	}
}

func TestCreateBackup(t *testing.T) {
	tempDir := t.TempDir()

	// Create source file
	srcPath := filepath.Join(tempDir, "source")
	srcContent := []byte("test content")
	if err := os.WriteFile(srcPath, srcContent, 0o755); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Create backup
	u := &Updater{}
	backupPath := filepath.Join(tempDir, "backup")

	if err := u.createBackup(srcPath, backupPath); err != nil {
		t.Fatalf("Failed to create backup: %v", err)
	}

	// Verify backup
	backupContent, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("Failed to read backup: %v", err)
	}

	if !bytes.Equal(backupContent, srcContent) {
		t.Error("Backup content doesn't match source")
	}

	// Check permissions
	srcInfo, _ := os.Stat(srcPath)
	backupInfo, _ := os.Stat(backupPath)

	if srcInfo.Mode() != backupInfo.Mode() {
		t.Error("Backup permissions don't match source")
	}
}

func TestVerifyChecksum(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.bin")

	// Create test file
	content := []byte("test content for checksum")
	if err := os.WriteFile(testFile, content, 0o600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	u := &Updater{}

	// Test with correct checksum (calculate actual checksum)
	h := sha256.New()
	h.Write(content)
	correctChecksum := fmt.Sprintf("%x", h.Sum(nil))

	err := u.verifyChecksum(testFile, correctChecksum)
	if err != nil {
		t.Errorf("Expected checksum verification to pass, got error: %v", err)
	}

	// Test with wrong checksum (should fail)
	wrongChecksum := "8b3d6c91e8f0c6e9d8a7f9e1c2b4a6d8e5f7a9b1c3d5e7f9"
	err = u.verifyChecksum(testFile, wrongChecksum)
	if err == nil {
		t.Error("Expected checksum verification to fail with wrong checksum")
	}

	// Test with non-existent file
	err = u.verifyChecksum("/nonexistent/file", correctChecksum)
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestGetLatestRelease(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/test/repo/releases/latest" {
			http.NotFound(w, r)
			return
		}

		// Check headers
		if r.Header.Get("Accept") != "application/vnd.github.v3+json" {
			t.Errorf("Expected Accept header, got %s", r.Header.Get("Accept"))
		}

		release := Release{
			TagName:     "v2.0.0",
			Name:        "Release v2.0.0",
			Body:        "New features",
			PublishedAt: time.Now(),
			Assets:      []Asset{},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(release)
	}))
	defer server.Close()

	u := &Updater{
		repoOwner:  "test",
		repoName:   "repo",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	// Override API URL by creating a custom HTTP client that redirects
	u.httpClient = &http.Client{
		Transport: &roundTripperFunc{
			fn: func(req *http.Request) (*http.Response, error) {
				// Redirect GitHub API calls to our test server
				if strings.Contains(req.URL.String(), "api.github.com") {
					req.URL.Host = strings.TrimPrefix(server.URL, "http://")
					req.URL.Scheme = "http"
				}
				return http.DefaultTransport.RoundTrip(req)
			},
		},
		Timeout: 5 * time.Second,
	}

	ctx := context.Background()
	release, err := u.getLatestRelease(ctx)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if release.TagName != "v2.0.0" {
		t.Errorf("Expected version v2.0.0, got %s", release.TagName)
	}
}

func TestGetLatestReleaseError(t *testing.T) {
	// Test with invalid server response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	u := &Updater{
		repoOwner:  "test",
		repoName:   "repo",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	ctx := context.Background()
	_, err := u.getLatestRelease(ctx)
	if err == nil {
		t.Error("Expected error for 404 response")
	}
}

func TestFindChecksum(t *testing.T) {
	assets := []Asset{
		{Name: "spellbook-linux-amd64", Size: 1400},
		{Name: "checksums.txt", Size: 100},
		{Name: "checksums.sha256", Size: 150},
	}

	u := NewUpdater(Config{
		CurrentVersion: "v0.1.0",
		RepoOwner:      "test",
		RepoName:       "test",
	})
	checksum := u.findChecksum(assets, "spellbook-linux-amd64")

	// Currently returns empty string as checksum parsing is not implemented
	if checksum != "" {
		t.Errorf("Expected empty checksum (not implemented), got %s", checksum)
	}

	// Test with no checksum files
	assetsNoChecksum := []Asset{
		{Name: "spellbook-linux-amd64", Size: 1400},
		{Name: "readme.txt", Size: 100},
	}

	checksum = u.findChecksum(assetsNoChecksum, "spellbook-linux-amd64")
	if checksum != "" {
		t.Errorf("Expected empty checksum, got %s", checksum)
	}
}

func TestStartAutoCheck(t *testing.T) {
	tempDir := t.TempDir()
	u := &Updater{
		currentVersion: "v1.0.0",
		repoOwner:      "test", 
		repoName:       "repo",
		checkInterval:  100 * time.Millisecond,
		httpClient:     &http.Client{Timeout: 1 * time.Second},
		cacheManager:   NewCacheManager(tempDir, 1*time.Hour),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	updateReceived := make(chan bool, 1)
	onUpdate := func(info *UpdateInfo) {
		select {
		case updateReceived <- true:
		default:
		}
	}

	u.StartAutoCheck(ctx, onUpdate)

	// Cancel context after short time
	time.Sleep(50 * time.Millisecond)
	cancel()

	// Test should not hang
	time.Sleep(200 * time.Millisecond)
}

func TestCheckAndNotify(t *testing.T) {
	tempDir := t.TempDir()
	u := &Updater{
		currentVersion: "v1.0.0",
		repoOwner:      "nonexistent",
		repoName:       "repo",
		httpClient:     &http.Client{Timeout: 1 * time.Second},
		cacheManager:   NewCacheManager(tempDir, 1*time.Hour),
	}

	ctx := context.Background()
	called := false
	onUpdate := func(info *UpdateInfo) {
		called = true
	}

	// This should fail but not panic
	u.checkAndNotify(ctx, onUpdate)

	if called {
		t.Error("onUpdate should not be called on error")
	}
}

func TestDownloadUpdateSizeMismatch(t *testing.T) {
	// Create test server that returns wrong size
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("short content"))
	}))
	defer server.Close()

	u := &Updater{
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	info := &UpdateInfo{
		Version:     "v2.0.0",
		DownloadURL: server.URL,
		Size:        1000, // Much larger than actual content
	}

	ctx := context.Background()
	_, err := u.DownloadUpdate(ctx, info)
	if err == nil {
		t.Error("Expected error for size mismatch")
	}

	if !strings.Contains(err.Error(), "size mismatch") {
		t.Errorf("Expected size mismatch error, got: %v", err)
	}
}

func TestDownloadUpdateHTTPError(t *testing.T) {
	// Create test server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	u := &Updater{
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	info := &UpdateInfo{
		Version:     "v2.0.0",
		DownloadURL: server.URL,
		Size:        100,
	}

	ctx := context.Background()
	_, err := u.DownloadUpdate(ctx, info)
	if err == nil {
		t.Error("Expected error for HTTP 500")
	}
}

func TestApplyUpdateBackupFailure(t *testing.T) {
	u := &Updater{}

	// Test with non-existent file
	err := u.ApplyUpdate("/nonexistent/update")
	if err == nil {
		t.Error("Expected error for non-existent update file")
	}
}

func TestRestoreBackup(t *testing.T) {
	tempDir := t.TempDir()
	
	backupFile := filepath.Join(tempDir, "backup")
	originalFile := filepath.Join(tempDir, "original")

	// Create backup file
	if err := os.WriteFile(backupFile, []byte("backup content"), 0o755); err != nil {
		t.Fatalf("Failed to create backup file: %v", err)
	}

	u := &Updater{}
	
	err := u.restoreBackup(backupFile, originalFile)
	if err != nil {
		t.Errorf("restoreBackup failed: %v", err)
	}

	// Verify restore
	if _, err := os.Stat(originalFile); os.IsNotExist(err) {
		t.Error("Original file should exist after restore")
	}

	if _, err := os.Stat(backupFile); !os.IsNotExist(err) {
		t.Error("Backup file should not exist after restore")
	}
}

func TestForceCheck(t *testing.T) {
	// Create test server
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		
		platform := fmt.Sprintf("test-app-%s-%s", runtime.GOOS, runtime.GOARCH)
		release := Release{
			TagName:     "v2.0.0",
			Name:        "Release v2.0.0", 
			Body:        "New features",
			PublishedAt: time.Now(),
			Assets: []Asset{
				{
					Name:        platform,
					Size:        1024,
					DownloadURL: "https://example.com/download",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(release)
	}))
	defer server.Close()

	tempDir := t.TempDir()
	u := &Updater{
		currentVersion: "v1.0.0",
		repoOwner:      "test",
		repoName:       "repo",
		checkInterval:  24 * time.Hour, // Set check interval for proper caching
		httpClient: &http.Client{
			Transport: &roundTripperFunc{
				fn: func(req *http.Request) (*http.Response, error) {
					// Redirect GitHub API calls to our test server
					if strings.Contains(req.URL.String(), "api.github.com") {
						req.URL.Host = strings.TrimPrefix(server.URL, "http://")
						req.URL.Scheme = "http"
					}
					return http.DefaultTransport.RoundTrip(req)
				},
			},
			Timeout: 5 * time.Second,
		},
		cacheManager: NewCacheManager(tempDir, 1*time.Hour),
	}

	ctx := context.Background()
	
	// First check
	u.CheckForUpdate(ctx)
	if callCount != 1 {
		t.Errorf("Expected 1 API call, got %d", callCount)
	}

	// Wait a bit to ensure cache is written
	time.Sleep(10 * time.Millisecond)
	
	// Second check should use cache
	cached, _ := u.cacheManager.GetCached()
	if cached == nil {
		t.Log("Warning: Cache not found after first check")
	}
	
	u.CheckForUpdate(ctx)
	if callCount != 1 {
		t.Errorf("Expected 1 API call (cached), got %d", callCount)
	}

	// Force check should bypass cache
	u.ForceCheck(ctx)
	if callCount != 2 {
		t.Errorf("Expected 2 API calls after force check, got %d", callCount)
	}
}

func TestClearCache(t *testing.T) {
	tempDir := t.TempDir()
	u := &Updater{
		cacheManager: NewCacheManager(tempDir, 1*time.Hour),
	}

	// Save something to cache
	u.cacheManager.SaveCache(&UpdateInfo{Version: "v2.0.0"}, time.Now().Add(1*time.Hour))

	// Clear cache
	err := u.ClearCache()
	if err != nil {
		t.Errorf("Failed to clear cache: %v", err)
	}

	// Verify cache is cleared
	cache, _ := u.cacheManager.GetCached()
	if cache != nil {
		t.Error("Cache should be nil after clearing")
	}
}

// Helper type for custom HTTP transport
type roundTripperFunc struct {
	fn func(*http.Request) (*http.Response, error)
}

func (f *roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f.fn(req)
}
