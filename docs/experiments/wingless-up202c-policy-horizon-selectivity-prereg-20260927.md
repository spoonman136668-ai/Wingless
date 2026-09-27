# Wingless UP-202C — sparse-schedule cross-policy horizon selectivity

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-201C eb199a3866190a78f20279f7cbcf892393ebdc9c.

## Question

Does the frozen sparse rescue rule prevent near-horizon failures more often than it causes genuinely unnecessary intervention across all four external pressure policies?

## Frozen policies
- no_refresh
- alternating_shield
- fixed_offset_refresh
- hostile_shield

## Frozen population
Heldout contiguous quartets:
- 0,1,2,3
- 4,5,6,7
- 8,9,10,11
- 12,13,14,15

Initial hand positions:
- 0..15

## Frozen rule
- monitoring cadence 2 writes
- first native critical warning triggers targeted refresh
- subsequent targeted refreshes every 4 monitoring intervals = every 8 writes
- maximum 8 actions
- endangered resident is the only target

No policy-specific changes.

## Frozen horizons
- 8,12,16,20,24,28,32,40,48,56 writes

Each policy × horizon is evaluated independently from the same initial state.

## Frozen grace window
- 8 writes

Definitions:
- baseline failure: no-action loss <= horizon
- prevented failure: baseline failure but treated survives beyond horizon
- apparent unnecessary action: baseline survives horizon but correction fires
- near-future action: apparent unnecessary arm whose baseline fails by horizon + 8
- true unnecessary action: apparent unnecessary arm whose baseline survives beyond horizon + 8

## Measurements
Per policy × horizon:
- arms
- baseline failures / survivors
- treated failures
- prevented failures
- total actions
- apparent unnecessary-action arms
- near-future action arms
- true unnecessary-action arms
- baseline survivors beyond grace
- prevented / true-unnecessary ratio where defined
- true-unnecessary fraction among beyond-grace survivors

## Interpretation

Strong prevention with few or zero true-unnecessary actions across policies would support a policy-robust, selective bounded correction rule. Policy-specific false alarms or misses would localize the remaining generalization boundary.

## Bounds

Counterfactual only. No policy-specific threshold/cadence/spacing/cap/target/cohort/grace/horizon tuning; no adaptive selection; no live activation or production authority.
