# Wingless UP-32: repeated-spectrum multiplicity dose response

Status: Windows-qualified ordered dose trend with preregistered material-endpoint gate negative; not activated, promoted, or connected to ckb-plane.

Parent research result: UP-31 qualified causal positive sealed at 601ac75b827d239b1dbdce1df16f6627877b31d4.

## Question

UP-30 showed a catastrophic loss when the six repeated transport copies were fully split.

UP-31 showed that exact coordinate alignment across those copies is not required: preserving the six-fold spectrum while hiding each copy in a different internal basis retained essentially perfect performance.

UP-32 asks whether capability degrades systematically as repeated-spectrum multiplicity is reduced.

## Anonymous split basis

The phase splits are not attached directly to the historical memory/anchor/pilot roles. Before offsets are applied, the six-copy axis is rotated with the existing dense deterministic 6 x 6 orthogonal channel mixer. The split eigendirections are therefore anonymous linear combinations of all six historical roles. The observer is not given this split basis.

## Dose arms

Five arms are tested: multiplicity 6 control; 4+2; 3+3; 2+2+2; and 1+1+1+1+1+1.

Every non-control arm has zero-mean offsets, fixed RMS offset magnitude 0.05 radians per step, the same 96-complex / 192-real scalar state budget, norm-preserving transport, and an exact real-orthogonal equivalent. Lower multiplicity is therefore not automatically paired with a larger average perturbation.

## Observer protocol

Every arm receives the same learning budget: discover 128 approximately commuting observables from that arm's transport; use depth-0 training states only; interaction-aware training-label-only selection; select 64 runtime observables; regress to the already learned UP-13 two-dimensional phase code; and use fresh 2-D calibration heads.

The observer receives no hidden factorization, split basis, held-out data, explicit depth, runtime unmix, inverse transport, or prototypes.

## Scientific gates

Structural validity requires, for every arm: real orthogonality error <= 1e-10; selected-bank commutator error <= 1e-5; feature drift <= 5e-3; and every non-control arm offset RMS equals 0.05 within 1e-12.

Base control requires held-out >= 0.99, mutable commit >= 0.99, exact final >= 0.95, and relational query >= 0.95.

Endpoint material effect requires multiplicity-6 minus multiplicity-1 held-out accuracy >= 0.25 and commit accuracy >= 0.25.

Ordered-dose criterion: at least 3 of 4 adjacent steps must be non-improving within a +0.03 tolerance for both held-out accuracy and mutable commit accuracy.

The raw curve is authoritative even if the preregistered dose-response gate is negative.

## Interpretation boundary

A positive ordered dose response would materially strengthen the hypothesis that repeated-spectrum multiplicity supplies useful transport-stable degrees of freedom for the observable algebra. A threshold-shaped response is also scientifically meaningful. A noisy or non-ordered response would indicate that multiplicity alone is not enough and that invariant-subspace arrangement or finite observer discovery also matters.

UP-29 already established exact real-orthogonal equivalence, so no outcome here supports uniquely complex or quantum computation claims.

## Plain speak

UP-30 told us that going from six repeated resonances to six different ones breaks the system. UP-31 told us those six repeated resonances do not need to line up in the same coordinates. UP-32 now turns the knob gradually: six repeated copies, then four, three, two, and finally one.

If performance falls as that number falls, we have much stronger evidence that the repeated dynamical structure itself is carrying useful memory geometry.

## Authority boundary

UP-32 is mathematical research only. It does not invoke a language model, activate Wingless in production, alter ckb-plane/KTRADE authority, access brokers or credentials, modify accepted refs, or perform promotion/deployment.


## Windows qualification result

Autonomous Windows qualification on workflow run `35875796731` completed successfully at source head `278f5b719d744be15b68a39f21f0120d9bb0a3c0`.

All four adjacent steps satisfied the preregistered +0.03 ordering tolerance for both held-out accuracy and mutable commit accuracy. The curve was:

- multiplicity 6: held 1.0, commit 1.0;
- 4+2: held 0.996826171875, commit 0.9908854166666666;
- 3+3: held 0.994873046875, commit 0.9947916666666666;
- 2+2+2: held 0.986328125, commit 0.96484375;
- singleton: held 0.90185546875, commit 0.34375.

The overall preregistered dose-response gate remained negative because its endpoint material-effect clause required both static and mutable endpoint drops to be at least 0.25. The mutable endpoint drop was 0.65625, but the held-out endpoint drop was 0.09814453125.

Interpretation: this is ordered graded evidence for repeated-spectrum structure, but not a gate-positive dose-response result. It also reveals a remaining confound: the tested partitions change both maximum multiplicity and total commutant capacity `sum(m_i^2)`. UP-33 should compare equal-capacity partitions with different maximum multiplicity.
