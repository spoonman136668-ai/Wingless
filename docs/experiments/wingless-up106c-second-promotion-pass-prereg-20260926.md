# Wingless UP-106C — repeated hot promotion pass

Status: preregistered scientific bounded-memory promotion experiment.

Scientific parent: sealed UP-105C 8534d020165f5f666e924db88f8a54dd2ce67f16.

## Question

UP-105C ruled out promotion-query cadence as the cause of missing hot entries. Can a second deterministic pass of the same observed hot writes close the underpromotion boundary without extra reads, future labels, or capacity expansion?

## Frozen mechanism

Identical to UP-105C:
- exact recall cap 16;
- two-bit-aging replacement;
- 1024-bit two-hash seen-once filter;
- generation interval 32;
- trigger threshold 8;
- post-trigger counter 8;
- second-write order alternating_ends;
- query every four promotion writes;
- query every four churn writes;
- 12,288 unique one-shot churn writes.

## Initial-order grid

1. ascending
2. reverse
3. evens_then_odds
4. odds_then_evens
5. rotate8

## Frozen arms

1. one_pass
   - exact UP-105C query_every4 control.

2. two_pass
   - execute the same 16-key alternating_ends promotion sequence a second time immediately after the first pass;
   - use the same query-every-four cadence in both passes;
   - values are unchanged;
   - no reset occurs between passes.

The second pass is additional observed reuse evidence, not a query label or future-demand oracle.

## Metrics

For each initial order x arm:
- resident hot count after pass 1;
- resident hot count after final promotion pass;
- trigger-fired episodes;
- false-positive churn admissions;
- final hot accuracy;
- final hot-set exact accuracy;
- recall entries used.

## Interpretation

If two_pass closes odds-then-evens or rotate8, the remaining failure is insufficient repeated write evidence. If not, the replacement topology cannot converge even with repeated hot observations.

## Bounds

No capacity expansion, no extra diagnostic reads, no filter/trigger change, no adaptive order, no future oracle, no result-informed retry, no live activation, no production authority.
