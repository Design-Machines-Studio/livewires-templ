import { test, expect } from '@playwright/test';

for (const width of [320, 375, 1440]) {
  test(`hardware loaders at ${width}px`, async ({ page }, info) => {
    await page.setViewportSize({ width, height: 900 });
    await page.goto('/examples/loaders');
    await expect(page.locator('body')).toHaveAttribute('data-demo-ready', 'true');
    await expect(page.locator('.loader-panel').first()).toBeVisible();
    await expect(page.locator('[role=progressbar]')).toHaveCount(3);
    for (const [i, value] of ['0', '40', '100'].entries()) {
      await expect(page.locator('[role=progressbar]').nth(i)).toHaveAttribute('aria-valuenow', value);
    }
    expect(await page.locator('html').evaluate(e => e.scrollWidth <= innerWidth)).toBeTruthy();
    await expect(page.locator('#hardware-region [data-loading-indicator]')).toBeHidden();
    await expect(page.locator('.save-status[data-save-state=saving] svg')).toBeVisible();
    await expect(page.locator('.save-status[data-save-state=saved] svg')).toHaveCount(0);
    const button = page.getByRole('button', { name: 'Preview page overlay' });
    await button.focus();
    await page.keyboard.press('Enter');
    const overlay = page.locator('[data-gallery-overlay]');
    await expect(overlay).toBeVisible();
    await expect(overlay.locator('.visually-hidden')).toBeFocused();
    expect(await overlay.evaluate(e => e.parentElement.hasAttribute('hidden'))).toBe(false);
    await expect(page.locator('main')).toHaveAttribute('inert', '');
    await page.keyboard.press('Escape');
    await expect(overlay).toBeHidden();
    await expect(button).toBeFocused();
    await expect(page.locator('main')).not.toHaveAttribute('inert');
    await page.screenshot({ path: info.outputPath(`loaders-${width}.png`), fullPage: true });
  });
}

test('loader reduced motion and forced colors retain progress', async ({ page }) => {
  await page.emulateMedia({ reducedMotion: 'reduce', forcedColors: 'active' });
  await page.goto('/examples/loaders');
  const heads = await page.locator('.loader .head').evaluateAll(nodes => nodes.map(n => getComputedStyle(n).animationName));
  expect(heads.every(name => name === 'none')).toBeTruthy();
  await expect(page.locator('[role=progressbar]').nth(1)).toHaveAttribute('aria-valuenow', '40');
  await expect(page.locator('[role=progressbar]').nth(1)).toContainText('40%');
});

test('loader native output needs no JavaScript', async ({ browser }) => {
  const context = await browser.newContext({ javaScriptEnabled: false });
  const page = await context.newPage();
  await page.goto(process.env.GALLERY_URL + '/examples/loaders');
  await expect(page.locator('.loader-panel').first()).toBeVisible();
  await expect(page.locator('[data-gallery-overlay]')).toBeHidden();
  await expect(page.locator('[role=progressbar]').nth(1)).toContainText('40%');
  await context.close();
});
