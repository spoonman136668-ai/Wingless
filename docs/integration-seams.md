# Integration seams — provisional

No interface here is authoritative for ckb-plane. Reconcile every seam before activating an adapter.

## worker

- **assumed external interface:** SubmitWork / Status / Cancel / Result / Health / Capabilities
- **reason:** narrow independently testable worker boundary
- **expected plane point:** future worker adapter
- **data contract:** Request + Outcome; in-memory IDs; shared deadline
- **authority boundary:** plane owns queue, dispatch and acceptance
- **failure semantics:** busy/duplicate/full rejected; terminal operator_blocked
- **risk if plane differs:** plane may require durable asynchronous resume
- **reconciliation action:** map IDs, retention, leases and recovery to authoritative plane contract

## broker

- **assumed external interface:** trusted route Policy and backend allowlist
- **reason:** models must not select their own authority
- **expected plane point:** dispatch policy
- **data contract:** Explicit, AllowedBackends, AllowDeep, AllowFallback, MaxRepairs <= 8
- **authority boundary:** caller config authorizes routing; output does not
- **failure semantics:** no allowed healthy capable backend => blocked
- **risk if plane differs:** plane authorization schema may differ
- **reconciliation action:** translate signed/validated plane envelope without widening permissions

## inference

- **assumed external interface:** InferenceBackend
- **reason:** vendor-neutral bounded generation
- **expected plane point:** worker inference operation
- **data contract:** Request carries IDs, role, context byte limit, token cap, absolute deadline, workspace identity, capabilities, resources
- **authority boundary:** inference outputs candidates only
- **failure semantics:** typed timeout/canceled/HTTP failure; no implicit retry
- **risk if plane differs:** plane may mandate artifacts rather than text
- **reconciliation action:** version result envelope and artifact exchange before integration

## resources

- **assumed external interface:** Provider.Snapshot + Policy
- **reason:** fail closed on unknown required floors
- **expected plane point:** pre-inference host budget check
- **data contract:** nullable metrics; RAM/VRAM/disk floors; CPU ceiling
- **authority boundary:** trusted provider; request policy cannot come from model
- **failure semantics:** missing required metric denies; provider error blocks
- **risk if plane differs:** plane may reserve resources or manage processes itself
- **reconciliation action:** reconcile reservation ownership; add context-aware Windows collector

## contextbroker

- **assumed external interface:** Provider.Retrieve
- **reason:** bounded source context
- **expected plane point:** workspace materialization output
- **data contract:** query + byte budget -> source snippets
- **authority boundary:** ICE advisory; source hashes checked
- **failure semantics:** stale or over-budget provider fails
- **risk if plane differs:** plane may supply revision-bound source provider
- **reconciliation action:** bind ICE root and source revision to isolated workspace lease

## toolboundary

- **assumed external interface:** Execute with caller-owned Policy
- **reason:** no model shell authority
- **expected plane point:** isolated tool executor
- **data contract:** read_file proposal; exact allowed paths; rooted read; byte/time policy
- **authority boundary:** caller grants tool/path; model only proposes
- **failure semantics:** deny escapes/nonregular/oversize/canceled reads
- **risk if plane differs:** plane may require compiler/edit tools and OS sandbox
- **reconciliation action:** delegate future writes/processes to stabilized plane executor; do not expand shell access

## broker verification

- **assumed external interface:** Verifier.Verify
- **reason:** bounded deterministic repair
- **expected plane point:** plane acceptance/test evidence
- **data contract:** trusted verifier error; at most MaxRepairs+1 attempts; shared deadline
- **authority boundary:** fixture verdict is not plane acceptance
- **failure semantics:** failed checks escalate only if authorized; otherwise blocked
- **risk if plane differs:** plane may own all retries
- **reconciliation action:** disable local repairs when plane owns retry budget

## inference lifecycle

- **assumed external interface:** Lifecycle.Move
- **reason:** explicit lifecycle transitions
- **expected plane point:** future local model supervisor
- **data contract:** STOPPED/STARTING/READY/BUSY/IDLE/DRAINING/FAILED/UNAVAILABLE
- **authority boundary:** no process launch permission implied
- **failure semantics:** invalid transition rejected
- **risk if plane differs:** plane may own runtime lifecycle
- **reconciliation action:** reconcile supervisor; attach transition state to managed backend before real deployment

## benchmark

- **assumed external interface:** Replay B/C/D plus future A/E/F/G
- **reason:** compare identical deterministic fixtures
- **expected plane point:** offline experiment harness
- **data contract:** fixture ID/prompt/expected; per-attempt result and route
- **authority boundary:** exact fixture comparison is research evidence only
- **failure semantics:** unsupported configurations return unavailable
- **risk if plane differs:** real coding acceptance more complex than exact text
- **reconciliation action:** add isolated coding fixtures and authoritative acceptance adapter

## modelhost — independent experiment

No plane interface is activated. Operator configuration starts one pinned local inference server. Lifecycle ownership must be reconciled against the stabilized plane before integration. No model/tool execution authority is added. Machine-readable details are in integration-seams.json.
