import { test, expect } from '@playwright/test';
import { createHash } from 'node:crypto';
import { writeFile } from 'node:fs/promises';

const names = ['relative-time', 'copy', 'loading', 'toast', 'sortable', 'inline-edit', 'combobox', 'save-status', 'page-header', 'comments', 'loaders'];
for (const [device, viewport] of [['desktop', { width: 1440, height: 1000 }], ['reflow', { width: 640, height: 1000 }], ['mobile', { width: 390, height: 844 }]]) {
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
      if (name === 'page-header') {
        await expect(page.getByRole('heading', { level: 1, name: 'Title-only header' })).toBeVisible();
      } else {
        await expect(page.locator('h1')).toHaveText(name);
      }
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

for (const [device, viewport] of [['desktop', { width: 1440, height: 1000 }], ['mobile', { width: 390, height: 844 }]]) {
  test(`page header ${device}: content reflows and actions remain keyboard accessible`, async ({ page }, info) => {
    await page.setViewportSize(viewport);
    await page.goto('/examples/page-header');
    await expect(page.getByRole('region', { name: 'Title-only header' })).toBeVisible();
    await expect(page.getByRole('heading', { level: 1, name: 'Title-only header' })).toBeVisible();
    await expect(page.getByRole('heading', { level: 1, name: 'A header with a lead' })).toBeVisible();
    await expect(page.getByRole('heading', { level: 1, name: 'Member directory (27)' })).toBeVisible();
    await expect(page.locator('#actions-heading #members-count')).toHaveText('(27)');
    const network = page.getByRole('region', { name: 'Network', exact: true });
    await expect(network.locator('header > [data-header-footer]')).toContainText('Last synced:');
    const networkRows = await network.locator('header > div').evaluateAll(rows => rows.map(row => row.getBoundingClientRect().toJSON()));
    expect(networkRows[1].top).toBeGreaterThanOrEqual(networkRows[0].bottom);
    const longHeading = page.getByRole('heading', { level: 1, name: 'Review the shared decisions and records for your cooperative before making changes' });
    await expect(longHeading).toBeVisible();
    await expect(page.getByText('Supplementary content remains owned by the consuming page.')).toBeVisible();
    const wrapped = await longHeading.evaluate(node => node.getBoundingClientRect().height > parseFloat(getComputedStyle(node).fontSize) * 1.5);
    expect(wrapped).toBeTruthy();
    expect(await page.locator('html').evaluate(e => e.scrollWidth <= innerWidth)).toBeTruthy();

    const action = page.getByRole('link', { name: 'Review members' });
    const titleAndAction = await Promise.all([
      page.getByRole('heading', { level: 1, name: 'Member directory (27)' }).boundingBox(),
      action.boundingBox(),
    ]);
    if (device === 'mobile') {
      expect(titleAndAction[1].y).toBeGreaterThan(titleAndAction[0].y);
    }
    await page.locator('body').focus();
    for (let i = 0; i < 8 && !(await action.evaluate(node => node === document.activeElement)); i++) {
      await page.keyboard.press('Tab');
    }
    await expect(action).toBeFocused();
    await expect(action).toHaveAttribute('href', '#header-action-result');
    const focusStyle = await action.evaluate(node => {
      const style = getComputedStyle(node);
      return { focusVisible: node.matches(':focus-visible'), outlineStyle: style.outlineStyle, outlineWidth: style.outlineWidth, outlineColor: style.outlineColor, boxShadow: style.boxShadow, color: style.color, background: style.backgroundColor };
    });
    expect(focusStyle.focusVisible).toBeTruthy();
    expect(focusStyle.outlineStyle !== 'none' || focusStyle.boxShadow !== 'none').toBeTruthy();
    await writeFile(info.outputPath('header-accessibility.json'), JSON.stringify({
      viewport, focusStyle, accessibility: await page.locator('main').ariaSnapshot(),
      headers: await page.locator('main header').evaluateAll(headers => headers.map(header => ({
        classes: header.className, rect: header.getBoundingClientRect().toJSON(),
        children: [...header.children].map(node => ({ tag: node.tagName, classes: node.className, display: getComputedStyle(node).display, wrap: getComputedStyle(node).flexWrap, gap: getComputedStyle(node).gap, rect: node.getBoundingClientRect().toJSON() })),
      }))),
    }, null, 2));
  });
}

for (const width of [390, 1440]) {
test(`save status at ${width}px keeps one atomic region per editor and preserves focus and placement`, async ({ page }, info) => {
  await page.setViewportSize({ width, height: 1000 });
  await page.goto('/examples/save-status');
  await expect(page.locator('body')).toHaveAttribute('data-demo-ready', 'true');
  const regions = page.locator('[data-save-region]');
  const first = page.locator('[data-save-region="proposal"]');
  const line = first.locator('[data-save-status]');
  expect(await page.locator('[role="status"]').count()).toBe(2);
  expect(await regions.locator('[role="status"] [role="status"]').count()).toBe(0);
  await expect(line).toHaveAttribute('data-save-state', 'ready');
  await expect(line).toHaveText('Ready');
  const before = await first.evaluate(node => {
    const rect = node.getBoundingClientRect();
    return { x: rect.x + scrollX, y: rect.y + scrollY };
  });
  await page.locator('#proposal-text').fill('A new text edit does not trigger a status announcement by itself.');
  await expect(line).toHaveAttribute('data-save-state', 'ready');

  await page.getByRole('button', { name: 'Edit', exact: true }).click();
  await expect(line).toHaveAttribute('data-save-state', 'dirty');
  await expect(line).toHaveText('Unsaved changes');
  await expect(page.getByRole('button', { name: 'Edit', exact: true })).toBeFocused();
  await page.getByRole('button', { name: 'Saving', exact: true }).click();
  await expect(line).toHaveAttribute('data-save-state', 'saving');
  await page.getByRole('button', { name: 'Confirm save', exact: true }).click();
  await expect(line).toHaveAttribute('data-save-state', 'saved');
  await expect(line).toHaveText('Draft saved');

  await page.getByRole('button', { name: 'Start pending save', exact: true }).click();
  await page.getByRole('button', { name: 'Edit while save is pending', exact: true }).click();
  await expect(line).toHaveAttribute('data-save-state', 'dirty');
  await expect(first).toContainText('A save is pending, and newer edits have arrived.');
  await page.getByRole('button', { name: 'Older completion; newer edits remain', exact: true }).click();
  await expect(line).toHaveAttribute('data-save-state', 'dirty');
  await expect(first).toContainText('An older completion arrived while newer edits remain unsaved.');

  await page.getByRole('button', { name: 'Failed save', exact: true }).click();
  await expect(line).toHaveAttribute('data-save-state', 'failed');
  await expect(line).toHaveText('Couldn’t save draft');
  await page.getByRole('button', { name: 'Unknown outcome', exact: true }).click();
  await expect(line).toHaveAttribute('data-save-state', 'unknown');
  await expect(line).toHaveText('Save not confirmed');

  await page.getByRole('button', { name: 'Confirm meeting note' }).click();
  await expect(page.locator('#meeting-save-status [data-save-status]')).toHaveText('Draft saved');
  await expect(line).toHaveText('Save not confirmed');
  await page.getByRole('button', { name: 'Text saved; settings pending' }).click();
  await expect(line).toHaveText('Text saved');
  await expect(first).toContainText('Rules and collaborators still need saving.');
  const after = await first.evaluate(node => {
    const rect = node.getBoundingClientRect();
    return { x: rect.x + scrollX, y: rect.y + scrollY };
  });
  expect(after.x).toBe(before.x);
  expect(after.y).toBe(before.y);
  await page.getByRole('button', { name: 'Long supporting text' }).click();
  await expect(first).toContainText('longer rules, collaborators and other organization settings remain pending');
  const longTextPosition = await first.evaluate(node => {
    const rect = node.getBoundingClientRect();
    return { x: rect.x + scrollX, y: rect.y + scrollY };
  });
  expect(longTextPosition).toEqual(before);
  await expect(page.getByRole('button', { name: 'Long supporting text' })).toBeFocused();
  expect(await page.locator('html').evaluate(e => e.scrollWidth <= innerWidth)).toBeTruthy();

  await page.getByRole('button', { name: 'Dark', exact: true }).click();
  await expect(page.locator('[data-save-scheme]')).toHaveClass(/scheme-dark/);
  await expect(page.getByRole('button', { name: 'Dark', exact: true })).toHaveAttribute('aria-pressed', 'true');
  expect(await page.getByRole('button', { name: 'Dark', exact: true }).evaluate(node => getComputedStyle(node).textDecorationLine)).toContain('underline');
  await page.getByRole('button', { name: 'Inherited', exact: true }).click();
  await expect(page.locator('[data-save-scheme]')).toHaveClass(/scheme-default/);
  await expect(page.getByRole('button', { name: 'Dark', exact: true })).toHaveAttribute('aria-pressed', 'false');
  expect(await page.getByRole('button', { name: 'Dark', exact: true }).evaluate(node => getComputedStyle(node).textDecorationLine)).not.toContain('underline');
  await writeFile(info.outputPath('save-status-placement.json'), JSON.stringify({ width, before, after, longTextPosition, focusedControl: 'Long supporting text (before scheme selection)', statusRegions: 2 }, null, 2));
});
}

test('save status remains legible at zoom and forced colors', async ({ page }) => {
  await page.setViewportSize({ width: 390, height: 844 });
  await page.goto('/examples/save-status');
  await page.emulateMedia({ colorScheme: 'light' });
  const lightColor = await page.locator('#proposal-save-status [data-save-status]').evaluate(node => getComputedStyle(node).color);
  await page.emulateMedia({ colorScheme: 'dark' });
  const darkColor = await page.locator('#proposal-save-status [data-save-status]').evaluate(node => getComputedStyle(node).color);
  expect(darkColor).not.toBe(lightColor);
  await page.evaluate(() => { document.documentElement.style.zoom = '200%'; });
  expect(await page.locator('html').evaluate(e => e.scrollWidth <= innerWidth)).toBeTruthy();
  await page.emulateMedia({ forcedColors: 'active' });
  expect(await page.evaluate(() => matchMedia('(forced-colors: active)').matches)).toBeTruthy();
  const color = await page.locator('#proposal-save-status [data-save-status]').evaluate(node => getComputedStyle(node).color);
  expect(color).toBeTruthy();
  await expect(page.locator('#proposal-save-status [data-save-status]')).toBeVisible();
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
    if (name === 'page-header') {
      await expect(page.getByRole('heading', { level: 1, name: 'Title-only header' })).toBeVisible();
    } else {
      await expect(page.locator('h1')).toHaveText(name);
    }
  }
  await page.goto(process.env.GALLERY_URL + '/examples/save-status');
  await expect(page.locator('[data-save-controls]')).toHaveCount(3);
  for (const controls of await page.locator('[data-save-controls]').all()) await expect(controls).toBeHidden();
  await expect(page.locator('#proposal-text')).toBeVisible();
  await expect(page.locator('#proposal-save-status')).toContainText('Ready');
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

test('combobox keeps native form values authoritative through typing, selection, reset and submission', async ({ page }) => {
  await page.goto('/examples/combobox');
  await expect(page.locator('body')).toHaveAttribute('data-demo-ready', 'true');
  await expect(page.locator('#remote-scenario-control')).toBeVisible();
  const localInput = page.locator('#local-choice input');
  const remoteInput = page.locator('#remote-choice input');
  const localSelect = page.locator('#local-zone');
  const remoteSelect = page.locator('#remote-zone');
  expect(await page.locator('[data-combobox]').count()).toBe(2);
  expect(await page.locator('[role=combobox]').count()).toBe(2);
  expect(await page.locator('select[name]').evaluateAll(nodes => nodes.map(node => node.name).sort())).toEqual(['interface_timezone', 'required_timezone']);
  expect(await localInput.getAttribute('name')).toBeNull();
  expect(await remoteInput.getAttribute('name')).toBeNull();
  expect(await remoteSelect.getAttribute('form')).toBe('combobox-form');

  const events = await page.evaluate(() => {
    const select = document.querySelector('#local-zone');
    window.comboboxEvents = { input: 0, change: 0, commit: 0 };
    select.addEventListener('input', () => window.comboboxEvents.input++);
    select.addEventListener('change', () => window.comboboxEvents.change++);
    document.querySelector('#local-choice').addEventListener('combobox-change', () => window.comboboxEvents.commit++);
    return [...new FormData(document.querySelector('#combobox-form'))];
  });
  expect(events).toEqual([['interface_timezone', ''], ['required_timezone', '']]);
  await localInput.fill('Makassar');
  expect(await localSelect.inputValue()).toBe('');
  expect(await page.evaluate(() => [...new FormData(document.querySelector('#combobox-form'))])).toEqual(events);
  expect(await page.evaluate(() => window.comboboxEvents)).toEqual({ input: 0, change: 0, commit: 0 });
  await localInput.press('ArrowDown');
  await localInput.press('Enter');
  expect(await localSelect.inputValue()).toBe('Asia/Makassar');
  expect(await page.evaluate(() => window.comboboxEvents)).toEqual({ input: 1, change: 1, commit: 1 });
  expect(await page.evaluate(() => [...new FormData(document.querySelector('#combobox-form'))])).toEqual([['interface_timezone', 'Asia/Makassar'], ['required_timezone', '']]);

  await page.locator('#combobox-reset').click();
  await expect(localInput).toHaveValue('Use organization default');
  expect(await localSelect.inputValue()).toBe('');
  await remoteInput.fill('Paris');
  await expect(page.locator('#remote-choice [role=option]')).toHaveText(['Paris']);
  await remoteInput.press('ArrowDown');
  await remoteInput.press('Enter');
  expect(await remoteSelect.inputValue()).toBe('Europe/Paris');

  await page.getByRole('button', { name: 'Submit choices' }).click();
  await page.waitForURL(/required_timezone=Europe%2FParis/);
  const submitted = new URL(page.url()).searchParams;
  expect([...submitted.keys()].sort()).toEqual(['interface_timezone', 'required_timezone']);
  expect(submitted.get('interface_timezone')).toBe('');
  expect(submitted.get('required_timezone')).toBe('Europe/Paris');
  await expect(page.locator('#demo-outcome')).toContainText('interface_timezone=""; required_timezone="Europe/Paris"');

  await localInput.fill('Paris');
  await localInput.press('Escape');
  expect(await localInput.inputValue()).toBe('Use organization default');
  await localInput.fill('Makassar');
  await localInput.press('ArrowDown');
  await localInput.press('Enter');
  await remoteInput.fill('Paris');
  await expect(page.locator('#remote-choice [role=option]')).toHaveText(['Paris']);
  await remoteInput.press('ArrowDown');
  await remoteInput.press('Enter');
  expect(await remoteSelect.inputValue()).toBe('Europe/Paris');
  await page.locator('#combobox-reset').click();
  await expect(localInput).toHaveValue('Use organization default');
  await expect(remoteInput).toHaveValue('Choose a timezone');
  expect(await localSelect.inputValue()).toBe('');
  expect(await remoteSelect.inputValue()).toBe('');
});

test('combobox required validation, fieldset disabling and external form ownership stay native', async ({ page }) => {
  await page.goto('/examples/combobox');
  await expect(page.locator('body')).toHaveAttribute('data-demo-ready', 'true');
  await page.getByRole('button', { name: 'Submit choices' }).click();
  expect(await page.locator('#remote-zone').evaluate(select => select.validity.valueMissing)).toBe(true);
  expect(await page.locator('#remote-choice input').evaluate(input => input === document.activeElement)).toBe(true);
  expect(await page.locator('#remote-choice input').getAttribute('aria-invalid')).toBe('true');
  await page.locator('#remote-scenario').selectOption('normal');
  const input = page.locator('#remote-choice input');
  await input.fill('Paris');
  await expect(page.locator('#remote-choice [role=option]')).toHaveText(['Paris']);
  await input.press('ArrowDown');
  await input.press('Escape');
  expect(await page.locator('#remote-zone').inputValue()).toBe('');
  await input.fill('Paris');
  await expect(page.locator('#remote-choice [role=option]')).toHaveText(['Paris']);
  await input.press('ArrowDown');
  await input.press('Enter');
  expect(await page.locator('#remote-zone').evaluate(select => select.checkValidity())).toBe(true);

  const fieldset = await page.evaluate(async () => {
    const { createCombobox } = await import('/assets/combobox.js');
    const form = document.createElement('form');
    form.innerHTML = '<fieldset disabled><legend><div data-combobox><label for="legend-zone">Legend timezone</label><select id="legend-zone" name="legend" required><option value="UTC">UTC</option></select></div></legend><fieldset><div data-combobox><label for="nested-zone">Nested timezone</label><select id="nested-zone" name="nested"><option value="UTC">UTC</option></select></div></fieldset></fieldset>';
    document.body.append(form);
    const controls = [...form.querySelectorAll('[data-combobox]')].map(root => createCombobox({ root }));
    await new Promise(resolve => setTimeout(resolve, 0));
    const before = { values: [...new FormData(form)], disabled: [...form.querySelectorAll('input')].map(input => input.disabled) };
    form.querySelector('fieldset').disabled = false;
    await new Promise(resolve => setTimeout(resolve, 0));
    const after = { values: [...new FormData(form)], disabled: [...form.querySelectorAll('input')].map(input => input.disabled) };
    controls.forEach(control => control.dispose());
    form.remove();
    return { before, after };
  });
  expect(fieldset.before).toEqual({ values: [['legend', 'UTC']], disabled: [false, true] });
  expect(fieldset.after).toEqual({ values: [['legend', 'UTC'], ['nested', 'UTC']], disabled: [false, false] });
});

test('combobox remote results reject stale replies, report failures, and retain selected values', async ({ page }) => {
  await page.clock.install();
  await page.goto('/examples/combobox');
  await expect(page.locator('body')).toHaveAttribute('data-demo-ready', 'true');
  const input = page.locator('#remote-choice input');
  const options = page.locator('#remote-choice [role=option]');
  await page.evaluate(() => {
    window.comboboxQueries = [];
    window.comboboxReplies = [];
    document.querySelector('#remote-choice').addEventListener('combobox-query', event => {
      window.comboboxQueries.push(event.detail.query);
    });
    const setTimeout = window.setTimeout.bind(window);
    window.setTimeout = (callback, delay, ...args) => {
      if (delay !== 900 && delay !== 100) return setTimeout(callback, delay, ...args);
      return setTimeout(() => {
        window.comboboxReplies.push(delay);
        callback(...args);
      }, delay);
    };
  });
  await page.locator('#remote-scenario').selectOption('normal');
  await input.fill('Choose');
  await expect(page.locator('#remote-choice [role=status]')).toHaveText('No results.');
  expect(await options.count()).toBe(0);
  await page.locator('#remote-scenario').selectOption('out-of-order');
  await page.clock.pauseAt(new Date(Date.now() + 60_000));
  await page.evaluate(() => { window.comboboxQueries = []; });
  await input.fill('pa');
  await page.clock.runFor(300);
  await expect.poll(() => page.evaluate(() => window.comboboxQueries)).toEqual(['pa']);
  await input.fill('par');
  await page.clock.runFor(300);
  await expect.poll(() => page.evaluate(() => window.comboboxQueries)).toEqual(['pa', 'par']);
  await page.clock.runFor(100);
  await expect(page.locator('#remote-choice [role=option]')).toHaveText(['Paris']);
  await page.clock.runFor(500);
  await expect.poll(() => page.evaluate(() => window.comboboxReplies)).toEqual([100, 900]);
  await expect(page.locator('#remote-choice [role=option]')).toHaveText(['Paris']);
  await page.clock.resume();

  await page.locator('#remote-scenario').selectOption('error');
  await input.fill('fail');
  await expect(page.locator('#remote-choice [role=status]')).toContainText('Search failed');
  await page.locator('#remote-scenario').selectOption('empty');
  await input.fill('zz');
  await expect(page.locator('#remote-choice [role=status]')).toHaveText('No results.');

  await page.locator('#remote-scenario').selectOption('delayed');
  await input.fill('Manila');
  await expect(page.locator('#remote-choice [role=status]')).toHaveText('Searching…');
  await expect(page.locator('#remote-choice [role=option]')).toHaveText(['Manila']);

  await page.locator('#remote-scenario').selectOption('normal');
  await input.fill('Makassar');
  await page.locator('#remote-choice [role=option]').filter({ hasText: 'Makassar' }).click();
  expect(await page.locator('#remote-zone').inputValue()).toBe('Asia/Makassar');
  await input.fill('Paris');
  await input.press('Escape');
  await expect(input).toHaveValue('Makassar');
  expect(await page.locator('#remote-zone').inputValue()).toBe('Asia/Makassar');
});

test('combobox local DOM updates preserve query and no-JavaScript fallback submits natively', async ({ page, browser }) => {
  await page.goto('/examples/combobox');
  await expect(page.locator('body')).toHaveAttribute('data-demo-ready', 'true');
  const input = page.locator('#local-choice input');
  await input.fill('sao');
  await input.evaluate(element => element.setSelectionRange(1, 2));
  await page.evaluate(() => {
    document.querySelector('#local-zone').add(new Option('São Tomé', 'Africa/Sao_Tome'));
    document.body.append(document.createElement('p'));
  });
  await expect(page.locator('#local-choice [role=option]')).toHaveCount(2);
  expect(await input.inputValue()).toBe('sao');
  expect(await input.evaluate(element => [element.selectionStart, element.selectionEnd])).toEqual([1, 2]);
  expect(await input.evaluate(element => element === document.activeElement)).toBe(true);

  const context = await browser.newContext({ javaScriptEnabled: false, viewport: { width: 390, height: 844 } });
  try {
    const native = await context.newPage();
    await native.goto('/examples/combobox');
    expect(await native.locator('[role=combobox]').count()).toBe(0);
    await expect(native.locator('#local-zone')).toBeVisible();
    await expect(native.locator('#remote-zone')).toBeVisible();
    await expect(native.locator('#remote-scenario-control')).toBeHidden();
    await native.locator('#remote-zone').selectOption('Europe/Paris');
    await native.getByRole('button', { name: 'Submit choices' }).click();
    await native.waitForURL(/required_timezone=Europe%2FParis/);
    const params = new URL(native.url()).searchParams;
    expect(params.get('interface_timezone')).toBe('');
    expect(params.get('required_timezone')).toBe('Europe/Paris');
    expect([...params.keys()].sort()).toEqual(['interface_timezone', 'required_timezone']);
  } finally {
    await context.close();
  }
});
