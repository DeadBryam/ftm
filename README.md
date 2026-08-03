<h1 align="center">ftm</h1>

<p align="center">
  <strong>Foundry Tunnel Manager</strong><br />
  Share your Foundry VTT world without port forwarding.
</p>

<p align="center">
  One binary. Several tunnel providers. Three ways in.<br />
  TUI, web dashboard, or native desktop — same tunnels, same config.
</p>

<p align="center">
  <a href="#install">Install</a>
  ·
  <a href="#sixty-seconds">Quick start</a>
  ·
  <a href="#interfaces">Interfaces</a>
  ·
  <a href="#providers">Providers</a>
  ·
  <a href="https://github.com/sthbryan/ftm/releases/latest">Releases</a>
</p>

---

You want players in your Foundry world. Your router, ISP, or CGNAT has other plans. Port forwarding is a pain; random free tunnels mean memorizing flags and babysitting processes.

**ftm starts and stops tunnels for you.** Pick a provider, point it at Foundry’s port, copy the public URL. When a tunnel dies or is about to expire, you get a notification — not a silent failure mid-session.

**It is a single static binary.** Providers are downloaded on demand. Config lives under your user directory. The web UI is embedded, so the same binary drives the TUI, the dashboard, and the desktop shell.

---

## Screenshots

| ![TUI](docs/tui.webp) | ![Web](docs/web.webp) | ![Desktop](docs/desktop.webp) |
| :--: | :--: | :--: |
| **TUI** — start/stop, logs, copy URL | **Web** — dashboard at `:40500` | **Desktop** — same UI in a native window |

---

## Interfaces

**TUI** — `ftm`  
Full-screen terminal UI. Navigate connections, open logs, copy the public URL, jump to the dashboard.

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate |
| `s` | Start / stop |
| `l` | Logs |
| `c` | Copy URL |
| `w` | Open web dashboard |
| `o` | Open config |
| `a` / `d` | Add / delete |
| `q` | Quit |

**Web** — always on when ftm runs  
Dashboard at `http://localhost:40500` (port auto-detects if busy). Themes, live status, manage tunnels from the browser.

```bash
ftm --web          # dashboard + open browser
ftm --server       # dashboard only
ftm --port 8080    # force a port
```

**Desktop** — release assets `ftm-desktop-*`  
Wails shell around the same embedded UI. Prefer this if you want a window without keeping a terminal open.

---

## Install

**macOS / Linux**

```bash
curl -fsSL https://raw.githubusercontent.com/sthbryan/ftm/main/install.sh | bash
```

Detects OS and arch, installs the CLI to `~/.local/bin` (override with `INSTALL_DIR`), and clears macOS quarantine when needed.

| | |
|---|---|
| **Install script** | The one-liner above — macOS and Linux |
| **Download** | [Releases](https://github.com/sthbryan/ftm/releases/latest) — CLI + desktop apps |
| **From source** | `make install` — needs Go 1.22+ and Bun |

Building from source installs to `$(go env GOPATH)/bin`. Point it somewhere else with `BINDIR`:

```bash
make install BINDIR=/usr/local/bin
```

### CLI assets

| Platform | File |
|----------|------|
| Windows | `ftm-windows-x64.exe` |
| Linux x64 | `ftm-linux-x64` |
| Linux ARM64 | `ftm-linux-arm64` |
| macOS Intel | `ftm-macos-x64` |
| macOS Apple Silicon | `ftm-macos-arm64` |

```bash
chmod +x ftm-*
# e.g. ~/.local/bin or /usr/local/bin
mv ftm-macos-arm64 ~/.local/bin/ftm
```

On macOS, if Gatekeeper blocks the binary:

```bash
xattr -d com.apple.quarantine ~/.local/bin/ftm
```

### Desktop apps

| Platform | File |
|----------|------|
| Windows | `ftm-desktop-windows.exe` |
| Linux | `ftm-desktop-linux` |
| macOS | `ftm-desktop-macos.app.zip` |

```bash
chmod +x ftm-desktop-linux
# macOS: unzip, then if needed:
xattr -d com.apple.quarantine ftm-desktop-macos.app
```

### Flags

| Flag | Effect |
|------|--------|
| *(none)* | TUI + web server |
| `--web` | Web only, open browser |
| `--server` | Web only, no browser |
| `--port N` | Fix the web port |
| `--version` | Print version |
| `--update` | Update to the latest release |
| `--check` | Check for updates only |
| `--uninstall` | Remove the install |

---

## Sixty seconds

```bash
ftm                  # TUI; dashboard at http://localhost:40500
# In the TUI: a → add a tunnel → pick a provider → Foundry port (usually 30000)
# s → start · c → copy public URL · share with players
```

Web-only session:

```bash
ftm --web
```

---

## Providers

| Provider | Notes |
|----------|--------|
| **Cloudflared** | Cloudflare Tunnel — solid default; auto-install |
| **Pinggy** | Quick public URLs |
| **Tunnelmole** | Lightweight; auto-install where available |
| **Bore** | Self-hostable / alternative relay |
| **localhost.run** | SSH-based, no extra binary |
| **Serveo** | SSH-based, no extra binary |

Missing binaries are downloaded when you start a tunnel that needs them. SSH providers only need a working `ssh` client.

---

## Development

Go 1.22+ and Bun. Desktop builds also need [Wails](https://wails.io).

```bash
make build-full   # web assets + bin/ftm with the UI embedded
make build        # Go only (uses existing embedded static files)
make web          # Svelte build → internal/web/static (+ desktop dist)
make test         # go test ./...
make test-race    # go test -race ./...
make vet          # go vet ./...
make fmt          # gofmt -w, then verify nothing is left unformatted
make install      # web + binary → $BINDIR (default $(go env GOPATH)/bin)
make desktop      # Wails app (no package)
make help         # full target list
```

Hot reload against a running ftm web server:

```bash
cd web-svelte && bun run dev
```

**macOS TUI notifications:** install `alerter` (`brew install alerter`) for native alerts. Desktop/GUI paths use the platform APIs.

---

## Design notes

- **One process, three surfaces** — TUI, embedded web, and Wails share the same tunnel manager and config.
- **Providers as plugins** — each tunnel backend implements a small interface; installers download only what you use.
- **Live status** — the dashboard tracks process state over WebSockets so start/stop/errors show up without refresh.
- **Local by default** — the web UI is meant for localhost; treat exposing it beyond your machine as out of scope.

---

<p align="center">
  <sub>MIT License · © 2026 Bryan Villafuerte</sub>
</p>
