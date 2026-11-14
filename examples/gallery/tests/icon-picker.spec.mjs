import { test, expect } from '@playwright/test';
import { writeFile } from 'node:fs/promises';

// Live Wires v2.29.0, the lock revision that supplies icon-picker.js and the Phosphor sprite.
const producer = 'd407d3edf652d6f4f754987a0ad9f48aaeb1faa4';
const route = '/examples/icon-picker';

const watchErrors = page => {
  const errors = [];
  page.on('pageerror', error => errors.push(error.message));
  page.on('requestfailed', req => errors.push(req.url() + ' ' + req.failure().errorText));
  page.on('response', res => { if (res.status() >= 400) errors.push(res.url() + ' HTTP ' + res.status()); });
  return errors;
};
const picker = (page, id) => ({
  root: page.locator(`#${id}-picker`),
  trigger: page.locator(`#${id}-picker .icon-picker-trigger`),
  panel: page.locator(`#${id}-panel`),
  search: page.locator(`#${id}-search`),
  status: page.locator(`#${id}-panel .icon-picker-status`),
  tile: name => page.locator(`#${id}-panel .icon-picker-body > section:not(.icon-picker-recent) input[value="${name}"]`),
});
const isOpen = locator => locator.evaluate(el => el.matches(':popover-open'));
const focusedValue = page => page.evaluate(() => document.activeElement?.value ?? null);

async function ready(page) {
  await page.goto(route);
  await expect(page.locator('body')).toHaveAttribute('data-demo-ready', 'true');
}

test('icon-picker: provenance, assets, no id collisions and unenhanced native section', async ({ page, request }) => {
  const errors = watchErrors(page);
  const health = await (await request.get('/healthz')).json();
  expect(health.source.assetRevision).toBe(producer);
  for (const file of ['icon-picker.js', 'icons/phosphor-bold.svg', 'icons/phosphor-catalog.json', 'icons/LICENSE']) {
    const asset = health.assets.find(a => a.file === file);
    expect(asset, file).toBeTruthy();
    expect((await request.get(asset.serve)).ok()).toBe(true);
  }
  expect((await request.get('/icons/phosphor-bold.svg')).headers()['content-type']).toBe('image/svg+xml');
  await ready(page);
  const dupes = await page.evaluate(() => {
    const ids = [...document.querySelectorAll('[id]')].map(n => n.id);
    return ids.filter((id, i) => ids.indexOf(id) !== i);
  });
  expect(dupes).toEqual([]);
  const names = await page.evaluate(() => [...document.querySelectorAll('[data-icon-picker]')].map(root => [...new Set([...root.querySelectorAll('input[type=radio]')].map(r => r.name))]));
  expect(names.every(group => group.length === 1)).toBe(true);
  expect(new Set(names.flat()).size).toBe(names.length);
  await expect(page.locator('[data-gallery-enhanced] [data-icon-picker][data-icon-picker-ready]')).toHaveCount(5);
  await expect(page.locator('[data-gallery-native] [data-icon-picker-ready]')).toHaveCount(0);
  await expect(page.locator('[data-gallery-native] .icon-picker-trigger[aria-expanded]')).toHaveCount(0);
  // The sprite resolves: a referenced symbol has geometry.
  const box = await page.locator('#space-picker .icon-picker-current svg').boundingBox();
  expect(box.width).toBeGreaterThan(8);
  expect(errors).toEqual([]);
});

