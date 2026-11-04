# CNC-02E — Three-arm original-weight comparison (frozen development)

This is an **extension of the previous CNC-02C numerical development experiment**,
not evidence of original-model training or scientific acceptance. The old
CNC-02C artifacts and independent inference PR #443 remain unchanged.

## Frozen source and cohorts
- Qwen/Qwen3-30B-A3B revision `ad44e777bcd18fa416d9da3bd8f70d33ebb85d39`.
- Reuse single G31 SHA-256-verified original BF16 source shard, G32 first
  complete expert layer, expert IDs 0 through 3 only.
- Original BF16 128×256 windows under G32 source seeds: 11, 23, 37, 53, 71.
- Synthetic **not original Qwen** Gaussian FIT seed 6011 and TRANSFER seed
  17201; 32 inputs each; no TRANSFER updates, labels or tuning.
- Frozen INT4 block64 FP32 scales, all arms initialized from the **same**
  prequantized/dequantized starting weights.
- Frozen exactly two training steps for all arms; direct packed code edits,
  FP32 SGD and straight-through-estimator-style quantization-aware training.
- FP32/QAT learning rate: 0.02 divided by mean input-squared norm.
  No post-result hyperparameter adjustment or selection on TRANSFER.
- 4 original expert windows × 5 original source seeds = exactly 20 receipts.

## Interpretation and measured constraints
- Equal step counts do **not** imply equal gradient FLOPs, parameter updates,
  optimizer cost, compute time, training memory or convergence opportunity.
- FP32 and QAT hold persistent FP32 master weights. Direct packed-code learner
  retains codes and scales only; temporary activations/gradients also consume
  runtime memory. Unmeasured transient peak RSS and VRAM remain UNVERIFIED.
- Report FIT and independent TRANSFER normalized errors for frozen INT4,
  direct-code, FP32 SGD and QAT; account persistent weight bytes separately.
- No original Qwen forward pass, original activation capture, task-quality
  validation, live deployment, external APIs or scientific promotion.

Execution: one-shot on isolated CNC02 branch after mandatory exact-head Linux
and Windows native regression, commit message marker `[cnc02e-real]`, original
SHA-256 verified source-custody path, and fail-closed artifact on any failure.
