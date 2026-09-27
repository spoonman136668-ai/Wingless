# Wingless UP-226C — critical-advanced action 5

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-225C 9f0985eecb8d0580f28d5399e930d37d0cc3cd1a.

## Question

Can a fresh native critical warning after action 4 advance only the already-budgeted fifth action soon enough to rescue the localized mid-sequence failure, without increasing total actions or disturbing the remainder of the committed schedule?

## Frozen environment

Policy sequence:
- alternating_shield / fixed_offset_refresh

Phase advance:
- 8 writes

Horizon:
- 80 writes

Population:
- heldout contiguous quartets
- initial hands 0..15

Monitoring cadence:
- 2 writes

## Frozen comparator

Accepted UP-225C cap-8 schedule:
- first action at initial adversarial critical warning;
- actions 2-7 on the committed 4-monitoring-interval / 8-write spacing;
- action 8 at first later adversarial critical warning after action 7;
- maximum total actions = 8.

## Primary intervention

Actions 1-4:
- identical to comparator.

Action 5:
- retain its original committed due boundary;
- after action 4, if a fresh adversarial_shield_horizon <= 2 appears before that due boundary, fire action 5 immediately at that monitoring boundary;
- action 5 may advance only once.

Actions 6-7:
- retain the original committed due boundaries they would have had under the comparator;
- an early action 5 does not compress or shift those later due boundaries.

Action 8:
- unchanged native critical trigger after action 7.

Maximum total actions remains 8.

## Measurements

Across all arms:
- comparator losses
- primary losses
- arms with advanced action 5
- rescued losses
- unnecessary advances on comparator survivors
- necessary advance attempts
- ineffective advance attempts
- action-5 advance amount in writes

Also preserve per-arm action steps for comparator and primary.

## Interpretation

Rescuing the localized action-4 failure without new actions would support bounded timing refinement: native danger can re-time an already-authorized correction rather than increase the resource budget.

## Bounds

Counterfactual only. No new action, no action budget >8, no threshold fitting, no cadence change, no schedule compression after action 5, no future-policy input, no live activation, no production authority.
