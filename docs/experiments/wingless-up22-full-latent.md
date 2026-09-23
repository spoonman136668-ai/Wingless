# Wingless UP-22: full 96-dimensional latent mixing

Status: research branch only; not activated, promoted, or connected to ckb-plane.

Parent research result: UP-21 qualified scientific pass sealed at `3c2181a95f7834b43ec58e5bf63abcb899d711b0`.

## Question

UP-21 removed the need to recover semantic memory/anchor/pilot channel identities, but its recurrent state still exposed six anonymous 16-dimensional blocks.

UP-22 asks whether that remaining block layout can be removed entirely.

## Full latent basis change

After the historical components are packed into one 96-dimensional complex state, UP-22 applies a fixed 96-point unitary DFT basis transform.

Every output coordinate receives equal-magnitude contribution from all 96 input coordinates:

`|Q[i,j]| = 1/sqrt(96)`.

Therefore no contiguous runtime coordinate block corresponds to memory, anchor, a pilot, an anonymous channel, or one original 16-dimensional transport space.

The learner is not given `Q`.

## Conjugated transport

The original one-step composite transport is:

`T = I_6 tensor U`.

UP-22 evolves the fully mixed latent state directly under:

`V = Q T Q^H`.

Depth operators `V^d` are precomputed for the preregistered train, held-out, and mutable depths.

There is no runtime unmix/re-mix cycle.

The matched non-unitary control receives its own exactly conjugated operator:

`V_control = Q T_control Q^H`.

## Oracle equivalence diagnostic

Only for diagnosis, UP-22 applies `Q^H` after latent evolution and checks that the result equals the historical packed-state evolution.

This unmix is never used by the learned observer or mutable runtime.

Gate:

- maximum recovered-state L2 error <= 1e-10.

## Direct full-state observer

Because no channel partition is visible, the observer uses the complete Hermitian quadratic representation of the 96-dimensional latent state:

- 96 self magnitudes;
- 4,560 complex pairwise correlations;
- 9,216 real features total.

This representation uses coordinate relations only. It does not group coordinates into historical channels.

Each entity learns:

1. a deterministic ridge map from 9,216 features to the two-dimensional UP-13 phase-code target;
2. a fresh 2D softmax calibration head.

The mixer, oracle unmix, semantic channel boundaries, explicit depth, and prototypes are absent from learning.

## Training

To prevent the learner from merely fitting one coordinate orientation at one depth, training includes two transport depths:

- depth 0;
- depth 216.

Training uses:

- all 128 balanced training tables;
- one deterministic noisy trial per table/depth;
- memory noise 0.05;
- memory-only global phase nuisance;
- ridge lambda 1e-5.

This gives 256 ridge rows.

## Held-out evaluation

Evaluation uses:

- all 128 disjoint held-out tables;
- depths 32, 128, 512, 1024;
- two noisy trials per table/depth.

The matched non-unitary path uses the exact same learned regressors and calibration heads.

## Mutable integration

The unitary full-latent model then runs:

- 48 scenarios;
- 16 writes each;
- 768 commit opportunities;
- held-out depths only;
- explicit irreversible overwrite/re-encode boundary;
- learned relation head.

Every recurrent segment exists and evolves only in the fully mixed 96-dimensional basis.

## Scientific gates

**Conjugation equivalence**

- diagnostic recovered-state error <= 1e-10.

**Full-latent learning**

- training accuracy >= 0.99.

**Unseen depth**

- held-out accuracy >= 0.99;
- maximum unitary norm drift <= 1e-10.

**Mutable integration**

- commit decode >= 0.99;
- exact final table >= 0.95;
- relational query >= 0.95;
- maximum unitary norm drift <= 1e-10.

The minimum row participation ratio of the fixed mixer is reported as a structural diagnostic; the DFT construction should be 96.

## Interpretation boundary

A positive result means the runtime representation no longer needs visible memory/anchor/pilot channels or even anonymous 16-dimensional channel blocks. Task decoding can operate from relational structure in one fully mixed 96-dimensional latent state.

This does not mean the encoder has stopped assembling the historically learned components before the basis transform. It removes runtime coordinate factorization, not the developmental/encoding lineage that created the state.

A negative result must distinguish transport-conjugation defects from observer generalization failure. It is not a reason to expose hidden block labels again.

No result here establishes phase-specific or unitary superiority. A real orthogonal transport comparator remains mandatory before stronger causal claims.

## Plain speak

UP-21 proved Wingless could use the scrambled box without knowing which drawer was which.

UP-22 removes the drawers themselves.

Every one of the 96 internal coordinates becomes a mixture of everything else, and the dynamics are rewritten so the mixed state evolves directly in that basis.

If this passes, there is no runtime place we can point to and say “this block is memory” or “this block is the compass.”

## Authority boundary

UP-22 is mathematical research only. It does not invoke a language model, activate Wingless in production, alter ckb-plane/KTRADE authority, access brokers or credentials, modify accepted refs, or perform promotion/deployment.
