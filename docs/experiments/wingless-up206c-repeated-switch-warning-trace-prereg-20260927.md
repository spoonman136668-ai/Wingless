# Wingless UP-206C — repeated-switch pre-failure warning trace

Status: preregistered scientific counterfactual diagnostic.

Scientific parent: sealed UP-205C 5f797af968dca22f1007e5035e64bc868ecd7b3b.

## Question

Does the existing native critical-warning signal reappear before the isolated repeated-switch rescue failure, early enough to identify the spacing-phase miss without changing the rescue rule?

## Frozen environment and rule

Exactly unchanged from UP-204C / UP-205C:
- heldout contiguous quartets
- initial hands 0..15
- four repeated policy schedules alternating every 8 writes
- cadence 2
- critical threshold: adversarial_shield_horizon <= 2
- initial targeted refresh on first critical warning
- later scheduled targeted refresh every 4 monitoring intervals = 8 writes
- cap 8
- horizon 56

## Diagnostic observation

At every monitoring boundary, before any scheduled action at that boundary:
- record current adversarial_shield_horizon.

This observation does not change state.

Per arm record:
- baseline loss step
- treated loss step
- trigger onset
- action writes
- final monitoring boundary before treated loss
- native horizon at that final monitoring boundary
- last critical-warning monitoring boundary before treated loss
- lead in writes from last critical warning to treated loss
- count of post-trigger critical-warning boundaries
- survival

## Aggregate measurements per schedule
- treated failures
- failures with a critical warning at the final monitoring boundary before loss
- failures with any critical warning after trigger and before loss
- mean/min/max critical-warning lead for treated failures

## Interpretation

A critical warning immediately before the isolated loss would show the warning can detect the maintenance phase miss even though it cannot replace the committed maintenance schedule wholesale. Absence of a warning would rule out that simple emergency pull-forward mechanism.

## Bounds

Diagnostic only. No rescue action, timing, threshold, cadence, spacing, cap, target, cohort, policy schedule, or horizon changes; no adaptive selection; no live activation or production authority.