test('icon-picker enhanced: open, focus order, arrows, search, empty state', async ({ page }, info) => {
  await page.setViewportSize({ width: 1280, height: 900 });
  const errors = watchErrors(page);
  await ready(page);
  const p = picker(page, 'space');
  await expect(p.trigger).toHaveAttribute('aria-expanded', 'false');
  await expect(p.trigger).toHaveAttribute('aria-label', 'Change space icon, Users three');
  await p.trigger.click();
  await expect.poll(() => isOpen(p.panel)).toBe(true);
  await expect(p.trigger).toHaveAttribute('aria-expanded', 'true');
  await expect(p.search).toBeVisible();
  await expect(p.search).toBeFocused();
  await expect(p.tile('users-three')).toBeChecked();
  await expect(p.tile('users-three').locator('..')).toBeInViewport();
  await expect(p.status).toHaveText('48 icons');
  await page.keyboard.press('Tab');
  expect(await focusedValue(page)).toBe('users-three');
  await page.screenshot({ path: info.outputPath('enhanced-open-1280.png') });

  await page.keyboard.press('ArrowRight');
  expect(await focusedValue(page)).toBe('handshake');
  const before = await page.evaluate(() => document.activeElement.closest('label').getBoundingClientRect().top);
  await page.keyboard.press('ArrowDown');
  const after = await page.evaluate(() => ({ top: document.activeElement.closest('label').getBoundingClientRect().top, value: document.activeElement.value }));
  expect(after.top).toBeGreaterThan(before);
  await page.keyboard.press('Home');
  expect(await focusedValue(page)).toBe('users-three');
  // Arrow movement does not choose; the original value stays checked.
  await expect(p.tile('users-three')).toBeChecked();

  await p.search.fill('coin');
  await expect(p.status).toHaveText('2 icons');
  await expect(p.tile('coins').locator('..')).toBeVisible();
  await p.search.fill('money');
  await expect(p.tile('piggy-bank').locator('..')).toBeVisible();
  await expect(p.tile('money').locator('..')).toBeVisible();
  await expect(page.locator('#space-panel section[aria-labelledby="space-people"]')).toBeHidden();
  await p.search.fill('zzzz');
  await expect(p.status).toHaveText('No icons match “zzzz”');
  await page.screenshot({ path: info.outputPath('enhanced-empty-1280.png') });
  await writeFile(info.outputPath('enhanced-accessibility-tree.json'), JSON.stringify({ panel: await p.panel.ariaSnapshot(), note: 'Accessibility-tree snapshot only. No screen reader was run.', screenReaderTested: false }, null, 2));
  expect(errors).toEqual([]);
});

test('icon-picker enhanced: Escape and outside click close without change; choose and Remove update the trigger', async ({ page }) => {
  await page.setViewportSize({ width: 1280, height: 900 });
  await ready(page);
  const p = picker(page, 'space');
  await p.trigger.click();
  await page.keyboard.press('Tab');
  await page.keyboard.press('ArrowRight');
  await page.keyboard.press('Escape');
  await expect.poll(() => isOpen(p.panel)).toBe(false);
  await expect(p.trigger).toBeFocused();
  await expect(p.tile('users-three')).toBeChecked();
  await expect(p.trigger).toHaveAttribute('aria-label', 'Change space icon, Users three');

  await p.trigger.click();
  await page.mouse.click(5, 5);
  await expect.poll(() => isOpen(p.panel)).toBe(false);
  await expect(p.tile('users-three')).toBeChecked();

  await p.trigger.click();
  await p.search.fill('coins');
  await p.tile('coins').locator('..').click();
  await expect.poll(() => isOpen(p.panel)).toBe(false);
  await expect(p.trigger).toBeFocused();
  await expect(p.trigger).toHaveAttribute('aria-label', 'Change space icon, Coins');
  await expect(p.root.locator('.icon-picker-current use')).toHaveAttribute('href', '/icons/phosphor-bold.svg#coins');

  await p.trigger.click();
  await page.locator('#space-panel .icon-picker-remove').click();
  await expect.poll(() => isOpen(p.panel)).toBe(false);
  await expect(p.trigger).toBeFocused();
  await expect(p.trigger).toHaveAttribute('aria-label', 'Choose space icon');
  await expect(p.root.locator('.icon-picker-placeholder')).toBeVisible();
  await expect(p.root.locator('.icon-picker-current')).toBeHidden();

  await p.trigger.click();
  await p.search.fill('gavel');
  await page.keyboard.press('Enter');
  await expect(p.trigger).toHaveAttribute('aria-label', 'Change space icon, Gavel');
  await page.locator('#icon-enhanced-form button[type=submit]').click();
  await expect(page).toHaveURL(/space_icon=gavel/);
  await expect(page.locator('#icon-picker-result')).toContainText('Space icon (enhanced) (space_icon): gavel.');
  await expect(picker(page, 'space').tile('gavel')).toBeChecked();
});

