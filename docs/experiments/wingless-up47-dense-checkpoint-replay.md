# Wingless UP-47 — frozen dense checkpoint replay

Status: Windows-qualified scientific result; dense replay selected step 42 and hard capability gates remained false.

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


## Authoritative Windows result

Workflow run: `35980735036`

Source head: `474ebb42082d1a9dd022517edb3c48c4f16dbd1d`

Runner: `WINGLESS-LINKDEADKB`

Artifact: `10800157359`

Artifact digest: `sha256:b43d7ab58f366b2d0b823c56967594d323daca824b02aae0edd78d14f59ed52a`

The production-priority guard, focused tests, deterministic double probe, and full repository regression passed.

Dense replay selected step `42`, not sparse step `24` and not exact-fused step `29`.

- dense-selected objective: `0.591814132893378`
- sparse step-24 objective: `0.5581278352889602`
- exact-fused step-29 objective: `0.5591129650448196`
- dense-minus-sparse objective: `+0.03368629760441788`
- selected exact capacity: `6`
- selected exact fusion: `false`
- held-out accuracy: `0.906982421875`
- mutable commit accuracy: `0.40625`
- exact final-table accuracy: `0.4791666666666667`
- relation accuracy: `0.5833333333333334`
- capability gates: `false`

## Scientific classification

The sparse UP-45 checkpoint schedule did hide a substantially better objective point, but the best recorded state was a non-fused step 42 and still failed the hard recurrent capability gates.

This separates two questions cleanly:

1. why step 42 maximizes the smooth objective without recovering hard capability; and
2. why exact-fused step 29 has markedly better hard recurrent capability than both step 24 and step 42 despite a lower smooth objective than step 42.

The post-UP47 program may therefore branch without changing the frozen result: mechanism work can compare the three frozen geometries, cognition work can test broader compositional capability independently, and stress work can probe scaling/failure boundaries independently.

## Plain speak

Looking at every recorded step fixed one mistake: step 24 was not actually the best point on the path.

But the new winner, step 42, still did not solve the real memory/relation task. That means the smooth score and the hard capability are still not aligned enough. We now have three useful frozen states—24, 29, and 42—that let us study that mismatch while other lanes push cognition and scale in parallel.
