# CNC-01Q-G9B — pinned-index exact tensor-class VRAM accounting

Status: **development-only analytical implementation**. No GPU, original
weights, native Qwen activation, inference, scientific admission, or promotion.

Pinned model: `Qwen/Qwen3-30B-A3B` revision
`ad44e777bcd18fa416d9da3bd8f70d33ebb85d39`.

This gate builds a deterministic 18,867-key Qwen3-MoE tensor inventory from
the pinned, SHA-256-verified config metadata. The inventory is reconciled to
the independently SHA-256-verified Safetensors index: exact key set, 16 shards,
and exact BF16 tensor payload sum. The original G27 model meta audit is a
separate prior implementation-independent identity cross-check; this method
does **not** read or authenticate individual original tensor headers.

Frozen scenarios (all **model weights resident** in exactly 12 GiB nominal):
BF16 nonexperts, all expert projections stored at 4/3/2-bit and block-64
FP32 per-row scales, plus hypothetical 2-bit block-128/256 FP16 scales.
Per-tensor 24-byte metadata follows G8 accounting. 8,192-token BF16 KV
cache, batch 1, and an **illustrative unmeasured 1 GiB workspace reserve**.
Model runtime fragmentation, actual workspace, GPU allocator, kernel state,
codebook/zero-point cost, and fidelity are **not** established.

The G8 block-64 FP32 INT4 operator is CPU weight-window only. **2/3-bit
operators and FP16-scale variants are not implemented or qualified**.
Positive budget headroom cannot establish resident inference. Negative
headroom under these assumptions is a mathematical infeasibility result
for that specific scenario, not a negative for unrelated compression.

Offline tests: `python research/nr1/gcen/run_focused.py`.
Read-only metadata execution after platform qualification:

```bash
python research/nr1/gcen/g9b_tensor_budget.py \
  --manifest research/nr1/gcen/upstream_full_manifest_unverified_weights.json \
  --source-sha "<exact checked-out Git SHA>" \
  --output evidence/g9b-budget.json
```

No automatic original-weight download or unapproved model runtime activation.
For original-source qualification, run the metadata-only CLI via an existing,
separately authorized execution path and preserve the original receipt.
