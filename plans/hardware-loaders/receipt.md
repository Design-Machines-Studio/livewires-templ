# Hardware loader port

User authorized adding the latest Live Wires loader and component updates. Found existing draft PR26 and reused its branch at 7c9ac2b1a4004b7d8852e4dcb925d4d7339bff39. Main base e38d79acaf77025e0e6d0bbfdcb0e8a4111b61a2. Producer v2.24.0: 3042ffe56c36b8e51e4e8a98f500b1aa04a43fd3.

Lean selection: settled producer contract; bounded overlay parity repair and gallery integration. No merge or release authorized. Original detached checkout and its untracked files preserved.

Changes: preserve outer overlay status while hidden applies to inner visual; focusable hidden label; current producer asset lock; maintained loaders gallery and browser checks. Existing PR adds Ratchet, Driver Line, Torque, busy Button, optional LoadingRegion and SaveStatus indicators, and legacy compatibility tests.

Verification: Docker templ generation, go test ./..., go vet ./... passed. Gallery verify passed including generated-file freshness and runner tests. Full Playwright gallery: 69/69 passed including 320/375/1440px loaders, focus restoration/Escape, reduced motion, forced colors and native fallback. These results cover dirty fingerprint 2fdb5f87dbf2f240e3b25608aeaf91edd8f8d12136afb704a49836943d9000b0. Final source identity/review results will be recorded separately.

Preview: /home/ned/assembly/livewires-templ-worktrees/hardware-loaders; examples/gallery/run.sh start|status|rebuild|verify|stop. Task-owned Docker-assigned port; use status for current URL. External gallery run /home/ned/.local/state/design-machines/depot/runs/gallery-gallery-0ef7bda81acb4edd839d89d43cd1d58c-91p85pmd. Evidence: examples/gallery/test-results and /tmp/templ-loaders-review-20261017.

Actual screen-reader and native browser zoom testing unverified. Designer acceptance pending. Hook activation unavailable: injected owner pointer has no accompanying native hook input; no session input fabricated. Formal producer/readiness validation pending; PR remains draft.

Review repair: pattern lane found gallery destroy/dispose mismatch; changed to dispose. Simplicity lane independently verified the repair with a Node VM check. Focused browser recheck: 6/6 passed, including new pagehide cleanup, on fingerprint 59aa25d2d535ba6b9f53e16e75d0fff06d8d54d58bff9a65776b11ca0271de66. Full suite evidence remains valid for unaffected cases.

User accepted loader appearance and overlay at 715b0a5, then requested the loading region update. The shared utility/gallery loading example now passes Driver Line and loader-panel to LoadingRegion; its standalone fixture CSS pin is v2.24.0. Existing controller semantics are unchanged. Five affected browser checks passed at fingerprint 6510e4407e00087f3a218d448b511895eaa8e762448a7360a7e811447eb46f4b. Original visual acceptance applies to unchanged loader/overlay output; the newly changed loading example is separately verified by browser automation.

UI source review completed via fallback after the initial attempt was terminated at approximately 14 minutes with no result. It found percentage floating-point artifacts (0.29 rendered as 28.999999999999996%). Fixed by rounding displayed and aria-valuenow percentages to two decimal places while retaining precise SVG progress. Regression cases cover 0.29 and 0.12345678, and both visible and accessible outputs.
