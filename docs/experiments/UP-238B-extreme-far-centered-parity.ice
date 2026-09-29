UP-238B PREREGISTRATION — EXTREME FAR FAMILY-CENTERED PARITY TRANSFER

Parent UP-237B qualified FAR_TRANSFER_PERFECT: the frozen family-centered near_zero_margin_count parity observable trained on phases 58..72 achieved 1.0 accuracy on far phases 185..216.

Question: does the same frozen observable remain perfect on a still farther disjoint phase window without refitting?

Freeze training exactly as UP-237B:
- phases 58..72;
- four frozen families;
- near_zero_margin_count only;
- family means and standardization from training only;
- one pooled EVEN and ODD centroid;
- no evaluation fitting or recentering;
- no phase input;
- no adaptive selection;
- no nonlinear classifier;
- no live activation.

Parent anchor phases 185..216 must reproduce FAR_TRANSFER_PERFECT with accuracy 1.0.
New extreme-far evaluation phases exactly 249..280 (32 consecutive phases), with no refitting or recentering.

Report overall and per-family extreme-far accuracy.

Classification:
EXTREME_TRANSFER_PERFECT if extreme-far accuracy=1.0.
EXTREME_TRANSFER_DECAY if extreme-far accuracy<1.0.
ANCHOR_NOT_REPRODUCED if the UP-237B parent anchor fails.

Scientific negatives are valid. No post-result tuning.
