# Plane reconciliation — Stage 3.5

State: **READY_FOR_ISOLATED_INTEGRATION_PROOF**

## Frozen and accepted plane input

Wingless reconciles against one immutable plane source only:

- repository: `spoonman136668-ai/ckb-plane`
- frozen source commit: `539d56fec273c8851aea131cf4d31425a62f7250`
- frozen source root: `snapshot/source/`
- source manifest: `snapshot/SOURCE-MANIFEST.sha256`
- snapshot timestamp: `2026-09-11T16:47:33.7746436Z`
- source files: 130
- Go files: 75
- source bytes: 758340

The exact source was re-qualified on authoritative Windows PowerShell and sealed as `acceptance/frozen-contract-closeout-20260911.json` in ckb-plane. The closeout established:

- frozen manifest verification: PASS;
- live contract source equals frozen snapshot: true;
- focused post-DT regressions: PASS;
- full `internal/plane` suite: PASS;
- frozen `cmd/ckb-plane` build: PASS;
- authority: `plane`;
- controller paused with fresh heartbeat;
- accepted-ref mutation disabled;
- original GitHub runner stopped.

The observed controller was operator-blocked on an unrelated product work order during verification. That state did not invalidate or mutate the frozen contract closeout.

Historical DN/DO/DP/DT evidence remains supplemental provenance; it is no longer required to infer source identity because the exact frozen source itself has a direct Windows acceptance binding.

## Reconciled ownership

Source inspection confirms the intended one-authority architecture:

| Concern | Authoritative owner | Reconciled disposition |
| --- | --- | --- |
| Work-order schema and protection gates | ckb-plane | `PLANE_OWNS` |
| Authorization contract | ckb-plane | `PLANE_OWNS` |
| Queue, dispatch, duplicate/recovery state | ckb-plane | `PLANE_OWNS` |
| Workspace identity/materialization/lifetime | ckb-plane | `PLANE_OWNS` |
| Work-order retry/reconciliation | ckb-plane | `PLANE_OWNS` |
| Cancel/control request state | ckb-plane | `PLANE_OWNS` |
| Acceptance and authoritative evidence | ckb-plane | `PLANE_OWNS` |
| Model routing/inference | Wingless | `WINGLESS_OWNS` |
| Resource measurement/admission | Wingless | `WINGLESS_OWNS`; measurement is not a plane allocation lease |
| ICE/context retrieval | Wingless | `WINGLESS_OWNS`, bound to plane workspace and baseline identity |
| Model lifecycle/telemetry | Wingless | `WINGLESS_OWNS`; cannot mutate plane lifecycle |
| Write/process tool authority | ckb-plane boundary | `ADAPTER_REQUIRED` |
| Worker protocol translation | shared boundary | `ADAPTER_REQUIRED` |

Key frozen source references are `internal/plane/model.go`, `worker.go`, `contracts.go`, `profile_workspace.go`, `profile_workspace_materialize.go`, `autonomy_retry.go` and `autonomy_controls.go`.

The plane worker contract is synchronous `Worker.Run(ctx, workspace, WorkOrder, stdout, stderr)`. Wingless's existing `worker.Boundary` (`SubmitWork/Status/Cancel/Result`) is therefore retained as an internal in-memory service only; it is not projected as plane lifecycle.

## Inactive adapter implemented

`integration/ckbplane` now contains a pure, non-default adapter/replay layer. It has no network listener, local-controller transport, queue write, worker launch, accepted-ref access, process execution or live-plane mutation.

The adapter:

- pins the frozen plane source identity in code;
- consumes a normalized subset of verified plane `WorkOrder` facts plus plane-owned workspace binding;
- deterministically maps work-order ID + attempt to a Wingless request ID;
- requires work-order/workspace ID agreement and exact baseline-SHA agreement;
- requires the plane clean-baseline and accepted-ref protection gates;
- rejects stale workspace revisions and unknown replay fields;
- passes trusted local resource/routing policy into Wingless without treating it as a plane reservation;
- forces `Repairs=0` and `MaxRepairs=0` for plane-driven work so Wingless cannot duplicate plane retry ownership;
- maps plane cancellation into cooperative cancellation of the corresponding Wingless request while leaving authoritative terminal state to plane;
- projects only candidate output with `external_required` acceptance;
- rejects any Wingless acceptance claim or unknown terminal state.

Replay tests cover valid request translation, workspace identity mismatch, stale revision, protection-gate refusal, missing trusted backend policy, resource/routing translation, zero duplicate retry budget, cancellation mapping, candidate-only success, blocked backend outcome, acceptance-claim rejection, unknown-field refusal and unknown-terminal-state refusal.

## Validation state

An isolated adapter compile/replay harness using the reconciled Wingless API shapes passed locally. This establishes the adapter's pure translation contract independently of live-plane state.

The repository's `native.yml` GitHub Actions workflow was also triggered by the adapter commits. Two attempts failed before step 1 on both `windows-latest` and `ubuntu-latest`: each job reported `runner_id=0` with an empty step list. No checkout, Go setup or test command ran. This is classified as hosted-runner infrastructure unavailability, not a source/test failure and not a green CI result.

No live ckb-plane tests were performed by Wingless, and no queue/control mutation was used for reconciliation.

## Reconciled seam classification

`docs/integration-seams.md`, `docs/integration-seams.json` and `.ice/integration-seams.json` now carry the source-backed classifications. The former all-`BLOCKED` inventory is retired.

Notable decisions:

- `worker`: `ADAPTER_REQUIRED`;
- `broker`, `inference`, `resources`, `contextbroker`, inference lifecycle, benchmark, modelhost and telemetry: `WINGLESS_OWNS` within their bounded intelligence/evidence role;
- `toolboundary`: `ADAPTER_REQUIRED` for any write/process operation;
- broker fixture verification as a plane retry/acceptance mechanism: `OBSOLETE`; internal offline research use remains allowed.

## Remaining blockers before live activation

These blockers are intentionally outside Stage 3.5 isolated-readiness:

1. **No plane Wingless worker kind exists.** Frozen `WorkOrder.Validate` permits only `fake` or `codex`. A future live integration must add an explicit plane-side Wingless worker contract; it must not alias Wingless to `codex` or bypass validation.
2. **No live transport is authorized or implemented.** The current adapter is pure translation/replay only.
3. **Write/process tools remain plane-owned.** Wingless's current tool boundary remains bounded/read-oriented.
4. **Hardware qualification remains separate.** Native Windows/NVIDIA/CUDA inference, residency, process GPU attribution, Windows Job Object runtime containment and energy telemetry remain Stage 3 qualification items.
5. **Activation requires a new isolated proof.** That proof must use mock/fixture plane inputs first and must not consume the live queue, change accepted refs, launch CKB, touch Coinbase/broker/credentials/production state or broaden authority.

## Next proof

The next permitted engineering step is an isolated integration proof that exercises the inactive adapter with recorded/mock plane fixtures and a Wingless mock backend. It should prove request translation, cancellation propagation, workspace/revision binding, authorized routing/fallback, backend failure evidence, candidate-only results and absence of duplicate retry ownership. Only after that proof should a proposal be made for the explicit plane-side Wingless worker contract.

READY_FOR_ISOLATED_INTEGRATION_PROOF
