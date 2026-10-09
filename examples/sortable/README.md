# Persistent sortable example

This small Go/Templ/Datastar consumer saves two orders in one JSON file. The
library supplies reusable markup; this application owns sessions, permissions,
validation, revisions and durable saving. It uses the existing module and
toolchain with no new Go dependency.

## Serving runbook

Application: `livewires-templ-sortable`. Selected physical checkout:
`/home/ned/assembly/livewires-templ-worktrees/sortable-components-9`. Feature
branch: `feat/sortable-components-9`. The approved component predecessor is
`ff99a5e5f0af82bb5bf57daaa2773ff35de18dc8`. The source head used for each build
comes from `git rev-parse HEAD` and is reported by `status`, `/healthz` and the
visible page marker. Commit implementation before the host starts the target.

The host/root allocates a fresh external root using the installed Workflow
Kernel `owned-run-start` helper and supplies its returned absolute path as
`SORTABLE_RUN_DIR`. Use the trusted installed launcher at
`/home/ned/.codex/plugins/cache/depot/workflow-kernel/0.27.0/skills/workflow-kernel/references/workflow-kernel-launcher.sh`.
For example, the host chooses a unique `SORTABLE_ALLOCATION_ID`, runs
`rtk bash "$WORKFLOW_KERNEL" owned-run-start --workflow pipeline --run-id "$SORTABLE_ALLOCATION_ID"`,
then sets `SORTABLE_RUN_DIR` to the returned `path`. Do not create or adopt a
plain directory. The lifecycle is deliberately limited to `ned` on `ned9000`
using the verified rootless daemon at `/home/ned/.local/share/docker`.

From the selected worktree, these are the implemented entry commands:

```sh
rtk bash examples/sortable/run.sh assets --run-dir "$SORTABLE_RUN_DIR"
rtk bash examples/sortable/run.sh start --run-dir "$SORTABLE_RUN_DIR"
rtk bash examples/sortable/run.sh status --run-dir "$SORTABLE_RUN_DIR"
rtk bash examples/sortable/run.sh rebuild --run-dir "$SORTABLE_RUN_DIR"
rtk bash examples/sortable/run.sh stop --run-dir "$SORTABLE_RUN_DIR"
```

`assets` records absent task directories through `owned-run-create`, downloads
only locked bytes and retains both MIT notices. Existing incorrect bytes fail
verification. It does not start Docker. `start` checks assets/configuration,
copies the relevant source into `build/source`, fingerprints those bytes and
embeds a receipt through Go linker flags. The source snapshot is the builder's
read-only `/app`; the application executes that resulting binary. The receipt
is compiled into the binary, and the lock and bridge are embedded Go resources.
The live source must still match the compiled receipt for status to say ready.
Changes to source, generated files, branch/head or relevant dirty state require
`rebuild`. Planning, Airlift and Kernel scratch files are excluded from the
source fingerprint.

The exact repository Docker boundary is the argv prefix
`docker compose --project-name <unique recorded project> --project-directory <physical checkout> -f <physical checkout>/examples/sortable/compose.yaml`.
Creating calls are `run --rm --build builder` and
`up --detach --wait --wait-timeout 30 app`. The wrapper sends each original argv
to Kernel `plan-compose` with the exact repository project binding, materializes
the returned label-only override, executes only the returned argv and immediately
calls `record-create`. Kernel inventory projections are retained before and
after each call, including partial failures. There is no raw creation fallback.
The wrapper initializes its unique shadow run under `.workflow-kernel/runs/`;
that scratch state and returned overrides must never be staged.

The container listens on port 8080. Docker assigns the host port on
`127.0.0.1`. Use only `actualOrigin`, `consumerUrl` and `referenceUrl` reported by
`status`; never guess a workstation address or fixed port. Status reports
checkout, branch, exact head, relevant source fingerprint/dirty state, compiled
receipt, readiness, project/registry location, producer revision, reference
HTML digest and loaded include/asset digest identities. It validates active
Kernel ownership before reporting a live target.

