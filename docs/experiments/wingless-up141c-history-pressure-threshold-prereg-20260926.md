# Wingless UP-141C — history-pressure threshold

Status: preregistered scientific bounded-memory consolidation experiment.

Scientific parent: sealed UP-140C 1ef8af85e7ee3d459e2147a3fb4611bffa78701e.

## Question

UP-140C showed that two 32-event windows with 16 qualified manipulation identities per window fill the 32-entry history table and evict an age-2 target, while lower qualified-entry densities preserve it. What is the exact insertion-pressure boundary?

## Frozen mechanism

Reuse the accepted age-eviction machine unchanged:
- exact recall cap 16;
- current/history tables 32/32;
- admission sidecar 128 bytes;
- maxAge = 2;
- maximum-history-age replacement;
- lowest-index tie break;
- no semantic/query/future priority.

Each arm uses one recurrence target.

## Frozen setup

Reuse UP-140C setup:
- expire original target history;
- rebuild fresh target history with a G3 same-window pair;
- enter manipulation with target history at age 0.

## Fixed manipulation

Execute exactly 64 unresolved candidate calls:
- 32 calls in window 1;
- 32 calls in window 2;
- window-2 identities are disjoint from window-1 identities;
- no manipulation identity can complete durable admission.

Within each 32-call window, create exactly q qualified current identities by presenting q keys twice. Fill the remaining 32-2q calls with one-shot unique keys.

q arms:
- q0
- q12
- q14
- q15
- q16

Thus all arms have identical unresolved-event count and candidate boundaries; only qualified-history insertion pressure changes.

Two targets:
- 100
- 103.

Seeds:
- 269000000
- 270000000

32 episodes per seed.

## Measurements

Per target × q:
- qualified identities per window;
- unresolved calls;
- target history presence/age after boundary 1;
- target history presence/age after boundary 2;
- history occupancy after each boundary;
- age0/age1/age2 evictions during manipulation;
- unexpected manipulation admissions;
- target admission rate/final accuracy;
- hot accuracy;
- max table occupancy;
- recall entries;
- panic rate.

## Interpretation

With target + q first-window history + q second-window history, q=15 predicts 31 total entries and survival; q=16 predicts 33 logical entries against capacity 32, forcing eviction of the age-2 target. The exact result is the boundary test.

No post-result threshold is introduced.

## Bounds

No memory increase, no replacement-rule change, no semantic labels, no query priority, no future oracle, no adaptive horizon, no result-informed retry, no live activation, no production authority.
