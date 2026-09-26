# Wingless UP-95C — generational filter-width frontier

Status: preregistered scientific bounded-memory admission experiment.

Scientific parent: sealed UP-94C 43727974dc4d24e3244e8eec0e069a1e8c16ae48.

## Question

UP-94C showed that increasing Bloom hash count does not improve the 256-bit continuous-generation filter. How much additional filter width is required to eliminate the remaining within-generation false admissions at the 1536-write horizon?

## Frozen mechanism

- exact recall cap 16;
- two-bit-aging replacement;
- legacy two-hash geometry generalized to power-of-two filter widths by taking the corresponding high hash bits;
- continuous generation interval exactly 32 filtered unseen writes;
- exact memory and aging state never reset;
- no query labels or future information.

## Frozen width arms

1. 256 bits
   - 32 filter bytes + 5 aging bytes + 1 generation counter byte = 38 metadata bytes.
   - exact UP-92C legacy geometry.

2. 512 bits
   - 64 filter bytes + 5 aging bytes + 1 counter = 70 metadata bytes.

3. 1024 bits
   - 128 filter bytes + 5 aging bytes + 1 counter = 134 metadata bytes.

Hash count remains exactly 2 in every arm.

## Frozen workload

- 32 original keys;
- 16 hot even keys;
- every hot key receives one second write before churn;
- 1536 unique one-shot churn writes;
- hot queries after every four second-hot writes and every four churn writes;
- value vocabulary 32;
- 64 episodes per arm;
- seeds 179M and 180M.

## Metrics

For every width:
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

This directly measures the memory cost required to suppress the remaining probabilistic membership errors under fixed exact-memory capacity and generation timing.

## Bounds

No exact-memory expansion, no hash-count change, no salt, no adaptive width, no query-derived labels, no future oracle, no interval tuning, no result-informed retry, no live activation, no production authority.
