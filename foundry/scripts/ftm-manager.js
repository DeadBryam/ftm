// scripts/ftm-manager.js
// Pure Node manager for the ftm binary. No Foundry globals, no DOM.
// Foundry v14 portable: dataPath comes from game.userData?.path in module.js,
// or from FOUNDRY_DATA_PATH env when running standalone (test.mjs).

import { spawn, execFileSync } from 'node:child_process';
import { promises as fs, createWriteStream, constants as fsConstants } from 'node:fs';
import path from 'node:path';
import os from 'node:os';
import https from 'node:https';

const REPO = 'sthbryan/ftm';
const API_BASE = 'http://localhost:40500';
const GITHUB_LATEST = `https://api.github.com/repos/${REPO}/releases/latest`;

const PLATFORM_MAP = { darwin: 'macos', linux: 'linux', win32: 'windows' };

function platformAsset() {
  const osName = PLATFORM_MAP[os.platform()] || os.platform();
  const arch = os.arch() === 'arm64' ? 'arm64' : 'x64';
  const ext = os.platform() === 'win32' ? '.exe' : '';
  return `ftm-${osName}-${arch}${ext}`;
}

function downloadToFile(url, dest) {
  return new Promise((resolve, reject) => {
    const follow = (u) => {
      https.get(u, { headers: { 'User-Agent': 'ftm-foundry-module' } }, (res) => {
        if (res.statusCode === 302 || res.statusCode === 301) {
          res.resume();
          return follow(res.headers.location);
        }
        if (res.statusCode !== 200) {
          res.resume();
          return reject(new Error(`HTTP ${res.statusCode} for ${u}`));
        }
        const out = createWriteStream(dest);
        res.pipe(out);
        out.on('finish', () => out.close(resolve));
        out.on('error', reject);
      }).on('error', reject);
    };
    follow(url);
  });
}

async function getLatestRelease() {
  const res = await fetch(GITHUB_LATEST, {
    headers: { Accept: 'application/vnd.github+json', 'User-Agent': 'ftm-foundry-module' },
  });
  if (!res.ok) throw new Error(`GitHub API ${res.status}`);
  return res.json();
}

export class FtmManager {
  constructor({ dataPath, port = 40500 } = {}) {
    this.dataPath = dataPath || process.env.FOUNDRY_DATA_PATH || process.cwd();
    this.binDir = path.join(this.dataPath, 'ftm');
    this.binName = os.platform() === 'win32' ? 'ftm.exe' : 'ftm';
    this.binPath = path.join(this.binDir, this.binName);
    this.port = port;
    this.child = null;
  }

  async isInstalled() {
    try {
      const s = await fs.stat(this.binPath);
      return s.isFile();
    } catch {
      return false;
    }
  }

  async version() {
    if (!await this.isInstalled()) return null;
    return new Promise((resolve) => {
      const p = spawn(this.binPath, ['--version'], { stdio: ['ignore', 'pipe', 'pipe'] });
      let out = '';
      p.stdout.on('data', d => out += d.toString());
      p.on('exit', () => resolve(out.trim() || null));
      p.on('error', () => resolve(null));
    });
  }

  async install() {
    await fs.mkdir(this.binDir, { recursive: true });

    const release = await getLatestRelease();
    const asset = release.assets.find(a => a.name === platformAsset());
    if (!asset) throw new Error(`No asset ${platformAsset()} in ${release.tag_name}`);

    const tmp = this.binPath + '.tmp';
    await downloadToFile(asset.browser_download_url, tmp);
    await fs.rename(tmp, this.binPath);
    if (os.platform() !== 'win32') await fs.chmod(this.binPath, 0o755);
    this.#stripMacOSQuarantine(this.binPath);

    return release.tag_name.replace(/^v/, '');
  }

  #stripMacOSQuarantine(p) {
    if (os.platform() !== 'darwin') return;
    try {
      execFileSync('xattr', ['-dr', 'com.apple.quarantine', p], { stdio: 'ignore' });
    } catch {}
    try {
      execFileSync('xattr', ['-d', 'com.apple.provenance', p], { stdio: 'ignore' });
    } catch {}
  }

  async startServer() {
    if (this.child && this.child.exitCode === null) return { alreadyRunning: true };
    if (!await this.isInstalled()) {
      const v = await this.install();
      await this._spawn();
      return { started: true, installedVersion: v };
    }
    await this._spawn();
    return { started: true };
  }

  async _spawn() {
    this.child = spawn(this.binPath, ['--web', '--server'], {
      stdio: ['ignore', 'pipe', 'pipe'],
      windowsHide: true,
      detached: false,
    });
    this.child.on('exit', () => { this.child = null; });
    this.child.stdout.on('data', d => process.stdout.write(`[ftm] ${d}`));
    this.child.stderr.on('data', d => process.stderr.write(`[ftm] ${d}`));

    // Wait up to 15s for API to respond
    const start = Date.now();
    while (Date.now() - start < 15000) {
      try {
        const r = await fetch(`${API_BASE}/api/status`, { signal: AbortSignal.timeout(1000) });
        if (r.ok) return;
      } catch {}
      await new Promise(r => setTimeout(r, 500));
    }
  }

  async stop() {
    if (!this.child) return { alreadyStopped: true };
    this.child.kill();
    this.child = null;
    return { stopped: true };
  }

  async getTunnels() {
    try {
      const r = await fetch(`${API_BASE}/api/tunnels`, { signal: AbortSignal.timeout(2000) });
      if (!r.ok) return [];
      return await r.json();
    } catch {
      return [];
    }
  }

  async status() {
    const installed = await this.isInstalled();
    let running = false;
    try {
      const r = await fetch(`${API_BASE}/api/status`, { signal: AbortSignal.timeout(2000) });
      running = r.ok;
    } catch {}
    return {
      installed,
      running,
      binPath: this.binPath,
      apiBase: API_BASE,
      port: this.port,
      platform: os.platform(),
      arch: os.arch(),
    };
  }
}
