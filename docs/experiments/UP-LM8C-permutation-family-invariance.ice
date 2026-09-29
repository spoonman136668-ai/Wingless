UP-LM8C PREREGISTRATION — PERMUTATION FAMILY INVARIANCE

Parent UP-LM7C qualified HELDOUT_ROTATION_EXACT_DEFICIT_LAW on all 108 frozen rows at held-out rotations {8,21}, with zero write_error and zero post_target_error.

Question: does the exact deficit law remain invariant when the arm permutation family changes from the previously tested identity/reverse/rotate2 family to a frozen family containing a one-step rotation and a non-cyclic fixed permutation?

Hypothesis:
- If the deficit law is substrate-invariant with respect to arm ordering, every LM8C row will retain:
  writes = max(0, pre_countdown - target)
  and post_countdown = target.
- Otherwise, at least one row will show nonzero write_error or post_target_error.

Frozen parent/baseline:
- exact parent qualification head: 8537c0a562d87d48cfcde366623c77a1c6d6941a
- parent classification required: HELDOUT_ROTATION_EXACT_DEFICIT_LAW

Single changed substrate dimension:
- parent permutations {identity,reverse,rotate2}
- LM8C permutations {identity,rotate1,random_fixed}
- rotate1 mapping over six arms: output[i] = input[(i+1)%6]
- random_fixed mapping: output = [input[3],input[0],input[5],input[1],input[4],input[2]]
- no runtime randomness and no seed-dependent generation.

Frozen controls:
- rotations={8,21}
- profiles={deferred_only,layout_only,hybrid_min}
- exactly six existing arms from uplm2xArms
- existing uplm3tTarget target rule
- existing prepressure write loop semantics
- 32-write ceiling
- correction policy unchanged
- no adaptive feature selection
- no post-result fitting
- no live activation
- no policy change

Evaluation substrate:
2 rotations x 3 profiles x 3 permutations x 6 arms = 108 rows.

For every row compute:
deficit=max(0, pre_countdown-target)
write_error=writes-deficit
post_target_error=post_countdown-target

Report:
- row count
- exact-law row count
- mismatch row count
- min/max write_error
- min/max post_target_error
- mismatch counts by permutation
- mismatch counts by rotation
- mismatch counts by profile

Classification:
PERMUTATION_FAMILY_EXACT_DEFICIT_LAW if the parent anchor reproduces and every one of 108 LM8C rows has write_error=0 and post_target_error=0.
PERMUTATION_DEPENDENT_DEFICIT_RESIDUAL if the parent anchor reproduces and any LM8C row has nonzero write_error or post_target_error.
ANCHOR_NOT_REPRODUCED if UP-LM7C no longer reproduces HELDOUT_ROTATION_EXACT_DEFICIT_LAW.

Scientific negatives are valid. No threshold, law, target function, permutation set, row set, write policy, or capacity may change after execution begins.
