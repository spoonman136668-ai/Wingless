# Wingless UP-183B — calibration drift

Status: preregistered scientific shadow-detection experiment.

Scientific parent: sealed UP-182B 5d5978457f6b1a25fc493f66fffed448f3941d6c.

## Question

How stable is the frozen phase-26 native failure-probability model across several additional held-out phases?

## Frozen model

Reuse UP-182B exactly:
- calibration data: phase 26 only;
- features: update class × fixed absolute-margin bin × temporal streak state;
- fixed 0.01 margin bins through .05, [.05,.10), [.10,+inf);
- Laplace estimate p=(crossings+1)/(slots+2);
- phase-26 global smoothed rate for unseen cells.

No refitting occurs after phase 26.

## Held-out phases

Evaluate unchanged on:
- phase 28;
- phase 29;
- phase 30.

## Measurements per phase

- actual crossings;
- summed predicted crossings and predicted/actual ratio;
- Brier score;
- ten-bin ECE;
- AUROC;
- average precision.

## Interpretation

Stable risk probabilities require predicted event mass and calibration metrics to remain consistent across phases, not merely high AUROC. Drift is a valid negative and blocks cost-based intervention thresholds.

## Bounds

Shadow only. No evaluation fitting, adaptive bins, threshold tuning, maintenance action, capacity increase, extra model calls, live activation, or production authority.
