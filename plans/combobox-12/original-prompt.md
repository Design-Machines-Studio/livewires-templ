Use $project-manager:assembly-coordinator to refresh ownership, then invoke the installed /pipeline workflow.

Implement only:
https://github.com/Design-Machines-Studio/livewires-templ/issues/12
“Forms: add optional combobox wrappers with accessible errors and hints”

Repository: Design-Machines-Studio/livewires-templ
Verified main/gallery merge:
284684633bcc5e026155b62fe5a1c21541f49099

Combobox producer, Live Wires PR #32:
5db085d2918375811ee30562bc4229b414f314e4

executorRole: builder-fast
executorCapabilities: [read-repository, write-repository, tool-use]
executorEffort: high
noMergeOnCompletion: true

Start
Confirm the remote, fetch main, verify the merged producer, read repository instructions and applicable skills, and check the agent-board inbox. Record exact execution heads.

Create:
 /home/ned/assembly/livewires-templ-worktrees/combobox-12
Branch: feat/combobox-12

Preserve unrelated worktrees and files. Refresh active ownership before changing shared form or gallery files.

Outcome
Expose the merged single-select combobox contract through additive CSS-free Templ components and demonstrate it in the maintained gallery.

Read Live Wires docs/combobox.md, its reference markup, source and tests at the producer commit. The native select remains the authoritative form control. Do not implement another combobox state machine.

Acceptance
- Follow Props/Class/Attrs and convenience-wrapper conventions, with attribute seams on the actual select.
- Support explicit IDs, label, options, selected value, hint, server error, required, disabled and form association.
- Preserve ordinary Select APIs and required-indicator behavior. Avoid unrelated textarea-hint work.
- Render the native select fallback. Its name/value own submission; the enhanced search input has no name. Prevent duplicate submitted values.
- Preserve empty-value semantics: optional “use organization default” is valid; required placeholders must still fail validation.
- Preserve selected labels and saved values absent from current search results, following the producer contract.
- Keep CSS/JS optional and external. Applications own remote queries, authorization, validation and saving.
- Test escaping, attribute placement, unique instances, selected options, hints/errors and compatibility.
- No multi-select, production search endpoint, Baseplate migration or Kanban.

Gallery
Use examples/gallery/README.md and the existing runner. Add a stable combobox example route with local and simulated remote search, deterministic delay/error/out-of-order controls, and native fallback.

Update the asset lock deliberately to an exact compatible Live Wires revision with verified checksums. Preserve existing examples and test them if shared assets change. Do not copy CSS/JS into the component package or inject temporary browser setup.

Use the documented commands from the task worktree:
 examples/gallery/run.sh start
 examples/gallery/run.sh status
 examples/gallery/run.sh rebuild
 examples/gallery/run.sh verify
 examples/gallery/run.sh stop

Read the actual URL from status; rebuild can change the port. For Codex, resolve the host-selected Kernel launcher and set WORKFLOW_KERNEL as documented. Do not hardcode a cache version into repository code.

Local Playwright is authorized. Follow the gallery’s pinned dependency/setup instructions and GALLERY_URL/GALLERY_EVIDENCE conventions. T3 reachability and a hosted domain are not prerequisites. Preserve other instances and clean up only owned resources.

Verification
Run Docker-only generation, tests and vet; commit generated Go. Run gallery verification and source-bound desktop/mobile browser checks.

Prove real FormData/submission, required validation, disabled fieldset behavior, reset, form association, keyboard navigation/Escape, stale results, failed searches, retained selection, two controls, DOM updates and no-JavaScript fallback. Verify no duplicate names or premature value changes while typing.

Record the exact built head/assets and preserve browser evidence. Distinguish synthetic composition and DOM accessibility assertions from real IME/screen-reader testing.

Consumer handoff
Document Baseplate’s account Timezone override field as the evidenced candidate:
 internal/baseplate/account/pages/interface.templ
 internal/timezone/options.go

Preserve its empty-string inheritance meaning and saved non-curated zones. Do not modify Baseplate in this task. Report the library/asset publication prerequisites for later adoption; no release is authorized.

Review and publication
Use Pipeline’s single integrated owner and proportional review/publication procedure. Preserve required plan/material-scope approval, independent review, applicable designer acceptance and affected-case rechecks.

Bind authentic source/review evidence during the run. Reuse unchanged valid coverage; do not repeat reviews merely while waiting for CI.

Commit/push, publish through the installed reviewed-PR procedure after its pre-PR gate, register the PR with T3 and verify its remote head. Keep it draft until applicable final-head CI, feedback settlement and UI acceptance pass.

Leave the PR unmerged. No tags, releases, issue closure or consumer rollout.

Update Project #12 to In progress or Review when accurate. Return the PR URL, public API, gallery review commands/URL, verification results, precise gaps and one next action. Include the actual model/cost report once.