| Task-owned path | Contents |
|---|---|
| `assets/` | Locked producer HTML/includes/CSS/modules, Datastar and notices |
| `data/orders.json` | Both keyed orders and independent revisions |
| `cache/` | Task-only module/build cache |
| `build/` | Frozen source, compiled binary |
| `evidence/` | Lifecycle binding, build/binary/image receipts, exact Docker argv, inventories and cleanup receipts |

`stop` plans reconciliation for this recorded run only and executes/records
guarded Kernel cleanup steps. It never executes the unrelated stale-sweep plan,
uses `compose down`, prunes Docker or removes unknown resources. It preserves
the order file, downloaded bytes and cache. `rebuild` performs that same owned
stop, recompiles and starts the same application with the same store. On SIGTERM
the Go server closes its listener and drains HTTP handlers before exiting;
writes are synchronous. A failed lifecycle retains its exact evidence for the
host to inspect and recover with `stop` then `start` or `rebuild`.

Host/root owns live serving, browser evidence, image removal and final filesystem
cleanup. After the final browser/restart proof, run `stop`, verify its receipt,
remove only the exact task image recorded in `evidence/image-receipt.json` if
still owned and unused, then finish the host's exact run using `owned-run-finish`.
Workers do not start this target or clean parent resources.

## Personas and request boundary

Open `consumerUrl`; `/sortable` redirects a missing session to `/session`.
GET `/session` issues an anonymous, random server-held session and CSRF token.
Submit the native persona form. POST validates origin, CSRF and its allowlist,
rotates the cookie/session and returns 303 to `/sortable`.

| Published fake persona | Reading | Member requirements |
|---|---|---|
| `editor` | Read and reorder | Read and reorder |
| `reading-editor` | Read and reorder | Read |
| `viewer` | Read | Read |

Anyone who can reach this isolated demo can choose any of these personas. This
is **not deployable authentication**. Sessions live in this process, expire after
eight hours and invalidate on restart. Cookies contain opaque random tokens,
use HttpOnly/SameSite Strict and intentionally omit Secure on the local HTTP
origin. No request actor field conveys permission. No Baseplate, Fixture,
install or operator API is involved.

Both mutation transports call the same validation/store operation. Native forms
POST to `/sortable/lists/reading/move` or
`/sortable/lists/requirements/move` with exactly one each of `itemId`, `before`,
`revision` and `csrf`. Include the `before` key even for empty append. Accepted
moves return 303; validation/stale/storage failures return a rejection page with
the authorized current list and refreshed boundaries. Cross-list IDs,
self-before, malformed/repeated fields, wrong revisions and forbidden actors
never write.

Enhanced requests use Datastar `@post`, with the unchanged producer event detail
`{itemId, before, requestId}`. HTTP-only fields add CSRF, revision, a per-request
connection marker and allowlisted presentation choices. `requestId` is echoed
only. It grants no authority and provides no idempotency. A successful change
increments only that list's revision; replay with the previous revision conflicts.
Two same-revision changes serialize to one commit and one conflict. A valid
current-revision no-op may succeed without writing or incrementing.

The file mutex encloses revision checking, ordering and temporary-file
write/sync/close, rename and directory sync. Memory changes only after the durable
result. After a post-rename error the store reloads the authoritative file and
reports uncertainty. If that reload fails, saving becomes unavailable and the
application sends no cached order as authoritative markup. Restart/reload is
required. A corrupt existing file fails startup instead of being reseeded. One
application process owns this store; running two processes against it is outside
this example's contract.

## Assets and confirmation adapter

The producer authority is
`Design-Machines-Studio/livewires@21cff3f08d156feabd9da4bbd3d366781c463f96`.
`assets.lock.json` pins actual source URLs, revisions, loaded versions, SHA-256
digests, sizes and notices. The application loads verified bytes at startup into
a read-only URL allowlist. It exposes no directory or order/session file server.
The unchanged reference snapshot appears at `referenceUrl`, including its original
inline synthetic responder, absolute include URLs and optional modules. That
responder never runs on `/sortable`. The reference reload returns static order;
only the application proves durable saving.

