# Wingless UP-222C — no-query critical action-9 trigger

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-221C 20f4fc3da1f7bce03df38942d8072b3527506e8b.

## Question

Can the cleaner no-query critical horizon replace adversarial critical as the action-9 trigger, preserving rescue while reducing unnecessary ninth actions?

## Frozen environment

Policy sequence:
- alternating_shield / fixed_offset_refresh

Phase advance:
- 8 writes

Horizons:
- 74,76,78,80 writes

Population:
- heldout contiguous quartets
- initial hands 0..15

Monitoring cadence:
- 2 writes

## Frozen actions 1-8

- actions 1-7 unchanged committed schedule
- action 8 fires at the first post-action-7 monitoring boundary where adversarial_shield_horizon <= 2

## Action-9 triggers compared

Adversarial critical:
- fresh later adversarial_shield_horizon <= 2

No-query critical:
- fresh later no_query_horizon <= 2

Action 9:
- cannot fire on the same monitoring boundary as action 8
- maximum total actions = 9

Comparator:
- cap 8, no action 9

No threshold fitting or cadence change.

## Measurements

Per horizon × trigger:
- cap8 losses
- cap9 losses
- ninth actions
- ninth rescues
- unnecessary ninth actions
- necessary ninth attempts
- ineffective ninth attempts
- rescue / unnecessary ratio where defined

## Interpretation

Equivalent rescue with fewer unnecessary ninth actions would support no-query critical as a more selective native trigger. Any lost rescue identifies a sensitivity cost. Neither signal is expected to recover the one UP-221C silent failure that emitted neither warning.

## Bounds

Counterfactual only. No threshold fitting, adaptive signal selection, future-policy input, action budget >9, schedule retuning, denser monitoring, live activation, or production authority.
