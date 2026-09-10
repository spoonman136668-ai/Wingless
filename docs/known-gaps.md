# Known gaps and nonclaims

- No live-plane integration, durable work queue, dispatcher leases, crash recovery, production service, installer or WebUI transport. Plane contracts await reconciliation.
- No real autonomous edit/build/test agent yet. The read-only tool is tested independently; model proposals are not automatically executed. Compiler/process tools and artifact application require a later sandboxed executor.
- No real model benchmark or Codex parity evidence. CLI demo/benchmark use fixtures. B/C/D policy replays work; A/E/F/G explicitly report unavailable. E needs a fixture-integrated context pipeline, F an expert cache, G speculation.
- No real host collector, resource reservation, GPU scheduler, model supervisor or weight downloads. Lifecycle is a transition model, not an attached process manager.
- LocalHTTP is a nonstreaming compatibility subset. It trusts the configured localhost server to honor max_tokens when token usage is absent; a separate byte bound still applies. Model IDs must match exactly.
- ICE test/call links are syntax hints, not type-resolved call graphs or proof of coverage. Explicit implementation assertions are compiler checked; unasserted implementations may be omitted. Source hash validation costs a repository scan; optimize only with safe invalidation evidence.
- Workspace roots must be caller-trusted and stable. Rooted reads prevent path traversal; regular file IO is checked for timeout before/after, not forcibly interruptible on a hung filesystem. No OS sandbox or adversarial mutable filesystem guarantee is claimed.
- Worker storage is in-memory and bounded; caller contexts must outlive work. Trusted plugin implementations must cooperate with deadlines. A malicious arbitrary Go plugin cannot be contained by this API.
- Telemetry captures total latency and available snapshots; optional fine-grained metrics remain null. Timing values vary between runs; routing and fixed-input serialization are deterministic.

Future seams: bounded expert cache and residency policy, asynchronous NVMe reads/prefetch, speculative drafting/verification, KV budgeting, process sleep/eviction, energy-aware scheduling. These are research directions, not implemented features.
