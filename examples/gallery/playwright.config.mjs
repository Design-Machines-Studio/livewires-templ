import { defineConfig } from '@playwright/test';
if (!process.env.GALLERY_URL) throw new Error('Set GALLERY_URL to the actual loopback URL from run.sh status.');
const url = new URL(process.env.GALLERY_URL);
if (url.hostname !== '127.0.0.1' || url.protocol !== 'http:') throw new Error('Use the task-owned loopback HTTP gallery.');
export default defineConfig({
  testDir: './tests', workers: 1, retries: 0, timeout: 30000,
  outputDir: process.env.GALLERY_EVIDENCE || './test-results',
  reporter: [['list'], ['json', { outputFile: (process.env.GALLERY_EVIDENCE || './test-results') + '/results.json' }]],
  use: { baseURL: url.origin, browserName: 'chromium', screenshot: 'only-on-failure', trace: 'retain-on-failure' },
});
