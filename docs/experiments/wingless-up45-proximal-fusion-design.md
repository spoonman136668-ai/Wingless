# Wingless UP-45 design: reversible convex exact fusion

Status: design staging only. This branch intentionally does not match the auto-running UP workflow and does not consume the Windows research runner.

Scientific parent: UP-44 seal `584373c3bd32fc305cffffeec16ad5e7ee07878c`.

## Motivation

UP-44 showed that the corrected square-root probability representation makes the task objective prefer a reversible symmetry-basin checkpoint. The selected step 24 had positive objective gain and provisional capacity 8, with nearest offset gap `0.00009309022967083583`, but the offsets were not exactly equal.

The next scientific question is therefore not whether the optimizer can reach a useful neighborhood. It is whether a principled exact-fusion operator can convert task-supported near-symmetry into exact equality without sticky snapping, target groups, or a partition menu.

## Proposed operator

Use a proximal-gradient update with the complete-graph convex clustering penalty

`R(x) = sum_{i<j} |x_i - x_j|`.

After the ordinary rolling task-gradient update proposes `y`, compute

`x = prox_{lambda R}(y)`

by solving

`min_x 0.5 * ||x-y||_2^2 + lambda * sum_{i<j}|x_i-x_j|`.

For six scalar spectral offsets this proximal problem has an exact deterministic solution:

1. sort the six proposed offsets;
2. in sorted order, use the identity

   `sum_{i<j}|x_j-x_i| = sum_i (2*i-n-1) x_i`

   for nondecreasing `x`;
3. shift each sorted observation by the corresponding linear penalty coefficient;
4. solve the resulting nondecreasing least-squares problem with deterministic pool-adjacent-violators isotonic regression;
5. unsort to original channel order;
6. apply the existing zero-mean/RMS normalization.

Isotonic pooling creates exact equal blocks when the convex optimum requires fusion. No explicit pair threshold determines membership.

## No new tuned free parameter

Freeze

`lambda = directOffsetLearningRate * directOffsetResourcePrice`

which under the existing constants is

`0.002 * 0.02 = 0.00004`.

This couples fusion strength to two already-frozen quantities rather than introducing a new value selected from UP-44's observed gap.

No lambda sweep is permitted in R1.

## Reversibility

There is no union-find memory and no sticky state.

Every optimization step starts from the current six offsets, applies the ordinary task-gradient proposal, and solves the convex proximal problem from scratch.

A pair fused on one step may separate on a later step if the subsequent task gradient makes separation preferable enough to overcome the regularizer. Exact equality is therefore an optimization result, not an irreversible architectural decree.

## Frozen inheritance from UP-44

Keep unchanged:

- square-root probability soft-memory representation;
- no global soft-memory renormalization;
- rolling seven-observation full-rank directional memory;
- 48 optimization updates;
- plus/minus perturbation 0.002;
- learning rate 0.002;
- maximum coordinate update 0.004;
- checkpoints 0, 12, 24, 48;
- 96/32 training-only split;
- untouched 128-table true held-out pool;
- phase/value/relation harmonic task objective;
- resource price 0.02;
- soft-capacity tau 0.0025;
- no known UP-35 target geometry;
- no finished partition menu;
- no hard merge search;
- no pair-attraction field;
- no held-out selection;
- unchanged hard recurrent final evaluation and capability gates.

## Exact-fusion reporting

At every step report:

- exact equality groups after proximal update;
- exact commutant capacity;
- whether any exact fusion exists;
- whether a previously exact pair later separates;
- proximal objective before/after;
- maximum proximal KKT residual or equivalent exact-solver residual.

The old distance-based provisional geometry may remain as a diagnostic only, not as authority for exact fusion.

## Frozen interpretation

- If the task-selected checkpoint has exact fusion, positive task-objective gain, and materially better hard capability than UP-44, continue testing whether exact symmetry is causally improving recurrent capability.
- If exact fusion appears and the smooth objective improves but hard capability remains essentially at UP-44 levels, the remaining bottleneck is downstream of approximate-vs-exact symmetry; investigate the surrogate/readout/capability seam rather than increasing fusion pressure.
- If proximal regularization causes exact fusion but the task objective selects a non-fused checkpoint, the task does not currently support paying for exact equality at this regularization strength. Seal the negative; do not tune lambda after seeing the result.
- If no exact fusion occurs, seal the negative and test whether the soft-capacity/resource term and convex clustering penalty are misaligned before considering any stronger coefficient.
- If the proximal path changes unrelated scientific mechanics or fails deterministic equivalence tests, treat that as an implementation defect, not a scientific result.

No post-result lambda change, gate relaxation, horizon extension, or threshold tuning is allowed.

## Plain speak

UP-44 got extremely close to making two internal channels truly identical, and the task score actually liked that region.

UP-45 should not simply say, "they are close enough, glue them forever." Instead it adds a small mathematical preference for simpler equal structure. The optimizer is then allowed to decide whether equality is worth it.

If equality is useful, the convex solver will naturally collapse values together. If later task evidence wants them apart again, they can separate because nothing is permanently locked.
