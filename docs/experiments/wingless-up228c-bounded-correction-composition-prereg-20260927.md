# Wingless UP-228C — bounded correction composition

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-227C 128e1909e0b4f460628de6806bca6f6da4ac4f04.

## Question

Do the two independently qualified bounded correction mechanisms compose cleanly: mid-sequence action-5 retiming and fresh-critical action-9 tail repair?

## Frozen mechanisms

### Mid-sequence retiming
- actions 1-4 unchanged committed schedule;
- action 5 may advance before its original due boundary only when adversarial_shield_horizon <= 2;
- actions 6-7 retain the original due boundaries derived from action 4;
- no schedule compression;
- this mechanism does not add an action.

### Tail repair
- action 8 remains the first post-action-7 adversarial critical action;
- action 9 may fire only on a later monitoring boundary with a fresh adversarial_shield_horizon <= 2;
- action 9 cannot fire on the same boundary as action 8;
- maximum total actions = 9.

The tail trigger is the accepted adversarial-critical trigger from UP-222C.

## Frozen policy arms

1. baseline:
- no action-5 retiming
- maximum 8 actions

2. retime_only:
- action-5 retiming enabled
- maximum 8 actions

3. tail9_only:
- no action-5 retiming
- fresh-critical action 9 enabled
- maximum 9 actions

4. combined:
- action-5 retiming enabled
- fresh-critical action 9 enabled
- maximum 9 actions

No thresholds or schedules differ between arms except the presence/absence of the two frozen mechanisms.

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

Population:
- heldout contiguous quartets
- initial hands 0..15

Monitoring cadence:
- 2 writes

16 condition cells.

## Measurements

Per condition cell:
- losses under all four policy arms;
- total actions under all four arms;
- action-5 advance events and advance writes;
- action-9 events;
- combined rescue beyond baseline;
- combined rescue beyond retime_only;
- combined rescue beyond tail9_only.

## Interpretation

If combined loss is never worse than either component and improves cells containing both mid-sequence and tail failure modes, the two corrections compose without destructive interaction. Any regression identifies an interaction boundary before iterative refinement is attempted.

## Bounds

Counterfactual only. No new trigger, no new action type, no action budget >9, no threshold fitting, no cadence change, no schedule compression, no future-policy input, no live activation, no production authority.
