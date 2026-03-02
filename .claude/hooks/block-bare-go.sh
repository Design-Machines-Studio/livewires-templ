#!/bin/bash
# block-bare-go.sh — Prevent Go/Templ commands from running outside Docker
#
# This enforces the project's critical rule: ALL Go toolchain commands
# must run inside Docker via `docker compose exec dev`.
#
# Exit code 2 = block the command and feed error to Claude

COMMAND=$(jq -r '.tool_input.command')

# Match bare go/templ commands that aren't wrapped in docker compose
if printf '%s\n' "$COMMAND" | grep -qE '(^|\&\&|\|\||;)\s*(go |templ )' && \
   ! printf '%s\n' "$COMMAND" | grep -q 'docker compose'; then
  echo "BLOCKED: Go/Templ commands must run inside Docker." >&2
  echo "" >&2
  echo "Use: docker compose exec dev <command>" >&2
  echo "Example: docker compose exec dev go test ./..." >&2
  echo "Example: docker compose exec dev templ generate" >&2
  exit 2
fi

exit 0
