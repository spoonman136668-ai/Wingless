# Wingless UP-19: direct anonymous Gram readout

Status: Windows-qualified scientific negative; not activated, promoted, or connected to ckb-plane.

Parent research result: UP-18 qualified scientific negative sealed at `070b759e6284830726e758e1b2ab9aae245ea5a5`.

## Question

UP-18 proved that the anonymously mixed 96-dimensional state still contains the full task-relevant representation, but both unconstrained and orthogonally constrained demixers failed to rediscover the historical memory/anchor/pilot roles from task loss.

UP-19 asks a different question:

**Does the observer need to reconstruct those old semantic channels at all?**

Instead of learning a demixer, UP-19 reads the anonymous state directly.

## Anonymous relational representation

The recurrent state remains the same single 96-dimensional state from UP-17/UP-18:

- six anonymous 16-dimensional channels;
- fixed dense orthogonal channel mixing;
- no semantic channel labels exposed;
- no mixer or inverse exposed to the learner.

For each state, UP-19 computes the complete Hermitian Gram geometry among the six anonymous channels:

- six channel self-inner-products;
- fifteen complex cross-inner-products.

This produces 36 real features.

The same transport block acts on every anonymous channel. Under unitary propagation, these pairwise inner products are therefore invariant with depth.

## Why a quadratic lift

The historical frame readout compares the phase of two semantic correlations.

After an unknown fixed linear channel mix:

1. each historical semantic correlation is a linear function of the anonymous Gram entries;
2. the phase-comparison numerator is the product of two such correlations;
3. that numerator is therefore quadratic in the anonymous Gram entries.

UP-19 consequently tests two readouts:

- a 36-feature linear Gram baseline;
- a complete degree-2 lift containing the 36 raw features plus all 666 unique quadratic products, for 702 features total.

Both use the existing deterministic linear-softmax learner. There is no learned demixer, hidden-mixer target, inverse transport, depth input, prototype lookup, or semantic channel location.

## Training

Unlike earlier observer experiments, UP-19 trains only at canonical depth before transport.

Training uses:

- all 128 balanced training tables;
- memory noise amplitude 0.05;
- memory-only global-phase nuisance;
- two deterministic noisy trials per table;
- linear baseline: 600 softmax steps at learning rate 1.0;
- quadratic readout: 400 softmax steps at learning rate 0.8.

The same trained heads are then evaluated after transport.

This makes unseen-depth evaluation stronger: depth is never represented in the training set.

## Held-out evaluation

Evaluation uses:

- all 128 disjoint held-out tables;
- depths 32, 128, 512, 1024;
- two noisy trials per table/depth.

The matched non-unitary arm uses the **same quadratic heads** learned on canonical states. It does not receive separately retrained heads.

## Mutable integration

If the quadratic unitary observer is usable, it is carried into:

- 48 scenarios;
- 16 writes each;
- 768 commit opportunities;
- held-out depths only;
- memory noise 0.05;
- memory-only global phase nuisance;
- explicit irreversible overwrite/re-encode boundary;
- learned relational query head.

No semantic channel reconstruction occurs at runtime.

## Scientific gates

**Feature invariance**

- maximum unitary quadratic-feature drift at depth 1024 <= 1e-12.

**Direct quadratic readout**

- canonical quadratic training accuracy >= 0.99.

**Unseen depth**

- unitary held-out accuracy >= 0.99;
- maximum unitary norm drift <= 1e-12.

**Mutable integration**

- commit decode >= 0.99;
- exact final table >= 0.95;
- relational query >= 0.95;
- maximum unitary norm drift <= 1e-12.

The linear Gram result and non-unitary result are diagnostics, not scientific pass shortcuts.

## Interpretation boundary

A positive UP-19 result would show that explicit recovery of the historical memory/anchor/pilot basis is unnecessary: a task learner can operate directly on relational geometry of the anonymous recurrent state.

That would remove the blind-demixing bottleneck from the architecture rather than solve it.

A negative UP-19 result would mean that either the complete degree-2 anonymous Gram representation is insufficient under the learned decoder/data budget, or that the old semantic factorization carries information in a form not captured by this direct observer. The result must be preserved rather than tuned until green.

The state still exposes six anonymous 16-dimensional blocks. A positive result would justify a later experiment that mixes channel and coordinate dimensions together.

This experiment still does not establish phase-specific or unitary superiority. A real orthogonal transport comparator remains required before any stronger causal claim.

## Plain speak

UP-17 and UP-18 tried to unscramble the box so it looked like our old labeled drawers again.

UP-19 stops doing that.

It gives the learner the full pattern of how the anonymous pieces relate to one another and asks it to solve the memory task directly.

If that works, Wingless does not need to figure out which hidden piece used to be called “memory,” “anchor,” or “pilot.” It can use the internal geometry without our labels.

## Authority boundary

UP-19 is mathematical research only. It does not invoke a language model, activate Wingless in production, alter ckb-plane/KTRADE authority, access brokers or credentials, modify accepted refs, or perform promotion/deployment.


## Windows qualification result

Autonomous Windows qualification on workflow run `35851476994` completed successfully at source head `7eb608cf291adb2dcb7718582d93fe148b0fbd56`.

The unitary anonymous Gram representation passed its depth-invariance gate: maximum quadratic-feature drift at depth 1024 was 1.389999226830696e-13, with state norm drift 3.3306690738754696e-14. The matched non-unitary feature drift was 1.077751562699128e25.

The learned readouts were scientific negatives. The 36-feature linear head reached 0.26220703125 held-out accuracy. The complete 702-feature degree-2 head reached only 0.2529296875 training accuracy and 0.25146484375 held-out accuracy. Mutable integration consequently failed.

Interpretation: the anonymous relational representation is stable, but the existing gradient softmax learner does not extract the task signal from it. Because the established semantic decision function is algebraically contained in the degree-2 feature span after fixed invertible channel mixing, the next experiment isolates the optimizer by keeping the exact same representation and replacing gradient softmax with a deterministic closed-form ridge readout.
