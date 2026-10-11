# Issue #10: utility and feedback hooks

## Goal

Expose relative-time, copy, operation-local loading, and managed toast markup
that follows the merged Live Wires producer contracts while preserving this
library's native fallbacks and existing API.

## Smallest usable implementation

Add Props and convenience wrappers for relative time, copy, loading regions,
and managed toast regions. Keep the library CSS-free and do not bundle producer
JavaScript. Reuse pinned optional Live Wires modules in the bounded example.
Keep existing Loader and Toast calls compatible.

## Approved scope and acceptance

The user's implementation request supplies this scope and approval source.

- Relative time retains useful absolute content, including date-only values.
- Copy retains visible selectable content and reports success only after the
  clipboard promise resolves; failure and missing clipboard support provide
  manual instructions.
- Loading state is local to its operation region. Status stays outside busy
  content, and pending, failed, cancelled, unknown and application-confirmed
  results remain distinct.
- Managed toast markup has stable region/list boundaries, stable occurrence
  IDs, semantic kind, optional duration, message hook, child actions and native
  dismiss buttons with correct type and labels. Standalone toast semantics
  remain intact; managed messages do not duplicate region announcements.
- Controller-owned attributes, lifecycle removal/reinsertion, ID reuse, sticky
  actions/errors, cleanup and focused-subtree preservation are documented.
- Applications own persistence, retries, permissions and confirmation. No
  change to Baseplate's approved Datastar Pro or clipboard policy.
- Add semantic render/compatibility tests and a bounded example using the
  producer assets pinned to merged revisions.

## Boundaries

No library CSS, copied JavaScript state machines, persistence behavior, retry
policy, permission policy, Baseplate adoption, releases, tags, issue closure,
or merge. The task must leave the PR draft and unmerged.

## Verification and review

Run Templ generation, tests and vet inside Docker, commit generated files, and
prove generation is clean. Compare rendered markup with producer references.
Run the declared relative-time, copy, loading, toast, responsive and
reduced-motion cases against the exact candidate and pinned assets. Preserve
host browser evidence. Reuse the no-finding independent source review at
`466486241d430a7b3ca0a06e36f14ebcb5e0c1e7`; perform no duplicate review loop.

## Execution controls

- Workflow: Lean, selected by the installed Pipeline selector. The user-set
  executor was `builder-fast` with `read-repository`, `write-repository`, and
  `tool-use` at high effort.
- Repository: `Design-Machines-Studio/livewires-templ`.
- Base: `main` at `5b274c89286845c98a329c0add42308ff2808cdf`.
- Branch: `feat/utility-feedback-10`, reuse expected head
  `466486241d430a7b3ca0a06e36f14ebcb5e0c1e7`.
- Rendered surface: required; target `/home/ned/assembly/assembly-atlas` at
  `https://dm026.asmbly.app`, as authorized by the user during this task.
- No merge on completion.
