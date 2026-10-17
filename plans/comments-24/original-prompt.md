Use $project-manager:assembly-coordinator to refresh dependencies and ownership, then invoke the installed /pipeline workflow.

Implement only:
https://github.com/Design-Machines-Studio/livewires-templ/issues/24

Repository: Design-Machines-Studio/livewires-templ
Main verified during prompt preparation:
5cf257e368dffe612b62f59301a6f4df486ebf3b

executorRole: builder-deep
executorCapabilities: [read-repository, write-repository, tool-use]
executorEffort: high
noMergeOnCompletion: true

Start
Read repository instructions and applicable Live Wires and accessibility guidance. Confirm the remote, fetch main, inspect worktrees and check the configured agent-board inbox. Record the actual execution base.

Use a clean task-owned worktree and branch:
 /home/ned/assembly/livewires-templ-worktrees/comments-24
 feat/comments-24

Preserve occupied checkouts, unrelated changes and serving instances. Page-header PR #23 has merged; refresh current ownership rather than assuming its old collision remains active.

Verified producer
Live Wires PR #36 merged:
https://github.com/Design-Machines-Studio/livewires/pull/36
Merge commit: 35c81f71e863db5479d56b513fd35217bb56efc7

Read its exact contract:
- src/css/6_components/comments.css
- public/reference/components/comments.html
- public/reference/components/comments-tests.html
- public/reference/components/comments-tests.js
- public/manual/components/comments.html

Verify that contract remains on trusted main. A merged reference is not proof of a published CSS artifact.

Goal
Extend the shared comment API additively so Governance proposal comments and Chat replies can use the same presentation.

Inspect component/comment.templ, its callers and tests, the maintained gallery, and current consumer sources:
- Governance: pages/workspace/proposals/discussion_comment.templ
- Chat: pages/workspace/conversations/thread.templ, PostArticle

Acceptance
- Match the merged Live Wires HTML/CSS contract exactly.
- Preserve existing Comment and CommentComponent callers, escaped plain Body text and existing date semantics.
- Provide an additive rich composition API using templ.Component/children. No raw HTML string escape hatch.
- Support stable IDs, author/avatar/profile composition, semantic datetime with caller-formatted visible text, optional timestamp permalink, edited metadata and action/footer slots.
- Follow repository Props conventions, including Class and Attrs. Add inner seams only where supported by consumer evidence.
- Keep author-only, unlinked-avatar and no-avatar variants accessible. Avoid redundant avatar announcements; linked avatars need accessible names.
- Keep action controls outside author/permalink links. Preserve valid nesting, focus and sensible reading order.
- Emit shared CSS hooks without baked layout primitives. This library contains no CSS.
- Do not add automatic live regions, domain flags, Markdown parsing, permission checks, routes, threading, editing behavior or new runtime dependencies.
- Keep Governance positions/decisions and Chat ordering/moderation in their Fixtures.

Gallery and tests
Add maintained examples for ordinary proposal comments and Chat replies using the same component, plus linked/unlinked authors, optional avatars, semantic time/permalinks, edited metadata, actions, multiple unique anchors and rich content.

Cover long names, long links, lists, code and images at desktop and 320/375px widths, zoom, light/dark/inherited schemes, keyboard navigation and no-JavaScript rendering.

Add meaningful render tests for escaping, old API compatibility, IDs, optional slots, semantic time and valid link/action composition.

Run Go/Templ generation, tests and vet inside Docker under repository instructions. Commit generated *_templ.go files and verify freshness.

Browser evidence
Discover the maintained gallery serving checkout, URL and build/restart commands through repository runbooks and the installed browser-target discovery workflow. These bindings were not verified during prompt preparation; resolve them before capture.

Use a task-owned preview when supported. Do not take over another owner's checkout or port. Record the exact checkout, branch/head, CSS revision and commands. Follow installed browser-review requirements and distinguish DOM checks from actual screen-reader testing.

Delivery
Use Pipeline's proportional integrated review and publication procedure. Preserve explicit plan/material-scope approval, applicable independent review, designer acceptance and affected-case rechecks.

Commit/push the verified candidate. Create/register the PR through the installed reviewed-PR procedure after its pre-PR gate. Verify the remote head. Keep it draft until final-head CI, feedback settlement and applicable UI acceptance pass. Reuse valid unchanged coverage while waiting.

Leave the PR unmerged. No tags, releases, issue closure, consumer migrations, deployment or shared reseeding.

Document the exported API, exact compatible CSS contract and proposed additive minor release; do not publish it.

Return the PR URL, gallery URL/review commands, API summary, verification results, remaining gaps and one next action. Include the actual model/cost report once.
