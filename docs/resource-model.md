# Resource and telemetry model

Nullable metrics describe RAM/VRAM total/free, CPU/GPU utilization, residency, disk free/read bytes and energy. Unknown is null, never a fabricated zero. No host collector or GPU driver integration is installed. The demo uses unknown metrics with no required floors. A required floor with missing telemetry fails closed. Provider errors block routing. CPU ceiling must be 0 (disabled) or (0,100].

This is a preflight check, not a resource reservation or continuous scheduler. Concurrent external applications can change resource availability. No pagefile, model loading, cache eviction or power control occurs. Windows telemetry, reservations, pressure monitoring and process supervision must precede live model use.

Events record backend/model, route/reason, repair iteration, usage if supplied, measured total latency, result and resource snapshot. TTFT, generation-only time, load time, before/peak/after samples, VRAM-hours/RAM-hours and energy per task are not measured yet. Null extension fields preserve those seams. Aggregate heavy-token cost by event route tier and successful deterministic fixture; do not treat mock usage as real inference performance.
