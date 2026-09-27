# Wingless UP-204C — repeated policy-switch transfer

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-203C dd4edb69ae2127725643f00a4564e7426ade8159.

## Question

Does the frozen sparse rescue rule remain effective when external pressure repeatedly switches during the same trajectory?

## Frozen population
Heldout contiguous quartets:
- 0,1,2,3
- 4,5,6,7
- 8,9,10,11
- 12,13,14,15

Initial hand positions:
- 0..15

## Frozen repeated schedules

Policy alternates every 8 writes, starting with the first named policy:
1. no_refresh <-> hostile_shield
2. hostile_shield <-> no_refresh
3. alternating_shield <-> fixed_offset_refresh
4. fixed_offset_refresh <-> alternating_shield

Blocks:
- writes 1–8: policy A
- 9–16: policy B
- 17–24: policy A
- 25–32: policy B
- 33–40: policy A
- 41–48: policy B
- 49–56: policy A

The warning receives no future schedule information.

## Frozen rescue rule
- monitoring cadence: 2 writes
- first native critical warning: adversarial_shield_horizon <= cadence
- immediate targeted refresh of endangered resident
- subsequent targeted refreshes every 4 monitoring intervals = every 8 writes
- maximum 8 actions
- evaluation horizon 56 writes

## Comparator
Paired no-action baseline under the exact same repeated policy schedule.

## Measurements
Per schedule:
- arms
- baseline losses
- treated losses
- prevented losses
- accelerated losses
- actions
- arms with action
- mean loss-step extension among treated failures

## Interpretation
Full rescue would show the warning-plus-sparse-maintenance rule tolerates repeated nonstationarity without policy-specific logic. Failures localized by schedule direction would define the boundary.

## Bounds
Counterfactual only. No future-schedule input to warning, no policy-specific correction, no adaptive switching, no threshold/cadence/spacing/cap/target/cohort/horizon tuning, no live activation or production authority.
