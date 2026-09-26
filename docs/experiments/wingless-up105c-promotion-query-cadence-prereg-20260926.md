# Wingless UP-105C — promotion-phase query cadence

Status: preregistered scientific bounded-memory promotion experiment.

Scientific parent: sealed UP-104C d34cdb9c75d5b0209969b9dd7244f734ee153ae8.

## Question

UP-104C showed that fixed alternating-ends promotion is exact for some initial orders but leaves hot keys missing for odds-then-evens and rotate8. Do queries interleaved during the promotion phase preserve incumbent entries strongly enough to block missing hot-key admission?

## Frozen memory and promotion order

Identical to UP-104C:
- exact recall cap 16;
- two-bit-aging replacement;
- 1024-bit two-hash seen-once filter;
- generation interval 32;
- trigger after 8 successful repeat admissions;
- post-trigger counter 8;
- hot set = all 16 even original keys;
- second-write order fixed to alternating_ends;
- no future oracle;
- no adaptive ordering.

## Initial-order grid

1. ascending
2. reverse
3. evens_then_odds
4. odds_then_evens
5. rotate8

## Frozen promotion query-cadence arms

1. query_every4
   - exact UP-104C control: query every hot key after each group of four promotion writes.

2. query_every8
   - query every hot key after each group of eight promotion writes.

3. query_after16
   - no promotion queries until all 16 hot second writes have completed, then issue exactly one hot sweep.

4. no_promotion_queries
   - no hot query during or immediately after promotion.

After promotion, all arms use the exact same churn phase:
- 12,288 unique one-shot churn writes;
- query every hot key after every four churn writes.

## Metrics

For every initial order x cadence:
- resident hot count immediately before churn;
- trigger-fired episodes;
- mean trigger index;
- false-positive churn admissions;
- final hot accuracy;
- final hot-set exact accuracy;
- recall entries used.

## Interpretation

Improved pre-churn occupancy when promotion queries are delayed or removed would show that query-driven aging is causally blocking missing-key promotion. No improvement would move the failure to write/replacement ordering itself.

## Bounds

No extra diagnostic reads beyond the declared cadence, no trigger/filter/capacity change, no adaptive ordering, no future oracle, no query-derived admission labels, no result-informed retry, no live activation, no production authority.
