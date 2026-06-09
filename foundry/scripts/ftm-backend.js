(function () {
  'use strict';

  const _require =
    typeof window !== 'undefined' &&
    (window.require || window.nodeRequire);

  const PLATFORM_ALIAS = { darwin: 'macos', linux: 'linux', win32: 'windows' };
  const ARCH_ALIAS = { x64: 'x64', arm64: 'arm64', aarch64: 'arm64' };
  const GITHUB_API = 'https://api.github.com';
  const DEFAULT_REPO = 'sthbryan/ftm';

  function resolveDataPath() {
    if (typeof process !== 'undefined' && process.env?.FOUNDRY_DATA_PATH) return process.env.FOUNDRY_DATA_PATH;
    if (typeof globalThis.game !== 'undefined' && globalThis.game?.userData?.path) {
      return globalThis.game.userData.path;
    }
    return null;
  }

  function platformAssetName() {
    const osName = PLATFORM_ALIAS[_require ? _require('os').platform() : 'win32'] ?? 'win32';
    const archName = ARCH_ALIAS[_require ? _require('os').arch() : 'x64'] ?? 'x64';
    return `ftm-${osName}-${archName}`;
  }

  class FtmManager {
    constructor(opts = {}) {
      this.dataPath = opts.dataPath ?? resolveDataPath();
      this.binDir = opts.binDir ?? (this.dataPath ? pathJoin(this.dataPath, 'ftm') : null);
      this.binName = 'ftm.exe';
      this.binPath = opts.binPath ?? (this.binDir ? pathJoin(this.binDir, this.binName) : null);
      this.apiBase = opts.apiBase ?? `http://localhost:${opts.port ?? 40500}`;
      this.repo = opts.repo ?? DEFAULT_REPO;
      this.fetchFn = opts.fetchFn ?? globalThis.fetch.bind(globalThis);
      this._child = null;
      this._hasNode = !!_require;
    }

    get canManageBinary() {
      return this._hasNode;
    }

    assertReady() {
      if (!this.binPath) {
        throw new Error('FtmManager: cannot resolve binPath. Set dataPath or binPath.');
      }
    }

    async isInstalled() {
      if (!this._hasNode) return false;
      this.assertReady();
      try {
        const fs = _require('fs/promises');
        await fs.stat(this.binPath);
        return true;
      } catch {
        return false;
      }
    }

    async version() {
      if (!this._hasNode) return '0.0.0 (browser mode)';
      this.assertReady();
      const { exec } = _require('child_process');
      const { promisify } = _require('util');
      const { stdout } = await promisify(exec)(`"${this.binPath}" --version`);
      return stdout.trim();
    }

    async install() {
      if (!this._hasNode) throw new Error('Binary management requires Foundry Desktop');
      this.assertReady();
      if (!this.fetchFn) throw new Error('no fetch function available');

      const releaseRes = await this.fetchFn(`${GITHUB_API}/repos/${this.repo}/releases/latest`, {
        headers: { Accept: 'application/vnd.github+json', 'User-Agent': 'ftm-foundry' },
      });
      if (!releaseRes.ok) throw new Error(`GitHub API status ${releaseRes.status}`);
      const release = await releaseRes.json();
      if (release.prerelease || release.draft) throw new Error('latest release is a prerelease/draft');

      const assetName = platformAssetName();
      const asset = release.assets.find((a) => a.name === assetName);
      if (!asset) throw new Error(`no asset ${assetName} in ${release.tag_name}`);

      const fs = _require('fs/promises');
      const path = _require('path');
      await fs.mkdir(this.binDir, { recursive: true });

      const tmpPath = this.binPath + '.download';
      const assetRes = await this.fetchFn(asset.browser_download_url, {
        headers: { 'User-Agent': 'ftm-foundry' },
      });
      if (!assetRes.ok) throw new Error(`asset download status ${assetRes.status}`);
      const buf = Buffer.from(await assetRes.arrayBuffer());
      await fs.writeFile(tmpPath, buf);

      try {
        await fs.rename(tmpPath, this.binPath);
      } catch (err) {
        await fs.writeFile(this.binPath, buf);
        await fs.unlink(tmpPath).catch(() => {});
        if (err.code !== 'ENOENT') throw err;
      }

      const os = _require('os');
      if (os.platform() !== 'win32') {
        await fs.chmod(this.binPath, 0o755);
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
      if (!this._hasNode) throw new Error('Binary management requires Foundry Desktop');
      if (this.isRunning()) return;
      this.assertReady();

      const { spawn } = _require('child_process');
      this._child = spawn(this.binPath, ['--web', '--server'], {
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
        } catch {}
        await new Promise((r) => setTimeout(r, 250));
      }
      throw new Error(`ftm server did not respond on ${this.apiBase} within ${timeoutMs}ms`);
    }

    async stop() {
      if (!this._child) return;
      return new Promise((resolve) => {
        this._child.once('exit', () => resolve());
        try { this._child.kill(); } catch { resolve(); }
        setTimeout(() => resolve(), 2000);
      });
    }

    async getTunnels() {
      const res = await this.fetchFn(`${this.apiBase}/api/tunnels`);
      if (!res.ok) throw new Error(`tunnels status ${res.status}`);
      return res.json();
    }
  }

  function pathJoin(...parts) {
    if (_require) return _require('path').join(...parts);
    return parts.join('/');
  }

  window.FtmManager = FtmManager;
})();
