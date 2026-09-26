# Wingless UP-127C — recurrence timing under second-wave overflow

Status: preregistered scientific bounded-memory timing experiment.

Scientific parent: sealed UP-126C 1b5a64de5f58d08ae37f6d14e92c82f71179da62.

## Question

UP-126C showed that after the stale age-2 reservoir is exhausted, even one additional qualified history candidate can evict still-useful pending target evidence. With memory and the age-only replacement rule unchanged, can target recurrence timing alone protect the valid evidence?

## Frozen mechanism

Use the exact UP-126C mechanism:
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

Reproduce the exact UP-126C saturation and first 16-candidate overflow wave.

Then test second-wave qualified counts Q2 in {0,1,4} under two timing arms:

1. recur_before_second
   - valid targets recur twice each at the start of post-first-wave generation A;
   - Q2 candidates then receive two sightings each;
   - unique one-shot writes fill generation A;
   - generation B is 32 unique one-shot writes.

2. recur_after_second
   - Q2 candidates receive two sightings each at the start of generation A;
   - unique one-shot writes fill generation A;
   - valid targets recur twice each in generation B;
   - unique one-shot writes fill generation B.

Both arms therefore execute the same two complete 32-write generations after the first overflow wave. Only target recurrence timing relative to the second qualified wave differs.

Seeds:
- 241000000;
- 242000000.

32 episodes per seed.

## Metrics

For every timing × Q2 cell:
- panic rate;
- first-wave and post-first-wave evictions;
- evictions by age;
- valid target admission/final accuracy;
- hot accuracy;
- target16 accuracy/exactness;
- nonpersistent admission/retention;
- one-shot false admissions;
- max current/history occupancy;
- recall entries used.

## Interpretation

If recurrence-before-second preserves valid targets while recurrence-after-second reproduces UP-126C loss, temporal recurrence is sufficient to move useful evidence into protected exact recall before destructive history pressure arrives. If both arms fail similarly, the boundary cannot be managed by timing alone under the frozen mechanism.

No numeric success threshold is introduced.

## Bounds

No memory increase, no replacement-rule change, no semantic labels, no query priority, no future oracle, no adaptive horizon, no result-informed retry, no live activation, no production authority.
