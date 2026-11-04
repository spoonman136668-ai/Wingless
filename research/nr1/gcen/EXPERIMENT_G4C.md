# GCEN G4C — development-only specification

Original NR-1A baseline: cf23d6e50c379ba2ffea58f7b995718a878e7dd9 (immutable).
Qwen full precision: Qwen/Qwen3-30B-A3B revision ad44e777bcd18fa416d9da3bd8f70d33ebb85d39.

Status: NOT QUALIFIED. All 16 Safetensors shard SHA-256 hashes, config and index hashes must be verified against downloaded bytes before real-weight measurements. Do not run inference or protected holdout.

Arms: original FP16; independent truncated SVD; shared right basis across expert matrices. Report stored bytes including factors, singular values, codebooks, metadata, residuals, runtime reconstruction buffers, per-expert error, latency and peak memory. Match physical storage budget; 4x expert-only compression target. Five frozen development seeds: 11, 23, 37, 53, 71. Preserve NR-1A source and measurements. Mandatory exact-head Windows and Linux full regression and focused tests precede acceptance. This specification is not a scientific acceptance receipt.
