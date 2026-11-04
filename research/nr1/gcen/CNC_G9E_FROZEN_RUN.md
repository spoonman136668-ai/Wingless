# CNC-G9E — frozen development-only original-weight INT2 fidelity

Source: Qwen/Qwen3-30B-A3B at
`ad44e777bcd18fa416d9da3bd8f70d33ebb85d39`.
Reuse G8/G31 source authentication and G32's first complete 16 distinct
original `down_proj.weight` projections in original shard
`model-00001-of-00016.safetensors`. Frozen five development window seeds
11, 23, 37, 53, 71 and 4/8/16 expert prefixes, 128×256 windows. The only
test inputs are independently seeded synthetic Gaussian batch-32 activation
probes: **NOT** original Qwen activations or task-quality evidence.

Compare original BF16 matrix multiplication with G8 block-64 FP32-scale
packed INT4 and G9C fixed-codebook INT2 block-128 and block-256 FP16-scale
direct CPU operators. Report relative output- and weight-squared errors,
per-expert error, packed codes/scales/metadata, 4/8/16 comparisons,
physical staged file count, clean temporary-source deletion and peak RSS.

This is **not matched storage** with INT4: evaluate error and storage together.
Do not retune INT2 codebook, scales, blocks, experts, seeds, windows, or
thresholds after observing authentic-source results. Prior G8 SVD negatives
remain unchanged. No GPU, full-model inference, original activation capture,
production use, sealed test cohorts, or scientific acceptance.

Existing Workflow `.github/workflows/native.yml` allows a *one-shot*
`[cnc-g9e-real]` original-weight development job only on the isolated
PR #443 development branch, gated by both exact-head Linux/Windows native
qualification jobs. Preserve FAIL and artifacts on source/network failures.
Use authenticated single-shard staging and delete it before returning the
receipt. No new controller.
