# Integration seams — Stage 3.5 reconciled inventory

Frozen plane source: `spoonman136668-ai/ckb-plane@539d56fec273c8851aea131cf4d31425a62f7250` (`snapshot/source/`). Exact-source Windows closeout is recorded in ckb-plane as `acceptance/frozen-contract-closeout-20260911.json`.

The source comparison is complete enough for an **isolated integration proof**. This does not authorize live activation. In particular, the frozen plane currently validates `worker.preferred` as only `fake` or `codex`; a live Wingless worker requires an explicit future plane-side contract change rather than aliasing or bypassing that rule.

| Seam | Classification | Owner / disposition | Reconciled rule |
| --- | --- | --- | --- |
| worker | `ADAPTER_REQUIRED` | Plane lifecycle; Wingless intelligence | Plane contract is synchronous `Worker.Run(ctx, workspace, WorkOrder, stdout, stderr)`. Wingless `SubmitWork/Status/Cancel/Result` remains internal. `integration/ckbplane` is inactive translation only. |
| broker | `WINGLESS_OWNS` | Wingless | Model-backend routing is not a plane queue/dispatch contract. For plane-driven work the adapter forces `Repairs=0`, `MaxRepairs=0`; plane owns work-order retry/reconciliation. |
| inference | `WINGLESS_OWNS` | Wingless | Inference produces candidates only. `broker.Outcome.Acceptance` must remain `external_required`; the adapter rejects an acceptance claim. |
| resources | `WINGLESS_OWNS` | Wingless measurement/admission | Frozen plane config/preflight owns toolchain/execution safety but exposes no RAM/VRAM/GPU allocation lease. Wingless resource policy is trusted local adapter configuration; measurement is not reservation. |
| contextbroker | `WINGLESS_OWNS` | Wingless, bound to plane identity | ICE retrieval is advisory. Adapter binding requires exact plane work-order ID, materialized workspace path and baseline SHA; revision mismatch fails closed. |
| toolboundary | `ADAPTER_REQUIRED` | Plane write/process authority | Plane work orders define allowed/forbidden paths and execute workers inside isolated workspaces. Wingless retains bounded reads only; future writes/processes must be delegated through an explicit plane-authorized seam. |
| broker verification | `OBSOLETE` for plane retry/acceptance | Keep only for internal research fixtures | Wingless deterministic fixture verification may remain offline, but it cannot consume a plane retry budget or imply plane acceptance. Plane-driven routing therefore has zero Wingless repairs. |
| inference lifecycle | `WINGLESS_OWNS` | Wingless backend lifecycle | STOPPED/STARTING/READY/BUSY/IDLE/DRAINING/FAILED/UNAVAILABLE describe inference health only and cannot mutate plane work-order state. |
| benchmark | `WINGLESS_OWNS` | Wingless research | Benchmark verdicts are empirical/research evidence only; they never map to plane acceptance. |
| modelhost | `WINGLESS_OWNS` | Wingless independent model process | Model residency/process supervision stays independent until an explicit plane worker contract exists. No plane queue/process authority is implied. |
| Stage 3 telemetry/output | `WINGLESS_OWNS` | Wingless evidence | GPU/process telemetry and output constraints are candidate metadata only; no lease, shell authority or acceptance grant is inferred. |

## Confirmed plane-owned contracts

The frozen source confirms the following ownership and removes earlier ambiguity:

- **Work-order schema and protection gates:** `snapshot/source/internal/plane/model.go`. `WorkOrder.Validate` requires a bounded allowed-path scope, accepted-ref protection and clean-baseline protection, and only accepts `fake` or `codex` workers.
- **Worker boundary:** `snapshot/source/internal/plane/worker.go`. `Worker.Run` receives the plane-selected workspace and immutable `WorkOrder`; `CodexAdapter` explicitly states that worker success never grants acceptance.
- **Authorization contract:** `snapshot/source/internal/plane/contracts.go`. Plane records an immutable work-order/config/tool contract and rejects changes with `AUTHORIZATION_CONTRACT_CHANGED`.
- **Workspace identity/materialization:** `snapshot/source/internal/plane/profile_workspace.go` and `profile_workspace_materialize.go`. Plane verifies the exact baseline SHA, creates/owns isolated workspace materialization, and protects the source repo and accepted refs.
- **Retry ownership:** `snapshot/source/internal/plane/autonomy_retry.go`. Queue/work-order retry, retry limits and operator-bound escape rules are plane-owned and evidence-bound.
- **Cancel/control ownership:** `snapshot/source/internal/plane/autonomy_controls.go`. Plane maps cancel to `Engine.Cancel`, guards authority generation and serializes controller requests.
- **Acceptance/evidence:** plane remains authoritative; worker/model output is candidate material only.

## Inactive adapter contract

`integration/ckbplane` intentionally has no network transport, process launcher, queue writer, accepted-ref access or live-plane mutation. It provides only:

- normalized replay input from the verified plane `WorkOrder`/workspace facts;
- deterministic `work-order + attempt -> Wingless request ID` mapping;
- exact work-order/workspace/baseline identity validation;
- trusted routing/resource-policy translation;
- hard-disabled Wingless repair budget for plane-driven work;
- cancellation mapping into an already-running Wingless request;
- candidate-only result projection with `external_required` acceptance;
- unknown replay fields and unknown terminal states rejected fail-closed.

Replay coverage includes valid translation, stale revision rejection, work-order/workspace mismatch, missing plane protection gates, untrusted backend policy, resource/routing preservation, zero local repair budget, cancellation mapping, blocked backend outcome, acceptance-claim rejection and unknown-state rejection.

## Remaining activation blockers

These do not block isolated replay/integration proof, but they block live plane activation:

1. Frozen ckb-plane has no `wingless` worker kind or explicit Wingless worker registration; `WorkOrder.Validate` allows only `fake` and `codex`.
2. No live transport from ckb-plane to Wingless is authorized or implemented.
3. Tool write/process execution remains plane-owned; Wingless currently has only its bounded read-side tool boundary.
4. Native Windows/NVIDIA/CUDA/model-residency qualification remains a separate Stage 3 hardware gate.
5. GitHub-hosted `native.yml` attempts for the first adapter commit failed before step 1 with `runner_id=0` on both Windows and Ubuntu. This is recorded as runner-infrastructure unavailability, not a test result.

No live queue, worker, CKB runtime, Coinbase/broker path, credentials, accepted refs or production state were touched by this reconciliation.
