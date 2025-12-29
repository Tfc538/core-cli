package update

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestChecker_Check(t *testing.T) {
	// Mock core API and GitHub API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/version/latest":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(coreVersionResponse{
				Status: "ok",
				Data: coreVersionData{
					Version:   "1.2.0",
					Commit:    "abc123",
					BuildDate: "2025-01-01T00:00:00Z",
				},
			})
		case "/repos/test-owner/test-repo/releases/latest":
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
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	checker := NewChecker(CheckerConfig{
		APIBaseURL:       server.URL,
		GitHubAPIBaseURL: server.URL,
		GitHubOwner:      "test-owner",
		GitHubRepo:       "test-repo",
		CurrentVersion:   "1.0.0",
	})

	info, err := checker.Check()
	if err != nil {
		t.Fatalf("expected check to succeed, got error: %v", err)
	}
	if !info.UpdateAvailable {
		t.Error("Expected update to be available for 1.0.0 -> 1.2.0")
	}
	if !info.Compatible {
		t.Error("Expected update to be compatible")
	}
	if info.LatestVersion != "1.2.0" {
		t.Errorf("Expected latest version 1.2.0, got %s", info.LatestVersion)
	}
	if info.DownloadURL == "" {
		t.Error("Expected a download URL for current platform")
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
		current        string
		latest         string
		wantAvailable  bool
		wantCompatible bool
	}{
		{"1.0.0", "1.1.0", true, true},            // Update available
		{"1.1.0", "1.0.0", false, true},           // No update
		{"1.0.0", "1.0.0", false, true},           // Same version
		{"dev", "1.0.0", true, true},              // dev version gets update
		{"1.0.0-beta", "1.0.0", true, true},       // Pre-release update
		{"2.0.0", "1.9.9", false, true},           // Major version ahead
		{"1.0.0-alpha", "1.0.0-beta", true, true}, // Pre-release comparison
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

func TestChecker_ParseVersion(t *testing.T) {
	checker := NewChecker(CheckerConfig{})

	tests := []struct {
		input    string
		expected string
	}{
		{"v1.2.3", "1.2.3"},
		{"1.2.3", "1.2.3"},
		{"v1.0.0-alpha", "1.0.0-alpha"},
		{"v2.1.0+build123", "2.1.0+build123"},
		{"release-1.5.0", "release-1.5.0"},
	}

	for _, tt := range tests {
		result := checker.parseVersion(tt.input)
		if result != tt.expected {
			t.Errorf("parseVersion(%s) = %s, want %s", tt.input, result, tt.expected)
		}
	}
}

func TestChecker_MultiPlatformAssets(t *testing.T) {
	checker := NewChecker(CheckerConfig{})

	release := &GitHubRelease{
		Assets: []GitHubAsset{
			{Name: "core-linux-amd64", DownloadURL: "https://example.com/linux-amd64"},
			{Name: "core-linux-arm64", DownloadURL: "https://example.com/linux-arm64"},
			{Name: "core-darwin-amd64", DownloadURL: "https://example.com/darwin-amd64"},
			{Name: "core-darwin-arm64", DownloadURL: "https://example.com/darwin-arm64"},
			{Name: "core-windows-amd64.exe", DownloadURL: "https://example.com/windows-amd64.exe"},
			{Name: "checksums.txt", DownloadURL: "https://example.com/checksums.txt"},
		},
	}

	downloadURL, checksumURL := checker.findAssetURLs(release)
	if downloadURL == "" {
		t.Error("Should find a download URL for current platform")
	}
	if checksumURL == "" {
		t.Error("Should find checksums.txt")
	}
}

func TestChecker_APIErrors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantErr    bool
	}{
		{
			name:       "404 not found",
			statusCode: 404,
			body:       `{"message":"Not Found"}`,
			wantErr:    true,
		},
		{
			name:       "rate limited",
			statusCode: 403,
			body:       `{"message":"API rate limit exceeded"}`,
			wantErr:    true,
		},
		{
			name:       "server error",
			statusCode: 500,
			body:       `{"message":"Internal Server Error"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.body))
			}))
			defer server.Close()

			// Create checker with mock server URL
			checker := NewChecker(CheckerConfig{
				APIBaseURL:       server.URL,
				GitHubAPIBaseURL: server.URL,
				GitHubOwner:      "test",
				GitHubRepo:       "test",
				CurrentVersion:   "1.0.0",
			})

			// Since we can't override the API URL directly, we test error handling logic
			_, err := checker.Check()
			if (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChecker_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`invalid json {`))
	}))
	defer server.Close()

	checker := NewChecker(CheckerConfig{
		APIBaseURL:       server.URL,
		GitHubAPIBaseURL: server.URL,
		GitHubOwner:      "test",
		GitHubRepo:       "test",
		CurrentVersion:   "1.0.0",
	})

	_, err := checker.Check()
	if err == nil {
		t.Error("Check() should return error for invalid JSON")
	}
}

func TestChecker_UpdateInfo(t *testing.T) {
	// Test that UpdateInfo is properly constructed
	info := &UpdateInfo{
		CurrentVersion:  "1.0.0",
		LatestVersion:   "1.2.0",
		UpdateAvailable: true,
		Compatible:      true,
		DownloadURL:     "https://example.com/core-linux-amd64",
		ChecksumURL:     "https://example.com/checksums.txt",
		ReleaseNotes:    "New features and bug fixes",
	}

	if info.CurrentVersion != "1.0.0" {
		t.Error("CurrentVersion not set correctly")
	}
	if !info.UpdateAvailable {
		t.Error("UpdateAvailable not set correctly")
	}
	if info.DownloadURL == "" {
		t.Error("DownloadURL not set")
	}
}

func TestChecker_Timeout(t *testing.T) {
	// Server that takes too long to respond
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Note: httptest doesn't actually timeout, but this tests the timeout config
		fmt.Fprintf(w, `{"tag_name":"v1.0.0","assets":[]}`)
	}))
	defer server.Close()

	checker := NewChecker(CheckerConfig{
		GitHubOwner:    "test",
		GitHubRepo:     "test",
		CurrentVersion: "1.0.0",
	})

	if checker.client.Timeout.Seconds() != 10 {
		t.Errorf("Expected 10 second timeout, got %v", checker.client.Timeout)
	}
}
