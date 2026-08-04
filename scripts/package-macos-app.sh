#!/usr/bin/env bash
# Wrap a macOS binary into a minimal .app and zip it for GitHub Releases.
#
# Output: <out-dir>/ftm-desktop-macos.app.zip containing "Foundry Tunnel
# Manager.app" (the friendly bundle name shown by Finder/Launchpad). The
# outer zip file keeps the asset-friendly slug for URL stability.
set -euo pipefail

BIN="${1:?usage: package-macos-app.sh <binary> [out-dir]}"
OUT_DIR="${2:-$(dirname "$BIN")}"
VERSION="${VERSION:-0.10.0}"

# Outer zip name (URL slug — never change without a release tag bump).
ZIP_SLUG="ftm-desktop-macos.app"
# Inner bundle name (what users see in /Applications).
BUNDLE_NAME="Foundry Tunnel Manager.app"

APP="$OUT_DIR/$BUNDLE_NAME"
ICNS_PATH="$APP/Contents/Resources/AppIcon.icns"

rm -rf "$APP"
mkdir -p "$APP/Contents/MacOS" "$APP/Contents/Resources"

cp "$BIN" "$APP/Contents/MacOS/ftm-desktop"
chmod +x "$APP/Contents/MacOS/ftm-desktop"

# Build AppIcon.icns from desktop/icon.png using sips + iconutil. macOS only.
if [ -f desktop/icon.png ]; then
  ICONSET="$(mktemp -d)/AppIcon.iconset"
  mkdir -p "$ICONSET"
  trap 'rm -rf "$(dirname "$ICONSET")"' EXIT

  # Generate all required sizes from the 1024x1024 source.
  declare -a SIZES=(16 32 64 128 256 512)
  for size in "${SIZES[@]}"; do
    sips -z "$size" "$size" desktop/icon.png --out "$ICONSET/icon_${size}x${size}.png" >/dev/null
    sips -z $((size * 2)) $((size * 2)) desktop/icon.png --out "$ICONSET/icon_${size}x${size}@2x.png" >/dev/null
  done
  # 512@2x = 1024 source, already in place, but iconutil wants it named explicitly.
  cp desktop/icon.png "$ICONSET/icon_512x512@2x.png"

  iconutil -c icns "$ICONSET" -o "$ICNS_PATH"
else
  echo "warning: desktop/icon.png not found; shipping without icon" >&2
fi

cat > "$APP/Contents/Info.plist" <<PLIST
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>CFBundleExecutable</key>
  <string>ftm-desktop</string>
  <key>CFBundleIdentifier</key>
  <string>com.sthbryan.ftm</string>
  <key>CFBundleName</key>
  <string>${BUNDLE_NAME%.app}</string>
  <key>CFBundleDisplayName</key>
  <string>${BUNDLE_NAME%.app}</string>
  <key>CFBundlePackageType</key>
  <string>APPL</string>
  <key>CFBundleIconFile</key>
  <string>AppIcon</string>
  <key>CFBundleShortVersionString</key>
  <string>${VERSION}</string>
  <key>CFBundleVersion</key>
  <string>${VERSION}</string>
  <key>LSMinimumSystemVersion</key>
  <string>10.15</string>
  <key>NSHighResolutionCapable</key>
  <true/>
</dict>
</plist>
PLIST

# Codesign + xattr only matter when running locally; CI also runs this on a
# macOS runner, so apply the same hardening here as the cask does post-install.
if [ "$(uname -s)" = "Darwin" ]; then
  codesign --force --deep --sign - "$APP" || true
  xattr -cr "$APP" || true
fi

(cd "$OUT_DIR" && zip -qr "${ZIP_SLUG}.zip" "$BUNDLE_NAME")
rm -rf "$APP"
rm -f "$BIN"
echo "Created $OUT_DIR/${ZIP_SLUG}.zip"
