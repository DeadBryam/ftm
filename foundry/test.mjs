#!/usr/bin/env node
// Standalone smoke test for FtmManager. No Foundry required.
// Run: FOUNDRY_DATA_PATH=/tmp/ftm-test node foundry/test.mjs

import { FtmManager } from './scripts/ftm-manager.js';

const dataPath = process.env.FOUNDRY_DATA_PATH || '/tmp/ftm-test';
console.log(`[test] dataPath = ${dataPath}`);

const m = new FtmManager({ dataPath });

console.log('[test] isInstalled:', await m.isInstalled());

console.log('[test] install...');
const v = await m.install();
console.log('[test] installed version:', v);

console.log('[test] version check:', await m.version());

console.log('[test] startServer...');
await m.startServer();
console.log('[test] server is up');

const tunnels = await m.getTunnels();
console.log(`[test] ${tunnels.length} tunnel(s) configured`);

console.log('[test] status:', await m.status());

console.log('[test] stop...');
await m.stop();
console.log('[test] stopped');

console.log('[test] OK');
