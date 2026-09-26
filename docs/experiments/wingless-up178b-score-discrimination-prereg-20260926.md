# Wingless UP-178B — score discrimination

Status: preregistered scientific shadow-detection experiment.

Scientific parent: sealed UP-177B c49343cffe447e23f4932845d7ee792d2e28b0fd.

## Question

Does the frozen native risk score meaningfully rank imminent correctness crossings above non-crossings even though no single tested threshold provided a usable trigger?

## Frozen predictors

Reuse without refitting:
- rank_plus_temporal probability from UP-175B/176B;
- rank_only probability from UP-176B.

No thresholds are selected in this experiment.

## Held-out evaluation

Use terminal phase 23 with the same six subjects, two 12-step cleanup paths, and 120 old examples per state.

For every example before every update:
1. compute static rank and accepted-band streak;
2. assign both frozen scores;
3. apply the unchanged update;
4. record whether correctness crosses.

## Measurements

For each predictor:
- AUROC with exact tie handling;
- average precision / AUPRC-style area over score groups;
- base crossing rate;
- predicted score sum;
- actual crossing count.

## Interpretation

AUROC/AUPRC above their chance baselines supports real ranking information even if threshold selectivity remains inadequate. A near-chance result would reject the current score as a useful failure-ranking signal.

## Bounds

Shadow only. No threshold tuning, intervention, maintenance trigger, update suppression, model fitting, extra training, capacity change, result-informed retry, live activation, or production authority.
