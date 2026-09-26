# Wingless UP-88C — seen-once filter width ablation

Status: preregistered scientific bounded-memory admission experiment.

Scientific parent: sealed UP-87C f28e7790d8808ce5e61c2dd3be59d60f4541cb47.

## Question

UP-87C extended exact hot retention from 12 to 13 with an 8-byte seen-once filter, but failed at 16. Is the remaining failure driven by collision pressure in the compact filter rather than the admission rule itself?

## Frozen exact memory and admission rule

Identical to UP-87C:
- exact recall cap exactly 16;
- two-bit-aging replacement metadata 5 bytes;
- existing keys always update;
- before exact memory fills, unseen writes admit normally;
- once full, an unseen key is admitted only if both of its two seen-once hash bits were already set;
- otherwise both bits are set and the write is rejected;
- no query labels or future information affect admission.

## Frozen filter-width arms

Seen-once filter width:
- 64 bits  = 8 bytes; total policy metadata 13 bytes.
- 128 bits = 16 bytes; total policy metadata 21 bytes.
- 256 bits = 32 bytes; total policy metadata 37 bytes.

Hash family and all other logic remain fixed.

## Frozen workload

Use the exact UP-87C workload at hot-set size 16:
- 32 original keys;
- every hot key receives one second write;
- cold original keys receive no second write;
- 24 unseen churn keys receive exactly one write;
- hot queries every four writes during second-write and churn phases;
- value vocabulary 32;
- 64 episodes per cell;
- seeds 165M and 166M.

## Metrics

For each filter width:
- hot query accuracy;
- hot-set exact accuracy;
- cold query accuracy;
- recall entries used;
- policy metadata bytes;
- total bounded-memory bytes;
- rejected one-shot writes;
- second-write admissions.

## Interpretation

Recovery with wider filters would isolate false-positive collision pressure as the UP-87C limit. Persistent failure would implicate the admission schedule or exact-memory interaction instead.

## Bounds

No exact-memory capacity expansion, no query-derived labels, no future oracle, no adaptive hash count, no result-informed tuning, no attention, no live activation, no production authority.
