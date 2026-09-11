# Plane reconciliation — Stage 3.5

State: **ISOLATED_INTEGRATION_PROOF_ACCEPTED**

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

The exact source was re-qualified on authoritative Windows PowerShell and sealed as `acceptance/frozen-contract-closeout-20260911.json` in ckb-plane. The closeout established manifest verification, exact live-source equality, focused post-DT regressions PASS, full `internal/plane` suite PASS, frozen `cmd/ckb-plane` build PASS, authority=`plane`, controller paused with fresh heartbeat, accepted-ref mutation disabled, and original GitHub runner stopped.

Historical DN/DO/DP/DT evidence remains supplemental provenance; exact source identity is now established directly by the frozen-source Windows closeout.

## Reconciled ownership

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

The frozen plane worker contract is synchronous `Worker.Run(ctx, workspace, WorkOrder, stdout, stderr)`. Wingless's `worker.Boundary` (`SubmitWork/Status/Cancel/Result`) is therefore internal only; it is not projected as plane lifecycle.

## Inactive adapter

`integration/ckbplane` is a pure, non-default translation/replay layer with no network listener, controller transport, queue write, worker launch, accepted-ref access, generic process authority, or live-plane mutation.

It pins the frozen plane source identity, consumes verified work-order/workspace identity, requires exact baseline agreement and plane protection gates, rejects stale identity and unknown fields, passes only trusted local routing/resource policy, forces `Repairs=0` and `MaxRepairs=0`, maps cancellation cooperatively, emits candidate-only output with `external_required`, and rejects any Wingless acceptance claim or unknown terminal state.

## Isolated integration proof — accepted

Authoritative Windows evidence is sealed in:

`acceptance/stage3-5-isolated-integration-20260911.json`

Accepted Wingless source commit:

`806740246fd4af5032b70f9025c43ec667fe8add`

Frozen plane source commit:

`539d56fec273c8851aea131cf4d31425a62f7250`

Windows toolchain:

`go version go1.25.5 windows/amd64`

Accepted proof results:

- focused adapter/proof tests: PASS;
- reconciled Wingless stack packages: PASS;
- full Wingless repository suite: PASS;
- Wingless build: PASS;
- candidate acceptance remains `external_required`;
- Wingless work-order retry budget remains `0`;
- live plane transport used: false;
- live plane queue mutated: false;
- accepted refs mutated: false;
- CKB runtime launched: false;
- Coinbase/broker used: false;
- credentials read: false;
- production state touched: false.

The accepted isolated proof exercised the actual adapter/broker/internal-worker stack with recorded/mock plane input and covered candidate success, authorized deep-to-fast fallback, backend failure with one attempt/no duplicate retry, and cancellation propagation.

GitHub hosted native CI remains separately classified as infrastructure-unavailable-before-step-1 because jobs were returned with `runner_id=0` and no executed steps. That infrastructure state does not replace or negate the authoritative Windows proof.

## Live activation blockers

The isolated proof does **not** authorize live activation. Remaining blockers are deliberate:

1. Frozen plane `WorkOrder.Validate` only accepts worker kinds `fake` and `codex`; no explicit `wingless` worker kind exists.
2. No live plane↔Wingless transport is implemented or authorized.
3. Plane must remain sole owner of queue, retry, cancellation terminal state, workspace lifetime, acceptance and evidence.
4. Wingless must not acquire generic shell/write/process authority; write/process execution remains behind the plane contract.
5. Hardware qualification remains separate: native NVIDIA/CUDA inference, model residency, process GPU attribution, Windows Job Object runtime containment and energy telemetry.
6. Any future activation requires a new isolated Windows proof against the explicit plane-side Wingless worker contract before touching a live queue.

## Next engineering stage

The next permitted step is to design and test an **explicit plane-side `wingless` worker contract** without activating it.

The contract must:

- add `wingless` as an explicit validated worker kind rather than aliasing `codex`;
- preserve the existing synchronous plane `Worker.Run` ownership model;
- pass only plane-owned immutable work-order/workspace identity plus bounded intelligence policy;
- keep plane-owned retry budget, cancellation terminal state, workspace lifecycle, acceptance and evidence unchanged;
- keep Wingless output candidate-only;
- fail closed on unsupported fields, stale identity, unknown state and unavailable backend;
- have no implicit live transport or default registration;
- include replay/unit coverage proving zero authority widening and zero duplicate retry ownership.

**ISOLATED_INTEGRATION_PROOF_ACCEPTED**
