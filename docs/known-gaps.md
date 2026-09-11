# Known gaps and nonclaims

## Plane-dependent: deliberately deferred

No live-plane integration, durable work queue, leases, crash recovery, dispatcher ownership, write/process tool authority, plane retry ownership or acceptance ownership. The service remains a provisional in-memory test boundary. No model output becomes acceptance. Reconcile contracts only after stabilized plane source is supplied.

## Wingless-internal

- Native Windows execution remains unproven: the hosted CI attempt failed before any test step. Windows code can be cross-compiled but that is not native acceptance.
- Host RAM/disk/CPU collection exists for Windows/Linux. GPU/VRAM, energy, container limits, per-process memory, peak sampling and resource reservations are unavailable.
- A pinned local llama-server can be started/stopped for bounded CPU-only experiments. No arbitrary command interface, model download manager, automatic restart, descendant containment, supervisor-crash recovery, GPU lifecycle, pressure eviction or graceful unload exists. Trusted stable runtime directories are required.
- SSE and nonstreaming LocalHTTP work with bounded payloads and cancellation. Authentication, remote endpoints, tool calls, logprobs and an exhaustive server compatibility matrix are not supported. Missing token usage stays null; byte bounds still apply.
- A real 0.5B model was tested on two strict microfixtures. Both failed formatting requirements. Transport/supervision success does not establish coding quality or Codex parity.
- No autonomous edit/build/test loop. Read-only tools remain separately authorized; no model proposal gets shell or write access.
- ICE maps are syntax hints, not type-resolved call graphs or measured coverage. Explicit interface assertions are compiler checked; other implementations can be absent. Validation rescans source. Generated indexes are advisory.
- Worker storage is bounded and in-memory. Trusted Go implementations must cooperate with deadlines. Rooted regular-file reads cannot forcibly interrupt a hung filesystem.

Future internal research: stronger output contracts, broader coding fixtures, measured resource peaks, runtime dependency integrity, GPU collectors, bounded expert caching, speculative drafting, KV budgets and energy-aware scheduling. These are not implemented or accepted features.

See [stage 2 evidence](internal-stage-2.md).
