# GCEN G21 — activation-output replay admission (development only)

## Purpose and known outcome

G20 original-weight sharing-scale result is `NO_GO` for the tested group-shared
SVD formulation: the independent SVD control won every frozen 4/8/16-expert
comparison. Retain the original G20 artifacts and frozen source; do not optimize
the rejected basis against the same source and claim an independent finding.

G21 introduces **activation-output distortion** as a diagnostic alongside
weight-matrix reconstruction error. Neither one by itself is task success,
next-token performance, inference throughput, peak VRAM, or scientific acceptance.

The code in `activation_replay.py` can score a **provided** trace of expert
inputs and full/approximated weight matrices. It does not generate or
authenticate a real model trace. Focused tests use synthetic fixtures; passing
them does not admit any scientific evidence.

## Required authentic capture before original-model scoring

- Model: `Qwen/Qwen3-30B-A3B`, exact upstream revision
  `ad44e777bcd18fa416d9da3bd8f70d33ebb85d39`.
- Pin immutable source commit and original checkpoint shard receipts; verify
  relevant actual Safetensors bytes prior to tensor extraction.
- Capture **expert input activations at the original model forward hook** for a
  pinned named expert matrix, with a frozen development-only task cohort and
  documented task provenance. Synthetic random input vectors are permitted only
  for negative-control engineering tests, not for authentic-science claims.
- Freeze disjoint calibration/evaluation record IDs. No evaluation IDs may be
  seen by fitting, tuning, or budget/rank selection. Never access the protected
  cohort before separate science authorization.
- Preserve an immutable workflow run ID, exact source commit, independent
  activation-trace archive SHA-256 and manifest-bound record IDs and byte
  checksum. Inspect the **actual** original captured artifact independently
  before lifting the evaluator's `source_authenticity=UNVERIFIED` warning.
- Release only aggregate measurement receipts; do not upload private prompts,
  raw tokens, personal information, or secrets to the public repo.
- Compare original FP16, the frozen independent SVD control, and any
  genuinely different candidate using exact matched storage accounting.
  Include factors, codebooks, residuals, metadata and reconstruction buffers,
  not just the float16 factor count.
- Measure actual peak process RSS and GPU VRAM if available using a separately
  documented repeatable instrument; analytic workspace lower bounds are not
  measured peak memory.
- Record the expert-output normalized squared error and any directly measured
  task outcome separately. Negative/near-zero output error is **not** evidence
  of language-model task quality without original-model task execution.

## G21 admission / rejection

Engineering qualification requires exact-head Windows/Linux mandatory
`scripts/Test-Native.ps1` regression plus the complete GCEN focused suite.

`G21_ENGINEERING_PASS` never implies `ORIGINAL_ACTIVATIONS_VERIFIED`.
Authentic scientific scoring remains `BLOCKED` until independently attested
original-forward activations exist, and separate acceptance remains blocked
pending required task/VRAM/throughput measurements and independent authority.
No new controller, model provider, protected execution, deployment, or model
release is authorized by this protocol.
