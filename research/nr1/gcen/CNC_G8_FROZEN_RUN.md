# CNC-01Q-G8 — frozen original-weight one-shot development comparison

Status: **development-only experiment**; no scientific acceptance, merge, activation,
whole-model inference, or protected held-out evaluation.

## Identity and custody

- Model: `Qwen/Qwen3-30B-A3B`
- Pinned revision: `ad44e777bcd18fa416d9da3bd8f70d33ebb85d39`
- Source: the one exact BF16 original Safetensors shard already used by G32.
- Byte authentication and source headers: reuse G31/G32 fail-closed SHA-256,
  index, shape, source-identity and single-shard staging checks.
- Scratch source must be removed before the receipt is sealed.

## Fixed development protocol

- Select the first fully co-resident group of 16 original `down_proj.weight`
  experts in the pinned shard, using the existing G32 rule.
- Preserve G32's five seeds `11,23,37,53,71`, 128×256 BF16 window,
  and 4/8/16 nested expert comparisons.
- Baselines: original BF16 weights and independent rank-14 truncated SVD.
- Candidate: per-row symmetric packed INT4, block 64, float32 scales.
- Inputs for output-distortion probes are **synthetic Gaussian** batch 32
  and are never described as original-model activations.
- Count codes, scales, metadata and SVD factor storage; retain per-expert
  numerical errors, aggregate normalized weight/output distortion and peak
  process RSS. Report no GPU VRAM or model throughput claim.
- Do not retune ranks, quantization blocks, expert groups, seeds, thresholds,
  or evaluation samples based on these results.

## Execution boundary

The existing `.github/workflows/native.yml` job `cnc_g8_real_weight_int4`
is enabled only by a one-shot push to the isolated CNC development branch
whose commit message includes `[cnc-g8-real]`. Its `needs: native`
requires exact-head Linux and Windows native/focused qualification first.
On any failure preserve FAIL/UNVERIFIED; do not relabel infrastructure or
source-custody failure as a scientific negative.

This original-weight experiment does not execute a full Qwen forward pass,
measure task quality, activate a model, mutate `main`, or alter the frozen
G32 experiment.
