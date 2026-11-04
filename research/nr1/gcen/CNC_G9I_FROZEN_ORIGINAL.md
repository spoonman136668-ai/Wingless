# CNC-G9I — independent authentic Qwen sparse-correction development trial

Prerequisite G9H exact-head Linux/Windows PASS `6337b5333ab81fb08bb2e595c6dfb8d2dfdbf1a2`.
Frozen original `Qwen/Qwen3-30B-A3B` source checkpoint revision
`ad44e777bcd18fa416d9da3bd8f70d33ebb85d39`;
G8/G9F single-shard SHA-256 custody validator, known 4-GB shard,
first complete 16 original expert down projections.

**Separately selected development seeds**: `107,163,229,271,347`,
mutually disjoint from G9E/G9F seeds 11,23,37,53,71. Original 128x256
BF16 windows, independently seeded synthetic Gaussian batch32 inputs
at `51000+seed`. Distinct original 4/8/16 expert group prefixes,
same G9F compressed INT2 block256 vs G9H fixed k=1,2,3 per-row
sparse corrections, plus earlier frozen G8 INT4 block64 and G9E
uniform INT2 block128/256 as higher/lower-byte controls. G9H
kernel/metadata frozen before access; no retuning after authentic output.

Compare original BF16 operator output-squared fidelity and weight error.
Record per-variant physical correction index, FP16 value and uint32
row-offset bytes (not just ideal bit rate). Full-model analytical budget
from G9H uses reconstructed source-authenticated architecture, not
fresh metadata download; **original weights are authenticated**
in this experiment with full SHA-256 verification of one original
shard, and staged source cleaned after receipt.

No original model activations, live GPU memory, full-model forward,
task quality, sealed heldout, scientific acceptance, protected main
mutation, automatic promotion or deployment. Trial is one-shot on
existing isolated PR branch only following exact-head native
Linux/Windows qualification; preserve failures and artifacts.
