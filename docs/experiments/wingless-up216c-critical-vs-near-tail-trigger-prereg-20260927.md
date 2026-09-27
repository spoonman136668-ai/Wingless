# Wingless UP-216C — critical vs near tail trigger

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-215C a4a74a4ff3c45f76466236c2350bb3b6f5b3ef1e.

## Question

Can the already-defined native critical threshold reduce unnecessary action-8 use while retaining the rescue benefit of the current native near-risk threshold?

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
- unchanged committed schedule from UP-215C.

Action 8 triggers compared:

1. critical:
- fire at first monitoring boundary after action 7 where adversarial_shield_horizon <= 2.

2. near:
- fire at first monitoring boundary after action 7 where adversarial_shield_horizon <= 4.
- this reproduces UP-215C.

Comparator:
- cap 7: no eighth action.

No threshold fitting.

## Measurements

Per schedule × phase × horizon × trigger:
- cap7 losses
- treated losses
- eighth actions
- tail rescues
- unnecessary eighth actions
- necessary eighth attempts
- ineffective eighth attempts
- no-eighth survivors
- rescue / unnecessary ratio where defined

## Interpretation

A critical trigger that materially reduces unnecessary action while retaining most or all rescues would improve selectivity using a threshold already established independently in the native risk work. Loss of rescue would show the near-risk lead time is necessary for action 8.

## Bounds

Counterfactual only. No threshold fitting, adaptive selection, future-policy input, action budget >8, live activation, or production authority.
