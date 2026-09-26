# Wingless UP-101C — promotion checkpoint diagnostic

Status: preregistered scientific bounded-memory diagnostic.

Scientific parent: sealed UP-100C 648c353adc5eed3e0b1ec43e28846f71405843a6.

## Question

UP-100C showed that trigger8 is not initial-order general at full 16-item hot occupancy: reverse order fails despite zero churn false admissions, while evens-first never fires. Are those failures already present at the end of the second-write promotion phase, before one-shot churn begins?

## Frozen memory mechanism

Use the exact UP-100C system:
- exact recall cap 16;
- two-bit-aging replacement;
- 1024-bit two-hash seen-once filter;
- generation interval 32;
- post-trigger counter 8;
- trigger threshold 8 successful repeat admissions;
- no explicit phase label;
- no query-derived admission labels;
- no future oracle.

## Frozen workload

Use only the maximally diagnostic full-hot case:
- 32 original keys;
- hot set = all 16 even keys;
- initial orders:
  1. ascending
  2. reverse
  3. evens_then_odds
- each hot key receives exactly one second write in ascending hot-key order;
- then 12,288 unique one-shot churn writes;
- same hot-query cadence and value vocabulary as UP-100C;
- 64 episodes per cell;
- seeds 191M and 192M.

## Frozen arms

For every initial order:
1. continuous_control
2. trigger8

## Checkpoints and metrics

At three checkpoints:
1. after initial 32-key loading;
2. immediately after all 16 hot second writes, before the first churn write;
3. after the full churn stream.

Report:
- resident hot-key count;
- hot-key query accuracy;
- hot-set exact accuracy;
- resident cold-key count;
- repeat admissions;
- trigger-fired episode count;
- mean trigger index;
- churn false-positive admissions;
- final recall entries used.

Checkpoint queries are diagnostic reads only and do not modify admission state beyond the existing memory policy's ordinary query behavior. The same checkpoint queries are applied to every arm.

## Interpretation

If a failing cell is already incomplete before churn, the next mechanism must repair promotion completeness rather than churn filtering. If promotion is complete but final exactness fails, the defect remains in churn protection. If evens-first promotion is complete but trigger8 does not fire, the trigger signal is structurally blind to already-resident hot evidence.

## Bounds

No policy change, no trigger change, no capacity expansion, no filter-width/hash/interval change, no adaptive checkpoint logic, no result-informed retry, no live activation, no production authority.
