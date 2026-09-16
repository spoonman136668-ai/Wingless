# CR-1A — Cognitive Runtime Foundation

Status: **candidate / external acceptance required**

Base: `model-qualification-v1` at `8152e8c57826007c1e84aa396107afd40ba28ef4`.

CR-1A is an intelligence-side research foundation only. It does not activate Wingless through ckb-plane, mutate accepted refs, execute stored model-generated code, grant filesystem/process authority to model output, or alter the sealed Stage-1 candidate boundary.

## Architecture

The new `cognitive` package is deliberately not wired into the live Wingless command or ckb-plane adapter.

```text
RunRequest
   |
   v
deterministic Router
   |---------------- exact memory ----------> candidate-only result
   |---------------- validated static skill -> candidate-only result
   |
   `---------------- inference
                         |
                         v
                 bounded pass policy
                         |
                         v
              external deterministic evaluator
                         |
               sufficient / hard stop
                         |
                         v
                 candidate-only result
```

All results retain `acceptance = external_required`.

### Bounded cognition

`Policy.MaxPasses` is a trusted integer bounded to `1..8`. Model output is never parsed as a request to increase the budget. Backend failures are terminal for the cognitive run and are not retried internally; ckb-plane retry/reconciliation remains separate.

Additional passes reuse the original request deadline and output-token bound. The previous candidate is appended only as explicitly labeled untrusted context. If the next pass would exceed `MaxContextBytes`, the run terminates as `cognitive_context_limit`.

Cancellation is checked before routing and before every inference pass and is propagated to the existing inference backend.

### Memory

`FileStore` provides a simple replaceable persistent memory abstraction with:

- semantic, episodic, procedural, and failure classes;
- versioned records;
- SHA-256 content-addressed IDs;
- exact-key deterministic retrieval;
- inspection-friendly JSON files;
- deletion by content ID;
- collision/tamper detection for existing content-addressed records;
- hard per-record size limits and hard record-count caps.

Memory is candidate context only. A memory value cannot change the runtime's `external_required` acceptance state.

The store root is supplied by trusted host configuration. No model output selects memory filesystem paths.

### Skills

Skills are persisted separately from memory and begin in `candidate` state.

A candidate cannot self-declare validation. Promotion requires explicit deterministic `ValidationEvidence` containing a validator identity and SHA-256 evidence digest.

Only `validated` `static_text` skills can bypass inference in CR-1A. Stored procedure and prompt-template skills are inspectable data only and are never executed. The runtime has no shell/code execution path.

### Routing

Routing is deterministic:

1. exact memory match, when memory is enabled;
2. validated `static_text` skill, when skill reuse is enabled;
3. bounded multi-pass inference when explicitly allowed and `MaxPasses > 1`;
4. otherwise single-pass inference.

Each decision has a stable SHA-256 decision ID and input digest for replay comparison.

### Evidence

Each run emits `wingless.cognitive-runtime-run.v1` evidence with:

- route decision;
- model-call count;
- pass count;
- per-pass reason;
- backend/model identity;
- input/output token counts when reported;
- latency;
- inference telemetry, including existing resource telemetry;
- evaluator reason;
- termination reason;
- total known token/latency cost;
- reuse flag;
- fixed `external_required` acceptance.

Unknown inference telemetry remains unknown; CR-1A does not manufacture missing measurements.

## Deterministic repeated-task fixture

`benchmark.RunCognitiveReuseFixture` is a research fixture for the reuse/accounting seam. It is **not** a claim about real-model capability.

The control performs the same deterministic fixture through the existing `inference.InferenceBackend` seam on every repetition, with no cognitive routing or reuse.

The experiment performs the first solution through inference, checks it against an exact deterministic expected result, then the harness — not the model — creates and promotes a `static_text` skill using deterministic validation evidence. Later related tasks reuse that validated skill without a model call.

The report exposes per-run cost checkpoints at repetitions 1, 10, 50, and 100 when present.

`cmd/cognitive-bench` prints the report as JSON.

## Files

- `cognitive/types.go`
- `cognitive/store.go`
- `cognitive/router.go`
- `cognitive/runtime.go`
- `cognitive/runtime_test.go`
- `cognitive/store_test.go`
- `benchmark/cognitive_reuse.go`
- `benchmark/cognitive_reuse_test.go`
- `cmd/cognitive-bench/main.go`
- `scripts/Test-CR1A.ps1`
- `docs/cognitive-runtime-cr1a.md`

## Authoritative Windows acceptance

Windows PowerShell remains authoritative.

Run from a clean checkout of the CR-1A branch:

```powershell
$ErrorActionPreference = 'Stop'

git status --short --branch
git rev-parse HEAD
git diff --stat
git diff
git diff --cached

# CR-1A adds source files, so regenerate advisory ICE before validating it.
go run -buildvcs=false ./cmd/ice build
go run -buildvcs=false ./cmd/ice validate

.\scripts\Test-CR1A.ps1
```

`Test-CR1A.ps1` runs the full race-enabled Go suite, the focused cognitive tests, the 100-task deterministic A/B fixture, and the Wingless build. It writes candidate evidence under `evidence/cr1a/`.

After the run, inspect the regenerated `.ice` diff and the evidence before any commit or acceptance action.

## Exit-gate interpretation

CR-1A is not accepted merely because this code exists.

The gate requires authoritative Windows evidence that:

- the existing full suite remains green;
- the new cognitive tests remain green;
- hard pass bounds hold;
- memory cannot mutate authority;
- non-static stored skills cannot execute;
- deterministic routing is replay-stable;
- baseline/non-cognitive behavior remains available because no live wiring changed;
- the A/B repeated-task fixture records first-run versus reuse cost;
- no live activation or accepted-ref mutation occurred.

Do not proceed to model surgery from this branch.
