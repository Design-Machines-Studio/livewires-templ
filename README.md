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

// Card (scheme controls color: "subtle", "dark", "accent", or "")
@component.Card("subtle") {
    <p>Card content</p>
}

// Clickable card link
@component.CardLink("/details", "")

// Stat card
@component.StatCardSimple("Revenue", "$12,450", "+12%", "")

// Info card
@component.InfoCard("Board Meeting", "Upcoming")

// Activity card
@component.ActivityCard("Proposal A", "/proposals/1")

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

// Toast
@component.Toast("Saved!", "success")

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

// Checkbox
@form.CheckboxSimple("agree", "I agree to the terms", false)

// Toggle switch
@form.Switch("notifications", "Enable notifications", true)

// Compact toggle switch
@form.SwitchSmall("compact-mode", "Compact mode", false)

// Textarea
@form.TextareaSimple("Message", "message", "", "Write...", "")

// Filter dropdown
@form.Filter(form.FilterProps{Title: "Status", Name: "status", Options: options})

// Date range
@form.DateRange(form.DateRangeProps{Legend: "Date range", StartName: "from", EndName: "to"})
```

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
lw.Initials("Jane Doe")                          // "JD"
lw.SingleInitial("Jane Doe")                     // "J"
lw.DateShort("2024-03-15")                        // "Mar 15, 2024"
```

## Packages

| Package | Purpose |
| ----------- | ----------- |
| `component` | UI components: button, badge, avatar, card, dialog, tabs, toast, etc. |
| `form` | Form elements: field, select, checkbox, radio, switch, textarea, search, filter, date range |
| `layout` | HTML base skeleton, page section |
| root | Class name utilities, text/date helpers |

## Security

The `Attrs templ.Attributes` field on every Props struct is spread directly into HTML output. templ escapes attribute values but does not restrict attribute names. **Never populate `Attrs` from untrusted user input** — it is designed for developer-controlled attributes like `data-*`, `aria-*`, and Datastar directives.

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
