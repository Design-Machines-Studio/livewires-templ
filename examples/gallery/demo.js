// Application-owned simulation only. Producer controllers own component state.
const name = document.body.dataset.example;
const lifecycle = new AbortController();
const controllers = [];
const timers = new Set();
const later = (fn, ms) => { const id = setTimeout(() => { timers.delete(id); fn(); }, ms); timers.add(id); };
const outcome = document.querySelector('#demo-outcome');
const unrelated = document.querySelector('#unrelated');
if (unrelated) {
  unrelated.hidden = false;
  unrelated.addEventListener('click', () => { document.querySelector('#unrelated-note').textContent = 'Unrelated note updated.'; }, { signal: lifecycle.signal });
}

async function setup() {
  let loading;
  let toasts;
  if (name === 'relative-time') await import('/assets/relative-time.js');
  if (name === 'copy') await import('/assets/copy-button.js');
  if (name === 'loading') {
    const { createLoadingOperation } = await import('/assets/loading.js');
    loading = createLoadingOperation({ root: document.querySelector('[data-loading-region]'), content: '[data-loading-content]', indicator: '[data-loading-indicator]', status: '[data-loading-status]', delay: 100, minDuration: 150 });
    controllers.push(loading);
  }
  if (name === 'toast') {
    const { createToastRegion } = await import('/assets/toast-region.js');
    toasts = createToastRegion({ root: document.querySelector('[data-toast-region]') });
    controllers.push(toasts);
  }
  if (name === 'inline-edit') {
    unrelated.textContent = 'Simulate external update';
    const { createInlineEditor } = await import('/assets/inline-edit.js');
    for (const root of document.querySelectorAll('[data-inline-edit]')) {
      const editor = createInlineEditor({ root });
      controllers.push(editor);
      root.addEventListener('lw-inline-edit-commit', (event) => {
        const detail = event.detail;
        const scenario = document.querySelector('#scenario').value;
        const settle = () => {
          const result = { ...detail, outcome: scenario === 'failure' ? 'rejected' : scenario === 'unknown' ? 'unknown' : 'confirmed', value: detail.value, message: scenario === 'failure' ? 'Simulated rejection. Your draft remains.' : scenario === 'unknown' ? 'Simulated unknown result. Cancel or reload to recover.' : 'Simulated confirmation. Nothing was saved.' };
          if (editor.applyResult(result)) outcome.textContent = result.message;
        };
        if (scenario === 'delayed') later(settle, 800); else settle();
      }, { signal: lifecycle.signal });
    }
    unrelated.addEventListener('click', () => {
      for (const editor of controllers) editor.updateConfirmed?.('Simulated external update');
      outcome.textContent = 'Simulated external update applied. Drafts are kept.';
    }, { signal: lifecycle.signal });
  }
  if (name === 'combobox') {
    const { createCombobox } = await import('/assets/combobox.js');
    const responseMode = document.querySelector('#remote-scenario');
    document.querySelector('#remote-scenario-control').hidden = false;
    for (const root of document.querySelectorAll('[data-combobox]')) {
      const remote = root.hasAttribute('data-search');
      const control = createCombobox({ root, mode: remote ? 'remote' : 'local' });
      controllers.push(control);
      if (!remote) continue;

      const choices = [...root.querySelectorAll('select option')].map(option => ({
        value: option.value,
        label: option.label,
        disabled: option.disabled,
      })).filter(option => option.value !== '');
      root.addEventListener('combobox-query', (event) => {
        const request = event.detail;
        const scenario = responseMode.value;
        const query = request.query.toLocaleLowerCase();
        const delay = scenario === 'delayed' ? 650
          : scenario === 'out-of-order' ? query === 'pa' ? 900 : query === 'par' ? 100 : 250
            : 80;
        later(() => {
          if (lifecycle.signal.aborted) return;
          if (scenario === 'error') {
            control.applyResults({ ...request, items: [], error: 'Try again.' });
            return;
          }
          if (scenario === 'empty') {
            control.applyResults({ ...request, items: [] });
            return;
          }
          const items = choices.filter(item => item.label.toLocaleLowerCase().includes(query) || item.value.toLocaleLowerCase().includes(query));
          control.applyResults({ ...request, items });
        }, delay);
      }, { signal: lifecycle.signal });
    }
  }
  const form = document.querySelector('#scenario-form');
  if (form) form.hidden = false;
  let occurrence = 0;
  let scenarioPending = false;
  form?.addEventListener('submit', async (event) => {
    event.preventDefault();
    if (name === 'sortable' || name === 'inline-edit') {
      return;
    }
    if (scenarioPending) return;
    scenarioPending = true;
    const runButton = form.querySelector('[name="run"]');
    runButton.disabled = true;
    const token = loading?.begin({ signal: lifecycle.signal });
    outcome.textContent = 'Simulated request pending.';
    try {
      const response = await fetch('/demo/result?scenario=' + encodeURIComponent(form.elements.scenario.value), { signal: lifecycle.signal });
      if (!response.ok) throw new Error('Demo transport failed');
      const result = await response.json();
      token?.finish(result.scenario === 'failure' ? 'failed' : 'unknown');
      outcome.textContent = result.message;
      if (toasts) {
        // The server renders each occurrence through the real Templ Toast API.
        const fragment = await fetch('/examples/toast?run=1&scenario=' + encodeURIComponent(result.scenario), { signal: lifecycle.signal });
        if (!fragment.ok) throw new Error('Demo fragment failed');
        const doc = new DOMParser().parseFromString(await fragment.text(), 'text/html');
        const node = doc.querySelector('[data-toast-id]');
        node.dataset.toastId = 'demo:outcome:' + (++occurrence);
        document.querySelector('[data-toast-list]').append(node);
        toasts.refresh();
      }
    } catch (error) {
      token?.finish(lifecycle.signal.aborted ? 'cancelled' : 'failed');
      if (!lifecycle.signal.aborted) outcome.textContent = 'Demo transport failed. Nothing was saved.';
    } finally {
      scenarioPending = false;
      runButton.disabled = false;
    }
  }, { signal: lifecycle.signal });
  if (name === 'sortable') {
    await import('/assets/sortable-list.js');
    const host = document.querySelector('lw-sortable-list');
    host.addEventListener('sortable-move-request', (event) => {
      const { requestId } = event.detail;
      const scenario = document.querySelector('#scenario').value;
      const settle = () => {
        if (scenario === 'unknown') { outcome.textContent = 'Simulated unknown result. Reload to recover.'; return; }
        host.resolveMove({ requestId, status: scenario === 'failure' ? 'rejected' : 'accepted', message: 'Simulated failure. Move was not applied.' });
        outcome.textContent = scenario === 'failure' ? 'Simulated failure. Nothing was saved.' : 'Simulated move confirmed. Reload resets this page.';
        // Refresh application-owned native form values after the move/rollback.
        const order = [...host.querySelectorAll('[data-sortable-item]')].map(node => node.dataset.itemId);
        for (const form of host.querySelectorAll('form')) {
          const i = order.indexOf(form.elements.itemId.value);
          const direction = form.querySelector('[data-sortable-move]').dataset.sortableMove;
          form.elements.order.value = order.join(',');
          form.elements.before.value = direction === 'up' ? order[Math.max(0, i - 1)] : (order[i + 2] || '');
          form.querySelector('button').disabled = direction === 'up' ? i === 0 : i === order.length - 1;
        }
      };
      if (scenario === 'delayed') later(settle, 800); else settle();
    }, { signal: lifecycle.signal });
  }
  document.body.dataset.demoReady = 'true';
}

function teardown() {
  lifecycle.abort();
  for (const timer of timers) clearTimeout(timer);
  timers.clear();
  for (const controller of controllers) controller.dispose();
}
window.addEventListener('pagehide', teardown, { once: true });
// Returning from the back/forward cache gets a fresh controller lifetime.
window.addEventListener('pageshow', (event) => { if (event.persisted) location.reload(); });
setup().catch((error) => { teardown(); throw error; });
