# Wingless UP-176C — second-action gate ablation

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-175C f878647f1722950e529d43b5189f5d665091964e.

## Question

Is cross-policy failure of the delayed two-shot rule caused by requiring a new warning before spending the bounded second action?

## Frozen first action

Reuse UP-175C:
- warning iff current adversarial-shield horizon <= cadence
- wait exactly one monitoring interval
- perform action one
- maximum two actions per arm

## Frozen second-action arms

- rewarning_gated: action two requires a later warning-positive interval, then waits one monitoring interval, exactly as UP-175C
- forced_second: after action one executes, action two executes exactly one monitoring interval later without requiring another warning

Both remain bounded to two actions maximum.

Target controls:
- targeted_refresh: endangered resident
- sham_refresh: non-endangered resident

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

Per policy × cadence × second-action gate × target:
- arms
- baseline losses
- actions
- treated losses
- prevented losses
- prevented losses per action
- mean loss delay among treated failures

## Interpretation

If forced_second restores prevention where rewarning_gated does not, the warning-repeat gate is too restrictive. If both fail, policy dynamics—not warning recurrence—limit the bounded correction.

## Bounds

Counterfactual only. Exactly the same first-action trigger and one-interval delay. Maximum two actions. No adaptive gate, delay, action budget, threshold, future-policy input, capacity change, live activation, or production authority.
