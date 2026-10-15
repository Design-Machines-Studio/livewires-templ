# Component gallery

The gallery renders this checkout's generated Go/Templ components. It is local
development tooling. Live Wires supplies the CSS and JavaScript; the Go component
package still contains neither. Demo requests and reorders are simulated.
Nothing is saved. Native sortable moves carry their order in the URL, so the
base example URL resets the page.

## Start and review

Run these commands from the repository root on NED as `ned`, using rootless
Docker. Go, Templ generation, compilation and serving stay inside Docker.
Python 3.12+, Docker Compose, Git and an installed Workflow Kernel are required.
Local browser review also needs Node.js and npm.

```sh
examples/gallery/run.sh start
examples/gallery/run.sh status
```

Startup acquires and verifies the locked assets, freezes the current worktree
inputs, generates and builds in Docker, then publishes an unused Docker-assigned
port on `127.0.0.1`. Read the actual `url` from status; no fixed port is reserved
or taken over. The index is `/`; example URLs are `/examples/relative-time`,
`/examples/copy`, `/examples/loading`, `/examples/toast`, `/examples/sortable`,
`/examples/inline-edit`, `/examples/combobox`, `/examples/save-status` and
`/examples/page-header`.

The runner uses the host-selected user installation from Claude's
`installed_plugins.json`, checks its plugin identity/version, and invokes the
documented Workflow Kernel launcher. The launcher resolves the compatible
runtime. Other hosts set `WORKFLOW_KERNEL` to their host-selected bundle's
`skills/workflow-kernel/references/workflow-kernel-launcher.sh`. Do not glob a
cache, select by modification time, or paste a versioned cache path into code.
Kernel is used only for Docker ownership and cleanup; startup reads no review
approval, findings or coverage state.

```sh
npm ci --prefix examples/gallery
(cd examples/gallery && npx playwright install chromium)
GALLERY_URL="$(examples/gallery/run.sh status | jq -r .url)" \
  npm test --prefix examples/gallery
```

Playwright runs locally in the gallery's network environment. T3 reachability,
a domain, tunnel or shared preview is unnecessary. The committed package lock
pins the test dependency and its browser installer. If Linux browser libraries
are missing, install them through the host's existing operations procedure;
do not run a development agent with system privileges.

Tests capture each example at 1440px desktop and 390px mobile widths, reject
page/network errors, verify served asset checksums, and exercise representative
enhanced and native behavior. Results, screenshots, traces on failure and
`source.json` files land under ignored `examples/gallery/test-results/`.
Set `GALLERY_EVIDENCE` to a fresh absolute evidence directory to preserve a run.
Screenshot names carry the example, device, source commit and content
fingerprint; the rendered footer records the same source and dirty state.
DOM assertions and synthetic keyboard events do not prove real IME or
screen-reader behavior. Record those manual checks separately.

## Rebuild and stop

```sh
examples/gallery/run.sh rebuild
examples/gallery/run.sh status
examples/gallery/run.sh verify
examples/gallery/run.sh stop
```

Rebuild stops only this instance's recorded Docker resources and serves a new
snapshot from this worktree. The assigned port can change; refresh `GALLERY_URL`.
Status compares the current checkout with the compiled receipt. A
`source-mismatch` means the running output needs a rebuild. `/healthz` reports
the immutable compiled source and asset lock, independently of runtime labels.

`verify` runs the Python runner regression checks, all Go tests and vet inside
the isolated Docker builder, then checks the complete generated-file inventory
and contents against the checkout, including newly generated files. It does not restart
a running gallery. Stale generated files are a failure; run `make generate`
inside the repository's Docker development service and commit the generated
files before retrying.

Each worktree has its own ignored `.gallery/instance.json`, randomly named
Compose project, external cache/build/evidence root and Kernel resource
registry. Two worktrees can run at once. Do not copy `.gallery` or
`.workflow-kernel` state into another checkout. Stop uses freshly inspected,
creation-recorded exact IDs through Kernel's cleanup commands. It never runs
Compose down against another project or executes the stale-resource sweep.
Restart with `start`; retained assets and caches make this cheap.

Stop retains the verified asset cache, compiled build, cache and lifecycle
receipts at the `run` path reported by status. Docker networks/containers are
released. Build images remain as Docker's reusable build cache. Keep evidence
until review is complete. For final disposal, stop first, then use Workflow
Kernel's `owned-run-finish --run-root <exact run path> --outcome succeeded`.
Remove only this worktree's `.gallery/instance.json` after successful disposal.
Do not prune shared Docker resources or delete another worktree's state.

## Locked assets

`assets.lock.json` records exact Live Wires revision
`50ce9ecd2dc8d102cd22e7444051932c6b450b08`, source URLs, SHA-256 hashes, sizes
and license provenance. Required files are fetched from that revision into the
instance's external development cache, never from a mutable producer checkout,
Vite server or floating CDN. The save-status example uses the producer's exact
source stylesheet because it is not yet a separate distribution file. No CSS
or producer runtime is committed here. The server verifies the files again
before listening. Recheck all maintained examples whenever a shared asset
changes.

`run.sh assets` acquires/verifies without starting the server. A missing or
unreachable source fails startup; no substitute or older cached revision is
used. A corrupt cached file fails closed. Remove only the named corrupt cache
file and retry `assets`. A checksum mismatch from the source requires checking
the producer commit and updating the lock deliberately; never accept new
bytes by suppressing verification. Keep the producer license with cached files.

## Add an example

1. Add a stable name to `examples` in `main.go` and a rendering branch in
   `page.templ`. Call actual Props components or a shared example template.
2. Add application-owned setup and teardown to `demo.js` using the producer's
   public APIs. Keep simulated transport/outcomes here, outside component code.
   Preserve visible content and native forms without JavaScript.
3. Add required asset entries at one compatible exact producer revision,
   including checksums, sizes and license notice. Fetch dependencies locally.
4. Add a rendering/route check and the representative browser case to
   `tests/gallery.spec.mjs`. Generate in Docker, commit `_templ.go`, run `verify`,
   rebuild and capture evidence from the resulting source.

Inline editing uses the actual wrappers merged in PR #18 at trusted main
66cf4fbb61c585943eb388ee9e2b7dbac8d2ce1b. Its application adapter uses matching
context/edit/request IDs, simulated confirmation/rejection/unknown results,
and late-result cancellation. It does not modify the producer or save text.

The combobox example uses the optional single-select assets from the exact
producer revision in the lock. It shows an optional timezone override, a
required select associated with the form from outside it, native reset and GET
submission, and application-simulated local/remote search. The remote response
control includes deterministic delay, out-of-order, error and empty-result
modes. No search endpoint is called and nothing is saved. The example is based
on Baseplate's `interface_timezone` field; an empty value means organization
default, and Baseplate's saved non-curated values remain application-owned.
