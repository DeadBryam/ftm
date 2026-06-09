const MODULE_ID = 'ftm-tunnel';
const BRIDGE_PORT = 40501;

function log(...args) { console.log(`[${MODULE_ID}]`, ...args); }

function bridgeUrl(path) {
  return `http://localhost:${BRIDGE_PORT}/${path}`;
}

async function bridge(action) {
  try {
    const res = await fetch(bridgeUrl(action), { signal: AbortSignal.timeout(5000) });
    return await res.json();
  } catch (err) {
    return { error: err.message, offline: true };
  }
}

Hooks.once('init', () => {
  game.settings.register(MODULE_ID, 'port', {
    name: 'Web port',
    hint: 'Port the ftm web server should bind to (default 40500)',
    scope: 'client', config: true, type: Number, default: 40500,
  });
  game.settings.register(MODULE_ID, 'autoStart', {
    name: 'Auto-start ftm',
    hint: 'Start the ftm server automatically when Foundry is ready',
    scope: 'client', config: true, type: Boolean, default: true,
  });
  game.modules.get(MODULE_ID).api = { bridge };
});

Hooks.once('ready', async () => {
  const status = await bridge('status');
  if (status?.offline) {
    log('bridge not running');
    ui.notifications.info('ftm: Run "node modules/ftm-tunnel/ftm-bridge.js" in a terminal for the dashboard to work');
  } else if (status?.success) {
    log('bridge connected:', status);
    game.ftm = { bridgeStatus: status };

    if (game.settings.get(MODULE_ID, 'autoStart')) {
      if (!status.installed) {
        log('installing...');
        const res = await bridge('install');
        if (res?.error) return ui.notifications.error(`ftm install: ${res.error}`);
        log(`installed v${res.version}`);
      }
      if (!status.running) {
        const res = await bridge('start');
        if (res?.error) return ui.notifications.error(`ftm start: ${res.error}`);
        ui.notifications.info('ftm tunnel started');
      }
    }
  }
});

Hooks.on('getSceneControlButtons', (controls) => {
  if (!Array.isArray(controls)) return;
  const tokenControls = controls.find((c) => c.name === 'token');
  if (!tokenControls?.tools) return;
  tokenControls.tools.push({
    name: 'ftm-tunnel',
    title: 'FTM Tunnel Dashboard',
    icon: 'fas fa-network-wired',
    onClick: () => game.ftm?.app?.render(true),
  });
});

Hooks.once('ready', () => {
  game.ftm = game.ftm || {};
  game.ftm.app = new FtmDashboard();
});

class FtmDashboard extends foundry.applications.api.HandlebarsApplicationMixin(
  foundry.applications.api.ApplicationV2,
) {
  static DEFAULT_OPTIONS = {
    id: 'ftm-tunnel-app',
    title: 'FTM Tunnel Manager',
    width: 600,
    height: 400,
    resizable: true,
    template: 'modules/ftm-tunnel/templates/app.hbs',
  };

  async _prepareContext() {
    const status = await bridge('status');
    return {
      bridgeOnline: !status?.offline,
      installed: status?.installed ?? false,
      running: status?.running ?? false,
      version: status?.version ?? '—',
      binPath: status?.binPath ?? '',
      apiBase: status?.apiBase ?? 'http://localhost:40500',
      platform: status?.platform ?? '',
      error: status?.error ?? null,
    };
  }

  _onRender(...args) {
    super._onRender(...args);
    this._listen('[data-action="install"]', 'Installing...', () => bridge('install'));
    this._listen('[data-action="start"]', 'Starting...', () => bridge('start'));
    this._listen('[data-action="stop"]', null, () => bridge('stop'));
    this._listen('[data-action="refresh"]', null, () => this.render());
  }

  _listen(sel, loadingText, fn) {
    const el = this.element?.querySelector(sel);
    if (!el) return;
    el.addEventListener('click', async () => {
      if (loadingText) { el.disabled = true; el.textContent = loadingText; }
      const res = await fn();
      if (res?.error && !res?.offline) ui.notifications.error(res.error);
      await this.render();
    });
  }
}
