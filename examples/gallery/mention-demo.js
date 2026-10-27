/*
 * Mention suggestions: gallery demonstration only.
 * Adapted from Live Wires v2.27.0, 1c52666402c0ec2de4ed64736ea679a8a8573838.
 * MIT. State/row fragments here come from the public Templ components.
 * No part of this adapter is imported by the component library.
 *
 * This page script shows one way a consumer can honor the contract in
 * docs/superpowers/specs/mention-suggestions.md. It is not shipped in /dist and
 * is not a supported controller. Baseplate owns the real trigger, insertion and
 * serialization. The caller supplies results, stable IDs and selection handling.
 */

const GAP = 4;

function scrollParent(node) {
  for (let el = node.parentElement; el && el !== document.body; el = el.parentElement) {
    const { overflowX, overflowY } = getComputedStyle(el);
    if (/(auto|scroll|hidden|clip)/.test(overflowX + overflowY)) return el;
  }
  return null;
}

function caretRect(range, editor) {
  const rects = range.getClientRects();
  if (rects.length) return rects[rects.length - 1];
  const rect = range.getBoundingClientRect();
  if (rect.width || rect.height) return rect;
  // An empty line has no glyph rect. Fall back to the editor's start edge.
  const box = editor.getBoundingClientRect();
  return { top: box.top, bottom: box.top + parseFloat(getComputedStyle(editor).lineHeight || 20), left: box.left, right: box.left };
}

/** Find "@query" immediately before a collapsed caret in one text node. */
function findToken(editor, trigger) {
  const selection = editor.ownerDocument.getSelection();
  if (!selection.rangeCount || !selection.isCollapsed) return null;
  const range = selection.getRangeAt(0);
  const node = range.startContainer;
  if (node.nodeType !== Node.TEXT_NODE || !editor.contains(node)) return null;
  const before = node.data.slice(0, range.startOffset);
  const t = trigger.replace(/[.*+?^${}()|[\]\\-]/g, '\\$&');
  const match = before.match(new RegExp(`(^|\\s)${t}([^\\s${t}]{0,40})$`));
  if (!match) return null;
  return { node, start: range.startOffset - match[2].length - trigger.length, end: range.startOffset, query: match[2] };
}

