# Wingless UP-128C — pre-overflow consolidation dose

Status: preregistered scientific bounded-memory consolidation experiment.

Scientific parent: sealed UP-127C 60ec52189e8230e19cee3b4a56e60e67d936adf9.

## Question

UP-127C showed that recurring valid targets before a second overflow wave completely protects them under unchanged bounded memory. What is the minimum pre-overflow recurrence evidence required to consolidate those targets into protected exact recall?

## Frozen mechanism

Use the exact UP-127C mechanism:
- exact recall cap 16;
- 32-entry current table;
- 32-entry history table;
- 128-byte admission sidecar;
- maximum-history-age replacement only;
- lowest-index tie break;
- no semantic priority;
- no query priority;
- no future oracle.

## Frozen workload

Reproduce exact saturation and the exact first 16-candidate overflow wave.

Fix second-wave pressure at Q2=4 qualified nonpersistent candidates.

Before that second wave, present every valid target:
- dose_0: 0 additional sightings;
- dose_1: 1 additional sighting;
- dose_2: 2 additional sightings.

Then:
- present the four Q2 candidates twice each;
- fill the current 32-write generation with unique one-shot writes;
- execute one additional full 32-write one-shot generation;
- apply unchanged 12,288-write churn and final evaluation.

Every arm therefore has two complete post-first-wave generations; only target recurrence dose differs.

Seeds:
- 243000000;
- 244000000.

32 episodes per seed.

## Metrics

Per dose:
- panic rate;
- first-wave and post-first-wave evictions;
- evictions by age;
- valid target admission and final accuracy;
- hot accuracy;
- target16 accuracy and exactness;
- nonpersistent admission/retention;
- one-shot false admissions;
- max current/history occupancy;
- recall entries used.

## Interpretation

A sharp change between one and two sightings would identify the exact temporal evidence threshold needed for consolidation under the frozen admission mechanism. A graded result would indicate that protection is not simply tied to the two-sighting qualification transition.

No numeric success threshold is introduced.

## Bounds

No memory increase, no replacement-rule change, no semantic labels, no query priority, no future oracle, no adaptive horizon, no result-informed retry, no live activation, no production authority.
