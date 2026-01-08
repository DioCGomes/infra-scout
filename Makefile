.PHONY: build run test clean install lint fmt

# Build variables
BINARY_NAME=infra-scout
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR=bin
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GORUN=$(GOCMD) run
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod
GOFMT=gofmt

# Default target
all: build

# Build the binary
build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/infra-scout

# Build for multiple platforms
build-all: build-linux build-darwin build-windows

build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/infra-scout
	GOOS=linux GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/infra-scout

build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/infra-scout
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/infra-scout

build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/infra-scout

# Run the application
run:
	$(GORUN) ./cmd/infra-scout $(ARGS)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

# Install to GOPATH/bin
install:
	$(GOCMD) install $(LDFLAGS) ./cmd/infra-scout

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Format code
fmt:
	$(GOFMT) -s -w .

# Lint code (requires golangci-lint)
lint:
	golangci-lint run ./...

# Check for vulnerabilities (requires govulncheck)
vulncheck:
	govulncheck ./...

# Development: build and run
dev: build
	./$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

# Help
help:
	@echo "Available targets:"
	@echo "  build        - Build the binary"
	@echo "  build-all    - Build for all platforms"
	@echo "  run          - Run the application (use ARGS= to pass arguments)"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  clean        - Remove build artifacts"
	@echo "  install      - Install to GOPATH/bin"
	@echo "  deps         - Download and tidy dependencies"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  vulncheck    - Check for vulnerabilities"
	@echo "  help         - Show this help"
