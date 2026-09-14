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

## Qwen3-Coder 30B-A3B Q3_K_M qualification — 2026-09-13

Windows PowerShell qualification used the hardened external acceptance gate on `model-qualification-v1`.

Pinned model identity:

- family: `Qwen3-Coder-30B-A3B-Instruct`;
- quantization: `Q3_K_M`;
- Hugging Face revision: `13a17f197c4d0d777e5eda136a47a17662c9f952`;
- GGUF SHA-256: `f67299d72124ed68b4cbf3a776079df289bd854a17dfe6febfc835d1e1fbdefb`;
- validated llama.cpp offload: `-ngl 32`.

Qualification status:

- resource / LocalHTTP compatibility probe: PASS;
- base real-patch qualification: PASS under the explicit byte-exact preservation contract;
- trust-ramp Stage 1 tool-boundary authorization repair: PASS;
- trust-ramp Stage 2 broker fallback repair: FAIL;
- trust ceiling: Stage 1 only;
- Stage 3: NOT AUTHORIZED.

Stage 2 is a substantive semantic failure, not a verifier or packaging artifact. The candidate retained the injected broken fallback guard (`tier == "deep"` instead of the authoritative `tier != "deep"`), changed capability filtering so a capability mismatch could still proceed to backend selection, and was not gofmt-clean. Visible `TestRoutes`, the persisted routine-fast semantic regression, broker package tests, integration acceptance, and the full suite all rejected the candidate. Full build alone remained green and is not acceptance.

No Stage 2 repair retry is authorized from this evidence. Model output remains candidate-only, deterministic external acceptance remains authoritative, and Qwen3-Coder must not be widened beyond Stage 1 trust based on its base-patch and Stage 1 success.

## Stage-1 local-model selection contract — 2026-09-13

The qualification branch now contains an unwired, fail-closed Stage-1 local-model eligibility seam in `broker/stage1_local_selection.go`. It does not invoke inference, alter `broker.Select`, apply edits, run tests, grant acceptance, or activate Wingless.

An allowed Stage-1 task must be exactly one bounded production-file repair with candidate-only structured output, exact preimage binding, immutable tests, external application, external acceptance, and zero model repair retries. Routing/architecture changes, authority changes, live actions, arbitrary shell, model-owned edits, and model-owned tests are rejected.

The preferred Stage-1 candidate generator is `qwen3-coder-30ba3b-q3km`; the secondary candidate generator is `devstral-small-2-24b-iq3m`. Both remain capped at Stage 1. This preference is a selection contract only and does not widen either model's authority.

Runtime wiring is intentionally absent. Any later caller must still pass the deterministic eligibility contract and preserve external application, testing, semantic verification, and promotion authority outside the model.

## Stage-1 authority projector — 2026-09-13

The qualification branch now contains an inactive Stage-1 authority projector in `integration/ckbplane/stage1_authority.go`. It binds Stage-1 eligibility metadata to the existing frozen plane `WorkOrderContract` and `WorkspaceBinding` before producing the broker's pure `Stage1LocalTask`.

The projector requires exactly one work-order-allowed production path expressed as a canonical repository-relative slash path, an exact lowercase SHA-256 preimage, a structured `write_text` JSON-schema candidate contract, immutable tests, external application, `external_required` acceptance, zero model repair retries, and all routing/architecture, authority, live-action, arbitrary-shell, model-edit, and model-test flags to remain false.

Repository-path validation is platform-independent: Git/repository paths are validated with slash semantics rather than host filesystem normalization. Backslashes, drive syntax, absolute paths, dot paths, and parent traversal are rejected.

The projector is intentionally inactive. It does not alter `Translate`, `broker.Select`, `Runner`, listener/service behavior, model invocation, edit application, testing, acceptance, queue authority, or live Wingless activation. Runtime wiring remains a separate future gate.

## Stage-1 candidate request builder — 2026-09-13

The qualification branch now contains an inactive Stage-1 candidate-request builder in `integration/ckbplane/stage1_candidate_request.go`.

The builder first revalidates the Stage-1 authority envelope through `ProjectStage1LocalTask`, then hashes the actual UTF-8 production-file preimage and requires it to equal the authority-bound SHA-256. It constructs a bounded `inference.Request` whose JSON schema permits exactly one `write_text` operation, exactly one production path, and exactly one expected preimage hash.

Work-order title and instructions are prompt context only and cannot widen authority. The schema and external verifier remain authoritative. The request explicitly states that application, testing, semantic verification, acceptance, and promotion are external.

This seam remains inactive: it does not select or invoke a backend, alter `broker.Select` or `Runner`, apply edits, run acceptance, mutate ckb-plane, or activate live Wingless.

## Stage-1 candidate validator — 2026-09-13

The qualification branch now contains an inactive strict Stage-1 candidate validator in `integration/ckbplane/stage1_candidate_validator.go`.

The validator accepts only one raw JSON object with exactly the four schema fields. Markdown fences, explanatory prefixes/suffixes, additional JSON values, arrays, unknown fields, non-`write_text` operations, path drift, preimage drift, empty/NUL content, no-op replacement content, stale current source, and tampered request/schema bindings are rejected.

Before validating the candidate operation, the validator rechecks the current source bytes against the authority-bound SHA-256 and rechecks that the built request's `json_schema` still exactly matches the authorized production path and preimage.

The validator is non-applying by design. It returns candidate data only and does not write files, invoke a backend, run tests, grant acceptance, mutate ckb-plane, or activate live Wingless.

## Stage-1 external application plan — 2026-09-13

The qualification branch now contains a non-executing Stage-1 external application-plan seam in `integration/ckbplane/stage1_application_plan.go`.

The planner revalidates the request binding, current source preimage, candidate operation/path/preimage/content, replacement size, and candidate postimage hash. It emits immutable plan data containing work-order/request identity, workspace, repository-relative path, expected before SHA-256, expected after SHA-256, replacement content, external application authority, `external_required` acceptance authority, and the Stage-1 trust ceiling.

This is deliberately not an application engine. There is no filesystem write, rename, chmod, shell, subprocess, test execution, acceptance, promotion, queue mutation, ckb-plane mutation, or live Wingless activation. The existing `toolboundary` package remains read-only.

## Stage-1 post-application verifier — 2026-09-13

The qualification branch now contains an inactive Stage-1 post-application verifier in `integration/ckbplane/stage1_postapplication.go`.

The verifier takes the already-bound Stage-1 request, the non-executing external application plan, and bytes observed after an external application step. It revalidates request identity, work-order identity, workspace, repository path, before hash, after hash, replacement content, external application authority, `external_required` acceptance authority, and the Stage-1 trust ceiling. It then requires the observed bytes to equal the planned content and planned SHA-256 exactly.

Success means only `postimage_verified`. Acceptance remains `external_required`; the verifier cannot emit accepted/promoted status.

This seam performs no file writes, model invocation, test execution, acceptance, promotion, queue mutation, ckb-plane mutation, or live Wingless activation.
