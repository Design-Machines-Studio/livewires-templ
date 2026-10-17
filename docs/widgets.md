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

Numeric values emit `.widget-metric-value.font-metric.tnum`. `Variant` accepts
`short` (text-4xl), `long` (text-lg), or empty (text-2xl); an unknown variant
uses the default. The producer supplies each variant's line-height and tracking.
Choose a variant after formatting. There is no digit counting, truncation,
scaling or automatic fit loop.

`ValueAttrs` and `ValueClass` affect only the numeric dd. Status and Detail do
not receive them. Set developer-owned CSS hooks there, on the frame, or on an
ancestor:

```go
ValueAttrs: templ.Attributes{
    "style": "--widget-metric-size: 2rem; --widget-metric-weight: 650; --widget-metric-width: 80%;",
},
```

Unset hooks inherit the producer defaults. Width means glyph width, separate
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
| Governance proposal widget header/body/footer and summary values | WidgetComponent and WidgetMetrics/WidgetMetricItem | DTOs, labels, deadlines, permissions and links |
| Atlas dashboard named sections and compact definitions | WidgetComponent and standalone/grouped metrics | Chart/map lifecycle, sharing, denominators, full data tables |
| Existing StatCard | Existing StatCard API, now with value-only font-metric | Existing card layout and formatting |

This work does not migrate consumers. Adoption remains Baseplate #1159 through
#1201, Governance #114, Atlas #60 and Fixture Jig #42, after a compatible Templ
release. Baseplate #1160 owns global font configuration; #1161 owns persisted
per-widget typography, validation, reset and draft/Save behavior. Gallery reloads
prove neither persisted theme configuration nor consumer adoption.
