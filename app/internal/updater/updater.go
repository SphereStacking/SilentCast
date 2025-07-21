package updater

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/SphereStacking/silentcast/internal/config"
	"github.com/SphereStacking/silentcast/pkg/logger"
)

// Updater manages application updates
type Updater struct {
	currentVersion string
	repoOwner      string
	repoName       string
	checkInterval  time.Duration
	httpClient     *http.Client
	cacheManager   *CacheManager
}

// Config holds updater configuration
type Config struct {
	CurrentVersion string
	RepoOwner      string
	RepoName       string
	CheckInterval  time.Duration
	AutoUpdate     bool
	ConfigDir      string
	CacheDuration  time.Duration
}

// Release represents a GitHub release
type Release struct {
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	Body        string    `json:"body"`
	Prerelease  bool      `json:"prerelease"`
	Draft       bool      `json:"draft"`
	PublishedAt time.Time `json:"published_at"`
	Assets      []Asset   `json:"assets"`
}

// Asset represents a release asset
type Asset struct {
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	DownloadURL string `json:"browser_download_url"`
	ContentType string `json:"content_type"`
}

// UpdateInfo contains information about an available update
type UpdateInfo struct {
	Version      string
	ReleaseNotes string
	PublishedAt  time.Time
	DownloadURL  string
	Size         int64
	Checksum     string
}

// NewUpdater creates a new updater instance
func NewUpdater(cfg Config) *Updater {
	if cfg.CheckInterval == 0 {
		cfg.CheckInterval = 24 * time.Hour
	}

	if cfg.CacheDuration == 0 {
		cfg.CacheDuration = 1 * time.Hour
	}

	return &Updater{
		currentVersion: cfg.CurrentVersion,
		repoOwner:      cfg.RepoOwner,
		repoName:       cfg.RepoName,
		checkInterval:  cfg.CheckInterval,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		cacheManager: NewCacheManager(cfg.ConfigDir, cfg.CacheDuration),
	}
}

// CheckForUpdate checks if a new version is available
func (u *Updater) CheckForUpdate(ctx context.Context) (*UpdateInfo, error) {
	// Check cache first
	if !u.cacheManager.ShouldCheck() {
		cache, err := u.cacheManager.GetCached()
		if err == nil && cache != nil {
			logger.Info("Using cached update check result")
			return cache.UpdateInfo, nil
		}
	}

	logger.Info("Checking for updates...")

	// Get latest release from GitHub
	release, err := u.getLatestRelease(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest release: %w", err)
	}

	// Compare versions
	var updateInfo *UpdateInfo
	if !u.isNewerVersion(release.TagName) {
		logger.Info("Already running latest version %s", u.currentVersion)
		updateInfo = nil
	} else {
		// Find appropriate asset for current platform
		asset, err := u.findPlatformAsset(release.Assets)
		if err != nil {
			return nil, fmt.Errorf("no suitable update found for platform: %w", err)
		}

		// Get checksum if available
		checksum := u.findChecksum(release.Assets, asset.Name)

		updateInfo = &UpdateInfo{
			Version:      release.TagName,
			ReleaseNotes: release.Body,
			PublishedAt:  release.PublishedAt,
			DownloadURL:  asset.DownloadURL,
			Size:         asset.Size,
			Checksum:     checksum,
		}
	}

	// Cache the result
	nextCheck := time.Now().Add(u.checkInterval)
	if err := u.cacheManager.SaveCache(updateInfo, nextCheck); err != nil {
		logger.Error("Failed to cache update check result: %v", err)
	}

	return updateInfo, nil
}

// DownloadUpdate downloads the update to a temporary file
func (u *Updater) DownloadUpdate(ctx context.Context, info *UpdateInfo) (string, error) {
	return u.DownloadUpdateWithProgress(ctx, info, nil)
}

// DownloadUpdateWithProgress downloads the update with optional progress reporting
func (u *Updater) DownloadUpdateWithProgress(ctx context.Context, info *UpdateInfo, progressWriter io.Writer) (string, error) {
	logger.Info("Downloading update %s...", info.Version)

	// Create temporary file
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, fmt.Sprintf("%s-update-%s", config.AppName, info.Version))

	// Download file
	out, err := os.Create(tmpFile)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer out.Close()

	req, err := http.NewRequestWithContext(ctx, "GET", info.DownloadURL, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status: %s", resp.Status)
	}

	// Copy with progress tracking
	var written int64
	if progressWriter != nil {
		// Use progress writer
		pw := NewProgressWriter(out, info.Size, progressWriter)
		written, err = io.Copy(pw, resp.Body)
		pw.Done()
	} else {
		// No progress tracking
		written, err = io.Copy(out, resp.Body)
	}
	
	if err != nil {
		return "", fmt.Errorf("failed to save update: %w", err)
	}

	if written != info.Size {
		return "", fmt.Errorf("download size mismatch: expected %d, got %d", info.Size, written)
	}

	// Verify checksum if available
	if info.Checksum != "" {
		if progressWriter != nil {
			fmt.Fprintln(progressWriter, "🔐 Verifying checksum...")
		}
		if err := u.verifyChecksum(tmpFile, info.Checksum); err != nil {
			os.Remove(tmpFile)
			return "", fmt.Errorf("checksum verification failed: %w", err)
		}
	}

	// Make executable using platform-specific method
	platform := GetPlatformUpdater()
	if err := platform.MakeExecutable(tmpFile); err != nil {
		return "", fmt.Errorf("failed to make executable: %w", err)
	}

	logger.Info("Update downloaded successfully to %s", tmpFile)
	return tmpFile, nil
}

