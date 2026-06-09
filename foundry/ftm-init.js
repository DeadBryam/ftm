//
// ftm-init.js — Inicio Rápido para Foundry VTT
//
// INSTRUCCIONES:
//   1. Ve a Settings → World Settings → Init Scripts
//   2. Pega este código
//   3. Guarda
//   4. Reinicia Foundry
//
// El binario se descarga solo y se ejecuta en segundo plano.
// Se apaga automáticamente cuando cierras Foundry.
//

const { spawn } = require('child_process');
const fs = require('fs');
const path = require('path');
const https = require('https');

const REPO = 'sthbryan/ftm';
const BIN_DIR = path.join(process.env.FOUNDRY_DATA_PATH || '.', 'ftm');
const BIN_NAME = process.platform === 'win32' ? 'ftm.exe' : 'ftm';
const BIN_PATH = path.join(BIN_DIR, BIN_NAME);

function platformAsset() {
  const map = { darwin: 'macos', linux: 'linux', win32: 'windows' };
  const osName = map[process.platform] || process.platform;
  const arch = process.arch === 'arm64' ? 'arm64' : 'x64';
  const ext = process.platform === 'win32' ? '.exe' : '';
  return `ftm-${osName}-${arch}${ext}`;
}

function download(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);
    https.get(url, (res) => {
      if (res.statusCode !== 200) {
        reject(new Error(`HTTP ${res.statusCode}`));
        return;
      }
      res.pipe(file);
      file.on('finish', () => file.close(resolve));
    }).on('error', reject);
  });
}

async function main() {
  // Create bin dir
  fs.mkdirSync(BIN_DIR, { recursive: true });

  // Download if missing
  if (!fs.existsSync(BIN_PATH)) {
    console.log('[ftm] downloading binary...');
    const url = `https://github.com/${REPO}/releases/latest/download/${platformAsset()}`;
    await download(url, BIN_PATH + '.tmp');
    fs.renameSync(BIN_PATH + '.tmp', BIN_PATH);
    if (process.platform !== 'win32') fs.chmodSync(BIN_PATH, 0o755);
    console.log('[ftm] downloaded');
  }

  // Start server
  const child = spawn(BIN_PATH, ['--web', '--server'], {
    stdio: 'ignore',
    windowsHide: true,
    detached: false,
  });
  console.log(`[ftm] started (pid ${child.pid})`);

  // Cleanup on Foundry shutdown
  process.on('exit', () => { child.kill(); });
}

main().catch(err => console.error('[ftm] error:', err.message));
