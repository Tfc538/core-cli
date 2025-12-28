# CORE CLI - Testing Guide

This document describes the comprehensive test suites for CORE CLI, including unit tests, integration tests, and end-to-end tests.

## Test Organization

```
Tests are organized into three categories:

1. Unit Tests       - Test individual functions and components in isolation
2. Integration Tests - Test multiple components working together
3. End-to-End Tests  - Test complete workflows from CLI to output
```

## Running Tests

### Run All Tests

```bash
# Run all tests (unit + integration + e2e)
go test -v ./...

# Run with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run Only Unit Tests

```bash
# Skip integration and e2e tests
go test -v ./... -short
```

### Run Only Integration Tests

```bash
# Run integration tests (marked with +build integration)
go test -v ./... -tags=integration
```

### Run Only End-to-End Tests

```bash
# Build the binary first
make build VERSION=0.1.0

# Run e2e tests (marked with +build e2e)
go test -v ./... -tags=e2e
```

### Run Specific Test

```bash
# Run a specific unit test
go test -v -run TestChecker_CompareVersions ./internal/engine/update

# Run all tests in a specific package
go test -v ./version
go test -v ./internal/engine/update
```

## Unit Tests

Unit tests focus on testing individual functions and components in isolation.

### Coverage

| Package | Tests | Coverage |
|---------|-------|----------|
| `version` | 4 | Version metadata, JSON marshaling, field names |
| `internal/engine/update` | 11 | Version comparison, checksum parsing, downloads |

### Key Test Files

- `internal/version/version_test.go` - Tests for version package
- `internal/engine/update/checker_test.go` - Tests for update checker
- `internal/engine/update/updater_test.go` - Tests for updater

### Unit Test Examples

#### Testing Version Metadata

```go
func TestInfo_String(t *testing.T) {
    info := version.Info{
        Version:   "1.0.0",
        GitCommit: "abc123",
        BuildDate: "2025-12-28T10:30:00Z",
    }

    str := info.String()
    if !strings.Contains(str, "CORE CLI") {
        t.Error("String should contain CORE CLI")
    }
}
```

#### Testing Version Comparison

```go
func TestChecker_CompareVersions(t *testing.T) {
    checker := NewChecker(CheckerConfig{})

    available, _ := checker.compareVersions("1.0.0", "1.1.0")
    if !available {
        t.Error("1.1.0 should be available for 1.0.0")
    }
}
```

#### Testing Checksum Parsing

```go
func TestUpdater_ParseChecksum(t *testing.T) {
    updater := NewUpdater(UpdaterConfig{})

    content := "abc123  core-linux-amd64\n"
    hash := updater.parseChecksum(content, "core-linux-amd64")

    if hash != "abc123" {
        t.Errorf("Expected abc123, got %s", hash)
    }
}
```

### Running Unit Tests

```bash
# Run all unit tests
go test -v -short ./...

# Run version package tests
go test -v -short ./version

# Run internal/engine/update tests
go test -v -short ./internal/engine/update
```

## Integration Tests

Integration tests verify that multiple components work correctly together.

### Features Tested

- **Full update workflow** - Check, download, verify, prepare apply
- **Version checking** - Version comparison across semver ranges
- **Multi-platform asset selection** - Correct binary for each OS/arch
- **Checksum validation** - Parsing and validating checksums
- **Error recovery** - Graceful handling of network errors
- **Progress tracking** - Progress reporting throughout workflow

### Key Test Files

- `internal/engine/update/integration_test.go` - Integration tests marked with `+build integration`

### Integration Test Examples

#### Testing Full Update Workflow

```go
func TestIntegration_FullUpdateWorkflow(t *testing.T) {
    // Create mock GitHub server
    // Check for updates
    // Download binary
    // Verify checksum
    // Verify progress events
}
```

#### Testing Multi-Platform Asset Selection

```go
func TestIntegration_MultiPlatformAssetSelection(t *testing.T) {
    checker := NewChecker(CheckerConfig{...})

    release := &GitHubRelease{
        Assets: []GitHubAsset{
            {Name: "core-linux-amd64", ...},
            {Name: "core-darwin-amd64", ...},
            {Name: "core-windows-amd64.exe", ...},
            {Name: "checksums.txt", ...},
        },
    }

    // Should select correct asset for current platform
    downloadURL, checksumURL := checker.findAssetURLs(release)
}
```

### Running Integration Tests

```bash
# Run all integration tests
go test -v -tags=integration ./...

# Run specific integration test
go test -v -run TestIntegration_FullUpdateWorkflow -tags=integration ./internal/engine/update

# Run integration tests with coverage
go test -cover -tags=integration ./...
```

## End-to-End Tests

End-to-end tests verify complete CLI workflows from command invocation to output.

### Features Tested

- **Version command** - Human-readable and JSON output
- **Help commands** - Proper help text for commands
- **Argument parsing** - Correct handling of various arguments
- **Exit codes** - Proper exit codes for success and failure
- **Output redirection** - Stdout/stderr handling
- **Binary execution** - Binary is executable and works

### Key Test Files

- `e2e_test.go` - End-to-end tests marked with `+build e2e`

### End-to-End Test Examples

#### Testing Version Command

```go
func TestE2E_VersionCommand(t *testing.T) {
    cmd := exec.Command("./core", "version")
    output, err := cmd.Output()

    if !strings.Contains(string(output), "CORE CLI") {
        t.Error("Expected CORE CLI in output")
    }
}
```

#### Testing JSON Output

```go
func TestE2E_VersionJSONOutput(t *testing.T) {
    cmd := exec.Command("./core", "version", "--json")
    output, err := cmd.Output()

    var info struct {
        Version   string `json:"version"`
        Commit    string `json:"commit"`
        BuildDate string `json:"build_date"`
    }

    json.Unmarshal(output, &info)
    // Verify fields are populated
}
```

#### Testing Help Command

```go
func TestE2E_HelpCommand(t *testing.T) {
    cmd := exec.Command("./core", "--help")
    output, err := cmd.Output()

    expected := []string{"CORE CLI", "Usage:", "Commands:"}
    for _, s := range expected {
        if !strings.Contains(string(output), s) {
            t.Errorf("Expected %s in help", s)
        }
    }
}
```

### Running End-to-End Tests

```bash
# Build the binary first
make build VERSION=0.1.0

