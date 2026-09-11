# Internal Stage 3 — Windows/GPU candidate

Status: source candidate with real CPU model evidence. **Not fully hardware-qualified.** Windows execution, NVIDIA measurements and CUDA inference are still unproven on a supported host. Plane-facing interfaces remain provisional.

## What changed

- A vendor-neutral GPUProvider and bounded, read-only NVIDIA query collect model/UUID, total/used/free VRAM, GPU/memory utilization, temperature and power where supported. N/A stays null. Required VRAM still fails closed. A trusted absolute nvidia-smi path is mandatory; arbitrary queries and driver mutations are not exposed.
- Owned-process probes provide working-set/RSS, OS-reported peak, lifetime and disk I/O where available. Windows includes aggregate process CPU time. Linux intentionally leaves aggregate CPU time null rather than mislabeling per-thread schedstat. Process identity changes are rejected. This sandbox's mismatched procfs namespace makes owned-process data unavailable; no unrelated PID is sampled.
- Windows amd64 supervision uses a non-inheritable Job Object with KILL_ON_JOB_CLOSE. Assignment failure kills/reaps the owned child. Closing the job cleans assigned descendants; parent exit closes the job handle. Assignment occurs just after process creation: **there is a pre-assignment race window**, so adversarial pre-assignment descendants are not guaranteed contained. Atomic launch-in-job and native descendant/parent-death proof remain necessary. Other Windows architectures refuse supervision. Linux remains direct-child-only. Shutdown is bounded hard termination; graceful server shutdown is not claimed.
- Trusted GPU layers, batch/microbatch, context, threads, device index and VRAM/RAM floors are bounded. GPU startup binds CUDA_VISIBLE_DEVICES to the observed NVIDIA UUID and requests CUDA0 explicitly. CPU fallback is not selected. A compatible trusted CUDA runtime is required. No model output chooses paths, arguments or environment.
- A schema-constraint request seam exists. Current adapters explicitly reject it as unavailable; no untested llama-specific grammar is silently applied.
- SSE regressions cover byte fragmentation including UTF-8 splits, malformed UTF-8, duplicate finish/usage/completion, missing completion, oversized events/health, cancellation and timeout. Inference waits for stream EOF after completion within the request deadline; peers that never close are refused on timeout.
- ICE `relationships` uses go/types and source importing without executing a compiler. Successfully checked active-platform non-test packages yield static call targets and same-package interface implementation relationships. Diagnostics preserve incomplete analysis. Existing test-to-source links remain syntactic, not measured coverage. An initial Linux run resolved 769 relationships with no diagnostics; subsequent source changes require rebuilding/requerying.

## Real 1.5B experiment

[Raw benchmark](experiments/qwen-1.5b-linux-stage3.json): Qwen2.5-Coder-1.5B-Instruct Q4_K_M, 1,117,320,768 bytes, CPU-only llama.cpp b10809, two threads, 1024 context, two repetitions.

Nine operations cover JSON, Go generation, repair, diagnosis, review, small multi-file reasoning, planning, plan-conditioned implementation and a structured read proposal. No generated source or tool proposal is executed.

| Measurement | Passes |
| --- | --- |
| Semantic checks | 15 / 18 |
| Protocol compliance | 2 / 18 |
| Strict correctness | 2 / 18 |

Semantic scoring recognizes at most one complete Markdown fence as a separate diagnostic view. Raw output is preserved; protocol and strict scoring are never fed stripped output. The run shows useful distinctions between content and format reliability, not broad coding proficiency. Caller-owned routing can classify protocol_failure; there is no automatic plane retry ownership.

The v2 report records model/runtime hashes, repository/revision/quantization, bytes, context, threads, GPU layers/batches, host/process telemetry, timing, usage, outcomes and per-task population latency variance. Unknown CPU model, GPU/process peaks, energy and unavailable I/O stay null. GPU telemetry is host/device-wide, not inference-attributed energy. RAM/VRAM-hours and energy/task are not fabricated.

CPU_ONLY and GPU_RESIDENT_SMALL validate configuration and deny incompatible mode selection. GPU residency is a requested experimental profile, not measured certification. FAST_LOCAL, DEEP_LOCAL and HYBRID report unavailable. External endpoints are marked EXTERNAL_UNVERIFIED and cannot assert a verified profile.

## Windows/NVIDIA acceptance

[Hosted run 34578254314](https://github.com/spoonman136668-ai/Wingless/actions/runs/34578254314) failed with zero executed steps for both jobs. Raw job evidence is retained alongside this report. Cause is not exposed by the available connector. No self-hosted CKB runner was restarted or repurposed.

From an independent Windows Wingless checkout with Go 1.25.5:

```powershell
$ErrorActionPreference = 'Stop'
.\scripts\Test-Native.ps1
# On a Windows NVIDIA host, additionally require real driver collection:
.\scripts\Test-Native.ps1 -NvidiaSMI 'C:\Windows\System32\nvidia-smi.exe'
# Supply your actual trusted CUDA runtime and modest GGUF paths:
.\scripts\New-LocalModelConfig.ps1 -Executable C:\models\llama-server.exe -Model C:\models\coder.gguf -Output .\gpu-model.json -GPULayers 128 -NvidiaSMI 'C:\Windows\System32\nvidia-smi.exe' -MinVRAM 2147483648
go run ./cmd/wingless local-benchmark --config .\gpu-model.json --profile GPU_RESIDENT_SMALL --repeats 2
if ($LASTEXITCODE -ne 0) { throw 'GPU experiment failed' }
```

Use the actual driver-installed path; the example is not a discovery assumption. Scripts create local evidence/configuration only. Computing a local file hash records a pin, not publisher authentication. Runtime directories and dependencies must be trusted and stable.

## Exit disposition

Linux race tests and Windows cross-build are source checks. No native Windows/GPU pass is claimed. GPU driver tests, CUDA inference/residency, Job Object descendant/parent-death execution and reliable owned-process sampling must be verified on independent supported hardware. Peak sampling, energy integration, constrained generation support and cross-package/type-resolved test coverage remain open. No production readiness or Codex parity claim.

References: [NVIDIA SMI](https://docs.nvidia.com/deploy/nvidia-smi/index.html), [Windows Job Objects](https://learn.microsoft.com/en-us/windows/win32/api/winnt/ns-winnt-jobobject_extended_limit_information).
