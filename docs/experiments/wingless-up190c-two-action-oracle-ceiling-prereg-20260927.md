# Wingless UP-190C — two-action targeted-refresh oracle ceiling

Status: preregistered scientific counterfactual diagnostic.

Scientific parent: sealed UP-189C 5a8f28303e12b2adaee3d588b0551211b87b08f5.

## Question

When the frozen correction rule fails to transfer across external pressure policies, is the problem poor timing or insufficient corrective budget?

## Frozen policies

- no_refresh
- alternating_shield
- fixed_offset_refresh
- hostile_shield

## Frozen cadences

- 2 writes
- 4 writes

## Frozen population

- four cohorts
- all initial hand positions 0..15

## Frozen action class

- targeted refresh of the endangered resident only
- maximum two actions
- actions may occur only at monitoring-interval starts
- no action changes the external pressure policy

## Offline oracle search

For each arm, policy, and cadence:
- evaluate the no-action baseline;
- evaluate every single targeted-refresh monitoring time;
- evaluate every ordered pair of distinct monitoring times t1 < t2;
- choose the schedule with the latest loss step;
- if any allowed schedule survives through write 64, classify the arm as oracle-preventable.

This search is an offline diagnostic upper bound only. It is not available to Wingless and is not an activation policy.

## Frozen comparator

Also evaluate the accepted UP-189C frozen two-stage correction rule unchanged.

## Measurements

Per policy × cadence:
- arms
- baseline losses
- frozen-rule prevented losses
- oracle-preventable losses
- oracle-preventable but frozen-rule missed
- arms not preventable by any <=2-action monitoring-grid schedule
- mean best-oracle loss-step extension over baseline
- maximum best-oracle loss-step extension

## Interpretation

If the oracle prevents substantially more arms than the frozen rule, the transfer failure is primarily timing/trigger placement. If the oracle also cannot prevent most arms, the two-action corrective budget is insufficient under that policy.

## Bounds

Counterfactual diagnostic only. No policy-specific rule activation, no threshold tuning, no action-budget increase, no arbitrary-write action timing, no future information exposed to Wingless, no live activation, or production authority.
