# Wingless UP-108C — phase-free resident-reuse checkpoint

Status: preregistered scientific bounded-memory control experiment.

Scientific parent: sealed UP-107C e96eec6518d83079735325a1728d3db4d191cc52.

## Question

UP-107C proved that clearing the seen-once filter and setting generation counter=8 immediately after complete promotion restores exact retention, but that diagnostic used an explicit promotion boundary. Can the same transition be emitted from local resident-write evidence only?

## Frozen memory/workload

Keep UP-107C unchanged:
- exact recall cap 16;
- two-bit-aging replacement;
- 1024-bit two-hash seen-once filter;
- generation interval 32;
- trigger threshold 8;
- alternating_ends hot promotion;
- two promotion passes;
- query every four promotion/churn writes;
- 12,288 unique one-shot churn writes;
- no future oracle or query-derived admission labels.

Initial orders:
- ascending
- reverse
- evens_then_odds
- odds_then_evens
- rotate8

## Phase-free detector

Add a 16-bit resident-reuse mask indexed by exact-memory slot.

Rules:
- on a write to a key already resident in exact memory, mark that resident slot as reused;
- whenever a slot is replaced by a different key, clear that slot's reuse bit before processing the new resident state;
- if all 16 current resident slots have been reused since their most recent replacement:
  - clear the seen-once filter;
  - set the generation counter to 8;
  - clear the resident-reuse mask;
  - increment checkpoint count.

The detector sees only current exact-memory residency and writes. It receives:
- no hot/cold label;
- no phase label;
- no query result;
- no future demand.

Metadata cost: 2 bytes for the reuse mask, in addition to existing policy metadata.

## Frozen arms

1. explicit_checkpoint
   - exact UP-107C post_promotion_filter_clear_counter8 diagnostic.

2. phasefree_reuse_checkpoint
   - no explicit promotion checkpoint;
   - use only the detector above.

3. no_checkpoint
   - exact UP-107C two-pass control.

## Metrics

- checkpoints fired per episode;
- first checkpoint filtered-write index;
- resident hot count before churn;
- false-positive churn admissions;
- final hot accuracy;
- final hot-set exact accuracy;
- total policy metadata.

## Interpretation

Matching the explicit checkpoint across initial orders would convert the C107 diagnostic mechanism into a local phase-free control rule. Premature or repeated checkpoints that damage retention are valid negatives.

## Bounds

No hot labels, no explicit phase label in the phase-free arm, no future oracle, no query-derived detector state, no capacity expansion, no adaptive threshold, no result-informed retry, no live activation, no production authority.
