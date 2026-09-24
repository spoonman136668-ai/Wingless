# Wingless post-UP-45 frozen decision tree

Status: design-only staging. No Windows qualification is triggered from this branch.

Parent experiment source: UP-45 source head `1c402ca58029941d1a5f9ebede2b9fd4fc1a5d0c`.

This decision tree is frozen before reading the UP-45 scientific result.

## Branch A — selected checkpoint has exact fusion and capability gates pass

Interpretation: the task objective, convex exact-fusion operator, and hard recurrent capability are aligned.

Next experiment: test persistence/generalization of the exact fused structure under a fresh disjoint scenario schedule and restart reconstruction. Do not change lambda, gates, horizon, representation, or optimizer first.

Question: does the task-selected exact symmetry survive restart and preserve capability on a new preregistered scenario schedule without re-optimizing the geometry?

## Branch B — selected checkpoint has exact fusion, objective gain is positive, but capability gates fail

Interpretation: exact equality is being selected by the smooth task objective, but exact symmetry alone is not sufficient for the hard recurrent capability target.

Next experiment: causal exact-vs-near ablation at the selected UP-45 geometry.

Freeze the selected exact fused geometry as the center. Construct a matched near-symmetry control by splitting only each exact fused block symmetrically while preserving zero mean and RMS. The split magnitude must be derived from an already-frozen scale before running the ablation; no fit to UP-45 outcome is permitted. Use the same trained/readout evaluation protocol and true held-out discipline.

Question: does exact equality itself improve hard recurrent capability relative to an otherwise matched near-symmetry state?

If exact and near perform equivalently, move downstream to the decoder/readout/surrogate seam rather than increasing fusion pressure.

If exact materially changes hard capability, investigate how the hard recurrence uses commutant multiplicity before changing the optimizer.

## Branch C — exact fusion occurs somewhere in the trajectory, but the task-selected checkpoint is not exactly fused

Interpretation: the convex operator can create equality, but the frozen task objective does not select paying for it.

Next experiment: objective-decomposition audit only. Compare phase, value-probability, relation-probability, soft-capacity penalty, and hard recurrence metrics at:
- the best selected non-fused checkpoint;
- the best objective-scored exactly fused checkpoint observed during the preregistered trajectory.

No lambda change, horizon extension, or new optimization is permitted.

Question: which objective component causes task selection to reject exact equality?

## Branch D — no exact fusion occurs anywhere

Interpretation: the derived convex regularizer is too weak or structurally misaligned to induce equality under the current task-gradient scale, but the coefficient must not be tuned after seeing the result.

Next experiment: scale audit only. Quantify, over the frozen UP-45 trajectory:
- task-gradient proposal magnitude;
- proximal displacement magnitude;
- pairwise penalty reduction;
- nearest pair gaps;
- exact solver residuals.

Question: is the lack of fusion explained by the derived regularizer being negligible relative to task-gradient motion, or by a geometric mismatch between the complete-graph penalty and task-supported symmetry?

No stronger lambda may be tried until that audit is sealed.

## Branch E — implementation or determinism defect

If the solver violates its exact-prox checks, deterministic replay, scientific boundaries, or unrelated mechanics, treat the run as an infrastructure/implementation defect. Repair only the defect and repeat the same preregistered UP-45 experiment. Do not reinterpret it scientifically.

## Global constraints

For every successor:
- untouched reference evaluator remains authoritative;
- no cached performance evaluator is adopted from HARNESS-PERF-R1;
- no post-result gate relaxation;
- no held-out selection;
- no UP-35 target geometry injected;
- no sticky unions;
- no production/CKB/ckb-plane/KTRADE/Nemotron integration;
- no activation or promotion.

## Plain speak

The next move is already decided by what kind of result UP-45 produces.

If exact symmetry works, test whether it persists.
If exact symmetry is selected but capability is still weak, isolate whether exactness itself matters.
If the task rejects exactness, inspect why.
If exactness never appears, measure whether the regularizer was too small or pointed the wrong way.
