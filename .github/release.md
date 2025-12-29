# CORE CLI v0.1.2

## Highlights
- Update checks now read the latest version from the Core API and fetch assets from GitHub Releases.
- New `CORE_GITHUB_API_BASE` env var for targeting a custom GitHub API endpoint.
- Update checker skips the Core API when configured to use GitHub directly.

## Binaries
- `core/core-linux-amd64`
- `core/core-linux-arm64`
- `core/core-darwin-amd64`
- `core/core-darwin-arm64`
- `core/core-windows-amd64.exe`
- `core-backend/core-backend-linux-amd64`
- `core-backend/core-backend-linux-arm64`
- `core-backend/core-backend-darwin-amd64`
- `core-backend/core-backend-darwin-arm64`
- `core-backend/core-backend-windows-amd64.exe`

## Checksums
- `checksums.txt` (SHA256)

## Notes
- Set `CORE_UPDATE_API_BASE` to point at the Core API (used for latest-version checks).
- Set `CORE_GITHUB_TOKEN` (or `GH_TOKEN`/`GITHUB_TOKEN`) for private repo update checks.
