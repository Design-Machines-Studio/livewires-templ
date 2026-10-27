import { test, expect } from '@playwright/test';
import { writeFile } from 'node:fs/promises';

const revision = '1c52666402c0ec2de4ed64736ea679a8a8573838';
const route = '/examples/mention-suggestions';
const editorOf = (page, id = 'mention') => page.locator(`#${id}-editor`);
const popupOf = (page, id = 'mention') => page.locator(`#${id}-popup`);

async function open(page, query = '@', id = 'mention') {
  const editor = editorOf(page, id);
  await editor.fill('');
  await editor.fill(query);
  await expect(popupOf(page, id)).toHaveAttribute('data-state', 'ready');
  await expect(editor).toBeFocused();
  return editor;
}
async function active(page, id = 'mention') {
  const editor = editorOf(page, id), row = popupOf(page, id).locator('[data-active]');
  await expect(row).toHaveCount(1);
  await expect(row).toHaveAttribute('aria-selected', 'true');
  await expect(row).toHaveClass('rounded-sm scheme-subtle');
  await expect(editor).toHaveAttribute('aria-activedescendant', await row.getAttribute('id'));
  return row;
}
async function contained(page, id = 'mention') {
  const box = await popupOf(page, id).boundingBox();
  const bounds = await page.evaluate(() => ({ width: document.documentElement.clientWidth, height: innerHeight, pageWidth: document.documentElement.scrollWidth }));
  expect(box.x).toBeGreaterThanOrEqual(-1);
  expect(box.y).toBeGreaterThanOrEqual(-1);
  expect(box.x + box.width).toBeLessThanOrEqual(bounds.width + 1);
  expect(box.y + box.height).toBeLessThanOrEqual(bounds.height + 1);
  expect(bounds.pageWidth).toBeLessThanOrEqual(bounds.width);
}

for (const width of [320, 375, 768, 1280]) {
  for (const scheme of ['light', 'dark', 'inherited']) {
    test(`mentions ${width}px ${scheme}: anatomy, focus, geometry and RTL`, async ({ page, request }, info) => {
      await page.setViewportSize({ width, height: 900 });
      await page.emulateMedia({ colorScheme: scheme === 'dark' ? 'dark' : 'light' });
      const errors = [];
      page.on('pageerror', e => errors.push(e.message));
      const source = (await (await request.get('/healthz')).json()).source;
      expect(source.assetRevision).toBe(revision);
      await page.goto(route);
      if (scheme === 'inherited') await page.evaluate(() => {
        document.body.className = 'scheme-dark';
        // The component's explicit scheme-default must inherit the current dark environment.
        document.querySelector('main').classList.add('scheme-subtle');
      });
      await expect(page.locator('#source-build')).toHaveAttribute('data-source-head', source.commit);
      await expect(page.locator('#source-build')).toHaveAttribute('data-source-fingerprint', source.fingerprint);
      const ids = await page.locator('[id]').evaluateAll(nodes => nodes.map(n => n.id));
      expect(new Set(ids).size).toBe(ids.length);
      for (const state of ['closed', 'loading', 'ready', 'empty', 'error']) {
        const root = page.locator(`#static-${state}-root`);
        await expect(root).toHaveAttribute('data-state', state);
        await expect(root.locator('p[role=status][aria-atomic=true]')).toHaveCount(1);
        if (state !== 'ready') await expect(root.locator('ul')).toHaveAttribute('hidden', '');
      }
      await expect(page.locator('#static-stale')).toHaveAttribute('aria-busy', 'true');
      await expect(page.locator('#static-stale [data-active]')).toHaveCount(0);
      await expect(page.locator('#static-ready .avatar img')).toHaveAttribute('alt', '');
      await expect(page.locator('#static-ready .avatar').first()).toHaveAttribute('aria-hidden', 'true');
      await expect(page.locator('#static-long [data-mention-name].truncate')).toHaveText('Maria-Fernanda Oyelaran-Castellanos de la Vega');

      await open(page);
      await active(page);
      await contained(page);
      const image = `mentions-${width}-${scheme}-${source.commit}-${source.fingerprint}.png`;
      await page.screenshot({ path: info.outputPath(image), fullPage: true });
      await writeFile(info.outputPath('source.json'), JSON.stringify({ source, width, scheme, aria: await popupOf(page).ariaSnapshot() }, null, 2));
      await editorOf(page).press('Escape');
      await open(page, '@', 'mention-rtl');
      await expect(popupOf(page, 'mention-rtl')).toHaveAttribute('dir', 'rtl');
      await active(page, 'mention-rtl');
      await contained(page, 'mention-rtl');
      await page.screenshot({ path: info.outputPath(`rtl-${width}-${scheme}.png`) });
      expect(errors).toEqual([]);
    });
  }
}

