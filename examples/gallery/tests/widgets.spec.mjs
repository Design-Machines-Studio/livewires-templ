import { test, expect } from '@playwright/test';
import { readFile, writeFile } from 'node:fs/promises';

const producer = 'aad7e3843898745a2f29833393b20896a204b90e';

for (const width of [1440, 768, 640, 375, 320]) {
  test(`widgets ${width}px: anatomy, inherited themes and complete wrapping`, async ({ page, request }, info) => {
    const errors = [];
    page.on('pageerror', e => errors.push(e.message));
    await page.setViewportSize({ width, height: 1000 });
    await page.goto('/examples/widgets');
    const observations = await page.evaluate(() => {
      const widgets = [...document.querySelectorAll('.widget')];
      const rect = n => n.getBoundingClientRect();
      const ids = [...document.querySelectorAll('[id]')].map(n => n.id);
      return {
        count: widgets.length,
        duplicateIDs: ids.filter((id, i) => ids.indexOf(id) !== i),
        unnamedSections: widgets.filter(n => n.tagName === 'SECTION' && !n.getAttribute('aria-label') && !n.getAttribute('aria-labelledby')).length,
        missingLabels: widgets.filter(n => n.hasAttribute('aria-labelledby') && !document.getElementById(n.getAttribute('aria-labelledby'))).length,
        unnamedBody: document.querySelector('#body-only').tagName,
        bodyHeader: !!document.querySelector('#body-only header'),
        live: widgets.filter(n => n.matches('[role="status"],[aria-live]') || n.querySelector('[role="status"],[aria-live]')).length,
        pageOverflow: document.documentElement.scrollWidth > innerWidth + 1,
        widgetOverflow: widgets.filter(n => rect(n).right > innerWidth + 1 || n.scrollWidth > n.clientWidth + 1).map(n => n.id || n.getAttribute('aria-labelledby')),
        badDefinitions: [...document.querySelectorAll('.widget dt,.widget dd')].filter(n => !(n.parentElement.matches('dl.widget-metric') || (n.parentElement.matches('div.widget-metric') && n.parentElement.parentElement.matches('dl.widget-metrics')))).length,
        statusFont: [...document.querySelectorAll('.widget-metric-status')].filter(n => n.classList.contains('font-metric') || n.classList.contains('widget-metric-value')).length,
        completeMoney: document.querySelector('#custom-b .widget-metric-value').textContent,
        plainBorder: getComputedStyle(document.querySelector('[aria-labelledby="empty-heading"] > header')).borderBottomWidth,
        table: {
          columns: document.querySelectorAll('[data-table-scroll] thead th[scope="col"]').length,
          rows: document.querySelectorAll('[data-table-scroll] tbody th[scope="row"]').length,
          cells: document.querySelectorAll('[data-table-scroll] tbody td').length,
          overflow: getComputedStyle(document.querySelector('[data-table-scroll]')).overflowX,
          scroll: document.querySelector('[data-table-scroll]').scrollWidth > document.querySelector('[data-table-scroll]').clientWidth,
          display: getComputedStyle(document.querySelector('[data-table-scroll] table')).display,
        },
        typography: ['custom-a', 'custom-b'].map(id => {
          const widget = document.getElementById(id), value = widget.querySelector('.widget-metric-value'), label = widget.querySelector('dt');
          const v = getComputedStyle(value), l = getComputedStyle(label);
          return { id, size: v.fontSize, weight: v.fontWeight, stretch: v.fontStretch, labelSize: l.fontSize, labelWeight: l.fontWeight, transform: v.transform, whiteSpace: v.whiteSpace, text: value.textContent };
        }),
        schemes: ['membership-heading','governance-heading','custom-a-heading'].map(id => {
          const n = document.getElementById(id).closest('.widget');
          return { id, color: getComputedStyle(n).color, background: getComputedStyle(n).backgroundColor, inherited: getComputedStyle(n).color === getComputedStyle(n.parentElement).color };
        }),
      };
    });
    expect(observations.count).toBe(15);
    for (const key of ['duplicateIDs','widgetOverflow']) expect(observations[key], key).toEqual([]);
    for (const key of ['unnamedSections','missingLabels','live','badDefinitions','statusFont']) expect(observations[key], key).toBe(0);
    expect(observations.unnamedBody).toBe('DIV');
    expect(observations.bodyHeader).toBe(false);
    expect(observations.pageOverflow).toBe(false);
    expect(observations.completeMoney).toBe('−Rp 1.234.567.890,50');
    expect(observations.plainBorder).toBe('0px');
    expect(observations.table.columns).toBe(4);
    expect(observations.table.rows).toBe(3);
    expect(observations.table.cells).toBe(9);
    expect(observations.table.overflow).toBe('auto');
    expect(observations.table.display).toBe('table');
    if (width <= 375) expect(observations.table.scroll).toBe(true);
    expect(observations.typography.map(n => n.weight)).toEqual(['650','500']);
    expect(observations.typography[0].labelSize).toBe(observations.typography[1].labelSize);
    expect(observations.typography[0].labelWeight).toBe(observations.typography[1].labelWeight);
    expect(observations.typography.every(n => n.transform === 'none' && n.whiteSpace !== 'nowrap')).toBe(true);
    expect(new Set(observations.schemes.map(n => n.background)).size).toBe(3);
    expect(observations.schemes.every(n => n.inherited)).toBe(true);
    expect(errors).toEqual([]);
    const source = await (await request.get('/healthz')).json();
    expect(source.source.assetRevision).toBe(producer);
    await writeFile(info.outputPath('widgets-dom.json'), JSON.stringify({ source, width, observations, accessibility: await page.locator('main').ariaSnapshot(), screenReaderTested: false },null,2));
    await page.screenshot({path:info.outputPath(`widgets-${width}.png`),fullPage:true});
  });
}

