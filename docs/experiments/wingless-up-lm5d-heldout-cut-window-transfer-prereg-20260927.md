# Wingless UP-LM5D — heldout asynchronous cut-window transfer

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM5C 0e2c683764c234b28289be43dc0a859fc00e14fb.

## Question

Does the frozen timing-aware ordered-path coordinate transfer to asynchronous resource-cut schedules not used in LM5A-LM5C?

## Heldout asynchronous windows

None of these six cut/restore combinations appeared in LM5A-LM5C:

- budget 2/3 ; throughput 1/4
- budget 1/4 ; throughput 2/3
- budget 2/4 ; throughput 3/5
- budget 3/5 ; throughput 2/4
- budget 2/3 ; throughput 1/5
- budget 1/5 ; throughput 2/3

These preserve the previously tested ordered restore-time pairs:
- early/late 3/4
- early/late 4/5
- early/late 3/5

but change when the temporary reductions begin.

## Frozen environment

Deadline profiles:
- deferred_only
- layout_only
- hybrid_min

Reductions:
- budget 1,2
- throughput 1,2

Topology:
- rotations 5,13
- permutations identity, reverse, rotate2

Resource grid:
- budgets 4,5,6,7
- starts 2,3,4,5
- throughput 2,3,4,5
- six global rounds
- earliest_deadline scheduling

Pooling:
- all six heldout schedules pooled per profile × budget reduction × throughput reduction × rotation × permutation
- 384 points per pooled cell
- 72 pooled cells

## Frozen coordinates

Comparator:
- ordered path without restore rounds

Primary:
- the accepted LM5C timing-aware ordered path:
  - final reachable actions
  - nominal throughput
  - final coverage end round
  - actions before earlier restoration
  - cumulative actions immediately after earlier restoration
  - actions before later restoration
  - earlier restoration round
  - later restoration round

Neither coordinate includes cut-round labels, resource-restoration identity, or full schedule identity.

## Measurements

Per pooled cell × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread

144 summaries.

## Interpretation

Zero primary spread would show that restore timing plus ordered path state generalizes across unseen cut timing. Residual spread would identify a genuine dependence on how the constraint begins rather than how much work is delivered through the path.

## Bounds

Diagnostic only. No adaptive schedule selection, cut-round labels in keys, resource identity in keys, coordinate search, resource tuning, topology/profile changes, capacity change after results, extra training, live activation, or production authority.
