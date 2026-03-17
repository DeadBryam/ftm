#!/bin/bash

if [ -z "$1" ]; then
    echo "Usage: ./scripts/bump-version.sh <version>"
    echo "Example: ./scripts/bump-version.sh 1.2.3"
    exit 1
fi

VERSION=$1
VERSION_NO_V=$(echo $VERSION | sed 's/^v//')

echo "Bumping version to $VERSION_NO_V..."

# Update desktop/package.json
if [ -f "desktop/package.json" ]; then
    sed -i.bak "s/\"version\": \"[^\"]*\"/\"version\": \"$VERSION_NO_V\"/" desktop/package.json
    rm -f desktop/package.json.bak
    echo "✓ Updated desktop/package.json"
fi

# Update Cargo.toml
if [ -f "desktop/src-tauri/Cargo.toml" ]; then
    sed -i.bak "s/^version = \"[^\"]*\"/version = \"$VERSION_NO_V\"/" desktop/src-tauri/Cargo.toml
    rm -f desktop/src-tauri/Cargo.toml.bak
    echo "✓ Updated desktop/src-tauri/Cargo.toml"
fi

# Update tauri.conf.json
if [ -f "desktop/src-tauri/tauri.conf.json" ]; then
    sed -i.bak "s/\"version\": \"[^\"]*\"/\"version\": \"$VERSION_NO_V\"/" desktop/src-tauri/tauri.conf.json
    rm -f desktop/src-tauri/tauri.conf.json.bak
    echo "✓ Updated desktop/src-tauri/tauri.conf.json"
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
