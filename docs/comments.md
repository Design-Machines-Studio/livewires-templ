# Shared comments

`CommentArticleComponent(CommentArticleProps)` renders the ordinary comment
presentation shared by proposal discussions and Chat replies. `CommentArticle`
is its escaped plain-text convenience wrapper:

```templ
@component.CommentArticle("Jane Doe", "Looks good to me.", "2026-03-15", "15 March")

@component.CommentArticleComponent(component.CommentArticleProps{
    ID: "post-815",
    Author: authorLink,
    Avatar: decorativeAvatar,
    AvatarHref: "/members/sam",
    AvatarLabel: "Sam Rivera's profile",
    DateTime: "2026-03-15T10:05:00-07:00",
    Timestamp: "Mar 15 at 10:05",
    Permalink: "#post-815",
    PermalinkLabel: "Mar 15 at 10:05, link to message 815",
    Edited: editedMetadata,
    Footer: actions,
    Class: "text-sm text-fg border-t pt-075 pb-2",
    BodyClass: "text-base mt-075",
    TimestampClass: "text-muted text-xs",
    EditedClass: "text-xs text-muted",
}) {
    <p>Caller-rendered content goes here.</p>
}
```

The lowercase values above are caller-owned `templ.Component` values. See the
maintained [`comments.templ`](../examples/gallery/comments.templ) examples for
complete author, avatar, edited metadata and action compositions.

`Content`, when non-nil, replaces children. Otherwise children render inside
`.comment-body.prose`. There is no HTML-string body property. The convenience
wrapper escapes the author, body and visible timestamp. It does not format the
date. Existing `Comment`, `CommentComponent` and `CommentProps` retain their
original output, escaped `Body`, `DateShort`, avatar, badge and size/alignment
behavior.

## Composition and attributes

- `ID` is caller-owned. Use a unique ID for each anchored item. `Permalink`
  supplies its actual destination; the library does not build routes.
- `Author` supplies inline content, normally an `a.author` or `span.author`.
  The library adds the enclosing `span.comment-author`. Keep controls in the
  footer. Profile links and timestamp links remain separate siblings.
- `Avatar` is optional decorative content. Give images `alt=""` and hide
  initials from assistive technology. With both `AvatarHref` and `AvatarLabel`,
  the library renders a named `a.comment-avatar`. Otherwise it renders an
  unlinked, hidden-from-assistive-technology `span.comment-avatar`. A URL without
  a label deliberately stays unlinked. No avatar produces no avatar wrapper.
- `DateTime` is a valid HTML machine datetime; `Timestamp` is the caller's
  visible text, unchanged. Empty `Timestamp` omits the time/permalink. The
  optional `PermalinkLabel` must begin with the visible timestamp so its
  accessible name includes its visible label. Without it, visible time names
  the link.
- `Edited` supplies inline content inside `p.comment-edited`, including any
  semantic time. `Footer` supplies content inside `footer.comment-footer`.
  For actions, supply a `ul.comment-actions` with native links/buttons in list
  items and an accessible list label. Both wrappers disappear when absent.
- `BodyClass`, `TimestampClass` and `EditedClass` support presentation already
  evidenced by the producer and consumers. `TimestampClass` applies to the
  permalink when linked, otherwise the time. No layout classes are automatic.
- `Class` adds classes to the article. `Attrs` adds developer-controlled outer
  attributes. Never take attribute names from untrusted input or duplicate
  generated `class`/`id` attributes through `Attrs`; generated attributes win
  by first occurrence. Slots must supply valid content for their element.

Rich-content rendering, sanitization, image descriptions and allowed URLs stay
with the caller. The component adds no live region, arrival announcement,
Markdown parser, permission checks, editing, threading or moderation. Governance
owns positions and decisions. Chat owns message order and reading position.

## CSS contract and release

The new API requires the shared comment contract from Live Wires **v2.23.0**,
whose tag resolves to commit
`35c81f71e863db5479d56b513fd35217bb56efc7` ([producer PR #36](https://github.com/Design-Machines-Studio/livewires/pull/36)).
The exact contract is in
[`comments.css`](https://github.com/Design-Machines-Studio/livewires/blob/35c81f71e863db5479d56b513fd35217bb56efc7/src/css/6_components/comments.css)
and the producer's
[manual](https://github.com/Design-Machines-Studio/livewires/blob/35c81f71e863db5479d56b513fd35217bb56efc7/public/manual/components/comments.html).

The emitted hierarchy is `article.comment.comment--layout`, optional
`.comment-avatar`, `.comment-content`, `header.comment-meta` containing
`.comment-author` and `.comment-timestamp` (optionally inside `.comment-permalink`),
`.comment-body.prose`, optional `p.comment-edited` and `footer.comment-footer`.
Caller action lists use `.comment-actions`. Legacy `.comment` alone does not
opt in to the new layout. This library bundles no CSS.

The gallery locks the published repository CSS artifact at that exact commit:
`public/dist/main.css`, 173886 bytes, SHA-256
`d2f7a3d50bcde6bb06f399f0ec0f9dad6d19f7dc5e6961e901234ee73b5b8061`.
This proves immutable public artifact bytes, not a GitHub Release or registry
publication. Consumers must select compatible CSS separately from the Go module.

Proposed next release: **v0.19.0**, an additive minor after current v0.18.0,
because this change exports a new Props type and component functions. No tag or
release is created by this work.

## Review

Use the [gallery runbook](../examples/gallery/README.md) and visit
`/examples/comments` at the URL returned by `run.sh status`. The examples cover
ordinary proposal comments, Chat replies, optional avatars/links, rich content,
edited metadata, actions, multiple anchors and inherited schemes.

Browser checks include desktop, 375px and 320px, 640px reflow, 200% CSS zoom,
light/dark appearance, keyboard focus/order and no-JavaScript reading and links.
CSS zoom and a narrow CSS viewport are not a native browser zoom test. DOM and
accessibility-tree checks do not establish actual screen-reader results.
