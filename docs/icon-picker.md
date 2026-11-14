# Icon picker

`form.IconPicker` renders the Live Wires icon picker: a trigger button and a
native popover panel of radio tiles for one Phosphor Bold icon. The library
renders markup only. It ships no CSS, no JavaScript, no icon catalog and no
Phosphor dependency. Your application chooses the icons, hosts the sprite,
optionally loads the Live Wires enhancer, and validates what the form submits.

Requires **Live Wires v2.29.0** (`d407d3edf652d6f4f754987a0ad9f48aaeb1faa4`).
The markup matches that release's
[reference page](https://github.com/Design-Machines-Studio/livewires/blob/d407d3edf652d6f4f754987a0ad9f48aaeb1faa4/public/reference/components/icon-picker.html)
and follows its
[spec](https://github.com/Design-Machines-Studio/livewires/blob/d407d3edf652d6f4f754987a0ad9f48aaeb1faa4/docs/superpowers/specs/icon-picker.md).

## API

```go
form.IconPickerComponent(form.IconPickerProps{
    ID:       "space",
    Name:     "space_icon",
    Label:    "Space icon",
    Value:    space.Icon, // "" for none
    Sections: sections,
})

form.IconPicker(name, label, value, sections)      // bordered
form.IconPickerGhost(name, label, value, sections) // borderless tiles, ghost trigger
```

| `IconPickerProps` | Meaning |
| --- | --- |
| `ID` | Id prefix. Renders `<ID>-picker` (root), `-panel`, `-heading`, `-search` and `<ID>-<section ID>`. Defaults to `Name`. Sanitized for id use. Use a distinct ID for each picker on a page. Section IDs must be unique and must not be `picker`, `panel`, `heading` or `search`. |
| `Name` | Submitted field name shared by every radio, including Remove. Required. |
| `Value` | Current icon name. Empty checks Remove and shows the placeholder. |
| `Label` | Panel heading, such as `Circle icon`. Defaults to `Icon`. |
| `ChooseLabel`, `ChangeLabel` | Trigger names. Default to `Choose <label>` and `Change <label>`, lowercasing the first letter unless the label starts with an acronym. With a value in `Sections` the trigger reads `Change circle icon, Gavel`, as the enhancer does. |
| `SpriteHref` | Sprite URL without a fragment. Defaults to `/icons/phosphor-bold.svg`. |
| `PlaceholderIcon` | Sprite name shown with no value. Defaults to `selection-plus`. |
| `Ghost` | Adds `.icon-picker--ghost` and swaps the trigger's `.button--subtle` for `.button--ghost`. |
| `HeadingLevel` | Panel heading level, 1–6, default 3. Section headings use the next level, capped at 6, so `6` renders both levels as `h6`. Prefer 1–5. |
| `Sections` | The rendered subset. |
| `Class`, `Attrs` | Extra classes and developer-controlled attributes on the root. |
| `InputAttrs` | Developer-controlled attributes repeated on every radio, including Remove, and on the enhancer's Recent copies. Rendered after the component's own; on a name collision the component's attribute wins. Do not repeat `type`, `name`, `value`, `class` or `checked`, and never pass `id` or other single-element attributes. |

`IconPickerSection` has `ID` (id suffix, default `section-<n>`), `Heading` and
`Icons`. `IconPickerIcon` has `Name`, `Label` and `Tags`:

- `Name` is the Phosphor kebab-case name, such as `users-three`. It is the
  submitted value and the sprite fragment. The library passes it through
  unchanged and does not check kebab case; templ escapes it.
- `Label` defaults to the name with hyphens as spaces and the first letter
  capitalised (`chat-circle-text` → `Chat circle text`), as in the reference.
- `Tags` join with spaces into `data-tags` for the enhancer's search. Phosphor's
  catalog tags make "home" find `house`. Empty tags omit the attribute.

The trigger never renders `aria-expanded`; the enhancer owns it. Search,
Recent and the status region render `hidden`, and the enhancer reveals them.
Exactly one `role="status"` region renders inside the fieldset.

## Subset and server validation

Render only the icons the field allows, and validate the submission against
the **same list**. The radios are ordinary form controls: anyone can submit any
string. An empty value means "no icon".

- Include the current value in `Sections`. A `Value` outside them checks no
  radio, so submitting the form unchanged omits the field entirely; treat a
  missing field as "unchanged", not "cleared". Without JavaScript the trigger
  still shows that icon. With the enhancer it shows the placeholder and the
  Choose label, because the enhancer reads the trigger from the checked radio.
- If a name appears in two sections, only the first copy is checked.
- `Name`, `Value`, section data and icon names are developer-controlled. Never
  build them, or `Attrs`/`InputAttrs`, from request input.

## Sprite and catalog hosting

Tiles reference one `<symbol>` sprite with `<use href="<SpriteHref>#<name>">`.
Serve Live Wires' `public/icons/phosphor-bold.svg` from the same origin, as
`image/svg+xml`, with normal caching; the full set downloads once. Browsers do
not load cross-origin sprites through `<use>`. `phosphor-catalog.json` lists
`name`, `category` and `tags` for building a subset on the server. The library
never reads it. Phosphor Icons are MIT licensed; ship `public/icons/LICENSE`
with the sprite.

## Optional enhancer

Load Live Wires `/dist/icon-picker.js` and call `createIconPicker({ root })` for
each `[data-icon-picker]` root. It adds:

- `aria-haspopup` and `aria-expanded` on the trigger, and positioning next to it;
  below a 30em viewport the panel becomes a full-width bottom sheet.
- Focus on search when the panel opens, with the checked icon scrolled into
  view. Tab then lands on the checked icon. Tab and Shift+Tab cycle Remove,
  search and the grid.
- Arrow keys, Home and End across the visible grid.
- Search by name and tag, a visible count, and `No icons match "…"`.
- Recent, stored in `localStorage` (`recent: false` turns it off).
- Choosing an icon or Remove updates the trigger, closes the panel and returns
  focus to the trigger. Escape and outside clicks close without a change.

Call `dispose()` before removing a root, and `refresh()` after replacing tiles.

## Without JavaScript

The trigger's `popovertarget` opens the native popover, centred. The radios
submit with the form. Search, Recent and the count stay hidden. The trigger
shows the saved value and changes only after the form is saved and the page
re-renders.

### Known v2.29.0 limitation

Without JavaScript, an arrow key that moves onto Remove checks it, which
clears the value. That happens with Left or Up from the first tile, and with
Right or Down wrapping from the last tile. Live Wires then hides the checked Remove
label (`icon-picker.css`), so keyboard focus is lost. With the enhancer, arrow
keys stay within the grid. The fix belongs in the Live Wires CSS.

## Heading levels

The panel heading labels the fieldset and the section headings label their
sections. Choose `HeadingLevel` one below the heading that contains the field.
The default (`h3` panel, `h4` sections) fits a field under a page `h2`.

## Live Wires v2.29.0 changes you inherit

Moving to v2.29.0 changes more than the icon picker:

- **Focus rings everywhere.** Links, inputs, selects, textareas, details,
  switches, checkboxes, tabs, dropdowns, file uploads, date filters, action
  menus and sortable lists move from a 2px `--color-accent` outline to
  `--focus-ring-width` (1px) in `--focus-ring-color`, falling back to
  `--color-muted`, offset by `--focus-ring-offset` (2px). In the gallery the
  ring measures 5.23:1 against the light background and 8.35:1 against the dark
  one. Themes that override `--color-muted` must keep it at 3:1 against their
  backgrounds, or set `--focus-ring-color`.
- **Card heading links** inside line clamps draw their ring at
  `-1 × --focus-ring-width` (was `-2px`).
- **New button variants** `.button--subtle` and `.button--icon`, and a new
  `input.input--small`. These are additive. `component.Button` reaches them
  through `Variant` and `Class`.
- **Widget metrics** gain `container-type: inline-size`, which changes their
  layout slightly.

## Gallery

`examples/gallery` (`/examples/icon-picker`) is demonstration code only:
`icon_picker.go` (submission handler and validation), `icon_picker.templ`
(forms, captions, layout classes and scheme boxes) and `icon_picker_subset.go`
(a static 48-icon subset generated from the pinned catalog). None of it is
part of the library API.
