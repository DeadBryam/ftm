#!/usr/bin/env bash
# Wrap a macOS binary into a minimal .app and zip it for GitHub Releases.
set -euo pipefail

BIN="${1:?usage: package-macos-app.sh <binary> [out-dir]}"
OUT_DIR="${2:-$(dirname "$BIN")}"
VERSION="${VERSION:-0.10.0}"
APP_NAME="ftm-desktop-macos.app"
APP="$OUT_DIR/$APP_NAME"

rm -rf "$APP"
mkdir -p "$APP/Contents/MacOS" "$APP/Contents/Resources"
cp "$BIN" "$APP/Contents/MacOS/ftm-desktop"
chmod +x "$APP/Contents/MacOS/ftm-desktop"

if [ -f desktop/icon.png ]; then
  cp desktop/icon.png "$APP/Contents/Resources/icon.png"
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
  <string>Foundry Tunnel Manager</string>
  <key>CFBundlePackageType</key>
  <string>APPL</string>
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

(cd "$OUT_DIR" && zip -qr "${APP_NAME}.zip" "$APP_NAME")
rm -rf "$APP"
rm -f "$BIN"
echo "Created $OUT_DIR/${APP_NAME}.zip"
