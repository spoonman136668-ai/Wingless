UP-LM6C PREREGISTRATION — PREPRESSURE DEFICIT LAW

Parent UP-LM6B qualified TARGET_DEPENDENT_EXTENSION across 108 frozen rows.

Question: is the variable number of prepressure writes completely explained by the prepressure countdown deficit, writes = max(0, pre_countdown - target), or is an additional profile/layout/rotation dependence required?

Freeze the exact 108 LM6B rows and substrate:
profiles={deferred_only,layout_only,hybrid_min}
rotations={5,13}
permutations={identity,reverse,rotate2}
6 arms each.
Do not alter the prepressure transform, target function, arm construction, selector policy, or resource schedule.
No live activation and no policy change.

For every row compute:
deficit = max(0, pre_countdown - target)
write_error = writes - deficit
post_target_error = post_countdown - target

Report:
- exact-law row count;
- mismatch row count;
- minimum and maximum write_error;
- minimum and maximum post_target_error;
- mismatch counts by profile.

Classification:
EXACT_DEFICIT_LAW if every row has writes=deficit and post_countdown=target.
WRITE_DEFICIT_EXACT_POST_RESIDUAL if writes=deficit for every row but any post_countdown differs from target.
TARGET_DEPENDENT_RESIDUAL if any row violates writes=deficit.
ANCHOR_NOT_REPRODUCED if LM6B no longer reproduces TARGET_DEPENDENT_EXTENSION.
OTHER_VALID_PATTERN otherwise.

Observational mechanism test only. Scientific negatives are valid. No post-result tuning.
