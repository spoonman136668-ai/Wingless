# Wingless UP-98C — post-promotion churn generations

Status: preregistered scientific bounded-memory admission experiment.

Scientific parent: sealed UP-97C a52f15992f6e0476f2f36eb88d222822694a9772.

## Question

UP-97C showed that generation-phase offsets trade collision avoidance against loss of legitimate hot-key promotion evidence because one filter serves both roles. If promotion is allowed to complete first, can a fresh churn-only filter generation eliminate the long-horizon collision while preserving all promoted hot keys?

## Frozen memory mechanism

- exact recall cap 16;
- two-bit-aging replacement;
- 1024-bit two-hash seen-once filter;
- continuous generation interval 32 filtered unseen writes;
- exact memory and aging state never reset;
- no query labels or future information.

## Frozen workload

- 32 original keys;
- 16 hot even keys;
- every hot key receives one second write before churn;
- hot queries after every four second-hot writes;
- after the entire hot second-write promotion phase completes:
  - clear only the seen-once filter;
  - exact memory and aging state remain untouched;
  - set the generation counter to the arm's fixed phase offset;
- then process 12288 unique one-shot churn writes;
- hot queries every four churn writes;
- value vocabulary 32;
- 64 episodes per arm;
- seeds 185M and 186M.

## Frozen churn-phase offsets

- 0
- 8
- 16
- 24

After the first churn reset, every generation is exactly 32 filtered unseen writes.

## Metrics

For every offset:
- hot query accuracy;
- hot-set exact accuracy;
- cold query accuracy;
- recall entries used;
- rejected one-shot writes;
- false-positive churn admissions;
- first false-positive churn index;
- filter reset count;
- policy metadata bytes;
- total bounded-memory bytes.

## Interpretation

An exact arm would show that promotion evidence and churn rejection can be separated cleanly, and that the remaining collision is generation alignment rather than exact-memory capacity. Failure across all offsets would point back to hash geometry.

## Bounds

No filter-width change, no interval change, no hash-count change, no salt, no adaptive phase, no query-derived labels, no future oracle, no result-informed retry, no live activation, no production authority.
