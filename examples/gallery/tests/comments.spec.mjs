import { test, expect } from '@playwright/test';
import { writeFile } from 'node:fs/promises';

for (const width of [1440, 1280, 640, 375, 320]) {
  for (const scheme of ['light', 'dark']) {
    test(`comments ${width}px ${scheme}: contract, wrapping and media bounds`, async ({ page, request }, info) => {
      await page.setViewportSize({ width, height: 1000 });
      await page.emulateMedia({ colorScheme: scheme });
      await page.goto('/examples/comments');
      await expect(page.locator('.comment--layout')).toHaveCount(8);
      await expect(page.locator('.comment-avatar img')).toHaveCount(4);
      const checks = await page.evaluate(() => {
        const comments = [...document.querySelectorAll('.comment--layout')];
        const anchored = comments.filter(c => c.id);
        const rect = n => n.getBoundingClientRect();
        const overlaps = comments.filter(c => {
          const a = c.querySelector('.comment-author'), t = c.querySelector('.comment-permalink, .comment-timestamp');
          if (!a || !t) return false;
          const ar = rect(a), tr = rect(t);
          return ar.left < tr.right && tr.left < ar.right && ar.top < tr.bottom && tr.top < ar.bottom;
        }).map(c => c.id);
        const inherited = document.querySelector('#comment-inherited-scheme');
        return {
          count: comments.length,
          duplicateIDs: anchored.map(c => c.id).filter((id, i, all) => all.indexOf(id) !== i),
          missingAnchors: [...document.querySelectorAll('.comment--layout a[href^="#"]')].filter(a => !document.getElementById(a.hash.slice(1))).map(a => a.hash),
          badTimes: [...document.querySelectorAll('.comment--layout time')].filter(t => !t.dateTime || Number.isNaN(Date.parse(t.dateTime))).length,
          badAvatarAlt: [...document.querySelectorAll('.comment-avatar img')].filter(img => img.getAttribute('alt') !== '').length,
          unnamedAvatarLinks: [...document.querySelectorAll('a.comment-avatar')].filter(a => !a.getAttribute('aria-label')).length,
          badNesting: document.querySelectorAll('.comment--layout a a, .comment--layout a button, .comment--layout button a, .comment--layout button button').length,
          misplacedActions: [...document.querySelectorAll('.comment-actions a, .comment-actions button')].filter(a => !a.closest('.comment-footer') || a.closest('.comment-meta')).length,
          liveRegions: comments.filter(c => c.closest('[aria-live], [role="status"], [role="log"]')).length,
          pageOverflow: document.documentElement.scrollWidth > innerWidth + 1,
          commentOverflow: comments.filter(c => c.scrollWidth > c.clientWidth + 1 || rect(c).right > innerWidth + 1).map(c => c.id),
          mediaOverflow: [...document.querySelectorAll('.comment-body img')].filter(img => rect(img).width > rect(img.parentElement).width + 1).length,
          imageLoaded: [...document.querySelectorAll('.comment-body img')].every(img => img.complete && img.naturalWidth > 0),
          overlaps,
          inheritedColor: getComputedStyle(inherited).color === getComputedStyle(inherited.parentElement).color,
          inheritedBackground: getComputedStyle(inherited).backgroundColor,
          code: [...document.querySelectorAll('.comment-body pre')].map(pre => ({ scroll: pre.scrollWidth > pre.clientWidth, overflow: getComputedStyle(pre).overflowX })),
          classes: comments.map(c => ({ id: c.id, class: c.className, children: [...c.children].map(n => ({ tag: n.tagName, class: n.className })) })),
        };
      });
      for (const name of ['badTimes', 'badAvatarAlt', 'unnamedAvatarLinks', 'badNesting', 'misplacedActions', 'liveRegions', 'mediaOverflow']) expect(checks[name], name).toBe(0);
      for (const name of ['duplicateIDs', 'missingAnchors', 'commentOverflow', 'overlaps']) expect(checks[name], name).toEqual([]);
      expect(checks.pageOverflow).toBe(false);
      expect(checks.imageLoaded).toBe(true);
      expect(checks.inheritedColor).toBe(true);
      expect(checks.inheritedBackground).toBe('rgba(0, 0, 0, 0)');
      if (width === 320) expect(checks.code.some(c => c.scroll && c.overflow === 'auto')).toBe(true);
      const source = await (await request.get('/healthz')).json();
      await writeFile(info.outputPath('comments-dom.json'), JSON.stringify({ width, scheme, checks, source, accessibility: await page.locator('main').ariaSnapshot(), screenReaderTested: false }, null, 2));
      await page.screenshot({ path: info.outputPath(`comments-${width}-${scheme}.png`), fullPage: true });
    });
  }
}

test('comments: keyboard follows author, permalink and separate footer actions', async ({ page }, info) => {
  await page.goto('/examples/comments');
  const reply = page.locator('#chat-reply-816');
  const avatar = reply.locator('a.comment-avatar');
  await avatar.focus();
  const order = [];
  for (let i = 0; i < 6; i++) {
    order.push(await page.evaluate(() => ({ tag: document.activeElement.tagName, text: document.activeElement.textContent.trim(), class: document.activeElement.className, name: document.activeElement.getAttribute('aria-label') })));
    expect(await page.evaluate(() => document.activeElement.matches(':focus-visible'))).toBe(true);
    const style = await page.evaluate(() => ({ outline: getComputedStyle(document.activeElement).outlineStyle, shadow: getComputedStyle(document.activeElement).boxShadow }));
    expect(style.outline !== 'none' || style.shadow !== 'none').toBe(true);
    await page.keyboard.press('Tab');
  }
  expect(order.map(x => x.class)).toEqual(['comment-avatar', 'author font-bold text-fg underline', 'comment-permalink text-muted text-xs', '', '', 'button button--small']);
  // Rich body link and focusable code precede the footer in reading/tab order.
  await expect(reply.getByRole('link', { name: 'Open linked comment' })).toBeFocused();
  await page.keyboard.press('Enter');
  await expect(page).toHaveURL(/#comment-long-values$/);
  await writeFile(info.outputPath('keyboard.json'), JSON.stringify({ order, screenReaderTested: false }, null, 2));
});

test('comments: no-JavaScript rendering and 200% CSS zoom reflow', async ({ browser }, info) => {
  const context = await browser.newContext({ javaScriptEnabled: false, viewport: { width: 375, height: 900 } });
  const page = await context.newPage();
  await page.goto('/examples/comments');
  await expect(page.locator('#post-815 time')).toHaveText('Mar 15 at 10:05');
  await expect(page.locator('#chat-reply-816 .comment-body li')).toHaveCount(2);
  await page.getByRole('link', { name: 'Reply to Sam' }).click();
  await expect(page).toHaveURL(/#post-815$/);
  await page.screenshot({ path: info.outputPath('comments-no-js.png'), fullPage: true });
  await context.close();

  const zoomContext = await browser.newContext({ viewport: { width: 1280, height: 1000 } });
  const zoomPage = await zoomContext.newPage();
  await zoomPage.goto('/examples/comments');
  await zoomPage.evaluate(() => { document.documentElement.style.zoom = '2'; });
  expect(await zoomPage.evaluate(() => document.documentElement.scrollWidth <= innerWidth + 1)).toBe(true);
  await zoomPage.screenshot({ path: info.outputPath('comments-css-zoom-200.png'), fullPage: true });
  await writeFile(info.outputPath('zoom-scope.json'), JSON.stringify({ cssZoom: '200%', reflowViewport: 640, nativeBrowserZoomTested: false, screenReaderTested: false }, null, 2));
  await zoomContext.close();
});
