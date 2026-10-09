# Rocket, PD Rockets, and Starbase: opportunities for Live Wires

Research date: 2026-10-09. This is an assessment, not an implementation plan approved for execution.

The strongest opportunity is an optional browser behavior layer around server-rendered Live Wires markup. Keep Go/Templ responsible for HTML, Live Wires responsible for CSS, and consuming applications responsible for persistence, permissions, and command confirmation.

## Evidence and limits

- Reviewed the supplied 82-minute transcript, with particular attention to 15–24 minutes (morphing and pending actions), 27–34 minutes (engineering workflow), and 56–70 minutes (component extraction and documentation).
- Read the current [Rocket reference](https://data-star.dev/reference/rocket).
- Inspected [PD Rockets at 4f411172](https://github.com/derekr/pd-rockets/tree/4f4111722763d2aa8f5851c397af026fcbe97e4f), including contracts, interaction implementations, and browser test source.
- Inspected [Starbase at c7ed1e76](https://github.com/zweiundeins/starbase/tree/c7ed1e76f8e7c867d32351d7045d6b66bf760b34), its live catalog, component documentation, selected implementations, and downstream patch collections.
- Compared against this kit's Kanban, checklist, dropdown, dialog, tabs, select, loader, and toast implementations; also inspected the sibling Live Wires JavaScript entry and popup-dialog implementation.
- No upstream test suites were executed. T3 browser navigation to Starbase succeeded, but two snapshot attempts failed; this is source/documentation research, not a browser accessibility or compatibility certification. No product code changed.

## Recommended candidates

| Candidate | Existing foundation | Recommendation |
| --- | --- | --- |
| Sortable list | Checklist and composable children | First substantial pilot. Useful for agenda order and ordered tasks. PD emits `{itemId, before}`; the application applies the move. |
| Kanban interactions | Board, column, and linked card markup | Add after the list pilot: pointer movement, keyboard movement, focus recovery, and optional edge scrolling. |
| Context/action menu | Dropdown trigger, menu, and item markup | Reuse behavior ideas for focus return, submenus, positioning, and handling replaced content. Keep a visible action button. |
| Inline editing | Heading and field primitives | Add explicit keyboard/button entry, commit/cancel events, validation, and draft preservation. Do not rely on double-click discovery. |
| Relative time | Comments, activities, timelines, date helpers | Small early win. Render an absolute date fallback; progressively enhance with localized relative text and a shared timer. |
| Copy button | Button primitive | Small early win for links and identifiers. Include a truthful failure/manual-copy path. |
| Busy state | Loader, progress bar, button | Extend behavior rather than duplicate visuals: scoped requests, delayed appearance, minimum display time, and announcements. |
| Toast region | Individual toast markup and roles | Add stacking, pause on focus/hover, and local dismissal that survives server updates. |
| Searchable/remote select | Native select with label, hint, error, required | Valuable for larger member/tag lists, but a separate enhanced control. Preserve the native select for ordinary forms. |
| Lazy tree | No general tree primitive found in this kit | Useful when a real hierarchical content use case arrives. Prefer this before draggable trees. |
| Virtual list / data table | Table and pagination primitives | Defer until measured volume warrants it. Requires server windowing, focus/selection rules, and an accessible alternative. |
| Bento dashboard | Cards and layout tools | Defer until there is a demonstrated need for user-arranged dashboards. Geometry, responsive layout, persistence, and keyboard resizing add substantial scope. |

Sources: [PD contracts](https://github.com/derekr/pd-rockets/tree/4f4111722763d2aa8f5851c397af026fcbe97e4f/contracts), [Starbase catalog](https://starbase.zweiundeins.gmbh/), and the individual component README files under [Starbase components](https://github.com/zweiundeins/starbase/tree/c7ed1e76f8e7c867d32351d7045d6b66bf760b34/components).

Count-up effects, typewriter effects, starfields, and general-purpose code editors are lower priorities for this kit. Existing buttons, cards, checkboxes, native dialogs, and ordinary selects do not need replacement merely because Rocket can implement them.

## Architecture to borrow

PD's sortable components use light DOM and setup/cleanup behavior without rendering a second copy of the server's markup. This is a particularly good fit for Templ. Rocket's default is open shadow DOM; light DOM must be an intentional choice. Starbase's styled controls often use shadow DOM and bundled CSS, so they are behavior references rather than direct replacements for this CSS-free library.

Suggested ownership:

1. **livewires-templ:** semantic HTML, stable IDs, props, attribute seams, optional host wrappers.
2. **Live Wires browser assets:** explicitly loaded behavior modules; styles remain in the CSS project. The sibling project already has custom elements, so this extends an existing approach.
3. **Application:** event-to-request bindings, authorization, persistence, request ordering, retries, and confirmed state.

Use semantic events such as “move this item before that item.” Do not embed application endpoint paths, transport policy, or member permissions in the reusable component.

Rocket is beta. Its combined bundle includes Datastar: a consumer should load one compatible runtime, not add another Datastar instance. Start with explicit module imports and pinned versions. An autoloader becomes useful when several optional components exist, especially when server updates introduce new tags. Use a trusted tag-to-module map and preserve useful fallback content if a module fails.

## Concrete integration work in our kit

- `KanbanBoardComponent` currently renders a `div`, not a Rocket host. A wrapper or adapter is needed.
- `KanbanColumnComponent` exposes attributes on the outer column but not its inner card container. An inner-container attribute seam may be necessary; existing `Attrs` alone is not a complete integration API.
- `KanbanCardComponent` renders the entire card as an anchor. A drag handle, action button, or inline input needs a non-link card container and an explicit detail link, rather than interactive descendants nested in an anchor.
- Existing components contain some layout classes despite the current contract prohibiting them. New work should follow the contract and coordinate reference markup/CSS; avoid copying those older patterns automatically.
- Native `Select` already supplies useful form semantics. Starbase's inspected select documentation explicitly lists gaps around required/validity, disabled fieldsets, external labels, and the `form` attribute. Its source uses form-data/reset listeners. Do not treat it as a drop-in native replacement.

## Learnings from the transcript

**Immediate feedback and confirmed state are different responsibilities.** A component can move a preview or show a pending marker immediately. An application still needs to know which command a server update confirms. Do not clear all pending work merely because any update arrives. For repeated moves, handle stale responses, rejection, and retry without duplicating commands.

**Shared rendering has a permission boundary.** Rendering once and distributing the same HTML is useful only for viewers authorized to receive the same content. CSS can express local selection and presentation, but cannot enforce access to private data.

**Choose update regions by measurement.** Start with understandable server-rendered regions. Window large lists when DOM size, rendering, or transfer costs justify it. The transcript's fast demos do not establish performance on our servers or users' networks.

**Make the component manual an interaction testbed.** Include keyboard-only use, slow and rejected requests, updates while dragging/editing, two simultaneous instances, cleanup after removal, and focus after replacement. Show browser work separately from server/network time. The PD benchmark explicitly excludes networking and painting; do not mistake its handler timings for end-to-end latency.

**Document contracts as well as props.** Each interactive component needs its required markup, emitted events, state ownership, keyboard behavior, and lifecycle guarantees. Rocket manifests and Starbase's source-linked examples show how to make these inspectable. A simulated SSE backend is useful for portable examples, supplemented by a real Go/Templ integration fixture.

## Starbase's patches are a valuable research source

Starbase carries its own [Rocket patches](https://github.com/zweiundeins/starbase/blob/c7ed1e76f8e7c867d32351d7045d6b66bf760b34/patches/rocket/README.md) and [PD patches](https://github.com/zweiundeins/starbase/tree/c7ed1e76f8e7c867d32351d7045d6b66bf760b34/patches/pd-rockets). Inspect these before selecting a runtime/component revision; a Starbase demonstration is not evidence that unmodified upstream behaves identically.

Examples include reduced-motion handling for movement animation, measuring all items before animating, keyboard entry to inline edit, context-menu focus and Tab dismissal, and menu cleanup when a server update removes its content. The inspected PD `core/flip.ts` does not contain the reduced-motion guard present in Starbase's patch. These are concrete adoption checks, not hypothetical concerns.

Retain the appropriate notices when vendoring: PD uses Beer-Ware; its Datastar runtime has a separate MIT notice; Starbase identifies exceptions for vendored components. Record the exact source revision and selected patches.

## Suggested first implementation prompt

> Build an isolated sortable-list proof of concept using Live Wires styling, Go/Templ-rendered items, and an optional PD Rockets/Rocket behavior module. Preserve the kit's CSS-free contract and existing APIs. Pin the runtime and component revisions, inspect Starbase's relevant downstream fixes, and retain upstream notices. Keep persistence and endpoint bindings in the example application. Provide pointer and keyboard reordering plus visible move controls, reduced-motion behavior, focus recovery, pending/error/retry feedback, and useful HTML before JavaScript loads. Test a delayed response, rejection, an intervening server update, two list instances, and removal/reinsertion. Document the markup/event contract and measured browser/network behavior. Do not tag or publish a release.


## Planning follow-up: current owner decision

Authenticated planning on 2026-10-09 found Baseplate #1025's explicit decision to retain Datastar Pro and not adopt Rocket, plus existing #1156/#1157/#1151 application work. The original suggestions above are research candidates, not approval to replace the application runtime. The [execution plan](interaction-kit-plan/README.md) preserves that boundary and starts with a kit-local behavior comparison.
