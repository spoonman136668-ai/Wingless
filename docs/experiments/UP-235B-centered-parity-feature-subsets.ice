UP-235B PREREGISTRATION — FAMILY-CENTERED PARITY FEATURE SUBSETS

Parent UP-234B qualified scientific result.
Observed parent anchor: family-centered pooled parity accuracy=1.0 across all 128 frozen evaluation states.

Question: what is the minimum subset of the four frozen native observables that preserves perfect family-centered pooled parity transfer?

Freeze exact UP-234B substrate:
- training phases 58..72;
- evaluation phases 121..152;
- families mixed4, observe4, store4, cross3;
- native features exactly native_correct_count, mean_absolute_margin, near_zero_margin_count, min_absolute_margin;
- family centering estimated from TRAINING rows only;
- standardization estimated from TRAINING centered rows only;
- one pooled EVEN centroid and one pooled ODD centroid;
- family identity may select only the frozen training-derived centering offset;
- phase/parity is never an inference feature;
- no evaluation-label fitting;
- no adaptive feature selection;
- no nonlinear classifier;
- no live activation.

Evaluate every non-empty feature subset exactly once: 2^4-1=15 subsets.
For each subset, refit only the frozen-form training standardization and pooled parity centroids on those selected coordinates, then evaluate the same 128 held-out states.

Report:
- parent four-feature anchor;
- all 15 subset masks, feature names, cardinality, accuracy, and per-family correct counts;
- minimum cardinality among subsets with accuracy=1.0;
- every perfect subset at that minimum cardinality.

Interpretation:
- cardinality 1 means a single native observable carries the centered parity axis;
- cardinality 2 or 3 means parity is genuinely distributed across a low-order native combination;
- cardinality 4 means all four frozen observables are jointly required.

This is exhaustive ablation, not adaptive feature search.
Scientific negatives are valid. No post-result tuning.
