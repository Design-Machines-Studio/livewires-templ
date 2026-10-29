# Widgets

Widgets share presentation across Baseplate, Governance and Atlas. This package
emits HTML only. Use Live Wires **v2.24.0**, commit
`3042ffe56c36b8e51e4e8a98f500b1aa04a43fd3`, or a later compatible release
containing the same widget and `font-metric` contract. The gallery retains its
existing exact asset lock. No extra stylesheet or JavaScript is required.

The authority is [the published widget contract](https://github.com/Design-Machines-Studio/livewires/blob/3042ffe56c36b8e51e4e8a98f500b1aa04a43fd3/docs/widget-references.md),
including its exact consumer inventory and source comparisons. The producer
merged in [Live Wires PR #31](https://github.com/Design-Machines-Studio/livewires/pull/31).

```templ
@component.WidgetComponent(component.WidgetProps{
    ID: "membership", Title: "Membership at a glance",
    HeadingID: "membership-title", HeadingLevel: "h2",
    Description: "Current shared observations", Footer: memberDirectoryLink(),
}) {
    @component.WidgetMetrics("") {
        @component.WidgetMetricItemComponent(component.WidgetMetricProps{
            Label: "Active", Value: "42", Variant: "short",
        })
        @component.WidgetMetricItemComponent(component.WidgetMetricProps{
            Label: "Probationary", Value: "0", Variant: "short",
        })
        @component.WidgetMetricItemComponent(component.WidgetMetricProps{
            Label: "Private membership", Status: "Permission denied",
        })
    }
}
```

`Widget(title, headingID)` is the h2 convenience wrapper. `WidgetComponent`
renders `.widget`, optional `header.widget-header`, `div.widget-content` with
children, and optional `footer.widget-footer`. The built-in header contains
`.widget-heading`, a `.widget-title` and optional `.widget-description`.
`HeadingLevel` accepts h1–h6 and defaults to h2. Choose it for the surrounding
page's heading hierarchy.

`Action` and `Footer` are caller-supplied `templ.Component` values. Put
`widget-action` on the action's root element; an anchor or button keeps its
native semantics. `Header` replaces the complete built-in header composition,
including Title, Description and Action. It renders inside `header.widget-header`.
`PlainHeader` emits `widget--plain-header`, which removes the divider and its
padding. It has no route-specific behavior.

`ID` names the frame. `HeadingID` names the generated heading and connects it
through `aria-labelledby`. Callers supply unique, stable IDs; this library does
not create them. A frame becomes a section when a generated title has a
HeadingID, or when `Attrs` supplies a nonempty `aria-label` or `aria-labelledby`.
Other frames use a div. A custom header must supply its naming relationship
through `Attrs`; its heading IDs remain caller-owned. A headingless body needs
no landmark unless the caller deliberately names it.

Every Props type exposes `Class` and `Attrs`. `HeadingAttrs` belongs to the
generated heading, `ContentAttrs` to `.widget-content`, and `FooterAttrs` to
the optional footer. No header/description attribute map is added. Attributes
are spread after component-generated attributes; on collision, the component's
first occurrence wins. Do not repeat generated class, ID or naming attributes
in these maps. Use Class and the explicit ID props instead. External labels
must resolve to real elements.

Attribute names and inline styles are **developer-controlled**. Templ escapes
values but does not restrict attribute names. Never build maps or CSS from
arbitrary user strings. Text fields accept ordinary untrusted text and escape it.

## Metrics

| API | Anatomy | Use |
| --- | --- | --- |
| `WidgetMetricComponent(WidgetMetricProps)` / `WidgetMetric(label, value)` | `dl.widget-metric > dt + dd` | Standalone metric |
| `WidgetMetricsComponent(WidgetMetricsProps)` / `WidgetMetrics(class)` | `dl.widget-metrics` with children | Group container |
| `WidgetMetricItemComponent(WidgetMetricProps)` / `WidgetMetricItem(label, value)` | `div.widget-metric > dt + dd` | Only inside the group dl |

Each group keeps its own term and definitions. Do not nest standalone lists in
WidgetMetrics or render WidgetMetricItem without a definition list. Metrics
have no card surface. Label is the dt; the primary definition and optional
Detail use dd. Terms come first in reading order. CSS moves numeric values
above their terms visually, matching the producer.

Value is the complete caller-formatted string. Currency, locale, signs, units,
precision and authorization stay with the application. `"0"` is a confirmed
numeric value. `Status` takes precedence over Value and uses
`.widget-metric-status` with ordinary typography. Empty Value also uses that
status presentation; supply explicit Status text for missing observations.
No missing value becomes a fabricated zero.

Numeric values emit `.widget-metric-value.font-metric.tnum`. Empty/auto and
unknown variants select the automatic sizing contract below. Explicit `short`
(text-4xl) and `long` (text-lg) variants retain their existing behavior. The
component counts formatted Unicode code points only; it does no numeric parsing.

`ValueAttrs` and `ValueClass` affect only the numeric dd. Status and Detail do
not receive them. Set developer-owned CSS hooks there:

```go
ValueAttrs: templ.Attributes{
    "style": "--widget-metric-size: 2rem; --widget-metric-weight: 650; --widget-metric-width: 80%;",
},
```

Prefer `ValueAttrs`. The hooks are inherited custom properties, and the
producer's `.font-metric` utility reads them, so a hook set on the frame or an
ancestor also resizes every `font-metric` value inside it, including a StatCard
value without a `ValueSize`. Unset hooks inherit the producer defaults. Width means glyph width, separate
from dashboard columns. The producer maps hooks to standard `font-size`,
`font-weight` and `font-stretch`; supported registered variable-font axes use
the font's declared range. Static fonts remain readable without a width axis.
A computed font-stretch value alone does not prove a change in glyph width.
No font is bundled, selected or downloaded by this API.

`.font-metric` uses `--font-metric`, falling back to `--font-body` and sans-serif.
The optional generic `--font-metric-size` / `--font-metric-width` initial settings
remain producer-owned. Widget-specific size/width hooks take precedence. With
them unset, `revert-layer` keeps underlying component presentation. Do not add
competing `text-*` or weight utilities to hook-controlled values. Size overrides
retain the selected variant's leading and tracking; use matching variants or
set line-height and letter-spacing on the value when needed.

The existing StatCard now adds `font-metric` only to its value. Its Props,
wrappers, sizing, links, status and progress behavior stay unchanged. Card and
EmptyState stay unchanged. Ordinary table numbers, labels, descriptions and
runtime version strings keep ordinary typography.

## Quiet messages and composition

`WidgetMessageComponent(WidgetMessageProps)` renders `.widget-message` with an
optional decorative Icon inside `.widget-message-icon[aria-hidden=true]`, then
`.widget-message-content` containing Message and optional
`.widget-message-detail`. `WidgetMessage(message)` is the simple wrapper.
Icons are caller-supplied components. They must contain no focusable controls;
put meaningful information in Message. No surface or live-region role is added.
Callers decide whether a dynamic replacement needs an announcement.

Use the content slot for ordinary lists, existing Timeline, StatusIndicator and
ProgressBar components, or a caller chart renderer. A static illustration or
unavailable message needs no chart runtime. Keep an independent complete table
with a caption, column scopes and row scopes. `table.widget-table` preserves
native columns at narrow widths. For wider data, use a named, focusable
`.overflow-x-auto` region and an explicit table minimum width. Do not combine
this variation with `table--lined` or remove columns to fit.

The maintained [gallery](../examples/gallery/README.md) at `/examples/widgets`
uses these public components. The exact producer reference is retained as a
test fixture for structural and rendered comparison, with its site-wide
includes and scripts outside the comparison.

## Consumer mapping and ownership

| Existing owner markup | Shared replacement | Remains caller-owned |
| --- | --- | --- |
| Baseplate `DashboardWidgetBody` | WidgetComponent, Footer link slot | Placement, account links, registry, grants and editor |
| Baseplate `WidgetMessage` | WidgetMessageComponent with a caller icon component | Icon library and dynamic announcement policy |
| Linked notification rows | Caller-written `ul > li.widget-message` with the same icon/content/detail classes | Links, times and list markup |
| Governance proposal widget header/body/footer and summary values | WidgetComponent and WidgetMetrics/WidgetMetricItem | DTOs, labels, deadlines, permissions and links |
| Atlas dashboard named sections and compact definitions | WidgetComponent and standalone/grouped metrics | Chart/map lifecycle, sharing, denominators, full data tables |
| Health facts (version strings, status indicators, linked details) | Caller-written `div.widget-metric > dt + dd` inside WidgetMetrics; plain `dd` without `font-metric` | Version text, StatusIndicator and links |
| Existing StatCard | Existing StatCard API, now with value-only font-metric | Existing card layout and formatting |

This work does not migrate consumers. Adoption remains Baseplate #1159 through
#1201, Governance #114, Atlas #60 and Fixture Jig #42, after a compatible Templ
release. Baseplate #1160 owns global font configuration; #1161 owns persisted
per-widget typography, validation, reset and draft/Save behavior. Gallery reloads
prove neither persisted theme configuration nor consumer adoption.

## Automatic metric sizing (unreleased)

Use `widget-metric-value--auto` with `data-widget-metric-length` set to the
upper formatted-length band: 4, 8, 12, 20, 32 or 64. The Templ component emits
both for its default/auto variant. Count Unicode code points in the complete
formatted value, including separators, currency, signs and units; do not parse
or reformat the number. Values longer than 32 use band 64 and may wrap.

The metric definition list/group establishes its own inline-size container,
so grouped metrics use each cell's available width rather than the entire card.
CSS estimates room at 0.7em per character in the band, caps the default at
`--text-3xl`, and keeps a 1rem floor. The same value grows in a wider cell;
longer strings receive a smaller size in the same cell. Browsers without
container units use text-3xl and ordinary wrapping.

This is a conservative layout estimate, not glyph measurement. Wide fonts,
full-width digits, combining sequences and large values can wrap. Every
character stays in the DOM: no precision loss, ellipsis, transform scaling,
JavaScript, storage, font downloads or measurement loop. Zoom and font changes
use normal CSS layout. Labels, descriptions and unavailable states stay
unchanged. An explicit widget size, generic metric size or size utility wins;
weight and font-width overrides retain their existing behavior. Explicit
short/long variants remain supported as deliberate fixed art direction.

This contract requires the paired Live Wires CSS and Templ changes to be
published and adopted together. Existing Baseplate/Atlas releases are unchanged.
After adoption, consumers should remove their default text-size utility so the
shared automatic rule can apply; explicit personal size choices still win.

Paired CSS implementation: [Live Wires PR #41](https://github.com/Design-Machines-Studio/livewires/pull/41).
