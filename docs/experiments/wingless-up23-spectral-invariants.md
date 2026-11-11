# Wingless UP-23: full-latent transport spectral invariants

Status: Windows-qualified scientific negative; not activated, promoted, or connected to ckb-plane.

Parent research result: UP-22 qualified scientific negative sealed at `e25b0312c1ee15ecda95cd56a105efa3a73bf060`.

## Question

UP-21 showed that the learned two-dimensional phase-code target can decode the anonymously mixed one-state representation perfectly.

UP-22 removed the remaining 6 × 16 block structure by mixing all 96 coordinates and conjugating the transport into that basis. The transport itself remained equivalent, but a coordinate-basis Hermitian observer failed to generalize across unseen depths.

UP-23 asks:

**Can the fully mixed 96-dimensional state be read through invariants defined by its own transport dynamics rather than by coordinate blocks?**

## Spectral-moment representation

Let the full-latent one-step unitary transport be `V`, and the current full-latent state be `psi`.

UP-23 computes:

`m_k = <psi, V^k psi>`

for `k = 1..16`.

Each complex moment contributes real and imaginary parts, for 32 real features.

For a unitary trajectory `psi_d = V^d psi`:

`<psi_d, V^k psi_d> = <psi, V^k psi>`.

Therefore these features are depth-invariant without:

- channel boundaries;
- hidden mixer access;
- unmixing;
- explicit depth;
- inverse transport;
- prototype lookup.

The transport operator itself is explicitly used as the reference structure. That architectural assumption is surfaced in the result schema.

## Training

The learner trains only at depth 0.

Training uses:

- all 128 balanced training tables;
- memory noise 0.05;
- memory-only global-phase nuisance;
- four deterministic noisy trials per table;
- 16 complex spectral moments / 32 real features;
- deterministic dual ridge regression to the UP-13 learned two-dimensional phase code;
- ridge lambda 1e-6;
- fresh two-dimensional softmax calibration heads.

The matched non-unitary arm receives the same training budget and uses its own conjugated one-step transport to construct its spectral moments.

## Held-out evaluation

Evaluation uses:

- all 128 disjoint held-out tables;
- depths 32, 128, 512, 1024;
- two noisy trials per table/depth.

No held-out depth is represented during training.

## Mutable integration

The learned unitary spectral model then runs:

- 48 scenarios;
- 16 writes each;
- 768 commit opportunities;
- held-out depths only;
- explicit irreversible overwrite/re-encode boundary;
- learned relational-query head.

Every recurrent segment remains a single fully mixed 96-dimensional state.

## Scientific gates

**Moment invariance**

- maximum unitary spectral-feature drift at depth 1024 <= 1e-10.

**Spectral learning**

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

A positive result means the fully mixed recurrent state can be decoded from basis-independent dynamical invariants without exposing historical channel blocks or the hidden coordinate mixer.

A negative result means scalar spectral moments of the state are insufficient to carry the task signal, even though they are depth-stable. The next experiment should then test a richer transport-defined invariant algebra rather than restoring visible blocks.

This experiment does not establish phase-specific or unitary superiority. A real orthogonal transport comparator remains mandatory before stronger causal claims.

## Plain speak

UP-22 scrambled all 96 internal coordinates so thoroughly that the old “drawers” disappeared, but our reader got lost because it was still looking at raw coordinates.

UP-23 stops looking at where information sits.

Instead it asks how the state resonates with its own dynamics.

Those resonance signatures stay the same no matter how long the unitary state evolves. If they still identify the memory values, Wingless can use the fully mixed state without needing internal drawer labels at all.

## Authority boundary

UP-23 is mathematical research only. It does not invoke a language model, activate Wingless in production, alter ckb-plane/KTRADE authority, access brokers or credentials, modify accepted refs, or perform promotion/deployment.


## Windows qualification result

Autonomous Windows qualification on workflow run `35864220728` completed successfully at source head `fb663a9a885a1dc8337e8792b2d8ad19fa021348`.

The spectral-moment invariance premise passed: maximum moment drift was 1.5197149094703377e-12. The learned observer nevertheless reached only 0.4091796875 training accuracy and 0.3720703125 held-out accuracy, so the negative is representational rather than an optimizer-only failure.

Interpretation: scalar moments of one transport operator preserve only aggregate weight across its eigenspaces. Because the full-latent transport retains the six-fold multiplicity inherited from the historical composite transport, these moments discard orientation inside degenerate multiplicity spaces. UP-24 therefore tests a generic algebra of observables that commute with the transport and resolve that missing multiplicity-space information without runtime unmixing or visible channel blocks.
