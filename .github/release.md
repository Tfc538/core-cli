# CORE CLI v0.1.1

## Highlights
- Monorepo layout with separate CLI and backend entrypoints.
- Backend v0 foundation with health/version endpoints and graceful shutdown.
- Docker image for backend with distroless runtime.
- Update checker now targets `https://api-cli.coreofficialhq.com`.

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
- Set `CORE_GITHUB_TOKEN` (or `GH_TOKEN`/`GITHUB_TOKEN`) for private repo update checks.
