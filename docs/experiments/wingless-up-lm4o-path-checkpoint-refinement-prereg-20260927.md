# Wingless UP-LM4O — revoke/restore path-checkpoint refinement

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM4N e85339d458dd4ea7989d33175fdafa4f840a1e4a.

## Question

Which minimal temporal checkpoint resolves the path dependence exposed by temporary budget revocation and restoration?

## Frozen scenario

Exactly reuse UP-LM4N:
- deadline profiles: deferred_only, layout_only, hybrid_min
- initial budgets 4,5,6,7
- action starts 2,3,4,5
- throughput 1,2,3,4
- six rounds
- earliest_deadline scheduling
- rotations 5 and 13
- permutations identity, reverse, rotate2
- revoke windows (2→4), (2→5), (3→5)
- revoke amounts 1 and 2

64 resource configurations per topology-condition cell.
108 topology-condition cells.

## Frozen base quantities

Under the actual revoke/restore cap schedule:
- dynamic_reachable: total actions that can actually be delivered
- dynamic_coverage_end: final round by which that reachable amount is delivered
- nominal throughput
- pre_cut_delivery: actions deliverable through the round immediately before revocation begins
- pre_restore_delivery: actions deliverable through the final revoked round, immediately before restoration

## Frozen candidate coordinates

A. dynamic base:
- dynamic_reachable + throughput + dynamic_coverage_end

B. early checkpoint:
- dynamic base + pre_cut_delivery

C. revoked-window checkpoint:
- dynamic base + pre_restore_delivery

D. both checkpoints:
- dynamic base + pre_cut_delivery + pre_restore_delivery

No other coordinate will be considered in this experiment.

## Measurements

Per profile × revoke schedule × revoke amount × rotation × permutation × coordinate:
- groups
- multi-member groups
- nonzero-spread groups
- maximum failure spread
- mean failure spread

432 coordinate summaries total.

## Interpretation

If pre_restore alone restores exact collapse, the critical missing state is how much resource can be delivered before restoration. If both checkpoints are needed, the path requires at least two temporal landmarks. If all remain nonzero, this small checkpoint representation is falsified.

## Bounds

Diagnostic only. No adaptive checkpoint selection, coordinate search, resource tuning, schedule/profile/topology selection, extra training, live activation, or production authority.
