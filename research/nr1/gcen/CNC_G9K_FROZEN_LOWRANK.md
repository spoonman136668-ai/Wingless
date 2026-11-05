# CNC-G9K — low-rank residual atop compressed-direct INT2

Unlike G9H/I sparse per-row corrections, G9K stores compact **dense low-rank
residual factors** for the G9F companded INT2 reconstruction error. This is
not the prior unsuccessful standalone rank14 SVD replacement. No original
model weights or activations are consumed by this first gate.

Frozen CPU mechanistic source: Gaussian 32×256 weight matrix, seed 401,
Gaussian activation batch6, and immutable correction ranks 1/2/3.
Encoder-only SVD factors the quantization residual. Runtime computes
`INT2_direct(x) + (x V) Uᵀ` from FP16 factors without reconstructing the
dense weight matrix. Include both factor arrays and 24-byte tensor header
in exact physical storage. Qualify malformed/zero/overflow boundaries.

Whole-model analytical budget: reconstruct the 18,867 original Qwen tensor
keys from the SHA256-authenticated G9B schema. For *each* 18,432 expert
weight matrix, account 2×rank×(rows+cols) FP16 factor bytes plus 24 bytes.
Retain BF16 nonexperts, 8K BF16 KV, nominal 12 GiB, and an **unmeasured**
1 GiB workspace reserve. Analytical fit is not GPU measurement.

Earlier G9I original BF16 windows and the G9J NO_GO for universally
distributed sparse corrections remain immutable. G9K is a distinct
mechanism hypothesis and must not claim superiority until tested on new
disjoint original-weight seed windows using the established SHA256
single-shard provenance and mandatory exact-head Linux/Windows regressions.
No full model forward, authentic Qwen activations, heldout scientific
acceptance, deployment, new controller, or protected main changes.
