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

// Avatar
@component.Avatar("Jane Doe", "")

// Stat card
@component.StatCardSimple("Revenue", "$12,450", "+12%", "")

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
@form.Switch("notifications", "notif", true)

// Textarea
@form.TextareaSimple("Message", "message", "", "Write...", "")
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
lw.DateShort("2024-03-15")                        // "Mar 15, 2024"
```

## Packages

| Package | Purpose |
| ----------- | ----------- |
| `component` | UI components: button, badge, avatar, card, dialog, tabs, toast, etc. |
| `form` | Form elements: field, select, checkbox, radio, switch, textarea, search, filter, date range |
| `layout` | HTML base skeleton, page section |
| root | Class name utilities, text/date helpers |

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
