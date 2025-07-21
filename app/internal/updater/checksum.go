package updater

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"strings"
)

// ChecksumVerifier handles checksum verification for updates
type ChecksumVerifier struct {
	httpClient *http.Client
}

// NewChecksumVerifier creates a new checksum verifier
func NewChecksumVerifier(client *http.Client) *ChecksumVerifier {
	return &ChecksumVerifier{
		httpClient: client,
	}
}

// GetChecksumForAsset downloads and parses checksums to find the checksum for a specific asset
func (cv *ChecksumVerifier) GetChecksumForAsset(ctx context.Context, assets []Asset, assetName string) (string, error) {
	// Look for checksum files
	var checksumAsset *Asset
	for _, asset := range assets {
		name := strings.ToLower(asset.Name)
		if strings.Contains(name, "checksum") || strings.Contains(name, "sha256") {
			checksumAsset = &asset
			break
		}
	}

	if checksumAsset == nil {
		return "", nil // No checksum file found
	}

	// Download checksum file
	req, err := http.NewRequestWithContext(ctx, "GET", checksumAsset.DownloadURL, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := cv.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to download checksums: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download checksums: status %s", resp.Status)
	}

	// Parse checksum file
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		// Format is usually: <checksum>  <filename>
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			checksum := parts[0]
			filename := parts[1]

			// Remove leading ./ or * from filename if present
			filename = strings.TrimPrefix(filename, "./")
			filename = strings.TrimPrefix(filename, "*")

			if filename == assetName {
				return checksum, nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to parse checksums: %w", err)
	}

	return "", nil // Checksum not found for asset
}
