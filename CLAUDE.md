# Live Wires Templ — Component Library

Shared Go/Templ component library that emits HTML with Live Wires CSS class names.
**This library has NO CSS. Zero. None.** CSS is a peer dependency managed by each consuming project.

## Quick Reference

```
import "github.com/Design-Machines-Studio/livewires-templ/component"
import "github.com/Design-Machines-Studio/livewires-templ/form"
import "github.com/Design-Machines-Studio/livewires-templ/layout"
import lw "github.com/Design-Machines-Studio/livewires-templ"
```

## Critical Rules

1. **Docker-only Go/Templ** — All Go and Templ commands MUST run inside Docker:
   ```
   docker compose run --rm dev templ generate
   docker compose run --rm dev go test ./...
   docker compose run --rm dev go vet ./...
   ```

2. **Commit generated files** — Always commit `*_templ.go` files. Consumers should only need `go get`, not the templ CLI.

3. **Props struct pattern** — Every component MUST have:
   - A `Props` struct with `Class string` and `Attrs templ.Attributes` fields
   - A full `*Component(props Props)` templ function
   - One or more convenience wrapper functions

4. **No CSS** — Only emit class names that Live Wires CSS defines. Never bundle CSS.

5. **Match reference HTML** — Output must match `/livewires/public/reference/` HTML examples.

6. **Security: Attrs must not come from untrusted input** — both `Attrs` (all `Props`) and `InputAttrs` (on wrapper-based form controls `Switch`/`Checkbox`/`FileUpload`, for the nested `<input>`) are spread directly into HTML; templ escapes values but does NOT restrict attribute names; never populate either from user input — developer-controlled attributes (`data-*`, `aria-*`, Datastar directives) only.

7. **No layout primitives in markup** — Never bake `stack`, `stack-compact`, `cluster`, or other layout classes into a component. Spacing is the CSS layer's decision. Emit semantic hooks (`.hint`, `.error`) and display utilities (`.block`) only.

8. **Never tag a release without being asked** — See Releasing below. Ask; do not assume.

9. **Two-map attribute seam on wrapped form controls** — wrapper-based form controls (`Switch`, `Checkbox`, `FileUpload`) expose both `Attrs` (spread on the outer `<label>`) and `InputAttrs templ.Attributes` (spread on the nested `<input>`, after component-generated attributes); this is the canonical two-map seam for future wrapped controls; on collision the component's attribute wins (first-occurrence); never duplicate through `InputAttrs` any attribute the component already emits (illustrative, not exhaustive: `type`, `role`, `name`, `id`, `value`, `checked`, `disabled`, `accept`, `multiple`, `aria-invalid`, `aria-describedby`), or the caller's copy is silently dropped.

## Releasing

**Do not tag or publish unless the user explicitly asks for a release.** Finish the work, push the commits, and say what version you would cut and why. Wait.

**A tag is permanent.** The moment any version is fetched, `proxy.golang.org` and `sum.golang.org` cache it forever. Deleting the tag does NOT recall it — the proxy keeps serving that version as the highest, so a "replacement" tagged below it is never resolved. There is no undo. Get the number right the first time.

**Choose the number from the Go API surface, not the rendered output.** This is a component library: its HTML changes on nearly every commit, so "output changed" is not a signal.

| Change | Bump |
| ------ | ---- |
| New exported prop, type, or function (e.g. adding `Hint string`) | minor |
| Renamed, removed, or retyped exported field or function | minor (pre-1.0) |
| Markup, class names, ids, ARIA wiring, generated output | **patch** |
| Bug fix, doc, test, tooling | patch |

Adding `Hint` to a Props struct is a minor. Swapping `stack` for `block` in the emitted markup is a patch, no matter how visible it looks in a browser.

**Batch the release.** One tag per unit of work the user asked for — not one per commit. Several related commits land under a single version.

## Directory Structure

```
livewires-templ/
├── component/          # UI components (badge, button, card, etc.)
├── form/               # Form components (field, select, checkbox, etc.)
├── layout/             # Structural layout (base, section)
├── classnames.go       # CSS class builder utilities (root package)
├── helpers.go          # Shared utility functions (root package)
├── go.mod
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── CLAUDE.md
```

## Component Design Pattern

Every component follows this pattern:

```go
// 1. Props struct — always has Class and Attrs
type BadgeProps struct {
    Label   string
    Variant string
    Size    string
    Class   string           // Consumer can add extra classes
    Attrs   templ.Attributes // Enables data-*, aria-*, etc.
}

// 2. Full component with Props
templ BadgeComponent(props BadgeProps) {
    <span class={ lw.ClassNames("badge", lw.VariantClass("badge", props.Variant), props.Class) }
          { props.Attrs... }
    >{ lw.Title(props.Label) }</span>
}

// 3. Convenience wrappers
templ Badge(label string, variant string) {
    @BadgeComponent(BadgeProps{Label: label, Variant: variant})
}
```

## CSS Class Naming Conventions

- **Component modifiers** use double-dash: `button--accent`, `badge--small`, `toast--success`
- **Layout modifiers** use single-dash: `stack-compact`, `box-tight`, `sidebar-reverse`
- **Color schemes**: `scheme-subtle`, `scheme-dark`, `scheme-accent`
- **Text sizes**: `text-sm`, `text-lg`, `text-6xl`

Use root package helpers:
- `lw.ClassNames(...)` — join non-empty class strings
- `lw.VariantClass("base", "variant")` → `"base--variant"`
- `lw.ModifierClass("base", "modifier")` → `"base-modifier"`
- `lw.SchemeClass("subtle")` → `"scheme-subtle"`
- `lw.If(cond, "val")` → `"val"` when true, `""` when false

## How to Add a New Component

1. Create `component/my_thing.templ` (or `form/` for form elements)
2. Define `MyThingProps` struct with `Class` and `Attrs` fields
3. Create `MyThingComponent(props)` templ function
4. Add convenience wrappers
5. Add render test in `component/my_thing_test.go`
6. Run `make generate && make test`
7. Commit both `.templ` and `*_templ.go` files

## Build Commands

```
make generate   # Run templ generate inside Docker
make test       # Run go test
make lint       # Run go vet
make check      # Verify generated files are fresh
make all        # generate + test + lint
```

## Packages

| Package | Import | Purpose |
|---------|--------|---------|
| root | `lw "github.com/.../livewires-templ"` | ClassNames, helpers, date formatting |
| component | `".../livewires-templ/component"` | UI components: button, badge, card, etc. |
| form | `".../livewires-templ/form"` | Form elements: field, select, checkbox, etc. |
| layout | `".../livewires-templ/layout"` | Base HTML skeleton, section wrapper |

## Dependencies

`github.com/a-h/templ` is the only runtime dependency. `golang.org/x/net` is a direct module dependency used exclusively by test infrastructure (`internal/testutil`, imported only from `_test.go` files), so consumer builds never link it. Keep dependencies otherwise minimal.
