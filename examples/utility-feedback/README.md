# Interaction utility example

`Page()` renders all four component hooks and points at the exact merged Live
Wires `public/dist` files listed in `docs/interaction-utilities.md`. It is a
small Go/Templ render fixture, not an application or a new shared server.

Run its render check inside the repository's Docker development service:

```sh
docker compose run --rm dev go test ./examples/utility-feedback
```

The repository has no maintained serving command or bound browser target for
this utility example. The render test does not establish that external modules
load, and it does not replace the producer-reference browser cases. Resolve a
maintained consumer/serving checkout and its exact source/domain binding before
claiming browser acceptance; do not start or repoint an unregistered local
server for this fixture.