export function attachMentionDemo({ editor, root, search, onSelect = () => {}, trigger = '@', loadingDelay = 400 }) {
  const doc = editor.ownerDocument;
  const list = root.querySelector('[role="listbox"]');
  const status = root.querySelector('[role="status"]');
  if (!list?.id || !status) throw new TypeError('Mount a root with an identified listbox and a status region.');

  // Editor contract: textbox semantics stay; these attributes are added while attached.
  editor.setAttribute('aria-autocomplete', 'list');
  editor.setAttribute('aria-haspopup', 'listbox');
  editor.setAttribute('aria-controls', list.id);

  let items = [];
  let active = -1;
  let token = null;
  let suppressed = null; // token start dismissed with Escape
  let request = 0;
  let open = false;
  let loadingTimer;
  let renderedQuery = '';
  const closedFragment = root.cloneNode(true);

  function say(message) {
    // Write only on change so typing does not repeat announcements.
    if (status.textContent !== message) status.textContent = message;
  }
  function setState(state) { root.dataset.state = state; }
  function options() { return [...list.querySelectorAll('[role="option"]')]; }

  function setActive(index) {
    const all = options();
    all.forEach((option, i) => {
      const on = i === index;
      option.toggleAttribute('data-active', on);
      option.classList.toggle('rounded-sm', on);
      option.classList.toggle('scheme-subtle', on);
      option.setAttribute('aria-selected', on ? 'true' : 'false');
    });
    active = index;
    if (index >= 0 && all[index]) {
      editor.setAttribute('aria-activedescendant', all[index].id);
      keepVisible(all[index]);
    } else {
      editor.removeAttribute('aria-activedescendant');
    }
  }

  // Scroll only the listbox. scrollIntoView could also scroll the page.
  function keepVisible(option) {
    const box = list.getBoundingClientRect();
    const row = option.getBoundingClientRect();
    if (row.top < box.top) list.scrollTop -= box.top - row.top;
    else if (row.bottom > box.bottom) list.scrollTop += row.bottom - box.bottom;
  }

  function move(step) {
    if (!items.length || list.getAttribute('aria-busy') === 'true') return;
    let next = active;
    for (let i = 0; i < items.length; i++) {
      next = (next + step + items.length) % items.length;
      if (!items[next]?.disabled) return setActive(next);
    }
  }

  // Preserve the mounted list/status nodes. Only rows come from the server.
  // A new status is written only when its message changes.
  function render(fragment) {
    const nextList = fragment.querySelector('[role="listbox"]');
    list.replaceChildren(...nextList.childNodes);
    list.hidden = nextList.hidden;
    if (nextList.hasAttribute('aria-busy')) list.setAttribute('aria-busy', 'true');
    else list.removeAttribute('aria-busy');
    setState(fragment.dataset.state);
    say(fragment.querySelector('[role="status"]').textContent);
    items = options().map(row => ({
      id: row.dataset.personId,
      name: row.querySelector('[data-mention-name]').textContent,
      detail: row.querySelector('[data-mention-detail]')?.textContent || '',
      disabled: row.getAttribute('aria-disabled') === 'true'
    }));
    list.scrollTop = 0;
  }

  function position() {
    if (!token) return;
    const range = doc.createRange();
    range.setStart(token.node, token.start);
    range.setEnd(token.node, token.start);
    const caret = caretRect(range, editor);
    const viewport = { top: 0, left: 0, right: doc.documentElement.clientWidth, bottom: innerHeight };
    const parent = scrollParent(editor);
    const clip = parent ? parent.getBoundingClientRect() : viewport;
    const bounds = {
      top: Math.max(viewport.top, clip.top), bottom: Math.min(viewport.bottom, clip.bottom),
      left: Math.max(viewport.left, clip.left), right: Math.min(viewport.right, clip.right)
    };
    root.style.removeProperty('--mention-suggestions-max-block');
    const preferred = parseFloat(getComputedStyle(root).maxBlockSize) || 240;
    const width = root.offsetWidth;
    const below = bounds.bottom - caret.bottom - GAP * 2;
    const above = caret.top - bounds.top - GAP * 2;
    const needed = root.offsetHeight; // Natural height, capped by the default limit.
    if (below >= needed || below >= above) {
      // Clamp so an off-screen caret never pushes the list out of the viewport.
      root.style.setProperty('--mention-suggestions-top', `${Math.max(GAP, Math.min(caret.bottom + GAP, viewport.bottom - GAP - needed))}px`);
      root.style.removeProperty('--mention-suggestions-bottom');
      root.style.setProperty('--mention-suggestions-max-block', `${Math.max(0, Math.min(preferred, below))}px`);
      root.dataset.placement = 'below';
    } else {
      root.style.setProperty('--mention-suggestions-bottom', `${Math.max(GAP, Math.min(viewport.bottom - caret.top + GAP, viewport.bottom - GAP - needed))}px`);
      root.style.removeProperty('--mention-suggestions-top');
      root.style.setProperty('--mention-suggestions-max-block', `${Math.max(0, Math.min(preferred, above))}px`);
      root.dataset.placement = 'above';
    }
    const rtl = getComputedStyle(editor).direction === 'rtl';
    root.dir = rtl ? 'rtl' : 'ltr';
    const ideal = rtl ? caret.right - width : caret.left;
    const min = Math.max(bounds.left, viewport.left + GAP);
    const max = Math.min(bounds.right, viewport.right - GAP) - width;
    root.style.setProperty('--mention-suggestions-left', `${Math.max(min, Math.min(ideal, max))}px`);
  }

  async function query(text) {
    const id = ++request;
    const stale = items.length > 0;
    if (stale) { list.setAttribute('aria-busy', 'true'); list.hidden = false; }
    else list.hidden = true;
    setActive(-1);
    clearTimeout(loadingTimer);
    let pending = true;
    const showLoading = async () => {
      try {
        const fragment = await search(renderedQuery, 'loading', stale);
        if (id !== request || !open || !pending) return;
        render(fragment); setActive(-1); position();
      } catch { /* Final response below owns the error state. */ }
    };
    loadingTimer = setTimeout(showLoading, loadingDelay);
    try {
      const fragment = await search(text);
      if (id !== request || !open) return; // Drop an older query's reply.
      pending = false;
      clearTimeout(loadingTimer);
      render(fragment);
      renderedQuery = text;
      setActive(root.dataset.state === 'ready' ? items.findIndex(p => !p.disabled) : -1);
    } catch {
      if (id !== request || !open) return;
      pending = false;
      clearTimeout(loadingTimer);
      // Network failure fallback: reuse server-rendered error anatomy supplied
      // in the page, without constructing option markup in JavaScript.
      render(document.getElementById('static-error-root').cloneNode(true));
      setActive(-1);
    }
    position();
  }

  function close() {
    if (!open) return;
    open = false;
    request++;
    clearTimeout(loadingTimer);
    token = null;
    setActive(-1);
    render(closedFragment.cloneNode(true));
    list.hidden = true;
    list.removeAttribute('aria-busy');
    setState('closed');
    say('');
  }

  function select(index) {
    const person = items[index];
    if (!person || person.disabled || !token) return;
    // Demo insertion: plain text. Baseplate owns real mention nodes and serialization.
    const range = doc.createRange();
    range.setStart(token.node, token.start);
    range.setEnd(token.node, token.end);
    range.deleteContents();
    const inserted = doc.createTextNode(`${trigger}${person.name} `);
    range.insertNode(inserted);
    const caret = doc.createRange();
    caret.setStartAfter(inserted);
    caret.collapse(true);
    const selection = doc.getSelection();
    selection.removeAllRanges();
    selection.addRange(caret);
    close();
    editor.normalize();
    onSelect(person);
  }

  function update() {
    const next = findToken(editor, trigger);
    if (!next) { suppressed = null; return close(); }
    if (suppressed && suppressed.node === next.node && suppressed.start === next.start) return;
    suppressed = null;
    const changed = !token || token.query !== next.query || !open;
    token = next;
    open = true;
    if (changed) query(next.query); else position();
  }

  function keydown(event) {
    // Handle keys only while the list is on screen; during the loading delay it is not.
    if (!open || root.dataset.state === 'closed' || event.isComposing || event.keyCode === 229) return;
    const modified = event.shiftKey || event.ctrlKey || event.metaKey || event.altKey;
    switch (event.key) {
      case 'ArrowDown':
      case 'ArrowUp':
        // With no rows on screen, arrows move the caret; leaving the token closes the list.
        if (list.hidden || modified) return;
        event.preventDefault();
        move(event.key === 'ArrowDown' ? 1 : -1);
        break;
      case 'Enter':
        event.preventDefault(); // Never send or submit while open, even with modifiers.
        if (active >= 0 && !modified) select(active); else close();
        break;
      case 'Tab':
        if (active >= 0 && !modified) { event.preventDefault(); select(active); } else close();
        break;
      case 'Escape':
        event.preventDefault();
        suppressed = token && { node: token.node, start: token.start };
        close();
        break;
      default: break;
    }
  }

  const onMouseDown = (event) => event.preventDefault(); // Keep editor focus and selection.
  function onClick(event) {
    const option = event.target.closest?.('[role="option"]');
    if (option) select(options().indexOf(option));
  }
  const onSelectionChange = () => { if (doc.activeElement === editor) update(); };
  const onFocusOut = () => close();
  const onCompositionEnd = () => update();
  const onInput = (event) => { if (!event.isComposing) update(); };
  const onViewport = () => { if (open) position(); };

  editor.addEventListener('keydown', keydown);
  editor.addEventListener('input', onInput);
  editor.addEventListener('compositionend', onCompositionEnd);
  editor.addEventListener('focusout', onFocusOut);
  doc.addEventListener('selectionchange', onSelectionChange);
  root.addEventListener('mousedown', onMouseDown);
  root.addEventListener('click', onClick);
  addEventListener('resize', onViewport);
  addEventListener('scroll', onViewport, true);
  setState('closed');
  list.hidden = true;

  return {
    get open() { return open; },
    get active() { return active; },
    close,
    dispose() {
      close();
      editor.removeEventListener('keydown', keydown);
      editor.removeEventListener('input', onInput);
      editor.removeEventListener('compositionend', onCompositionEnd);
      editor.removeEventListener('focusout', onFocusOut);
      doc.removeEventListener('selectionchange', onSelectionChange);
      root.removeEventListener('mousedown', onMouseDown);
      root.removeEventListener('click', onClick);
      removeEventListener('resize', onViewport);
      removeEventListener('scroll', onViewport, true);
      for (const name of ['aria-autocomplete', 'aria-haspopup', 'aria-controls', 'aria-activedescendant']) editor.removeAttribute(name);
    }
  };
}


