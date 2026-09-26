# Wingless UP-107C — post-promotion admission-state checkpoint

Status: preregistered diagnostic bounded-memory experiment.

Scientific parent: sealed UP-106C 84ddc7a82455820c20efa64c9647f332c92c395a.

## Question

UP-106C showed that a second hot promotion pass restores all 16 hot entries for odds-then-evens and rotate8, but one later false-positive admission still evicts a hot key. Is stale seen-once filter state after promotion the sole residual cause?

## Frozen memory/workload

Identical to UP-106C:
- exact recall cap 16;
- two-bit-aging replacement;
- 1024-bit two-hash seen-once filter;
- generation interval 32;
- trigger threshold 8;
- post-trigger counter 8;
- alternating_ends promotion order;
- query every four promotion/churn writes;
- two promotion passes for the experimental hard cases;
- 12,288 unique one-shot churn writes.

Initial orders:
- odds_then_evens
- rotate8
- evens_then_odds
- ascending control

## Frozen arms

1. two_pass_control
   - exact UP-106C two-pass state.

2. post_promotion_filter_clear
   - after the second promotion pass, clear only the seen-once filter and set its generation counter to zero;
   - exact memory entries and two-bit-aging state are untouched.

3. post_promotion_filter_clear_counter8
   - same filter clear;
   - set generation counter to 8, matching the established trigger-aligned generation offset.

This experiment intentionally uses the known promotion boundary as a diagnostic phase label. It is not a phase-free production mechanism.

## Metrics

- resident hot count before churn;
- false-positive churn admissions;
- trigger-fired episodes;
- final hot accuracy;
- final hot-set exact accuracy;
- recall entries used.

## Interpretation

If either checkpoint arm restores exact retention for odds-then-evens/rotate8, stale post-promotion admission evidence is confirmed as the sole residual cause. Evens-then-odds remains a separate trigger-blind control.

## Bounds

No exact-memory reset, no capacity expansion, no extra reads, no adaptive checkpoint timing, no future oracle, no result-informed retry, no live activation, no production authority.
