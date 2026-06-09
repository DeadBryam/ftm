import { FtmApp } from './ftm-app.js';

const MODULE_ID = 'ftm-tunnel';
const SETTING_AUTO_START = 'autoStart';
const SETTING_PORT = 'port';
const SETTING_MANUAL_URL = 'manualUrl';

function log(...args) {
  console.log(`[${MODULE_ID}]`, ...args);
}

Hooks.once('init', () => {
  game.settings.register(MODULE_ID, SETTING_AUTO_START, {
    name: 'Auto-start ftm',
    hint: 'Start the ftm server automatically when Foundry is ready',
    scope: 'client',
    config: true,
    type: Boolean,
    default: true,
  });

  game.settings.register(MODULE_ID, SETTING_PORT, {
    name: 'Web port',
    hint: 'Port the ftm web server should bind to (default 40500)',
    scope: 'client',
    config: true,
    type: Number,
    default: 40500,
  });

  game.settings.register(MODULE_ID, SETTING_MANUAL_URL, {
    name: 'Manual ftm URL',
    hint: 'If using Portable, enter the URL of a running ftm instance (e.g. http://localhost:40500)',
    scope: 'client',
    config: true,
    type: String,
    default: '',
  });

  const manager = new (window.FtmManager)({ port: game.settings.get(MODULE_ID, SETTING_PORT) });
  game.modules.get(MODULE_ID).api = { FtmManager: window.FtmManager, FtmApp };
  game.ftm = { manager, app: null };
});

Hooks.once('ready', async () => {
  const manager = game.ftm.manager;

  if (manager.canManageBinary) {
    log('running on Foundry Desktop — full binary management available');

    if (!game.settings.get(MODULE_ID, SETTING_AUTO_START)) {
      log('autoStart disabled, skipping');
      return;
    }

    try {
      if (!(await manager.isInstalled())) {
        log('installing ftm binary...');
        ui.notifications.info('Installing ftm binary (first run)...');
        const v = await manager.install();
        log(`installed ftm v${v}`);
      }

      log('starting ftm server...');
      await manager.startServer();
      log(`ftm-ready on ${manager.apiBase}`);

      ui.notifications.info('ftm ready — open the Tunnel sidebar button');
    } catch (err) {
      log('startup failed:', err);
      ui.notifications.error(`ftm: ${err.message}`);
    }
  } else {
    log('running on Foundry Portable — binary management not available in browser');
    log('user must run ftm manually or configure Manual URL in settings');

    const manualUrl = game.settings.get(MODULE_ID, SETTING_MANUAL_URL);
    if (manualUrl) {
      manager.apiBase = manualUrl;
      log(`using manual URL: ${manualUrl}`);
    } else {
      ui.notifications.info('ftm: Configure Manual URL in module settings, or run ftm on Desktop');
    }
  }
});

Hooks.on('getSceneControlButtons', (controls) => {
  const tokenControls = controls.find((c) => c.name === 'token');
  if (!tokenControls?.tools) return;
  tokenControls.tools.push({
    name: 'ftm-tunnel',
    title: 'Open FTM Tunnel Dashboard',
    icon: 'fas fa-network-wired',
    onClick: () => {
      if (!game.ftm?.app) {
        game.ftm = game.ftm ?? {};
        game.ftm.app = new FtmApp(game.ftm.manager);
      }
      game.ftm.app.render(true);
    },
  });
});

Hooks.on('closeGame', async () => {
  if (game.ftm?.manager?.canManageBinary) {
    log('shutting down ftm');
    await game.ftm.manager.stop();
  }
});
