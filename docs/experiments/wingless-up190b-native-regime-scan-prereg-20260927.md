# Wingless UP-190B — native regime scan

Status: preregistered scientific shadow-calibration experiment.

Scientific parent: sealed UP-189B 62429a84c67c974773efe5ddb51d7324be880507.

## Question

Do later untouched phases expose any native calibration state outside the repeating 620/623 correct-count regime seen through UP-189B, and if so does the frozen UP-188B mass correction remain calibrated there?

## Frozen model

Reuse unchanged:
- pooled probability model trained on phases 26..30;
- native context threshold 621;
- context_A factor = 1.0103904235372925;
- context_B factor = 0.8983526543771945.

No factor, threshold, feature bin, or probability table is recomputed.

## Untouched scan

Evaluate phases 43..58 inclusive.

For each phase record:
- native correct count;
- native context label;
- pooled and corrected aggregate mass error;
- Brier score;
- ECE;
- AUROC;
- average precision.

Also report the unique native correct counts encountered.

## Interpretation

If a new native correct-count state appears, its unchanged corrected calibration is the first genuine transfer test beyond the known two-state regime.

If only 620 and 623 recur, the result is still informative: the present phase generator itself does not expose an out-of-distribution native state, so later calibration claims must not be described as OOD transfer.

## Bounds

Shadow only. No held-out fitting, phase/parity predictor input, context-threshold change, factor update, maintenance action, capacity change, extra model calls, live activation, or production authority.
