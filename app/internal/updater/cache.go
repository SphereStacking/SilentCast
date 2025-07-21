package updater

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// UpdateCache caches update check results to avoid excessive API calls
type UpdateCache struct {
	LastCheck   time.Time   `json:"last_check"`
	UpdateInfo  *UpdateInfo `json:"update_info,omitempty"`
	NextCheck   time.Time   `json:"next_check"`
	CacheExpiry time.Time   `json:"cache_expiry"`
}

// CacheManager manages update check result caching
type CacheManager struct {
	cacheDir      string
	cacheFile     string
	cacheDuration time.Duration
}

// NewCacheManager creates a new cache manager
func NewCacheManager(configDir string, cacheDuration time.Duration) *CacheManager {
	if cacheDuration == 0 {
		cacheDuration = 1 * time.Hour // Default cache duration
	}

	cacheDir := filepath.Join(configDir, "cache")
	return &CacheManager{
		cacheDir:      cacheDir,
		cacheFile:     filepath.Join(cacheDir, "update_check.json"),
		cacheDuration: cacheDuration,
	}
}

// GetCached retrieves cached update check result
func (cm *CacheManager) GetCached() (*UpdateCache, error) {
	// Ensure cache directory exists
	if err := os.MkdirAll(cm.cacheDir, 0755); err != nil {
		return nil, err
	}

	// Read cache file
	data, err := os.ReadFile(cm.cacheFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No cache exists
		}
		return nil, err
	}

	var cache UpdateCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}

	// Check if cache is still valid
	if time.Now().After(cache.CacheExpiry) {
		return nil, nil // Cache expired
	}

	return &cache, nil
}

// SaveCache saves update check result to cache
func (cm *CacheManager) SaveCache(info *UpdateInfo, nextCheck time.Time) error {
	// Ensure cache directory exists
	if err := os.MkdirAll(cm.cacheDir, 0755); err != nil {
		return err
	}

	cache := UpdateCache{
		LastCheck:   time.Now(),
		UpdateInfo:  info,
		NextCheck:   nextCheck,
		CacheExpiry: time.Now().Add(cm.cacheDuration),
	}

	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cm.cacheFile, data, 0644)
}

// ClearCache removes cached update check result
func (cm *CacheManager) ClearCache() error {
	err := os.Remove(cm.cacheFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// ShouldCheck determines if an update check should be performed
func (cm *CacheManager) ShouldCheck() bool {
	cache, err := cm.GetCached()
	if err != nil || cache == nil {
		return true // Check if no cache or error
	}

	return time.Now().After(cache.NextCheck)
}
