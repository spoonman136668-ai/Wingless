# Wingless UP-127C — recurrence timing protection

Status: preregistered scientific bounded-memory identifiability experiment.

Scientific parent: sealed UP-126C 1b5a64de5f58d08ae37f6d14e92c82f71179da62.

## Question

UP-126C showed that age-only eviction is safe while genuinely stale age-2 evidence exists, but after that reservoir is exhausted a second four-candidate overflow wave destroys all four pending valid targets. With memory, evidence, and the replacement rule unchanged, is recurrence timing alone sufficient to protect those targets?

## Frozen mechanism

Use the exact UP-126C mechanism:
- exact recall cap 16;
- current table 32 entries;
- history table 32 entries;
- admission sidecar 128 bytes;
- history ages 0/1/2;
- maximum-age eviction;
- lowest-index tie-break;
- no semantic priority;
- no query-derived priority;
- no future oracle;
- no memory increase.

## Frozen sequence

Both arms reproduce the same UP-126C state:

1. establish 12 hot exact residents;
2. generation 1: 16 qualified nonpersistent set-A candidates;
3. generation 2: four valid targets + 12 qualified nonpersistent set-B candidates;
4. first overflow wave: 16 fresh qualified nonpersistent candidates, consuming all stale set-A age-2 slots;
5. second overflow wave: exactly four fresh qualified nonpersistent candidates;
6. four valid targets recur exactly twice total;
7. unchanged 12,288-write churn;
8. final evaluation.

The arms differ only in the position of step 5 relative to target recurrence:

- **before_overflow**: valid targets recur immediately after the first overflow wave, then the four-candidate second overflow wave occurs.
- **after_overflow**: four-candidate second overflow wave occurs first, then valid targets recur. This is the UP-126C destructive control.

Each target therefore receives exactly the same number of observations in both arms.

## Seeds and episodes

- seeds: 241000000 and 242000000;
- 32 episodes per seed;
- 64 episodes per arm.

## Metrics

Per arm:
- panic rate;
- first-wave and second-wave eviction counts;
- evictions by age;
- valid-target admission rate;
- valid-target final accuracy;
- 12-hot accuracy;
- target16 accuracy and exact rate;
- nonpersistent admission/retention;
- one-shot false admissions;
- maximum current/history occupancy;
- exact recall entries used.

## Interpretation

If before_overflow preserves the four valid targets while after_overflow reproduces their loss, recurrence timing provides sufficient causal evidence for the existing bounded mechanism to protect useful history. The distinction would arise from observed recurrence, not semantic labels or future knowledge.

If both arms fail, timing alone is insufficient under the current admission mechanism.

No numeric threshold is introduced after execution; exact measured differences are the result.

## Bounds

No memory increase, no replacement-rule change, no semantic labels, no query priority, no future oracle, no adaptive timing, no extra target sightings, no result-informed retry, no live activation, no production authority.
