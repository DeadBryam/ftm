# ftm-foundry

Foundry VTT v14 module that exposes the local server to the internet via [ftm](https://github.com/sthbryan/ftm).

## Install

1. In Foundry VTT: **Add-on Modules → Install Module**
2. Paste this manifest URL:
   ```
   https://github.com/sthbryan/ftm/releases/latest/download/module.json
   ```
3. Enable the module in your world

On **world ready** the module:
- Downloads the latest `ftm` binary for your OS into `<foundry-data>/ftm/`
- Spawns `ftm --web --server` and waits for `:40500` to respond
- Adds a sidebar control that opens the FTM Tunnel dashboard

## Architecture

```
foundry/
├── module.json
├── scripts/
│   ├── module.js         Foundry v14 hooks + ApplicationV2 dashboard
│   └── ftm-manager.js    pure Node: install / start / stop / status
├── templates/
│   └── app.hbs           dashboard template (iframe to :40500)
├── styles/
│   └── module.css        dashboard styling
└── lang/
    └── en.json
```

## Test on Mac (no Foundry required)

The manager is pure Node — no Foundry globals, no DOM. Smoke test:

```bash
FOUNDRY_DATA_PATH=/tmp/ftm-test node foundry/test.mjs
```

Mocks the data path, downloads the real binary, spawns it, hits the API, and cleans up.

## Test in Foundry v14 (portable)

1. Open Foundry VTT portable / demo mode
2. Install with the manifest URL above
3. On world ready, watch the console (F12) for `[ftm-tunnel]` logs
4. Click the sidebar **network icon** to open the dashboard

## Settings

(none exposed yet — module starts/stops ftm automatically on ready/close)
