# Wingless UP-20: closed-form anonymous quadratic readout

Status: research branch only; not activated, promoted, or connected to ckb-plane.

Parent research result: UP-19 qualified scientific negative sealed at `c511ff937a02692c5d5bec0d1d397d1a3f3df22b`.

## Question

UP-19 established two facts:

1. the anonymous Gram representation remains numerically stable under deep unitary transport;
2. gradient-trained linear and complete quadratic softmax heads remain near chance even on the training set.

UP-20 isolates the remaining ambiguity:

**Is the UP-19 failure caused by the anonymous quadratic representation, or by the gradient learner?**

The state representation is unchanged.

## Exact quadratic witness

The six anonymous channels are related to the historical semantic channels by a fixed real orthogonal mixer.

For a semantic Gram entry `H[a,b]`, the inverse congruence

`H = M^T G M`

makes each semantic correlation a linear function of the 36 real anonymous Gram coordinates.

The established frame numerator for entity `e` is

`<pilot_e,memory> * conj(<anchor,memory>)`.

That quantity is therefore exactly quadratic in the anonymous Gram coordinates.

UP-20 constructs the corresponding degree-2 coefficient vector analytically from the hidden mixer **only as a diagnostic witness**. It verifies that the 702-dimensional feature map reproduces the historical numerator to numerical precision.

The witness is never supplied to the learned ridge path.

## Oracle-witness classification control

The exact witness projects the anonymous quadratic features to two real coordinates: real and imaginary parts of the historical phase numerator.

A fresh two-feature softmax head is trained on those diagnostic coordinates and evaluated on unseen depths.

This control answers whether the quadratic feature basis actually contains a task-sufficient direction.

## Closed-form learner

The learned path receives only:

- the same 702 anonymous quadratic features from UP-19;
- task class labels.

It receives no:

- semantic channel labels;
- hidden mixer;
- oracle witness weights;
- mixer inverse;
- depth;
- prototype lookup.

Instead of iterative softmax gradient descent, UP-20 uses deterministic linear ridge regression in dual form.

Training rows:

- all 128 balanced training tables;
- memory noise 0.05;
- memory-only global-phase nuisance;
- four deterministic noisy trials per table;
- 512 total training rows.

Ridge regularization is frozen at:

`lambda = 1e-6`.

The dual solve is converted back into ordinary primal 4 × 702 weights plus biases for each entity. Training examples are **not** retained at runtime.

## Evaluation

The same learned ridge heads are evaluated on:

- all 128 disjoint held-out tables;
- depths 32, 128, 512, 1024;
- two noisy trials per table/depth.

The matched non-unitary path receives the exact same learned heads.

## Mutable integration

The learned ridge heads then run:

- 48 scenarios;
- 16 writes each;
- 768 commit opportunities;
- held-out depths only;
- explicit irreversible overwrite/re-encode boundary;
- learned relational-query head.

## Scientific gates

**Oracle witness**

- exact quadratic numerator error <= 1e-12;
- oracle-witness held-out accuracy >= 0.99.

**Closed-form learning**

- ridge canonical training accuracy >= 0.99.

**Unseen depth**

- ridge unitary held-out accuracy >= 0.99;
- maximum unitary norm drift <= 1e-12.

**Mutable integration**

- commit decode >= 0.99;
- exact final table >= 0.95;
- relational query >= 0.95;
- maximum unitary norm drift <= 1e-12.

## Interpretation boundary

If the oracle witness passes and ridge succeeds, UP-19 was primarily an optimizer failure: semantic demixing is unnecessary and a direct anonymous quadratic readout is learnable.

If the oracle witness passes but ridge fails, the feature span is sufficient but generic one-hot ridge learning still does not identify the useful low-rank phase direction. The next experiment should learn that low-rank relational factorization directly rather than increasing polynomial degree.

If the oracle witness fails, the algebraic premise or implementation is defective and must be repaired before interpreting the experiment.

No positive result here establishes phase-specific or unitary superiority. A real orthogonal transport comparator remains mandatory before stronger causal claims.

## Plain speak

UP-19 gave the learner all the right puzzle pieces but its normal training method could not assemble them.

UP-20 first proves, mathematically and numerically, that the answer really is present in those pieces.

Then it replaces gradual trial-and-error learning with a direct solve.

If that works, the problem was not the anonymous state. It was the way we were teaching the reader.

## Authority boundary

UP-20 is mathematical research only. It does not invoke a language model, activate Wingless in production, alter ckb-plane/KTRADE authority, access brokers or credentials, modify accepted refs, or perform promotion/deployment.
