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
