#!/bin/bash

if [ -z "$1" ]; then
    echo "Usage: ./scripts/bump-version.sh <version>"
    echo "Example: ./scripts/bump-version.sh 1.2.3"
    exit 1
fi

VERSION=$1
VERSION_NO_V=$(echo $VERSION | sed 's/^v//')

echo "Bumping version to $VERSION_NO_V..."

# Version string for the binary comes from `git describe` at build time
# (see Makefile). Keep a sensible default in version.go for plain `go build`.
if [ -f "internal/version/version.go" ]; then
    sed -i.bak "s/^var Version = .*/var Version = \"$VERSION_NO_V\"/" internal/version/version.go
    rm -f internal/version/version.go.bak
    echo "✓ Updated internal/version/version.go"
fi

# Update web-svelte/package.json
if [ -f "web-svelte/package.json" ]; then
    sed -i.bak "s/\"version\": \"[^\"]*\"/\"version\": \"$VERSION_NO_V\"/" web-svelte/package.json
    rm -f web-svelte/package.json.bak
    echo "✓ Updated web-svelte/package.json"
fi

# Update desktop/wails.json
if [ -f "desktop/wails.json" ]; then
    sed -i.bak "s/\"version\": \"[^\"]*\"/\"version\": \"$VERSION_NO_V\"/" desktop/wails.json
    rm -f desktop/wails.json.bak
    echo "✓ Updated desktop/wails.json"
fi

# Commit changes
git add -A
git commit -m "chore(release): Bump version to $VERSION"

# Create tag
git tag -a "v$VERSION_NO_V" -m "Release v$VERSION_NO_V"

echo ""
echo "Version bumped to $VERSION_NO_V"
echo "To publish the release, run:"
echo "  git push origin main --tags"
