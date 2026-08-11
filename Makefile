# FTM Makefile

# Variables
BINARY_NAME=ftm
# git describe includes the leading "v" from tags; strip it so callers can print "v$(Version)".
VERSION=$(shell (git describe --tags --always --dirty 2>/dev/null || echo dev) | sed 's/^v//')
LDFLAGS=-ldflags "-X github.com/sthbryan/ftm/internal/version.Version=$(VERSION)"
CGO_ENABLED=0
# Windows resources need a strictly numeric version, so drop any git describe suffix.
WINRES_VERSION=$(shell echo $(VERSION) | sed 's/[-+].*//')
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

# Desktop (Wails v3) — needs CGO + platform webview libs. Builds for the host OS.
DESKTOP_OUT := $(DESKTOP_DIR)/build/bin
DESKTOP_BIN := ftm-desktop
ifeq ($(OS),Windows_NT)
DESKTOP_BIN := ftm-desktop.exe
endif
# Wails v3 defaults to webkitgtk-6.0 (GTK4) on Linux, but the bundled AppImage
# helpers target WebKitGTK 4.1 (GTK3), so force the gtk3 tag on Linux hosts.
UNAME_S := $(shell uname -s)
DESKTOP_DEV_TAGS_ARG :=
DESKTOP_PROD_TAGS_ARG := production
ifeq ($(UNAME_S),Linux)
DESKTOP_DEV_TAGS_ARG := gtk3
DESKTOP_PROD_TAGS_ARG := production,gtk3
endif

.PHONY: desktop-dev
desktop-dev: web
	@echo "Building desktop (dev)..."
	@mkdir -p $(DESKTOP_OUT)
	CGO_ENABLED=1 go build $(if $(DESKTOP_DEV_TAGS_ARG),-tags $(DESKTOP_DEV_TAGS_ARG),) -o $(DESKTOP_OUT)/$(DESKTOP_BIN) ./desktop

.PHONY: desktop
desktop: web
	@echo "Building desktop (production tags)..."
	@mkdir -p $(DESKTOP_OUT)
	CGO_ENABLED=1 go build $(if $(DESKTOP_PROD_TAGS_ARG),-tags $(DESKTOP_PROD_TAGS_ARG),) -ldflags "-s -w" -o $(DESKTOP_OUT)/$(DESKTOP_BIN) ./desktop

.PHONY: desktop-winres
desktop-winres:
	@echo "Generating Windows resources (icon, manifest, version info)..."
	go run github.com/tc-hib/go-winres@v0.3.3 make \
		--in $(DESKTOP_DIR)/winres/winres.json \
		--out $(DESKTOP_DIR)/rsrc \
		--arch amd64,arm64 \
		--product-version $(VERSION) \
		--file-version $(WINRES_VERSION)

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

# Release: bump version, commit, tag, push.
# Forward all args to scripts/bump-version.sh, e.g.
#   make version                 # interactive menu
#   make version ARGS="patch"    # bump patch
#   make version ARGS="0.11.0 --dry-run"
.PHONY: version
version:
	@./scripts/bump-version.sh $(ARGS)

.PHONY: release
release: version

# Help target
.PHONY: help
help:
	@echo "FTM Makefile targets:"
	@echo "  build                     - Build CLI for current platform → $(BUILD_DIR)/$(BINARY_NAME)"
	@echo "  web                       - Build Svelte UI into internal/web/static (+ desktop dist)"
	@echo "  build-full                - Build web UI then the Go binary"
	@echo "  desktop                   - Build desktop app (Wails v3, production, host OS)"
	@echo "  desktop-dev               - Build desktop without production tags"
	@echo "  run / dev                 - go run the CLI"
	@echo "  test                      - Run go tests"
	@echo "  test-race                 - Run tests with -race detector"
	@echo "  vet                       - Run go vet"
	@echo "  fmt                       - Run gofmt -w and verify formatting"
	@echo "  clean                     - Remove build artifacts"
	@echo "  install                   - Build web + binary, copy to \$$BINDIR (BINDIR=... override; default \$$(go env GOPATH)/bin)"
	@echo "  uninstall                 - Remove installed binary from \$$BINDIR"
	@echo "  version [ARGS=...]        - Bump version (interactive). ARGS forwarded to scripts/bump-version.sh"
	@echo "  release                   - Alias of version"
	@echo "  help                      - Show this help message"
