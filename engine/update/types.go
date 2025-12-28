package update

// UpdateInfo contains information about a potential update.
type UpdateInfo struct {
	CurrentVersion  string `json:"current_version"`
	LatestVersion   string `json:"latest_version"`
	UpdateAvailable bool   `json:"update_available"`
	Compatible      bool   `json:"compatible"`
	DownloadURL     string `json:"download_url"`
	ChecksumURL     string `json:"checksum_url,omitempty"`
	ReleaseNotes    string `json:"release_notes,omitempty"`
}

// CheckerConfig contains configuration for the update checker.
type CheckerConfig struct {
	GitHubOwner    string
	GitHubRepo     string
	CurrentVersion string
}

// UpdateProgress represents the progress of a download or update operation.
type UpdateProgress struct {
	Stage      string // "downloading", "verifying", "replacing", "complete", "failed"
	Percent    int    // 0-100
	BytesTotal int64
	BytesDone  int64
	Error      error
}

// UpdaterConfig contains configuration for the updater.
type UpdaterConfig struct {
	DownloadURL string // URL to the binary to download
	ChecksumURL string // Optional URL to checksum file
	TargetPath  string // Path to current binary (usually os.Executable())
}

// ProgressCallback is called to report progress during updates.
type ProgressCallback func(UpdateProgress)
