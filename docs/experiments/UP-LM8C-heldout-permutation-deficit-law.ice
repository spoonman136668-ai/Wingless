UP-LM8C PREREGISTRATION — HELD-OUT PERMUTATION DEFICIT-LAW TRANSFER

Parent UP-LM7C qualified HELDOUT_ROTATION_EXACT_DEFICIT_LAW on all 108 frozen rows at held-out rotations {8,21}. The exact law remained:
writes = max(0, pre_countdown - target)
and post_countdown = target,
with zero mismatches.

Planning provenance:
- advisory model: nvidia/nemotron-3-ultra-550b-a55b:free
- advisory qualification identity: 086c5c3cef88ef60b2c852be724aa88fe21c47d80194abbfeacbff15d128196a
- advisory response SHA-256: 3683059d8ccf9be094a9db9eaabef1cc82dc256786afeab0ec10ae41c25a3aca
- deterministic validation before freeze corrected two advisory defects: ambiguous "e.g." permutation was frozen exactly, and row count was corrected from 18 combinations to 108 arm-level rows.

Question: is the exact deficit law invariant to held-out permutation structure when every other substrate dimension and the correction policy are frozen?

Freeze:
- profiles={deferred_only,layout_only,hybrid_min};
- rotations={8,21};
- six arms produced by the existing uplm2xArms construction;
- existing uplm3tTarget target rule;
- existing prepressure write loop and 32-write ceiling;
- no policy change;
- no live activation.

Change exactly one substrate dimension:
- replace parent permutations {identity,reverse,rotate2} with {identity,rotate1,random_fixed};
- rotate1 is the deterministic six-arm cyclic ordering out[i]=in[(i+1)%6];
- random_fixed is exactly [3,0,5,1,4,2], interpreted as out[i]=in[random_fixed[i]];
- no random sampling occurs at runtime.

For each of 3 profiles × 2 rotations × 3 permutations × 6 arms = 108 rows compute:
deficit=max(0, pre_countdown-target)
write_error=writes-deficit
post_target_error=post_countdown-target

The parent UP-LM7C anchor must reproduce before evaluating the held-out permutation family.

Report:
- row count;
- exact-law row count;
- mismatch row count;
- min/max write_error;
- min/max post_target_error;
- mismatch counts by permutation and profile.

Classification:
HELDOUT_PERMUTATION_EXACT_DEFICIT_LAW if every row has write_error=0 and post_target_error=0.
PERMUTATION_DEPENDENT_DEFICIT_RESIDUAL if the parent anchor reproduces and any row has nonzero write_error or post_target_error.
ANCHOR_NOT_REPRODUCED if UP-LM7C no longer reproduces HELDOUT_ROTATION_EXACT_DEFICIT_LAW.

No change to the law, target function, arm construction, thresholds, rotations, profiles, or write policy after observing results. Scientific negatives are valid.
