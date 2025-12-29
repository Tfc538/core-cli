.PHONY: build build-all build-cli build-backend test clean checksums help

# Version configuration
VERSION ?= dev
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

# Package path
PACKAGE := github.com/Tfc538/core-cli

# LDFLAGS for version injection
LDFLAGS := -X '$(PACKAGE)/internal/version.Version=$(VERSION)' \
           -X '$(PACKAGE)/internal/version.GitCommit=$(GIT_COMMIT)' \
           -X '$(PACKAGE)/internal/version.BuildDate=$(BUILD_DATE)'

# Output directory
DIST_DIR := dist

help:
	@echo "CORE CLI Build Targets:"
	@echo ""
	@echo "  make build              Build CLI + backend for current platform"
	@echo "  make build-cli          Build CLI for current platform"
	@echo "  make build-backend      Build backend for current platform"
	@echo "  make build-all          Build all executables for all supported platforms"
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

# Build for current platform (CLI + backend)
build: clean build-cli build-backend

# Build CLI for current platform
build-cli:
	@echo "Building CORE CLI v$(VERSION)..."
	@mkdir -p $(DIST_DIR)/core
	go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/core/core ./cmd/core
	@echo "✓ Built: $(DIST_DIR)/core/core"

# Build for all supported platforms
build-all: clean
	@echo "Building CORE executables v$(VERSION) for all platforms..."
	@mkdir -p $(DIST_DIR)/core $(DIST_DIR)/core-backend
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/core/core-linux-amd64 ./cmd/core
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/core/core-linux-arm64 ./cmd/core
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/core/core-darwin-amd64 ./cmd/core
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/core/core-darwin-arm64 ./cmd/core
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/core/core-windows-amd64.exe ./cmd/core
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/core-backend/core-backend-linux-amd64 ./cmd/core-backend
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/core-backend/core-backend-linux-arm64 ./cmd/core-backend
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/core-backend/core-backend-darwin-amd64 ./cmd/core-backend
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/core-backend/core-backend-darwin-arm64 ./cmd/core-backend
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/core-backend/core-backend-windows-amd64.exe ./cmd/core-backend
	@echo "✓ Built all binaries in ./dist"

# Build backend for current platform
build-backend:
	@echo "Building CORE Backend v$(VERSION)..."
	@mkdir -p $(DIST_DIR)/core-backend
	go build -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/core-backend/core-backend ./cmd/core-backend
	@echo "✓ Built: $(DIST_DIR)/core-backend/core-backend"

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...
	@echo "✓ Tests passed"

# Generate checksums for dist binaries
checksums: build-all
	@echo "Generating checksums..."
	@cd $(DIST_DIR) && find core core-backend -maxdepth 1 -type f -print0 | xargs -0 sha256sum > checksums.txt
	@echo "✓ Checksums generated in ./dist/checksums.txt"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(DIST_DIR)
	@echo "✓ Cleaned"