// ApplyUpdate applies the downloaded update
func (u *Updater) ApplyUpdate(updatePath string) error {
	logger.Info("Applying update...")

	// Get current executable path
	currentPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get current executable: %w", err)
	}

	// Resolve any symlinks
	currentPath, err = filepath.EvalSymlinks(currentPath)
	if err != nil {
		return fmt.Errorf("failed to resolve executable path: %w", err)
	}

	// Create backup of current version
	backupPath := currentPath + ".backup"
	if err := u.createBackup(currentPath, backupPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Replace current executable
	if err := u.replaceExecutable(updatePath, currentPath); err != nil {
		// Restore backup on failure
		if restoreErr := u.restoreBackup(backupPath, currentPath); restoreErr != nil {
			logger.Error("Failed to restore backup: %v", restoreErr)
		}
		return fmt.Errorf("failed to apply update: %w", err)
	}

	// Clean up
	os.Remove(backupPath)
	os.Remove(updatePath)

	logger.Info("Update applied successfully!")
	return nil
}

// StartAutoCheck starts periodic update checks
func (u *Updater) StartAutoCheck(ctx context.Context, onUpdate func(*UpdateInfo)) {
	go func() {
		ticker := time.NewTicker(u.checkInterval)
		defer ticker.Stop()

		// Initial check after 1 minute
		time.Sleep(1 * time.Minute)
		u.checkAndNotify(ctx, onUpdate)

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				u.checkAndNotify(ctx, onUpdate)
			}
		}
	}()
}

// checkAndNotify checks for updates and notifies if available
func (u *Updater) checkAndNotify(ctx context.Context, onUpdate func(*UpdateInfo)) {
	info, err := u.CheckForUpdate(ctx)
	if err != nil {
		logger.Error("Update check failed: %v", err)
		return
	}

	if info != nil && onUpdate != nil {
		onUpdate(info)
	}
}

// getLatestRelease fetches the latest release from GitHub
func (u *Updater) getLatestRelease(ctx context.Context) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest",
		u.repoOwner, u.repoName)

	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return nil, err
	}

	// GitHub recommends including Accept header
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := u.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status: %s", resp.Status)
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

// isNewerVersion compares version strings
func (u *Updater) isNewerVersion(newVersion string) bool {
	// Handle dev versions
	if u.currentVersion == "dev" || newVersion == "dev" {
		return false
	}

	// Remove 'v' prefix if present
	current := strings.TrimPrefix(u.currentVersion, "v")
	newer := strings.TrimPrefix(newVersion, "v")

	// Split versions into parts
	currentParts := strings.Split(current, ".")
	newParts := strings.Split(newer, ".")

	// Compare each part numerically
	for i := 0; i < len(currentParts) && i < len(newParts); i++ {
		currentNum, err1 := strconv.Atoi(currentParts[i])
		newNum, err2 := strconv.Atoi(newParts[i])

		// If parsing fails, fall back to string comparison
		if err1 != nil || err2 != nil {
			if newParts[i] > currentParts[i] {
				return true
			} else if newParts[i] < currentParts[i] {
				return false
			}
			continue
		}

		// Numeric comparison
		if newNum > currentNum {
			return true
		} else if newNum < currentNum {
			return false
		}
	}

	// If all compared parts are equal, the longer version is newer
	return len(newParts) > len(currentParts)
}

// findPlatformAsset finds the appropriate asset for current platform
func (u *Updater) findPlatformAsset(assets []Asset) (*Asset, error) {
	platform := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)

	for _, asset := range assets {
		name := strings.ToLower(asset.Name)

		// Check for exact platform match
		if strings.Contains(name, platform) {
			// Skip archives, we want the binary
			if !strings.HasSuffix(name, ".tar.gz") &&
				!strings.HasSuffix(name, ".zip") {
				return &asset, nil
			}
		}
	}

	return nil, fmt.Errorf("no asset found for platform %s", platform)
}

// findChecksum finds checksum for the given asset
func (u *Updater) findChecksum(assets []Asset, assetName string) string {
	// Create checksum verifier
	verifier := NewChecksumVerifier(u.httpClient)
	
	// Try to get checksum, but don't fail if we can't
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	checksum, err := verifier.GetChecksumForAsset(ctx, assets, assetName)
	if err != nil {
		logger.Warn("Failed to fetch checksum: %v", err)
		return ""
	}
	
	return checksum
}

// verifyChecksum verifies file checksum
func (u *Updater) verifyChecksum(filePath, expectedChecksum string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	h := sha256.New()
	if _, err := io.Copy(h, file); err != nil {
		return err
	}

	actualChecksum := fmt.Sprintf("%x", h.Sum(nil))
	if actualChecksum != expectedChecksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s",
			expectedChecksum, actualChecksum)
	}

	return nil
}

// createBackup creates a backup of current executable
func (u *Updater) createBackup(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, copyErr := io.Copy(dstFile, srcFile); copyErr != nil {
		return copyErr
	}

	// Copy permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, srcInfo.Mode())
}

// replaceExecutable replaces the current executable with update
func (u *Updater) replaceExecutable(src, dst string) error {
	platform := GetPlatformUpdater()
	return platform.ReplaceExecutable(src, dst)
}

// restoreBackup restores backup in case of failure
func (u *Updater) restoreBackup(backup, original string) error {
	return os.Rename(backup, original)
}

// ClearCache clears the update check cache
func (u *Updater) ClearCache() error {
	return u.cacheManager.ClearCache()
}

// ForceCheck performs an update check ignoring cache
func (u *Updater) ForceCheck(ctx context.Context) (*UpdateInfo, error) {
	// Clear cache first to force fresh check
	if err := u.cacheManager.ClearCache(); err != nil {
		// Log error but continue with check anyway
		// Cache clearing failure shouldn't prevent update check
		// The check will simply use any remaining cached data
	}
	return u.CheckForUpdate(ctx)
}
