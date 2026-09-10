# Provisional worker protocol

`worker.Boundary` exposes SubmitWork, Status, Cancel, Result, Health and Capabilities. The Go API is the test boundary; no HTTP/RPC listener is exposed. IDs are unique per service lifetime. Exactly one job runs at once, at most 256 completed records are retained, and records are lost on process exit. No queue consumption or acceptance mutation occurs.

Request fields: request_id, parent_work_id, role, context, max_context_bytes (1..1 MiB), max_output_tokens (1..32768), absolute deadline (future, at most 15 minutes), workspace identity, required capabilities, resource floors. Workspace identity in inference is metadata, not permission to edit. Service-owned route policy is configured independently from model content.

Result status is result_ready, verified_candidate or operator_blocked. Acceptance is always external_required. Events carry route, rejected candidates, attempt, backend/model, text, available usage/resource metrics, latency, termination class and optional verification/error. Fixture success does not accept a plane work order. A canceled parent cancels this provisional worker; plane integration must supply a properly owned operation context, not a short-lived HTTP request context.

Cancel requests cooperative backend cancellation. No generic forced process kill is implemented. Callers must retain the parent context through completion. Duplicate IDs and simultaneous jobs are rejected; no implicit retry on service restart.
