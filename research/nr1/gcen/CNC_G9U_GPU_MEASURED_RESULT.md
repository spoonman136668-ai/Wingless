# G9U — Authenticated original-expert GPU measurements (negative speed result)

**Engineering: PASS. Scientific: NO_GO for a claimed speedup of the current tiled INT2 operator. Deployment: NOT AUTHORIZED.**

This is the immutable development result of one original Qwen expert benchmark, not model inference or full-checkpoint fit. The source experiment is SHA `e5ee6de7ef7bfb09947753d631358815c99e74ce`, run [#42092133021](https://github.com/spoonman136668-ai/Wingless/actions/runs/42092133021), attempt **2**, archived artifact [#13456072646](https://github.com/spoonman136668-ai/Wingless/actions/runs/42092133021/artifacts/13456072646). ZIP SHA-256 `c954a6223f25922a0bced0933060c4396e03ece3b795d6298d112f9e43582a01`. The machine-readable extracted, provenance-bound summary is [the G9U development receipt](receipts/cnc01q_g9u_original_gpu_development.json).

## Frozen comparison and custody

- Original pinned BF16 `Qwen/Qwen3-30B-A3B` revision `ad44e777bcd18fa416d9da3bd8f70d33ebb85d39`; shard SHA-256 `454e77b346a61bfb201d54df60e15158838cf930617ee135113556204f2802b5`. In-process shard verification and original expert extraction were confirmed in the archived receipt.
- One layer-0 expert-0 down projection, `2048 × 768`; batch 8, frozen synthetic Gaussian inputs (not authentic activations), seeds **751/811**, 3 warmups and 7 CUDA-event timings per arm per seed.
- Matched FP32 dense CUDA GEMM vs existing companded INT2 GPU **tiled decode plus FP32 GEMM**. This is **not** a fused packed kernel.

| Frozen seed | Dense FP32 median GPU ms | INT2 tiled median GPU ms | Packed speedup over dense | Relative squared output error |
|---|---:|---:|---:|---:|
| 751 | 1.939648 | 6.831328 | 0.2839× | 0.119193 |
| 811 | 1.865088 | 1.238432 | 1.5060× | 0.118092 |

Mean seed median: dense **1.902368 ms**; packed **4.034880 ms**. Packed was **2.12098× as slow** on this small matched sample; the per-seed disagreement prohibits a robust speed conclusion. Independent packed-operator CPU/GPU parity passed with relative squared differences around `1e-13`, but it does **not** establish parity to the original BF16 output. Mean compression relative squared output error: **0.1186425**.

Original BF16 expert stored weight bytes: **3,145,728**. Packed INT2 stored bytes including codec metadata: **405,504** (**7.7576×** smaller). Measured **joint-residency CUDA allocator peak**: **24,429,056 bytes**, not packed-only footprint, device-wide usage, full-model fit, or a 12-GiB launch result. Runtime scale conversion and temporary allocations must be accounted for separately.

## Decision and next falsifiable question

Retain the packed codec as a **storage feasibility baseline only**. **Do not promote the current tiled operator as a speed optimization** or infer quality on real prompts. Before investing in a new GPU kernel, first obtain separately authorized, resource-bounded phase timings of the existing operator (GPU decode/indexing, scaling, GEMM and allocations), with the same authenticated expert, frozen seeds and matched dense reference. Pre-register a phase budget and repeat the experiment independently; compare both seeds and include timing-overhead negative controls. If the variability persists or the compression fidelity remains inadequate on authentic activations, stop this speedup track instead of multiplying variants.

No original activations, protected heldout, full-model forward, scientific acceptance, or deployment were authorized or established by this run. The original shard is kept isolated on `WINGLESS-UP-B` and is **not** committed to this repository.
