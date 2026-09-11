# Known gaps — Stage 3 candidate

See [Stage 3 source changes, measurements and acceptance](internal-stage-3.md).

Native Windows/NVIDIA/CUDA and containment execution remain hardware-blocked. Source implementations and CPU benchmark evidence must not be represented as native qualification. Job assignment has a documented post-start race window. Process telemetry rejects incompatible procfs namespaces; process CPU on Linux and inference-attributed GPU memory remain unavailable. Peak/time-integrated memory and energy metrics remain open. Constrained generation is an explicitly unavailable capability seam. ICE static call targets and same-package implementation checks are stronger than syntax hints, but test links and cross-package implementation coverage remain incomplete.

No durable queue, leases, plane retry/acceptance ownership, generic shell/compiler/write authority, live-plane integration or production service installation is implemented. All plane-facing seams are provisional.

The 1.5B CPU experiment scored 15/18 semantic, 2/18 protocol and 2/18 strict. This is limited fixture evidence, not production coding readiness or Codex parity. Earlier Stage 2 measurements remain preserved.
