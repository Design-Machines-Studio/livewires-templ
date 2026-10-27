# Inline mention suggestions

`component.MentionSuggestionsComponent` renders a suggestion root, listbox and
one status region. `component.MentionOptionComponent` renders a person row.
These components render caller state. They contain no controller, CSS, directory
lookup, filtering, permissions, callbacks, persistence or notifications.

Load **Live Wires v2.27.0 or later**, including both `main.css` and the optional
`combobox.css`. Mention styles ship in `combobox.css`. The producer contract is
[mention-suggestions.md at v2.27.0](https://github.com/Design-Machines-Studio/livewires/blob/1c52666402c0ec2de4ed64736ea679a8a8573838/docs/superpowers/specs/mention-suggestions.md).

## Render API

```templ
@component.MentionSuggestionsComponent(component.MentionSuggestionsProps{
    ID: "reply-mentions", Label: "People", State: "ready",
    Status: "2 people found", Dir: "ltr",
    Attrs: templ.Attributes{"style": "--mention-suggestions-top: 212px; --mention-suggestions-left: 48px"},
}) {
    @component.MentionOptionComponent(component.MentionOptionProps{
        ID: component.MentionOptionID("reply-mentions", person.ID),
        Name: person.Name, Detail: person.CircleName, Active: true,
        Avatar: &component.AvatarProps{Name: person.Name, Src: person.AvatarURL},
    })
    @component.MentionOption(component.MentionOptionID("reply-mentions", other.ID), other.Name, "", false)
}
```

`MentionSuggestionsProps.ID` identifies the **listbox**, not the root. Supply a
nonempty, page-unique ID and a localized `Label`. `State` accepts `closed`,
`loading`, `ready`, `empty` and `error`; zero or unknown values render `closed`.
`Busy` emits `aria-busy="true"` on the list. `Status` is caller-supplied text.
`Dir` defaults to `ltr`; pass the editor's direction. `Class` and `Attrs` apply
to the root. Use `Attrs["id"]` if the root also needs an ID. Children supply rows.

`MentionOptionProps` supplies `ID`, `Name`, optional `Detail`, optional `Avatar`,
`Active`, `Disabled`, `Truncate`, `Class` and `Attrs`. Detail must be a real
identifier supplied by the caller, such as a circle, role or email. Names never
produce handles. Two people with the same name need supplied details to look
different. Without details they look alike.

`Avatar` uses only `AvatarProps.Name` and `Src` as data. A nil pointer omits it;
an empty source renders initials, using the row name when the avatar name is
empty. Images have `alt=""`; the avatar span has `aria-hidden="true"`. Size,
shape, `ShowName`, `Class` and `Attrs` from AvatarProps are intentionally unused.
The existing `AvatarComponent` emits different role/alt markup and is unchanged.

`Active` adds `data-active`, `aria-selected="true"` and
`rounded-sm scheme-subtle` together. Inactive rows emit `aria-selected="false"`
and omit those classes. The caller supplies at most one active enabled row;
disabled rows must never be active. `Disabled` emits `aria-disabled="true"`.
`Truncate` adds `truncate` to the name and any detail; their complete text stays
in the DOM. Rows have no tab stop. The CSS supplies their grid layout; callers
can add classes through `Class` without changing the library's layout policy.

The wrappers are `MentionSuggestions(id, label, state, status)` with children
and `MentionOption(id, name, detail, active)`.

`Attrs` is developer-controlled, never copied from user input. Root attributes
apply only to the root; option attributes apply to the `li`, after component
attributes. The component's first occurrence wins on a collision in HTML.
Do not repeat component attributes in `Attrs`; use Props and `Class`. In
particular, do not inject active-state classes or ARIA attributes through it.

## Stable IDs

`MentionOptionID(listboxID, personID)` returns `listboxID + "-" + personID`,
matching the producer script: `reply-mentions` and `sam-2` become
`reply-mentions-sam-2`. Pass canonical, nonempty, whitespace-free IDs. The caller
must keep listbox IDs unique across the page and person IDs unique in each list,
including after concatenation. The helper performs no lossy normalization,
random generation or name-based slugging. Keep IDs unchanged across searches.

## States and announcements

| State | Rows | Caller status |
| --- | --- | --- |
| `closed` | Hidden, children omitted | Empty |
| `loading`, `Busy: false` | Hidden, children omitted | Searching people… |
| `loading`, `Busy: true` | Retained stale children, `aria-busy="true"` | Searching people… |
| `ready` | Children visible; first enabled row active | Result count, e.g. 5 people found |
| `empty` | Hidden, children omitted | No matching people |
| `error` | Hidden, children omitted | Caller retry message |

Callers can localize the text. CSS clips the closed root while leaving it
mounted, visually hides the ready count and shows loading/empty/error text.
Keep the root, list and its single `p[role="status"][aria-atomic="true"]` mounted
for the editor's lifetime. Update status text only when its message changes.
Delay the loading announcement about 400 ms so quick replies announce only the
count. Repeated typing with the same count must not rewrite the status node.

When retaining stale rows, pass `State: "loading", Busy: true` and the previous
results with every row inactive. Remove the editor's active descendant before
waiting. The consumer may allow click/tap selection of explicitly named stale
rows. Discard older request replies before rendering. `Busy` does not cancel
requests or clear row state for the caller.

## Editor and consumer contract

The library does not render an editor. Baseplate must wire these attributes on
the Milkdown contenteditable:

| Attribute | When |
| --- | --- |
| `role="textbox"`, `aria-multiline="true"`, accessible name | Always |
| `aria-autocomplete="list"`, `aria-haspopup="listbox"` | While attached |
| `aria-controls="<listbox ID>"` | While attached |
| `aria-activedescendant="<active option ID>"` | Only while open with an active row |
| `aria-expanded` | Never |

Remove `aria-activedescendant` when closed, loading, stale, empty or error.
ARIA 1.2 does not permit `aria-expanded` on `textbox`. Changing a multiline rich
text editor to `combobox` would lose its multiline textbox semantics. The
status count and active descendant convey the open state instead.

Keep focus and selection in the editor. Its highest-priority ProseMirror keymap
owns arrows (wrap and skip disabled rows), Enter/Tab selection, Escape dismissal,
modifier keys and composition pass-through. While open, Enter always prevents
default, including with modifiers and without an active row. Every send/submit
handler must respect `defaultPrevented`, `isComposing` and legacy keyCode 229.
Shift+Tab passes through; Tab without an active row closes and passes through.
With no visible rows, arrows pass through. Escape keeps the token dismissed
until a new trigger or until the caret leaves and returns.

Cancel popup `mousedown` to preserve focus; select on `click` for mouse, touch
and pen. Hover does not move the active row. Scroll only the list's `scrollTop`
to reveal the active option; `scrollIntoView` can move the page.

Consumers measure the caret, choose above/below placement and set
`--mention-suggestions-top` **or** `--mention-suggestions-bottom`, plus physical
`--mention-suggestions-left` and optionally `--mention-suggestions-max-block`.
Clamp to the viewport and nearest scroll container, reposition on scroll/resize,
and align the right edge in RTL. Mount outside the editor, in the same dialog or
top-layer context, outside transformed/filtered/layout-contained ancestors.

Baseplate owns Milkdown triggers, token rules, caret geometry, keymap, mention
insertion and serialization. Chat owns access-filtered search, request identity
and consumer acceptance. `[[` record references can use the same rows with a
different label; implementing them is outside this API adoption.

## Gallery and verification limits

`/examples/mention-suggestions` demonstrates the public components, static states,
two independent editors, RTL and a narrow container. `mention-demo.js` adapts
the producer fixture for this gallery. `/demo/mentions` returns Templ-rendered
fragments from fictional fixture data. The adapter keeps the mounted list and
status, swaps row fragments, and updates status only when its text changes.
All demo token parsing, searching, timing, positioning and plain-text insertion
stay in `examples/gallery`; nothing is saved and no directory exists.

Render and Playwright tests check markup, DOM/accessibility-tree relationships,
keyboard/pointer/touch emulation, state, geometry and status mutations. They do
not establish actual screen-reader support or platform touch acceptance.
`aria-activedescendant` on contenteditable varies by screen reader, particularly
VoiceOver. Test the real Milkdown integration with VoiceOver/Safari and
NVDA/Firefox before claiming consumer assistive-technology acceptance.

This adds exported types and functions over v0.23.1: the proposed next release
is **v0.24.0**. No release is published by this change.
