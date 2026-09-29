UP-LM9C PREREGISTRATION — PROFILE FAMILY INVARIANCE

Parent UP-LM8C qualified PERMUTATION_FAMILY_EXACT_DEFICIT_LAW on 108/108 frozen rows with zero write_error and zero post_target_error.

Question: does the exact deficit law
writes = max(0, pre_countdown - target)
and
post_countdown = target
remain exact when the profile family is changed while every other substrate dimension is frozen?

Frozen parent/baseline:
- accepted parent SHA: f7874c5262f34906261b67d28c9d1d3f19c62665
- rotations={8,21}
- permutations={identity,rotate1,random_fixed}
- random_fixed=[3,0,5,1,4,2]
- six existing arms from uplm2xArms
- existing prepressure write loop
- 32-write ceiling
- no correction-policy change
- no live activation

Existing target primitives remain exact:
- deferred_only target = uplm3hDeferredTarget(a)
- layout_only target = uplm3hLayoutTarget(a)
- hybrid_min target = min(uplm3hDeferredTarget(a), uplm3hLayoutTarget(a))

Change exactly one substrate dimension:
- replace parent profile family {deferred_only,layout_only,hybrid_min}
- with {deferred_only,layout_only,hybrid_max}

Freeze the new profile exactly:
- hybrid_max target = max(uplm3hDeferredTarget(a), uplm3hLayoutTarget(a))
- no learned, adaptive, fitted, or tunable parameter is introduced.

For every arm row compute only the existing deficit-law measurements:
- deficit=max(0, pre_countdown-target)
- write_error=writes-deficit
- post_target_error=post_countdown-target

Frozen evaluation size:
2 rotations x 3 profiles x 3 permutations x 6 arms = 108 rows.

Report:
- row count
- exact-law row count
- mismatch row count
- min/max write_error
- min/max post_target_error
- mismatch counts by profile, permutation, and rotation

Classification:
- PROFILE_FAMILY_EXACT_DEFICIT_LAW if all 108 rows have write_error=0 and post_target_error=0.
- PROFILE_DEPENDENT_DEFICIT_RESIDUAL if the parent anchor reproduces and any LM9C row has nonzero write_error or post_target_error.
- ANCHOR_NOT_REPRODUCED if UP-LM8C no longer reproduces PERMUTATION_FAMILY_EXACT_DEFICIT_LAW.

No change to the deficit law, existing target primitives, arm construction, thresholds, rotations, permutations, write policy, or measurements after observing results. Scientific negatives are valid. No post-result tuning.
