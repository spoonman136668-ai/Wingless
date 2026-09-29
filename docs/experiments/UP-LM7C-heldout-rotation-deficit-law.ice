UP-LM7C PREREGISTRATION — HELD-OUT ROTATION DEFICIT-LAW TRANSFER

Parent UP-LM6C qualified EXACT_DEFICIT_LAW on all 108 frozen rows:
writes = max(0, pre_countdown - target)
and post_countdown = target,
with zero mismatches across profiles {deferred_only,layout_only,hybrid_min}, rotations {5,13}, permutations {identity,reverse,rotate2}, and six arms.

Question: is the exact deficit law invariant to held-out rotation values when every other substrate dimension and the correction policy are frozen?

Freeze:
- profiles={deferred_only,layout_only,hybrid_min};
- permutations={identity,reverse,rotate2};
- six arms produced by the existing uplm2xArms construction;
- existing uplm3tTarget target rule;
- existing prepressure write loop and 32-write ceiling;
- no policy change;
- no live activation.

Change exactly one substrate dimension:
- replace the parent rotation set {5,13} with held-out rotations {8,21}.

For every resulting row compute:
deficit=max(0, pre_countdown-target)
write_error=writes-deficit
post_target_error=post_countdown-target

The parent UP-LM6C anchor must reproduce before evaluating the held-out rotations.

Report:
- row count;
- exact-law row count;
- mismatch row count;
- min/max write_error;
- min/max post_target_error;
- mismatch counts by held-out rotation and profile.

Classification:
HELDOUT_ROTATION_EXACT_DEFICIT_LAW if every held-out-rotation row has write_error=0 and post_target_error=0.
ROTATION_DEPENDENT_DEFICIT_RESIDUAL if the parent anchor reproduces and any held-out-rotation row has nonzero write_error or post_target_error.
ANCHOR_NOT_REPRODUCED if UP-LM6C no longer reproduces EXACT_DEFICIT_LAW.

No change to the law, target function, arm construction, thresholds, or write policy after observing results. Scientific negatives are valid.
