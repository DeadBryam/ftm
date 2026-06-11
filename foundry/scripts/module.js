// scripts/module.js — Foundry v14 module entrypoint
// Auto-installs ftm binary, spawns the web server, registers a sidebar
// control that opens the FTM dashboard as a native Foundry app.

import { FtmManager } from './ftm-manager.js';

const MODULE_ID = 'ftm-tunnel';
const FTM_WEB = 'http://localhost:40500';
const SETTINGS = {
  AUTO_START: 'autoStart',
  PORT: 'port',
};

function log(...args) { console.log(`[${MODULE_ID}]`, ...args); }

let manager = null;

Hooks.once('init', () => {
  game.settings.register(MODULE_ID, SETTINGS.AUTO_START, {
    name: game.i18n.localize('FTM_SETTINGS.AUTO_START'),
    hint: game.i18n.localize('FTM_SETTINGS.AUTO_START_HINT'),
    scope: 'world',
    config: true,
    type: Boolean,
    default: true,
  });
  game.settings.register(MODULE_ID, SETTINGS.PORT, {
    name: game.i18n.localize('FTM_SETTINGS.PORT'),
    scope: 'world',
    config: true,
    type: Number,
    default: 40500,
  });

  game.modules.get(MODULE_ID).api = {
    manager: () => manager,
    status: () => manager?.status(),
  };
});

Hooks.once('ready', async () => {
  // Foundry v14 portable: data path lives in game.userData.path
  const dataPath = game?.userData?.path || globalThis.FOUNDRY_DATA_PATH;

  manager = new FtmManager({ dataPath, port: game.settings.get(MODULE_ID, SETTINGS.PORT) });

  const s = await manager.status();
  log(`platform=${s.platform}/${s.arch} installed=${s.installed} running=${s.running}`);

  const autoStart = game.settings.get(MODULE_ID, SETTINGS.AUTO_START);
  if (!autoStart) {
    log('autoStart disabled, skipping install/start');
    return;
  }

  if (!s.installed) {
    log('installing ftm binary...');
    try {
      const v = await manager.install();
      log(`installed ftm v${v}`);
    } catch (err) {
      log('install failed:', err.message);
      ui.notifications.error(game.i18n.format('FTM_INSTALL_FAILED', { error: err.message }));
      return;
    }
  }

  if (!s.running) {
    log('starting ftm server...');
    try {
      await manager.startServer();
      log('ftm ready on ' + FTM_WEB);
    } catch (err) {
      log('start failed:', err.message);
      ui.notifications.error(game.i18n.format('FTM_START_FAILED', { error: err.message }));
      return;
    }
  }

  ui.notifications.info(game.i18n.localize('FTM_TUNNEL') + ' — ' + FTM_WEB);
});

Hooks.on('getSceneControlButtons', (controls) => {
  if (!Array.isArray(controls)) return;
  const t = controls.find(c => c.name === 'token');
  if (!t?.tools) return;
  t.tools.push({
    name: 'ftm-tunnel',
    title: game.i18n.localize('FTM_OPEN'),
    icon: 'fas fa-network-wired',
    onClick: () => game.ftm?.app?.render({ force: true }),
    button: true,
  });
});

// Foundry v14: ApplicationV2 with HandlebarsMixin for templates.
const { HandlebarsApplicationMixin } = foundry.applications.api;
const { ApplicationV2 } = foundry.applications;

class FtmDashboard extends HandlebarsApplicationMixin(ApplicationV2) {
  static DEFAULT_OPTIONS = {
    id: 'ftm-tunnel-app',
    classes: ['ftm-tunnel'],
    window: { title: 'FTM_TUNNEL', icon: 'fas fa-network-wired' },
    position: { width: 720, height: 560 },
  };

  static PARTS = {
    app: { template: 'modules/ftm-tunnel/templates/app.hbs' },
  };

  async _prepareContext() {
    const s = manager ? await manager.status() : { running: false, apiBase: FTM_WEB };
    return {
      connected: s.running,
      apiBase: FTM_WEB,
      installed: s.installed,
      version: s.binPath,
    };
  }

  _onRender(context, options) {
    super._onRender(context, options);
    const root = this.element;
    if (!root) return;

    root.querySelector('[data-action="refresh"]')?.addEventListener('click', () => this.render());
    root.querySelector('[data-action="stop"]')?.addEventListener('click', async () => {
      await manager?.stop();
      ui.notifications.info(game.i18n.localize('FTM_STOPPING'));
      this.render();
    });
    root.querySelector('[data-action="start"]')?.addEventListener('click', async () => {
      ui.notifications.info(game.i18n.localize('FTM_STARTING'));
      try {
        await manager?.startServer();
      } catch (err) {
        ui.notifications.error(err.message);
      }
      this.render();
    });
  }
}

Hooks.once('ready', () => {
  game.ftm = game.ftm || {};
  game.ftm.app = new FtmDashboard();
});
