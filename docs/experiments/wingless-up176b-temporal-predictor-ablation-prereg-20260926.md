# Wingless UP-176B — temporal predictor ablation

Status: preregistered scientific held-out predictor ablation.

Scientific parent: sealed UP-175B 6fc3ee7e5d2535654db6c75b2806b82e60bdfe90.

## Question
Does temporal streak state improve held-out failure-risk prediction beyond static vulnerability rank alone?

## Frozen predictors
Evaluate both on terminal phase 21.

1. composite: exact class × rank-stratum × temporal-state probabilities frozen in UP-175B.
2. rank_only: class × rank-stratum probabilities formed before this run by pooling the sealed UP-174B cells across temporal state:
- STORE rank1_4 = 80/240; rank5_10 = 24/360
- OBSERVE rank1_4 = 67/240; rank5_10 = 28/360
- REPORT rank1_4 = 16/96
- all other states use the same sealed background probability 0.00043793793793793793.

No probabilities are fit on phase 21.

Freeze warning threshold 0.20 for both predictors.

## Measurements
For each predictor:
- Brier score;
- predicted crossing sum;
- TP/FP/FN/TN;
- precision and recall.

Primary comparison: lower held-out Brier score.
Secondary comparison: precision/recall at the shared frozen threshold.

## Bounds
Shadow only. No intervention, adaptive threshold, refitting, extra training, capacity change, result-informed retry, live activation, or production authority.
