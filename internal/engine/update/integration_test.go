// +build integration

package update

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TestIntegration_FullUpdateWorkflow tests the complete update workflow.
func TestIntegration_FullUpdateWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Create a mock GitHub release server
	binaryContent := []byte("updated binary content v1.2.0")
	checksumContent := fmt.Sprintf("%x  core-linux-amd64\n", sha256.Sum256(binaryContent))

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/releases/latest" {
			release := GitHubRelease{
				TagName: "v1.2.0",
				Body:    "## Version 1.2.0\n- New features\n- Bug fixes",
				Assets: []GitHubAsset{
					{
						Name:        "core-linux-amd64",
						DownloadURL: fmt.Sprintf("%s/download/binary", server.URL),
					},
					{
						Name:        "checksums.txt",
						DownloadURL: fmt.Sprintf("%s/download/checksums", server.URL),
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(release)
		} else if r.URL.Path == "/download/binary" {
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(binaryContent)))
			w.Write(binaryContent)
		} else if r.URL.Path == "/download/checksums" {
			w.Write([]byte(checksumContent))
		}
	}))
	defer server.Close()

	tmpDir := t.TempDir()

	// Step 1: Check for updates
	checker := NewChecker(CheckerConfig{
		GitHubOwner:    "test",
		GitHubRepo:     "test",
		CurrentVersion: "1.0.0",
	})

	updateInfo, err := checker.Check()
	if err != nil {
		t.Fatalf("Check() failed: %v", err)
	}

	// Verify update is available
	if !updateInfo.UpdateAvailable {
		t.Error("Update should be available")
	}
	if updateInfo.LatestVersion != "1.2.0" {
		t.Errorf("Expected version 1.2.0, got %s", updateInfo.LatestVersion)
	}

	// Step 2: Download and verify
	updater := NewUpdater(UpdaterConfig{
		DownloadURL: fmt.Sprintf("%s/download/binary", server.URL),
		ChecksumURL: fmt.Sprintf("%s/download/checksums", server.URL),
		TargetPath:  tmpDir + "/core",
	})

	// Track progress events
	progressEvents := []UpdateProgress{}
	updater.SetProgressCallback(func(up UpdateProgress) {
		progressEvents = append(progressEvents, up)
	})

	// Download the update (we won't actually apply it to avoid modifying system files)
	downloadedFile, err := updater.download()
	if err != nil {
		t.Fatalf("download() failed: %v", err)
	}
	defer os.Remove(downloadedFile)

	// Verify download
	content, err := os.ReadFile(downloadedFile)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}
	if string(content) != string(binaryContent) {
		t.Error("Downloaded content mismatch")
	}

	// Step 3: Verify checksum
	err = updater.verifyChecksum(downloadedFile)
	if err != nil {
		t.Errorf("Checksum verification failed: %v", err)
	}

	// Verify progress was reported
	if len(progressEvents) == 0 {
		t.Error("No progress events reported")
	}

	downloadingEvents := 0
	for _, event := range progressEvents {
		if event.Stage == "downloading" {
			downloadingEvents++
		}
	}

	if downloadingEvents == 0 {
		t.Error("No download progress events reported")
	}
}

// TestIntegration_UpdateCheckAndComparison tests version checking and comparison.
func TestIntegration_UpdateCheckAndComparison(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	checker := NewChecker(CheckerConfig{
		GitHubOwner:    "Tfc538",
		GitHubRepo:     "core-cli",
		CurrentVersion: "1.0.0",
	})

	tests := []struct {
		name        string
		current     string
		latest      string
		expectErr   bool
		expectAvail bool
	}{
		{
			name:        "newer version available",
			current:     "1.0.0",
			latest:      "1.1.0",
			expectErr:   false,
			expectAvail: true,
		},
		{
			name:        "already on latest",
			current:     "1.1.0",
			latest:      "1.1.0",
			expectErr:   false,
			expectAvail: false,
		},
		{
			name:        "ahead of latest",
			current:     "1.2.0",
			latest:      "1.1.0",
			expectErr:   false,
			expectAvail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			available, compatible := checker.compareVersions(tt.current, tt.latest)

			if available != tt.expectAvail {
				t.Errorf("compareVersions() available = %v, want %v", available, tt.expectAvail)
			}

			if !compatible {
				t.Error("compareVersions() should be compatible")
			}
		})
	}
}

