# Binary output name
BINARY_NAME=andtls
MAIN_PACKAGE=./cmd/andtls
BUILD_DIR=bin
VERSION=0.1.0
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS=-s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildDate=$(BUILD_DATE)

.PHONY: all build run clean install uninstall fmt vet deps help

# Default target builds the binary
all: build

# Build binary into root and bin directory
build:
	@echo "Building $(BINARY_NAME) v$(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BINARY_NAME) $(MAIN_PACKAGE)
	@cp $(BINARY_NAME) $(BUILD_DIR)/$(BINARY_NAME)
	@echo "Build successful: ./$(BINARY_NAME)"

# Run application directly
run:
	go run $(MAIN_PACKAGE)

# Format Go source code
fmt:
	@echo "Formatting Go source files..."
	go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	go mod download

# Install binary to ~/.local/bin or /usr/local/bin
install: build
	@echo "Installing $(BINARY_NAME)..."
	@if [ -w /usr/local/bin ]; then \
		cp $(BINARY_NAME) /usr/local/bin/$(BINARY_NAME); \
		echo "Installed to /usr/local/bin/$(BINARY_NAME)"; \
	else \
		mkdir -p $(HOME)/.local/bin; \
		cp $(BINARY_NAME) $(HOME)/.local/bin/$(BINARY_NAME); \
		echo "Installed to $(HOME)/.local/bin/$(BINARY_NAME)"; \
		echo "Ensure $(HOME)/.local/bin is in your PATH"; \
	fi

# Uninstall binary
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f /usr/local/bin/$(BINARY_NAME)
	@rm -f $(HOME)/.local/bin/$(BINARY_NAME)
	@echo "Uninstalled $(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete"

# Show available targets
help:
	@echo "andtls Makefile targets:"
	@echo "  make build      - Build the binary"
	@echo "  make run        - Run application directly"
	@echo "  make fmt        - Format Go source code"
	@echo "  make vet        - Run go vet analysis"
	@echo "  make install    - Install binary to PATH"
	@echo "  make uninstall  - Remove binary from PATH"
	@echo "  make clean      - Remove build files"
	@echo "  make deps       - Download Go module dependencies"