test('widgets: keyboard links, disclosure and table scrolling retain visible focus', async ({ page },info) => {
  await page.setViewportSize({width:375,height:900});
  await page.goto('/examples/widgets');
  const action = page.locator('.widget-action');
  await action.focus();
  await page.keyboard.press('Tab');
  await expect(page.getByRole('link',{name:"Review the community workshop's equipment sharing agreement",exact:true})).toBeFocused();
  const summary=page.locator('[aria-labelledby="chart-heading"] summary');
  await summary.focus();
  await page.keyboard.press('Enter');
  await expect(page.locator('[aria-labelledby="chart-heading"] details')).toHaveAttribute('open','');
  const region=page.locator('[data-table-scroll]');
  await region.focus();
  expect(await region.evaluate(n=>n.matches(':focus-visible'))).toBe(true);
  const focus=await region.evaluate(n=>({outline:getComputedStyle(n).outlineStyle,shadow:getComputedStyle(n).boxShadow}));
  expect(focus.outline!=='none'||focus.shadow!=='none').toBe(true);
  await page.keyboard.press('ArrowRight');
  await expect.poll(()=>region.evaluate(n=>n.scrollLeft)).toBeGreaterThan(0);
  await page.screenshot({path:info.outputPath('widgets-keyboard.png'),fullPage:true});
});

test('widgets: no JavaScript, CSS zoom and readable missing/static-font fallback',async({browser},info)=>{
  const native=await browser.newContext({javaScriptEnabled:false,viewport:{width:375,height:900}});
  const page=await native.newPage();await page.goto('/examples/widgets');
  await expect(page.locator('#membership .widget-metric-value').nth(1)).toHaveText('0');
  await expect(page.locator('.widget-metric-status').first()).toHaveText('Unavailable');
  await page.locator('[aria-labelledby="chart-heading"] summary').click();
  await expect(page.locator('[aria-labelledby="chart-heading"] table')).toBeVisible();
  await page.screenshot({path:info.outputPath('widgets-no-js.png'),fullPage:true});await native.close();
  const context=await browser.newContext({viewport:{width:1280,height:1000}}), zoom=await context.newPage();
  await zoom.goto('/examples/widgets');
  await zoom.evaluate(()=>{document.documentElement.style.zoom='2';});
  expect(await zoom.evaluate(()=>document.documentElement.scrollWidth<=innerWidth+1)).toBe(true);
  await zoom.screenshot({path:info.outputPath('widgets-css-zoom-200.png'),fullPage:true});
  const fallbacks=[];
  for(const family of ['"Missing Widget Font", serif','monospace']){
    fallbacks.push(await zoom.evaluate(family=>{
      document.documentElement.style.setProperty('--font-metric',family);
      const n=document.querySelector('#custom-b .widget-metric-value');
      const range=document.createRange();range.selectNodeContents(n);
      const r=range.getBoundingClientRect();
      return {family:getComputedStyle(n).fontFamily,text:n.textContent,rect:{width:r.width,height:r.height},overflow:n.scrollWidth>n.clientWidth+1};
    },family));
  }
  expect(fallbacks.every(n=>n.rect.width>0&&n.rect.height>0&&!n.overflow&&n.text==='−Rp 1.234.567.890,50')).toBe(true);
  await writeFile(info.outputPath('font-zoom-scope.json'),JSON.stringify({fallbacks,cssZoom:'200%',nativeBrowserZoomTested:false,screenReaderTested:false,variableFontGlyphWidthTested:false,reason:'No supported variable font is bundled or required; computed stretch alone is not glyph-width proof.'},null,2));
  await context.close();
});

