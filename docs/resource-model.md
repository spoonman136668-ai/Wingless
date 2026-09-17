# Resource and telemetry model

Nullable host metrics describe RAM/VRAM total/free, CPU/GPU utilization, the legacy residency flag, disk free/read bytes and energy. Unknown is null, never a fabricated zero. Host collectors use native Windows APIs and Linux /proc plus statfs. There is no GPU driver integration. The demo uses unknown metrics with no required floors. A required floor with missing telemetry fails closed. Provider errors block routing. CPU ceiling must be 0 (disabled) or (0,100].

Host/process telemetry is not model attribution. Process RSS is not labeled model RAM, GPU allocation is not labeled model residency, and a model file size is not treated as an active working set. Existing `resources.Policy` admission semantics remain unchanged.

## Model-attributable telemetry

CR-1A adds the versioned `wingless.model-resource-telemetry.v1` evidence shape. It is schema readiness only; it does not implement paging, cache eviction, MoE conversion, model surgery, storage scheduling, NVMe streaming, DirectStorage, or any other neural-residency mechanism.

The optional model telemetry can represent:

- capacity: `logical_model_bytes`, `resident_model_bytes`, `active_model_bytes`;
- memory tiers: `model_vram_bytes`, `model_host_ram_bytes`, `model_storage_bytes`;
- movement: `host_to_device_bytes`, `device_to_host_bytes`, `storage_to_host_bytes`, `storage_to_device_bytes`;
- working set: `working_set_bytes_current`, `working_set_bytes_peak`, plus an explicitly marked active-parameter estimate;
- future cache/paging counters: `model_cache_hits`, `model_cache_misses`, `model_page_faults`, `model_pages_loaded`, `model_bytes_loaded`.

Every numeric field is nullable. Null means unavailable. Zero means a measured/reported zero. No CR-1A code substitutes zero for unknown. The top-level model telemetry is optional on general inference telemetry to preserve existing consumers; where CR-1A benchmark/session evidence exposes the field directly, unavailable telemetry is serialized as null.

Telemetry groups may carry typed provenance identifying direct measurement, backend report, or runtime report and the provider identity. An active parameter count uses a dedicated estimate type that serializes `estimated=true` and requires a textual basis; token count and file size are not accepted as implicit parameter-activity evidence.

The conceptual future tiers are:

```text
HOT   -> device memory / VRAM
WARM  -> host/system RAM
COLD  -> storage-backed model state
```

These labels describe future accounting topology only. CR-1A does not assign pages, move weights, reserve tiers, or alter admission decisions.

This remains a preflight check, not a resource reservation or continuous scheduler. Concurrent external applications can change resource availability. The resource governor performs no pagefile, model-loading, cache-eviction or power-control operations; modelhost separately owns its bounded local model process. Reservations, pressure monitoring, GPU telemetry and stronger supervision remain needed for production use.

Events record backend/model, route/reason, repair iteration, usage if supplied, measured total latency, result and resource snapshot. Streaming measures first-content and content-generation intervals. The real benchmark samples host resources before/after. Null extension fields preserve unavailable seams. Aggregate heavy-token cost by event route tier and successful deterministic fixture; do not treat mock usage as real inference performance.

Collector references: [GlobalMemoryStatusEx](https://learn.microsoft.com/en-us/windows/win32/api/sysinfoapi/nf-sysinfoapi-globalmemorystatusex), [GetSystemTimes](https://learn.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-getsystemtimes), [GetDiskFreeSpaceExW](https://learn.microsoft.com/en-us/windows/win32/api/fileapi/nf-fileapi-getdiskfreespaceexw). Windows CPU counters cover the primary processor group on hosts with more than 64 logical processors. CPU is null on the first sample or if no counter interval elapsed.