test('icon-picker enhanced: bottom sheet below 30em', async ({ page }, info) => {
  await page.setViewportSize({ width: 320, height: 700 });
  await ready(page);
  const p = picker(page, 'circle');
  await p.trigger.click();
  await expect.poll(() => isOpen(p.panel)).toBe(true);
  const rect = await p.panel.boundingBox();
  expect(Math.round(rect.x)).toBeLessThanOrEqual(1);
  expect(Math.round(rect.width)).toBeGreaterThanOrEqual(318);
  expect(Math.round(rect.y + rect.height)).toBeGreaterThanOrEqual(698);
  await page.screenshot({ path: info.outputPath('bottom-sheet-320.png') });
});

test('icon-picker without JavaScript: popover opens, radios submit, trigger waits for save, rejection shows error', async ({ browser }, info) => {
  const context = await browser.newContext({ javaScriptEnabled: false, viewport: { width: 1280, height: 900 } });
  const page = await context.newPage();
  await page.goto(route);
  const p = picker(page, 'native-circle');
  await expect(p.trigger).toHaveAttribute('aria-label', 'Choose circle icon');
  await expect(p.trigger).not.toHaveAttribute('aria-expanded', /.*/);
  await p.trigger.click();
  await expect.poll(() => isOpen(p.panel)).toBe(true);
  await expect(page.locator('#native-circle-panel .icon-picker-search')).toBeHidden();
  await expect(page.locator('#native-circle-panel .icon-picker-recent')).toBeHidden();
  await expect(p.status).toBeHidden();
  const rect = await p.panel.boundingBox();
  const centre = rect.x + rect.width / 2;
  expect(Math.abs(centre - 640)).toBeLessThan(20);
  await page.screenshot({ path: info.outputPath('native-open-1280.png') });
  await p.tile('gavel').locator('..').click();
  await expect(p.tile('gavel')).toBeChecked();
  // The trigger does not update until the form is saved.
  await expect(p.trigger).toHaveAttribute('aria-label', 'Choose circle icon');
  await p.trigger.click();
  await expect.poll(() => isOpen(p.panel)).toBe(false);
  await page.locator('#icon-native-form button[type=submit]').click();
  await expect(page).toHaveURL(/native_circle_icon=gavel/);
  await expect(page.locator('#icon-picker-result')).toContainText('Circle icon (without JavaScript) (native_circle_icon): gavel.');
  await expect(picker(page, 'native-circle').trigger).toHaveAttribute('aria-label', 'Change circle icon, Gavel');

  await page.goto(route + '?native_circle_icon=jeep&space_icon=rocket');
  await expect(page.locator('#icon-picker-error')).toContainText('Circle icon (without JavaScript) (native_circle_icon): "jeep" is not one of this picker’s icons. Nothing was saved.');
  await expect(page.locator('#icon-picker-error')).toContainText('Space icon (enhanced) (space_icon): "rocket"');
  await expect(picker(page, 'native-circle').trigger).toHaveAttribute('aria-label', 'Choose circle icon');
  await page.screenshot({ path: info.outputPath('native-rejected-1280.png'), fullPage: true });
  await context.close();
});

for (const width of [1280, 640, 375, 320]) {
  for (const colorScheme of ['light', 'dark']) {
    test(`icon-picker ${width}px ${colorScheme}: reflow and open panels in every scheme`, async ({ page }, info) => {
      await page.setViewportSize({ width, height: 900 });
      await page.emulateMedia({ colorScheme });
      const errors = watchErrors(page);
      await ready(page);
      expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth + 1)).toBe(true);
      await page.screenshot({ path: info.outputPath(`page-${width}-${colorScheme}.png`), fullPage: true });
      for (const id of ['space', 'scheme-dark', 'scheme-light']) {
        const p = picker(page, id);
        await p.trigger.scrollIntoViewIfNeeded();
        await p.trigger.click();
        await expect.poll(() => isOpen(p.panel)).toBe(true);
        const rect = await p.panel.boundingBox();
        expect(rect.x).toBeGreaterThanOrEqual(-1);
        expect(rect.x + rect.width).toBeLessThanOrEqual(width + 1);
        await page.screenshot({ path: info.outputPath(`open-${id}-${width}-${colorScheme}.png`) });
        await page.keyboard.press('Escape');
        await expect.poll(() => isOpen(p.panel)).toBe(false);
      }
      expect(errors).toEqual([]);
    });
  }
}