test('mentions real keyboard: wrap, skip disabled, select, dismiss, modifiers and IME boundary', async ({ page }) => {
  await page.goto(route);
  const editor = await open(page), root = popupOf(page);
  await expect(editor).toHaveAttribute('role', 'textbox');
  await expect(editor).toHaveAttribute('aria-multiline', 'true');
  await expect(editor).toHaveAttribute('aria-controls', 'mention-list');
  await expect(editor).not.toHaveAttribute('aria-expanded');
  await editor.press('ArrowDown'); await expect(await active(page)).toContainText('Bo Lindqvist');
  await editor.press('ArrowDown'); await expect(await active(page)).toContainText('Sam Rivera');
  await editor.press('ArrowDown'); await expect(await active(page)).toContainText('Di Ng');
  await editor.press('ArrowUp'); await expect(await active(page)).toContainText('Sam Rivera');
  await editor.press('Home'); // Unhandled editor key; selectionchange may close the token.
  await open(page);
  const before = await page.evaluate(() => ({ x: scrollX, y: scrollY }));
  await editor.press('ArrowUp'); await expect(await active(page)).toContainText('Rosa Jiménez');
  expect(await root.locator('ul').evaluate(n => n.scrollTop)).toBeGreaterThan(0);
  expect(await page.evaluate(() => ({ x: scrollX, y: scrollY }))).toEqual(before);
  const inside = await root.evaluate(r => {
    const list = r.querySelector('ul').getBoundingClientRect(), row = r.querySelector('[data-active]').getBoundingClientRect();
    return row.top >= list.top - 1 && row.bottom <= list.bottom + 1;
  });
  expect(inside).toBeTruthy();
  await editor.press('ArrowDown'); await expect(await active(page)).toContainText('Ada Okafor');
  await editor.press('Enter');
  await expect(editor).toContainText('@Ada Okafor');
  await expect(root).toHaveAttribute('data-state', 'closed');
  await expect(editor).not.toHaveAttribute('aria-activedescendant');
  await expect(editor).toBeFocused();
  await expect(page.locator('[data-mention-log]').first()).toContainText('Selected Ada');
  await open(page, '@sam'); await editor.press('Tab');
  await expect(editor).toContainText('@Sam Rivera'); await expect(editor).toBeFocused();
  for (const key of ['Shift+Enter', 'Control+Enter', 'Meta+Enter', 'Alt+Enter']) {
    await open(page, '@ada'); await editor.press(key);
    await expect(root).toHaveAttribute('data-state', 'closed');
    await expect(editor).toHaveText('@ada');
    await expect(page.locator('[data-mention-log]').first()).not.toContainText('Simulated send');
  }
  await open(page, '@a'); await editor.press('Escape');
  await editor.press('d'); await page.waitForTimeout(200);
  await expect(root).toHaveAttribute('data-state', 'closed');
  await editor.press('Space'); await editor.press('@');
  await expect(root).toHaveAttribute('data-state', 'ready');
  // This explicitly synthetic composition check proves event guards, not real IME acceptance.
  const guarded = await editor.evaluate(e => {
    const event = new KeyboardEvent('keydown', { key: 'Enter', isComposing: true, bubbles: true, cancelable: true });
    e.dispatchEvent(event); return !event.defaultPrevented;
  });
  expect(guarded).toBeTruthy(); await expect(root).toHaveAttribute('data-state', 'ready');
  await editor.press('Shift+Tab'); await expect(root).toHaveAttribute('data-state', 'closed'); await expect(editor).not.toBeFocused();
});

