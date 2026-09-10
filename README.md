# Wingless

Independent, offline-first Go foundation for a resource-efficient coding/reasoning worker. Intended to sit beneath ckb-plane after its interfaces stabilize. This stage supplies testable broker contracts and mock execution, not a production autonomous coding agent or evidence of Codex parity.

## Run locally

Requires Go 1.25 or later. No model downloads, API keys, package dependencies or live services are needed.

```powershell
$ErrorActionPreference = "Stop"
go test ./... -count=1
if ($LASTEXITCODE -ne 0) { throw "tests failed" }
go run ./cmd/ice validate
if ($LASTEXITCODE -ne 0) { throw "ICE validation failed" }
go run ./cmd/wingless demo
if ($LASTEXITCODE -ne 0) { throw "demo failed" }
go run ./cmd/wingless benchmark
if ($LASTEXITCODE -ne 0) { throw "benchmark failed" }
go build -o .\bin\wingless.exe ./cmd/wingless
if ($LASTEXITCODE -ne 0) { throw "build failed" }
```

`demo` returns a mock candidate with `acceptance: external_required`. `benchmark` compares a mock response against a fixed expected value. It does not run a model or touch CKB. There is no installable daemon yet.

## Context retrieval

```sh
go run ./cmd/ice query Select
go run ./cmd/ice query InferenceBackend
go run ./cmd/ice query TestFloors
go run ./cmd/ice query Boundary
```

Read only the returned file ranges, expand as needed, verify against source. After source or documentation changes: `go run ./cmd/ice build`. Validation checks source hashes and all generated maps, including deletions. Indexes are advisory and cannot authorize acceptance. `implementations.json` records explicit compiler-checked assertions; test/call maps are syntax references, not measured test coverage.

See [architecture](docs/architecture.md), [integration seams](docs/integration-seams.md) and [known gaps](docs/known-gaps.md).
