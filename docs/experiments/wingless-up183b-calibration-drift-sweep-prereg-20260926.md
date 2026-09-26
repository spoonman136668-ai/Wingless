# Wingless UP-183B — calibration drift sweep

Status: preregistered scientific shadow-evaluation experiment.

Scientific parent: sealed UP-182B 5d5978457f6b1a25fc493f66fffed448f3941d6c.

## Question

How stable is the frozen phase-26 continuous-margin-plus-temporal probability table across several additional held-out phases when no refitting is allowed?

## Frozen model

Reuse UP-182B exactly:
- class × absolute-margin bin × temporal-state cells;
- Laplace-smoothed cell probabilities from phase 26 only;
- global phase-26 fallback probability;
- no parameter or threshold fitting after phase 26.

## Held-out phases

Evaluate the same frozen model independently on phases:
- 27;
- 28;
- 29;
- 30.

For each phase measure:
- actual crossings;
- predicted crossing mass;
- predicted/actual ratio;
- Brier score;
- ECE with the same 10 equal-width probability bins;
- AUROC;
- average precision.

Also report mean and range of predicted/actual ratio across phases.

## Interpretation

Stable ratios near 1 with stable Brier/ECE support numeric risk calibration. Persistent phase-to-phase drift despite strong AUROC/AP means Wingless has a robust ordering signal but not a stable absolute probability model.

## Bounds

Shadow only. No intervention, no refitting, no adaptive bins, no threshold tuning, no added model calls, no extra training, no capacity change, no live activation.
