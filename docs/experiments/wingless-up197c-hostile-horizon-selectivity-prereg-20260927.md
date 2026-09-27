# Wingless UP-197C — hostile sparse-schedule horizon selectivity

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-196C dd33632854139ef0ff360ecbcbc0802e3d059977.

## Question

On the heldout cohort topology, does the frozen sparse hostile-rescue schedule prevent near-horizon failures more often than it causes genuinely unnecessary intervention?

## Frozen environment

Policy:
- hostile_shield

Heldout cohorts:
- contiguous quartets 0,1,2,3 / 4,5,6,7 / 8,9,10,11 / 12,13,14,15

All initial hand positions:
- 0..15

Monitoring cadence:
- 2 writes

## Frozen correction rule

- first native critical warning triggers targeted refresh of endangered resident
- subsequent refreshes every 4 monitoring intervals = every 8 writes
- maximum 8 actions

No retuning.

## Frozen evaluation horizons

- 8,12,16,20,24,28,32,36,40,44,48,52,56 writes

Each horizon is evaluated independently from the same initial state.

## Frozen grace window

- 8 writes, equal to one sparse correction interval.

Definitions:
- baseline failure: no-action loss <= evaluation horizon
- prevented failure: baseline failure but treated arm survives beyond horizon
- apparent unnecessary action: baseline survives horizon but correction fires before horizon
- near-future action: apparent unnecessary arm whose baseline fails within the next 8 writes
- true unnecessary action: apparent unnecessary arm whose baseline survives beyond horizon + 8 writes

## Measurements

Per horizon:
- arms
- baseline failures and survivors
- treated failures
- prevented failures
- total actions
- apparent unnecessary-action arms
- near-future action arms
- true unnecessary-action arms
- prevented / true-unnecessary ratio where defined
- true-unnecessary fraction among baseline arms surviving beyond horizon + grace

## Interpretation

Positive prevention with few or no true unnecessary actions supports a bounded, selective correction rule rather than indiscriminate maintenance.

## Bounds

Counterfactual only. No threshold, schedule, spacing, cap, cohort, grace, or horizon tuning after results; no new action type; no live activation or production authority.
