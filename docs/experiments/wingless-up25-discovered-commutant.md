# Wingless UP-25: dynamics-discovered commuting observables

Status: research branch only; not activated, promoted, or connected to ckb-plane.

Parent research result: UP-24 qualified scientific pass sealed at `08acd3e00736069e3682813c4a905691697ed505`.

## Question

UP-24 restored perfect full-latent decoding with a transport-commuting Weyl observable algebra, but that algebra was hand-constructed from the known six-fold multiplicity structure.

UP-25 asks:

**Can useful commuting observables be discovered from the full 96-dimensional transport itself without using the hidden multiplicity factorization or the hidden full-state mixer?**

## Discovery seeds

UP-25 creates 32 deterministic dense rank-one complex seed observables directly in the 96-dimensional latent coordinate system.

The seed vectors depend only on:

- the seed index;
- the 96 coordinate indices.

They do not depend on:

- semantic channel identities;
- the historical six-channel factorization;
- the full-latent mixer;
- task labels;
- the UP-24 Weyl basis.

## Binary conjugation averaging

Let `V` be the full-latent one-step unitary transport and `A` a seed observable.

The commuting projection is approximated by repeatedly applying:

`A <- (A + V^(2^r) A V^(-2^r)) / 2`

for 20 rounds.

This is algebraically equivalent to averaging the seed over an orbit of:

`2^20 = 1,048,576`

transport conjugations, but requires only 20 doubling rounds.

The transport adjoint is used **only during offline observable discovery**.

The resulting matrices are normalized and frozen.

Runtime feature extraction does not apply the transport adjoint or an inverse.

## Representation

For each state, the observer measures the complex expectation of each of the 32 discovered observables:

- 64 real raw features.

A complete degree-2 lift gives:

- 2,144 total features.

The task learner then performs:

- deterministic ridge regression to the already learned two-dimensional UP-13 phase code;
- fresh two-dimensional softmax calibration heads.

## Training

Training uses only depth 0:

- all 128 balanced training tables;
- two deterministic noisy trials per table;
- memory noise 0.05;
- memory-only global-phase nuisance;
- ridge lambda 1e-6.

No held-out depth is represented during training.

## Held-out evaluation

Evaluation uses:

- all 128 disjoint held-out tables;
- depths 32, 128, 512, 1024;
- two noisy trials per table/depth.

The matched non-unitary path receives the exact same learned regressors, classifiers, and discovered unitary observables.

## Mutable integration

The learned unitary model runs:

- 48 scenarios;
- 16 writes each;
- 768 commit opportunities;
- held-out depths only;
- explicit irreversible overwrite/re-encode boundary;
- learned relational-query head.

Every recurrent segment remains one fully mixed 96-dimensional state.

## Scientific gates

**Discovery commutator**

- maximum entrywise commutator error <= 1e-5.

**Feature invariance**

- maximum degree-2 feature drift at depth 1024 <= 5e-3.

**Phase-code learning**

- depth-0 training accuracy >= 0.99.

**Unseen depth**

- unitary held-out accuracy >= 0.99;
- maximum unitary norm drift <= 1e-10.

**Mutable integration**

- commit decode >= 0.99;
- exact final table >= 0.95;
- relational query >= 0.95;
- maximum unitary norm drift <= 1e-10.

## Interpretation boundary

A positive result means Wingless can derive a useful stable measurement structure from the transport dynamics themselves without being handed the six-channel multiplicity algebra.

A negative result must be split:

- poor commutator / poor invariance means the finite discovery projection is inadequate;
- good commutator / good invariance but poor task accuracy means generic discovered commutant directions do not capture enough task-relevant orientation.

No negative result authorizes increasing rounds, seed count, or changing thresholds after inspection. A successor experiment must isolate the identified failure mode.

No result here establishes phase-specific or unitary superiority. A real orthogonal comparator remains required before stronger causal claims.

## Plain speak

UP-24 gave Wingless a perfect set of stable internal measuring tools, but we built those tools for it.

UP-25 starts with arbitrary dense measuring tools and lets the dynamics repeatedly wash away every part that does not stay aligned with the system’s own motion.

What remains should be a set of naturally stable measurements discovered from the dynamics themselves.

If those measurements still recover memory perfectly, the system no longer needs us to tell it where its internal “compass” lives.

## Authority boundary

UP-25 is mathematical research only. It does not invoke a language model, activate Wingless in production, alter ckb-plane/KTRADE authority, access brokers or credentials, modify accepted refs, or perform promotion/deployment.
