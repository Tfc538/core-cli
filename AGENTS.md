# Repository Guidelines

## Project Structure & Module Organization
- `cmd/core/main.go`: entry point, routes CLI vs TUI.
- `cmd/core-backend/main.go`: backend entry point.
- `internal/cli/`: Cobra commands and output formatting.
- `internal/engine/update/`: update checker/updater business logic.
- `internal/backend/`: API handlers, services, and storage/telemetry stubs.
- `internal/config/`: backend config loading.
- `internal/tui/`: Bubble Tea UI and Lip Gloss styles.
- `internal/version/`: build-time version metadata.
- Tests live alongside packages (`*_test.go`), plus `e2e_test.go` and `internal/testutil/`.

## Build, Test, and Development Commands
- `make build VERSION=0.2.0`: build current platform binary as `./core`.
- `make build-backend VERSION=0.2.0`: build backend binary as `./core-backend`.
- `make build-all VERSION=0.2.0`: cross-compile into `./dist/`.
- `make test`: run all tests (`go test -v ./...`).
- `make checksums`: build all binaries and create `./dist/checksums.txt`.
- `go test -v -tags=integration ./...` and `go test -v -tags=e2e ./...`: run tagged suites.

## Coding Style & Naming Conventions
- Follow existing patterns; max line length is 100 (guideline).
- Naming: `camelCase` for locals, `PascalCase` for exported types/functions.
- Use descriptive package names (e.g., `version`, `update`).
- Comment exported symbols; avoid obvious comments.
- Wrap errors with context (e.g., `fmt.Errorf("...: %w", err)`).

## Testing Guidelines
- Framework: Go `testing` with `httptest` for HTTP mocks.
- Test names: `TestFunction_Condition_Expected`.
- Coverage targets: `internal/version` 100%, `internal/engine/update` 90%+, `internal/cli` 80%+, `internal/tui` 60%+, `internal/backend` baseline coverage for handlers/services.
- Check coverage: `go test -coverprofile=coverage.out ./...` then `go tool cover -html=coverage.out`.

## Commit & Pull Request Guidelines
- Commit messages: concise, imperative, pattern like `Add/Implement/Fix/Refactor X`.
  Example: `Add GitHub token support for update checks`.
- PRs: one feature/fix per PR, clear title, describe why/what, link issues, include testing notes.
- Ensure tests pass before submitting.

## Security & Configuration Tips
- Use `CORE_GITHUB_TOKEN` (or `GH_TOKEN`/`GITHUB_TOKEN`) for update checks against GitHub Releases.
- Don’t log sensitive data; validate inputs at boundaries.
