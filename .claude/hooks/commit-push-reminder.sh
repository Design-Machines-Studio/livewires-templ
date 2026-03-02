#!/bin/bash
# commit-push-reminder.sh — Nudge toward frequent commits and pushes
#
# PostToolUse hook on Edit|Write. After file edits, checks:
# 1. Uncommitted file count → suggest commit at 3+ files, insist at 5+
# 2. Unpushed commit count → suggest push at 3+ commits

INPUT=$(cat)

HEAD=$(git -C "${CLAUDE_PROJECT_DIR}" rev-parse --short HEAD 2>/dev/null)

# Count uncommitted changes (unstaged + staged, deduplicated)
CHANGED=$(git -C "${CLAUDE_PROJECT_DIR}" diff --name-only 2>/dev/null)
STAGED=$(git -C "${CLAUDE_PROJECT_DIR}" diff --cached --name-only 2>/dev/null)
TOTAL=$(printf "%s\n%s" "$CHANGED" "$STAGED" | sort -u | grep -c -v '^$')

# Count unpushed commits (0 if no upstream tracking)
UNPUSHED=$(git -C "${CLAUDE_PROJECT_DIR}" log @{u}..HEAD --oneline 2>/dev/null | wc -l | tr -d ' ')
[ -z "$UNPUSHED" ] && UNPUSHED=0

MSG=""

# Strong nudge at 3+ files
if [ "$TOTAL" -ge 3 ]; then
  MSG="You have $TOTAL uncommitted file changes. Stop and commit now with a focused message. Keep commits to 1-4 files."

# Gentle nudge at 2+ files
elif [ "$TOTAL" -ge 2 ]; then
  MARKER="/tmp/livewires-templ-commit-nudge-${HEAD}"
  if [ ! -f "$MARKER" ]; then
    touch "$MARKER"
    MSG="$TOTAL files changed since last commit. Commit before making more changes."
  fi
fi

# Push nudge at 2+ unpushed commits
if [ "$UNPUSHED" -ge 2 ]; then
  PUSH_MARKER="/tmp/livewires-templ-push-nudge-${UNPUSHED}"
  if [ ! -f "$PUSH_MARKER" ]; then
    touch "$PUSH_MARKER"
    PUSH_MSG="You have $UNPUSHED unpushed commits — push to remote."
    if [ -n "$MSG" ]; then
      MSG="$MSG $PUSH_MSG"
    else
      MSG="$PUSH_MSG"
    fi
  fi
fi

if [ -n "$MSG" ]; then
  MSG_JSON=$(echo "$MSG" | jq -Rs '.')
  echo "{\"systemMessage\": $MSG_JSON}"
fi

exit 0
