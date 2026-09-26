# Wingless UP-175B — held-out hazard calibration

Status: preregistered scientific shadow-detection experiment.

Scientific parent: sealed UP-174B 5fdc7debe728fe385baa1f4c9e28ab628220da6f.

## Question

Can a completely frozen class-aware risk estimate derived from prior sealed phases predict next-update correctness crossings in a held-out phase, without retuning?

## Frozen predictor

Use the exact class × static-rank-stratum × temporal-state crossing densities measured in UP-174B as probabilities. Use the sealed UP-173B ALL out-of-band density 0.00043793793793793793 as the background probability for all states outside those cells.

Static strata:
- rank_1_4;
- rank_5_10.

Temporal states:
- current_only = streak 1;
- transition_to_persistent = streak 2;
- established_persistent = streak 3+.

Class-specific probabilities are copied verbatim from the sealed UP-174B result. No fitting occurs in this experiment.

Freeze a warning threshold at predicted probability >= 0.20 before phase-20 execution.

## Held-out evaluation

Use terminal phase 20 with the same six subjects, two 12-step cleanup paths, and 120 old examples per state. For every example before every update:
1. compute static rank;
2. update consecutive accepted-band streak;
3. assign frozen predicted probability;
4. observe whether the immediately following update changes correctness.

## Measurements

- total example-update slots;
- actual crossings;
- sum of predicted probabilities;
- Brier score;
- warning TP, FP, FN, TN;
- warning precision and recall;
- one-update lead time;
- calibration by frozen predictor cell.

## Bounds

Shadow only. No intervention, maintenance trigger, update suppression, adaptive threshold, model fitting, extra training, capacity change, result-informed retry, live activation, or production authority.
