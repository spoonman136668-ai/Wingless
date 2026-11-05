# G9V — One original-expert GPU phase experiment: STOP

**Engineering PASS · Scientific NO_GO for the current INT2 GPU speedup track · Deployment NOT AUTHORIZED.**

The one-shot diagnostic used original BF16 Qwen/Qwen3-30B-A3B weights from revision `ad44e777bcd18fa416d9da3bd8f70d33ebb85d39`, original shard SHA256 `454e77b346a61bfb201d54df60e15158838cf930617ee135113556204f2802b5`, in-process custody verification and an unchanged G9U packed-tile operator. Original model activations and full-model forward were not run.

- Experiment exact source SHA: `d2a458f9b89c5b7511ad08bb73195220dd68b847`, [Actions run 42122916246](https://github.com/spoonman136668-ai/Wingless/actions/runs/42122916246) attempt 1; Linux, Windows, G9T and the isolated original-expert G9V GPU job all PASS.
- Evidence: [artifact 13459562694](https://github.com/spoonman136668-ai/Wingless/actions/runs/42122916246/artifacts/13459562694), ZIP SHA-256 `212c0c3576f11dff283f9a2ea7112a6bd5d40f1d6bf8723e603390083a64476a`. This permanent [machine-readable result](receipts/cnc01q_g9v_original_gpu_phase_development.json) includes per-seed samples and phase medians.
- Frozen comparison: 2048×768 original BF16 expert-0 down projection, batch 8 synthetic Gaussian activations, seeds 751/811, 3 warmups, 7 CUDA-event timing repeats, matched dense FP32 vs existing companded INT2 decode-plus-GEMM, event-only negative control, independent CPU packed parity.

| Seed | Dense FP32 median | INT2 packed median | Packed slowdown | Unpack+decode share | Relative squared output error |
|---|---:|---:|---:|---:|---:|
| 751 | 0.111712 ms | 0.744512 ms | 6.66× | 55.66% | 0.119193 |
| 811 | 0.188160 ms | 0.738496 ms | 3.92× | 57.31% | 0.118092 |

Mean seed medians: dense **0.149936 ms**, INT2 **0.741504 ms** (**4.95× slower** in this run). The per-seed packed medians are similar here, but dense timing and prior G9U absolute latencies differ materially across runs. This is not evidence of a robust speed gain.

Median phase durations (ms), from intrusive instrumented trace: seed 751 allocation 0.1630, unpack 0.3960, decode 0.1511, GEMM 0.1404, scale/accumulate 0.1324; seed 811 allocation 0.1482, unpack 0.3999, decode 0.1531, GEMM 0.1345, scale/accumulate 0.1292. The trace overhead ratios to uninstrumented packed medians were **1.565×** and **1.539×**, above the preregistered 1.50× attribution limit. Event-only control medians were 0.0097 and 0.0020 ms. Thus the phase fractions identify work but cannot be treated as reliable uninstrumented hot-path allocation.

## Frozen decision

Primary preregistered screen: synthetic-input normalized squared output error must not exceed **0.02** for either seed. Both measured about **0.119**, a decisive STOP. Secondary checks: combined unpack+decode share must reach 60% on both seeds; actual fractions were 55.66% and 57.31%. Trace inflation must not exceed 1.5×; both exceeded it. This is not evidence of acceptable model inference, despite verified packed CPU/GPU parity (~1e-13). INT2 storage compression around 7.76× remains a storage research result only.

**Recommendation:** Stop additional kernel work on this specific companded INT2 expert implementation; retain the negative evidence and, if memory compression becomes a higher priority later, consider fidelity-first mixed precision and authenticated real-activation tests as a **new separately authorized study**, not an automatic next gate. No new checkpoint download, FORGE changes, live deployment or protected-heldout access occurred.
