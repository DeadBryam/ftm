import { exec, spawn } from 'node:child_process';
import { promisify } from 'node:util';
import { chmod, mkdir, rename, writeFile, unlink, stat } from 'node:fs/promises';
import { createRequire } from 'node:module';
import path from 'node:path';
import os from 'node:os';

const execp = promisify(exec);
const require = createRequire(import.meta.url);

const GITHUB_API = 'https://api.github.com';
const DEFAULT_REPO = 'sthbryan/ftm';
const PLATFORM_ALIAS = { darwin: 'macos', linux: 'linux', win32: 'windows' };
const ARCH_ALIAS = { x64: 'x64', arm64: 'arm64', aarch64: 'arm64' };

function resolveDataPath() {
  if (process.env.FOUNDRY_DATA_PATH) return process.env.FOUNDRY_DATA_PATH;
  if (typeof globalThis.game !== 'undefined' && globalThis.game?.userData?.path) {
    return globalThis.game.userData.path;
  }
  return null;
}

function platformAssetName(tag) {
  const osName = PLATFORM_ALIAS[os.platform()] ?? os.platform();
  const archName = ARCH_ALIAS[os.arch()] ?? os.arch();
  return `ftm-${osName}-${archName}`;
}

export class FtmManager {
  constructor(opts = {}) {
    this.dataPath = opts.dataPath ?? resolveDataPath();
    this.binDir = opts.binDir ?? (this.dataPath ? path.join(this.dataPath, 'ftm') : null);
    this.binName = os.platform() === 'win32' ? 'ftm.exe' : 'ftm';
    this.binPath = opts.binPath ?? (this.binDir ? path.join(this.binDir, this.binName) : null);
    this.apiBase = opts.apiBase ?? `http://localhost:${opts.port ?? 40500}`;
    this.repo = opts.repo ?? DEFAULT_REPO;
    this.fetchFn = opts.fetchFn ?? globalThis.fetch.bind(globalThis);
    this.execFn = opts.execFn ?? execp;
    this.spawnFn = opts.spawnFn ?? spawn;
    this._child = null;
  }

  assertReady() {
    if (!this.binPath) {
      throw new Error('FtmManager: cannot resolve binPath. Set dataPath or binPath.');
    }
  }

  async isInstalled() {
    this.assertReady();
    try {
      await stat(this.binPath);
      return true;
    } catch {
      return false;
    }
  }

  async version() {
    this.assertReady();
    const { stdout } = await this.execFn(`"${this.binPath}" --version`);
    return stdout.trim();
  }

  async install() {
    this.assertReady();
    if (!this.fetchFn) throw new Error('no fetch function available');

    const releaseRes = await this.fetchFn(`${GITHUB_API}/repos/${this.repo}/releases/latest`, {
      headers: { Accept: 'application/vnd.github+json', 'User-Agent': 'ftm-foundry' },
    });
    if (!releaseRes.ok) {
      throw new Error(`GitHub API status ${releaseRes.status}`);
    }
    const release = await releaseRes.json();
    if (release.prerelease || release.draft) {
      throw new Error('latest release is a prerelease/draft');
    }

    const asset = release.assets.find((a) => a.name === platformAssetName());
    if (!asset) {
      throw new Error(`no asset ${platformAssetName()} in ${release.tag_name}`);
    }

    await mkdir(this.binDir, { recursive: true });

    const tmpPath = this.binPath + '.download';
    const assetRes = await this.fetchFn(asset.browser_download_url, {
      headers: { 'User-Agent': 'ftm-foundry' },
    });
    if (!assetRes.ok) {
      throw new Error(`asset download status ${assetRes.status}`);
    }
    const buf = Buffer.from(await assetRes.arrayBuffer());
    await writeFile(tmpPath, buf);

    try {
      await rename(tmpPath, this.binPath);
    } catch (err) {
      await writeFile(this.binPath, buf);
      await unlink(tmpPath).catch(() => {});
      if (err.code !== 'ENOENT') throw err;
    }

    if (os.platform() !== 'win32') {
      await chmod(this.binPath, 0o755);
    }

    return release.tag_name.replace(/^v/, '');
  }

  async ensureInstalled() {
    if (await this.isInstalled()) return await this.version();
    return await this.install();
  }

  isRunning() {
    return this._child !== null && this._child.exitCode === null;
  }

  async startServer() {
    if (this.isRunning()) return;
    this.assertReady();

    this._child = this.spawnFn(this.binPath, ['--web', '--server'], {
      stdio: ['ignore', 'pipe', 'pipe'],
      windowsHide: true,
      detached: false,
    });

    this._child.on('exit', (code) => {
      this._child = null;
    });

    await this.waitForApi();
  }

  async waitForApi(timeoutMs = 15000) {
    const deadline = Date.now() + timeoutMs;
    while (Date.now() < deadline) {
      try {
        const res = await this.fetchFn(`${this.apiBase}/api/status`, {
          signal: AbortSignal.timeout(1000),
        });
        if (res.ok) return;
      } catch {
        // not ready yet
      }
      await new Promise((r) => setTimeout(r, 250));
    }
    throw new Error(`ftm server did not respond on ${this.apiBase} within ${timeoutMs}ms`);
  }

  async stop() {
    if (!this._child) return;
    return new Promise((resolve) => {
      this._child.once('exit', () => resolve());
      try {
        this._child.kill();
      } catch {
        resolve();
      }
      setTimeout(() => resolve(), 2000);
    });
  }

  async getTunnels() {
    const res = await this.fetchFn(`${this.apiBase}/api/tunnels`);
    if (!res.ok) throw new Error(`tunnels status ${res.status}`);
    return res.json();
  }
}
