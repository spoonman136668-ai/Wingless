UP-237B PREREGISTRATION — FAR FAMILY-CENTERED NEAR-ZERO PARITY TRANSFER

Parent UP-236B qualified CENTERING_REQUIRED: raw single-feature accuracy=0.875, family-centered near_zero_margin_count accuracy=1.0 on phases 121..152.

Question: does the same training-derived family-centered near_zero_margin_count parity observable transfer to a substantially farther phase window?

Freeze training exactly as UP-236B: phases 58..72, four frozen families, near_zero_margin_count only, family means and standardization from training only, one pooled EVEN and ODD centroid, no evaluation fitting, no phase input, no adaptive selection, no nonlinear classifier, no live activation.

Parent evaluation anchor remains phases 121..152 and must reproduce centered accuracy=1.0.
New far evaluation phases exactly 185..216 (32 consecutive phases), with no refitting or recentering.

Report overall and per-family far accuracy.
Classification:
FAR_TRANSFER_PERFECT if far accuracy=1.0.
FAR_TRANSFER_DECAY if far accuracy<1.0.
ANCHOR_NOT_REPRODUCED if the UP-236B parent anchor fails.
Scientific negatives are valid. No post-result tuning.
