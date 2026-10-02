# Wingless UP-24: commuting Weyl observable algebra

Status: Windows-qualified scientific pass; not activated, promoted, or connected to ckb-plane.

Parent research result: UP-23 qualified scientific negative sealed at `8e54fba4e58e632c37d580e2d9cf4b58e92ba4f7`.

## Question

UP-23 proved that scalar transport moments are depth-invariant, but they discard orientation inside the six-fold repeated eigenspaces inherited from the historical composite transport.

UP-24 tests the causal diagnosis:

**Can a generic algebra of observables that commutes with the full 96-dimensional transport recover the missing multiplicity-space information while keeping the runtime state fully mixed and block-free?**

## Commuting observable algebra

The historical composite transport has the form:

`T = I_6 tensor U`.

Any operator of the form:

`B tensor I_16`

commutes with `T`.

UP-24 uses the complete six-dimensional discrete Weyl operator basis:

`X^a Z^b`, for `a,b in {0,...,5}`.

There are 36 such operators.

Each is lifted to the 96-dimensional historical composite space as:

`(X^a Z^b) tensor I_16`

and then conjugated into the fully mixed latent basis by the same fixed full-coordinate mixer used by UP-22:

`A_ab = Q [(X^a Z^b) tensor I_16] Q^H`.

The recurrent state stays in the fully mixed basis at all times.

The observer never applies `Q^H` to the state and never receives semantic channel labels.

## Feature map

For each full-latent state `psi`, UP-24 measures:

`<psi, A_ab psi>`

for all 36 commuting Weyl observables.

Each complex expectation contributes real and imaginary parts:

- 72 raw real features.

A complete degree-2 lift then appends all unique pairwise products:

- 2,700 total features.

The phase-code learner is the same conceptual target that closed UP-21:

- deterministic ridge regression to the already task-learned UP-13 2-D phase alphabet;
- fresh 2-D softmax calibration heads.

## Why this is only a diagnostic bridge

The observable algebra is **hand-constructed** using the known transport multiplicity dimension 6.

That assumption is explicit in the schema and is not being hidden.

A positive result would show that the UP-23 failure really was caused by collapsing multiplicity-space orientation, and that a transport-commuting observable algebra can recover it while the recurrent state remains fully mixed.

It would not yet prove that Wingless can discover that algebra autonomously.

If positive, the next experiment must attempt to learn/discover the commuting algebra from transport/task evidence rather than being handed the Weyl construction.

## Training

Training uses only depth 0:

- all 128 balanced training tables;
- two deterministic noisy trials per table;
- memory noise amplitude 0.05;
- memory-only global-phase nuisance;
- ridge lambda 1e-6.

No held-out depth is represented during training.

## Held-out evaluation

Evaluation uses:

- all 128 disjoint held-out tables;
- depths 32, 128, 512, 1024;
- two noisy trials per table/depth.

The matched non-unitary arm uses the exact same learned unitary regressor and calibration heads.

## Mutable integration

The learned unitary model then runs:

- 48 scenarios;
- 16 writes each;
- 768 commit opportunities;
- held-out depths only;
- explicit irreversible overwrite/re-encode boundary;
- learned relational-query head.

Every recurrent segment remains a single fully mixed 96-dimensional state.

## Scientific gates

**Commutator gate**

- maximum entrywise commutator error between each observable and the full-latent one-step transport <= 1e-10.

**Feature invariance**

- maximum degree-2 feature drift at depth 1024 <= 1e-9.

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

A positive result would support the degeneracy diagnosis from UP-23 and show that the fully mixed state is usable without runtime unmixing or visible block boundaries when the observer has access to a sufficiently rich commuting observable algebra.

A negative result would mean that the hand-constructed commutant basis plus degree-2 phase-code learner is still insufficient, and the next step must reassess either the observable algebra or the encoding rather than simply adding training depths.

No result here establishes phase-specific or unitary superiority. A real orthogonal comparator remains required before stronger causal claims.

## Plain speak

UP-23 listened only to the state’s overall resonance and lost track of how information was arranged inside repeated resonance bands.

UP-24 gives the reader a full set of “polarizers” that all stay aligned with the dynamics.

They do not unscramble the 96-dimensional state. They only measure it in ways that remain stable while it evolves.

If that restores perfect memory, we have located the missing ingredient: not labeled drawers, but a stable internal measurement algebra.

## Authority boundary

UP-24 is mathematical research only. It does not invoke a language model, activate Wingless in production, alter ckb-plane/KTRADE authority, access brokers or credentials, modify accepted refs, or perform promotion/deployment.


## Windows qualification result

Autonomous Windows qualification on workflow run `35868771822` completed successfully at source head `98e173b3c785526f40f7d5a94f1c92ecdc865719`.

All preregistered scientific gates passed. Maximum observable commutator error was 1.622242763424964e-14 and maximum feature drift was 1.7917889394425401e-12. The unitary path achieved 1.0 training accuracy, 1.0 held-out accuracy, and 1.0 / 1.0 / 1.0 mutable commit, exact-final, and relation accuracy. Minimum value margin was 0.6370882275394993.

Interpretation: UP-23 failed because scalar spectral moments collapsed the six-fold multiplicity-space orientation. A sufficiently rich commuting observable algebra preserves that orientation without runtime unmixing or visible channel blocks. The remaining cheat is explicit: UP-24 was handed the Weyl algebra using the known multiplicity structure. UP-25 must discover commuting observables from the transport itself.
