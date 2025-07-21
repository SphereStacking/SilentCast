package updater

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCacheManager_SaveAndGet(t *testing.T) {
	tempDir := t.TempDir()
	cm := NewCacheManager(tempDir, 1*time.Hour)

	// Test saving cache
	updateInfo := &UpdateInfo{
		Version:      "v2.0.0",
		ReleaseNotes: "Test release",
		PublishedAt:  time.Now(),
		DownloadURL:  "https://example.com/download",
		Size:         1024,
	}

	nextCheck := time.Now().Add(24 * time.Hour)
	err := cm.SaveCache(updateInfo, nextCheck)
	if err != nil {
		t.Fatalf("Failed to save cache: %v", err)
	}

	// Test getting cache
	cache, err := cm.GetCached()
	if err != nil {
		t.Fatalf("Failed to get cache: %v", err)
	}

	if cache == nil {
		t.Fatal("Expected cache, got nil")
	}

	if cache.UpdateInfo == nil {
		t.Fatal("Expected update info in cache")
	}

	if cache.UpdateInfo.Version != updateInfo.Version {
		t.Errorf("Expected version %s, got %s", updateInfo.Version, cache.UpdateInfo.Version)
	}

	// Verify cache file exists
	cacheFile := filepath.Join(tempDir, "cache", "update_check.json")
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		t.Error("Cache file should exist")
	}
}

func TestCacheManager_Expiry(t *testing.T) {
	tempDir := t.TempDir()
	// Create cache with very short duration
	cm := NewCacheManager(tempDir, 100*time.Millisecond)

	// Save cache
	updateInfo := &UpdateInfo{
		Version: "v2.0.0",
	}

	nextCheck := time.Now().Add(1 * time.Hour)
	err := cm.SaveCache(updateInfo, nextCheck)
	if err != nil {
		t.Fatalf("Failed to save cache: %v", err)
	}

	// Cache should be valid immediately
	cache, err := cm.GetCached()
	if err != nil || cache == nil {
		t.Error("Cache should be valid immediately after saving")
	}

	// Wait for cache to expire
	time.Sleep(150 * time.Millisecond)

	// Cache should be expired
	cache, err = cm.GetCached()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if cache != nil {
		t.Error("Cache should be expired")
	}
}

func TestCacheManager_ShouldCheck(t *testing.T) {
	tempDir := t.TempDir()
	cm := NewCacheManager(tempDir, 1*time.Hour)

	// Should check when no cache exists
	if !cm.ShouldCheck() {
		t.Error("Should check when no cache exists")
	}

	// Save cache with future next check
	updateInfo := &UpdateInfo{Version: "v2.0.0"}
	nextCheck := time.Now().Add(1 * time.Hour)
	err := cm.SaveCache(updateInfo, nextCheck)
	if err != nil {
		t.Fatalf("Failed to save cache: %v", err)
	}

	// Should not check when next check is in the future
	if cm.ShouldCheck() {
		t.Error("Should not check when next check is in the future")
	}

	// Save cache with past next check
	nextCheck = time.Now().Add(-1 * time.Hour)
	err = cm.SaveCache(nil, nextCheck)
	if err != nil {
		t.Fatalf("Failed to save cache: %v", err)
	}

	// Should check when next check is in the past
	if !cm.ShouldCheck() {
		t.Error("Should check when next check is in the past")
	}
}

func TestCacheManager_ClearCache(t *testing.T) {
	tempDir := t.TempDir()
	cm := NewCacheManager(tempDir, 1*time.Hour)

	// Save cache
	updateInfo := &UpdateInfo{Version: "v2.0.0"}
	nextCheck := time.Now().Add(1 * time.Hour)
	err := cm.SaveCache(updateInfo, nextCheck)
	if err != nil {
		t.Fatalf("Failed to save cache: %v", err)
	}

	// Verify cache exists
	cache, _ := cm.GetCached()
	if cache == nil {
		t.Fatal("Cache should exist before clearing")
	}

	// Clear cache
	err = cm.ClearCache()
	if err != nil {
		t.Fatalf("Failed to clear cache: %v", err)
	}

	// Verify cache is cleared
	cache, _ = cm.GetCached()
	if cache != nil {
		t.Error("Cache should be nil after clearing")
	}

	// Clear non-existent cache should not error
	err = cm.ClearCache()
	if err != nil {
		t.Error("Clearing non-existent cache should not error")
	}
}

func TestCacheManager_InvalidJSON(t *testing.T) {
	tempDir := t.TempDir()
	cm := NewCacheManager(tempDir, 1*time.Hour)

	// Create cache directory
	cacheDir := filepath.Join(tempDir, "cache")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		t.Fatalf("Failed to create cache dir: %v", err)
	}

	// Write invalid JSON
	cacheFile := filepath.Join(cacheDir, "update_check.json")
	err := os.WriteFile(cacheFile, []byte("invalid json"), 0o644)
	if err != nil {
		t.Fatalf("Failed to write invalid cache: %v", err)
	}

	// Getting cache should fail
	cache, err := cm.GetCached()
	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
	if cache != nil {
		t.Error("Cache should be nil for invalid JSON")
	}
}

func TestCacheManager_NoUpdateInfo(t *testing.T) {
	tempDir := t.TempDir()
	cm := NewCacheManager(tempDir, 1*time.Hour)

	// Save cache with nil update info (no update available)
	nextCheck := time.Now().Add(1 * time.Hour)
	err := cm.SaveCache(nil, nextCheck)
	if err != nil {
		t.Fatalf("Failed to save cache: %v", err)
	}

	// Get cache
	cache, err := cm.GetCached()
	if err != nil {
		t.Fatalf("Failed to get cache: %v", err)
	}

	if cache == nil {
		t.Fatal("Expected cache, got nil")
	}

	if cache.UpdateInfo != nil {
		t.Error("UpdateInfo should be nil when no update is available")
	}
}
