package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
)

const defaultAPIBaseURL = "https://api-cli.coreofficialhq.com"

// Checker is responsible for checking GitHub Releases for updates.
type Checker struct {
	config CheckerConfig
	client *http.Client
}

// NewChecker creates a new update checker.
func NewChecker(config CheckerConfig) *Checker {
	if strings.TrimSpace(config.APIBaseURL) == "" {
		config.APIBaseURL = defaultAPIBaseURL
	}

	return &Checker{
		config: config,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Check performs the update check against GitHub Releases.
func (c *Checker) Check() (*UpdateInfo, error) {
	// Fetch latest release from GitHub API
	release, err := c.getLatestRelease()
	if err != nil {
		return nil, fmt.Errorf("failed to check for updates: %w", err)
	}

	// Parse version from release tag
	latestVersion := c.parseVersion(release.TagName)
	currentVersion := c.config.CurrentVersion

	// Check if update is available
	updateAvailable, compatible := c.compareVersions(currentVersion, latestVersion)

	// Find download URL for current platform
	downloadURL, checksumURL := c.findAssetURLs(release)

	return &UpdateInfo{
		CurrentVersion:  currentVersion,
		LatestVersion:   latestVersion,
		UpdateAvailable: updateAvailable,
		Compatible:      compatible,
		DownloadURL:     downloadURL,
		ChecksumURL:     checksumURL,
		ReleaseNotes:    release.Body,
	}, nil
}

// GitHubRelease represents a GitHub release response.
type GitHubRelease struct {
	TagName string        `json:"tag_name"`
	Body    string        `json:"body"`
	Assets  []GitHubAsset `json:"assets"`
}

// GitHubAsset represents a release asset.
type GitHubAsset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
}

// getLatestRelease fetches the latest release from GitHub.
func (c *Checker) getLatestRelease() (*GitHubRelease, error) {
	baseURL := strings.TrimRight(c.config.APIBaseURL, "/")
	url := fmt.Sprintf("%s/repos/%s/%s/releases/latest",
		baseURL, c.config.GitHubOwner, c.config.GitHubRepo)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub request: %w", err)
	}

	token := strings.TrimSpace(c.config.GitHubToken)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Accept", "application/vnd.github+json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch GitHub release: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned %d: %s", resp.StatusCode, string(body))
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to parse GitHub response: %w", err)
	}

	return &release, nil
}

// parseVersion extracts a semantic version from a git tag.
func (c *Checker) parseVersion(tag string) string {
	// Remove leading 'v' if present
	if strings.HasPrefix(tag, "v") {
		tag = tag[1:]
	}
	return tag
}

// compareVersions compares two semantic versions.
// Returns (updateAvailable, compatible).
func (c *Checker) compareVersions(currentStr, latestStr string) (bool, bool) {
	currentVersion, err := semver.NewVersion(currentStr)
	if err != nil {
		// If current version is unparseable (e.g., "dev"), always consider update compatible
		return true, true
	}

	latestVersion, err := semver.NewVersion(latestStr)
	if err != nil {
		// If latest version is unparseable, no update available
		return false, false
	}

	// Update is available if latest > current
	updateAvailable := latestVersion.GreaterThan(currentVersion)

	// Compatible if major version matches (basic compatibility check)
	// For now, always consider compatible if versions are parseable
	compatible := true

	return updateAvailable, compatible
}

// findAssetURLs locates the correct binary for the current platform.
func (c *Checker) findAssetURLs(release *GitHubRelease) (downloadURL, checksumURL string) {
	goos := runtime.GOOS
	goarch := runtime.GOARCH

	// Normalize architecture names if needed
	if goarch == "amd64" {
		goarch = "amd64"
	}

	// Build expected binary name patterns
	patterns := []string{
		fmt.Sprintf("core-%s-%s", goos, goarch),
		fmt.Sprintf("core-%s-%s.exe", goos, goarch),
		fmt.Sprintf("core-%s-%s.tar.gz", goos, goarch),
		fmt.Sprintf("core-%s-%s.zip", goos, goarch),
	}

	for _, asset := range release.Assets {
		// Look for binary matching current platform
		for _, pattern := range patterns {
			if strings.Contains(asset.Name, pattern) {
				downloadURL = asset.DownloadURL
				break
			}
		}

		// Look for checksums file
		if strings.Contains(asset.Name, "checksums") {
			checksumURL = asset.DownloadURL
		}
	}

	return downloadURL, checksumURL
}