test('icon-picker 200% zoom: CSS zoom and a 2x-density half-width viewport reflow', async ({ browser }, info) => {
  for (const colorScheme of ['light', 'dark']) {
    const context = await browser.newContext({ viewport: { width: 640, height: 450 }, deviceScaleFactor: 2, colorScheme });
    const page = await context.newPage();
    await ready(page);
    expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth + 1)).toBe(true);
    const p = picker(page, 'space');
    await p.trigger.click();
    await expect.poll(() => isOpen(p.panel)).toBe(true);
    const rect = await p.panel.boundingBox();
    expect(rect.x + rect.width).toBeLessThanOrEqual(641);
    await page.screenshot({ path: info.outputPath(`zoom-equivalent-${colorScheme}.png`) });
    await context.close();
  }
  const context = await browser.newContext({ viewport: { width: 1280, height: 900 } });
  const page = await context.newPage();
  await ready(page);
  await page.evaluate(() => { document.documentElement.style.zoom = '2'; });
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= innerWidth + 1)).toBe(true);
  await page.screenshot({ path: info.outputPath('css-zoom-200.png'), fullPage: true });
  await writeFile(info.outputPath('zoom-scope.json'), JSON.stringify({ cssZoom: '200%', equivalentViewport: '640x450 at 2x', nativeBrowserZoomTested: false, screenReaderTested: false }, null, 2));
  await context.close();
});

test('icon-picker forced colors: focus ring visible; checked state recorded for producer review', async ({ page, browserName }, info) => {
  test.skip(browserName !== 'chromium', 'forced-colors emulation is Chromium-only in Playwright');
  await page.setViewportSize({ width: 1280, height: 900 });
  await page.emulateMedia({ forcedColors: 'active' });
  await ready(page);
  const p = picker(page, 'space');
  await p.trigger.click();
  await page.keyboard.press('Tab');
  const focus = await page.evaluate(() => {
    const label = document.activeElement.closest('label');
    const style = getComputedStyle(label);
    return { outline: style.outlineStyle, outlineWidth: style.outlineWidth, border: style.borderStyle, value: document.activeElement.value };
  });
  await page.screenshot({ path: info.outputPath('forced-colors-open.png') });
  // Move focus away so the checked tile is judged on its own.
  await page.keyboard.press('ArrowRight');
  const tileStyle = value => page.locator(`#space-panel input[value="${value}"]`).evaluate(input => {
    const style = getComputedStyle(input.closest('label'));
    return { background: style.backgroundColor, border: style.borderColor, borderWidth: style.borderWidth, outline: style.outlineStyle, color: style.color };
  });
  const checked = await tileStyle('users-three');
  const unchecked = await tileStyle('baby');
  await page.screenshot({ path: info.outputPath('forced-colors-checked-unfocused.png') });
  await writeFile(info.outputPath('forced-colors.json'), JSON.stringify({ focus, checked, unchecked, computedStylesDiffer: JSON.stringify(checked) !== JSON.stringify(unchecked), visualReview: 'Producer gap (Live Wires v2.29.0 icon-picker.css): with focus elsewhere the checked tile is visually identical to unchecked tiles. Its Canvas background matches the panel, and border and color are the same. See forced-colors-checked-unfocused.png.', note: 'Emulated forced-colors in Chromium; not a Windows High Contrast run.' }, null, 2));
  expect(focus.value).toBe('users-three');
  expect(focus.outline).not.toBe('none');
  expect(parseFloat(focus.outlineWidth)).toBeGreaterThan(0);
});
