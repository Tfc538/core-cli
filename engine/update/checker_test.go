package update

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChecker_Check(t *testing.T) {
	// Mock GitHub API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		release := GitHubRelease{
			TagName: "v1.2.0",
			Body:    "## Changes\n- New feature\n- Bug fixes",
			Assets: []GitHubAsset{
				{
					Name:        "core-linux-amd64",
					DownloadURL: "https://github.com/user/repo/releases/download/v1.2.0/core-linux-amd64",
				},
				{
					Name:        "core-darwin-amd64",
					DownloadURL: "https://github.com/user/repo/releases/download/v1.2.0/core-darwin-amd64",
				},
				{
					Name:        "checksums.txt",
					DownloadURL: "https://github.com/user/repo/releases/download/v1.2.0/checksums.txt",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(release)
	}))
	defer server.Close()

	// Override the API URL for testing (we'll use a custom checker instead)
	checker := &Checker{
		config: CheckerConfig{
			GitHubOwner:    "test-owner",
			GitHubRepo:     "test-repo",
			CurrentVersion: "1.0.0",
		},
	}

	// Test version comparison
	updateAvailable, compatible := checker.compareVersions("1.0.0", "1.2.0")
	if !updateAvailable {
		t.Error("Expected update to be available for 1.0.0 -> 1.2.0")
	}
	if !compatible {
		t.Error("Expected update to be compatible")
	}

	// Test no update needed
	updateAvailable, compatible = checker.compareVersions("1.2.0", "1.0.0")
	if updateAvailable {
		t.Error("Expected no update to be available for 1.2.0 -> 1.0.0")
	}

	// Test version parsing
	version := checker.parseVersion("v1.2.3")
	if version != "1.2.3" {
		t.Errorf("Expected parsed version '1.2.3', got '%s'", version)
	}

	version = checker.parseVersion("1.2.3")
	if version != "1.2.3" {
		t.Errorf("Expected parsed version '1.2.3', got '%s'", version)
	}
}

func TestChecker_FindAssetURLs(t *testing.T) {
	checker := NewChecker(CheckerConfig{})

	release := &GitHubRelease{
		Assets: []GitHubAsset{
			{
				Name:        "core-linux-amd64",
				DownloadURL: "https://example.com/core-linux-amd64",
			},
			{
				Name:        "core-darwin-amd64",
				DownloadURL: "https://example.com/core-darwin-amd64",
			},
			{
				Name:        "checksums.txt",
				DownloadURL: "https://example.com/checksums.txt",
			},
		},
	}

	downloadURL, checksumURL := checker.findAssetURLs(release)

	if checksumURL != "https://example.com/checksums.txt" {
		t.Errorf("Expected checksum URL, got: %s", checksumURL)
	}

	// downloadURL will depend on the runtime OS
	if downloadURL == "" {
		t.Error("Expected to find a download URL for current platform")
	}
}

func TestChecker_CompareVersions(t *testing.T) {
	checker := NewChecker(CheckerConfig{})

	tests := []struct {
		current         string
		latest          string
		wantAvailable   bool
		wantCompatible  bool
	}{
		{"1.0.0", "1.1.0", true, true},   // Update available
		{"1.1.0", "1.0.0", false, true},  // No update
		{"1.0.0", "1.0.0", false, true},  // Same version
		{"dev", "1.0.0", true, true},     // dev version gets update
		{"1.0.0-beta", "1.0.0", true, true}, // Pre-release update
	}

	for _, tt := range tests {
		available, compatible := checker.compareVersions(tt.current, tt.latest)
		if available != tt.wantAvailable {
			t.Errorf("compareVersions(%s, %s): available=%v, want %v",
				tt.current, tt.latest, available, tt.wantAvailable)
		}
		if compatible != tt.wantCompatible {
			t.Errorf("compareVersions(%s, %s): compatible=%v, want %v",
				tt.current, tt.latest, compatible, tt.wantCompatible)
		}
	}
}
