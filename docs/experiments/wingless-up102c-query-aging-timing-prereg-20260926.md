# Wingless UP-102C — diagnostic-query aging timing

Status: preregistered scientific bounded-memory state-sensitivity experiment.

Scientific parent: sealed UP-101C 0deaa954944277a63e1b96c01d2769bba82b4b73.

## Question

UP-101C showed that the reverse-order trigger8 failure from UP-100C disappears when diagnostic hot-key reads are inserted, proving that two-bit-aging query state can change final retention. Which checkpoint read timing causes that change?

## Frozen memory mechanism

Use the exact UP-100C trigger8 mechanism:
- exact recall cap 16;
- two-bit-aging replacement;
- 1024-bit two-hash seen-once filter;
- generation interval 32;
- trigger after 8 successful repeat admissions;
- post-trigger counter 8;
- no query-derived admission labels;
- no future oracle.

## Frozen workload

- 32 original keys;
- all 16 even keys are hot;
- initial orders:
  1. ascending
  2. reverse
  3. evens_then_odds
- one second write to every hot key in ascending hot-key order;
- 12,288 unique one-shot churn writes;
- retain the ordinary UP-100C hot-query cadence after every four promotion writes and every four churn writes;
- value vocabulary 32;
- 64 episodes per cell;
- seeds 193M and 194M.

## Frozen read-timing arms

Every arm uses trigger8. The only difference is an additional full hot-key query sweep:

1. no_extra_reads
   - exact UP-100C trigger8 behavior.

2. after_initial
   - one extra full hot-key query sweep immediately after the initial 32-key load.

3. after_promotion
   - one extra full hot-key query sweep immediately after all second writes and before churn.

4. both
   - both extra query sweeps.

The extra reads are ordinary memory queries and therefore may update two-bit-aging state. They do not change the seen-once filter, trigger counter, capacity, or values.

## Metrics

For each initial-order x read-timing arm:
- resident hot-key count immediately before churn, measured without a query;
- trigger-fired episode count;
- mean trigger index;
- false-positive churn admissions;
- final hot query accuracy;
- final hot-set exact accuracy;
- final recall entries used.

## Interpretation

A specific read timing that flips reverse-order failure to success isolates the aging-state dependency. Evens-first remains a trigger-blind control. This experiment is diagnostic; extra reads are not proposed as a production memory policy.

## Bounds

No trigger change, no capacity expansion, no filter change, no adaptive reads, no result-informed retry, no live activation, no production authority.
