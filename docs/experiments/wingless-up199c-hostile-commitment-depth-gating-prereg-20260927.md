# Wingless UP-199C — hostile commitment-depth gating

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-198C c8eb2c79ca5fe09371799e9c29031e3051b7c11b.

## Question

How many scheduled hostile-rescue actions must remain unconditionally committed before fresh native warning can safely gate the remaining action opportunities?

## Frozen environment

Policy:
- hostile_shield

Heldout cohorts:
- contiguous quartets 0,1,2,3 / 4,5,6,7 / 8,9,10,11 / 12,13,14,15

Initial hand positions:
- 0..15

Monitoring cadence:
- 2 writes

Spacing:
- one scheduled opportunity every 4 monitoring intervals = every 8 writes

Maximum actions:
- 8

## Frozen commitment depths

Total action counts that may fire without a new warning after initial trigger:
- 1
- 2
- 4
- 8

Interpretation:
- depth 1: only the initial critical-warning action is committed; every later opportunity requires a fresh critical warning. This reproduces UP-198C gating.
- depth 2: first two total actions are committed.
- depth 4: first four total actions are committed.
- depth 8: all actions are committed up to the cap. This reproduces the successful comparator.

After the committed depth is reached:
- each remaining scheduled opportunity acts only if adversarial_shield_horizon <= cadence;
- a skipped opportunity does not shift the future schedule.

## Frozen horizons
- 16,24,32,40,48,56 writes

## Measurements

Per commitment depth × horizon:
- baseline failures
- treated failures
- prevented failures
- total actions
- arms with action
- accelerated failures
- skipped opportunities

## Interpretation

The smallest depth preserving full rescue identifies how much open-loop commitment is necessary before native warning can safely resume control. If only depth 8 preserves rescue, the tested native warning is unsuitable for maintenance-efficiency gating.

## Bounds

Counterfactual only. No post-result depth selection for this experiment, no threshold/cadence/spacing/cap/target/cohort/horizon tuning, no new action type, no live activation or production authority.
