# Wingless UP-46 — objective decomposition audit

Status: staged diagnostic only. Qualification is intentionally not triggered while the UP-44 FIXA seed-matched causal control is active.

Scientific source state: UP-45 seal `48fb886f517db0cb21fdbc757cd7bc9500ed9f59`.

## Frozen question

UP-45's reversible convex proximal operator created exact fusion at step 29, but the unchanged task objective selected non-fused step 24.

UP-46 asks which measured component explains that rejection and whether the rejection agrees with hard recurrent capability.

This is a diagnostic comparison, not a new optimization.

## Frozen states

Selected UP-45 step 24:

`[-0.05841298929708432, -0.03887820934987345, 0.0051658188506906524, -0.03702321187042985, 0.07759973733795059, 0.05154885432874637]`

Exact-fused UP-45 step 29:

`[-0.045664698831175514, -0.04501676812740249, 0.007218659931181362, -0.045664698831175514, 0.07894943527059196, 0.05017807058798021]`

Step 29 contains the exact fused pair `[0,3]` and exact capacity 8. Step 24 has exact capacity 6.

## Frozen mechanics

For both states, independently recompute:

- normalized phase score;
- mean correct value probability;
- mean correct relation probability;
- harmonic task score;
- soft capacity;
- resource penalty;
- final smooth objective;
- held-out static accuracy;
- mutable commit accuracy;
- exact final-table accuracy;
- relational query accuracy.

Use the untouched reference evaluator and the same UP-45 configuration.

No optimizer is run.
No lambda is changed.
No state is perturbed or selected.
No checkpoint search is performed.
No held-out result influences state selection.

Report all deltas as `step29 exact-fused - step24 selected`.

## Frozen interpretation

### Case 1 — smooth objective rejects fusion and hard capability also worsens

The objective is directionally aligned with hard capability for this exact-fusion event. Exact equality at this state is not beneficial enough to justify stronger fusion pressure. Next work should investigate what structural change is required for symmetry to support recurrence rather than forcing equality.

### Case 2 — smooth objective rejects fusion but hard capability improves

The surrogate objective is misaligned with the capability we ultimately care about. Identify which smooth component contributes the rejection, then redesign only that surrogate component in a separately preregistered experiment. Do not increase lambda.

### Case 3 — smooth objective improves at exact fusion but step 29 was not selected

That would indicate checkpoint scheduling rather than objective preference caused the apparent rejection, because step 29 was not one of the selection checkpoints. The next experiment must test a preregistered denser checkpoint schedule on a replay/frozen trajectory, not rerun optimization with result-informed choices.

### Case 4 — soft task components improve but resource penalty makes total objective worse

The resource term and exact-symmetry geometry are misaligned. Audit the soft-capacity proxy against exact commutant capacity before changing its coefficient.

### Case 5 — relation probability dominates the rejection

The relational surrogate becomes the leading downstream seam. Test whether exact symmetry damages the relation head itself or merely the probability representation supplied to it.

### Case 6 — value probability dominates the rejection

The mutable-memory surrogate becomes the leading downstream seam. Test decoder probability calibration/readout at the exact-fused geometry without changing fusion pressure.

### Case 7 — phase score dominates the rejection

The learned phase observable is not invariant to the exact structure the regularizer creates. Audit observable selection/transport at the two frozen states before any optimizer change.

## Seed-matched causal-control dependency

UP-46 is an internal decomposition of the already observed UP-45 trajectory and remains mathematically well-defined regardless of the UP-44 FIXA outcome.

However, its interpretation must not be used to claim that `sqrt(probability)` itself caused the UP-43 -> UP-44 geometry alignment unless the seed-matched UP-44 FIXA control independently restores that causal claim.

Qualification must therefore wait for the active seed-matched control to seal.

## Plain speak

UP-45 made exact equality, but the scoring system chose an earlier non-equal state.

UP-46 freezes those two states and takes the scoring system apart piece by piece. It tells us whether equality was genuinely worse, whether our score misunderstood a useful state, or whether the penalty term itself discouraged the structure we were trying to learn.

Nothing is tuned during this audit.
