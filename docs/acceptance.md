# Initial-stage verification — 2026-09-10

Verified using Go 1.25.5 on Linux:

- `go test -buildvcs=false -race ./... -count=1 -timeout 60s`: passed, all eight tested packages (18 top-level tests plus subtests).
- `go run -buildvcs=false ./cmd/ice build` and `validate`: passed.
- ICE queries were used to locate routing and runner symbols before consequential follow-up edits.
- `go run -buildvcs=false ./cmd/wingless demo`: result_ready, external_required.
- `go run -buildvcs=false ./cmd/wingless benchmark`: fixed mock fixture correct; external_required remains unchanged.
- `GOOS=windows GOARCH=amd64 go build -buildvcs=false -o bin/wingless.exe ./cmd/wingless`: passed cross-compilation.

Native Windows execution has not been tested. No real inference benchmark or live-plane proof was performed. Tests use fake loopback HTTP servers and fixture data. The initial cancellation-server test fixture was corrected to bound its own wait and drain the HTTP request body; final regular and race suites passed.

Reproduce on Windows using the README PowerShell block. `-buildvcs=false` is only needed in environments with incomplete VCS metadata; do not change global Go settings.

Rollback: revert the relevant coherent Git commits, or discard an unused Wingless checkout. Nothing is installed into ckb-plane, no service is registered, and no CKB state migration exists. Do not use rollback to change live plane authority.
## Local-model qualification hardened acceptance — 2026-09-13

Windows PowerShell qualification exposed a false-green class: a model candidate can pass the visible targeted suite, full suite, and build while still broadening routing semantics. The qualification gate therefore fails closed.

For known synthetic single-mutation qualification scenarios, acceptance requires all of the following:

- schema-valid structured candidate operation;
- exact preimage match before external application;
- immutable tests;
- bounded production scope;
- targeted deterministic acceptance;
- scenario-specific semantic counterexample acceptance;
- complete suite and build;
- gofmt cleanliness for Go source;
- exact source restoration to the known-good source for the injected single-mutation scenario.

A visible-suite pass is insufficient when any hidden semantic counterexample or exact-restoration gate fails. Model output remains candidate-only and acceptance remains external_required.

The Stage 2 broker fallback scenario proved this requirement. A Devstral bounded repair changed the fallback break guard to pass == 1; the ordinary suite passed, but a counterexample proved that routine fast routing could be retried as a fallback pass. Authoritative source permits fallback only from deep selection to fast selection.

Qualification status for Devstral Small 2 IQ3_M after the trust ramp:

- base real-patch qualification: PASS;
- trust-ramp Stage 1: PASS;
- trust-ramp Stage 2: REJECTED after bounded repair due latent routing-semantic drift;
- Stage 3: NOT AUTHORIZED.

Do not widen model trust from visible-suite success alone.
