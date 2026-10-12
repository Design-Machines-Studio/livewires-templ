Use $project-manager:assembly-coordinator to refresh ownership, then invoke the installed /pipeline workflow.

Implement:
https://github.com/Design-Machines-Studio/livewires-templ/issues/17
“Add a browser-reviewable component gallery”

Repository: Design-Machines-Studio/livewires-templ
Verified main: 0b579fca98af2740c8159395cff58d4fbb9cc216

executorRole: builder-deep
executorCapabilities: [read-repository, write-repository, tool-use]
executorEffort: high
noMergeOnCompletion: true

Goal
Give this library a small, maintained browser gallery that renders actual Templ components from the current checkout, uses pinned Live Wires assets, and supports repeatable local Playwright review. Eliminate temporary renderers and manually injected browser hookups.

This prompt explicitly authorizes a task-owned, local Docker gallery server and local Playwright on NED. T3 reachability, a hosted domain, tunnels and shared preview infrastructure are not required.

Start and ownership
Confirm the remote, fetch current main, read repository instructions and applicable skills, and check the agent-board inbox. Record the exact execution base.

Create:
 /home/ned/assembly/livewires-templ-worktrees/component-gallery-17
Branch:
 feat/component-gallery-17

Preserve unrelated files and existing worktrees.

Inline-edit #16 is being finished separately at candidate 572c83cf350bba7076df7a29427391d7dd3004ed. Start infrastructure and existing examples from main. Do not cherry-pick, copy or modify that owner’s implementation. Incorporate its example after its merge and refresh the final base. If it remains unmerged, report only that example’s acceptance as pending; continue independent gallery work.

Design
- Keep the gallery in this repository as development tooling. Live Wires remains the CSS/JS source; Baseplate remains a consumer integration check.
- Render through actual generated Go/Templ components from the selected worktree, not copied reference HTML or installed older module versions.
- Provide a simple index and stable example URLs. Begin with relative-time, copy, loading, toast, sortable and, once merged, inline editing. Reuse existing example components where possible.
- Clearly distinguish simulated application outcomes from real persistence. Do not add production authentication, database storage or application services.
- Keep the UI simple and useful for component review. No full documentation-site redesign or exhaustive catalog migration.

Assets
- Pin the exact compatible Live Wires revision and record asset checksums/provenance.
- Serve the required assets locally from a reproducible development cache or verified source. Document acquisition and failure handling.
- Do not depend on a mutable producer checkout, running Vite server, floating CDN URL or manually copied assets.
- Do not vendor a parallel CSS/JS library or add CSS/runtime assets to the consumer-facing Go component package.
- Reuse existing asset-lock conventions where suitable.

Docker and lifecycle
- All Go/Templ generation, compilation and serving run inside Docker.
- Provide documented start, status, rebuild and stop commands, using existing repository conventions where practical.
- Bind only to loopback. Allocate or explicitly configure an unused port; never take over another process’s port.
- Isolate Compose/resources by worktree/run. Stop and cleanup must affect only resources owned by this gallery instance.
- Support two worktrees running simultaneously and restarting one without affecting the other.
- Record checkout, branch, commit, dirty state, asset revision and actual URL. Tie browser evidence to the exact built source; don't label dirty output as clean-commit proof.
- Resolve installed workflow tools using their documented supported procedure. Do not hardcode user-specific, versioned plugin-cache paths.
- Inspect examples/sortable/run.sh, which currently hardcodes an old Kernel path. Reuse useful behavior without copying that coupling. Avoid a broad runner rewrite unless required for gallery integration.
- Keep ordinary gallery startup independent of review approval/coverage state. Preserve any required Docker ownership protections without turning the gallery into a new workflow engine.

Interaction examples
Maintain repository-owned setup and teardown code using the public producer APIs. No temporary Playwright injection to make an otherwise inert page work.

Examples should expose deterministic scenarios needed for review: normal operation, delayed response, failure/unknown result, recovery and unrelated updates where relevant. Retain meaningful native/no-JavaScript fallback. Keep demo transport and state clearly separate from production component code.

Browser acceptance
Use local Playwright in the same environment/network as the gallery. This is explicitly authorized even if T3 cannot reach NED localhost. Do not spend the task repairing T3 connectivity.

Commit the reusable browser test script/configuration and document its setup and command. Reuse installed tooling where suitable; make any required test dependency reproducible.

Prove:
1. A clean worktree can acquire pinned assets, generate/build and start the gallery using documented commands.
2. Each selected example loads its actual generated markup and required assets without page errors.
3. Desktop and mobile screenshots identify the tested example and source revision.
4. Representative interactions and native/no-JavaScript fallback work. Reuse existing component cases instead of duplicating every producer test.
5. Rebuild serves updated worktree output.
6. Two worktree instances remain independent.
7. Stop releases only owned resources; restarting works.

Local browser proof is sufficient. Report real IME and screen-reader testing separately from simulated events and DOM assertions.

Run applicable Docker tests/vet and generated-file checks. Keep tests focused on rendering, routing, asset availability, lifecycle and meaningful interactions.

Workflow and closeout
Use Pipeline’s proportional integrated review/publication procedure, preserving required plan/material-scope approval, independent review and applicable designer acceptance.

Resolve current installed bundles rather than retaining vanished cache paths. Reuse valid exact-head evidence and perform only affected rechecks after repairs. Do not repeat reviews merely while waiting for CI.

Publish through the installed reviewed-PR procedure after its pre-PR gate. Register the PR with T3 and verify the remote head. Keep it draft until applicable final-head CI, feedback settlement and UI acceptance pass.

Leave the PR unmerged. No releases, tags, consumer rollout, shared preview changes or unrelated repairs.

Finish with:
- PR URL and exact head.
- Copy-ready start, browser-test, rebuild and stop commands.
- Actual local URL and whether the task-owned server remains running.
- Asset revision, completed examples, verification results and precise gaps.
- A short explanation of how the next component author adds a gallery example.
- One next action and the actual model/cost report once.

Preserve required evidence and any agreed running preview. Remove only owned disposable resources.
