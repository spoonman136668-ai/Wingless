UP-236B PREREGISTRATION — NEAR-ZERO MARGIN CENTERING NECESSITY

Parent UP-235B qualified scientific result.
Observed parent anchor: minimum perfect feature cardinality=1 and the only minimum perfect mask is 4, corresponding exactly to near_zero_margin_count.

Question: does near_zero_margin_count alone carry a globally aligned parity signal, or is the family-centering correction required to expose it?

Freeze exact UP-235B substrate:
- training phases 58..72;
- evaluation phases 121..152;
- families mixed4, observe4, store4, cross3;
- selected feature exactly near_zero_margin_count only;
- training parity labels only;
- no phase/parity inference feature;
- no evaluation-label fitting;
- no adaptive feature selection;
- nearest-centroid classifier only;
- no live activation.

Evaluate exactly two pooled single-feature decoders:
RAW:
- standardize near_zero_margin_count from all TRAINING rows without family centering;
- fit one EVEN and one ODD pooled centroid.

FAMILY_CENTERED:
- subtract the TRAINING mean for the corresponding family;
- standardize centered TRAINING rows;
- fit one EVEN and one ODD pooled centroid.

At evaluation, family identity may only select the frozen family mean in FAMILY_CENTERED. It may not select a classifier or centroid.

Report both overall accuracies and per-family correct counts.

Classification:
RAW_SINGLE_FEATURE_SUFFICIENT if RAW accuracy=1.0 and FAMILY_CENTERED reproduces 1.0.
CENTERING_REQUIRED if RAW accuracy<1.0 and FAMILY_CENTERED reproduces 1.0.
ANCHOR_NOT_REPRODUCED if FAMILY_CENTERED accuracy<1.0.
OTHER_VALID_PATTERN otherwise.

Scientific negatives are valid. No post-result tuning.
