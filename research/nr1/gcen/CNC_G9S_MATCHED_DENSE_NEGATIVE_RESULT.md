# G9S — Authentic-expert matched dense-versus-packed CPU qualification

**Decision: NO_GO for promoting the G9R Numba packed-INT2 CPU operator as a speed optimization.**

The comparison used authentic BF16 source weights from original pinned Qwen
`Qwen/Qwen3-30B-A3B` revision `ad44e777bcd18fa416d9da3bd8f70d33ebb85d39`,
but **synthetic Gaussian input vectors**, not original-model activations.

## Exact-source receipt

- Branch: `research/cnc01q-direct-expert-development-20261103` in `spoonman136668-ai/Wingless`.
- Qualified experiment SHA: `b2956f2aa35fa04ca04079d963326990420fe071`.
- Native Windows / Ubuntu + isolated original-expert job: **PASS**,
  GitHub Actions run `41927892073`.
- Immutable uploaded artifact ID `13370003471`,
  archive SHA-256 `80e6f7b332c0c2cf9f6faa9655a92a3419579bea0b2127c9b4d1bf4f89f7c95b`.
- Original shard bytes verified: `3,999,417,504`.
  Original shard SHA-256:
  `454e77b346a61bfb201d54df60e15158838cf930617ee135113556204f2802b5`.
- One full original BF16 expert weight tensor, `2048 × 768`;
  two frozen synthetic seeds `751` and `811`, batch 8, seven
  timed repeats per operator per seed. CPU-only; two Numba threads,
  one OpenBLAS thread. No model forward or GPU measurements.

## Recorded median times (seconds)

| Frozen seed | Dense FP32 | Packed NumPy INT2 | Packed Numba INT2 | Numba/dense speedup |
|---|---:|---:|---:|---:|
| 751 | 0.001211844 | 0.005113191 | 0.005170518 | 0.23438× |
| 811 | 0.001172941 | 0.005114964 | 0.005142315 | 0.22810× |

Mean of per-seed medians: dense FP32 **1.1924 ms**,
packed NumPy INT2 **5.1141 ms**, packed Numba INT2 **5.1564 ms**.
Numba was about **4.32× slower than dense** in this matched run,
and slightly slower than the existing packed NumPy operator.
JIT compilation and first execution cost **1.3225 s**, excluded
from the warm timings.

## Fidelity, storage and interpretation

- Relative squared output error of packed Numba INT2 against dense
  original-BF16-weight / FP32-compute reference:
  `0.1181868666` and `0.1197556383`; mean **0.1189712525**.
- Original expert BF16 weight bytes: `3,145,728`.
  Packed INT2 weight bytes, accounting for existing packed format metadata:
  `405,528` (**7.76×** reduction). Additional transient FP32
  runtime scales: `24,576` bytes, so stored compressed bytes are not
  the whole runtime-memory budget.
- Max absolute parity difference versus the existing packed NumPy
  operator: `1.55e-6` and `8.94e-7`.
- Prior G9R receipt from run `41926494531` reported approximately
  1.805× improvement over packed NumPy; **that improvement did not
  reproduce** in this run under the matched dense-reference benchmark.
  No claim of stable cross-run latency acceleration is warranted.
- Both seeds and all three operators used the same source expert and
  matched synthetic inputs. These results apply only to this tested
  CPU implementation and environment. They do **not** establish model
  task quality, GPU speed, context behavior, or end-to-end inference.

## Next independent direction

Do **not** multiply variants of the same weight-only row selection or
promote this CPU kernel. Keep the negative result and pursue a bounded
GPU-packed execution comparison only after independently verifying
available GPU resources, runtime/ownership and execution authority.
Separately, genuine activation-aware quality is blocked until the
existing G9M original-forward capture custody requirements are met.

`scientific_acceptance=NOT_GRANTED`. No model launch, deployment,
protected-main change, or promotion is implied by this receipt.
