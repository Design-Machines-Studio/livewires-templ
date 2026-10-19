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
  if (name === 'people-picker') {
    for (const form of document.querySelectorAll('form[data-demo-picker]')) {
      const rows = [...form.querySelectorAll('tr[data-demo-person]')];
      const groups = [...form.querySelectorAll('[data-demo-group]')];
      const direct = new Set(rows.filter(row => row.dataset.demoDirect === 'true').map(row => row.dataset.demoPerson));
      const excluded = new Set(rows.filter(row => row.dataset.demoExcluded === 'true').map(row => row.dataset.demoPerson));
      const locked = new Set(rows.filter(row => row.dataset.demoLocked === 'true').map(row => row.dataset.demoPerson));
      const selectedGroups = new Set(groups.filter(button => button.getAttribute('aria-pressed') === 'true').map(button => button.dataset.demoGroup));
      const organizerState = new Map([...form.querySelectorAll('[data-demo-organizer]')].map(input => [input.dataset.demoOrganizer, input.checked]));
      const status = form.querySelector('[data-demo-status]');
      const circleLabels = new Map(groups.map(button => [button.dataset.demoGroup, button.textContent.trim()]));
      // Clear first so repeating an identical message is announced again.
      const announce = text => {
        if (!status) return;
        status.textContent = '';
        later(() => { status.textContent = text; }, 60);
      };
      const selectedCount = () => rows.filter(row => row.querySelector('input[type="checkbox"][id]').checked).length;
      const countText = () => { const n = selectedCount(); return `${n} ${n === 1 ? 'person' : 'people'} selected.`; };

      const render = () => {
        for (const row of rows) {
          const id = row.dataset.demoPerson;
          const input = row.querySelector(`input[type="checkbox"][id]`);
          const memberships = (row.dataset.demoCircles || '').split(',').filter(Boolean);
          const fromGroups = memberships.some(group => selectedGroups.has(group)) && !excluded.has(id);
          const selected = locked.has(id) || direct.has(id) || fromGroups;
          input.checked = selected;
          for (const organizer of row.querySelectorAll('[data-demo-organizer]')) {
            const isLocked = locked.has(id);
            organizer.disabled = isLocked || !selected;
            organizer.checked = isLocked ? false : selected && Boolean(organizerState.get(id));
          }
          // Keep the detail line true to the simulated state. Locked rows keep their fixed reason.
          const detail = row.querySelector('.people-picker-detail');
          if (detail && !locked.has(id)) {
            const circles = fromGroups ? memberships.filter(group => selectedGroups.has(group)).map(group => circleLabels.get(group)).filter(Boolean) : [];
            const parts = [];
            if (selected && organizerState.get(id)) parts.push('Organizer');
            parts.push(...circles);
            if (selected && circles.length === 0 && direct.has(id)) parts.push('Added directly');
            const text = excluded.has(id) ? 'Removed' : parts.length ? parts.join(', ') : 'Not selected';
            detail.textContent = text;
            if (detail.hasAttribute('title')) detail.title = text;
          }
        }
        for (const button of groups) button.setAttribute('aria-pressed', String(selectedGroups.has(button.dataset.demoGroup)));
      };

      form.addEventListener('change', event => {
        const input = event.target;
        if (!(input instanceof HTMLInputElement)) return;
        if (input.matches('[data-demo-organizer]')) {
          organizerState.set(input.dataset.demoOrganizer, input.checked);
          render();
          return;
        }
        if (!input.matches('input[type="checkbox"][id]')) return;
        const row = input.closest('tr[data-demo-person]');
        if (!row || locked.has(row.dataset.demoPerson)) return;
        const id = row.dataset.demoPerson;
        if (input.checked) { direct.add(id); excluded.delete(id); }
        else { direct.delete(id); excluded.add(id); }
        render();
      }, { signal: lifecycle.signal });

      form.addEventListener('submit', event => {
        const button = event.submitter;
        if (!(button instanceof HTMLButtonElement)) return;
        if (button.hasAttribute('data-demo-group')) {
          event.preventDefault();
          const group = button.dataset.demoGroup;
          if (selectedGroups.has(group)) selectedGroups.delete(group); else selectedGroups.add(group);
          render();
          announce(`${button.textContent.trim()} circle ${selectedGroups.has(group) ? 'selected' : 'removed'}. ${countText()}`);
        } else if (button.hasAttribute('data-demo-clear')) {
          event.preventDefault();
          selectedGroups.clear();
          direct.clear();
          excluded.clear();
          render();
          announce(`Selections cleared. Locked people remain selected. ${countText()}`);
        }
      }, { signal: lifecycle.signal });

      render();
      for (const control of form.querySelectorAll('[data-demo-clear], [data-demo-group]')) {
        if (!control.hasAttribute('data-demo-unavailable')) control.disabled = false;
      }
    }
  }
  if (name === 'loaders') {
    const overlay = document.querySelector('[data-gallery-overlay]');
    const status = overlay.parentElement;
    // Keep the status wrapper outside the inert content, as in the producer reference.
    document.body.append(status);
    const changed = new Set();
    let previous;
    let timer;
    const close = () => {
      clearTimeout(timer);
      overlay.hidden = true;
      for (const node of changed) node.inert = false;
      changed.clear();
      if (previous?.isConnected) previous.focus();
      previous = null;
    };
    document.querySelector('[data-preview-overlay]').addEventListener('click', () => {
      if (!overlay.hidden) return;
      previous = document.activeElement;
      overlay.hidden = false;
      overlay.querySelector('.visually-hidden').focus();
      for (const node of document.body.children) {
        if (node !== status && !node.inert) { node.inert = true; changed.add(node); }
      }
      timer = setTimeout(close, 2000);
    }, { signal: lifecycle.signal });
    document.addEventListener('keydown', event => {
      if (event.key === 'Escape' && !overlay.hidden) { event.preventDefault(); close(); }
    }, { signal: lifecycle.signal });
    controllers.push({ dispose: close });
  }
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
  if (name === 'save-status') {
    const schemeRoot = document.querySelector('[data-save-scheme]');
    for (const controls of schemeRoot.querySelectorAll('[data-save-controls]')) controls.hidden = false;
    for (const button of document.querySelectorAll('[data-save-scheme-class]')) {
      button.addEventListener('click', () => {
        schemeRoot.className = `stack ${button.dataset.saveSchemeClass}`;
        for (const other of document.querySelectorAll('[data-save-scheme-class]')) {
          other.setAttribute('aria-pressed', String(other === button));
          other.classList.toggle('underline', other === button);
        }
      }, { signal: lifecycle.signal });
    }
    for (const button of document.querySelectorAll('[data-save-demo]')) {
      button.addEventListener('click', () => {
        const editor = document.querySelector(`[data-save-editor="${CSS.escape(button.dataset.saveDemo)}"]`);
        const line = editor?.querySelector('[data-save-status]');
        const support = line?.nextElementSibling;
        if (!line || !support) return;
        line.dataset.saveState = button.dataset.state;
        line.textContent = button.dataset.label || 'Save not confirmed';
        support.textContent = button.dataset.support || '';
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
