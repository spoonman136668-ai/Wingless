# Wingless UP-97C — generation-phase alignment

Status: preregistered scientific bounded-memory admission experiment.

Scientific parent: sealed UP-96C 1a9d4ff092ab58db8741d5a526ea6b352867956e.

## Question

UP-96C showed exactly one false admission per episode beginning between 1536 and 3072 churn writes, with no additional false admissions through 12288. Does that isolated collision depend on the deterministic alignment of 32-write filter generations?

## Frozen memory mechanism

Identical to UP-96C:
- exact recall cap 16;
- two-bit-aging replacement;
- 1024-bit two-hash seen-once filter;
- continuous generation interval 32 filtered unseen writes;
- exact memory and aging state never reset;
- no query labels or future information.

## Frozen phase arms

Generation counter initial value before the first filtered unseen write:
- offset 0
- offset 8
- offset 16
- offset 24

After the first reset, every subsequent generation is exactly 32 filtered unseen writes.

No arm uses workload phase labels or result feedback.

## Frozen workload

- 32 original keys;
- 16 hot even keys;
- every hot key receives one second write before churn;
- 12288 unique one-shot churn writes;
- hot queries after every four second-hot writes and every four churn writes;
- value vocabulary 32;
- 64 episodes per arm;
- seeds 183M and 184M.

## Metrics

For every phase offset:
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

An offset that removes the isolated failure would identify generation alignment as the remaining deterministic collision mechanism. Failure across all offsets would move the frontier to hash geometry rather than phase.

## Bounds

No filter-width change, no interval change, no hash-count change, no salt, no adaptive phase, no query-derived labels, no future oracle, no result-informed retry, no live activation, no production authority.