test('mentions state fragments: mounted status, delayed/stale/late replies, empty/error, no repeated count writes', async ({ page }) => {
  await page.goto(route);
  const editor = editorOf(page), root = popupOf(page), status = root.locator('[role=status]');
  await root.evaluate(r => { window.originalMentionStatus = r.querySelector('[role=status]'); });
  await open(page, '@fa');
  await status.evaluate(s => {
    window.mentionWrites = [];
    window.mentionObserver = new MutationObserver(() => window.mentionWrites.push(s.textContent));
    window.mentionObserver.observe(s, { childList: true, subtree: true, characterData: true });
  });
  await editor.press('t'); await page.waitForTimeout(220);
  await editor.press('i'); await page.waitForTimeout(220);
  expect(await page.evaluate(() => window.mentionWrites)).toEqual([]);
  expect(await status.evaluate(s => s === window.originalMentionStatus)).toBeTruthy();
  await page.selectOption('#mention-scenario', 'slow');
  await open(page, '@');
  await editor.press('s');
  await expect(editor).not.toHaveAttribute('aria-activedescendant');
  await expect(root.locator('ul')).toHaveAttribute('aria-busy', 'true');
  await expect(root).toHaveAttribute('data-state', 'loading');
  await expect(status).toHaveText('Searching people…');
  await expect(root.locator('[data-active]')).toHaveCount(0);
  await editor.press('a'); await editor.press('m');
  await expect(root).toHaveAttribute('data-state', 'ready');
  await expect(status).toHaveText('2 people found');
  await page.waitForTimeout(1700);
  await expect(root.locator('[role=option]')).toHaveCount(2);
  await expect(root.locator('ul')).not.toHaveAttribute('aria-busy');
  await page.selectOption('#mention-scenario', 'quick');
  await editor.fill('@zz'); await expect(root).toHaveAttribute('data-state', 'empty');
  await expect(status).toHaveText('No matching people'); await expect(editor).not.toHaveAttribute('aria-activedescendant');
  await editor.press('Enter'); await expect(root).toHaveAttribute('data-state', 'closed');
  await expect(page.locator('[data-mention-log]').first()).not.toContainText('Simulated send');
  await page.selectOption('#mention-scenario', 'error');
  await editor.fill('@ada'); await expect(root).toHaveAttribute('data-state', 'error');
  await expect(status).toHaveText('People could not be loaded. Keep typing to try again.');
  await expect(editor).not.toHaveAttribute('aria-activedescendant');
  await editor.press('Tab'); await expect(root).toHaveAttribute('data-state', 'closed'); await expect(editor).not.toBeFocused();
  expect(await status.evaluate(s => s === window.originalMentionStatus)).toBeTruthy();
});

test('mentions pointer, hover and touch selection retain editor focus; disabled click does nothing', async ({ browser }) => {
  const context = await browser.newContext({ viewport: { width: 375, height: 812 }, hasTouch: true });
  const page = await context.newPage(); await page.goto(process.env.GALLERY_URL + route);
  await open(page, '@sam');
  const editor = editorOf(page), root = popupOf(page);
  const previous = await editor.getAttribute('aria-activedescendant');
  await root.locator('[aria-disabled=true]').hover();
  expect(await editor.getAttribute('aria-activedescendant')).toBe(previous);
  await root.locator('[aria-disabled=true]').click({ force: true }); await expect(editor).toBeFocused();
  await expect(root).toHaveAttribute('data-state', 'ready');
  await root.locator('[role=option]').first().tap();
  await expect(editor).toBeFocused(); await expect(editor).toContainText('@Sam Rivera');
  await open(page, '@di'); await root.locator('[role=option]').click();
  await expect(editor).toBeFocused(); await expect(editor).toContainText('@Di Ng');
  await context.close();
});

for (const scheme of ['light', 'dark', 'inherited']) {
  test(`mentions 200% CSS zoom ${scheme}, forced colors, reduced motion and flip`, async ({ page }, info) => {
    await page.setViewportSize({ width: 1280, height: 900 });
    await page.emulateMedia({ colorScheme: scheme === 'dark' ? 'dark' : 'light', reducedMotion: 'reduce' });
    await page.goto(route);
    await page.evaluate(scheme => {
      document.body.style.zoom = '2';
      if (scheme === 'inherited') { document.body.className = 'scheme-dark'; document.querySelector('main').classList.add('scheme-subtle'); }
    }, scheme);
    await open(page, '@sam'); await contained(page);
    await page.screenshot({ path: info.outputPath(`zoom-200-${scheme}.png`) });
    await editorOf(page).press('Escape');
    await page.evaluate(() => { document.body.style.zoom = ''; });
    await page.emulateMedia({ forcedColors: 'active' });
    await open(page); await contained(page);
    const color = await (await active(page)).evaluate(n => ({ bg: getComputedStyle(n).backgroundColor, color: getComputedStyle(n).color }));
    expect(color.bg).not.toBe(color.color);
    await page.screenshot({ path: info.outputPath(`forced-colors-${scheme}.png`) });
    await editorOf(page).press('Escape');
    // A clipped scroll container gives the caret more space above than below.
    await editorOf(page).evaluate(editor => {
      const form = editor.closest('form');
      form.style.cssText = 'height:300px;overflow:auto';
      const spacer = document.createElement('div'); spacer.style.height = '700px';
      editor.before(spacer); editor.scrollIntoView({ block: 'end' });
    });
    await open(page);
    await expect(popupOf(page)).toHaveAttribute('data-placement', 'above');
    await contained(page);
  });
}
