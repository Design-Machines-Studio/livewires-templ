import { test, expect } from '@playwright/test';
import { readFile, writeFile } from 'node:fs/promises';

const producer = 'aad7e3843898745a2f29833393b20896a204b90e';
const producerReference = new URL('../../../component/testdata/livewires/people-picker-reference.html', import.meta.url);
const schemes = ['light', 'dark', 'no-preference'];

for (const width of [320, 375, 480, 1024]) {
  for (const colorScheme of schemes) {
    test(`people picker ${width}px ${colorScheme}: presentation matches pinned Live Wires reference`, async ({ browser }, info) => {
      const context = await browser.newContext({ javaScriptEnabled: false, viewport: { width, height: 1000 }, colorScheme });
      const target = await context.newPage();
      const reference = await context.newPage();
      await target.goto('/examples/people-picker');

      const source = await readFile(producerReference, 'utf8');
      await reference.setContent(source.replace(/<script\b[^>]*>[\s\S]*?<\/script>/gi, '').replace(/<html-include\b[^>]*><\/html-include>/gi, ''));
      await reference.evaluate((href) => {
        const link = document.createElement('link');
        link.rel = 'stylesheet';
        link.href = href;
        link.dataset.testCss = 'true';
        document.head.append(link);
        for (const image of document.querySelectorAll('.people-picker img')) image.src = '/assets/comment-image.jpg';
      }, new URL('/assets/main.css', process.env.GALLERY_URL).href);
      await reference.locator('link[data-test-css]').evaluate(link => new Promise(resolve => {
        if (link.sheet) resolve();
        else link.addEventListener('load', resolve, { once: true });
      }));

      // Give the reference the same container width as the gallery card so widths are comparable.
      const targetBoxWidth = await target.locator('form[data-demo-picker="space"]').evaluate(node => node.parentElement.getBoundingClientRect().width);
      await reference.locator('fieldset.people-picker').first().evaluate((node, boxWidth) => {
        const box = node.closest('.box');
        box.style.boxSizing = 'border-box';
        box.style.inlineSize = `${boxWidth}px`;
        box.style.maxInlineSize = `${boxWidth}px`;
      }, targetBoxWidth);
      const actualPicker = target.locator('form[data-demo-picker="space"] fieldset.people-picker');
      const sourcePicker = reference.locator('fieldset.people-picker').first();
      await expect(actualPicker).toBeVisible();
      await expect(sourcePicker).toBeVisible();
      const inspect = node => {
        const rect = element => {
          const box = element.getBoundingClientRect();
          return { x: Math.round(box.x), y: Math.round(box.y), width: Math.round(box.width), height: Math.round(box.height) };
        };
        const style = element => {
          const value = getComputedStyle(element);
          return { display: value.display, width: value.width, fontSize: value.fontSize, lineHeight: value.lineHeight, gap: value.gap, flexWrap: value.flexWrap, padding: value.padding, color: value.color, background: value.backgroundColor };
        };
        const groups = node.querySelector('.people-picker-groups');
        const table = node.querySelector('.people-picker-table');
        const row = node.querySelector('tbody tr');
        const identity = row.querySelector('.people-picker-person');
        const trailing = row.querySelector('.people-picker-trailing');
        const name = row.querySelector('.people-picker-name');
        return {
          fieldset: { rect: rect(node), style: style(node) },
          groups: { rect: rect(groups), style: style(groups), labels: [...groups.querySelectorAll('button')].map(button => button.textContent.trim()) },
          table: { rect: rect(table), style: style(table), headers: [...table.querySelectorAll('thead th')].map(th => th.textContent.trim()) },
          identity: { rect: rect(identity), style: style(identity), name: name.textContent.trim() },
          trailing: { rect: rect(trailing), style: style(trailing), dataLabel: trailing.dataset.th },
          overflow: document.documentElement.scrollWidth > innerWidth + 1,
        };
      };
      const sourceLayout = await sourcePicker.evaluate(inspect);
      const targetLayout = await actualPicker.evaluate(inspect);
      expect(targetLayout.overflow).toBe(false);
      expect(sourceLayout.overflow).toBe(false);
      expect(targetLayout.groups.style).toEqual(sourceLayout.groups.style);
      expect(targetLayout.groups.labels).toEqual(sourceLayout.groups.labels);
      expect(targetLayout.table.style).toEqual(sourceLayout.table.style);
      expect(targetLayout.table.headers).toEqual(sourceLayout.table.headers);
      expect(targetLayout.identity.style).toEqual(sourceLayout.identity.style);
      expect(targetLayout.trailing.style).toEqual(sourceLayout.trailing.style);
      expect(targetLayout.trailing.dataLabel).toBe('Organizer');
      expect(targetLayout.identity.name).toBe(sourceLayout.identity.name);
      expect(targetLayout.identity.name).toBe('Ada OkaforOrganizer, Board');

      const targetImage = info.outputPath(`people-picker-templ-${width}-${colorScheme}.png`);
      const referenceImage = info.outputPath(`people-picker-producer-${width}-${colorScheme}.png`);
      await actualPicker.screenshot({ path: targetImage });
      await sourcePicker.screenshot({ path: referenceImage });
      await info.attach('Templ composition', { path: targetImage, contentType: 'image/png' });
      await info.attach('Pinned producer composition', { path: referenceImage, contentType: 'image/png' });
      await writeFile(info.outputPath(`people-picker-parity-${width}-${colorScheme}.json`), JSON.stringify({ producer, width, colorScheme, sourceLayout, targetLayout, domAccessibilitySnapshot: await actualPicker.ariaSnapshot(), actualScreenReaderTested: false }, null, 2));
      await context.close();
    });
  }
}

