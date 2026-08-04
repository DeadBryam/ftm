#!/usr/bin/env python3
import pathlib
import sys


def main() -> int:
    if len(sys.argv) != 3:
        print("usage: relocate-webkit.py <libwebkit2gtk.so> <WK_LINK>", file=sys.stderr)
        return 2

    lib = pathlib.Path(sys.argv[1])
    new_link = sys.argv[2].encode()
    blob = lib.read_bytes()

    found = set()
    marker = b"webkit2gtk-4.1"
    start = 0
    while True:
        i = blob.find(marker, start)
        if i < 0:
            break
        # Expand to the enclosing NUL-delimited C string.
        beg = blob.rfind(b"\x00", 0, i) + 1
        end = blob.find(b"\x00", i)
        frag = blob[beg:end]
        if frag.startswith(b"/") and (
            frag.endswith(b"/webkit2gtk-4.1")
            or frag.endswith(b"/webkit2gtk-4.1/injected-bundle/")
        ):
            found.add(frag)
        start = i + 1

    if not found:
        print(
            f"relocate-webkit: no absolute '.../webkit2gtk-4.1' helper dir found in {lib} -- "
            "this WebKitGTK lays out its helpers somewhere unexpected. "
            "The AppImage would ship with a broken renderer.",
            file=sys.stderr,
        )
        return 2

    new_link_plain = new_link.decode()
    for full in sorted(found, key=lambda s: -len(s)):
        idx = full.index(marker) + len(marker)
        new_bytes = new_link + full[idx:]
        if len(new_bytes) > len(full):
            print(
                f"relocate-webkit: {new_bytes!r} ({len(new_bytes)} bytes) longer than "
                f"compiled-in {full!r} ({len(full)} bytes); cannot grow in place.",
                file=sys.stderr,
            )
            return 2
        blob = blob.replace(full, new_bytes + b"\0" * (len(full) - len(new_bytes)))

    lib.write_bytes(blob)
    print(
        f"relocate-webkit: redirected {sorted(found)} -> {new_link_plain} in {lib}",
        file=sys.stderr,
    )
    return 0


if __name__ == "__main__":
    sys.exit(main())