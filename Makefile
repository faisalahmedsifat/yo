# yo CLI

.PHONY: build install test clean run

# Build variables
VERSION := 0.1.0
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
LDFLAGS := -ldflags "-X github.com/faisalahmedsifat/yo/cmd.Version=$(VERSION) -X github.com/faisalahmedsifat/yo/cmd.BuildDate=$(BUILD_DATE) -X github.com/faisalahmedsifat/yo/cmd.GitCommit=$(GIT_COMMIT)"

# Default target
all: build

# Build the binary
build:
	go build $(LDFLAGS) -o yo .

# Install to $GOPATH/bin
install:
	go install $(LDFLAGS) .

# Run tests
test:
	go test -v ./...

# Run unit tests only
test-unit:
	go test -v ./internal/...

# Run integration tests
test-integration:
	go test -v ./tests/...

# Clean build artifacts
clean:
	rm -f yo
	go clean

# Run the CLI
run:
	go run . $(ARGS)

# Format code
fmt:
	go fmt ./...

# Vet code
vet:
	go vet ./...

# Lint (requires golangci-lint)
lint:
	golangci-lint run

# Tidy dependencies
tidy:
	go mod tidy

# Cross-compile for multiple platforms
build-all:
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/yo-linux-amd64 .
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/yo-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/yo-darwin-arm64 .

# Help
help:
	@echo "Available targets:"
	@echo "  build        - Build the yo binary"
	@echo "  install      - Install to GOPATH/bin"
	@echo "  test         - Run all tests"
	@echo "  test-unit    - Run unit tests only"
	@echo "  test-integration - Run integration tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  fmt          - Format code"
	@echo "  vet          - Vet code"
	@echo "  lint         - Run linter"
	@echo "  tidy         - Tidy dependencies"
	@echo "  build-all    - Cross-compile for all platforms"
