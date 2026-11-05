# CNC-G9P — hardware-corrected capacity gate (historical, not live)

This gate reuses existing numerical/storage engines G9B/G9D/G9H/G9K,
without building another controller or inference backend. The original
source-of-truth GPU hardware receipt was independently retrieved from
`spoonman136668-ai/FORGE` Actions run `41030389354`,
artifact `12966370600`; verified ZIP SHA256
`c6d368b70e6697e4ec699c7943100401214eed5fcea17688e3ed6351b564bac2`.
Original JSON records one RTX 5070 Ti Laptop GPU, driver 596.13,
physical `12227 MiB` and **historically available** `11812 MiB`.

**Crucial:** `12227 MiB` is **61 MiB below** the `12288 MiB` original
G9B nominal ceiling, and free VRAM at that historical moment was
**476 MiB below** nominal. These are *not current measurements* and
cannot certify runtime fit or original-model launch.

All G9P cases reuse the exact G9B 18,867-tensor architectural byte
census, including BF16 nonexperts, 8,192-token BF16 KV and the G9B
hypothetical, unmeasured 1 GiB workspace reserve. Rebase the G9F
INT2 block256 baseline, G9H k1/2/3 sparse corrections, G9K rank1/2/3
low-rank corrections, and frozen G9D fixed INT4-protection masks.
Preserve failed cases, all actual physical index/factor metadata, and
signed remaining-byte deficits. **Analytical budget alone is NEVER
CUDA admission**; fragmentation, context allocation and runtime copy
buffers are not measured.

G9P deliberately does not install CUDA/Torch, discover current runner
state, download checkpoint shards, start the full Qwen model, inspect
secrets or workloads, or claim authentic activation quality. Previous
G9N/G9O boundaries and G9M capture integrity rules remain in force.

Next eligible original-model runner action needs independently
verified current availability, workstation owner isolation/scheduling,
a source-pinned runtime and 18-file checkpoint-custody receipt,
real peak free VRAM/workspace measurements and task/activation evidence.
No deployment, protected-main change or scientific acceptance.
