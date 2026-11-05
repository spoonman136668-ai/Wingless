# CNC-G9L — frozen original-weight low-rank residual development gate

G9K distinct mechanism qualified exact-head Linux/Windows at
`97c34fa79bc664a5ccd725ad19d9a151899e9d9d` (Actions
`41807892629`). Reuse existing one-source-shard SHA-256 custody validator
for Qwen/Qwen3-30B-A3B pinned `ad44e777bcd18fa416d9da3bd8f70d33ebb85d39`.

Independent frozen development seeds **541, 563, 587, 613, 641**, disjoint
from prior G9E/F `11,23,37,53,71` and G9I
`107,163,229,271,347`. Original expert BF16 128×256 windows, first
16 distinct complete source experts (4/8/16 prefixes), identical seeded
synthetic Gaussian activation probes batch 32 using 71000+seed.
Reference original BF16, G8 INT4 block64 FP32 scale, G9F companded INT2
block256 FP16 scale, and rank1/2/3 FP16 low-rank factors of the INT2
quantization residual. Each factor stores 2*rank*(rows+cols) additional
bytes and 24-byte tensor header; no hidden metadata. SVD happens only at
development encode time; compressed forward uses two factor matmuls.

Run only one original 4-GB shard at a time on the existing isolated
GitHub runner, gated on exact-head dual-platform full native regression.
No frozen method retuning after authentic results; preserve negatives
and all unpromising ranks. Original weights authenticated on this trial,
but actual Qwen model activations, full-model forward, 12-GiB GPU
residency, scientific heldout, reasoning quality and deployment remain
unverified. No merge or production changes.
