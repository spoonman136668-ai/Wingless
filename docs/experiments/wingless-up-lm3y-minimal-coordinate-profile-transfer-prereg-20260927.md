# Wingless UP-LM3Y — minimal coordinate profile transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM3X b0bb842b995260571a01dfcb649e47b1dbbf0a53.

## Question

Does the minimal temporal resource coordinate identified by endpoint ablation transfer across distinct deadline structures on the broader resource grid?

## Frozen deadline profiles

- deferred_only
- layout_only
- hybrid_min

## Frozen pressure

Exactly four extra pre-window writes.

## Frozen resource grid

- budgets: 3,4,5,6,7
- action start rounds: 2,3,4,5
- throughput: 1,2,3,4
- rotations: 5 and 13
- permutations: identity, reverse, rotate2
- six global rounds
- earliest_deadline policy

80 resource configurations per deadline profile.

## Frozen coordinates

Prior rich:
- reachable_budget + capacity_by_round3 + capacity_by_round4

Minimal candidate:
- reachable_budget + capacity_by_round2 + capacity_by_round3 + capacity_by_round4

Full comparator:
- reachable_budget + capacity_by_round2 + capacity_by_round3 + capacity_by_round4 + capacity_by_round5

No coordinate changes after results.

## Measurements

Per deadline profile:
- global outcome range
- distinct outcome count

Per deadline profile × coordinate:
- groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread

## Interpretation

If the minimal candidate has zero spread wherever the full comparator has zero spread, round-5 capacity is genuinely redundant and the compressed scheduler state transfers across deadline structures.

## Bounds

Diagnostic only. No adaptive coordinate search, profile selection, resource tuning, pressure tuning, topology search, capacity change, extra training, live activation, or production authority.
