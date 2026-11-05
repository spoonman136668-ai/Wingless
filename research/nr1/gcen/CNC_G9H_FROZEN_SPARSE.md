# CNC-G9H — fixed sparse residual over G9F INT2

**Development only.** Follow G9G immutable receipt and G9F/Qwen original
positive. Freeze independent synthetic seed 811, 16×256 weights,
batch7 Gaussian activations and 20x outliers at column 15 for the
unit-level mechanistic screen. Use integer corrections-per-row
`k in [1,2,3]`; select top absolute *weight residual* per row with
stable tie breaking. Encode one uint16 column index plus one FP16
correction value for each chosen position; include a uint32 per-row
CSR offset table, including the terminal entry. All bytes must be
accounted, not merely theoretical 2 bits per parameter.

**12 GiB planning** uses exact Qwen architectural tensor census from
G9B, prior SHA256-authenticated source metadata independent of this
synthetic fixture, 8K BF16 KV and unmeasured 1-GiB workspace reserve.
Check each 1/2/3-per-row candidate for actual physical sparse index,
value and pointer overhead across *all* 18,432 expert projections.
No claiming original index was freshly fetched in this gate.

The CPU decoder operates blockwise G9F INT2 and applies only sparse
corrections; it never reconstructs a full model or GPU representation.
No source data, protected heldout, scientific acceptance, real model
activations, actual GPU memory measurement, task quality, or promotion.

Next G9I needs new, independently frozen original BF16 expert windows,
not the existing G9E/G9F five seeds, and failure-closed SHA256 custody.
The synthetic tests must never be read as authentic-Qwen evidence.
