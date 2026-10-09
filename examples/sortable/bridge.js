// Datastar owns HTTP and SSE. The producer owns all sorting and focus behavior.
// This adapter provides only request/connection correlation and confirmation.
const current = new WeakMap();

function removed(records) {
  for (const record of records) {
    for (const node of record.removedNodes) {
      if (!(node instanceof Element)) continue;
      const hosts = node.matches('lw-sortable-list') ? [node] : [];
      hosts.push(...node.querySelectorAll('lw-sortable-list'));
      for (const host of hosts) {
        current.delete(host);
        host.removeAttribute('data-active-request');
      }
    }
  }
}

// Observe disconnection only, never child changes as evidence of acceptance.
const connections = new MutationObserver(removed);
connections.observe(document.documentElement, { childList: true, subtree: true });

function matches(host, payload) {
  removed(connections.takeRecords());
  return host.isConnected && current.get(host) === payload
    && host.dataset.activeRequest === payload.marker;
}

window.SortableBridge = Object.freeze({
  payload(host, detail) {
    removed(connections.takeRecords());
    if (!host.isConnected || !detail || typeof detail.requestId !== 'string'
      || typeof detail.itemId !== 'string' || typeof detail.before !== 'string') return null;
    const order = host.querySelector('[data-sortable-list]');
    if (!order) return null;
    const payload = {
      itemId: detail.itemId, before: detail.before, requestId: detail.requestId,
      marker: crypto.randomUUID(), revision: order.dataset.revision,
      csrf: host.dataset.csrf, variant: host.dataset.variant,
      density: host.dataset.density, controls: host.dataset.controls,
    };
    current.set(host, payload);
    host.dataset.activeRequest = payload.marker;
    return payload;
  },

  result(host, result) {
  // A scalar signal replaces the complete envelope. Datastar's recursive
  // object merge must never fill a missing status from an earlier response.
    if (typeof result !== 'string') return;
    try { result = JSON.parse(result); } catch { return; }
    const payload = current.get(host);
    if (!payload || !matches(host, payload) || !result
      || result.list !== host.dataset.list || result.requestId !== payload.requestId
      || result.marker !== payload.marker
      || !['accepted', 'rejected'].includes(result.status)
      || typeof result.revision !== 'string' || !/^[1-9][0-9]*$/.test(result.revision)) return;
    const order = host.querySelector('[data-sortable-list]');
    if (!order || order.dataset.delivery !== payload.requestId
      || order.dataset.revision !== result.revision) return;
    if ([...order.querySelectorAll('[name="revision"]')]
      .some(field => field.value !== result.revision)) return;
    // A matching explicit result plus the expected delivered child revision is
    // required. Invalid status, unrelated patches and old replies leave pending.
    if (!host.resolveMove({ requestId: payload.requestId, status: result.status, message: result.message })) return;
    current.delete(host);
    host.removeAttribute('data-active-request');
    host.dataset.settledRevision = result.revision;
  },

  failure(host, payload, error) {
    // The pinned action rejects FetchFailed after its zero retry budget on a
    // failed request transport. It wraps the exhausted-budget string without
    // retaining that string in Error metadata, so do not test for that text.
    // Successful streams with invalid/missing status do not throw and remain
    // pending. Other action error types do not enter this failure path.
    if (!matches(host, payload) || !(error instanceof Error)
      || !error.message.startsWith('FetchFailed\n')) return;
    host.resolveMove({ requestId: payload.requestId, status: 'rejected',
      message: 'Saving could not be confirmed. Reload to see the authoritative order.' });
    current.delete(host);
    host.removeAttribute('data-active-request');
  },
});
