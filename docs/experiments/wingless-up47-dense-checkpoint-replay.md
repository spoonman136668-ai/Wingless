# Wingless UP-47 — frozen dense checkpoint replay

Status: preregistered scientific diagnostic.

Scientific parent: UP-46 seal `3a4a5c209818a803b5451d2ca4adf74c62d16680`.

## Question

UP-46 established that UP-45's exact-fused step 29 scores better than sparse-selected step 24 on the same smooth task objective and also performs substantially better on hard recurrent capability.

UP-45 did not select step 29 because selection was only evaluated at steps 0, 12, 24, and 48.

UP-47 asks what the unchanged objective would have selected if **every already-recorded UP-45 trajectory state** had been eligible.

## Frozen data

The complete 49-state sequence consists of UP-45 initial state 0 plus the authoritative output offsets from steps 1 through 48 from Windows run `35968534136`.

The states are embedded verbatim in source and protected by tests at steps 0, 24, 29, and 48.

No optimizer is rerun.

## Evaluation

For every frozen state 0..48:

- evaluate the unchanged UP-45 smooth task objective with the untouched reference evaluator;
- report phase, value probability, relation probability, harmonic task score, soft capacity, resource penalty, objective, exact capacity, and exact fusion.

Select the maximum training-side objective across all 49 points, using the existing lower-soft-capacity tie break only for exact objective ties.

Only after that selection is frozen may the true 128-table held-out pool be used to evaluate the selected state.

## Frozen interpretation

- If dense selection chooses an exact-fused state with objective above sparse step 24, sparse checkpoint scheduling materially hid task-supported self-reorganization. If hard capability also improves, next isolate exact equality causally against a matched near-symmetry control at the dense-selected geometry.
- If dense selection chooses an exact-fused state and the full capability gates pass, proceed instead to restart persistence/generalization before any further mechanism change.
- If dense selection chooses a non-fused state with objective above step 24, sparse checkpointing still lost objective quality, but exact fusion is not necessary for the best point on this trajectory. Next compare the dense-selected geometry with step 29 to identify the structural feature responsible.
- If dense selection still chooses step 24 despite the independently reproduced step-29 objective being higher, treat that as an implementation/evaluation inconsistency and repair only the defect.
- If another state dominates both step 24 and step 29, accept that state as the dense replay result; do not privilege exact fusion merely because step 29 was previously interesting.

No optimizer, lambda, objective component, state, seed, gate, or held-out rule may be changed after seeing the result.

## Causal boundary

UP-44 FIXA at `013339e102080ce57e7600ebacb8ab086fd53467` showed that square-root probability encoding improves the smooth objective under the seed-matched control but did not replicate the stronger symmetry-selection claim.

Therefore UP-47 is an internal analysis of the observed UP-45 trajectory. It does not restore a representation-only causal claim.

## Plain speak

We already have the path the system took. We are not letting it train again.

We simply score every point on that recorded path instead of looking only every 12 steps.

That tells us whether the system actually found a better self-organization and our sparse inspection schedule failed to notice it.
