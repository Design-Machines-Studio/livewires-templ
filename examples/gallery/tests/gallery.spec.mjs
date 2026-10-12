import { test, expect } from '@playwright/test';
import { createHash } from 'node:crypto';
import { writeFile } from 'node:fs/promises';

const names = ['relative-time', 'copy', 'loading', 'toast', 'sortable', 'inline-edit'];
for (const [device, viewport] of [['desktop', { width: 1440, height: 1000 }], ['mobile', { width: 390, height: 844 }]]) {
  for (const name of names) {
    test(`${name} ${device}: generated markup, verified assets and provenance`, async ({ page, request }, info) => {
      await page.setViewportSize(viewport);
      const errors = [];
      page.on('pageerror', error => errors.push(error.message));
      page.on('requestfailed', req => errors.push(req.url() + ' ' + req.failure().errorText));
      page.on('response', res => { if (res.status() >= 400) errors.push(res.url() + ' HTTP ' + res.status()); });
      const health = await (await request.get('/healthz')).json();
      await page.goto('/examples/' + name);
      await expect(page.locator('body')).toHaveAttribute('data-demo-ready', 'true');
      await expect(page.locator('h1')).toHaveText(name);
      await expect(page.locator('#source-build')).toHaveAttribute('data-source-head', health.source.commit);
      await expect(page.locator('#source-build')).toHaveAttribute('data-source-fingerprint', health.source.fingerprint);
      await page.getByText('Build provenance', { exact: true }).click();
      for (const asset of health.assets.filter(a => a.serve)) {
        const response = await request.get(asset.serve);
        expect(response.ok()).toBeTruthy();
        expect(createHash('sha256').update(await response.body()).digest('hex')).toBe(asset.sha256);
      }
      expect(await page.locator('html').evaluate(e => e.scrollWidth <= innerWidth)).toBeTruthy();
      expect(errors).toEqual([]);
      const screenshot = info.outputPath(`${name}-${device}-${health.source.commit.slice(0, 12)}-${health.source.fingerprint.slice(0, 12)}.png`);
      await page.screenshot({ path: screenshot, fullPage: true });
      await writeFile(info.outputPath('source.json'), JSON.stringify({ name, device, url: page.url(), source: health.source, assets: health.assets }, null, 2));
      await info.attach('Source-bound screenshot', { path: screenshot, contentType: 'image/png' });
    });
  }
}

test('relative time enhances valid date and retains invalid fallback', async ({ page }) => {
  await page.goto('/examples/relative-time');
  await expect(page.locator('[data-relative-time-output]').first()).toBeVisible();
  await expect(page.getByText('Unknown date; absolute fallback remains visible.')).toBeVisible();
});
test('clipboard reports actual granted write and preserves selectable source', async ({ page, context }) => {
  await context.grantPermissions(['clipboard-read', 'clipboard-write']);
  await page.goto('/examples/copy');
  await page.getByRole('button', { name: 'Copy address' }).click();
  await expect(page.locator('[data-copy-status]')).toContainText('Copied');
  expect(await page.evaluate(() => navigator.clipboard.readText())).toBe('https://example.test/meeting');
  await expect(page.locator('[data-copy-source]')).toHaveAttribute('readonly', '');
});
test('loading delay, unrelated editing, failure, unknown and recovery', async ({ page }) => {
  await page.goto('/examples/loading');
  await expect(page.locator('body')).toHaveAttribute('data-demo-ready', 'true');
  await page.selectOption('#scenario', 'delayed');
  await page.getByRole('button', { name: 'Run scenario' }).click();
  await expect(page.locator('[data-loading-content]')).toHaveAttribute('aria-busy', 'true');
  await expect(page.getByRole('button', { name: 'Run scenario' })).toBeDisabled();
  await page.locator('#scenario-form').evaluate(form => form.requestSubmit());
  await expect(page.locator('#demo-outcome')).toHaveText('Simulated request pending.');
  await page.locator('#unrelated-input').fill('Still editable');
  await page.getByRole('button', { name: 'Update unrelated note' }).click();
  await expect(page.locator('[data-loading-content]')).toHaveAttribute('aria-busy', 'true');
  await expect(page.locator('#demo-outcome')).toContainText('Simulated success');
  for (const scenario of ['failure', 'unknown', 'normal']) {
    await page.selectOption('#scenario', scenario);
    await page.getByRole('button', { name: 'Run scenario' }).click();
    await expect(page.locator('#demo-outcome')).toContainText(scenario === 'normal' ? 'Simulated success' : 'Simulated ' + scenario);
    await expect(page.locator('[data-loading-content]')).toHaveAttribute('aria-busy', 'false');
  }
});
test('toast renders new occurrences and dismisses through producer API', async ({ page }) => {
  await page.goto('/examples/toast');
  await expect(page.locator('body')).toHaveAttribute('data-demo-ready', 'true');
  await page.getByRole('button', { name: 'Dismiss example message' }).click();
  await expect(page.locator('[data-toast-id="demo:initial"]')).toBeHidden();
  await page.selectOption('#scenario', 'failure');
  await page.getByRole('button', { name: 'Run scenario' }).click();
  await expect(page.locator('[data-toast-id="demo:outcome:1"]')).toContainText('Simulated failure');
  await expect(page.locator('[data-toast-id="demo:outcome:1"]')).toHaveClass(/toast--error/);
});
test('sortable confirms, rejects, delays and recovers unknown without persistence', async ({ page }) => {
  await page.goto('/examples/sortable');
  await expect(page.locator('body')).toHaveAttribute('data-demo-ready', 'true');
  const order = () => page.locator('[data-sortable-item]').evaluateAll(rows => rows.map(e => e.dataset.itemId));
  await page.getByRole('button', { name: 'Move Field notes down', exact: true }).click();
  expect(await order()).toEqual(['paper', 'notes', 'index']);
  await page.selectOption('#scenario', 'failure');
  await page.getByRole('button', { name: 'Move Field notes down', exact: true }).click();
  await expect(page.locator('#demo-outcome')).toContainText('Simulated failure');
  expect(await order()).toEqual(['paper', 'notes', 'index']);
  await page.selectOption('#scenario', 'delayed');
  await page.getByRole('button', { name: 'Move Field notes down', exact: true }).click();
  await expect(page.locator('lw-sortable-list')).toHaveAttribute('data-pending', '');
  await page.getByRole('button', { name: 'Update unrelated note' }).click();
  await expect(page.locator('lw-sortable-list')).not.toHaveAttribute('data-pending', '');
  await page.selectOption('#scenario', 'unknown');
  await page.getByRole('button', { name: 'Move Field notes up', exact: true }).click();
  await expect(page.locator('#demo-outcome')).toContainText('Simulated unknown');
  await expect(page.locator('lw-sortable-list')).toHaveAttribute('data-pending', '');
  await page.getByRole('link', { name: 'Reload and recover' }).click();
  expect(await order()).toEqual(['notes', 'paper', 'index']);
});
test('native fallback: content, URL-carried move and simulated response', async ({ browser }) => {
  const context = await browser.newContext({ javaScriptEnabled: false });
  const page = await context.newPage();
  for (const name of names) {
    await page.goto(process.env.GALLERY_URL + '/examples/' + name);
    await expect(page.locator('h1')).toHaveText(name);
  }
  await page.goto(process.env.GALLERY_URL + '/examples/sortable');
  await page.getByRole('button', { name: 'Move Field notes down', exact: true }).click();
  await expect(page.locator('[data-sortable-item]').first()).toHaveAttribute('data-item-id', 'paper');
  await page.goto(process.env.GALLERY_URL + '/examples/loading');
  await page.selectOption('#scenario', 'failure');
  await page.getByRole('button', { name: 'Run scenario' }).click();
  await expect(page.locator('#demo-outcome')).toContainText('Simulated failure');
  await page.goto(process.env.GALLERY_URL + '/examples/copy');
  await expect(page.locator('[data-copy-source]')).toHaveValue('https://example.test/meeting');
  await page.goto(process.env.GALLERY_URL + '/examples/relative-time');
  await expect(page.getByText('October 10, 2026', { exact: true })).toBeVisible();
  await page.goto(process.env.GALLERY_URL + '/examples/inline-edit');
  await page.getByLabel('First title', { exact: true }).fill('Fictional native title');
  await page.getByRole('button', { name: 'Save', exact: true }).click();
  await expect(page.locator('#demo-outcome')).toContainText('Native example submission');
  expect(new URL(page.url()).searchParams.get('title-a')).toBe('Fictional native title');
  await context.close();
});

