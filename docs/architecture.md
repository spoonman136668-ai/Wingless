# Architecture

`worker.Service` provides a provisional, single-concurrency in-memory boundary. `broker.Runner` checks request bounds, obtains trusted resource metrics, chooses an explicitly allowed healthy/capable backend and invokes it within the shared deadline. A trusted optional verifier can request bounded repair; outputs remain candidates.

`broker.Select` prefers fast inference. Caller-classified hard failure, ambiguity, low confidence or explicit deep request selects deep when authorized. An explicit backend ID wins over automatic routing but never bypasses capability, health or resource checks. Deep-to-fast fallback requires explicit policy. `PlanThenImplement` uses deep planning and fast implementation with the plan as untrusted bounded context. It is a separate caller-selected two-phase operation, not an automatic policy override.

`inference.Mock` and `LocalHTTP` implement the backend contract. Backend class names permit later Codex and local model adapters, but registering a class name does not implement an adapter. `Lifecycle` validates state transitions independently; process management is deferred.

`contextbroker.ICE` retrieves current source ranges through a versioned hash-validated index. `toolboundary` provides a separate read-only authorized tool primitive. Neither is automatically granted to model output. Future model/tool loops must explicitly parse proposals and invoke the tool boundary; there is no generic shell executor.

No live-plane, CKB, broker, credential, deployment, accepted-ref or queue integration exists.
