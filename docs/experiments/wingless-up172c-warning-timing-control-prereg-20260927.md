# Wingless UP-172C — warning timing control

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-171C a82baf0cb475c054621f289774d7d563dfc642e1.

## Question

Is the timing supplied by Wingless's native warning actually useful, or would the same two targeted refreshes work as well when spent on a fixed early schedule?

## Frozen state space

Reuse the expanded UP-171C state space:
- four cohorts;
- all initial hand positions 0..15;
- evaluate every arm where the frozen trigger-state construction succeeds.

## Frozen environment

- no_refresh pressure only;
- monitoring cadence 2;
- 64-write evaluation horizon.

This is the regime where UP-171C established a reproducible two-shot prevention effect.

## Frozen interventions

All targeted modes refresh the same endangered resident and use at most two actions.

- warning_targeted: reuse UP-170C/171C exactly; refresh at the first two warning-positive monitoring intervals.
- fixed_early_targeted: ignore warning state and refresh before writes 1 and 3.
- warning_sham: reuse the warning timing but spend actions on a non-endangered resident.
- no_action: baseline.

The fixed schedule is preregistered before this run and is not adapted per arm.

## Measurements

Per intervention:
- arms;
- actions;
- eventual losses;
- prevented losses versus no-action baseline;
- prevented losses per action;
- mean loss delay among arms that still fail.

Also report paired counts:
- warning_targeted prevents while fixed_early_targeted does not;
- fixed_early_targeted prevents while warning_targeted does not;
- both prevent;
- neither prevents.

## Interpretation

If warning-triggered targeting outperforms the equal-budget fixed schedule, Wingless's endogenous warning timing contributes causal value. If fixed-early matches or exceeds it, the corrective mechanism is useful but the warning timing is not yet necessary.

## Bounds

Counterfactual only. Maximum two actions per arm. No adaptive timing, action-budget tuning, threshold tuning, repeated maintenance beyond two actions, future-policy schedule input, capacity change, live activation, or production authority.
