# Wingless UP-96C — 1024-bit generational horizon

Status: preregistered scientific bounded-memory admission experiment.

Scientific parent: sealed UP-95C 7f4192c8341908834d727ad97678c1080c865c8d.

## Question

UP-95C showed that a 1024-bit two-hash filter with continuous 32-write generations preserves all 16 hot items exactly through 1536 unique one-shot churn writes. Does that exactness persist at substantially longer horizons?

## Frozen mechanism

Exact UP-95C 1024-bit arm:
- exact recall cap 16;
- two-bit-aging replacement;
- 1024-bit seen-once filter;
- legacy two-hash geometry generalized with the 10 high hash bits;
- continuous generation interval exactly 32 filtered unseen writes;
- one-byte generation counter;
- 134 policy metadata bytes total;
- exact memory and aging state never reset;
- no query labels or future information.

## Frozen workload

- 32 original keys;
- 16 hot even keys;
- every hot key receives one second write before churn;
- unique one-shot churn ladder:
  - 1536
  - 3072
  - 6144
  - 12288 writes
- hot queries after every four second-hot writes and every four churn writes;
- value vocabulary 32;
- 64 episodes per cell;
- seeds 181M and 182M.

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

Stable exact retention across this ladder would show that the widened generational filter has a long deterministic operating horizon under unique churn. The first failure identifies the next collision boundary without changing exact-memory capacity.

## Bounds

No filter-width change, no generation-interval tuning, no hash-count change, no salt, no adaptive reset, no query-derived labels, no future oracle, no result-informed retry, no live activation, no production authority.
