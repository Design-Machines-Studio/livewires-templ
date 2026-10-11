# Optional interaction utility markup

These components render the semantic HTML hooks used by the Live Wires
relative-time, copy, loading and toast modules. This Go package remains
CSS-free. It does not bundle browser code, import Datastar or manage application
requests, persistence, authorization, retries or success confirmation.

All `Attrs` values are developer-controlled. Templ escapes values, but does not
restrict attribute names. Do not populate `Attrs` from request or member input.
Component-owned attributes come first; avoid repeating them in `Attrs` because
HTML keeps the first duplicate.

## Relative time

```templ
@component.RelativeTimeComponent(component.RelativeTimeProps{
    DateTime: "2026-10-10T15:30:00-04:00",
    Fallback: "October 10, 2026 at 3:30 PM EDT",
    Lang: "en", TimeZone: "America/New_York",
})
```

The wrapper emits `<time datetime>` with visible absolute fallback text and a
hidden relative output. Without the optional module, or when a value, locale or
time zone cannot be enhanced safely, the absolute date remains visible.
Date-only values keep their calendar meaning; timestamp values should retain
their original UTC offset. The optional controller uses one shared timer and
releases it when the last instance is removed. It does not move focus.

## Copy

```templ
@component.CopyComponent(component.CopyProps{
    ID: "share-address", Label: "Share address",
    Value: "https://example.test/meeting", ButtonLabel: "Copy address",
})
```

The source remains read-only, visible and selectable without JavaScript. The
optional copy controller reports success only after the browser resolves the
clipboard write. Denial, missing clipboard support or a rejected write leaves
the value available and announces manual selection instructions. It does not
claim success early, provide a hidden clipboard fallback, or change an
application's approved clipboard or Datastar policy.

## Loading operations

`LoadingRegion` provides an operation-local root, a busy content boundary, a
decorative hidden skeleton and a polite status outside the busy content. Keep
the status outside the content and keep unrelated fields/controls usable.
Load the optional Live Wires loading module only in the consuming application.
The existing `Loader` component and its calls keep their current behavior.

```templ
@component.LoadingRegionComponent(component.LoadingRegionProps{ID: "editor-save", Label: "Editor save"}) {
    <p>Current editor content stays available.</p>
}
```

Create one controller for one finite operation region. Each overlapping
request gets its own token; finish only that token. Delay and minimum display
time affect presentation, not busy state or save results. A completed request
is `unknown` until the application confirms its matching command/version.
Report pending, failed, cancelled and unknown distinctly. Never treat a fetch
completion, unrelated patch, SSE connection or timeout as saved. The module
does not disable controls, move focus, abort application work or retry requests.

The controller owns the operation content's `aria-busy` and
`data-loading-state`, the indicator's `hidden` state, and status text. Do not
attach two controllers to the same targets. Keep a stable root; replacement
content, indicator and status targets are rebound without clearing active
tokens. Removing an initiating button does not finish its request. Removing a
required target or root disposes the operation. A reinserted root needs fresh
setup. On unmount, the application removes its listeners and aborts its own
requests as appropriate, disconnects any adapter, then calls `dispose()`.
Disposal clears controller timers/listeners and restores owned attributes, but
does not cancel or roll back application work.

## Toast regions

Standalone `Toast` keeps its polite/assertive semantics. For managed behavior,
use one stable `ToastRegion` and managed `ToastComponent` occurrences:

```templ
@component.ToastRegionComponent(component.ToastRegionProps{Label: "Editor feedback"}) {
    @component.ToastComponent(component.ToastProps{
        ID: "save-42:confirmed", Message: "Your changes were saved.",
        Variant: "success", Duration: &sixSeconds, Dismissible: true,
        DismissLabel: "Dismiss saved message", Managed: true,
    })
}
```

Here `sixSeconds` is an application-owned `int` value set to `6000` before the
template is rendered.

