# Wingless UP-227C — action-5 retiming transfer

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-226C 16da4b7a1cf7e4c12dd63a8b84d423973ab0c52a.

## Question

Does the frozen native critical-warning retiming rule selectively improve outcomes across heldout policy sequences and phase advances without increasing the eight-action budget?

## Frozen rule

Actions 1-4:
- unchanged committed schedule.

Action 5:
- original committed due boundary remains fixed;
- after action 4, if adversarial_shield_horizon <= 2 appears before the due boundary, fire the already-budgeted action 5 immediately;
- action 5 may advance once.

Actions 6-7:
- remain on their original committed due boundaries;
- no schedule compression after an early action 5.

Action 8:
- unchanged critical trigger after action 7.

Maximum actions:
- 8.

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
- comparator losses
- primary losses
- advanced-action5 arms
- rescued losses
- unnecessary advances
- necessary advance attempts
- ineffective advance attempts
- total advance writes

## Interpretation

Transfer with positive rescue and low unnecessary retiming would support bounded timing refinement beyond the single localized C226 arm. Failure in specific policy/phase cells would identify the retiming rule's generalization boundary.

## Bounds

Counterfactual only. No new actions, no action budget >8, no threshold fitting, no cadence change, no post-retiming schedule compression, no future-policy input, no live activation, no production authority.
