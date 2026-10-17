Use $project-manager:assembly-coordinator to refresh dependencies and ownership, then invoke the installed /pipeline workflow.

Implement only:
https://github.com/Design-Machines-Studio/livewires-templ/issues/13
“Widgets: add shared frame, metric and message Templ components”

Repository: Design-Machines-Studio/livewires-templ
Verified main and v0.20.0:
bfc09a97db1da706484b3e03f8e2c858ca3adf4b

executorRole: builder-deep
executorCapabilities: [read-repository, write-repository, tool-use]
executorEffort: high
noMergeOnCompletion: true

Start
Read repository instructions and applicable Live Wires/accessibility guidance.
Confirm the remote, fetch main, inspect existing worktrees and check the configured agent-board inbox. Record the actual execution base and check for existing implementation work before creating anything.

Use a clean task-owned worktree and branch:
 /home/ned/assembly/livewires-templ-worktrees/widgets-13
 feat/widgets-13

Preserve other branches, unrelated files and running previews. The canonical checkout may be detached and occupied; do not switch it. Coordinate shared gallery files with any active owner.

Goal
Provide small reusable widget frame, metric and quiet-message components matching the merged Live Wires reference. Extract presentation shared by Baseplate, Governance and Atlas.

Keep dashboard placement, registry, permissions, editor controls, saved settings and domain data in their existing owners.

Producer authority
Live Wires PR #31 merged:
https://github.com/Design-Machines-Studio/livewires/pull/31
Merge commit: d6479198c0d2d1c54682ab983be2926e2e5d8e98

The current gallery already pins compatible Live Wires v2.24.0:
3042ffe56c36b8e51e4e8a98f500b1aa04a43fd3

Read these files at the pinned producer revision:
- docs/widget-references.md
- public/reference/components/widgets.html
- public/manual/components/widgets.html
- src/css/6_components/widgets.css
- src/css/7_utilities/typography.css
- public/reference/components/cards.html

Verify the widget contract remains present in that published revision. Reuse the current compatible gallery asset lock; do not downgrade it or introduce floating assets.

Read the producer document’s consumer inventory and exact source comparisons. Refresh only the relevant Baseplate/Governance/Atlas sources needed to settle an API choice. Historical screenshots and running installations do not establish current source identity.

Inspect existing Templ Card, StatCard, EmptyState, Timeline, StatusIndicator, ProgressBar and heading patterns before designing new APIs.

Component acceptance
- Mirror the exact merged HTML/CSS contract.
- Follow Props conventions, including Class, Attrs and convenience wrappers.
- Provide a widget frame/body with children and optional header/action/footer composition. Support caller-owned instance/heading IDs and appropriate heading levels.
- Preserve optional descriptions, headingless compositions and accessible naming. Do not emit an unnamed section landmark when a non-landmark body is appropriate.
- Expose only evidenced heading/body/footer attribute seams.
- Support the producer’s plain-header variation without route-specific exceptions.
- Icons are caller-supplied templ.Component values. No Baseplate icons, DTOs, Fixture SDK or runtime imports.
- Emit semantic CSS hooks without baked layout primitives. This library contains no CSS.
- Do not add automatic status/live-region roles. Callers decide whether dynamic updates need announcements.
- Preserve existing Card, StatCard and EmptyState APIs and behavior except the explicitly requested value-only font-metric hook.

Metric semantics and typography
- Match standalone dl.widget-metric and grouped dl.widget-metrics > div.widget-metric anatomy. Keep dt/dd relationships valid; avoid nested cards and div-only definition lists.
- Accept caller-formatted strings. Preserve complete currency symbols, localized amounts, signs, units, precision and meaningful distinctions.
- Distinguish a confirmed zero from missing, denied or unavailable data. Status text must not receive numeric-value styling or font-metric.
- Apply font-metric to prominent shared metric values, including the existing StatCard value where the reference requires it. Leave labels, descriptions, ordinary numbers and runtime version strings unchanged.
- Support the producer’s default/short/long value variations.
- Provide value-scoped composition/attribute hooks for the existing --widget-metric-size, --widget-metric-weight and --widget-metric-width contract. Avoid creating a second typography configuration system.
- Preserve inherited defaults when unset. Width means glyph width, separate from widget columns.
- Use the producer’s standard font properties and static-font fallbacks. No transform scaling, truncation, automatic digit counting or shrink-to-fit loop.
- Do not add competing size/weight utilities that defeat hook-controlled values.
- Treat Attrs/style hooks as developer-controlled. Do not turn arbitrary user strings into CSS.

