# Wingless UP-LM5I — deployment-onset transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM5H 253f24027071915a6f089459ed36c155edb7ea48.

## Question

Does a derived first-deliverable-round state variable resolve the exact 24-cell residual alias exposed by UP-LM5H?

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

This is exactly the 24-cell residual boundary from UP-LM5H.

## Frozen schedules

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

384 points per condition cell.

## Frozen comparator coordinate

Exactly the failed LM5G primary:
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

## Primary coordinate

Comparator plus:
- first_deliverable_round

Definition:
- the first global round at which the actual frozen scheduler delivers at least one action under the resource path;
- -1 only if no action is delivered.

This is a derived deployability state, not the raw configured action-start parameter.

## Measurements

Per residual condition cell × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread

48 summaries.

## Interpretation

Zero primary spread would show that deployment onset is the missing resource-path state. Residual spread would reject onset timing alone and require a deeper diagnostic.

## Bounds

Diagnostic only. No raw action-start field in the key, no raw budget field, no cut-round labels, no adaptive coordinate search, no resource tuning, no topology/profile expansion, no capacity change, no extra training, no live activation, or production authority.
