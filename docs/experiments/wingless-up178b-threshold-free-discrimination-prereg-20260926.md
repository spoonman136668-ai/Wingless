# Wingless UP-178B — threshold-free discrimination

Status: preregistered scientific shadow-detection experiment.

Scientific parent: sealed UP-177B c49343cffe447e23f4932845d7ee792d2e28b0fd.

## Question

Does the frozen rank-plus-temporal native risk score intrinsically discriminate imminent correctness crossings better than frozen rank-only risk on a fresh held-out phase, independent of choosing a warning threshold?

## Frozen predictors

Reuse without refitting:
- rank_plus_temporal probabilities from UP-175B/176B;
- rank_only probabilities from UP-176B.

No score or threshold is changed after seeing phase-23 outcomes.

## Held-out evaluation

Use terminal phase 23 with the same six subjects, two 12-step cleanup paths, and 120 old examples per state.

For every example-update slot, record:
- frozen predictor score before the update;
- whether correctness changes after the update.

## Metrics

Per predictor:
- AUROC using grouped score thresholds;
- AUPRC using non-interpolated average precision over score groups;
- maximum-score-group positive rate;
- base crossing rate;
- maximum-score-group enrichment over base rate;
- Brier score.

## Interpretation

If rank_plus_temporal exceeds rank_only on AUROC/AUPRC and enrichment, temporal history provides threshold-independent discrimination. If not, its earlier calibration gain does not translate into meaningful ordering.

## Bounds

Shadow only. No intervention, threshold selection, refit, extra training, update suppression, capacity change, result-informed retry, live activation, or production authority.
