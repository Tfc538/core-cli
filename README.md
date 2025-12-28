# CORE CLI

A local, intent-driven developer control plane with a CLI-first interface and an optional, premium TUI.

## Features

- **Args-First Design**: All functionality accessible via command-line arguments for headless usage
- **Optional TUI**: Launch an interactive terminal UI when invoked without arguments
- **Version Management**: Display current version with build metadata
- **Update Checking**: Check for new versions from GitHub Releases
- **Safe Self-Updates**: Atomic binary replacement with automatic rollback on failure
- **Cross-Platform**: Supports Linux, macOS, and Windows

## Installation

Download the latest release for your platform from [GitHub Releases](https://github.com/Tfc538/core-cli/releases).

For Linux/macOS:
```bash
chmod +x core-linux-amd64  # or core-darwin-amd64
./core-linux-amd64 version
```

For Windows:
```cmd
core-windows-amd64.exe version
```

## Usage

### Interactive Mode (No Arguments)

Launch the premium TUI by running the binary without arguments:

```bash
core
```

This opens an interactive terminal UI where you can:
- View current version and update status
- Press 'u' to check for or apply updates
- Press 'q' to quit

### CLI Mode (With Arguments)

All commands work headlessly via arguments:

#### Check Version

```bash
# Human-readable format
core version

# JSON output
core version --json
```

Example output:
```
CORE CLI v0.1.0
Commit: a1b2c3d
Built: 2025-12-28T19:26:19Z
```

JSON output:
```json
{
  "version": "0.1.0",
  "commit": "a1b2c3d",
  "build_date": "2025-12-28T19:26:19Z"
}
```

#### Check for Updates

```bash
# Check if a newer version is available
core update check

# JSON output (useful for scripts)
core update check --json
```

For private repos or higher rate limits, set `CORE_GITHUB_TOKEN` (or `GH_TOKEN`/`GITHUB_TOKEN`) with access to the repo.

Example output:
```
Current version: 0.1.0
Latest version:  0.2.0

✓ Update available!

Run 'core update apply' to update.
```

#### Apply Update

```bash
# Apply the latest version (with confirmation prompt)
core update apply

# Skip confirmation and apply immediately
core update apply --yes
```

Example output:
```
Update Available

Current version     : 0.1.0
Latest version      : 0.2.0
Target location     : /usr/local/bin/core

Continue with update? [y/N]: y

Starting update

⬇  Downloading... 45% (4/10 MB)
✓  Checksum verified
⬇  Replacing binary...
✓  Update complete!

✓ CORE CLI updated to v0.2.0
```

## Backend Service

The repo also ships a minimal backend service for local development and future distribution metadata.

### Run Locally

```bash
go run ./cmd/core-backend
```

Configuration (optional):

- `CORE_BACKEND_HOST` (default `127.0.0.1`)
- `CORE_BACKEND_PORT` (default `8080`)
- `CORE_BACKEND_SHUTDOWN_TIMEOUT` (default `5s`)

### Endpoints

- `GET /healthz`
- `GET /api/v1/version/latest`
- `GET /api/v1/version/{version}`

## Building from Source

### Prerequisites

- Go 1.23 or later
- Make (optional, but recommended)

### Quick Build

```bash
make build VERSION=0.2.0
```

This builds the binary for your current platform as `./core`.

### Build Backend

```bash
make build-backend VERSION=0.2.0
```

This builds the backend binary as `./core-backend`.

### Build for All Platforms

```bash
make build-all VERSION=0.2.0
```

Outputs binaries to `./dist/`:
- `core-linux-amd64`
- `core-linux-arm64`
- `core-darwin-amd64`
- `core-darwin-arm64`
- `core-windows-amd64.exe`

### Generate Checksums

```bash
make checksums
```

Creates `./dist/checksums.txt` with SHA256 hashes for all binaries.

### Testing

```bash
# Run all tests
make test

# Or use go directly
go test -v ./...
```

## Architecture

CORE CLI follows a clean, layered architecture:

```
core-cli/
├── cmd/core/             # CLI entry point (args routing)
├── cmd/core-backend/     # Backend entry point
└── internal/
    ├── backend/          # Backend HTTP/API + services
    ├── cli/              # Cobra command definitions (CLI layer)
    ├── config/           # Env-based configuration
    ├── engine/update/    # Pure business logic for updates
    ├── tui/              # Bubble Tea TUI (presentation layer)
    └── version/          # Version metadata package
```

### Design Principles

1. **Args-First**: Commands always work via arguments, TUI is optional
2. **Separation of Concerns**: Engine logic separate from CLI/TUI presentation
3. **No Duplication**: CLI and TUI both use identical engine code
4. **Explicit Updates**: Users must opt-in to updates, no background automation
5. **Safe Defaults**: Confirmation prompts, atomic replacement, automatic rollback

### Key Components

- **internal/version/**: Build-time version injection via ldflags
- **internal/backend/**: API handlers, services, and storage/telemetry stubs
- **internal/engine/update/checker**: GitHub Releases API integration
- **internal/engine/update/updater**: Download, verify, and atomic binary replacement
- **internal/cli/**: Cobra command structure with consistent output formatting
- **internal/tui/**: Bubble Tea application with status bar and update view

## Updating

### Checking for Updates

The CLI respects the `--no-update-check` flag and will not check for updates in scripted/non-interactive contexts.

When running the TUI, update checks happen automatically in the background on startup.

### Applying Updates

Updates are **never applied automatically**. To update:

1. **CLI**: `core update apply`
2. **TUI**: Press 'u' when an update is available

Updates are applied safely with:
- SHA256 checksum verification
- Atomic binary replacement
- Automatic rollback on failure
- Binary backup (`.old` extension on Unix)

## Development

### Project Structure

```
cmd/core/main.go                    # CLI entry point with CLI/TUI routing
cmd/core-backend/main.go            # Backend entry point

internal/backend/api/               # HTTP handlers + middleware
internal/backend/service/           # Backend business logic
internal/backend/storage/           # Interfaces for storage backends
internal/backend/telemetry/         # Telemetry stubs
internal/config/backend.go          # Backend env config
internal/version/version.go         # Version constants and Info struct
internal/version/version_test.go

internal/engine/update/types.go     # UpdateInfo, UpdateProgress types
internal/engine/update/checker.go   # GitHub Releases API integration
internal/engine/update/checker_test.go
internal/engine/update/updater.go   # Download, verify, replace logic
internal/engine/update/updater_test.go

internal/cli/root.go                # Root command setup
internal/cli/version.go             # 'core version' command
internal/cli/update.go              # 'core update' parent command
internal/cli/update_check.go        # 'core update check' command
internal/cli/update_apply.go        # 'core update apply' command
internal/cli/output.go              # Output formatting utilities

internal/tui/app.go                 # Main Bubble Tea app
internal/tui/styles.go              # Lip Gloss styles
internal/tui/statusbar.go           # Status bar component
internal/tui/update_view.go         # Update progress view

Makefile                            # Build automation
```

### Adding New Commands

1. Create `internal/cli/mycmd.go` with `NewMyCmd()` function
2. Add to root command in `internal/cli/root.go`
3. Follow existing patterns for output and error handling

### Testing

- Unit tests exist for version, checker, and updater
- Backend unit and integration tests live under `internal/backend/`
- Tests use mock HTTP servers for network operations
- All tests can run offline without external dependencies

## FAQ

### How is the version set?

Version is embedded at build time using Go's `-ldflags` flag:

```bash
-X github.com/Tfc538/core-cli/internal/version.Version=1.0.0
-X github.com/Tfc538/core-cli/internal/version.GitCommit=$(git rev-parse --short HEAD)
-X github.com/Tfc538/core-cli/internal/version.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)
```

The Makefile handles this automatically.

### Can I disable update checks?

Yes, pass `--no-update-check` flag (to be implemented in future versions).

### What platforms are supported?

- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

### What happens if an update fails?

The updater automatically:
1. Preserves a backup of the current binary (`.old`)
2. Restores it if replacement fails
3. Never corrupts or removes your binary

### How often are releases published?

Release schedule depends on the maintainers. Check GitHub Releases for latest versions.

### Can I use this in scripts?

Yes! Use the `--json` flag for machine-readable output:

```bash
core version --json | jq .version
core update check --json | jq .updateAvailable
```

## License

[Your License Here]

## Contributing

[Contributing guidelines go here]

## Support

For issues, bug reports, or feature requests, please open an issue on [GitHub](https://github.com/Tfc538/core-cli/issues).
