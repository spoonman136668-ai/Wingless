# Wingless UP-27: training-selected discovered commutant

Status: research branch only; not activated, promoted, or connected to ckb-plane.

Parent research result: UP-26 qualified strong partial positive sealed at `e9b94494eac451e94c97a5b17a26cd53ba0ced64`.

## Question

UP-26 showed that widening a generic dynamics-discovered observable bank from 32 to 64 raises unseen-depth accuracy from 92.48% to 99.88%, while the mutable commit gate remains narrowly negative at 98.4375%.

UP-27 does **not** increase runtime breadth.

It asks:

**Can training-only task relevance choose a better 64-observable subset from a larger dynamics-discovered candidate bank and close the remaining recurrent robustness gap?**

## Candidate discovery

UP-27 discovers 128 candidate observables using the exact UP-25/UP-26 dynamics-only construction:

- deterministic dense 96-dimensional coordinate-only seeds;
- 20 binary transport-conjugation averaging rounds;
- effective orbit size 1,048,576;
- no hidden six-channel multiplicity;
- no hidden full-state mixer;
- no semantic roles;
- no task labels during observable discovery.

The transport adjoint remains an offline discovery operation only.

## Selection data

Selection uses **only the same depth-0 training tables** used by the task learner:

- all 128 balanced training tables;
- two deterministic noisy trials per table;
- memory noise 0.05;
- memory-only global-phase nuisance.

No held-out table and no held-out depth participates in selection.

## Relevance score

For each candidate observable, UP-27 records its real and imaginary expectation values over the frozen training rows.

The task target is the eight real components formed by the UP-13 learned two-dimensional phase code for each of the four entities.

Each candidate receives a deterministic normalized linear relevance score:

- squared normalized covariance of its real component with all eight target components;
- plus squared normalized covariance of its imaginary component with all eight target components.

Candidates are sorted by descending score, with candidate index as the deterministic tie-breaker.

The top 64 are frozen as the runtime observable bank.

This is supervised feature selection, not held-out optimization.

## Prefix-64 control

The first 64 candidates from the 128-candidate discovery form a separate baseline arm.

Because discovery is seed-local and deterministic, this prefix must reproduce UP-26's 64-observable held-out accuracy:

`0.998779296875`

within absolute tolerance `1e-12`.

This verifies that any difference in the selected arm comes from selection rather than a changed discovery procedure.

## Learning

Both the prefix and selected arms use:

- exactly 64 observables;
- 128 raw real features;
- complete degree-2 lift: 8,384 features;
- ridge lambda 1e-6;
- UP-13 learned two-dimensional phase-code supervision;
- fresh two-dimensional softmax calibration heads;
- depth-0 training only.

## Held-out evaluation

Evaluation uses:

- all 128 disjoint held-out tables;
- depths 32, 128, 512, 1024;
- two deterministic noisy trials per table/depth.

The matched non-unitary control uses the selected 64-observable model unchanged.

## Mutable integration

The selected model runs:

- 48 scenarios;
- 16 writes each;
- 768 commit opportunities;
- held-out depths only;
- explicit irreversible overwrite/re-encode boundary;
- learned relational-query head.

## Scientific gates

**Prefix reproduction**

- absolute held-out difference from UP-26 primary64 <= 1e-12.

**Selection integrity**

- selection uses training labels: true;
- selection uses held-out data: false;
- selection uses explicit depth: false;
- candidate count: 128;
- runtime count: 64.

**Discovery commutator**

- maximum selected-observable commutator entry error <= 1e-5.

**Feature invariance**

- maximum selected degree-2 feature drift at depth 1024 <= 5e-3.

**Phase-code learning**

- selected depth-0 training accuracy >= 0.99.

**Unseen depth**

- selected held-out accuracy >= 0.99;
- maximum unitary norm drift <= 1e-10.

**Mutable integration**

- commit decode >= 0.99;
- exact final table >= 0.95;
- relational query >= 0.95;
- maximum unitary norm drift <= 1e-10.

## Interpretation boundary

A positive result means the dynamics-discovered candidate bank already contains enough stable task-relevant structure and that training-only selection can recover a robust 64-observable runtime bank without widening runtime state observation.

A negative result with a reproduced prefix control means raw single-observable relevance is insufficient. The next experiment should test task-learned **linear combinations or joint interaction-aware selection** of discovered commuting directions, not increase candidate count blindly.

No post-result changes to selector, lambda, data, rounds, runtime width, or thresholds are authorized.

No result establishes phase-specific or unitary superiority. A real orthogonal comparator remains required before stronger causal claims.

## Plain speak

UP-26 almost worked with 64 stable measuring tools, but 12 of 768 write checkpoints were still decoded incorrectly.

UP-27 does not give Wingless more tools at runtime.

Instead, it discovers a larger toolbox offline and uses only the training task to choose the 64 tools that are most related to the learned memory code.

If that closes the final gap, the problem was not capacity alone. It was choosing the right stable measurements.

## Authority boundary

UP-27 is mathematical research only. It does not invoke a language model, activate Wingless in production, alter ckb-plane/KTRADE authority, access brokers or credentials, modify accepted refs, or perform promotion/deployment.
