# Wingless UP-LM4M — mid-window budget-revocation transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4L f598b8acf5c912d22c6f30ce53dbbe3714095157.

## Question

Does the compact resource law remain sufficient when some remaining action budget is revoked after scheduling has begun, and does an actual-schedule deployability coordinate recover any break?

## Frozen deadline profiles
- deferred_only
- layout_only
- hybrid_min

## Frozen budget revocations

Initial budgets:
- 4,5,6,7 actions

At a fixed global cut round:
- 2
- 3
- 4

the total action cap is permanently reduced by:
- 1 action
- 2 actions

If actions already spent exceed the reduced cap, they are not undone; no further actions may occur until the cap again exceeds spent actions, which never happens in this experiment.

## Frozen resource grid
- initial budgets 4,5,6,7
- action starts 2,3,4,5
- throughput 1,2,3,4
- six global rounds
- earliest_deadline scheduling

64 resource configurations per topology-condition cell.

## Frozen topology
- rotations 5 and 13
- permutations identity, reverse, rotate2

## Frozen coordinates

Static comparator:
- reachable_budget from original budget/start/throughput
- nominal throughput
- original coverage_end_round

Dynamic primary:
- actual reachable actions under the budget-cap schedule
- nominal throughput
- actual coverage_end_round under the budget-cap schedule

Cut round and revocation amount are fixed within each condition cell and are not grouping-key inputs.

## Measurements

Per deadline profile × cut round × revoke amount × rotation × permutation × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread

108 topology-condition cells; 216 coordinate summaries.

## Interpretation

If the static coordinate breaks while the dynamic coordinate remains exact, actual remaining action authority—not nominal initial budget—is the necessary resource state. If both remain exact, the compact law is robust to bounded budget loss. If dynamic also breaks, resource history beyond reachable amount and coverage is required.

## Bounds

Diagnostic only. No adaptive cut timing/amount, budget restoration, resource tuning, coordinate modification, topology/profile selection, extra training, live activation, or production authority.
