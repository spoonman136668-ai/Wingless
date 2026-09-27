# Wingless UP-178C — warning-stage action ablation

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-177C 5dd0dd36b223926006a798baeb069c9660ff529f.

## Question

Does two-action prevention require cumulative reinforcement from both warning-triggered actions, or is one action timed from the later warning stage sufficient?

## Frozen environment

- no_refresh pressure
- cadences 2 and 4
- four cohorts
- all initial hand positions 0..15
- 64-write horizon

## Frozen warning timing

Warning iff current adversarial-shield horizon <= cadence.
Every action waits exactly one monitoring interval after its triggering warning-positive interval.

## Preregistered arms

- first_only: targeted action after the first warning-positive interval; maximum one action.
- second_warning_only: skip action after the first warning-positive interval; targeted action after the second warning-positive interval; maximum one action.
- both: targeted actions after the first and next later warning-positive intervals; maximum two actions.
- sham_both: same two-stage timing as both, but non-endangered target.

## Measurements

Per cadence × arm:
- actions
- eventual losses
- prevented losses
- prevented losses per action
- mean loss delay among treated failures

## Interpretation

If second_warning_only reproduces most two-action prevention, later timing dominates. If it does not, cumulative reinforcement across two warning-action cycles is necessary.

## Bounds

Counterfactual only. No adaptive stage selection, delay, action budget, threshold, policy input, capacity change, live activation, or production authority.
