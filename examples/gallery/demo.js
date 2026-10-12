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
  const form = document.querySelector('#scenario-form');
  let occurrence = 0;
  form?.addEventListener('submit', async (event) => {
    if (name === 'sortable') return;
    event.preventDefault();
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
