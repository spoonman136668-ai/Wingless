# Wingless UP-230C — action-9 hysteresis

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-229C f0c4678f08038b23942d5e6fe4d942378723f7e4.

## Question

Can requiring two consecutive native critical warnings before action 9 reduce unnecessary tail intervention on interleaved cohorts without materially reducing rescue?

## Frozen environment

Heldout interleaved cohorts:
- 0,5,10,15
- 1,6,11,12
- 2,7,8,13
- 3,4,9,14

Policy pairs:
- no_refresh / hostile_shield
- hostile_shield / no_refresh
- alternating_shield / fixed_offset_refresh
- fixed_offset_refresh / alternating_shield

Phase advances:
- 8 writes
- 12 writes

Horizon:
- 80 writes

Monitoring cadence:
- 2 writes

Population:
- initial hands 0..15

Eight condition cells.

## Frozen actions 1-8

Exactly the UP-229C combined policy:
- first action at initial adversarial critical warning;
- actions 2-4 on committed 8-write spacing;
- action 5 may advance before its due boundary only when adversarial_shield_horizon <= 2;
- actions 6-7 retain original committed due boundaries;
- no schedule compression;
- action 8 at first later adversarial critical warning after action 7.

## Action-9 arms

1. none:
- no action 9;
- maximum 8 actions.

2. single_critical:
- current accepted action-9 trigger;
- after action 8, fire at first later monitoring boundary where adversarial_shield_horizon <= 2.

3. two_consecutive_critical:
- after action 8, maintain a critical-warning streak across monitoring boundaries;
- streak increments only when adversarial_shield_horizon <= 2;
- noncritical boundary resets streak to zero;
- action 9 fires when streak reaches 2;
- action 9 cannot fire on the action-8 boundary.

Maximum total actions in action-9 arms remains 9.

## Measurements

Per condition cell:
- cap8 losses
- single-critical losses
- hysteresis losses
- ninth actions under each action-9 arm
- rescued cap8 failures
- unnecessary ninth actions on cap8 survivors
- ineffective ninth attempts on cap8 failures
- rescue / unnecessary ratio where defined

## Interpretation

A persistence gate that preserves most rescue while reducing unnecessary ninth actions would support hysteresis as an efficiency mechanism. A large rescue loss would show the second warning arrives too late.

## Bounds

Counterfactual only. No threshold change, no action budget >9, no new action type, no cohort-specific tuning, no cadence change, no future-policy input, no live activation, no production authority.
