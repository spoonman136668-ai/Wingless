# Wingless UP-229C — interleaved-cohort transfer

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-228C f6e4274d318708d735e7e5544a79d4e8dbb26b19.

## Question

Does the frozen combined bounded-correction policy transfer to a disjoint interleaved cohort topology without retuning?

## Frozen combined policy

Exactly the accepted UP-228C combined arm:

- first action at initial adversarial critical warning;
- actions 2-4 on committed 8-write spacing;
- action 5 may advance before its original due boundary only when adversarial_shield_horizon <= 2;
- actions 6-7 retain their original committed due boundaries;
- no schedule compression after an advanced action 5;
- action 8 at the first later adversarial critical warning after action 7;
- action 9 only at a fresh later adversarial critical warning;
- action 9 cannot fire on the same monitoring boundary as action 8;
- maximum total actions = 9.

No thresholds or timings change.

## Heldout cohort topology

Four interleaved quartets partition keys 0..15:

- cohort A: 0,5,10,15
- cohort B: 1,6,11,12
- cohort C: 2,7,8,13
- cohort D: 3,4,9,14

These differ from the contiguous quartets used to qualify the policy.

Initial hands:
- 0..15

## Preregistered transfer grid

Policy pairs:
- no_refresh / hostile_shield
- hostile_shield / no_refresh
- alternating_shield / fixed_offset_refresh
- fixed_offset_refresh / alternating_shield

Phase advances:
- 8 writes
- 12 writes

Horizons:
- 74 writes
- 80 writes

Monitoring cadence:
- 2 writes

16 condition cells.

## Comparator

Frozen cap-8 baseline schedule with no action-5 retiming and no action 9.

## Measurements

Per condition cell:
- arms
- baseline losses
- combined losses
- baseline actions
- combined actions
- combined action-5 advance events
- combined action-5 advance writes
- combined ninth actions
- rescued baseline failures
- intervention-bearing baseline survivors

## Interpretation

Low combined loss with selective intervention on this disjoint topology supports generalization of the bounded correction mechanism beyond the original cohort geometry. Regressions localize a topology-transfer boundary.

## Bounds

Counterfactual only. No cohort-specific tuning, no new trigger, no new action type, no action budget >9, no threshold fitting, no cadence change, no schedule compression, no future-policy input, no live activation, no production authority.