test('people picker keeps overlapping circles, direct selection, exclusions and locked rows separate', async ({ page }) => {
  await page.setViewportSize({ width: 1024, height: 1000 });
  await page.goto('/examples/people-picker');
  const form = page.locator('form[data-demo-picker="space"]');
  const ada = form.locator('#space-ada');
  const bo = form.locator('#space-bo');
  const cy = form.locator('#space-cy');
  const di = form.locator('#space-di');
  const ed = form.locator('#space-ed');
  const board = form.getByRole('button', { name: 'Board', exact: true });
  const officers = form.getByRole('button', { name: 'Officers', exact: true });
  const clear = form.getByRole('button', { name: 'Clear selection' });
  const status = form.locator('[role="status"]');
  await expect(form.locator('[role="status"]')).toHaveAttribute('aria-atomic', 'true');
  await expect(status).toHaveText('');
  await expect(ada).toBeChecked();
  await expect(bo).toBeChecked();
  await expect(cy).not.toBeChecked();
  await expect(di).toBeChecked();
  await expect(di).toBeDisabled();

  await officers.click();
  await expect(officers).toHaveAttribute('aria-pressed', 'true');
  await expect(ed).toBeChecked();
  await expect(status).toHaveText('officers circle selected.');
  await expect(officers).toBeFocused();
  await board.click();
  await expect(board).toHaveAttribute('aria-pressed', 'false');
  await expect(bo).not.toBeChecked();
  await expect(ada).toBeChecked();
  await expect(cy).not.toBeChecked();
  await expect(ed).toBeChecked();

  await form.locator('label[for="space-ada"] .avatar').click();
  await expect(ada).not.toBeChecked();
  await expect(ada).toBeEnabled();
  await form.locator('label[for="space-ada"] .people-picker-name').click();
  await expect(ada).toBeChecked();
  await expect(ada).toBeEnabled();
  await ada.focus();
  await page.keyboard.press('Space');
  await expect(ada).not.toBeChecked();
  await page.keyboard.press('Space');
  await expect(ada).toBeChecked();

  await clear.click();
  await expect(ada).not.toBeChecked();
  await expect(bo).not.toBeChecked();
  await expect(cy).not.toBeChecked();
  await expect(ed).not.toBeChecked();
  await expect(di).toBeChecked();
  await expect(di).toBeDisabled();
  await expect(form.locator('[data-demo-organizer]')).toHaveCount(5);
  for (const organizer of await form.locator('[data-demo-organizer]').all()) {
    await expect(organizer).not.toBeChecked();
    await expect(organizer).toBeDisabled();
  }
  await expect(status).toHaveText('Selections cleared. Locked people remain selected.');
  await expect(clear).toBeFocused();

  const other = page.locator('form[data-demo-picker="plain-space"]');
  await expect(other.locator('#plain-space-ada')).toBeChecked();
  expect(await page.locator('[id]').evaluateAll(nodes => nodes.map(node => node.id).filter((id, index, all) => all.indexOf(id) !== index))).toEqual([]);
  const role = page.getByRole('textbox', { name: 'Role in this circle for Ada Okafor' });
  await role.fill('Facilitator');
  await expect(role).toHaveValue('Facilitator');
  await expect(page.locator('#circle-ada')).toBeChecked();
});