The region owns the stable `data-toast-region` boundary and the patchable
`data-toast-list`. Each managed toast has a stable occurrence ID, kind,
optional duration, message hook and a child-content slot for native actions.
Region state is local to one controller;
two independent regions may use the same occurrence ID. A managed toast omits
its standalone live role so the controller can announce through its own stable
polite/assertive nodes outside the list. Dismiss controls are native
`type="button"` buttons with an accessible label.

An ID identifies one immutable feedback occurrence during a controller
lifetime. Replacing a node or updating its content/configuration with the same
ID preserves dismissal, announcement history and remaining duration. A new
meaningful outcome needs a new ID. Removal pauses and suspends its timer;
reinsertion resumes remaining time without announcing again. Tombstones remain
until explicitly forgotten or the region is disposed. To intentionally reuse
an ID, remove every node with it, let the controller reconcile (or call
`refresh()`), then call `forget(id)` while no canonical node remains.

No duration means sticky. Pending and error messages, plus messages with active
native actions, remain until explicit dismissal/application removal even if a
duration is supplied. Hover and focus pause time independently. Escape dismisses
only the focused toast; Enter and Space activate its native button. Arrival and
expiry do not move focus. Explicit dismissal returns focus to another dismiss
button, an application-provided `[data-toast-focus-return]` target or the
region. Preserve focused action subtrees in application patches; the controller
does not reparent or restore a subtree the application removed.

The controller owns each toast's `hidden`, `data-toast-state`, `role` and
`aria-live`, the list's review state, and generated controls/announcement nodes.
Keep IDs, kind, duration, message content and action listeners application
owned. Do not set live roles on managed messages or nest live regions inside
them. On unmount call `dispose()`: it releases timers, observers and delegated
listeners, removes generated nodes and restores controller-owned attributes.
Removing the root or list disposes it; reinsertion needs a new controller and
starts a new identity scope. Dismissal is not durable persistence or a server
acknowledgement.

## Pinned producer assets and proof level

The optional modules and references below are pinned to the exact merged
Live Wires commits. The source repository tracks these built distribution
files under `public/dist`; these pins establish file provenance only. They do
not prove a package release, CDN publication or a running consumer has loaded
them. Consumers must verify their own served bytes.

| Module | Source commit | Pinned file |
| --- | --- | --- |
| Relative time | `f73b8a26d9c449ede4069edfe13eb50c2876a732` | `relative-time.js`, SHA-256 `3668407d4eb37e8eac33d51ec04ef0e6d1d085a7bac27d983700b3a38ff514af` |
| Copy | `f73b8a26d9c449ede4069edfe13eb50c2876a732` | `copy-button.js`, SHA-256 `06fac73364cd2fb1721fbd6ccfbdb204f1208e184dd8140b34badb5f7815b903` |
| Loading | `98b890dc804813b0b36b9a484bf49e3bdc7f196b` | `loading.js`, SHA-256 `5de053802835444b6a782d8a5a58e0c33a9aa9fa1490c5dad819340a23a253f6` |
| Toast region | `6e0ff82c222c644fdda687a46bcd24816b7f7dc1` | `toast-region.js`, SHA-256 `d3c0bb52a15dffcd7964bc8335ad6a431022662fdfc4e524e00a33d2093facca` |
| Live Wires CSS (example only) | `6e0ff82c222c644fdda687a46bcd24816b7f7dc1` | `main.css`, SHA-256 `45f8ef695178ced30487c5fb1352bbb546273cffc0a04caadc54249d43eac34c` |

The pinned URLs used by the example are the corresponding
`https://raw.githubusercontent.com/Design-Machines-Studio/livewires/<commit>/public/dist/<file>` paths.

The source markup references are `public/reference/components/relative-time.html`
and `copy-button.html` at the first commit, `loading.html` at the loading
commit, and `toasts.html` at the toast commit. Browser lifecycle fixtures are
`interaction-utilities-tests.html`, `loading-tests.html` and
`toasts-tests.html` at those same respective commits. These source references
are parity targets; they are not a claim that this example has been served and
compared in a browser.

Load each module only where its behavior is used. Do not duplicate these
modules' JavaScript state machines or vendor CSS in this library.
