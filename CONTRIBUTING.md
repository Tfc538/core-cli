# Contributing to CORE CLI

Thank you for your interest in contributing! This document provides guidelines for development, testing, and submitting changes.

## Getting Started

### Prerequisites

- Go 1.23 or later
- Git
- Make (optional)

### Setup

```bash
git clone https://github.com/Tfc538/core-cli.git
cd core-cli
go mod download
```

### Building

```bash
# Build for current platform
make build

# Build with specific version
make build VERSION=0.2.0

# Build for all platforms
make build-all VERSION=0.2.0
```

The binary will be created as `./core` (or `./core.exe` on Windows).

### Testing

```bash
# Run all tests
make test

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test -v ./version
go test -v ./engine/update
```

All tests must pass before submitting changes.

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/my-feature
```

### 2. Make Changes

Follow these principles:

- **Keep it focused**: One feature or fix per PR
- **Follow existing patterns**: Match the code style of existing files
- **Add tests**: New functionality should include tests
- **Document intent**: Use clear comments for non-obvious code

### 3. Write Tests

Tests are required for:

- New functions in `engine/` packages
- CLI commands
- Version/update logic

Example test pattern:

```go
func TestMyFunction(t *testing.T) {
    result := MyFunction(input)

    if result != expected {
        t.Errorf("MyFunction(%v) = %v, want %v", input, result, expected)
    }
}
```

Mock HTTP responses for network-dependent tests:

```go
server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Respond with test data
    w.Write([]byte(`{"key": "value"}`))
}))
defer server.Close()
```

### 4. Commit Changes

Use clear, descriptive commit messages:

```
Add feature: implement X functionality

- Include specific implementation details
- Reference any related issues
- Keep messages concise but informative
```

Commits should:

- Compile successfully
- Leave tests passing
- Follow the pattern: "add/implement/fix/refactor X"

Example good commits:

```
- "Add version package with build-time metadata"
- "Implement GitHub Releases API integration"
- "Fix progress reporting in download loop"
- "Refactor update checker error handling"
```

### 5. Push and Create PR

```bash
git push origin feature/my-feature
```

Then create a Pull Request with:

- Clear title describing the change
- Description explaining the "why" and "what"
- Link to any related issues
- Testing notes if applicable

## Code Style

### Naming Conventions

- Use `camelCase` for variables and functions
- Use `PascalCase` for exported types and functions
- Use descriptive names: `currentVersion` not `cv`
- Prefer explicit names over abbreviations

### Structure

- Maximum line length: 100 characters (aim for this, not strict)
- Use clear package names: `version`, `update`, not `util`, `helper`
- Place related functionality in same package
- Use interfaces for abstraction points

### Comments

Comment public functions and types:

```go
// MyFunction does something specific.
// It returns X when Y condition is met.
func MyFunction(input string) string {
    // ...
}
```

Avoid obvious comments:

```go
// BAD: Just restates the code
x = x + 1 // increment x

// GOOD: Explains the intent
x += 1 // Account for header offset
```

### Error Handling

Wrap errors with context:

```go
// BAD
return err

// GOOD
return fmt.Errorf("failed to download binary: %w", err)
```

## Architecture Guidelines

### Engine Layer

- Pure business logic, no CLI/TUI dependencies
- Testable with mocks (no real HTTP calls)
- Reusable by both CLI and TUI

### CLI Layer

- Thin wrappers around engine functionality
- Handle user input validation
- Format output for human consumption
- Use `OutputHelper` for consistent formatting

### TUI Layer

- Presentation only
- Use same engine APIs as CLI
- Handle Bubble Tea events
- Use Lip Gloss for styling

## Testing Guidelines

### Unit Tests

- Test public functions and types
- Include edge cases and error conditions
- Use descriptive test names: `TestFunctionX_WithYCondition`
- Run all tests before committing

### Mocking

- Mock external HTTP calls
- Use `httptest.NewServer` for network tests
- Don't mock internal functions

### Coverage

- Aim for >80% coverage on critical paths
- Cover error cases and edge conditions
- Check with `go test -cover ./...`

## Performance Considerations

- Keep startup fast (important for CLI)
- Cache results when appropriate
- Stream large downloads instead of buffering
- Profile before optimizing

## Security

- Verify checksums for downloads
- Use HTTPS for all external connections
- Validate user input at boundaries
- Don't execute untrusted code
- Never log sensitive information

## Documentation

### Code Documentation

- Document exported functions and types
- Explain the "why" not just the "what"
- Include examples for complex functions

### README Updates

Update README.md if your changes:

- Add new commands
- Change existing command behavior
- Affect build process
- Impact users

### Changelog

Maintain CHANGELOG.md with:

- New features
- Breaking changes
- Bug fixes
- Security updates

## Submitting Changes

### Before Submitting

- [ ] Run `make test` - all tests pass
- [ ] Run `make build` - builds successfully
- [ ] Code follows existing style
- [ ] Changes are focused and minimal
- [ ] Commit messages are clear
- [ ] README/docs updated if needed

### PR Review Process

1. Code review for correctness and style
2. Testing verification
3. Performance impact assessment
4. Security review if applicable
5. Documentation review
6. Merge to main branch

### After Merge

- Celebrate! 🎉
- Your code will be included in the next release
- Monitor for any issues in the wild

## Common Tasks

### Adding a New Command

1. Create `cli/mycmd.go`:
   ```go
   func NewMyCmd() *cobra.Command {
       return &cobra.Command{
           Use: "mycommand",
           Short: "Description",
           RunE: func(cmd *cobra.Command, args []string) error {
               // Implementation
               return nil
           },
       }
   }
   ```

2. Add to root command in `cli/root.go`:
   ```go
   rootCmd.AddCommand(NewMyCmd())
   ```

3. Test with `./core mycommand`

### Adding Engine Functionality

1. Create package in `engine/`: `engine/myfeature/`
2. Implement logic in `myfeature.go`
3. Add tests in `myfeature_test.go`
4. Wire CLI commands to use the engine package

### Adding to TUI

1. Create component or update `tui/app.go`
2. Use Bubble Tea patterns for state and messages
3. Apply styles from `tui/styles.go`
4. Test interaction in TUI

## Troubleshooting

### Tests Failing

```bash
# Run verbose to see details
go test -v ./path/to/package

# Run specific test
go test -run TestName -v ./path
```

### Build Issues

```bash
# Clean build
make clean && make build

# Check Go version
go version

# Update dependencies
go mod tidy
go mod download
```

### Import Errors

```bash
# Update Go path if needed
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

## Questions?

- Open an issue on GitHub
- Check existing issues for similar questions
- Look at commit history for examples

## Code of Conduct

- Be respectful and inclusive
- Assume good intent
- Provide constructive feedback
- Help newer contributors

Thank you for contributing to CORE CLI! 🚀