The required `main.js`, `main.css`, sortable CSS/module and print CSS were inspected
at the exact producer revision. Their distributed files are self-contained;
no transitive module import, CSS import or URL resource was found. `main.js`
expands `html-include`; the reference must actually expand the three includes
before visual comparison. The original head names two optional favicon URLs
whose pinned sources return 404. No substitute bytes are supplied. Google Flow
Block font requests are unrelated external requests and remain unpinned, as
recorded in the lock. No typography overrides are authored here.

The selected official free bundle is Datastar **v1.0.4**, distributed at
`https://cdn.jsdelivr.net/gh/starfederation/datastar@v1.0.4/bundles/datastar.js`,
SHA-256 `727844adfc825ee651fb93c544a2a739986f9a21820a94524b35f0cac470cf91`,
source revision `1efcdc3cb336ec3e9139491604e770e0329657bc`. Its actual bytes were
checked against the official [Actions](https://data-star.dev/reference/actions)
and [SSE](https://data-star.dev/reference/sse_events) references before directives
were written. The verified options are `payload`, `retry: 'never'`,
`retryMaxCount: 0` and `requestCancellation: 'disabled'`. The zero retry budget
also covers this bundle's unconditional network-error retry branch. Both lists
use distinct endpoints/signals and independent request markers.

The server emits `datastar-patch-elements` with `mode replace` targeting only the
`ol`/`tbody` child under a CSS-encoded active-request marker. It preserves the
host, table headings, status and event bindings. Replacement rows carry pending
control cues, refreshed native fields, revision and delivery marker. The following
`datastar-patch-signals` carries list identity, request ID, explicit status,
revision and marker. The result signal is a scalar JSON envelope, replaced as a
whole. The bridge decodes that explicit result; recursively merging result
objects could inherit an old accepted status when a later response omits status.
Bundle inspection shows synchronous element-patch handling before the following
signal patch can trigger `data-effect`. **Host browser proof
of that actual delivery order is still required.**

`bridge.js` is the sole authored JavaScript. The Datastar substitution table has
no generic custom-element confirmation entry. This narrow adapter records the
host/request/connection marker, provides the HTTP payload and hands an explicitly
delivered result to `resolveMove`. It requires a matching current host, marker,
ID, accepted/rejected status and rendered child revision/delivery marker. A
disconnection observer invalidates markers; it never treats child changes as
confirmation. Missing/invalid status, stale IDs and unrelated patches leave
pending intact. Datastar performs all HTTP/SSE; the producer performs dragging,
keyboard moves, animation, rollback and focus restoration.

A rejected Datastar `FetchFailed` from its exhausted zero retry budget can settle
only its matching current request. The message says saving could not be confirmed
and reload is needed. It never claims a server write was undone. Successful
responses with missing/invalid application status are not transport failures.

## Deterministic browser handoff

Use T3 `preview_status` first, then `preview_open` if no automation-capable preview
is attached. Bind evidence to the URLs from the current status receipt. Before
capture, verify `/healthz`, the source/build marker, registered rootless ownership,
reference HTML digest, include expansion and each actual loaded stylesheet/module
against the lock. Use the same origin for both routes at **375 and 1440 pixels**.
Capture source/DOM/screenshots and targeted computed font-size, font-weight,
color, padding, margin, background-color and border evidence for the reading and
table fragments. Application session/presentation/case/build controls sit outside
those fragments. Global DOM/pixel/theme equality is not the acceptance test.

Initial reading IDs are `essay-01` (Field notes), `essay-02` (On paper) and
`essay-03` (An index). Initial requirement IDs are `orientation` (Member
orientation), `agreements` (Shared agreements) and `governance` (Governance
participation). Seed occurs only for an absent task-owned store. Do not delete,
reset or reseed an existing store for evidence; start from the observed order and
record the exact starting order/revisions. An up arrow moves before the previous
item; a down arrow moves before the item two positions ahead or appends.

Select response cases using the per-list buttons after loading the page. They
POST through Datastar with the same session/reorder/CSRF gates as moves. Each
case affects only the next move for that session/list. `normal` supersedes a
selected case. Reload clears selected cases and returns current disk order;
it is the recovery for interrupted or unknown confirmation. Diagnostic results
finish automatically after three seconds; `delay` waits five seconds outside
the store lock. These test controls require JavaScript; ordinary arrows remain
native forms.

| Case | Exact steps | Expected result |
|---|---|---|
| 1. Durable two-list moves | Choose editor. Put essay-02 before essay-01 and governance before orientation using handles or repeated arrows as needed. Reload. Host runs stop/start with the same run directory, then chooses editor again. | Both orders survive reload/restart; each list advances only its own revision. |
| 2. Rejected write | Select reading `reject`, then request a reading move. | Authoritative order stays unchanged; pending clears only after matching rejected result and announces rejection. |
| 3. Concurrent lists | Select reading `delay`, move reading, then move requirements before five seconds elapse. | Reading shows provisional order and blocks repeat controls; requirements saves before reading settles. |
| 4. Invalid replies | Separately select reading `stale`, `missing`, then `invalid`; request a move after each selection. Observe the first result while pending. | Old request ID, missing status and unknown status do not settle. Matching real patch/result after three seconds settles. |
| 5. Unrelated children | Select reading `unrelated`, request a move. Retain row-node identities before/after the first patch. | Unrelated authoritative replacement does not settle. Controls retain pending cues. Matching rejection after three seconds keeps replacement nodes rather than rolling back original disconnected nodes. |
| 6. Keyboard staging | Focus a handle; Space or Enter starts staging, arrows change insertion and Space/Enter commits. Repeat with Escape, then repeat and tab out of the staged row. | Commit saves; Escape and leaving the row cancel without writing. |
| 7. Focus | Move with a focused arrow/handle, observing immediate movement and replacement. In arrows-only mode move a row to a boundary. During a delayed request deliberately focus a control in the other list. | Immediate movement keeps its focused descendant. Replacement restores the same enabled arrow, otherwise handle/enabled arrow/nearby handle/status. Deliberate focus elsewhere is never stolen. |
| 8. Nested/pointer controls | Use both reading links and the An index input. Drag a handle; cancel the pointer and test release outside the list. Send an additional pointer during an active gesture. | Nested controls stay usable; cancellation/outside drop does not write; the extra pointer cannot replace the gesture. Label synthetic pointer evidence synthetic. |
| 9. Disconnect | Select reading `delay`, request a move; remove and reinsert that same host through the isolated browser tool. Start a new reading request before the old response arrives. Record markers and request IDs. | Old child selector/result cannot overwrite or settle the new request. Listeners are not duplicated. The old server write may still succeed; a new stale revision is honestly rejected with current state. |
| 10. Native fallback | Disable JavaScript before navigation; choose editor through /session. Use arrows in both list/table. Inspect POST/303, reload, form boundaries and revisions. Submit a cross-list or stale POST. | Orders persist, boundaries update, invalid POST does not write. Handle-only is excluded from this no-JavaScript claim. |
| 11. Authorization | Choose reading-editor and then viewer. Inspect absent forbidden controls and attempt native/enhanced POST bypasses. Try missing/expired/restarted cookies, bad/missing CSRF and disallowed Origin. | Reading-editor cannot reorder requirements; viewer cannot reorder either. Denied requests disclose no forbidden list order and do not write. |
| 12. Presentations | Use the native selectors for cards/dividers, default/compact and handle/arrows/both at both widths. Default matched case is reading cards plus lined table, combined controls. | Hierarchy remains readable, controls usable, semantic table rows/headings/caption/rules intact and actions visible. Handle-only is documented as requiring enhancement. |

On the reference route also exercise moves, reload, staging/cancel and focus.
Its responder settles synchronously and resets on reload; server latency and
durability belong only to this application. Record the permitted persistence,
session/authorization/validation/CSRF/URL and table-native-form differences.
Retain physical-touch, active reduced-motion, direct-paint and formal
Playwright/WebKit/Gecko gaps unless actually tested. T3 evidence keeps its actual
transport identity. Failed required assets/declared routes leave matched visual
parity incomplete after bounded recovery; source inspection cannot substitute.

## Auth Boundary Map

| Surface | Actual action/resource and gates | Default deny / tests |
|---|---|---|
| GET `/session` | Public local persona page; issues anonymous server-held session and CSRF context | Anonymous context has no read/reorder grant. Rotation/restart coverage: `TestSessionAndGrantBoundary`, `TestHealthReceiptAndRestartSessions`. |
| POST `/session` | Existing session, exact origin/CSRF, closed fields and allowlisted persona; rotate cookie/session | Unknown/stale session, persona or invalid CSRF denies. `TestSessionValidationAndOrigin`. |
| GET `/sortable` | Handler requires actual actor read grants for reading and requirements | Missing/unknown/stale actor goes to persona page. Move controls use the same `grant(actor,"reorder",list)` as POST. `TestSessionAndGrantBoundary`. |
| POST reading `/move` | Handler requires read/reorder reading; strict fields, CSRF/origin, reading membership and per-list revision; rechecks session after delays | Denies unknown actor/list/action, viewer, stale session and malformed inputs. `TestStrictValidationAndCSRF`, `TestSessionAndGrantBoundary`, `TestOrderedSSEAndCorrelation`. |
| POST requirements `/move` | Same handler gates, using requirements grants/membership/revision | Also denies reading-editor. Native and enhanced paths share store operation. Same tests plus `TestConcurrentEndpointsAndExampleCases`. |
| POST each list `/case` | Existing session and actual reorder grant for that list; closed csrf/case fields, origin/CSRF and bounded allowlist | Controls are absent without permission. `TestSessionAndGrantBoundary`, `TestSessionValidationAndOrigin`, `TestConcurrentEndpointsAndExampleCases`. |
| GET `/healthz`, locked assets, `/bridge.js` and reference | Public local build identity/documentation; fixed read-only asset mapping | No order/session file traversal or mutation. `TestReferenceBytesAndAllowlist`, `TestHealthReceiptAndRestartSessions`. |

Handlers own these gates directly; there is no separate authentication middleware
or Baseplate Authorizer action/resource pair. Exact tests live in
[http_test.go](http_test.go) and [store_test.go](store_test.go). Unknown action/list
checks fail closed. Residual limits are the openly selectable fake personas,
process-local session invalidation, loopback HTTP and one process/file owner.
There are no operator/install edges to claim or test.

## Verification and delivery state

```sh
rtk docker compose run --rm dev templ generate
rtk docker compose run --rm dev go test ./examples/sortable
rtk docker compose run --rm dev go test -race ./examples/sortable
rtk docker compose run --rm dev go test ./...
rtk docker compose run --rm dev go vet ./...
```

Commit generated `page_templ.go`; never hand-edit it. Preserve the existing
generator pin despite its informational runtime-version warning. Root compares
generation freshness again against the integrated candidate. The focused tests
cover persistence/reopen, duplicate/no-op/concurrent writes, real temporary-file
I/O/rename/directory-sync failures, auth/CSRF, strict fields, native boundaries,
independent list requests, ordered SSE and reference byte equality/allowlisting.
These are source/HTTP results. They do not prove browser patch delivery, focus,
computed styles or a live owned serving lifecycle.

Ambiguity receipt: `ambiguity_resolved: true`. Chose one-shot per-session/list
diagnostic cases and replacement of only the order container. Rejected a replay
ledger, generic confirmation framework and replacement of the outer sortable
host. Those alternatives add machinery or violate the producer contract.

Host/root must capture all twelve browser cases and actual assets/start/status/
rebuild/stop receipts before integrated final review/publication. Missing runtime
or browser evidence remains NOT-COVERED. Root owns review, bounded repairs,
publication/CI, T3 registration and exact-owned cleanup; leave the PR unmerged.
