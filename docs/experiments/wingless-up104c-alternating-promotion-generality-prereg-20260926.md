# Wingless UP-104C — alternating-ends promotion generality

Status: preregistered scientific bounded-memory ordering experiment.

Scientific parent: sealed UP-103C 63a7d329d9adee4fde510ddf3279143ab271c553.

## Question

UP-103C showed that a fixed alternating-ends hot second-write order preserves all 16 hot entries for both ascending and reverse initial orders. Does that deterministic promotion schedule remain robust across broader initial permutations?

## Frozen memory mechanism

Identical to UP-103C:
- exact recall cap 16;
- two-bit-aging replacement;
- 1024-bit two-hash seen-once filter;
- generation interval 32;
- trigger after 8 successful repeat admissions;
- post-trigger counter 8;
- no extra diagnostic reads;
- no adaptive ordering;
- no future oracle.

## Frozen hot second-write schedule

Use only:
- alternating_ends: 0,30,2,28,...,14,16

This order is fixed for every cell.

## Initial-order grid

1. ascending
2. reverse
3. evens_then_odds
4. odds_then_evens
5. interleaved_low_high: 0,31,1,30,...,15,16
6. rotate8: 8,9,...,31,0,1,...,7

All orders are deterministic and preregistered.

## Workload

- 32 original keys;
- all 16 even keys hot;
- one second write per hot key using alternating_ends;
- then 12,288 unique one-shot churn writes;
- ordinary hot-query cadence after every four promotion/churn writes;
- 64 episodes per cell;
- seeds 197M and 198M.

## Metrics

For each initial order:
- resident hot count immediately before churn;
- trigger-fired episodes;
- mean trigger index;
- false-positive churn admissions;
- final hot accuracy;
- final hot-set exact accuracy.

## Interpretation

Broad success would establish alternating-ends as a deterministic order-robust promotion schedule. Failure only when all hot keys start resident would keep the separate trigger-blind problem isolated.

## Bounds

No order adaptation, no extra reads, no trigger change, no filter change, no capacity expansion, no result-informed retry, no live activation, no production authority.
