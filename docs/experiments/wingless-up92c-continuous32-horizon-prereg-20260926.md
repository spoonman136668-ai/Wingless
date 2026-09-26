# Wingless UP-92C — continuous-32 long-horizon stability

Status: preregistered scientific bounded-memory admission experiment.

Scientific parent: sealed UP-91C db52aa6b7ede98945d7b22e065b1b2a23aec4b53.

## Question

UP-91C achieved exact preservation of all 16 hot items through 192 unique one-shot churn writes using a continuous 32-write seen-once generation. Does that fixed mechanism remain exact at substantially longer churn horizons?

## Frozen memory and admission system

Identical to the UP-91C continuous_32 arm:
- exact recall cap 16;
- two-bit-aging replacement;
- 256-bit two-hash seen-once filter;
- one-byte continuous generation counter;
- 38 bytes policy metadata total;
- generation counter begins on the first filtered unseen write after exact memory is full;
- existing exact-memory updates do not increment the counter;
- every filtered unseen write increments the counter;
- after processing the 32nd filtered write in a generation, clear only the filter and reset the counter;
- exact memory and aging state are never reset;
- no query labels or future information.

## Frozen workload

- 32 original keys;
- 16 hot even keys;
- each hot key receives one second write before churn;
- unique one-shot churn ladder:
  - 192
  - 384
  - 768
  - 1536 writes
- hot queries after every four second-hot writes and every four churn writes;
- value vocabulary 32;
- 64 episodes per cell;
- seeds 173M and 174M.

## Metrics

For every churn depth:
- hot query accuracy;
- hot-set exact accuracy;
- cold query accuracy;
- recall entries used;
- rejected one-shot writes;
- false-positive churn admissions;
- filter reset count;
- policy metadata bytes;
- total bounded-memory bytes.

## Interpretation

Stable exact retention across the extended horizon would support a bounded generational admission mechanism whose effectiveness does not decay merely with elapsed unique churn. Failure identifies the horizon where within-generation collisions accumulate despite deterministic clearing.

## Bounds

No interval tuning, no exact-memory expansion, no adaptive reset, no query-derived labels, no future oracle, no threshold tuning, no result-informed retry, no attention, no live activation, no production authority.
