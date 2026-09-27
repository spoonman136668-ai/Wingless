# Wingless UP-215C — native tail-action selectivity

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-214C 1bbd8fc636474781a998a65015865e2be9ed7d28.

## Question

Does the frozen native near-risk trigger for action 8 rescue cap-7 failures more often than it causes an eighth action on arms that cap 7 would already survive?

## Frozen environment

Policy pairs:
- no_refresh / hostile_shield
- hostile_shield / no_refresh
- alternating_shield / fixed_offset_refresh
- fixed_offset_refresh / alternating_shield

Phase advances:
- 0,4 writes

Horizons:
- 66,68,70,72 writes

Population:
- heldout contiguous quartets
- initial hands 0..15

## Frozen correction

Actions 1-7:
- unchanged committed schedule from UP-214C.

Action 8 primary:
- after action 7, fire at first monitoring boundary where adversarial_shield_horizon <= 4.

Comparator:
- cap 7: no eighth action.

## Paired classifications

Per arm:
- tail_rescue: cap7 loses by horizon; native-tail survives.
- unnecessary_eighth: cap7 survives by horizon; native-tail fires action 8.
- necessary_eighth_attempt: cap7 loses; native-tail fires action 8.
- ineffective_eighth: cap7 loses; native-tail fires action 8 but still loses.
- no_eighth_survivor: cap7 survives; native-tail does not fire action 8.

## Measurements

Per schedule × phase × horizon:
- cap7 losses
- native-tail losses
- eighth actions
- all paired classifications
- tail_rescue / unnecessary_eighth ratio where defined.

## Interpretation

More rescues than unnecessary eighth actions supports the frozen native tail trigger as net-selective correction. A high unnecessary-action burden would bound its utility even if survival improves.

## Bounds

Counterfactual only. No threshold fitting, adaptive selection, future-policy input to trigger, action budget >8, live activation, or production authority.
