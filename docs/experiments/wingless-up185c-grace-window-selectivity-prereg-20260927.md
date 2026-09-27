# Wingless UP-185C — grace-window selectivity

Status: preregistered scientific counterfactual selectivity diagnostic.

Scientific parent: sealed UP-184C 0f9ca8328048efb421b7d2bb2f06ef02b842e47c.

## Question

After separating near-future baseline failures from genuinely long-lived survivors, how much true unnecessary intervention remains under the frozen correction rule?

## Frozen environment and correction rule

Unchanged:
- no_refresh pressure
- four cohorts
- all initial hand positions 0..15
- cadences 2 and 4
- stage one targeted immediately at first warning
- stage two requires a fresh warning and acts one monitoring interval later
- maximum two actions

## Frozen evaluation horizons

Every integer horizon from 16 through 29 writes inclusive.

## Frozen grace windows

- 0 writes
- 2 writes
- 4 writes
- 8 writes

For an arm that survives the evaluation horizon but receives an action before that horizon:
- near_future_action_arm: no-action baseline fails on or before horizon + grace;
- true_unnecessary_action_arm: no-action baseline survives beyond horizon + grace.

Grace 0 reproduces the original horizon-only definition.

## Measurements

Per cadence × horizon × grace:
- baseline failures by horizon;
- prevented failures by horizon;
- apparent unnecessary-action arms at the horizon;
- near-future action arms reclassified by the grace window;
- true unnecessary-action arms;
- prevented / true-unnecessary ratio where defined;
- true-unnecessary fraction among baseline arms surviving beyond horizon + grace.

## Interpretation

If modest fixed grace windows remove most or all unnecessary intervention, the rule is primarily early rather than nonspecific.

## Bounds

Diagnostic only. No warning-threshold tuning, adaptive grace selection, horizon selection after results, action-budget change, policy input, capacity change, live activation, or production authority.
