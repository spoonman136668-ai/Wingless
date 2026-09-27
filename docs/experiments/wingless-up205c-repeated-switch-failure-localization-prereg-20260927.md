# Wingless UP-205C — repeated-switch failure localization

Status: preregistered scientific counterfactual diagnostic.

Scientific parent: sealed UP-204C 922d8cdf5a6a06a5298e0e01472b5362c742ba05.

## Question

Why does one arm escape rescue under the fixed_offset_refresh <-> alternating_shield repeated-switch schedule?

## Frozen population, schedules, and rule

Exactly unchanged from UP-204C:
- heldout contiguous quartets
- initial hands 0..15
- four repeated policy schedules alternating every 8 writes
- cadence 2
- critical native trigger
- spacing 4 monitoring intervals = 8 writes
- action cap 8
- targeted endangered-resident refresh
- horizon 56

No correction behavior changes.

## Per-arm diagnostics

For every valid arm:
- schedule
- cohort index
- initial hand
- endangered key
- no-action baseline loss step
- treated loss step
- trigger onset monitoring write
- action count
- exact action monitoring writes
- last action write
- whether action cap 8 was reached
- loss delay relative to baseline
- whether treatment survives horizon

## Aggregate diagnostics per schedule

- arms
- treated failures
- treated failures that reached action cap
- mean action count among treated failures
- mean writes from last action to treated failure
- minimum/maximum treated failure step

## Interpretation

Cap exhaustion before the escape would support a resource-coverage boundary. Failure before cap exhaustion would instead point to policy-switch timing or target dynamics. This experiment does not authorize either repair.

## Bounds

Diagnostic only. No rule, threshold, cadence, spacing, cap, target, cohort, policy schedule, or horizon changes; no adaptive selection; no live activation or production authority.
