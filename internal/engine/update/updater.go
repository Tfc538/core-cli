package update

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/selfupdate"
)

// Updater handles downloading and applying updates.
type Updater struct {
	config   UpdaterConfig
	client   *http.Client
	progress ProgressCallback
}

// NewUpdater creates a new updater.
func NewUpdater(config UpdaterConfig) *Updater {
	return &Updater{
		config: config,
		client: &http.Client{
			Timeout: 5 * time.Minute,
		},
		progress: func(UpdateProgress) {}, // Default no-op callback
	}
}

// SetProgressCallback sets the callback function for progress updates.
func (u *Updater) SetProgressCallback(cb ProgressCallback) {
	u.progress = cb
}

// Apply downloads the update and applies it, replacing the current binary.
func (u *Updater) Apply() error {
	if u.config.DownloadURL == "" {
		return fmt.Errorf("download URL not specified")
	}

	if u.config.TargetPath == "" {
		return fmt.Errorf("target path not specified")
	}

	// Download binary to temporary file
	tmpFile, err := u.download()
	if err != nil {
		u.progress(UpdateProgress{
			Stage: "failed",
			Error: err,
		})
		return fmt.Errorf("download failed: %w", err)
	}
	defer os.Remove(tmpFile)

	// Verify checksum if available
	if u.config.ChecksumURL != "" {
		if err := u.verifyChecksum(tmpFile); err != nil {
			u.progress(UpdateProgress{
				Stage: "failed",
				Error: err,
			})
			return fmt.Errorf("checksum verification failed: %w", err)
		}

		u.progress(UpdateProgress{
			Stage: "verifying",
		})
	}

	// Apply the update using selfupdate
	u.progress(UpdateProgress{
		Stage: "replacing",
	})

	if err := u.replace(tmpFile); err != nil {
		u.progress(UpdateProgress{
			Stage: "failed",
			Error: err,
		})
		return fmt.Errorf("failed to apply update: %w", err)
	}

	u.progress(UpdateProgress{
		Stage: "complete",
	})

	return nil
}

// download downloads the binary from the configured URL to a temporary file.
func (u *Updater) download() (string, error) {
	u.progress(UpdateProgress{
		Stage: "downloading",
	})

	// Create temporary file
	tmpFile, err := os.CreateTemp("", "core-update-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()

	// Download binary
	resp, err := u.client.Get(u.config.DownloadURL)
	if err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		os.Remove(tmpPath)
		return "", fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	// Stream to file and report progress
	out, err := os.Create(tmpPath)
	if err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	totalSize := resp.ContentLength
	var downloaded int64

	// Create progress reader
	reader := &progressReader{
		reader: resp.Body,
		total:  totalSize,
		update: func(current int64) {
			downloaded = current
			percent := int((downloaded * 100) / totalSize)
			u.progress(UpdateProgress{
				Stage:      "downloading",
				Percent:    percent,
				BytesTotal: totalSize,
				BytesDone:  downloaded,
			})
		},
	}

	if _, err := io.Copy(out, reader); err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// Make file executable
	if err := os.Chmod(tmpPath, 0755); err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("failed to make file executable: %w", err)
	}

	return tmpPath, nil
}

// verifyChecksum verifies the SHA256 checksum of the downloaded file.
func (u *Updater) verifyChecksum(filePath string) error {
	if u.config.ChecksumURL == "" {
		return nil // No checksum to verify
	}

	// Download checksum file
	resp, err := u.client.Get(u.config.ChecksumURL)
	if err != nil {
		// Warn but don't fail if checksum fetch fails
		fmt.Printf("Warning: failed to download checksum: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Warn but don't fail if checksum not found
		fmt.Printf("Warning: checksum file returned status %d\n", resp.StatusCode)
		return nil
	}

	// Parse checksum file (assuming sha256sum format: "hash  filename")
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Warning: failed to read checksum file: %v\n", err)
		return nil
	}

	expectedHash := u.parseChecksum(string(body), filepath.Base(u.config.TargetPath))
	if expectedHash == "" {
		// Couldn't find matching checksum, warn but allow to proceed
		fmt.Println("Warning: could not find matching checksum in file")
		return nil
	}

	// Calculate actual hash
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file for hashing: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Errorf("failed to hash file: %w", err)
	}

	actualHash := fmt.Sprintf("%x", hash.Sum(nil))

	if actualHash != expectedHash {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedHash, actualHash)
	}

	return nil
}

// parseChecksum extracts the hash for a given filename from a checksum file.
func (u *Updater) parseChecksum(checksumContent, filename string) string {
	for _, line := range strings.Split(checksumContent, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 {
			// Format: hash filename
			hash := parts[0]
			name := parts[len(parts)-1]

			// Match filename (handle both exact match and with/without path)
			if name == filename || strings.HasSuffix(name, "/"+filename) {
				return hash
			}
		}
	}
	return ""
}

// replace performs the atomic binary replacement using selfupdate.
func (u *Updater) replace(newBinaryPath string) error {
	// Open the new binary
	newBinary, err := os.Open(newBinaryPath)
	if err != nil {
		return fmt.Errorf("failed to open new binary: %w", err)
	}
	defer newBinary.Close()

	// Apply the update using minio/selfupdate which handles atomic replacement
	err = selfupdate.Apply(newBinary, selfupdate.Options{
		TargetPath: u.config.TargetPath,
	})

	if err != nil {
		return fmt.Errorf("failed to apply update: %w", err)
	}

	return nil
}

// progressReader wraps a reader and reports progress.
type progressReader struct {
	reader io.Reader
	total  int64
	read   int64
	update func(int64)
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	if n > 0 {
		pr.read += int64(n)
		pr.update(pr.read)
	}
	return n, err
}
