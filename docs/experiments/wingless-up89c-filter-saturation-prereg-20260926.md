# Wingless UP-89C — 256-bit seen-once filter saturation lifetime

Status: preregistered scientific bounded-memory admission experiment.

Scientific parent: sealed UP-88C db9f1c6d6478cd5c2f14268163d7dd9e996c01da.

## Question

UP-88C restored exact retention of all 16 hot items with a 256-bit, two-hash seen-once filter. How long does that fixed filter remain effective as unique one-shot churn accumulates without resetting it?

## Frozen memory and admission rule

Identical to the UP-88C 256-bit arm:
- exact recall cap 16;
- two-bit-aging replacement metadata 5 bytes;
- 256-bit seen-once filter = 32 bytes;
- total policy metadata 37 bytes;
- two deterministic hashes;
- existing keys update normally;
- once exact memory is full, first-seen unseen keys are rejected and mark both filter bits;
- an unseen key whose two bits are already set is admitted using ordinary two-bit-aging replacement;
- no filter reset during an episode;
- no query labels or future information.

## Frozen workload

- 32 original keys;
- hot set = all 16 even original keys;
- hot keys receive one second write before churn;
- cold original keys receive no second write;
- unique one-shot churn ladder:
  - 24
  - 48
  - 96
  - 192 writes
- after every four second-hot writes and every four churn writes, query every hot key;
- missing keys are not reconstructed;
- value vocabulary 32;
- 64 episodes per cell;
- seeds 167M and 168M.

## Metrics

For each churn depth:
- hot query accuracy;
- hot-set exact accuracy;
- cold query accuracy;
- recall entries used;
- rejected one-shot writes;
- false-positive churn admissions;
- policy metadata bytes;
- total bounded-memory bytes.

## Interpretation

The first churn depth that breaks exact retention estimates the practical saturation boundary of this compact admission memory.

## Bounds

No filter reset, no capacity expansion, no adaptive hash count, no query-derived labels, no future oracle, no threshold tuning, no result-informed retry, no live activation, no production authority.
