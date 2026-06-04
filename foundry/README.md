# ftm-foundry

Foundry VTT module that exposes the local server to the internet via [ftm](https://github.com/sthbryan/ftm).

## Install (user-facing)

1. In Foundry VTT: **Add-on Modules → Install Module**
2. Paste this manifest URL:
   ```
   https://github.com/sthbryan/ftm/releases/latest/download/module.json
   ```
3. Enable the module in your world

On first run the module:
- Downloads the latest `ftm` binary for your OS from GitHub
- Saves it to `<foundry-data>/ftm/ftm{,.exe}`
- Spawns `ftm --web --server` and waits for `:40500` to respond
- Registers a sidebar button that opens the tunnel dashboard

## Architecture

```
foundry/                    (sibling of desktop/ and web-svelte/)
├── module.json
├── scripts/
│   ├── module.js          hooks, settings, sidebar button
│   ├── ftm-manager.js     pure Node: isInstalled / install / start / stop / api
│   └── ftm-app.js         ApplicationV2 wrapping an iframe to :40500
├── templates/
│   └── app.hbs            single iframe element
├── styles/
│   └── module.css         minimal frame styling
└── lang/
    └── en.json
```

## Test on Mac (no Foundry required)

The manager is pure Node — no Foundry globals, no DOM. Run a quick smoke test:

```bash
FOUNDRY_DATA_PATH=/tmp/ftm-test node foundry/test.mjs
```

This mocks `game.userData.path`, downloads the real binary, spawns it, hits the API, and cleans up.

## Test in Foundry (Windows recommended for user)

1. Open Foundry VTT in demo mode (no license key → watermark, but modules work)
2. Install the module with the manifest URL above
3. On world ready, watch the console for `ftm-ready`
4. Click the new sidebar button "Tunnel" to open the dashboard

## Settings

- `ftm.autoStart` (default `true`) — spawn `ftm` on Foundry ready
- `ftm.port` (default `40500`) — port ftm should bind to
