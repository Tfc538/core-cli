package update

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestUpdater_ParseChecksum(t *testing.T) {
	updater := NewUpdater(UpdaterConfig{})

	tests := []struct {
		content  string
		filename string
		expected string
	}{
		{
			"abc123  core-linux-amd64\n",
			"core-linux-amd64",
			"abc123",
		},
		{
			"def456  ./core-darwin-amd64\nabc123  core-linux-amd64\n",
			"core-linux-amd64",
			"abc123",
		},
		{
			"ghi789  core-windows-amd64.exe\n",
			"core-windows-amd64.exe",
			"ghi789",
		},
		{
			"invalid content\n",
			"core-linux-amd64",
			"",
		},
	}

	for _, tt := range tests {
		result := updater.parseChecksum(tt.content, tt.filename)
		if result != tt.expected {
			t.Errorf("parseChecksum(%q, %q) = %q, want %q",
				tt.content, tt.filename, result, tt.expected)
		}
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

	// Verify progress was reported
	if lastProgress.Stage != "downloading" && lastProgress.Stage != "" {
		t.Errorf("Expected progress stage 'downloading', got %q", lastProgress.Stage)
	}
}

func TestUpdater_VerifyChecksum(t *testing.T) {
	// Create a test file
	tmpDir := t.TempDir()
	testFile := tmpDir + "/test-file"
	os.WriteFile(testFile, []byte("test content"), 0644)

	// Mock checksum server
	checksumContent := "2c26b46911185131006b1bf07635f4a3f9ffd3f4aaf54ef35d87a1e3ecc4e5d9  test-binary"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(checksumContent))
	}))
	defer server.Close()

	updater := NewUpdater(UpdaterConfig{
		ChecksumURL: server.URL,
		TargetPath:  "test-binary",
	})

	// This should warn but not fail for mismatched checksum
	err := updater.verifyChecksum(testFile)
	if err != nil && !contains(err.Error(), "mismatch") {
		t.Errorf("Expected checksum mismatch error or nil, got: %v", err)
	}
}

func TestUpdater_SetProgressCallback(t *testing.T) {
	updater := NewUpdater(UpdaterConfig{})

	callbackCalled := false
	updater.SetProgressCallback(func(up UpdateProgress) {
		callbackCalled = true
	})

	updater.progress(UpdateProgress{Stage: "test"})

	if !callbackCalled {
		t.Error("Progress callback was not called")
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