test('people picker keeps selection and role inputs native when JavaScript is disabled', async ({ browser }) => {
  const context = await browser.newContext({ javaScriptEnabled: false, viewport: { width: 375, height: 900 } });
  const page = await context.newPage();
  await page.goto('/examples/people-picker');
  const form = page.locator('form[data-demo-picker="space"]');
  await expect(form.getByRole('button', { name: 'Board', exact: true })).toBeDisabled();
  await expect(form.getByRole('button', { name: 'Clear selection' })).toBeDisabled();
  const checkbox = form.locator('#space-bo');
  await checkbox.focus();
  await page.keyboard.press('Space');
  await expect(checkbox).not.toBeChecked();
  const role = page.getByRole('textbox', { name: 'Role in this circle for Ada Okafor' });
  await role.fill('Secretary');
  await expect(role).toHaveValue('Secretary');
  expect(await form.locator('[role="status"]').count()).toBe(1);
  expect(await page.locator('html').evaluate(node => node.scrollWidth <= innerWidth + 1)).toBe(true);
  await context.close();
});

test('people picker exposes trailing labels, long accessible names, narrow wrapping, forced colors and RTL', async ({ page }) => {
  await page.setViewportSize({ width: 1024, height: 1000 });
  await page.goto('/examples/people-picker');
  const narrow = page.locator('form[data-demo-picker="narrow-truncate"]');
  const input = narrow.locator('#narrow-long');
  const labelID = await input.getAttribute('aria-labelledby');
  await expect(page.locator(`#${labelID}`)).toHaveText('Select Maximiliana Featherstonehaugh-Abernathy');
  await expect(narrow.locator('.people-picker-name > .block.truncate')).toContainText('Maximiliana Featherstonehaugh-Abernathy');
  await expect(narrow.locator('.people-picker-trailing')).toHaveAttribute('data-th', 'Organizer');
  await expect(narrow.locator('thead th').last()).toHaveText('Organizer');
  const narrowBox = narrow.locator('xpath=..');
  expect(await narrowBox.evaluate(node => getComputedStyle(node).maxWidth)).not.toBe('none');
  // Narrow means a constrained card, not a collapsed one: the name must not wrap one character per line.
  expect((await narrowBox.boundingBox()).width).toBeGreaterThanOrEqual(240);
  const wrapBox = page.locator('#narrow-wrap').locator('xpath=ancestor::div[contains(@class,"max-w-xs")]');
  expect((await wrapBox.boundingBox()).width).toBeGreaterThanOrEqual(240);
  expect((await page.locator('label[for="narrow-wrap"] .people-picker-name').boundingBox()).height).toBeLessThan(120);
  expect(await page.locator('html').evaluate(node => node.scrollWidth <= innerWidth + 1)).toBe(true);

  await page.emulateMedia({ forcedColors: 'active' });
  expect(await page.evaluate(() => matchMedia('(forced-colors: active)').matches)).toBe(true);
  const checkboxWrapper = narrow.locator('label.checkbox:has(#narrow-long)');
  await expect(checkboxWrapper).toBeVisible();
  const hit = await checkboxWrapper.boundingBox();
  expect(hit.width).toBeGreaterThanOrEqual(24);
  expect(hit.height).toBeGreaterThanOrEqual(24);
  await page.locator('html').evaluate(node => { node.dir = 'rtl'; });
  await expect(narrow.locator('legend')).toBeVisible();
  expect(await page.locator('html').evaluate(node => node.scrollWidth <= innerWidth + 1)).toBe(true);
  await page.screenshot({ path: test.info().outputPath('people-picker-forced-colors-rtl.png'), fullPage: true });
});
