# Wingless UP-180C — stage-specific delay matrix

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-179C 88fe5debf5424c205e0274bee40a9c4b61ff1eaa.

## Question

Which stage of the two-warning corrective sequence benefits from the one-monitoring-interval delay?

## Frozen environment

- no_refresh pressure
- cadences 2 and 4
- four cohorts
- all initial hand positions 0..15
- maximum two targeted actions
- action two requires a fresh later warning

## Frozen stage-delay matrix

Evaluate all four combinations:
- stage1 delay 0, stage2 delay 0
- stage1 delay 0, stage2 delay 1
- stage1 delay 1, stage2 delay 0
- stage1 delay 1, stage2 delay 1

Delay 0 means act in the warning-positive interval before that interval's writes.
Delay 1 means act at the next monitoring interval.

Targets are always the endangered resident.

## Measurements

Per cadence × delay pair:
- actions at each stage
- eventual losses
- prevented losses
- prevented losses per action
- mean loss delay among treated failures

## Interpretation

The matrix localizes whether the timing benefit comes from delaying the bridging first action, the survival second action, or their combination.

## Bounds

Counterfactual only. No adaptive delay, target, action budget, warning threshold, policy input, capacity change, live activation, or production authority.
