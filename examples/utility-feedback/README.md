# Interaction utility example

`Page()` renders all four component hooks and points at the exact merged Live
Wires `public/dist` files listed in `docs/interaction-utilities.md`. It is a
small Go/Templ render fixture, not an application or a new shared server.

Run its render check inside the repository's Docker development service:

```sh
docker compose run --rm dev go test ./examples/utility-feedback
```

For maintained local serving and browser review, use the [component gallery](../gallery/README.md). It shares these example templates and loads verified local assets.
