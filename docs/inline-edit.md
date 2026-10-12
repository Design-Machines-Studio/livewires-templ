# Inline editing wrappers

`form.InlineEditCompact` and `form.InlineEditDisclosure` render the native markup
hooks used by Live Wires inline editing. Their `Props` keep values, IDs, labels,
form attributes and input attributes explicit. These components do not implement
the interaction controller or save data.

Each instance needs a unique `ID` for its HTML IDs and a stable application-owned
`ContextID` for the producer controller. Pass the application's `action`,
`method` and CSRF field. `CSRF` accepts a rendered `templ.Component`, so the
application can use its existing CSRF component. `FormAttrs` land on `<form>`;
`InputAttrs` land on the actual `<input>`; `Attrs` land on the editor root.
Only pass developer-controlled attributes to these maps.
`CancelLabel` labels the controller's cancel button. When `CancelHref` is set,
`CancelLinkLabel` labels the native return link; if omitted, it reads “Return
without saving” so the link stays distinct from the button.

`TriggerLabel` is the visible edit control text. In the compact variant,
`EditLabel` supplies the link's accessible name, so a short visible “Edit” can
have a more specific accessible label. Disclosure values use a neutral
paragraph element; consumers can place the component under their own heading
without forcing an inappropriate heading level.

Choose `Class`, `RowClass`, `FormClass` and `ActionsClass` values in the consuming
application. The components do not add layout primitives. The compact form
starts visible: without JavaScript, its Edit link targets the form and the
browser submits it normally. The disclosure uses native `<details>` and
`<summary>` around the same kind of form. An application can provide
`CancelHref` for a no-JavaScript back or cancel destination; the Cancel button
is an enhancement hook for the optional controller. The compact form groups
its Save, Cancel and optional return link in an actions `<div>` that receives
`ActionsClass`; the producer reference places these controls directly in the
form. This wrapper is the class seam for consumers who want to group the
actions.

## Optional producer assets

The producer CSS and JavaScript are optional external assets. The bounded
example pins them, and the base stylesheet, to merged Live Wires commit
[`cf6ef2e4627336c09a84825e50d6c97a1b265aff`](https://github.com/Design-Machines-Studio/livewires/commit/cf6ef2e4627336c09a84825e50d6c97a1b265aff).
Consumers should follow the producer's [inline-edit contract](https://github.com/Design-Machines-Studio/livewires/blob/cf6ef2e4627336c09a84825e50d6c97a1b265aff/docs/inline-edit.md)
and only include the assets they use. Templ emits markup hooks only; it has no
CSS or JavaScript dependency.

## Application result contract

The application creates one controller for each stable editor root and listens
for its commit event. A result is valid only when `contextId`, `editId` and
`requestId` all match the active request. Treat these as correlation IDs, never
as credentials. Foreign, duplicate, cancelled, disposed, stale and late
responses must not update the editor. A matching later result may settle an
`unknown` request while that request remains active; late results for a request
that has been replaced or disposed must be ignored. The application must also
protect persistence from stale writes; local response matching cannot resolve
server conflicts.

Return one of the producer's outcomes:

- **confirmed** only after the application has saved and can return the
  authoritative text;
- **rejected** when the application knows the change did not apply, with a
  message while retaining the draft for correction;
- **unknown** when a transport failure or timeout cannot establish whether the
  save applied. A matching later result may still settle an unknown request.

Cancellation ends local editing or observation. It does not cancel an
application request or guarantee rollback; the server may still have applied
the save. The application owns persistence, authorization, validation,
conflict resolution and retry policy. Do not add those behaviors to these
wrappers or duplicate the producer's controller.

The controller retains the draft on blur; leaving the field never saves or
cancels. Enter submits a single-line draft, while Escape cancels. The producer
guards composition events, `isComposing`, legacy key code 229 and the immediate
composition-end task so IME text is not accidentally submitted. A simulated
composition event is not a real IME test; verify real IME behavior separately.

Keep the editor root and required child nodes stable for the controller's
lifetime. The controller keeps the input node in place and avoids rewriting an
unchanged value so the browser can retain focus and selection. Application DOM
morphs must preserve the focused input node, its live value and selection. Use
the producer's `updateConfirmed` method for external values during an active
draft; do not replace the focused editor subtree. Removing a required node
disposes the controller, and reinsertion requires fresh setup.
