UP-239B PREREGISTRATION — HELD-OUT COMPOSITION PARITY TRANSFER

Parent UP-238B qualified EXTREME_TRANSFER_PERFECT: the frozen family-centered near_zero_margin_count parity observable trained on phases 58..72 transferred with 1.0 accuracy through both the far window 185..216 and extreme-far window 249..280 on the four parent compositions.

Question: does that parity observable transfer to composition families that were already defined as held-out profiles earlier in the B lineage but were not used by the UP-238B parity family set?

Freeze the parent training/evaluation mechanics:
- original training families remain exactly:
  mixed4={0,5,1,6}
  observe4={5,6,7,8}
  store4={0,1,2,3}
  cross3={0,5,13}
- training phases remain exactly 58..72;
- feature remains near_zero_margin_count only;
- pooled EVEN/ODD centroids and pooled standard deviation are learned only from the original four training families;
- no phase input;
- no adaptive feature selection;
- no nonlinear classifier;
- no evaluation-label fitting;
- no live activation.

Held-out composition families are frozen from the previously defined UP-193B evaluation profiles:
- store3={0,1,2}
- observe3={5,6,7}
- mixed3={0,5,1}
- report2={13,14}

Because the observable is family-centered, each held-out family receives exactly one unlabeled calibration statistic: its mean near_zero_margin_count over phases 58..72. No parity labels from the held-out family are used to estimate that mean. The parent pooled standard deviation and EVEN/ODD centroids remain frozen from the original four families.

Evaluate each held-out family on both frozen disjoint windows:
- far phases 185..216;
- extreme-far phases 249..280.

Also reproduce the UP-238B parent anchor exactly.

Report overall and per-family accuracy for both windows.

Classification:
HELDOUT_COMPOSITION_TRANSFER_PERFECT if every held-out family is 32/32 correct in both windows.
HELDOUT_COMPOSITION_TRANSFER_DECAY if the parent anchor reproduces but any held-out family is below 32/32 in either window.
ANCHOR_NOT_REPRODUCED if UP-238B no longer reproduces EXTREME_TRANSFER_PERFECT with accuracy 1.0.

No threshold tuning, family search, family replacement, refitting, or post-result retry is allowed. Scientific negatives are valid.
