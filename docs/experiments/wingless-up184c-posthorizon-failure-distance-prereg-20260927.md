# Wingless UP-184C — post-horizon failure distance

Status: preregistered scientific counterfactual selectivity diagnostic.

Scientific parent: sealed UP-183C 86cf3e5da3b7e3ab2bdcfabfff333d3f63d3e17c.

## Question

When the frozen correction rule acts on an arm that would have survived the chosen evaluation horizon, is that intervention a true false positive or an early warning of a failure shortly beyond the horizon?

## Frozen environment and correction rule

Unchanged from UP-183C:
- no_refresh pressure
- four cohorts
- all initial hand positions 0..15
- cadences 2 and 4
- stage one targeted immediately at first warning
- stage two requires a fresh warning and acts one monitoring interval later
- maximum two actions

## Frozen evaluation horizons

Every integer horizon from 16 through 29 writes inclusive.

## Frozen baseline follow-up

For every arm classified as an unnecessary-action arm at a horizon:
- continue using its already-frozen no-action baseline loss step through write 64;
- post-horizon failure distance = baseline loss step - evaluation horizon.

No intervention result is used to define this distance.

## Measurements

Per cadence × horizon:
- baseline survivors;
- unnecessary-action arms;
- unnecessary arms failing within +1, +2, +4, +8, and +16 writes;
- unnecessary arms still surviving beyond write 64;
- minimum, maximum, and mean post-horizon failure distance for unnecessary-action arms with observed baseline failure by 64.

## Interpretation

A concentration of unnecessary-action arms at short post-horizon distances means the rule is primarily firing early on near-future failures rather than producing arbitrary false positives.

## Bounds

Diagnostic only. No threshold tuning, horizon selection after results, action-budget change, policy input, capacity change, live activation, or production authority.
