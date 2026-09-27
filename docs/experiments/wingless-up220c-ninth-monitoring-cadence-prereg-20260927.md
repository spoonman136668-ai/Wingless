# Wingless UP-220C — action-9 monitoring cadence ablation

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-219C c56186d088c0a58ff44d9e1e5df94c5034aea60f.

## Question

Is the persistent cap-9 residual failure caused by the normal two-write monitoring cadence missing a short-lived critical-warning opportunity after action 8?

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
- action 8 fires at first post-action-7 normal monitoring boundary with adversarial_shield_horizon <= 2
- actions 1-8 retain the normal two-write monitoring cadence

## Action-9 monitoring modes

1. cadence2
- evaluate fresh critical warning only at subsequent normal monitoring boundaries
- reproduces UP-219C critical action-9 rule

2. every_write
- after action 8, evaluate adversarial_shield_horizon <= 2 immediately before each subsequent write
- action 9 cannot fire on the same monitoring boundary as action 8
- maximum total actions remains 9

Threshold is unchanged.

## Comparator

Cap 8, no action 9.

## Measurements

Per horizon × monitoring mode:
- cap8 losses
- cap9 losses
- ninth actions
- ninth rescues
- unnecessary ninth actions
- necessary ninth attempts
- ineffective ninth attempts
- rescue / unnecessary ratio

## Interpretation

If every-write monitoring removes the persistent residual without a large unnecessary-action increase, the boundary is temporal sampling of a valid native warning. If it does not, the residual is not recoverable by denser observation of the same warning signal.

## Bounds

Counterfactual only. No threshold fitting, adaptive cadence selection, future-policy input, action budget >9, schedule retuning, live activation, or production authority.
