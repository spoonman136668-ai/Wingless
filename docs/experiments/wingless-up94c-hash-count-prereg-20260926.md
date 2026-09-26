# Wingless UP-94C — seen-once hash-count ablation

Status: preregistered scientific bounded-memory admission experiment.

Scientific parent: sealed UP-93C 37acecd7c314f2c86e40ecff6d687382db864864.

## Question

UP-93C showed that generation salting does not solve the long-horizon boundary and can increase false positives. With width and generation size fixed, does a Bloom-style multi-hash count closer to the expected optimum reduce within-generation false admissions?

## Frozen memory and workload

- exact recall cap 16;
- two-bit-aging replacement;
- seen-once filter width 256 bits;
- continuous generation interval 32 filtered unseen writes;
- exact memory and aging state never reset;
- no query labels or future information;
- 32 original keys;
- 16 hot even keys;
- every hot key receives one second write before churn;
- 1536 unique one-shot churn writes;
- hot queries after every four second-hot writes and every four churn writes;
- value vocabulary 32;
- 64 episodes per arm;
- seeds 177M and 178M.

## Frozen arms

1. legacy2
   - exact UP-92C two-hash geometry.

2. multi2
3. multi4
4. multi6

For multi-k arms:
- derive k deterministic independent bit positions from sq0Mix64(key XOR fixed per-hash constant);
- a key is considered seen only when all k bits are already set;
- first-seen rejection sets all k bits;
- hash family is identical across multi2/multi4/multi6 except for how many positions are used;
- no generation salt.

Policy metadata remains 38 bytes for every arm; hash count is fixed by the experiment arm and requires no runtime storage.

## Metrics

For every arm:
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

Improvement as k increases would identify membership false-positive rate as the remaining C boundary. Degradation at high k would expose filter overfill from excessive bit setting.

## Bounds

No filter-width increase, no exact-memory expansion, no adaptive hash count, no salt, no query-derived labels, no future oracle, no interval tuning, no result-informed retry, no live activation, no production authority.
