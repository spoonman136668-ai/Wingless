# Integration seams — Stage 3.5 blocked inventory

All plane-facing contracts remain provisional. The stabilized plane source has not been supplied. BLOCKED classifies the comparison, not the existence of internal Wingless functionality. No match to the plane is asserted.

See [reconciliation intake and required inputs](plane-reconciliation-stage.md).

| Seam | Status | Disposition | Plane reference |
| --- | --- | --- | --- |
| worker | BLOCKED | ADAPTER_ONLY | Not supplied |
| broker | BLOCKED | KEEP_INTERNAL | Not supplied |
| inference | BLOCKED | KEEP_INTERNAL | Not supplied |
| resources | BLOCKED | KEEP_INTERNAL | Not supplied |
| contextbroker | BLOCKED | KEEP_INTERNAL | Not supplied |
| toolboundary | BLOCKED | ADAPTER_ONLY | Not supplied |
| broker verification | BLOCKED | KEEP_INTERNAL | Not supplied |
| inference lifecycle | BLOCKED | KEEP_INTERNAL | Not supplied |
| benchmark | BLOCKED | KEEP_INTERNAL | Not supplied |
| modelhost (independent experiment) | BLOCKED | KEEP_INTERNAL | Not supplied |
| stage3 internal telemetry and output contracts | BLOCKED | KEEP_INTERNAL | Not supplied |

## worker

Provisional interface: SubmitWork / Status / Cancel / Result / Health / Capabilities.

Current documented data: Request + Outcome; in-memory IDs; shared deadline.

Required authority boundary: plane owns queue, dispatch and acceptance. This is not a verified plane source finding.

Failure semantics in current inventory: busy/duplicate/full rejected; terminal operator_blocked.

Required reconciliation: map IDs, retention, leases and recovery to authoritative plane contract.

Remaining uncertainty: plane may require durable asynchronous resume. Translation rules and plane replay coverage are blocked until frozen definitions and fixtures are supplied.

## broker

Provisional interface: trusted route Policy and backend allowlist.

Current documented data: Explicit, AllowedBackends, AllowDeep, AllowFallback, MaxRepairs <= 8.

Required authority boundary: caller config authorizes routing; output does not. This is not a verified plane source finding.

Failure semantics in current inventory: no allowed healthy capable backend => blocked.

Required reconciliation: translate signed/validated plane envelope without widening permissions.

Remaining uncertainty: plane authorization schema may differ. Translation rules and plane replay coverage are blocked until frozen definitions and fixtures are supplied.

## inference

Provisional interface: InferenceBackend.

Current documented data: Request carries IDs, role, context byte limit, token cap, absolute deadline, workspace identity, capabilities, resources.

Required authority boundary: inference outputs candidates only. This is not a verified plane source finding.

Failure semantics in current inventory: typed timeout/canceled/HTTP failure; no implicit retry.

Required reconciliation: version result envelope and artifact exchange before integration.

Remaining uncertainty: plane may mandate artifacts rather than text. Translation rules and plane replay coverage are blocked until frozen definitions and fixtures are supplied.

## resources

Provisional interface: Provider.Snapshot + Policy.

Current documented data: nullable metrics; RAM/VRAM/disk floors; CPU ceiling.

Required authority boundary: trusted provider; request policy cannot come from model. This is not a verified plane source finding.

Failure semantics in current inventory: missing required metric denies; provider error blocks.

Required reconciliation: Reconcile reservation ownership; Stage 3 already supplies Windows host collectors, with native proof still outstanding..

Remaining uncertainty: plane may reserve resources or manage processes itself. Translation rules and plane replay coverage are blocked until frozen definitions and fixtures are supplied.

## contextbroker

Provisional interface: Provider.Retrieve.

Current documented data: query + byte budget -> source snippets.

Required authority boundary: ICE advisory; source hashes checked. This is not a verified plane source finding.

Failure semantics in current inventory: stale or over-budget provider fails.

Required reconciliation: Bind ICE root, workspace identity and source revision using the actual plane contract; do not assume a lease exists..

Remaining uncertainty: plane may supply revision-bound source provider. Translation rules and plane replay coverage are blocked until frozen definitions and fixtures are supplied.

## toolboundary

Provisional interface: Execute with caller-owned Policy.

Current documented data: read_file proposal; exact allowed paths; rooted read; byte/time policy.

Required authority boundary: caller grants tool/path; model only proposes. This is not a verified plane source finding.

Failure semantics in current inventory: deny escapes/nonregular/oversize/canceled reads.

Required reconciliation: delegate future writes/processes to stabilized plane executor; do not expand shell access.

Remaining uncertainty: plane may require compiler/edit tools and OS sandbox. Translation rules and plane replay coverage are blocked until frozen definitions and fixtures are supplied.

## broker verification

Provisional interface: Verifier.Verify.

Current documented data: trusted verifier error; at most MaxRepairs+1 attempts; shared deadline.

Required authority boundary: fixture verdict is not plane acceptance. This is not a verified plane source finding.

Failure semantics in current inventory: failed checks escalate only if authorized; otherwise blocked.

Required reconciliation: disable local repairs when plane owns retry budget.

Remaining uncertainty: plane may own all retries. Translation rules and plane replay coverage are blocked until frozen definitions and fixtures are supplied.

## inference lifecycle

Provisional interface: Lifecycle.Move.

Current documented data: STOPPED/STARTING/READY/BUSY/IDLE/DRAINING/FAILED/UNAVAILABLE.

Required authority boundary: no process launch permission implied. This is not a verified plane source finding.

Failure semantics in current inventory: invalid transition rejected.

Required reconciliation: reconcile supervisor; attach transition state to managed backend before real deployment.

Remaining uncertainty: plane may own runtime lifecycle. Translation rules and plane replay coverage are blocked until frozen definitions and fixtures are supplied.

## benchmark

Provisional interface: Replay B/C/D plus future A/E/F/G.

Current documented data: fixture ID/prompt/expected; per-attempt result and route.

Required authority boundary: exact fixture comparison is research evidence only. This is not a verified plane source finding.

Failure semantics in current inventory: unsupported configurations return unavailable.

Required reconciliation: add isolated coding fixtures and authoritative acceptance adapter.

Remaining uncertainty: real coding acceptance more complex than exact text. Translation rules and plane replay coverage are blocked until frozen definitions and fixtures are supplied.

## modelhost (independent experiment)

Provisional interface: none activated; operator-owned Config only.

Current documented data: Pinned executable/model; bounded CPU and experimental GPU configuration; process evidence; hardware qualification separate..

Required authority boundary: no model-proposed execution and no plane authority. This is not a verified plane source finding.

Failure semantics in current inventory: startup/runtime failure stops owned child; no automatic retry.

Required reconciliation: reconcile lifecycle ownership before exposing any plane adapter.

Remaining uncertainty: plane may own inference runtime lifecycle. Translation rules and plane replay coverage are blocked until frozen definitions and fixtures are supplied.

## stage3 internal telemetry and output contracts

Provisional interface: none activated; GPUProvider, ProcessProbe, OutputConstraint remain internal/provisional.

Current documented data: nullable hardware/process samples; bounded output constraint; independent benchmark scores.

Required authority boundary: no resource lease, generic process tool or acceptance grant. This is not a verified plane source finding.

Failure semantics in current inventory: required missing telemetry denies; unsupported constraints/profiles report unavailable.

Required reconciliation: map only after stabilized source; preserve one authority.

Remaining uncertainty: plane may own lifecycle, resource reservations and output schemas. Translation rules and plane replay coverage are blocked until frozen definitions and fixtures are supplied.

