# External PageHeader consumer

This nested Go module compiles outside the library's module boundary and imports
`component.PageHeaderComponent` and `component.PageHeaderProps`, which proves a
Fixture can render the public API without importing an `internal` package.

The local `replace` directive is only for this repository's producer check. A
consumer should pin the published `livewires-templ` version after this change is
released. No Live Wires CSS or JavaScript is included here; the host supplies
its compatible Live Wires assets and caller-owned layout classes.

## Content and settings headers

`Title` remains escaped text. Use `TitleSuffix` for inline content inside the
H1, such as a member-count span with an ID used by live updates. The component
adds a space before the suffix. Keep block content out of this slot.

`Supplement` stays below the lead in the title column. Use `Footer` for a
separate row after the entire title/action layout, such as Atlas's lead and
last-synced line. Leave `Lead` empty when that row supplies its own lead.
The caller owns the row's markup and classes. Both slots render without actions.

| Existing header | Public props |
| --- | --- |
| Content/settings title and introduction | `Title`, `Lead`, `Class` |
| Member directory with a live count | `Title`, `TitleSuffix`, `HeadingAttrs` |
| Dashboard or list with actions | `Actions`, `LayoutClass`, `ContentClass` |
| Notifications with context beside actions | `Supplement` within the title column |
| Atlas title/actions and separate lead/sync row | `Actions`, `LayoutClass`, `Footer` |
| Governance proposal summary and metadata | `Title`, `Lead`, `Supplement`, `HeadingAttrs` |

Profile and organization identity headers remain separate components. Page
navigation, breadcrumbs, forms, permissions and transport stay with the caller.
`Attrs` and `HeadingAttrs` accept developer-controlled attributes only; never
construct attribute names from request or member input.

