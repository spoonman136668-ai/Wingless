# CNC-02G — Frozen authentic weight-window consensus evaluation

Development-only test, NOT scientific admission, model inference, activation or deployment.

Source custody: pinned Qwen/Qwen3-30B-A3B revision
`ad44e777bcd18fa416d9da3bd8f70d33ebb85d39`, original BF16
`model-00001-of-00016.safetensors` authenticated under the existing G31
SHA-256/header/index verification. No complete model loading. Retain G32
source-window seeds 11,23,37,53,71, first four original experts from the
first complete layer and 128x256 weight windows (20 comparisons).

New, disjoint synthetic input cohorts (not Qwen model activations):
FIT-A seed **26061**, FIT-B seed **26087**, TRANSFER seed **26111**;
32 samples per cohort. No TRANSFER use for candidate generation, rollback
decisions or hyperparameter choices. Freeze exactly two compressed
coordinate proposals selected on FIT-A and accepted only if each proposal
reduces FIT-B squared output error. Rejecting a proposal restores original
packed codes. Frozen per-block INT4 scales never update.

Controls: frozen INT4 without updates; FP32 SGD starting from the same
quantized reconstruction, two steps on FIT-A at the existing frozen learning
rate. Equal step count is not equal parameters, compute, optimizer cost or
transient memory. Record cohort-specific normalized squared errors,
accepted/rejected proposals and persistent representation bytes. Do not
claim speed, peak GPU memory, model task accuracy or genuine long-range
generalization.

The evaluation is opt-in only after exact-head Linux/Windows regression,
with one source-staged original shard and the existing fail-closed
custody protocol. The trigger `[cnc02g-real]` applies only on the isolated
`research/cnc02-compressed-expert-learning-20261104` branch. The
independent compressed-inference PR #443 is untouched.

If transfer error is not better than frozen INT4 and the FP32 comparator,
retain a NO_GO. Do not use this cohort to tune new thresholds after viewing
results. A separate authentic model-forward and resource test would be needed
for scientific or practical inference/training claims.
