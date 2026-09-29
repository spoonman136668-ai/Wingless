UP-234B PREREGISTRATION — FAMILY-CENTERED POOLED PARITY

Parent UP-233B qualified scientific result.
Observed parent anchor: pooled parity accuracy=0.75; family-conditioned parity accuracy=1.0.

Question: is the family-conditioned parity signal a shared native parity direction hidden by family-specific offsets, or does each family require genuinely different parity geometry?

Freeze exact UP-233B substrate:
- training phases 58..72;
- evaluation phases 121..152;
- families mixed4, observe4, store4, cross3;
- native features exactly native_correct_count, mean_absolute_margin, near_zero_margin_count, min_absolute_margin;
- no evaluation-label fitting;
- no adaptive feature selection;
- no nonlinear classifier;
- no live activation.

Preregistered transform:
1. Compute one four-coordinate mean for each family using TRAINING phases only.
2. Subtract that frozen family training mean from each training and evaluation vector.
3. Standardize the centered coordinates using TRAINING centered rows only.
4. Fit exactly one pooled EVEN centroid and one pooled ODD centroid across all four families.
5. At inference, family identity may only select the frozen training-derived centering offset. It may not select a classifier or parity centroid.
6. Phase/parity is never an inference feature.

Report:
- parent pooled and family-conditioned anchors;
- family-centered pooled parity accuracy over all 128 evaluation states;
- per-family correct counts.

Interpretation:
- accuracy=1.0 supports a shared native parity direction after family-offset removal;
- improvement above 0.75 but below 1.0 supports partially shared geometry;
- no improvement supports genuinely family-relative parity geometry beyond a simple offset.

Scientific negatives are valid. No post-result tuning.
