# Wingless UP-183C — dense horizon selectivity

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-182C 4526fdb17de1bbdd873e79968177f6416a79225e.

## Question

Where exactly does the frozen correction rule transition from highly selective intervention to broad intervention as the evaluation horizon approaches population-wide failure?

## Frozen environment and rule

Unchanged from UP-182C:
- no_refresh pressure
- four cohorts
- all initial hand positions 0..15
- cadences 2 and 4
- stage one targeted immediately at first warning
- stage two requires fresh warning and acts one monitoring interval later
- maximum two actions

## Frozen evaluation horizons

Every integer horizon from 16 through 31 writes inclusive.

## Measurements

Per cadence × horizon:
- baseline failures and survivors
- intervention actions
- arms with any action
- prevented failures
- unnecessary-action arms
- necessary-action arms
- no-action survivors
- prevented / unnecessary-action-arm ratio where defined
- unnecessary-action fraction among baseline survivors

## Bounds

Counterfactual only. No threshold tuning, adaptive horizon selection, action-budget change, policy input, capacity change, live activation, or production authority.