Quiet messages and composition
- Match widget-message, optional decorative icon, content and detail anatomy.
- No extra card/surface, fabricated number or forced live region.
- Demonstrate composition with existing Timeline, StatusIndicator, ProgressBar and semantic table content.
- Preserve captions, column/row scopes and complete data in chart fallback tables.
- Do not introduce chart/map runtimes, provider-specific widget APIs, a video configuration system, flip controls or a dashboard engine.
- Not every producer reference composition needs a new exported component. Prefer existing primitives and caller slots.

Gallery and documentation
Add a maintained widgets example using the actual public components.

Cover:
- Frame with and without description, action, header and footer.
- Long headings/actions and multiple unique heading IDs.
- Standalone/grouped metrics, zero, negative values, percentages, counts, long localized money and unavailable/denied text.
- Default/short/long variants and independent typography overrides on two widgets.
- Quiet empty/unavailable messages with optional icons/details.
- Representative list, timeline, progress and chart-with-table compositions.
- Narrow/wide containers and inherited light/dark/accent schemes.

Document the public API, exact CSS requirement, safe attribute ownership and the mapping from existing consumer markup to shared components.

Baseplate #1160 owns global metric-font configuration; #1161 owns persisted per-widget settings, reset and draft/Save behavior. Do not implement those features here or claim theme save/reload proof from a Templ gallery.

Downstream adoption remains:
- Baseplate #1159, coordinated through #1201
- Governance #114
- Atlas #60
- Fixture Jig #42

Do not modify those repositories or start their migrations.

Verification
Run Go/Templ generation, tests and vet inside Docker. Commit generated *_templ.go files and verify freshness.

Add meaningful render tests for:
- Escaping and developer attribute placement.
- Optional slots and stable unique IDs.
- Heading/accessibility relationships.
- Valid standalone/grouped definition-list markup.
- Numeric zero versus unavailable text.
- Value-only font-metric and typography hooks.
- Legacy component compatibility.

Use the maintained gallery runner:
 examples/gallery/run.sh start
 examples/gallery/run.sh status
 examples/gallery/run.sh verify
 examples/gallery/run.sh rebuild
 examples/gallery/run.sh stop

Run it from the task-owned checkout. It allocates its own loopback port; record the actual URL, source head/fingerprint and build commands. Do not take over the loader preview or another owner’s instance.

Follow current T3 browser rules and the repository’s supported Playwright workflow. Verify desktop, 320/375px containers, reflow/zoom, keyboard focus, no-JavaScript rendering, themes, long-value wrapping and table overflow. Distinguish DOM checks and CSS zoom from actual screen-reader/native browser zoom testing.

Prove variable-font width behavior only with an available supported font and correct declared axis ranges. A computed font-stretch value alone is not proof that glyphs changed. Verify a readable static/missing-font fallback; report any unavailable real-font coverage honestly.

Compare rendered Templ output against the exact producer reference. Reuse valid existing gallery coverage; broaden checks only when changed shared assets or source justify it.

Workflow and delivery
Use Pipeline’s proportional workflow with the existing implementation request as scope authority. Ask only for a genuinely unresolved design decision or material scope change.

Initialize supported source-bound review evidence before dispatch. Keep prompts, exact source/patches, actual reviewer receipts, browser evidence and repairs attributable to their real heads. Do not leave evidence assembly until after review or invent missing provenance.

Use bounded review packets: complete relevant source for each lane, with generated code handled according to the installed review contract. Preserve required review coverage without unrelated repository audits or repeated unchanged-head reviews.

Complete applicable independent review, fixes, affected-case rechecks and formal evidence validation. Commit/push the verified candidate and create/register the PR through the installed reviewed-PR procedure. Verify its remote head.

Keep the PR draft until applicable final-head checks, feedback settlement and designer acceptance pass. Preserve the task-owned preview and evidence while acceptance is pending.

Leave the PR unmerged. No tag, release, deployment, shared reseed or consumer migration is authorized. Determine the proposed minor release from the actual exported API and current tags; do not hardcode a version prematurely.

Return the PR URL, working gallery URL, API summary, producer comparison, verification results, precise remaining gaps and one next action. Include the actual routed model/cost report once, with unmeasured usage identified.
