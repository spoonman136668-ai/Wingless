# Wingless UP-LM5H — residual collision trace

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM5G ef2c4e50db44007c089c14ae928799c8fdbf1a3f.

## Question

Which exact resource configurations still alias under the frozen LM5G state vector while producing different outcomes?

## Frozen residual boundary

Deadline profiles:
- layout_only
- hybrid_min

Budget reduction:
- 2

Throughput reductions:
- 1
- 2

Topology:
- rotations 5,13
- permutations identity, reverse, rotate2

This is exactly the 24-cell residual boundary from LM5F/LM5G.

## Frozen schedules

The six LM5E-LM5G heldout windows remain unchanged:
- budget 0/2 ; throughput 1/3
- budget 1/3 ; throughput 0/2
- budget 0/2 ; throughput 2/4
- budget 2/4 ; throughput 0/2
- budget 0/2 ; throughput 3/5
- budget 3/5 ; throughput 0/2

## Frozen resource grid

- budgets 4,5,6,7
- action starts 2,3,4,5
- throughput 2,3,4,5
- six global rounds
- earliest_deadline scheduling

384 points per residual condition cell.

## Frozen failed coordinate

Exactly the LM5G primary:
- final reachable actions
- nominal throughput
- final coverage end round
- actions before earlier restoration
- cumulative actions after earlier restoration
- actions before later restoration
- cumulative actions after later restoration
- earlier restore round
- later restore round
- early-restored resource identity

No new field is added.

## Measurements

For every nonzero-spread collision group:
- condition cell
- frozen coordinate key
- minimum and maximum failure outcome
- every member's:
  - nominal budget
  - action start
  - nominal throughput
  - cut/restore window
  - derived frozen-coordinate values
  - outcome

Also report collision-group count and member count per condition.

## Interpretation

The first scheduler fact that differs consistently across opposite-outcome members becomes the next preregistered candidate variable. This experiment itself does not select or add one.

## Bounds

Diagnostic only. No coordinate modification, adaptive feature selection, resource tuning, topology/profile expansion, capacity change, extra training, live activation, or production authority.
