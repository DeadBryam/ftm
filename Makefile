# FTM Makefile

# Variables
BINARY_NAME=ftm
# git describe includes the leading "v" from tags; strip it so callers can print "v$(Version)".
VERSION=$(shell (git describe --tags --always --dirty 2>/dev/null || echo dev) | sed 's/^v//')
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILT=$(shell date +%Y-%m-%d)
LDFLAGS=-ldflags "-X github.com/sthbryan/ftm/internal/version.Version=$(VERSION)"
CGO_ENABLED=0
BUILD_DIR=bin
PKG=./cmd/ftm
DESKTOP_DIR=./desktop

# Install location (override with `make install BINDIR=/path/to/bin`)
GOBIN := $(shell go env GOPATH)/bin
BINDIR ?= $(GOBIN)
DESTDIR ?=
INSTALL_PATH=$(DESTDIR)$(BINDIR)/$(BINARY_NAME)

# Default target
.DEFAULT_GOAL:=build

# Frontend (Svelte) → internal/web/static + desktop/frontend/dist
.PHONY: web
web:
	@echo "Building web UI..."
	./scripts/build-web-assets.sh
	mkdir -p $(DESKTOP_DIR)/build
	cp $(DESKTOP_DIR)/icon.png $(DESKTOP_DIR)/build/appicon.png

# Build for current platform
.PHONY: build
build:
	@echo "Building $(BINARY_NAME) for current platform... ($(VERSION))"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(PKG)

.PHONY: build-full
build-full: web build

# Build for all platforms
.PHONY: build-all
build-all: build-darwin-amd64 build-darwin-arm64 build-linux-amd64 build-linux-arm64 build-windows-amd64

.PHONY: build-darwin-amd64
build-darwin-amd64:
	@echo "Building $(BINARY_NAME) for darwin/amd64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(PKG)

.PHONY: build-darwin-arm64
build-darwin-arm64:
	@echo "Building $(BINARY_NAME) for darwin/arm64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=$(CGO_ENABLED) go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(PKG)

.PHONY: build-linux-amd64
build-linux-amd64:
	@echo "Building $(BINARY_NAME) for linux/amd64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(PKG)

.PHONY: build-linux-arm64
build-linux-arm64:
	@echo "Building $(BINARY_NAME) for linux/arm64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=$(CGO_ENABLED) go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(PKG)

.PHONY: build-windows-amd64
build-windows-amd64:
	@echo "Building $(BINARY_NAME) for windows/amd64..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 CGO_ENABLED=$(CGO_ENABLED) go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(PKG)

# Desktop (Wails)
.PHONY: desktop-dev
desktop-dev: web
	@echo "Building desktop (dev, no package)..."
	cd $(DESKTOP_DIR) && wails build -s -nopackage -devtools

.PHONY: desktop
desktop: web
	@echo "Building desktop..."
	cd $(DESKTOP_DIR) && wails build -s -nopackage

.PHONY: desktop-package
desktop-package: web
	@echo "Building desktop package..."
	cd $(DESKTOP_DIR) && wails build -s

.PHONY: desktop-darwin-universal
desktop-darwin-universal: web
	@echo "Building desktop for darwin/universal..."
	cd $(DESKTOP_DIR) && wails build -s -nopackage -platform darwin/universal

.PHONY: desktop-darwin-arm64
desktop-darwin-arm64: web
	@echo "Building desktop for darwin/arm64..."
	cd $(DESKTOP_DIR) && wails build -s -nopackage -platform darwin/arm64

.PHONY: desktop-darwin-amd64
desktop-darwin-amd64: web
	@echo "Building desktop for darwin/amd64..."
	cd $(DESKTOP_DIR) && wails build -s -nopackage -platform darwin/amd64

.PHONY: desktop-linux-amd64
desktop-linux-amd64: web
	@echo "Building desktop for linux/amd64..."
	cd $(DESKTOP_DIR) && wails build -s -nopackage -platform linux/amd64

.PHONY: desktop-windows-amd64
desktop-windows-amd64: web
	@echo "Building desktop for windows/amd64..."
	cd $(DESKTOP_DIR) && wails build -s -nopackage -platform windows/amd64

.PHONY: desktop-all
desktop-all: web
	@echo "Building desktop for all platforms..."
	cd $(DESKTOP_DIR) && wails build -s -nopackage -platform darwin/universal
	cd $(DESKTOP_DIR) && wails build -s -nopackage -platform linux/amd64
	cd $(DESKTOP_DIR) && wails build -s -nopackage -platform windows/amd64

# Run without installing
.PHONY: run
run:
	go run $(PKG)

.PHONY: dev
dev: run

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

# Run tests with race detector
.PHONY: test-race
test-race:
	@echo "Running tests with -race..."
	go test -race ./...

# Run go vet
.PHONY: vet
vet:
	@echo "Running go vet..."
	go vet ./...

# Format Go source files with gofmt (CI checks this with `gofmt -l .`)
.PHONY: fmt
fmt:
	@echo "Running gofmt -w..."
	gofmt -w .
	@echo "Checking for unformatted files..."
	@test -z "$$(gofmt -l .)" || (echo "ERROR: gofmt found unformatted files:" && gofmt -l . && exit 1)
	@echo "All files formatted."

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f $(BINARY_NAME) ftm-*
	rm -rf $(DESKTOP_DIR)/build/bin/*
	rm -rf $(DESKTOP_DIR)/frontend/dist
	@go clean -cache -testcache 2>/dev/null || true

# Install binary (builds web UI + binary, then copies it to $(BINDIR)).
.PHONY: install
install: web build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	install -d $(DESTDIR)$(BINDIR)
	install -m 0755 $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_PATH)
	@echo "Installed: $(INSTALL_PATH)"

# Uninstall binary
.PHONY: uninstall
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@rm -f $(INSTALL_PATH)

# Help target
.PHONY: help
help:
	@echo "FTM Makefile targets:"
	@echo "  build                     - Build CLI for current platform → $(BUILD_DIR)/$(BINARY_NAME)"
	@echo "  web                       - Build Svelte UI into internal/web/static (+ desktop dist)"
	@echo "  build-full                - Build web UI then the Go binary"
	@echo "  build-all                 - Build CLI for all platforms (darwin/amd64, darwin/arm64, linux/amd64, linux/arm64, windows/amd64)"
	@echo "  desktop                   - Build desktop app (Wails, no package)"
	@echo "  desktop-dev               - Build desktop with devtools"
	@echo "  desktop-package           - Build packaged desktop app"
	@echo "  desktop-all               - Build desktop for darwin/universal, linux/amd64, windows/amd64"
	@echo "  desktop-darwin-universal  - Desktop darwin/universal"
	@echo "  desktop-darwin-arm64      - Desktop darwin/arm64"
	@echo "  desktop-darwin-amd64      - Desktop darwin/amd64"
	@echo "  desktop-linux-amd64       - Desktop linux/amd64"
	@echo "  desktop-windows-amd64     - Desktop windows/amd64"
	@echo "  run / dev                 - go run the CLI"
	@echo "  test                      - Run go tests"
	@echo "  test-race                 - Run tests with -race detector"
	@echo "  vet                       - Run go vet"
	@echo "  fmt                       - Run gofmt -w and verify formatting"
	@echo "  clean                     - Remove build artifacts"
	@echo "  install                   - Build web + binary, copy to \$$BINDIR (BINDIR=... override; default \$$(go env GOPATH)/bin)"
	@echo "  uninstall                 - Remove installed binary from \$$BINDIR"
	@echo "  help                      - Show this help message"
