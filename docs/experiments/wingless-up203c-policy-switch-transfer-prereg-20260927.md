# Wingless UP-203C — sparse rescue under mid-run policy switches

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-202C d53d84720a7706eba18e2ebfdb08c103f28e060c.

## Question

Does the frozen sparse rescue rule remain effective when the external pressure policy changes during the same trajectory?

## Frozen population

Heldout contiguous cohorts:
- 0,1,2,3
- 4,5,6,7
- 8,9,10,11
- 12,13,14,15

Initial hand positions:
- 0..15

## Frozen policy schedules

All switch once at write 24:
1. no_refresh -> hostile_shield
2. hostile_shield -> no_refresh
3. alternating_shield -> fixed_offset_refresh
4. fixed_offset_refresh -> alternating_shield

Policy A applies for writes < 24.
Policy B applies for writes >= 24.

The warning does not receive the future policy schedule.

## Frozen rescue rule

Monitoring cadence:
- 2 writes

Native trigger:
- first monitoring boundary where adversarial_shield_horizon <= cadence

After trigger:
- immediate targeted refresh of endangered resident
- one targeted refresh every 4 monitoring intervals = every 8 writes
- maximum 8 actions

Evaluation horizon:
- 56 writes

## Comparator

Paired no-action baseline under the same policy-switch schedule.

## Measurements

Per policy schedule:
- arms
- baseline losses
- treated losses
- prevented losses
- accelerated losses
- total actions
- arms with action
- mean loss-step extension among treated failures

## Interpretation

Successful transfer would show that the current warning-plus-sparse-maintenance rule is robust to nonstationary external pressure, not only to fixed policy classes. Failure localized to a switch direction would identify where future policy context becomes necessary.

## Bounds

Counterfactual only. No policy-specific rule, no future-schedule input to warning, no threshold/cadence/spacing/cap/target/cohort/horizon tuning, no adaptive selection, no live activation or production authority.
