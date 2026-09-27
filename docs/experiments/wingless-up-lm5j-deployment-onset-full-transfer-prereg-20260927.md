# Wingless UP-LM5J — deployment-onset full transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM5I 40ebd27d1d2bb2ffd5877ffc1e23c5ec2ec8b012.

## Question

Does the deployment-onset state variable that resolved the exact LM5I residual boundary remain sufficient across the full heldout round-2 restore-time set, including cells that were already unambiguous?

## Frozen heldout schedules

Exactly the six UP-LM5E schedules:
- budget 0/2 ; throughput 1/3
- budget 1/3 ; throughput 0/2
- budget 0/2 ; throughput 2/4
- budget 2/4 ; throughput 0/2
- budget 0/2 ; throughput 3/5
- budget 3/5 ; throughput 0/2

## Frozen deadline profiles
- deferred_only
- layout_only
- hybrid_min

## Frozen reductions
- budget reduction 1,2
- throughput reduction 1,2

## Frozen topology
- rotations 5,13
- permutations identity, reverse, rotate2

## Frozen resource grid
- budgets 4,5,6,7
- action starts 2,3,4,5
- throughput 2,3,4,5
- six global rounds
- earliest_deadline scheduling

384 points per pooled cell; 72 pooled cells total.

## Frozen comparator coordinate

Exactly the LM5I frozen coordinate:
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

## Frozen primary coordinate

Comparator plus:
- first_deliverable_round

Definition is unchanged from LM5I:
- first global round at which the frozen scheduler actually delivers at least one action
- -1 only if no action is delivered

No raw configured action-start value is included.

## Measurements

Per pooled cell × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread

144 summaries.

## Interpretation

Zero primary spread across all 72 cells would establish deployment onset as a transferable completion of the current path-state representation for the heldout round-2 restore regime. New residual spread would bound that conclusion to the LM5I residual subset.

## Bounds

Diagnostic only. No raw action-start field, no raw budget field, no cut-round labels, no adaptive coordinate search, no resource tuning, no topology/profile changes, no capacity change, no extra training, no live activation, or production authority.
