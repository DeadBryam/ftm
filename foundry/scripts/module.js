const MODULE_ID = 'ftm-tunnel';
const FTM_API = 'http://localhost:40500';

function log(...args) { console.log(`[${MODULE_ID}]`, ...args); }

async function checkFtm() {
  try {
    const res = await fetch(`${FTM_API}/api/status`, { signal: AbortSignal.timeout(3000) });
    return res.ok ? { connected: true, apiBase: FTM_API } : { connected: false };
  } catch {
    return { connected: false };
  }
}

async function bridge(cmd) {
  try {
    const res = await fetch(`http://localhost:40501/${cmd}`, { signal: AbortSignal.timeout(8000) });
    return await res.json();
  } catch {
    return { error: 'bridge offline' };
  }
}

Hooks.once('init', () => {
  game.modules.get(MODULE_ID).api = { checkFtm, bridge };
});

Hooks.once('ready', async () => {
  const status = await checkFtm();
  log(status.connected ? 'ftm running' : 'ftm not running');
  game.ftm = game.ftm || {};
  game.ftm.app = new FtmDashboard();
});

Hooks.on('getSceneControlButtons', (controls) => {
  if (!Array.isArray(controls)) return;
  const t = controls.find(c => c.name === 'token');
  if (!t?.tools) return;
  t.tools.push({
    name: 'ftm-tunnel', title: 'FTM Tunnel',
    icon: 'fas fa-network-wired',
    onClick: () => game.ftm?.app?.render(true),
  });
});

class FtmDashboard extends foundry.applications.api.HandlebarsApplicationMixin(
  foundry.applications.api.ApplicationV2,
) {
  static DEFAULT_OPTIONS = {
    id: 'ftm-tunnel-app', title: 'FTM Tunnel Manager',
    width: 600, height: 500, resizable: true,
    template: 'modules/ftm-tunnel/templates/app.hbs',
  };

  async _prepareContext() {
    const s = await checkFtm();
    return {
      connected: s.connected,
      apiBase: s.apiBase || '',
    };
  }

  _onRender(...args) {
    super._onRender(...args);
    ['stop', 'refresh'].forEach(a => {
      const el = this.element?.querySelector(`[data-action="${a}"]`);
      if (!el) return;
      el.addEventListener('click', async () => {
        if (a === 'stop') await bridge('stop');
        await this.render();
      });
    });
  }
}
