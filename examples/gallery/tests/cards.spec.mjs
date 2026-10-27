import { test, expect } from '@playwright/test';
import { readFile, writeFile } from 'node:fs/promises';

const producer = '1c52666402c0ec2de4ed64736ea679a8a8573838';
const referenceFixture = new URL('../../../component/testdata/livewires/card-summary-reference.html', import.meta.url);

const measureLines = node => {
  const style = getComputedStyle(node);
  return { lines: Math.round(node.clientHeight / parseFloat(style.lineHeight)), clamped: node.scrollHeight > node.clientHeight + 1, clampValue: style.webkitLineClamp, overflow: style.overflow };
};

for (const width of [375, 1280]) {
  for (const colorScheme of ['light', 'dark']) {
    test(`cards ${width}px ${colorScheme}: provenance, layout, clamps, focus and header rules`, async ({ browser, request }, info) => {
      const context = await browser.newContext({ viewport: { width, height: 900 }, colorScheme });
      const page = await context.newPage();
      const errors = [];
      page.on('pageerror', error => errors.push(error.message));
      page.on('requestfailed', req => errors.push(req.url() + ' ' + req.failure().errorText));
      page.on('response', res => { if (res.status() >= 400) errors.push(res.url() + ' HTTP ' + res.status()); });
      const health = await (await request.get('/healthz')).json();
      expect(health.source.assetRevision).toBe(producer);
      await page.goto('/examples/cards');
      await expect(page.locator('h1')).toHaveText('cards');
      await expect(page.locator('#source-build')).toHaveAttribute('data-source-head', health.source.commit);
      await expect(page.locator('#source-build')).toHaveAttribute('data-source-fingerprint', health.source.fingerprint);
      expect(await page.evaluate(() => matchMedia('(prefers-color-scheme: dark)').matches)).toBe(colorScheme === 'dark');
      expect(await page.locator('html').evaluate(e => e.scrollWidth <= innerWidth)).toBeTruthy();

      // Existing cards keep their markup and header rule when no footer follows.
      const existing = page.locator('[data-cards=existing]');
      await expect(existing.locator('> .card')).toHaveCount(4);
      expect(await existing.locator('.card:not(.card--media):not(.box) > .header').evaluate(h => getComputedStyle(h).borderBlockEndWidth)).toBe('1px');

      // Header rule is removed only when a footer follows the header directly.
      const rules = await page.evaluate(() => {
        const width = el => getComputedStyle(el).borderBlockEndWidth;
        const top = el => getComputedStyle(el).borderBlockStartWidth;
        const bodyCase = document.querySelector('[data-case=header-body-footer]');
        const directCase = document.querySelector('[data-case=header-footer]');
        const summary = document.querySelector('[data-cards=members] > li');
        return {
          headerBeforeBody: width(bodyCase.querySelector(':scope > .header')),
          footerAfterBody: top(bodyCase.querySelector(':scope > .footer')),
          headerBeforeFooter: width(directCase.querySelector(':scope > .header')),
          footerAfterHeader: top(directCase.querySelector(':scope > .footer')),
          summaryHeader: width(summary.querySelector(':scope > .header')),
          summaryFooter: top(summary.querySelector(':scope > .footer')),
          directSibling: directCase.querySelector(':scope > .header').nextElementSibling === directCase.querySelector(':scope > .footer'),
        };
      });
      expect(rules).toEqual({ headerBeforeBody: '1px', footerAfterBody: '1px', headerBeforeFooter: '0px', footerAfterHeader: '1px', summaryHeader: '0px', summaryFooter: '1px', directSibling: true });

      // Equal heights and aligned footers in every row of every summary grid.
      const grids = await page.evaluate(() => [...document.querySelectorAll('ul[data-cards]')].map(grid => ({
        name: grid.dataset.cards,
        rows: [...grid.querySelectorAll(':scope > li')].map(li => {
          const box = li.getBoundingClientRect();
          const footer = li.querySelector(':scope > .footer').getBoundingClientRect();
          return { top: Math.round(box.top + scrollY), height: Math.round(box.height * 100) / 100, footerTop: Math.round((footer.top - box.top) * 100) / 100, footerBottom: Math.round((footer.bottom - box.top) * 100) / 100 };
        }),
      })));
      const rowReport = [];
      for (const { name, rows } of grids) {
        const byRow = new Map();
        for (const card of rows) byRow.set(card.top, [...(byRow.get(card.top) ?? []), card]);
        for (const [top, cards] of byRow) {
          rowReport.push({ grid: name, top, cards });
          for (const card of cards) {
            expect(card.height, `${name} row ${top} equal heights`).toBe(cards[0].height);
            expect(card.footerBottom, `${name} row ${top} footer bottom at card bottom`).toBe(card.height);
            // Footers sit on the card bottom. Their tops also line up when the footer content is the same shape.
            if (name === 'author' || name === 'author-members') expect(card.footerTop, `${name} row ${top} aligned footers`).toBe(cards[0].footerTop);
          }
        }
      }
      if (width === 1280) expect(rowReport.some(r => r.cards.length > 1)).toBe(true);

      // Three- and four-line limits, with the CSS custom properties honoured on request.
      const limits = page.locator('[data-cards=limits] > li');
      const longHeading = await limits.nth(0).locator('.card-heading').evaluate(measureLines);
      const longDescription = await limits.nth(0).locator('.description').evaluate(measureLines);
      const customHeading = await limits.nth(1).locator('.card-heading').evaluate(measureLines);
      const customDescription = await limits.nth(1).locator('.description').evaluate(measureLines);
      expect(longHeading).toMatchObject({ lines: 3, clamped: true, clampValue: '3' });
      expect(longDescription).toMatchObject({ lines: 4, clamped: true, clampValue: '4' });
      expect(customHeading).toMatchObject({ lines: 2, clamped: true, clampValue: '2' });
      expect(customDescription).toMatchObject({ lines: 2, clamped: true, clampValue: '2' });
      // The full text stays in the DOM.
      expect((await limits.nth(0).locator('.card-heading').textContent()).trim()).toMatch(/^Festival submissions 2027: .* until a deadline is two weeks away$/);
      expect((await limits.nth(0).locator('.description').textContent()).trim()).toMatch(/Anything that needs a decision goes through a proposal first\.$/);

      // The heading link's focus ring is drawn inside the clamped text.
      const link = limits.nth(0).locator('.card-heading a');
      await link.focus();
      await page.keyboard.press('Shift+Tab');
      await page.keyboard.press('Tab');
      await expect(link).toBeFocused();
      const focus = await link.evaluate(a => {
        const heading = a.closest('.card-heading').getBoundingClientRect();
        const style = getComputedStyle(a);
        const rects = [...a.getClientRects()].map(r => ({ top: r.top, bottom: r.bottom, left: r.left, right: r.right }));
        return { outlineStyle: style.outlineStyle, outlineWidth: style.outlineWidth, outlineOffset: style.outlineOffset, matchesFocusVisible: a.matches(':focus-visible'), inside: rects.filter(r => r.top >= heading.top - 0.5 && r.bottom <= heading.bottom + 0.5 && r.left >= heading.left - 0.5 && r.right <= heading.right + 0.5).length, rects: rects.length, visibleRects: rects.filter(r => r.bottom > heading.top && r.top < heading.bottom).length };
      });
      expect(focus.matchesFocusVisible).toBe(true);
      expect(focus.outlineStyle).not.toBe('none');
      expect(parseFloat(focus.outlineWidth)).toBeGreaterThan(0);
      expect(focus.outlineOffset).toBe('-2px');
      expect(focus.inside).toBe(focus.visibleRects);
      const ringShot = info.outputPath(`cards-focus-${width}-${colorScheme}.png`);
      await limits.nth(0).locator('.header').screenshot({ path: ringShot });
      await info.attach('Focus ring inside the clamp', { path: ringShot, contentType: 'image/png' });

      // Accessibility structure.
      const a11y = await page.evaluate(() => {
        const cards = [...document.querySelectorAll('ul[data-cards] > li')];
        return {
          groupsHidden: cards.flatMap(li => [...li.querySelectorAll('.avatar-group')]).every(g => g.getAttribute('aria-hidden') === 'true'),
          avatarsNamed: [...document.querySelectorAll('ul[data-cards] .avatar')].filter(a => a.hasAttribute('role') || a.hasAttribute('aria-label') || a.querySelector('img:not([alt=""])')).length,
          interactivePerHeading: cards.map(li => li.querySelector('.card-heading').querySelectorAll('a, button, input, select, textarea, [tabindex]').length),
          timesWithDatetime: [...document.querySelectorAll('ul[data-cards] time')].every(t => /^\d{4}-\d{2}-\d{2}$/.test(t.getAttribute('datetime') ?? '')),
          authorNames: [...document.querySelectorAll('[data-cards=author] .footer .truncate')].map(n => n.textContent.trim()),
          hiddenText: [...document.querySelectorAll('ul[data-cards] .visually-hidden')].map(n => n.textContent.trim()).filter((v, i, all) => all.indexOf(v) === i),
        };
      });
      expect(a11y.groupsHidden).toBe(true);
      expect(a11y.avatarsNamed).toBe(0);
      expect(a11y.interactivePerHeading.every(n => n === 1)).toBe(true);
      expect(a11y.timesWithDatetime).toBe(true);
      expect(a11y.authorNames).toContain('By Wiremu Te Rangi-Hēnare-Whitehead');
      expect(a11y.hiddenText).toEqual(expect.arrayContaining(['By', '5 members', 'conversations, 1 unread', 'Last updated', 'members']));
      const truncated = await page.locator('[data-cards=author] .footer .truncate').last().evaluate(n => ({ overflow: getComputedStyle(n).textOverflow, clipped: n.scrollWidth > n.clientWidth, textLength: n.textContent.length }));
      expect(truncated.overflow).toBe('ellipsis');
      expect(errors).toEqual([]);

      const shot = info.outputPath(`cards-${width}-${colorScheme}-${health.source.commit.slice(0, 12)}-${health.source.fingerprint.slice(0, 12)}.png`);
      await page.screenshot({ path: shot, fullPage: true });
      await info.attach('Source-bound screenshot', { path: shot, contentType: 'image/png' });
      await writeFile(info.outputPath(`cards-metrics-${width}-${colorScheme}.json`), JSON.stringify({ producer, width, colorScheme, source: health.source, rules, rowReport, longHeading, longDescription, customHeading, customDescription, focus, truncated, a11y, scrollWidth: await page.evaluate(() => document.documentElement.scrollWidth) }, null, 2));
      await context.close();
    });

    test(`cards ${width}px ${colorScheme}: producer reference comparison`, async ({ browser, request }, info) => {
      const context = await browser.newContext({ viewport: { width, height: 900 }, colorScheme, javaScriptEnabled: false });
      const reference = await context.newPage();
      const target = await context.newPage();
      const css = await (await request.get('/assets/main.css')).text();
      const html = await readFile(referenceFixture, 'utf8');
      await reference.setContent(`<!doctype html><html lang="en"><head><meta name="viewport" content="width=device-width, initial-scale=1"></head><body class="scheme-default"><main class="section center stack w-auto">${html}</main></body></html>`);
      await reference.evaluate(text => { const style = document.createElement('style'); style.textContent = text; document.head.append(style); }, css);
      await target.goto('/examples/cards');
      const read = grid => [...grid.querySelectorAll(':scope > li')].map(li => {
        const pick = (el, props) => el ? Object.fromEntries(props.map(p => [p, getComputedStyle(el)[p]])) : null;
        const text = el => el?.textContent.replace(/\s+/g, ' ').trim() ?? null;
        return {
          heading: text(li.querySelector('.card-heading')), footerText: li.querySelector('.footer').innerText.replace(/\s+/g, ' ').trim(),
          card: pick(li, ['display', 'flexDirection', 'borderTopWidth', 'borderTopLeftRadius', 'backgroundColor', 'minInlineSize']),
          header: pick(li.querySelector('.header'), ['paddingTop', 'paddingLeft', 'borderBottomWidth', 'rowGap', 'display']),
          headingStyle: pick(li.querySelector('.card-heading'), ['fontSize', 'fontWeight', 'lineHeight', 'overflow', 'webkitLineClamp', 'display']),
          description: pick(li.querySelector('.description'), ['fontSize', 'lineHeight', 'webkitLineClamp', 'overflow']),
          footer: pick(li.querySelector('.footer'), ['display', 'flexWrap', 'justifyContent', 'alignItems', 'columnGap', 'paddingTop', 'paddingLeft', 'borderTopWidth', 'borderTopColor', 'fontSize']),
          avatar: pick(li.querySelector('.footer .avatar'), ['inlineSize', 'blockSize', 'borderTopLeftRadius']),
          time: pick(li.querySelector('.footer time'), ['fontSize']),
          badge: pick(li.querySelector('.badge'), ['backgroundColor', 'color', 'fontSize']),
        };
      });
      // Gallery grids whose copy equals a producer grid: author = 2nd, author-members = 3rd.
      const pairs = [['author', 1], ['author-members', 2]];
      const results = [];
      for (const [name, index] of pairs) {
        const actual = await target.locator(`ul[data-cards=${name}]`).evaluate(read);
        const source = await reference.locator('ul.grid').nth(index).evaluate(read);
        expect(actual).toEqual(source);
        results.push({ name, cards: actual.length });
      }
      // Members grid: cards 1, 3 and 4 share producer copy (card 2 uses a longer heading to hit the clamp).
      const members = await target.locator('ul[data-cards=members]').evaluate(read);
      const membersSource = await reference.locator('ul.grid').nth(0).evaluate(read);
      for (const index of [0, 2, 3]) expect(members[index]).toEqual(membersSource[index]);
      results.push({ name: 'members', cards: 3 });
      const first = info.outputPath(`producer-reference-${width}-${colorScheme}.png`);
      const second = info.outputPath(`gallery-${width}-${colorScheme}.png`);
      await reference.screenshot({ path: first, fullPage: true });
      await target.screenshot({ path: second, fullPage: true });
      await info.attach('Producer reference', { path: first, contentType: 'image/png' });
      await info.attach('Gallery', { path: second, contentType: 'image/png' });
      if (process.env.PRODUCER_MANUAL) {
        const manual = await context.newPage();
        const manualHtml = (await readFile(process.env.PRODUCER_MANUAL, 'utf8')).replace(/<script\b[^>]*>[\s\S]*?<\/script>/g, '').replace(/<html-include\b[^>]*><\/html-include>/g, '').replace(/<link[^>]*stylesheet[^>]*>/g, '');
        await manual.setContent(manualHtml);
        await manual.evaluate(text => { const style = document.createElement('style'); style.textContent = text; document.head.append(style); }, css);
        const manualShot = info.outputPath(`producer-manual-${width}-${colorScheme}.png`);
        await manual.screenshot({ path: manualShot, fullPage: true });
        await info.attach('Producer manual', { path: manualShot, contentType: 'image/png' });
        expect(await manual.evaluate(() => document.documentElement.scrollWidth <= innerWidth + 1)).toBeTruthy();
      }
      await writeFile(info.outputPath(`producer-parity-${width}-${colorScheme}.json`), JSON.stringify({ producer, width, colorScheme, results, adaptations: ['Producer site scripts, includes and stylesheets removed; the gallery locked main.css injected.', 'Members grid card 2 uses a longer heading than the producer to hit the 3-line limit.', 'Computed-style and text parity, not pixel parity: the gallery column width differs from the producer page.'] }, null, 2));
      await context.close();
    });
  }
}

test('cards render without JavaScript', async ({ browser }) => {
  const context = await browser.newContext({ javaScriptEnabled: false });
  const page = await context.newPage();
  await page.goto('/examples/cards');
  await expect(page.locator('h1')).toHaveText('cards');
  await expect(page.locator('ul[data-cards=members] > li')).toHaveCount(4);
  await context.close();
});
