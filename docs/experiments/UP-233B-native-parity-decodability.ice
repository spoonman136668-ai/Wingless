UP-233B PREREGISTRATION — NATIVE PARITY DECODABILITY
Parent UP-232B qualified scientific result.
Observed parent anchor: baseline family accuracy=0.875, true-parity-conditioned accuracy=1.0, flipped-parity accuracy=0.625; mixed4 true-parity=32/32 and flipped-parity=0/32.

Question: is the parity variable that rescued far-phase family classification already encoded in the four native observables, or does the decoder require an external phase clock?

Freeze exact UP-232B substrate:
- training phases 58..72;
- evaluation phases 121..152;
- families mixed4, observe4, store4, cross3;
- native features exactly native_correct_count, mean_absolute_margin, near_zero_margin_count, min_absolute_margin;
- training-set standardization only;
- no evaluation-label fitting;
- no adaptive feature selection;
- no nonlinear classifier;
- no live activation.

Preregistered decoders:
1. POOLED parity decoder: standardized nearest centroid for EVEN vs ODD using all training families together.
2. FAMILY-CONDITIONED parity diagnostic: standardized nearest centroid for EVEN vs ODD within each frozen family.

Phase parity may be used only as a TRAINING LABEL and for evaluation scoring. Phase/parity must not be supplied as an inference feature.

Report:
- pooled parity accuracy across all 128 evaluation states;
- family-conditioned parity accuracy across all 128 evaluation states;
- per-family pooled and family-conditioned correct counts.

Interpretation:
- perfect pooled transfer supports a globally native parity observable;
- family-conditioned-only transfer supports a family-relative parity observable;
- weak transfer supports continued dependence on an external phase coordinate.

Scientific negatives are valid. No post-result tuning.