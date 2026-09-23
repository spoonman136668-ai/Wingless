# Wingless UP-18: orthogonal anonymous-channel demixer

Status: research branch only; not activated, promoted, or connected to ckb-plane.

Parent research result: UP-17 qualified scientific negative sealed at `652cefb5ce4b629a048f10c9ed82c622695c35b1`.

## Question

UP-17 showed that the row-normalized 6 × 6 task-learned readout projection did not recover hidden semantic channel roles.

UP-18 separates two explanations:

1. did the dense anonymous channel mixer destroy task-relevant information?
2. or did information survive while the UP-17 learner used a poor optimization geometry?

## Frozen substrate

UP-18 freezes:

- the UP-13 learned value-phase alphabet;
- the UP-14 learned pilot support;
- the UP-15 learned anchor;
- the UP-16 single 96-dimensional recurrent state;
- the exact UP-17 dense orthogonal anonymous-channel mixer;
- the same table pools, transport depths, memory noise, and global-phase nuisance.

No harder scrambling is added.

## Oracle recoverability control

Because the fixed anonymous-channel mixer is orthogonal, its transpose is its exact inverse.

UP-18 uses that transpose only as a **diagnostic control**.

The oracle:

- is not supplied to the task learner;
- does not initialize learned angles;
- is not used in the learned runtime path;
- is not an acceptance shortcut.

It answers only whether the information needed by the established observer is still present after scrambling.

Gates:

- maximum recovered feature error <= 1e-12;
- oracle unitary held-out accuracy >= 0.99;
- oracle unitary norm drift <= 1e-12.

If this fails, UP-17 scrambling itself damaged the usable representation.

If it passes, UP-17 was a learning problem rather than an information-preservation problem.

## Structured learned demixer

UP-17 learned 36 independent real matrix entries and normalized rows. That parameterization did not preserve orthogonality.

UP-18 instead uses the exact dimension of O(6):

**15 trainable Givens angles.**

The learned readout projection remains orthogonal for every parameter value.

The parameterization begins at identity. The true inverse is representable in the family, but its target angles are never exposed to learning.

Learning uses only classification loss.

Parameters:

- 64 balanced learning tables;
- memory noise 0.05;
- memory-only global-phase nuisance;
- two deterministic noisy trials per table;
- 100 outer steps;
- 12 head updates per outer step;
- angle learning rate 0.18;
- central-difference epsilon 1e-4.

Post-hoc matrix distance and role alignment are diagnostics only.

## Fresh evaluation

After angle learning, alternating heads are discarded.

Fresh entity decoders train on:

- all 128 balanced training tables;
- depths 8, 24, 72, 216, 432, 648;
- four noisy trials per table/depth;
- 1,200 softmax steps.

Held-out evaluation uses:

- all 128 disjoint held-out tables;
- unseen depths 32, 128, 512, 1024;
- two noisy trials per table/depth.

The matched non-unitary path receives the same learned orthogonal projection and decoder budget.

## Mutable integration

The learned unitary path then runs:

- 48 scenarios;
- 16 writes each;
- 768 commit opportunities;
- held-out depths only;
- memory noise 0.05;
- memory-only global phase nuisance;
- explicit irreversible overwrite/re-encode boundary.

## Scientific gates

**Oracle recoverability**

- feature error <= 1e-12;
- held-out accuracy >= 0.99;
- norm drift <= 1e-12.

**Structured learning**

- identity initial capacity < 0.70;
- final task loss < initial task loss;
- angle vector moves >= 0.50 L2 from identity;
- final unitary train accuracy >= 0.99.

**Unseen depth**

- learned unitary held-out accuracy >= 0.99;
- learned unitary norm drift <= 1e-12.

**Mutable integration**

- commit decode >= 0.99;
- exact final table >= 0.95;
- relational query >= 0.95;
- norm drift <= 1e-12.

## Interpretation boundary

A positive oracle control with a negative structured learner means the representation remains recoverable but task loss still does not discover the hidden roles under this learner.

A positive structured learner means UP-17 failed primarily because its unconstrained demixer geometry was unsuitable.

Neither outcome justifies phase-specific superiority claims. The real orthogonal transport comparator remains mandatory before stronger causal claims.

## Plain speak

UP-17 mixed all the labeled drawers together and our first learner could not sort them back out.

UP-18 first checks whether the contents are actually still there. Then, instead of giving the learner a floppy 36-number matrix, we give it a rigid 15-joint rotation mechanism that can only rotate the internal space without stretching or collapsing it.

If that works, the problem was the learner's steering mechanism—not missing information.

## Authority boundary

UP-18 is mathematical research only. It does not invoke a language model, activate Wingless in production, alter ckb-plane/KTRADE queue or acceptance authority, access brokers or credentials, or promote/deploy code.
