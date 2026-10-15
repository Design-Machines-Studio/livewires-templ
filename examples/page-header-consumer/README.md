# External PageHeader consumer

This nested Go module compiles outside the library's module boundary and imports
`component.PageHeaderComponent` and `component.PageHeaderProps`, which proves a
Fixture can render the public API without importing an `internal` package.

The local `replace` directive is only for this repository's producer check. A
consumer should pin the published `livewires-templ` version after this change is
released. No Live Wires CSS or JavaScript is included here; the host supplies
its compatible Live Wires assets and caller-owned layout classes.
