# Wingless UP-3: depth retention stress

Status: Windows-qualified research primitive; not activated, promoted, or connected to ckb-plane.

Parent research qualification: UP-2 Windows qualification sealed at `fa5e9e80737d90e9c1ad901319bcd2b045c3c94f`.

## Question

Does the unitary substrate preserve usable latent information under repeated propagation, distractor dimensions, and bounded perturbation as depth increases?

UP-3 is a stress/measurement experiment. It does not require the matched non-unitary control to fail.

## Construction

The latent state expands from four to sixteen complex dimensions:

- four dimensions contain the four-way Walsh relative-phase class code from UP-2;
- twelve dimensions contain deterministic class-independent distractor state;
- every clean state is normalized.

A deterministic butterfly-style propagation block contains 32 pair couplings spanning strides 1, 2, 4, and 8. The exact same pair topology and scalar values are used by both paths.

The paths are:

1. **unitary** — complex Givens rotations;
2. **matched non-unitary** — the same unconstrained mixer used by the UP-2 control.

The block is repeated to depths:

`0, 1, 8, 32, 128`.

Thirty-two noisy held-out samples cover all four classes, multiple unseen global phases, and deterministic perturbations across every latent dimension.

## Measurements

At every depth the probe records:

- global-phase-invariant nearest-prototype classification accuracy;
- maximum norm drift from the normalized clean-state baseline;
- maximum error in the prototype fidelity/Gram geometry;
- maximum perturbation gain relative to the original input perturbation.

At maximum depth the unitary path is also inverted all the way back to the input and reports round-trip error.

## Pass conditions

For the unitary path at every tested depth:

- classification accuracy remains 1.0;
- maximum norm drift <= 1e-12;
- maximum Gram/fidelity error <= 1e-12;
- perturbation gain remains within 1e-12 of 1.0.

Additionally:

- unitary depth-128 forward/inverse round-trip error <= 1e-11;
- all matched-control metrics remain finite;
- both paths are identical at depth zero;
- the complete result is deterministic;
- the full existing Wingless regression remains green after advisory ICE rebuild.

No gate requires the non-unitary control to degrade. Its observed behavior is evidence, not a scripted failure condition.

## Interpretation boundary

A pass would establish that the current unitary primitive can carry a structured relative-phase code through thousands of individual local rotations without losing classification geometry or amplifying bounded perturbations beyond floating-point tolerance.

That is a memory/transport result, not evidence of language reasoning. A positive UP-3 result permits UP-4: a sequential memory/composition task in which state must be updated repeatedly and queried later, with unitary and matched recurrent baselines trained under comparable budgets.

## Authority boundary

UP-3 is mathematical research only. It does not register an inference backend, invoke a language model, activate a Wingless worker/listener, execute model-selected tools, or alter ckb-plane queue, retry, workspace, acceptance, promotion, broker, credential, or production authority.


## Windows qualification

Authoritative operator proof on 2026-09-22 against source head `b4115bdcf58125b856851ad0fbd90355b40de86c` passed the focused UP-3 suite and complete Wingless regression.

At depth 128 the unitary path retained 1.0 classification accuracy with maximum norm drift 7.771561172376096e-15, maximum Gram/fidelity error 2.7755575615628914e-15, perturbation gain 1.0000000000000098, and complete forward/inverse round-trip error 9.92884257223707e-15.

The matched non-unitary path also retained 1.0 classification accuracy, but its maximum norm drift reached 26.751397615116673, maximum Gram/fidelity error reached 0.0781818758548447, and perturbation gain reached 5.541839982219787.

Interpretation remains bounded: UP-3 establishes a deep transport/stability distinction under this matched stress construction. It does not establish language reasoning or a general performance advantage.
