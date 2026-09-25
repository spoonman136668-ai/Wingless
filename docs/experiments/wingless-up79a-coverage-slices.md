# Wingless UP-79A — 36-state structural-slice replication

Status: preregistered scientific replication experiment.

Scientific parent: sealed UP-78A Windows evidence `2a35f521daa7898d714959a602994f5b04b15580`.

## Question

UP-78A localized the one-step coverage transition between 27 and 36 states. Its 36-state pass added the fixed secondary/tertiary slice `s=1,t=0` to the 27-state `s=0` base. UP-79A tests whether the 36-state result reflects generic structural coverage or depends specifically on that chosen nine-state tertiary slice.

## Frozen design

Unchanged from UP-78A:

- the same 243 five-role ternary states;
- the same primary residue split;
- the same 162 held-out states with primary residue != 0;
- raw factorized one-hot representation, 15 features;
- role-factorized simplex representation, 10 features;
- five independent 3-class linear-softmax heads;
- zero initialization;
- learning rate 1.0;
- exactly one optimization step;
- per-role held-out gate 0.98.

The frozen training conditions are exactly:

- base `s=0` plus `s=1,t=0`;
- base `s=0` plus `s=1,t=1`;
- base `s=0` plus `s=1,t=2`.

Each condition contains exactly 36 training states: the same 27-state base plus one disjoint nine-state tertiary slice.

No adaptive stopping, extra step, learning-rate change, threshold change, seed search, post-result slice selection, representation change, promotion, or activation is permitted.

## Interpretation

Every slice/representation cell is sealed exactly as observed. If all three slices pass, that supports a generic 36-state structural-coverage threshold under the frozen task. If only some pass, the UP-78A transition is slice-dependent and the specific added structure matters. Harness acceptance is independent of the scientific result.
