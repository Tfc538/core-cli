.PHONY: build build-all test clean checksums help

# Version configuration
VERSION ?= dev
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# Package path
PACKAGE := github.com/Tfc538/core-cli

# LDFLAGS for version injection
LDFLAGS := -X '$(PACKAGE)/version.Version=$(VERSION)' \
           -X '$(PACKAGE)/version.GitCommit=$(GIT_COMMIT)' \
           -X '$(PACKAGE)/version.BuildDate=$(BUILD_DATE)'

# Output directory
DIST_DIR := dist

help:
	@echo "CORE CLI Build Targets:"
	@echo ""
	@echo "  make build              Build for current platform"
	@echo "  make build-all          Build for all supported platforms"
	@echo "  make VERSION=1.0.0      Build with specific version (default: dev)"
	@echo "  make test               Run tests"
	@echo "  make clean              Remove build artifacts"
	@echo "  make checksums          Generate SHA256 checksums"
	@echo ""
	@echo "Supported platforms:"
	@echo "  - linux-amd64"
	@echo "  - linux-arm64"
	@echo "  - darwin-amd64"
	@echo "  - darwin-arm64"
	@echo "  - windows-amd64"

# Build for current platform
build: clean
	@echo "Building CORE CLI v$(VERSION)..."
	go build -ldflags "$(LDFLAGS)" -o core ./main.go
	@echo "✓ Built: ./core"

# Build for all supported platforms
build-all: clean
	@echo "Building CORE CLI v$(VERSION) for all platforms..."
	@mkdir -p $(DIST_DIR)
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/core-linux-amd64
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/core-linux-arm64
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/core-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/core-darwin-arm64
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/core-windows-amd64.exe
	@echo "✓ Built all binaries in ./dist"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...
	@echo "✓ Tests passed"

# Generate checksums for dist binaries
checksums: build-all
	@echo "Generating checksums..."
	@cd $(DIST_DIR) && sha256sum core-* > checksums.txt
	@echo "✓ Checksums generated in ./dist/checksums.txt"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -f core
	rm -rf $(DIST_DIR)
	@echo "✓ Cleaned"
