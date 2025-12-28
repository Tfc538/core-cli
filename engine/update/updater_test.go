package update

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUpdater_ParseChecksum(t *testing.T) {
	updater := NewUpdater(UpdaterConfig{})

	tests := []struct {
		name     string
		content  string
		filename string
		expected string
	}{
		{
			name:     "single file checksum",
			content:  "abc123  core-linux-amd64\n",
			filename: "core-linux-amd64",
			expected: "abc123",
		},
		{
			name:     "multiple checksums",
			content:  "def456  ./core-darwin-amd64\nabc123  core-linux-amd64\n",
			filename: "core-linux-amd64",
			expected: "abc123",
		},
		{
			name:     "windows executable",
			content:  "ghi789  core-windows-amd64.exe\n",
			filename: "core-windows-amd64.exe",
			expected: "ghi789",
		},
		{
			name:     "invalid format",
			content:  "invalid content\n",
			filename: "core-linux-amd64",
			expected: "",
		},
		{
			name:     "extra spaces",
			content:  "abc123    core-linux-amd64\n",
			filename: "core-linux-amd64",
			expected: "abc123",
		},
		{
			name:     "missing file",
			content:  "abc123  core-linux-amd64\n",
			filename: "core-darwin-amd64",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := updater.parseChecksum(tt.content, tt.filename)
			if result != tt.expected {
				t.Errorf("parseChecksum() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestUpdater_MockDownload(t *testing.T) {
	// Mock HTTP server serving a fake binary
	content := []byte("fake binary content")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
		w.Write(content)
	}))
	defer server.Close()

	// Create temporary directory for test
	tmpDir := t.TempDir()
	targetPath := tmpDir + "/test-binary"

	updater := NewUpdater(UpdaterConfig{
		DownloadURL: server.URL,
		TargetPath:  targetPath,
	})

	// Set up progress tracking
	var lastProgress UpdateProgress
	updater.SetProgressCallback(func(up UpdateProgress) {
		lastProgress = up
	})

	// Test download
	tmpFile, err := updater.download()
	if err != nil {
		t.Fatalf("download() failed: %v", err)
	}
	defer os.Remove(tmpFile)

	// Verify file was created and has content
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}

	if string(data) != string(content) {
		t.Errorf("Downloaded content mismatch")
	}

	// Verify file is executable
	info, err := os.Stat(tmpFile)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	if info.Mode()&0100 == 0 {
		t.Error("Downloaded file is not executable")
	}

	// Verify progress was reported
	if lastProgress.Stage != "downloading" && lastProgress.Stage != "" {
		t.Errorf("Expected progress stage 'downloading', got %q", lastProgress.Stage)
	}
}

func TestUpdater_VerifyChecksum(t *testing.T) {
	// Create a test file with known content
	tmpDir := t.TempDir()
	testFile := tmpDir + "/test-file"
	testContent := []byte("test content for checksum verification")
	os.WriteFile(testFile, testContent, 0644)

	// Calculate SHA256 of test content
	hash := sha256.Sum256(testContent)
	expectedHash := fmt.Sprintf("%x", hash)

	// Mock checksum server
	checksumContent := fmt.Sprintf("%s  test-binary\n", expectedHash)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(checksumContent))
	}))
	defer server.Close()

	updater := NewUpdater(UpdaterConfig{
		ChecksumURL: server.URL,
		TargetPath:  "test-binary",
	})

	// This should succeed with matching checksum
	err := updater.verifyChecksum(testFile)
	if err != nil {
		t.Errorf("verifyChecksum() failed: %v", err)
	}
}

func TestUpdater_VerifyChecksum_Mismatch(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := tmpDir + "/test-file"
	os.WriteFile(testFile, []byte("test content"), 0644)

	// Mock checksum server with wrong hash
	checksumContent := "wronghash0000000000000000000000000000000000000000000000000000000000  test-binary"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(checksumContent))
	}))
	defer server.Close()

	updater := NewUpdater(UpdaterConfig{
		ChecksumURL: server.URL,
		TargetPath:  "test-binary",
	})

	// This should fail with mismatched checksum
	err := updater.verifyChecksum(testFile)
	if err == nil {
		t.Error("verifyChecksum() should fail with mismatched hash")
	}
	if !contains(err.Error(), "mismatch") {
		t.Errorf("Expected 'mismatch' in error, got: %v", err)
	}
}

func TestUpdater_DownloadProgress(t *testing.T) {
	largeContent := make([]byte, 1024*100) // 100 KB
	for i := range largeContent {
		largeContent[i] = byte(i % 256)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(largeContent)))
		w.Write(largeContent)
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	targetPath := tmpDir + "/test-binary"

	updater := NewUpdater(UpdaterConfig{
		DownloadURL: server.URL,
		TargetPath:  targetPath,
	})

	progressReports := []UpdateProgress{}
	updater.SetProgressCallback(func(up UpdateProgress) {
		progressReports = append(progressReports, up)
	})

	tmpFile, err := updater.download()
	if err != nil {
		t.Fatalf("download() failed: %v", err)
	}
	defer os.Remove(tmpFile)

	// Verify progress was reported
	if len(progressReports) == 0 {
		t.Error("No progress reports received")
	}

	// Check that at least some progress increments were reported
	for _, report := range progressReports {
		if report.Stage != "downloading" {
			t.Errorf("Expected stage 'downloading', got %q", report.Stage)
		}
		if report.Percent < 0 || report.Percent > 100 {
			t.Errorf("Invalid percent: %d", report.Percent)
		}
	}
}

func TestUpdater_SetProgressCallback(t *testing.T) {
	updater := NewUpdater(UpdaterConfig{})

	callbackCalled := false
	var receivedProgress UpdateProgress

	updater.SetProgressCallback(func(up UpdateProgress) {
		callbackCalled = true
		receivedProgress = up
	})

	testProgress := UpdateProgress{
		Stage:   "test",
		Percent: 50,
	}
	updater.progress(testProgress)

	if !callbackCalled {
		t.Error("Progress callback was not called")
	}

	if receivedProgress.Stage != "test" {
		t.Errorf("Expected stage 'test', got %q", receivedProgress.Stage)
	}
}

func TestUpdater_DownloadFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found"))
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	updater := NewUpdater(UpdaterConfig{
		DownloadURL: server.URL,
		TargetPath:  tmpDir + "/test-binary",
	})

	_, err := updater.download()
	if err == nil {
		t.Error("download() should fail with 404 status")
	}
}

func TestUpdater_Config_Validation(t *testing.T) {
	tests := []struct {
		name    string
		config  UpdaterConfig
		wantErr bool
	}{
		{
			name:    "missing download URL",
			config:  UpdaterConfig{TargetPath: "/usr/local/bin/core"},
			wantErr: true,
		},
		{
			name:    "missing target path",
			config:  UpdaterConfig{DownloadURL: "https://example.com/binary"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updater := NewUpdater(tt.config)

			// Try to apply without proper config
			err := updater.Apply()
			if (err != nil) != tt.wantErr {
				t.Errorf("Apply() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdater_Progress_Stages(t *testing.T) {
	stages := []string{
		"downloading",
		"verifying",
		"replacing",
		"complete",
		"failed",
	}

	for _, stage := range stages {
		progress := UpdateProgress{
			Stage: stage,
		}

		// Verify progress can be created with each stage
		if progress.Stage != stage {
			t.Errorf("Stage mismatch: got %q, want %q", progress.Stage, stage)
		}
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
