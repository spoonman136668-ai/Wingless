# Wingless UP-28: interaction-aware discovered commutant selection

Status: research branch only; not activated, promoted, or connected to ckb-plane.

Parent research result: UP-27 qualified selector-negative result sealed at `04074246d0ca13ae7714bbf26e40eb86ed14d7cd`.

## Question

UP-27 showed that ranking discovered observables by their individual linear association with the phase target is insufficient.

The runtime decoder is not linear in raw observable expectations: it uses a complete degree-2 lift.

UP-28 asks:

**Can training-only selection based on the full quadratic task model identify a more robust 64-observable runtime bank?**

## Candidate discovery

UP-28 reuses the frozen dynamics-only candidate construction:

- 128 deterministic dense discovered observables;
- 20 binary transport-conjugation averaging rounds;
- effective orbit size 1,048,576;
- no hidden multiplicity;
- no hidden full-state mixer;
- no semantic labels during discovery.

## Training-only selector

Selection uses only:

- all 128 balanced training tables;
- depth 0 only;
- two deterministic noisy trials per table;
- memory noise 0.05;
- memory-only global-phase nuisance;
- the already learned UP-13 2-D phase target.

No held-out table or held-out depth participates.

## Interaction-aware importance

For all 128 candidates, UP-28 forms 256 raw real expectation features.

The selector fits a deterministic ridge model from the **complete degree-2 feature map** to the eight phase-target components.

The dual Gram matrix is evaluated with the exact degree-2 polynomial feature inner product, avoiding an implementation mismatch with the runtime decoder.

After solving the training-only quadratic model, UP-28 reconstructs the primal linear and pairwise coefficients.

Each observable receives importance from:

- squared linear coefficients involving its real/imaginary feature components;
- squared quadratic coefficients involving those components;
- pairwise interaction energy split equally across the two observables involved.

The 64 highest-scoring observables are frozen for runtime.

## Controls

**Prefix64**

The first 64 candidates must reproduce UP-26 held-out accuracy:

`0.998779296875`

within absolute tolerance `1e-12`.

**UP-27 overlap**

The selected bank reports its overlap with the UP-27 marginal-correlation bank as a diagnostic only.

## Runtime model

The selected runtime bank remains exactly 64 observables:

- 128 raw real features;
- 8,384 complete degree-2 features;
- ridge lambda 1e-6;
- UP-13 learned 2-D phase-code supervision;
- fresh 2-D softmax calibration heads.

Runtime width is unchanged from UP-26 and UP-27.

## Held-out evaluation

Evaluation uses:

- all 128 disjoint held-out tables;
- depths 32, 128, 512, 1024;
- two deterministic noisy trials per table/depth.

The matched non-unitary control uses the same selected bank and learned model.

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

**Selector integrity**

- training labels used: true;
- held-out data used: false;
- explicit depth used: false;
- quadratic interactions used: true;
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

A positive result means task-relevant quadratic interactions provide enough information to choose a robust 64-observable runtime bank from the dynamics-discovered commutant.

A negative result with reproduced controls means hard subset selection itself is likely the limiting abstraction. The next experiment should learn task-supervised **linear combinations** of commuting candidates so multiple weak but complementary directions can contribute without widening runtime width.

No post-result changes to candidate count, runtime width, lambda, training data, rounds, or thresholds are authorized.

No result establishes phase-specific or unitary superiority. A real orthogonal comparator remains required before stronger causal claims.

## Plain speak

UP-27 judged each measuring tool mostly on its own.

But the successful decoder works from relationships between measurements.

UP-28 judges each tool by how much it contributes to the full web of pairwise relationships used by the task.

If this closes the remaining write errors, we have a principled way to pick the right 64 stable measurements instead of simply keeping the first 64 discovered.

## Authority boundary

UP-28 is mathematical research only. It does not invoke a language model, activate Wingless in production, alter ckb-plane/KTRADE authority, access brokers or credentials, modify accepted refs, or perform promotion/deployment.
