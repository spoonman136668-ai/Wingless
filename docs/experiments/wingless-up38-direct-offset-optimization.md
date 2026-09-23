# Wingless UP-38: direct continuous spectral-offset optimization

Status: research branch only; not activated, promoted, or connected to ckb-plane.

Parent research result: UP-37 qualified near-positive gate-negative sealed at `79d2638496cfa94591d332c6369ad66b1062402f`.

## Question

UP-37 showed that adaptive continuous symmetry formation can grow useful exact degeneracy from capacity 6 to capacity 12 and recover nearly all static and mutable capability, missing only the frozen relational-query gate.

UP-38 removes the hand-designed pair-attraction field and the hand-designed continuous flow rule.

It asks:

**Can the six spectral offsets themselves be optimized directly under a smooth task/resource objective, starting from independent directions, such that useful exact symmetry emerges without a partition menu, hard merge search, or pair-affinity field?**

## Direct architecture parameters

The trainable architecture parameters are the six anonymous spectral offsets themselves.

They remain constrained to:

- zero mean;
- RMS magnitude 0.05 radians per step;
- the same fully mixed 96-dimensional recurrent state;
- norm-preserving transport;
- no hidden semantic channel assignments.

The optimizer starts from the six distinct singleton offsets used throughout the multiplicity-ablation lineage.

## Smooth task objective

UP-38 uses the validation mean phase cosine from the accepted transport-discovered observer as its smooth task signal.

No true held-out data participates.

## Smooth resource objective

The discrete commutant-capacity penalty used by UP-34/UP-35 is replaced by a smooth approximation:

`C_soft = 6 + 2 * sum_{i<j} exp(-(delta_i-delta_j)^2 / (2 tau^2))`

with `tau = 0.0025`.

This has the correct limiting behavior:

- well-separated singleton offsets approach capacity 6;
- six exactly equal offsets equal capacity 36;
- partial near-degeneracies interpolate smoothly.

The resource penalty remains:

`0.02 * C_soft / 36`.

The optimized objective is:

`mean_phase_cosine - resource_penalty`.

Thus the task must earn any movement toward higher symmetry.

## Optimizer

UP-38 uses deterministic SPSA-style direct optimization:

- 12 optimization steps;
- perturbation magnitude 0.002;
- learning rate 0.002;
- maximum per-coordinate update 0.004;
- one deterministic zero-mean perturbation direction per step;
- two objective evaluations per gradient estimate;
- all six offsets updated simultaneously.

No pair-attraction matrix is constructed.

No discrete merge candidate is evaluated.

## Exact symmetry formation

After each direct continuous parameter update:

- offsets are re-centered to zero mean;
- RMS is restored to 0.05;
- any already fused group remains exactly fused;
- if two current group centroids come within the existing frozen 0.0025 tolerance, they are projected to exact equality and remain sticky.

The sticky projection is retained explicitly as scaffolding. UP-38 removes the hand-designed attraction dynamics, not yet the exact-fusion projection.

## Data discipline

The same balanced training-side split is retained:

- 96 inner-fit tables;
- 32 inner-validation tables;
- untouched 128-table true held-out pool.

Optimization and checkpoint selection use only the training-side split.

## Model selection

Each optimized checkpoint is scored by the same smooth task/resource objective.

The best checkpoint is frozen before true held-out evaluation.

## Scientific gates

Preflight:

- balanced inner split;
- no true held-out selection.

Optimization:

- at least one optimized checkpoint must beat initialization;
- exact symmetry must emerge;
- selected commutant capacity must be > 6 and < 36.

Final capability:

- held-out accuracy >= 0.985;
- mutable commit accuracy >= 0.95;
- exact final table accuracy >= 0.90;
- relational query accuracy >= 0.95.

Retention versus full capacity:

- held-out loss <= 0.02;
- mutable-commit loss <= 0.05.

## Interpretation boundary

A positive result would remove a major scaffold from UP-37: task-relevant symmetry could emerge by directly optimizing the transport parameters under a smooth objective rather than through a hand-designed attraction field.

A negative result would separate two possibilities:

- the direct optimizer is insufficient despite a useful smooth objective;
- the smooth mean-phase objective does not adequately represent the mutable/relational task pressure that the discrete and adaptive experiments captured.

A likely negative-result successor would therefore compare direct optimization under richer smooth task losses rather than return to manual architecture selection.

UP-38 still retains:

- six predefined spectral directions;
- fixed RMS normalization;
- a fixed optimizer;
- sticky exact-fusion projection;
- the learned phase target.

No result establishes uniquely complex or quantum-computation advantage; UP-29 already established exact real-orthogonal equivalence.

## Plain speak

UP-35 let the system choose which walls to knock down.

UP-36 and UP-37 replaced the hammer with attraction forces.

UP-38 removes those forces too.

Now the positions of the walls themselves are the parameters being optimized. The task directly pushes on those positions while a resource cost resists unnecessary symmetry.

If some walls still converge and the resulting structure works on unseen data, we will have removed another layer of hand-designed developmental machinery.

## Authority boundary

UP-38 is mathematical research only. It does not invoke a language model, activate Wingless in production, alter ckb-plane/KTRADE authority, access brokers or credentials, modify accepted refs, or perform promotion/deployment.