// Page-only attachment. Search is a fixture endpoint, never a real directory.
const controllers = [];
for (const editor of document.querySelectorAll('[data-mention-demo]')) {
  const root = document.getElementById(editor.dataset.mentionDemo);
  const form = editor.closest('form');
  const log = form.querySelector('[data-mention-log]');
  // Keep fixed geometry outside the editor and layout/paint containment.
  document.body.append(root);
  const listID = root.querySelector('[role="listbox"]').id;
  controllers.push(attachMentionDemo({
    editor, root,
    async search(q, state = 'ready', stale = false) {
      const mode = form.querySelector('[name="mention-scenario"]').value;
      const params = new URLSearchParams({ list: listID, q, state, stale: String(stale), mode });
      const response = await fetch(`/demo/mentions?${params}`);
      if (!response.ok) throw new Error('Demonstration response unavailable');
      const doc = new DOMParser().parseFromString(await response.text(), 'text/html');
      const fragment = doc.querySelector('[data-mention-suggestions]');
      if (!fragment) throw new Error('Demonstration fragment missing');
      return fragment;
    },
    onSelect(person) {
      log.textContent = `Selected ${person.name}${person.detail ? ` (${person.detail})` : ''}, ID ${person.id}. Nothing saved.`;
    }
  }));
  form.addEventListener('keydown', event => {
    if (event.key === 'Enter' && !event.shiftKey && !event.ctrlKey && !event.metaKey && !event.altKey &&
        !event.isComposing && event.keyCode !== 229 && !event.defaultPrevented && event.target === editor) {
      event.preventDefault(); form.requestSubmit();
    }
  });
  form.addEventListener('submit', event => {
    event.preventDefault();
    log.textContent = `Simulated send: ${editor.textContent.trim() || '(empty)'}. Nothing saved.`;
  });
}
addEventListener('pagehide', () => controllers.forEach(controller => controller.dispose()), { once: true });
