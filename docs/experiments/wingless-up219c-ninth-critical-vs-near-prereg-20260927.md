# Wingless UP-219C — action-9 critical vs near warning

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-218C 0c3385e8c39eba77a10478cf5f4f9f448036687e.

## Question

On the one residual cap-9 boundary, does using the already-defined near-risk warning for action 9 eliminate the remaining warning-lead misses at acceptable action cost?

## Frozen residual environment

Policy sequence:
- alternating_shield / fixed_offset_refresh

Phase advance:
- 8 writes

Horizons:
- 74,76,78,80 writes

Population:
- heldout contiguous quartets
- initial hands 0..15

## Frozen actions 1-8

- actions 1-7 unchanged committed schedule
- action 8 fires at first post-action-7 monitoring boundary with adversarial_shield_horizon <= 2

## Action-9 triggers compared

Critical:
- fresh later adversarial_shield_horizon <= 2

Near:
- fresh later adversarial_shield_horizon <= 4

Action 9 cannot fire on the same monitoring boundary as action 8.
Maximum total actions = 9.

Comparator:
- cap 8, no action 9.

No threshold fitting.

## Measurements

Per horizon × trigger:
- cap8 losses
- cap9 losses
- ninth actions
- ninth rescues
- unnecessary ninth actions
- necessary ninth attempts
- ineffective ninth attempts
- rescue / unnecessary ratio

## Interpretation

If near warning removes the residual critical-trigger miss with a bounded increase in unnecessary action, it identifies a specific lead-time requirement for the second iterative tail repair.

## Bounds

Counterfactual only. No threshold fitting, adaptive selection, future-policy input, action budget >9, schedule retuning, live activation, or production authority.