# Run all e2e tests
go test -v -tags=e2e ./...

# Run specific e2e test
go test -v -run TestE2E_VersionCommand -tags=e2e

# Run e2e tests with coverage
go test -cover -tags=e2e ./...
```

## Test Utilities

Test utilities in `internal/testutil/helpers.go` provide common test functions:

```go
// Create test files and directories
tmpFile := testutil.CreateTestFile(t, []byte("content"))
tmpDir := testutil.CreateTestDir(t)

// Calculate SHA256 hash
hash := testutil.CalculateHash([]byte("content"))

// Create mock servers
server := testutil.CreateMockServer([]byte("content"), "filename")
errorServer := testutil.CreateErrorServer(http.StatusNotFound)

// Assertions
testutil.AssertEqual(t, expected, actual, "message")
testutil.AssertError(t, err, "message")
testutil.AssertStringContains(t, str, substr, "message")
testutil.AssertFileExists(t, path, "message")
```

## Test Coverage Goals

| Component | Target Coverage |
|-----------|-----------------|
| `version` | 100% |
| `internal/engine/update` | 90%+ |
| `cli` | 80%+ |
| `tui` | 60%+ (limited by Bubble Tea interaction model) |

### Checking Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage in browser
go tool cover -html=coverage.out

# Check coverage for specific package
go test -cover ./version
go test -cover ./internal/engine/update
```

## Mocking Strategies

### Mock HTTP Servers

For tests that require HTTP interactions:

```go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}))
defer server.Close()
```

### Mock File System

For tests that work with files:

```go
tmpDir := t.TempDir()  // Automatically cleaned up
tmpFile, _ := os.CreateTemp("", "test-*")  // Manual cleanup needed
```

### Progress Tracking

For tests that verify progress callbacks:

```go
var progressEvents []UpdateProgress
updater.SetProgressCallback(func(up UpdateProgress) {
    progressEvents = append(progressEvents, up)
})
```

## Testing Best Practices

1. **Isolation**: Each test should be independent and not rely on other tests
2. **Cleanup**: Use `t.TempDir()` or defer cleanup for resources
3. **Assertions**: Use helper functions from `testutil` for consistent assertions
4. **Naming**: Use descriptive test names: `TestFunction_Condition_Expected`
5. **Coverage**: Aim for high coverage on critical paths, especially security-related code
6. **Mocking**: Mock external dependencies (HTTP, file system) to ensure reproducibility
7. **Edge Cases**: Test boundary conditions and error cases

## Continuous Integration

The test suite is designed to run in CI/CD environments:

```bash
# Run all tests (short mode for quick feedback)
go test -short ./...

# Run with coverage
go test -cover ./...

# Run integration tests
go test -tags=integration ./...

# Build and run e2e tests
make build VERSION=$(git describe --tags)
go test -tags=e2e ./...
```

## Troubleshooting Tests

### Port Already in Use

If you see "address already in use" errors, try:
- Wait a moment before retrying tests
- Use unique port numbers (httptest.NewServer uses random ports by default)

### File Permission Errors

On Unix systems, ensure the core binary is executable:
```bash
chmod +x ./core
```

### Timing Issues

Some tests may be flaky due to timing:
- Use `t.TempDir()` which handles cleanup automatically
- Avoid hardcoded sleep durations
- Use channels or conditions for synchronization

### Coverage Gaps

To find untested code:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Adding New Tests

When adding new functionality:

1. **Write unit tests first** - Test the component in isolation
2. **Add integration tests** - Test with other components
3. **Add e2e tests** - Test from the CLI perspective
4. **Update documentation** - Add examples to TESTING.md

Example workflow for new feature:

```go
// 1. Unit test
func TestMyFeature_BasicFunctionality(t *testing.T) {
    result := MyFeature()
    AssertEqual(t, expected, result, "basic case")
}

// 2. Integration test (integration_test.go)
func TestIntegration_MyFeatureWithOtherComponent(t *testing.T) {
    // Test interaction with other parts
}

// 3. E2E test (e2e_test.go)
func TestE2E_MyFeatureFromCLI(t *testing.T) {
    cmd := exec.Command("./core", "mycommand")
    // Verify output
}
```

## Test Statistics

As of the latest update:

- **Total Test Files**: 5
- **Total Test Functions**: 30+
- **Code Coverage**: 85%+ on critical paths
- **Average Test Time**: <1 second
- **Integration Test Time**: <5 seconds
- **E2E Test Time**: <10 seconds

## Further Reading

- [Go Testing Package](https://pkg.go.dev/testing)
- [httptest Package](https://pkg.go.dev/net/http/httptest)
- [Go Benchmark Testing](https://pkg.go.dev/testing#hdr-Benchmarks)
- [Table-Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)
