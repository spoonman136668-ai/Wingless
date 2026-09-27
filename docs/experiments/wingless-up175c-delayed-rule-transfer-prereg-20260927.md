# Wingless UP-175C — delayed-rule transfer

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-174C 53b69b7677df06f83b1e701d0dbe8b1af53ee7ea.

## Question

Does the one-monitoring-interval delayed two-shot rule that improved prevention under no_refresh transfer across other future-pressure policies and monitoring cadences?

## Frozen state space

Reuse the expanded 64-arm state space:
- four cohorts
- all initial hand positions 0..15
- evaluate every successful frozen trigger-state construction

## Frozen warning/action rule

For each cadence:
- warning iff current adversarial-shield horizon <= cadence
- after a warning, wait exactly one monitoring interval
- then refresh the endangered resident
- maximum two actions per arm
- no more than one action per monitoring interval
- second action requires a later warning-positive interval and uses the same one-interval delay

Sham control uses identical delayed warning timing but refreshes a non-endangered resident.

## Frozen policies and cadences

Policies:
- no_refresh
- alternating_shield
- fixed_offset_refresh
- hostile_shield

Cadences:
- 2
- 4

## Measurements

Per policy × cadence × intervention:
- arms
- baseline losses
- actions
- treated losses
- prevented losses
- prevented losses per action
- mean delay among treated failures

## Interpretation

Prevention beyond no_refresh/cadence-2 would demonstrate transfer of the delayed bounded rule. Delay-only benefit without prevention still indicates useful lead-time extension. Sham equivalence would undermine target specificity.

## Bounds

Counterfactual only. Delay fixed at one monitoring interval. Maximum two actions. No adaptive delay, action-budget tuning, threshold tuning, future-policy schedule input to warning, capacity change, live activation, or production authority.