// TestIntegration_MultiPlatformAssetSelection tests asset selection for different platforms.
func TestIntegration_MultiPlatformAssetSelection(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	checker := NewChecker(CheckerConfig{
		GitHubOwner:    "test",
		GitHubRepo:     "test",
		CurrentVersion: "1.0.0",
	})

	release := &GitHubRelease{
		TagName: "v1.1.0",
		Assets: []GitHubAsset{
			{Name: "core-linux-amd64", DownloadURL: "https://example.com/linux-amd64"},
			{Name: "core-linux-arm64", DownloadURL: "https://example.com/linux-arm64"},
			{Name: "core-darwin-amd64", DownloadURL: "https://example.com/darwin-amd64"},
			{Name: "core-darwin-arm64", DownloadURL: "https://example.com/darwin-arm64"},
			{Name: "core-windows-amd64.exe", DownloadURL: "https://example.com/windows-amd64"},
			{Name: "checksums.txt", DownloadURL: "https://example.com/checksums.txt"},
		},
	}

	downloadURL, checksumURL := checker.findAssetURLs(release)

	if downloadURL == "" {
		t.Error("Should find a download URL for current platform")
	}

	if checksumURL != "https://example.com/checksums.txt" {
		t.Errorf("Expected checksums.txt URL, got: %s", checksumURL)
	}

	// Verify we got a reasonable URL
	if !contains(downloadURL, "example.com") {
		t.Errorf("Expected example.com in URL, got: %s", downloadURL)
	}
}

// TestIntegration_ChecksumValidation tests checksum parsing and validation.
func TestIntegration_ChecksumValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tmpDir := t.TempDir()

	// Create test files with different content
	file1 := tmpDir + "/core-linux-amd64"
	file2 := tmpDir + "/core-darwin-amd64"

	content1 := []byte("binary for linux amd64")
	content2 := []byte("binary for darwin amd64")

	os.WriteFile(file1, content1, 0644)
	os.WriteFile(file2, content2, 0644)

	// Create checksum file
	hash1 := fmt.Sprintf("%x", sha256.Sum256(content1))
	hash2 := fmt.Sprintf("%x", sha256.Sum256(content2))

	checksumContent := fmt.Sprintf("%s  core-linux-amd64\n%s  core-darwin-amd64\n", hash1, hash2)

	updater := NewUpdater(UpdaterConfig{})

	// Test parsing multiple checksums
	found1 := updater.parseChecksum(checksumContent, "core-linux-amd64")
	found2 := updater.parseChecksum(checksumContent, "core-darwin-amd64")

	if found1 != hash1 {
		t.Errorf("First checksum mismatch: got %s, want %s", found1, hash1)
	}

	if found2 != hash2 {
		t.Errorf("Second checksum mismatch: got %s, want %s", found2, hash2)
	}

	// Test missing checksum
	missing := updater.parseChecksum(checksumContent, "core-windows-amd64.exe")
	if missing != "" {
		t.Errorf("Should not find checksum for missing file, got: %s", missing)
	}
}

// TestIntegration_ErrorRecovery tests error handling and recovery.
func TestIntegration_ErrorRecovery(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		server  func() *httptest.Server
		wantErr bool
	}{
		{
			name: "network error recovery",
			server: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))
			},
			wantErr: true,
		},
		{
			name: "missing file error",
			server: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.server()
			defer server.Close()

			updater := NewUpdater(UpdaterConfig{
				DownloadURL: server.URL,
				TargetPath:  tmpDir + "/core",
			})

			_, err := updater.download()
			if (err != nil) != tt.wantErr {
				t.Errorf("download() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestIntegration_ProgressTracking tests progress reporting throughout the workflow.
func TestIntegration_ProgressTracking(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Create large binary to ensure progress updates
	binaryContent := make([]byte, 1024*500) // 500 KB
	for i := range binaryContent {
		binaryContent[i] = byte(i % 256)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(binaryContent)))
		w.Write(binaryContent)
	}))
	defer server.Close()

	tmpDir := t.TempDir()

	updater := NewUpdater(UpdaterConfig{
		DownloadURL: server.URL,
		TargetPath:  tmpDir + "/core",
	})

	progressEvents := []UpdateProgress{}
	updater.SetProgressCallback(func(up UpdateProgress) {
		progressEvents = append(progressEvents, up)
	})

	_, err := updater.download()
	if err != nil {
		t.Fatalf("download() failed: %v", err)
	}

	// Verify progress events
	if len(progressEvents) == 0 {
		t.Fatal("No progress events received")
	}

	// Check stage
	for _, event := range progressEvents {
		if event.Stage != "downloading" {
			t.Errorf("Expected stage 'downloading', got %q", event.Stage)
		}

		// Verify percent is reasonable
		if event.Percent < 0 || event.Percent > 100 {
			t.Errorf("Invalid percent: %d", event.Percent)
		}

		// Verify bytes are non-negative
		if event.BytesDone < 0 || event.BytesTotal < 0 {
			t.Error("Negative byte counts in progress")
		}

		if event.BytesDone > event.BytesTotal {
			t.Error("BytesDone exceeds BytesTotal")
		}
	}

	// Verify final progress is at 100%
	if progressEvents[len(progressEvents)-1].Percent != 100 {
		t.Errorf("Final progress should be 100%%, got %d%%", progressEvents[len(progressEvents)-1].Percent)
	}
}

func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
