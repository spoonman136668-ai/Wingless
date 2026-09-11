# Wingless internal worker protocol

`worker.Boundary` exposes `SubmitWork`, `Status`, `Cancel`, `Result`, `Health` and `Capabilities`. This is an **internal Wingless service boundary**, not the ckb-plane worker protocol and not a durable queue. The Go API has no HTTP/RPC listener. IDs are unique per service lifetime; exactly one job runs at once, at most 256 completed records are retained, and records are lost on process exit. It never consumes the plane queue or mutates plane acceptance.

For ckb-plane reconciliation, `integration/ckbplane` is the inactive adapter boundary. The frozen plane contract is `spoonman136668-ai/ckb-plane@539d56fec273c8851aea131cf4d31425a62f7250`, where the authoritative worker shape is synchronous `Worker.Run(ctx, workspace, WorkOrder, stdout, stderr)`. The Wingless in-memory job lifecycle must therefore never be projected back as plane work-order lifecycle.

A plane-driven replay is normalized into an adapter envelope containing the verified work-order subset, plane-owned workspace binding, bounded context and trusted Wingless inference policy. Translation requires exact work-order ID and baseline-SHA agreement between the order and workspace. Missing accepted-ref/clean-baseline protection, stale workspace revision, unknown replay fields or invalid bounded inference fields fail closed.

The translated `inference.Request` contains request ID, parent work ID, role, context, context-byte cap, output-token cap, absolute deadline, plane-owned workspace path, capabilities and trusted Wingless resource policy. Workspace identity is metadata/binding, not permission to create or mutate workspaces. The adapter deterministically derives a Wingless request ID from plane work-order ID plus attempt.

For plane-driven work, broker routing remains Wingless-owned but the adapter forces `Repairs=0` and `MaxRepairs=0`; ckb-plane owns retry, retry limits and reconciliation. Trusted adapter configuration may authorize fast/deep routing or fallback, but model output cannot widen that policy.

Result status remains `result_ready`, `verified_candidate` or `operator_blocked`. Acceptance remains `external_required`. `integration/ckbplane.CandidateFromOutcome` rejects any Wingless outcome that attempts to claim acceptance and rejects unknown terminal states. Fixture verification remains internal research evidence, not plane acceptance.

Cancellation maps the plane work-order attempt to the deterministic Wingless request ID and invokes cooperative Wingless cancellation. The caller/plane still owns the authoritative work-order terminal state and any retry decision.

Live activation remains forbidden. The frozen plane currently permits only `fake` and `codex` values for `worker.preferred`; no `wingless` worker kind or live transport exists yet. Adding one requires an explicit plane-side contract change and separate isolated integration proof before any live queue use.
