# CNC-G9G — bounded robustness and selective INT4 protection

**Prerequisites:** pinned draft Wingless PR #443; G9F original-weight
checksum-custodied numerical result at `45ddc4c983f7cd9c4144960e0e005aa61a3110ce`.
This is an **independent synthetic development stress screen**, not
another original-shard download. It preserves the authentic G9E negative
and G9F positive without re-tuning either representation.

Frozen, non-retuned synthetic weight distributions: (a) standard
Gaussian, (b) Student-t degrees of freedom 3 scaled by sqrt(3)
to unit variance in expectation, (c) Gaussian with independent
1% probability of 16x outlier amplification. Seeds 101,211,307,
128x256 matrix, Gaussian activation batch32 seeded 60000+seed.
Compare fixed G9F nonuniform RMS INT2 block128/256, frozen G9C peak
INT2 at identical physical bytes, and G8 INT4 block64 at higher bytes.
Measure relative output **squared** error, mean/worst family error,
both storage costs, and number of seeds where companding outperforms
uniform. All alternatives receive exactly the same inputs/weights.

Mixed-precision analysis reuses G9D frozen original-index-compatible
tensor-class census and precommitted masks (first/last/both boundary
layers and first 1/2 expert IDs each layer), with G9B exact storage
accounting, BF16 nonexpert weights, 8K BF16 KV, and an illustrative
**unmeasured** 1GiB runtime reserve. These masks are **not actual
router or sensitivity rankings**. A positive analytical fit does not
verify GPU residency or full-model task/reasoning quality.

**Gate output** is a reproducible synthetic stress diagnosis and
existing-mask memory feasibility. No authentic weights, no protected
held-out, no model forward, no GPU activation, no scientific acceptance,
no protected-main mutation. Pass means source integrity and controlled
diagnostics only; a worse stress-case result must remain visible,
not be tuned away. Next meaningful advancement requires independently
frozen authentic weight- and activation-sensitivity evidence.
