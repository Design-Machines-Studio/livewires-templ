# Live Wires Templ — Lessons

## Patterns
- Test files in templ component packages need explicit `import "github.com/a-h/templ"` for the `templ.Component` type
- Root package helpers are exported (uppercase) since they're used across sub-packages
- Assembly components use unexported helpers (lowercase) since they're within a single package
