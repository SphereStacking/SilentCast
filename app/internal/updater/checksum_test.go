package updater

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestGetChecksumForAsset(t *testing.T) {
	tests := []struct {
		name          string
		assets        []Asset
		assetName     string
		checksumFile  string
		statusCode    int
		expectedHash  string
		expectedError bool
	}{
		{
			name: "Valid checksum found",
			assets: []Asset{
				{Name: "silentcast-linux-amd64", DownloadURL: "http://example.com/binary"},
				{Name: "checksums.txt", DownloadURL: "http://example.com/checksums"},
			},
			assetName: "silentcast-linux-amd64",
			checksumFile: `abc123def456  silentcast-linux-amd64
789xyz012345  silentcast-darwin-amd64`,
			statusCode:   http.StatusOK,
			expectedHash: "abc123def456",
		},
		{
			name: "Checksum with ./ prefix",
			assets: []Asset{
				{Name: "app.zip", DownloadURL: "http://example.com/app"},
				{Name: "SHA256SUMS", DownloadURL: "http://example.com/sums"},
			},
			assetName: "app.zip",
			checksumFile: `123456  ./app.zip
789012  ./other.zip`,
			statusCode:   http.StatusOK,
			expectedHash: "123456",
		},
		{
			name: "Checksum with * prefix",
			assets: []Asset{
				{Name: "release.tar.gz", DownloadURL: "http://example.com/release"},
				{Name: "checksums.sha256", DownloadURL: "http://example.com/sha256"},
			},
			assetName:    "release.tar.gz",
			checksumFile: `abcdef  *release.tar.gz`,
			statusCode:   http.StatusOK,
			expectedHash: "abcdef",
		},
		{
			name: "No checksum file",
			assets: []Asset{
				{Name: "binary", DownloadURL: "http://example.com/binary"},
			},
			assetName:    "binary",
			expectedHash: "",
		},
		{
			name: "Checksum not found for asset",
			assets: []Asset{
				{Name: "app", DownloadURL: "http://example.com/app"},
				{Name: "checksums.txt", DownloadURL: "http://example.com/checksums"},
			},
			assetName:    "app",
			checksumFile: `123456  other-app`,
			statusCode:   http.StatusOK,
			expectedHash: "",
		},
		{
			name: "HTTP error",
			assets: []Asset{
				{Name: "app", DownloadURL: "http://example.com/app"},
				{Name: "checksums.txt", DownloadURL: "http://example.com/checksums"},
			},
			assetName:     "app",
			statusCode:    http.StatusNotFound,
			expectedError: true,
		},
		{
			name: "Invalid checksum format",
			assets: []Asset{
				{Name: "app", DownloadURL: "http://example.com/app"},
				{Name: "checksums.txt", DownloadURL: "http://example.com/checksums"},
			},
			assetName:    "app",
			checksumFile: `invalid-format-no-spaces`,
			statusCode:   http.StatusOK,
			expectedHash: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.statusCode != 0 {
					w.WriteHeader(tt.statusCode)
				}
				if tt.checksumFile != "" {
					fmt.Fprint(w, tt.checksumFile)
				}
			}))
			defer server.Close()

			// Update asset URLs to use test server
			for i := range tt.assets {
				if strings.Contains(strings.ToLower(tt.assets[i].Name), "checksum") || 
				   strings.Contains(strings.ToLower(tt.assets[i].Name), "sha256") {
					tt.assets[i].DownloadURL = server.URL
				}
			}

			verifier := NewChecksumVerifier(http.DefaultClient)
			hash, err := verifier.GetChecksumForAsset(context.Background(), tt.assets, tt.assetName)

			if tt.expectedError {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if hash != tt.expectedHash {
					t.Errorf("Expected hash %q, got %q", tt.expectedHash, hash)
				}
			}
		})
	}
}

func TestGetChecksumForAsset_ContextCancellation(t *testing.T) {
	// Test context cancellation
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Delay to allow context cancellation
		select {
		case <-r.Context().Done():
			return
		case <-time.After(100 * time.Millisecond):
			fmt.Fprint(w, "should not reach here")
		}
	}))
	defer server.Close()

	assets := []Asset{
		{Name: "app", DownloadURL: "http://example.com/app"},
		{Name: "checksums.txt", DownloadURL: server.URL},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	verifier := NewChecksumVerifier(http.DefaultClient)
	_, err := verifier.GetChecksumForAsset(ctx, assets, "app")
	if err == nil {
		t.Error("Expected error due to context cancellation")
	}
}