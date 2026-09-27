# Wingless UP-196C — hostile minimal-schedule heldout cohort transfer

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-195C 9fa909c1cc8f7503daea221e395e9d563b2a11a1.

## Question

Does the minimal full-rescue hostile schedule discovered on the original strided cohorts transfer unchanged to a distinct cohort topology?

## Frozen policy
- hostile_shield

## Frozen monitoring cadence
- 2 writes

## Frozen primary schedule
- first native critical warning triggers targeted refresh
- post-warning spacing: every 4 monitoring intervals = every 8 writes
- action cap: 8
- target: endangered resident

No retuning.

## Heldout cohort topology

Discovery cohorts were strided modulo-4 quartets:
- 0,4,8,12
- 1,5,9,13
- 2,6,10,14
- 3,7,11,15

Heldout evaluation cohorts are contiguous quartets:
- 0,1,2,3
- 4,5,6,7
- 8,9,10,11
- 12,13,14,15

All initial hand positions 0..15 are evaluated.

## Frozen comparators

1. paired no-action hostile baseline.
2. dense commitment comparator:
   - cadence 2
   - spacing 1 interval = every 2 writes
   - cap 16
   - same first critical trigger and target.

## Measurements

- valid arms
- baseline losses
- primary prevented losses
- primary actions
- primary accelerated losses
- dense comparator prevented losses
- dense comparator actions
- mean/max loss extension for failures

## Interpretation

Strong heldout transfer supports the claim that temporal coverage, not discovery-cohort idiosyncrasy, governs hostile rescue. Failure would bound the schedule to the original cohort structure.

## Bounds

Counterfactual only. No cohort-specific tuning, threshold change, adaptive spacing/cap, new action type, future-policy input to warning, live activation, or production authority.
