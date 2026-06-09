/**
 * ftm-bridge — Sidecar for Foundry VTT module
 * Manages the ftm binary lifecycle (download, start, stop).
 *
 * Usage: node ftm-bridge.js
 * The Foundry module connects to http://localhost:40501
 */
const http = require('http');
const { spawn } = require('child_process');
const fs = require('fs/promises');
const path = require('path');
const os = require('os');

const PORT = 40501;
const GITHUB_API = 'https://api.github.com';
const REPO = 'sthbryan/ftm';

const BIN_DIR = path.join(__dirname, 'bin');
const BIN_NAME = os.platform() === 'win32' ? 'ftm.exe' : 'ftm';
const BIN_PATH = path.join(BIN_DIR, BIN_NAME);

// ---- Auto-install to Windows Startup ----
function installStartup() {
  if (os.platform() !== 'win32') return;
  const startupDir = path.join(
    os.homedir(),
    'AppData', 'Roaming', 'Microsoft', 'Windows', 'Start Menu', 'Programs', 'Startup',
  );
  const batPath = path.join(startupDir, 'ftm-bridge.bat');
  const batContent = `@echo off\nstart /MIN "" node "${__filename}"\n`;
  fs.mkdir(startupDir, { recursive: true }).then(() => {
    fs.writeFile(batPath, batContent).then(() => {
      console.log(`  ✓ Added to Windows Startup: ${batPath}`);
    }).catch(() => {});
  }).catch(() => {});
}
installStartup();

const PLATFORM_ALIASES = { darwin: 'macos', linux: 'linux', win32: 'windows' };
const ARCH_ALIASES = { x64: 'x64', arm64: 'arm64' };

let child = null;

function platformAsset() {
  const osName = PLATFORM_ALIASES[os.platform()] || os.platform();
  const arch = ARCH_ALIASES[os.arch()] || os.arch();
  const ext = os.platform() === 'win32' ? '.exe' : '';
  return `ftm-${osName}-${arch}${ext}`;
}

function json(res, data, status = 200) {
  res.writeHead(status, {
    'Content-Type': 'application/json',
    'Access-Control-Allow-Origin': '*',
    'Access-Control-Allow-Methods': 'GET, POST, OPTIONS',
    'Access-Control-Allow-Headers': 'Content-Type',
  });
  res.end(JSON.stringify(data));
}

async function handleInstall() {
  await fs.mkdir(BIN_DIR, { recursive: true });

  const releaseRes = await fetch(`${GITHUB_API}/repos/${REPO}/releases/latest`, {
    headers: { Accept: 'application/vnd.github+json', 'User-Agent': 'ftm-bridge' },
  });
  if (!releaseRes.ok) throw new Error(`GitHub API ${releaseRes.status}`);
  const release = await releaseRes.json();

  const assetName = platformAsset();
  const asset = release.assets.find(a => a.name === assetName);
  if (!asset) throw new Error(`No asset ${assetName} in ${release.tag_name}`);

  const tmp = BIN_PATH + '.download';
  const dl = await fetch(asset.browser_download_url, { headers: { 'User-Agent': 'ftm-bridge' } });
  if (!dl.ok) throw new Error(`Download ${dl.status}`);

  const buf = Buffer.from(await dl.arrayBuffer());
  await fs.writeFile(tmp, buf);
  try { await fs.rename(tmp, BIN_PATH); } catch {
    await fs.writeFile(BIN_PATH, buf);
    await fs.unlink(tmp).catch(() => {});
  }
  if (os.platform() !== 'win32') await fs.chmod(BIN_PATH, 0o755);

  return release.tag_name.replace(/^v/, '');
}

async function handleStart() {
  if (child && child.exitCode === null) return { alreadyRunning: true };
  await fs.stat(BIN_PATH);
  child = spawn(BIN_PATH, ['--web', '--server'], {
    stdio: ['ignore', 'pipe', 'pipe'],
    windowsHide: true,
    detached: false,
  });
  child.on('exit', code => { console.log(`ftm exited (${code})`); child = null; });
  child.stdout.on('data', d => process.stdout.write(d));
  child.stderr.on('data', d => process.stderr.write(d));
  return { success: true };
}

function handleStop() {
  if (!child) return { alreadyStopped: true };
  child.kill();
  child = null;
  return { success: true };
}

async function handleStatus() {
  let installed = false;
  try { await fs.stat(BIN_PATH); installed = true; } catch {}

  let running = false;
  try {
    const res = await fetch('http://localhost:40500/api/status', { signal: AbortSignal.timeout(2000) });
    running = res.ok;
  } catch {}

  return {
    installed,
    running,
    binPath: BIN_PATH,
    apiBase: 'http://localhost:40500',
    bridgePort: PORT,
    platform: os.platform(),
    arch: os.arch(),
  };
}

const handlers = {
  install: async () => ({ success: true, version: await handleInstall() }),
  start: handleStart,
  stop: handleStop,
  status: handleStatus,
};

const server = http.createServer(async (req, res) => {
  if (req.method === 'OPTIONS') return json(res, {});

  const url = new URL(req.url, `http://localhost:${PORT}`);
  const action = url.pathname.replace(/^\/+/, '') || 'status';

  const handler = handlers[action];
  if (!handler) return json(res, { error: `Unknown action: ${action}` }, 404);

  try {
    const result = await handler();
    json(res, result);
  } catch (err) {
    console.error(`${action} error:`, err.message);
    json(res, { error: err.message }, 500);
  }
});

server.listen(PORT, () => {
  console.log(`\n  ftm-bridge running on http://localhost:${PORT}`);
  console.log(`  Binary: ${BIN_PATH}`);
  console.log(`  Commands: /status, /install, /start, /stop\n`);
});