for(const width of [1440,375,320]){
  test(`widgets ${width}px: exact published producer comparison`,async({browser,request},info)=>{
    const context=await browser.newContext({javaScriptEnabled:false,viewport:{width,height:1000}});
    const reference=await context.newPage(), target=await context.newPage();
    const html=await readFile(new URL('../../../component/testdata/livewires/widgets-reference.html',import.meta.url),'utf8');
    const css=await(await request.get('/assets/main.css')).text();
    await reference.setContent(html.replace(/<script\b[^>]*>[\s\S]*?<\/script>/g,'').replace(/<html-include\b[^>]*><\/html-include>/g,''));
    await reference.evaluate(css=>{const style=document.createElement('style');style.textContent=css;document.head.append(style);},css);
    await target.goto('/examples/widgets');
    const results=[];
    for(const id of ['membership-heading','governance-heading','empty-heading','unavailable-heading','chart-heading']){
      const inspect=n=>{
        const w=n.closest('.widget'), style=getComputedStyle(w);
        return {tag:w.tagName,classes:w.className,children:[...w.children].map(x=>({tag:x.tagName,class:x.className})),title:n.textContent,
          definitions:[...w.querySelectorAll('dt,dd')].map(x=>({tag:x.tagName,class:x.className,text:x.textContent})),
          message:w.querySelector('.widget-message-content')?.textContent.trim()??null,
          display:style.display,direction:style.flexDirection,padding:style.padding,gap:style.gap,
          values:[...w.querySelectorAll('.widget-metric-value')].map(x=>({size:getComputedStyle(x).fontSize,weight:getComputedStyle(x).fontWeight,leading:getComputedStyle(x).lineHeight,tracking:getComputedStyle(x).letterSpacing}))};
      };
      const source=await reference.locator('#'+id).evaluate(inspect), actual=await target.locator('#'+id).evaluate(inspect);
      // The gallery's static chart has no producer footer; the frame's optional
      // footer seam is proved by the other four paired examples and render tests.
      if(id==='chart-heading'){source.children=source.children.filter(x=>x.tag!=='FOOTER');}
      // Whitespace between caller paragraphs is immaterial; compare text per p.
      if(source.message) source.message=source.message.replace(/\s+/g,' ');
      if(actual.message) actual.message=actual.message.replace(/\s+/g,' ');
      expect(actual).toEqual(source);
      await reference.locator('#'+id).locator('..').locator('..').locator('..').screenshot({path:info.outputPath(`${id}-producer-${width}.png`)});
      await target.locator('#'+id).locator('..').locator('..').locator('..').screenshot({path:info.outputPath(`${id}-templ-${width}.png`)});
      results.push({id,source,actual});
    }
    await writeFile(info.outputPath('producer-parity.json'),JSON.stringify({producer,width,results,adaptations:['Producer site scripts/includes disabled and same locked CSS injected.','Caller heading levels adapt to gallery hierarchy; chart footer omitted intentionally.','Caller-provided icon contents can differ; shared decorative wrapper retained.']},null,2));
    await context.close();
  });
}
