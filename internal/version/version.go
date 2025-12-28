package version

import "fmt"

var (
	// Version is the semantic version of CORE CLI.
	// Injected at build time via -X flag: -X github.com/Tfc538/core-cli/internal/version.Version=1.0.0
	Version = "dev"

	// GitCommit is the short git commit hash.
	// Injected at build time: -X github.com/Tfc538/core-cli/internal/version.GitCommit=$(git rev-parse --short HEAD)
	GitCommit = "unknown"

	// BuildDate is the timestamp when the binary was built (RFC3339 format).
	// Injected at build time: -X github.com/Tfc538/core-cli/internal/version.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)
	BuildDate = "unknown"
)

// Info represents the version information for CORE CLI.
type Info struct {
	Version   string `json:"version"`
	GitCommit string `json:"commit"`
	BuildDate string `json:"build_date"`
}

// Get returns the current version information.
func Get() Info {
	return Info{
		Version:   Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
	}
}

// String returns a human-readable version string.
func (i Info) String() string {
	return fmt.Sprintf("CORE CLI v%s\nCommit: %s\nBuilt: %s", i.Version, i.GitCommit, i.BuildDate)
}
