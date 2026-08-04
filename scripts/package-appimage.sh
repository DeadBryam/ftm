#!/usr/bin/env bash
# Wrap the Wails Linux binary into an AppImage.
#
# One file that runs on any distro with glibc >= the build host's, with the
# icon and .desktop entry inside it, so launchers stop showing a generic icon.
# GTK and WebKitGTK are bundled by linuxdeploy, which is why this must run on
# the oldest distro we intend to support (the workflow pins ubuntu-22.04).
#
# Output: <out-dir>/ftm-desktop-linux-x86_64.AppImage
#
# Usage: package-appimage.sh <binary> [out-dir]
set -euo pipefail

BIN="${1:?usage: package-appimage.sh <binary> [out-dir]}"
OUT_DIR="${2:-$(dirname "$BIN")}"
VERSION="${VERSION:-$(sed -nE 's/^var Version = "([^"]+)"$/\1/p' internal/version/version.go)}"

export ARCH="${ARCH:-x86_64}"
# GitHub runners have no FUSE, so the AppImage tools must self-extract instead.
export APPIMAGE_EXTRACT_AND_RUN=1

BIN="$(cd "$(dirname "$BIN")" && pwd)/$(basename "$BIN")"
mkdir -p "$OUT_DIR"
OUT_DIR="$(cd "$OUT_DIR" && pwd)"

WORK="$(mktemp -d)"
trap 'rm -rf "$WORK"' EXIT

APPDIR="$WORK/AppDir"
mkdir -p "$APPDIR/usr/bin" "$APPDIR/usr/share/applications"

cp "$BIN" "$APPDIR/usr/bin/ftm-desktop"
chmod +x "$APPDIR/usr/bin/ftm-desktop"

cat > "$APPDIR/usr/share/applications/ftm-desktop.desktop" <<DESKTOP
[Desktop Entry]
Type=Application
Name=Foundry Tunnel Manager
GenericName=Tunnel Manager
Comment=Share your Foundry VTT world without port forwarding
Exec=ftm-desktop
Icon=ftm-desktop
Categories=Network;Game;
Keywords=foundry;vtt;tunnel;
Terminal=false
StartupWMClass=ftm-desktop
DESKTOP

for size in 64 128 256 512; do
  dir="$APPDIR/usr/share/icons/hicolor/${size}x${size}/apps"
  mkdir -p "$dir"
  if command -v convert >/dev/null 2>&1; then
    convert desktop/icon.png -resize "${size}x${size}" "$dir/ftm-desktop.png"
  else
    cp desktop/icon.png "$dir/ftm-desktop.png"
  fi
done

TOOLS="$WORK/tools"
mkdir -p "$TOOLS"
base="https://github.com/linuxdeploy/linuxdeploy/releases/download/continuous"
gtk="https://raw.githubusercontent.com/linuxdeploy/linuxdeploy-plugin-gtk/master/linuxdeploy-plugin-gtk.sh"

curl -fsSL -o "$TOOLS/linuxdeploy" "$base/linuxdeploy-${ARCH}.AppImage"
curl -fsSL -o "$TOOLS/linuxdeploy-plugin-gtk.sh" "$gtk"
chmod +x "$TOOLS/linuxdeploy" "$TOOLS/linuxdeploy-plugin-gtk.sh"
export PATH="$TOOLS:$PATH"

"$TOOLS/linuxdeploy" \
  --appdir "$APPDIR" \
  --executable "$APPDIR/usr/bin/ftm-desktop" \
  --desktop-file "$APPDIR/usr/share/applications/ftm-desktop.desktop" \
  --icon-file "$APPDIR/usr/share/icons/hicolor/256x256/apps/ftm-desktop.png" \
  --plugin gtk

# WebKitGTK launches its network and web processes from a path baked in at
# compile time, which does not exist inside the AppImage. Without this hook the
# window comes up permanently blank.
mkdir -p "$APPDIR/apprun-hooks"
cat > "$APPDIR/apprun-hooks/webkit.sh" <<'HOOK'
for libexec in "$APPDIR/usr/libexec/webkit2gtk-4.1" \
               "$APPDIR/usr/lib/x86_64-linux-gnu/webkit2gtk-4.1" \
               "$APPDIR/usr/lib/webkit2gtk-4.1"; do
  if [ -d "$libexec" ]; then
    export WEBKIT_EXEC_PATH="$libexec"
    break
  fi
done
export WEBKIT_DISABLE_COMPOSITING_MODE=1
HOOK

OUTPUT_NAME="ftm-desktop-linux-${ARCH}.AppImage"
(
  cd "$WORK"
  VERSION="$VERSION" "$TOOLS/linuxdeploy" --appdir "$APPDIR" --output appimage
)

built="$(find "$WORK" -maxdepth 1 -name '*.AppImage' -print -quit)"
[ -n "$built" ] || { echo "linuxdeploy produced no AppImage" >&2; exit 1; }

mv "$built" "$OUT_DIR/$OUTPUT_NAME"
chmod +x "$OUT_DIR/$OUTPUT_NAME"
rm -f "$BIN"

echo "Created $OUT_DIR/$OUTPUT_NAME"
