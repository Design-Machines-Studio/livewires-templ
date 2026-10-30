# Live Wires Templ

Go/Templ component library for [Live Wires](https://github.com/Design-Machines-Studio/livewires) CSS framework. Emits semantic HTML with Live Wires class names. **No CSS bundled** — CSS is a peer dependency.

## Install

```sh
go get github.com/Design-Machines-Studio/livewires-templ
```

## Usage

```go
import (
    "github.com/Design-Machines-Studio/livewires-templ/component"
    "github.com/Design-Machines-Studio/livewires-templ/form"
    "github.com/Design-Machines-Studio/livewires-templ/layout"
    lw "github.com/Design-Machines-Studio/livewires-templ"
)
```

### Components

The sortable API composes `SortableListComponent`, `SortableItemComponent`,
`SortableHandleComponent` and `SortableMoveComponent`. Supply an `ol` or semantic
table, stable item IDs, native form endpoints and application-owned CSRF/revision
fields. Attribute keys must remain developer-controlled. Optional Live Wires
sortable CSS/JavaScript are peer assets; the library handles no persistence or
authorization. See the [persistent Go/Templ/Datastar example and serving runbook](examples/sortable/README.md)
for two lists, native arrows, explicit move confirmation and isolated serving.

`SaveStatusComponent` renders one persistent `role="status"` region containing
the short `.save-status` line and optional supporting text. The host supplies
explicit state and any scope-aware label; unknown or invalid states render
“Save not confirmed.” The component does not save, publish, or manage editor
state. Keep actions and surrounding layout outside the region. See the
[maintained save-status gallery example](examples/gallery/README.md).

```go
@component.SaveStatusComponent(component.SaveStatusProps{
    ID:             "proposal-save-status",
    State:          "saved",
    Label:          "Text saved",
    SupportingText: "Rules and collaborators still need saving.",
})
```

```go
// Button
@component.Button("Save", "accent")

// Badge
@component.Badge("Active", "green")

// Avatar (sizes: "xs", "sm", "lg", "xl", "2xl", "3xl", "4xl", or "" for default)
@component.Avatar("Jane Doe", "")
@component.AvatarSmall("Jane Doe")
@component.AvatarImage("/img/avatar.jpg", "Jane Doe", "lg")

// Avatar with name displayed alongside (wraps in cluster layout)
@component.AvatarComponent(component.AvatarProps{Name: "Jane Doe", ShowName: true, Size: "sm"})
@component.AvatarComponent(component.AvatarProps{Name: "Jane Doe", Src: "/img/avatar.jpg", ShowName: true})

// Square avatar (composable with any size)
@component.AvatarComponent(component.AvatarProps{Name: "Acme Co", Square: true, Size: "sm"})

// Aspect ratio avatar (adds aspect-*, flex centering, and h-auto classes)
@component.AvatarComponent(component.AvatarProps{Name: "Logo", Src: "/img/logo.png", Size: "xl", AspectRatio: "4-3"})

// Card (scheme controls color: "subtle", "dark", "accent", or "")
@component.Card("subtle") {
    <p>Card content</p>
}

// Clickable card link
@component.CardLink("/details", "")

// Stat card
@component.StatCardSimple("Revenue", "$12,450", "+12%", "")

// Stat card with progress bar
@component.StatCardComponent(component.StatCardProps{
    Label:    "Storage used",
    Value:    "72%",
    Progress: 72,
})

// Info card
@component.InfoCard("Board Meeting", "Upcoming")

// Activity card
@component.ActivityCard("Proposal A", "/proposals/1")

// Card footer (any card) and summary card (Live Wires v2.26.0, see docs/card-summary.md)
@component.CardFooter("cluster cluster-between items-center text-sm") {
    <span>Footer</span>
}
@component.CardSummaryComponent(component.CardSummaryProps{Title: "Board room", Href: "/spaces/board", Scheme: "subtle", Class: "rounded", HeaderClass: "stack"}) {
    @component.CardFooter("cluster cluster-between items-center text-sm") { ... }
}

// Empty state
@component.EmptyState("No items found")

// Comment
@component.Comment("Jane Doe", "Looks great!", "2024-03-15", "Member")

// Checklist
@component.Checklist(items)

// Progress checklist
@component.ProgressChecklist("Setup", 3, 5, items, "3 of 5 complete")

// Kanban board
@component.KanbanBoard("board-1") {
    @component.KanbanColumn("To Do") {
        @component.KanbanCard("/task/1") {
            <p>Task content</p>
        }
    }
}

// Dialog
@component.Dialog("Confirm") {
    <p>Are you sure?</p>
}

// Dialog with custom heading level
@component.DialogComponent(component.DialogProps{
    Title:        "Edit Item",
    HeadingLevel: "h3",
})

// Toast
@component.Toast("Saved!", "success")

// Panel dropdown (details.dropdown--panel, form.menu, scrollable .body)
@component.DropdownPanelComponent(component.DropdownPanelProps{
    Label:       "Status",
    Footer:      applyButtons,
    FooterClass: "cluster cluster-between",
}) {
    ...body content...
}

// Icon-only panel trigger: Label is visually hidden text, no chevron
@component.DropdownPanelComponent(component.DropdownPanelProps{
    Label:        "Notifications",
    Icon:         bellIcon,
    IconOnly:     true,
    TriggerClass: "button--subtle",
    TriggerAttrs: templ.Attributes{"aria-controls": "notifications-body"},
    MenuTag:      "div",
    Variant:      "end",
}) {
    ...body content...
}

// Checkbox filter panel
@component.CheckboxFilter("Circles", "circle", options)

// Checkbox filter with summary, footer, and consumer-owned layout classes
@component.CheckboxFilterComponent(component.CheckboxFilterProps{
    Label:       "Circles",
    Summary:     "All",
    GroupLabel:  "Filter by circle",
    Name:        "circle",
    CountLabel:  "members",
    Options:     options,
    MenuTag:     "div", // client-side filtering or nested in another form
    Footer:      clearButton,
    FooterClass: "cluster cluster-between",
})

// Calendar month grid (six weeks, Monday first)
@component.CalendarComponent(component.CalendarProps{
    Month:      month,
    Today:      today,
    RangeStart: from,
    RangeEnd:   to,
    DayAttrs:   func(d time.Time) templ.Attributes { return templ.Attributes{"data-date": d.Format("2006-01-02")} },
})

// Two months side by side, navigation on the outer edges
@component.CalendarGroup(component.CalendarProps{Month: month, Today: today}, 2)

// Date filter: compose the body from Presets, DateRangeInputs, and calendars
@component.DateFilterComponent(component.DateFilterProps{
    Label:       "Jan 8 – Jan 14, 2027",
    Sidebar:     true,
    Footer:      applyButtons,
    FooterClass: "cluster cluster-between",
}) {
    @component.Presets("range", presets)
    @component.CalendarGroup(component.CalendarProps{Month: month, RangeStart: from, RangeEnd: to}, 2)
}
@component.DateRangeInputsComponent(component.DateRangeInputsProps{StartName: "from", EndName: "to", Size: "small"})

// Full props control
@component.ButtonComponent(component.ButtonProps{
    Variant:  "accent",
    Type:     "submit",
    Disabled: false,
    Class:    "mt-1",
    Attrs:    templ.Attributes{"data-on-click": "$post('/save')"},
})
```

### Forms

```go
// Text field
@form.FieldText("Name", "name", "", true)

// Select
@form.SelectSimple("Country", "country", "Choose...", options, false)

// Optional searchable select with native fallback. Applications load and
// initialize the optional Live Wires assets separately.
@form.ComboboxComponent(form.ComboboxProps{
    ID: "timezone", RootID: "timezone-choice",
    Label: "Timezone override", Name: "interface_timezone",
    Placeholder: "Use organization default",
    Options: []form.SelectOption{
        {Value: "Asia/Makassar", Label: "Makassar"},
        {Value: "Europe/Paris", Label: "Paris"},
    },
    Hint: "Leave empty to inherit the organization timezone.",
    InputAttrs: templ.Attributes{"data-testid": "timezone-select"},
})

// Remote mode only marks the wrapper for an application-owned adapter.
@form.ComboboxRemote(form.ComboboxProps{
    ID: "country", RootID: "country-choice", Label: "Country", Name: "country",
    Options: countries,
})

// Checkbox
@form.CheckboxSimple("agree", "I agree to the terms", false)

// Checkbox with a hint
@form.Checkbox(form.CheckboxProps{
    Name:  "agree",
    Label: "I agree to the terms",
    Hint:  "You can opt out later",
})

// Input-level attributes land on the nested <input>, not the wrapper label.
@form.Checkbox(form.CheckboxProps{
    Name:  "agree",
    Label: "I agree",
    InputAttrs: templ.Attributes{
        "data-on:change": "@post('/agree')",
    },
})

// Toggle switch with a hint
@form.SwitchComponent(form.SwitchProps{
    Name:  "notifications",
    Label: "Enable notifications",
    Hint:  "We send a weekly digest",
})

// Autosave toggle — Datastar directives live on InputAttrs, not the label.
@form.SwitchComponent(form.SwitchProps{
    Name:    "autosave",
    Label:   "Autosave",
    Checked: settings.Autosave,
    InputAttrs: templ.Attributes{
        "data-bind":                 "autosave",
        "data-on:change":            "@post('/settings/autosave')",
        "data-indicator:autosavePending": "",
        "data-attr:disabled":        "$autosavePending ? true : false",
    },
})

// No visible label — provide an accessible name via InputAttrs.
@form.SwitchComponent(form.SwitchProps{
    Name: "notifications",
    InputAttrs: templ.Attributes{
        "aria-label": "Enable notifications",
    },
})

// Radio
@form.RadioSimple("plan", "pro", "Pro", false)

// Radio group
@form.RadioGroup(form.RadioGroupProps{
    Name:     "plan",
    Selected: "pro",
    Options: []form.RadioOption{
        {Value: "free", Label: "Free"},
        {Value: "pro", Label: "Pro"},
    },
})

// Radio group with per-option hints. Each hint renders as supporting text
// beneath its label and is wired to that input via aria-describedby.
@form.RadioGroup(form.RadioGroupProps{
    Name:     "plan",
    Selected: "pro",
    Options: []form.RadioOption{
        {Value: "free", Label: "Free", Hint: "No credit card required"},
        {Value: "pro", Label: "Pro", Hint: "Billed annually"},
    },
})

// Single radio with a hint
@form.Radio(form.RadioProps{
    Name:  "plan",
    Value: "pro",
    Label: "Pro",
    Hint:  "Billed annually",
})

// Toggle switch
@form.Switch("notifications", "Enable notifications", true)

// Compact toggle switch
@form.SwitchSmall("compact-mode", "Compact mode", false)

// Textarea
@form.TextareaSimple("Message", "message", "", "Write...", "")

// Filter dropdown
@form.Filter(form.FilterProps{Title: "Status", Name: "status", Options: options})

// Date range
@form.DateRange(form.DateRangeProps{Label: "Date range", StartName: "from", EndName: "to"})

// Field with error
@form.Field(form.FieldProps{
    Label: "Email",
    Name:  "email",
    Type:  "email",
    Error: "Please enter a valid email",
})

// Hints are announced, not just displayed. Every form control wires its Hint
// to the input via aria-describedby, alongside the error message when both
// are present. Ids are derived from the field name; a name that is not
// id-safe ("user[email]") is sanitized for the id and left intact for the
// submitted name.
@form.Field(form.FieldProps{
    Label: "Email",
    Name:  "email",
    Type:  "email",
    Hint:  "We never share it",
})

// File upload
@form.FileUploadSimple("resume", "Upload resume")

// File upload drop zone
@form.FileUploadZone("docs", "Browse files", "or drag and drop files here")

// File upload with full props
@form.FileUpload(form.FileUploadProps{
    Name:     "images",
    ButtonText: "Upload images",
    Text:     "PNG, JPG, or SVG",
    Variant:  "accent",
    Zone:     true,
    Accept:   "image/*",
    Multiple: true,
})

// FileUpload exposes InputAttrs for directives on the actual <input type="file">.
@form.FileUpload(form.FileUploadProps{
    Name:       "avatar",
    ButtonText: "Upload avatar",
    InputAttrs: templ.Attributes{
        "data-on:change": "@post('/upload')",
    },
})
```

`ComboboxProps.Attrs` lands on the `data-combobox` wrapper. `InputAttrs` lands
on the native `<select>`. The optional producer creates a separate search input
without a name, so typed query text never becomes a submitted value. Load the
matching optional `combobox.css` and `combobox.js` in the application and
initialize Live Wires' `createCombobox` API. The application owns remote
queries, authorization, validation and saving.

If `Value` is not in `Options`, the component adds it as the selected native
option so saved values remain submitted. Set `SelectedLabel` to provide a
friendly label for that value; it defaults to the value itself.

The evidenced future consumer candidate is Baseplate's account Timezone
override field, `interface_timezone`, in
[`internal/baseplate/account/pages/interface.templ`](https://github.com/Design-Machines-Studio/assembly-baseplate/blob/6accb8a27c746d1cbc3e19a734d8303b8e9191b0/internal/baseplate/account/pages/interface.templ)
and [`internal/timezone/options.go`](https://github.com/Design-Machines-Studio/assembly-baseplate/blob/6accb8a27c746d1cbc3e19a734d8303b8e9191b0/internal/timezone/options.go).
An empty value means “use organization default”; saved zones outside the
curated list remain application-owned and can use `SelectedLabel` for display.
This issue does not change Baseplate.
Later adoption requires a merged and published Go module version containing
the additive exported API, plus compatible published Live Wires assets and
verified checksums. The new exported API requires a minor version under this
repository's release policy.

### Layout

```go
// HTML skeleton
@layout.Base(layout.BaseProps{
    Title:   "My App",
    CSSPath: "/css/livewires.css",
}) {
    @layout.Section(layout.SectionProps{Heading: "Welcome"}) {
        <p>Hello, world!</p>
    }
}
```

### Utilities

```go
lw.ClassNames("button", "button--accent", "")  // "button button--accent"
lw.VariantClass("badge", "green")               // "badge--green"
lw.SchemeClass("dark")                           // "scheme-dark"
lw.AspectRatioClass("4-3")                       // "aspect-4-3"
lw.If(true, "active")                            // "active" (or "" when false)
lw.Initials("Jane Doe")                          // "JD"
lw.SingleInitial("Jane Doe")                     // "J"
lw.CapitalizeFirst("alert")                      // "Alert"
lw.DateShort("2024-03-15")                        // "Mar 15, 2024"
```

## Packages

| Package | Purpose |
| ----------- | ----------- |
| `component` | UI components: button, badge, avatar, card, dialog, tabs, toast, etc. |
| `form` | Form elements: field, select, combobox, checkbox, radio, switch, textarea, search, filter, date range, file upload |
| `layout` | HTML base skeleton, page section |
| root | Class name utilities, text/date helpers |

## Attrs vs InputAttrs

`Attrs` renders on the wrapper/label of wrapped form controls (Switch, Checkbox, FileUpload) and on the Combobox wrapper — use it for wrapper hooks like `data-testid`. `InputAttrs` renders on the nested `<input>` for Switch, Checkbox, and FileUpload, and on the native `<select>` for Combobox, **after** component-generated attributes. Use it for developer-controlled hooks such as `aria-label` and Datastar directives (`data-bind`, `data-on:*`, `data-indicator:*`, `data-attr:*`).

On a collision the component's own attribute wins (HTML first-occurrence), so never pass an attribute the component already emits — this includes `type`, `role`, `name`, `id`, `value`, `checked`, `disabled`, `required`, `form`, `accept`, `multiple`, `aria-invalid`, and `aria-describedby`. Treat that list as illustrative, not exhaustive: as a rule, never duplicate through `InputAttrs` any attribute the component renders itself, or the caller's copy is silently dropped. Neither map may come from untrusted input.

Non-wrapped components like `Button` spread their single `Attrs` onto their own element, so `Attrs` is already their element-level seam — there is no separate `InputAttrs` for them.

## Security

Both the `Attrs` field (on every `Props` struct, spread onto the wrapper or own element) and the `InputAttrs` field (on Switch/Checkbox/FileUpload's nested `<input>` and Combobox's native `<select>`) are inserted directly into HTML output. templ escapes attribute values but does not restrict attribute names. **Never populate `Attrs` or `InputAttrs` from untrusted user input** — they are for developer-controlled attributes like `data-*`, `aria-*`, and Datastar directives.

## Design Pattern

Every component uses a Props struct with `Class` and `Attrs` fields for extensibility:

```go
@component.ButtonComponent(component.ButtonProps{
    Variant: "accent",
    Class:   "my-custom-class",
    Attrs:   templ.Attributes{"aria-label": "Save changes"},
})
```

## Development

Requires Docker.

```sh
make generate   # Run templ generate
make test       # Run tests
make lint       # Run go vet
make all        # generate + test + lint
```

## License

Unlicense

Required Field, Select, and Textarea labels display muted `Required` text from
the existing `Required` prop. Native required semantics remain authoritative;
the decorative indicator is hidden from assistive technology. Custom labels and
group legends can compose `form.RequiredIndicator(required)`. Markdown adapters
can use the static `form.RequiredIndicatorMarkup(required)` with escaped labels.
Consumers must provide readable contrast for Live Wires `text-muted`.

## Browser gallery

Serve the generated components locally with `examples/gallery/run.sh start`.
See the [gallery runbook](examples/gallery/README.md) for pinned assets, Playwright,
source receipts and isolated rebuild/stop commands.

For rich proposal comments and Chat replies, use the additive [shared comment article API](docs/comments.md). Existing `Comment` and `CommentComponent` behavior stays compatible.

The card footer region and summary cards (Live Wires v2.26.0) compose with caller-supplied layout classes, icons and short date labels. Heading and description line limits stay in CSS. See [footer and summary card composition](docs/card-summary.md).

Hardware loaders support Ratchet and Torque in standalone status/progress, busy buttons, LoadingRegion and SaveStatus. Driver Line stays standalone, centered in LoaderPanel or LoaderOverlay. See [loader composition and props](docs/loaders.md).

Inline editor suggestion surfaces: see [Mention suggestions](docs/mention-suggestions.md) for the render API and the consumer-owned editor contract (Live Wires v2.27.0).

Automatic metric sizing is prepared but unreleased: the shared default adapts to its cell width and formatted-value length, with explicit size overrides preserved. See [the widget contract](docs/widgets.md).