test('inline edit retains draft/focus, rejects, recovers and ignores cancelled delayed result', async ({ page }) => {
  await page.goto('/examples/inline-edit');
  await expect(page.locator('body')).toHaveAttribute('data-demo-ready', 'true');
  const root = page.locator('#title-a');
  const input = page.getByLabel('First title', { exact: true });
  await page.getByRole('link', { name: 'Edit first title' }).click();
  await input.fill('Fictional new title');
  await input.press('ArrowLeft');
  const caret = await input.evaluate(e => e.selectionStart);
  await page.getByRole('button', { name: 'Simulate external update' }).click();
  await expect(input).toHaveValue('Fictional new title');
  expect(await input.evaluate(e => e.selectionStart)).toBe(caret);
  await page.selectOption('#scenario', 'failure');
  await root.getByRole('button', { name: 'Save', exact: true }).click();
  await expect(root).toHaveAttribute('data-inline-edit-state', 'rejected');
  await expect(input).toHaveValue('Fictional new title');
  await page.selectOption('#scenario', 'normal');
  await root.getByRole('button', { name: 'Save', exact: true }).click();
  await expect(root).toHaveAttribute('data-inline-edit-state', 'confirmed');
  await expect(root.locator('[data-inline-edit-value]')).toHaveText('Fictional new title');
  await page.getByRole('link', { name: 'Edit first title' }).click();
  await input.fill('Delayed fictional title');
  await page.selectOption('#scenario', 'delayed');
  await root.getByRole('button', { name: 'Save', exact: true }).click();
  await expect(root).toHaveAttribute('data-inline-edit-state', 'pending');
  await root.getByRole('button', { name: 'Cancel', exact: true }).click();
  await expect(root).toHaveAttribute('data-inline-edit-state', 'confirmed');
  await page.getByRole('link', { name: 'Edit first title' }).click();
  await input.fill('New draft after cancellation');
  await page.waitForTimeout(900); // Bounded delayed-response observation.
  await expect(input).toHaveValue('New draft after cancellation');
  await expect(root).toHaveAttribute('data-inline-edit-state', 'editing');
  await page.selectOption('#scenario', 'unknown');
  await root.getByRole('button', { name: 'Save', exact: true }).click();
  await expect(root).toHaveAttribute('data-inline-edit-state', 'unknown');
  await page.getByRole('link', { name: 'Reload and recover' }).click();
  await expect(root.locator('[data-inline-edit-value]')).toHaveText('Community gathering');
  await page.getByText('Edit second title', { exact: true }).click();
  const second = page.getByLabel('Second title', { exact: true });
  await second.fill('Fictional second title');
  await second.press('Enter');
  await expect(page.locator('#title-b')).toHaveAttribute('data-inline-edit-state', 'confirmed');
});
