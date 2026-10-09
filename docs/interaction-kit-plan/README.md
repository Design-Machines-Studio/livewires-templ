# Live Wires interaction plan

## Widget follow-up — 2026-10-09

**Current recommended next chunk: [W00 global metric typography](https://github.com/Design-Machines-Studio/assembly-baseplate/issues/1160).** W01 then establishes the shared widget reference. The user requested unifying existing Baseplate/Fixture widget designs after the original interaction plan. [The source inventory](widgets.md) establishes real repeated consumers; this track does not depend on Rocket, sortable lists or a new dashboard engine. The original sortable issue remains independently Next, but do not start simultaneous Live Wires styling/review work without checking shared paths and serving ownership. No parallel lane is recommended initially.

| ID | Native issue | Depends on | Fresh-session prompt |
| --- | --- | --- | --- |
| W00 | [Global font-metric token and utility](https://github.com/Design-Machines-Studio/assembly-baseplate/issues/1160) | none | [W00](W00.md) |
| W01 | [Widgets: establish shared frame, metric and message references](https://github.com/Design-Machines-Studio/livewires/issues/23) | W00 | [W01](W01.md) |
| W02 | [Widgets: add shared frame, metric and message Templ components](https://github.com/Design-Machines-Studio/livewires-templ/issues/13) | W01 | [W02](W02.md) |
| W03 | [Widgets: adopt shared Live Wires presentation without changing the editor](https://github.com/Design-Machines-Studio/assembly-baseplate/issues/1159) | W00, W02 | [W03](W03.md) |
| W04 | [Widgets: replace Governance body duplication with shared components](https://github.com/Design-Machines-Studio/assembly-governance/issues/114) | W02, W03, W07 | [W04](W04.md) |
| W05 | [Widgets: unify Atlas frames while preserving chart and data semantics](https://github.com/Design-Machines-Studio/assembly-atlas/issues/60) | W02, W03, W07 | [W05](W05.md) |
| W06 | [Widgets: teach the shared widget contract in the Fixture template](https://github.com/Design-Machines-Studio/assembly-fixture-jig/issues/42) | W02, W03 | [W06](W06.md) |
| W07 | [Per-widget size, weight and variable-font width](https://github.com/Design-Machines-Studio/assembly-baseplate/issues/1161) | W00, W03 | [W07](W07.md) |

W00 is Next; W01–W07 are Blocked on named producer contracts/publication. W03 additionally depends on W00; W04/W05 additionally depend on W07 for typography acceptance. All are P3 / Design in Project 1. Production adoption is a separate stage from component publication. The original dashboard discovery #22 now evaluates remaining gaps in the existing editor after W03; its former assumption of no demonstrated dashboard consumer is superseded. No old application/Fixture issue was closed.

The sections below preserve the original 16-task planning snapshot; use current native Issue/Project state and the widget follow-up for selection.

The original interaction track starts with [Live Wires #11](https://github.com/Design-Machines-Studio/livewires/issues/11): prove a sortable list and settle the optional behavior contract. No parallel implementation lane is recommended before that contract is established. It remains Next within that track; W00 is the current recommendation after the widget and typography requests.

Parent issue: https://github.com/Design-Machines-Studio/livewires-templ/issues/8. Project: https://github.com/orgs/Design-Machines-Studio/projects/1. Native issues and dependency edges are authoritative; these prompts are handoffs, not a second task database. The explicitly requested 16-task set and its parent are projected in Project 1 under Design / P3; this does not import the repositories' unrelated backlog or create a beta/release gate.

## Work and prompts

| ID | Native issue | Owner repository | Initial status | Dependencies | Fresh session |
| --- | --- | --- | --- | --- | --- |
| 01 | [Interactions: prove a sortable list and settle the optional behavior contract](https://github.com/Design-Machines-Studio/livewires/issues/11) | livewires | Next | none | [Prompt](01.md) |
| 02 | [Components: add sortable-list markup and a real Go/Templ interaction example](https://github.com/Design-Machines-Studio/livewires-templ/issues/9) | livewires-templ | Blocked | 01 | [Prompt](02.md) |
| 03 | [Utilities: add relative timestamps and truthful copy feedback](https://github.com/Design-Machines-Studio/livewires/issues/12) | livewires | Blocked | 01 | [Prompt](03.md) |
| 04 | [Feedback: add scoped loading behavior with stable announcements](https://github.com/Design-Machines-Studio/livewires/issues/13) | livewires | Blocked | 01 | [Prompt](04.md) |
| 05 | [Feedback: keep toast dismissal and timing stable through updates](https://github.com/Design-Machines-Studio/livewires/issues/14) | livewires | Blocked | 04 | [Prompt](05.md) |
| 06 | [Components: expose relative-time, copy and feedback behavior hooks](https://github.com/Design-Machines-Studio/livewires-templ/issues/10) | livewires-templ | Blocked | 03, 04, 05 | [Prompt](06.md) |
| 07 | [Menus: add accessible action-menu behavior around server markup](https://github.com/Design-Machines-Studio/livewires/issues/15) | livewires | Blocked | 01 | [Prompt](07.md) |
| 08 | [Editing: add inline-edit behavior with draft and focus preservation](https://github.com/Design-Machines-Studio/livewires/issues/16) | livewires | Blocked | 01, 04 | [Prompt](08.md) |
| 09 | [Kanban: add accessible movement and composable card markup](https://github.com/Design-Machines-Studio/livewires/issues/17) | livewires | Blocked | 01, 07, 08 | [Prompt](09.md) |
| 10 | [Components: expose action menus, inline editing and enhanced Kanban markup](https://github.com/Design-Machines-Studio/livewires-templ/issues/11) | livewires-templ | Blocked | 07, 08, 09, 02 | [Prompt](10.md) |
| 11 | [Forms: add a searchable and remote combobox with native form guarantees](https://github.com/Design-Machines-Studio/livewires/issues/18) | livewires | Blocked | 01, 04, 07 | [Prompt](11.md) |
| 12 | [Forms: add optional combobox wrappers with accessible errors and hints](https://github.com/Design-Machines-Studio/livewires-templ/issues/12) | livewires-templ | Blocked | 11 | [Prompt](12.md) |
| 13 | [Components: load optional modules on demand when measured use warrants it](https://github.com/Design-Machines-Studio/livewires/issues/19) | livewires | Blocked | 01, 03, 04, 07 + admission condition | [Prompt](13.md) |
| 14 | [Discovery: scope a lazy tree only against a real navigation use case](https://github.com/Design-Machines-Studio/livewires/issues/20) | livewires | Blocked | 01 + admission condition | [Prompt](14.md) |
| 15 | [Discovery: measure whether a virtual list or data table is needed](https://github.com/Design-Machines-Studio/livewires/issues/21) | livewires | Blocked | 01 + admission condition | [Prompt](15.md) |
| 16 | [Discovery: assess remaining gaps in the existing dashboard editor](https://github.com/Design-Machines-Studio/livewires/issues/22) | livewires | Blocked | W03 + measured remaining gap | [Prompt](16.md) |

Tasks 14–16 are discovery only. A missing real consumer or measured need should produce a defer decision; it must not produce a speculative component. Task 13 similarly allows a no-change outcome if explicit imports remain simpler. Every implementation prompt invokes Pipeline explicitly. Do not launch dependent prompts until their native blockers are cleared at the stated evidence level.

## Sequence and ownership

- 01 establishes the contract and reference testbed. 02 proves the real Go/Templ consumer.
- 03–05 add utilities/loading/toasts; 06 packages their Go markup.
- 07–09 add menus/editing/boards; 10 packages their Go markup.
- 11–12 add an enhanced select while retaining the native control.
- 13 measures optional loading; 14–16 test admission for larger candidates.

All new issues are unassigned; their future executor claims one branch and issue. Initially safe parallel lanes: **None**. After 01 merges, the coordinator may choose 02 (Templ owner) alongside one independent Live Wires task, after checking build/asset/serving collisions. One shared Live Wires review instance must be serialized. No present lane is assigned to another developer or agent.

## Live coordination evidence: 2026-10-09

| Repository | Exact fetched origin/main | Relevant state |
| --- | --- | --- |
| Design-Machines-Studio/livewires-templ | `227c3d6c063eefcbadf6a027597d4f4c72e14b6d` | No open PRs; unrelated issue #2 owns textarea hint placement. |
| Design-Machines-Studio/livewires | `c67bacb24ba4c90e3873b5a50ac1b7404f524687` | PR #10 and #3 remain owned background work; no existing native issue duplicates this set. |

The originating Templ checkout is detached at `119341903daee1ec516aa26a3cab9e17d80db837`, behind trusted main. Its untracked research document is preserved. This planning branch was created in a clean worktree from the fetched Templ main, not from that stale checkout. Templ PR #7 head `d3d95eaebcb8f3df873307cc450960ed1247eacf` is merged as `227c3d6c063eefcbadf6a027597d4f4c72e14b6d`; required-label changes are already present and must be retained. That is merged-source proof, not new publication or consumer proof.

Live Wires PR #10 head `70e50ceadc3ab4b5fd477346ccd5eb9951c1b87a` is open, non-draft, mergeable, with no review decision; Cloudflare Pages passed and Macroscope was skipped. The canonical checkout is on its feature branch. PR #3 head `237c72c30b069fc72a6fc8ea422c812f9dda4e20` is open, non-draft, conflicting, with no review decision or reported checks. Its package/build/default-entry/token/design-panel paths overlap common integration surfaces. Neither PR was changed, reviewed for approval, merged or registered as work in this thread.

Bounded agent-board inboxes for both kit repositories were empty. Authenticated Project 1 inspection covered all 466 pre-existing items. Existing Baseplate #1156 and #1157 already cover transcript-derived poor-network and DOM-boundary work; #1151 owns Fixture toast persistence. They remain untouched. #1157 is blocked on #1156. Baseplate #1025 is closed as not planned and explicitly preserves Datastar Pro with no planned Rocket adoption. The first task evaluates an optional kit implementation; it does not overturn that owner decision.

## Browser and publication prerequisites

The Live Wires checkout is `/home/ned/assembly/livewires`. Its README/Vite document `npm run dev`, localhost port 3000, and `npm run build`. A complete maintained public-domain/service binding has not been verified. The catalog's DM-003 code is not proof of a running review target. The first executor must resolve source/target/ownership using the installed dm-review discovery contract. Do not steal the PR #10 serving checkout or repoint infrastructure. Templ has no standalone browser server; its consumer example needs a documented, task-owned serving binding.

No browser checks, product builds, upstream test runs or release checks were performed in this planning session. No external application or package publication is authorized. Source dependencies require exact merged trusted-main evidence. A release-dependent consumer must separately verify its published artifact and consumer behavior; a local replace is only local proof.

## Research coverage and deferrals

The [source assessment](../rocket-component-research.md) covers the transcript, PD Rockets and Starbase. The following learnings are acceptance requirements in the tasks rather than separate process projects: semantic events; light-DOM/progressive enhancement where suitable; cleanup and focus; explicit pending versus confirmation; permission-safe server rendering; tests under delay/rejection and incoming updates; reference/source documentation; source pins and downstream patch review; and separate browser/network measurements.

Existing buttons/cards/native dialogs, ordinary form controls and date/calendar components remain in place. Decorative count-up/typewriter/starfield effects and general code editors are not proposed work. Drawer/popover mechanisms may be used inside the menu/select references when needed, not added as a parallel component system. No new architecture, universal retry service, CRDT/editor engine or shared-render permission bypass is planned.

No tracked coordination document was repaired. `tasks/todo.md` has older future entries for already-existing components, but rewriting that unrelated historical list is outside this plan. The new documents add durable research and execution handoffs. Planning branch: `docs/interaction-kit-plan-20261009`; no planning PR is required or requested. The final pushed SHA is recorded in the session handoff and native parent issue.
