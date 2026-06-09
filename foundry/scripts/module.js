const MODULE_ID = 'ftm-tunnel';
const BRIDGE = 'http://localhost:40501';

function log(...args) { console.log(`[${MODULE_ID}]`, ...args); }

async function bridge(action) {
  try {
    const res = await fetch(`${BRIDGE}/${action}`, { signal: AbortSignal.timeout(5000) });
    return await res.json();
  } catch (err) {
    return { error: err.message, offline: true };
  }
}

Hooks.once('init', () => {
  game.settings.register(MODULE_ID, 'port', {
    name: 'Web port', hint: 'Port the ftm web server (default 40500)',
    scope: 'client', config: true, type: Number, default: 40500,
  });
  game.modules.get(MODULE_ID).api = { bridge };
});

Hooks.once('ready', async () => {
  const status = await bridge('status');
  if (status?.offline) {
    log('bridge not running');
    ui.notifications.info('ftm: Run "node modules/ftm-tunnel/ftm-bridge.js" once — it auto-installs to Windows Startup');
  } else {
    log('bridge online:', status);
    game.ftm = { bridgeStatus: status };

    if (!status.installed) {
      log('auto-installing...');
      const res = await bridge('install');
      if (res?.error) return ui.notifications.error(`ftm install: ${res.error}`);
      log(`installed v${res.version}`);
    }
  }
});

Hooks.on('getSceneControlButtons', (controls) => {
  if (!Array.isArray(controls)) return;
  const token = controls.find(c => c.name === 'token');
  if (!token?.tools) return;
  token.tools.push({
    name: 'ftm-tunnel', title: 'FTM Tunnel Dashboard',
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
    id: 'ftm-tunnel-app', title: 'FTM Tunnel Manager',
    width: 600, height: 400, resizable: true,
    template: 'modules/ftm-tunnel/templates/app.hbs',
  };

  async _prepareContext() {
    const s = await bridge('status');
    return {
      bridgeOnline: !s?.offline,
      installed: s?.installed ?? false,
      running: s?.running ?? false,
      version: s?.version ?? '—',
      binPath: s?.binPath ?? '',
      apiBase: s?.apiBase ?? 'http://localhost:40500',
      error: s?.error ?? null,
    };
  }

  _onRender(...args) {
    super._onRender(...args);
    ['install', 'start', 'stop', 'refresh'].forEach(a => {
      const el = this.element?.querySelector(`[data-action="${a}"]`);
      if (!el) return;
      el.addEventListener('click', async () => {
        if (a === 'install') { el.disabled = true; el.textContent = 'Installing...'; }
        const res = await bridge(a);
        if (res?.error && !res?.offline) ui.notifications.error(res.error);
        await this.render();
      });
    });
  }
}
