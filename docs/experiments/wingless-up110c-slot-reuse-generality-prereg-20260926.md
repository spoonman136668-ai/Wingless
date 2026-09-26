# Wingless UP-110C — phase-free slot-reuse generality

Status: preregistered scientific bounded-memory stress experiment.

Scientific parent: sealed UP-109C 5a328c4115d82061704a01b6431fb6b6bcb50dcb.

## Question

UP-109C converted the explicit post-promotion checkpoint into a local slot-reuse detector and restored exact retention across four tested orders. Does that detector remain exact across broader initial orders and substantially longer unique churn, and do redundant checkpoint emissions remain bounded?

## Frozen mechanism

Exact UP-109C slot_reuse_confirmation:
- exact recall cap 16;
- two-bit-aging replacement;
- 1024-bit two-hash seen-once filter;
- generation interval 32;
- 2-byte reuse-confirmation slot mask;
- checkpoint action = clear filter and set counter to 8;
- slot confirmation on existing-key write or successful seen-before admission;
- slot bit cleared on replacement;
- no phase label;
- no query-derived detector evidence;
- no future oracle.

## Initial-order grid

- ascending
- reverse
- evens_then_odds
- odds_then_evens
- interleaved_low_high
- rotate8

## Churn horizon ladder

- 12,288
- 24,576
- 49,152 unique one-shot writes

Promotion remains the exact two-pass alternating-ends sequence.

## Controls

1. explicit_checkpoint
2. slot_reuse_confirmation

## Sampling

- seeds 207M and 208M;
- 32 episodes per seed per cell;
- 64 episodes total per cell.

## Metrics

For every arm x order x horizon:
- resident hot count before churn;
- checkpoint count per episode;
- first checkpoint filtered-write index;
- false-positive admissions;
- final hot query accuracy;
- final hot-set exact accuracy;
- recall entries used;
- policy metadata bytes.

## Interpretation

Exact retention across all cells would establish a local phase-free checkpoint mechanism robust to order and long churn. Divergence at longer horizons would expose detector/checkpoint reuse interactions rather than the original phase-boundary defect.

## Bounds

No detector change, no interval tuning, no capacity expansion, no adaptive ordering, no query evidence, no future oracle, no result-informed retry, no live activation, no production authority.
