# Wingless UP-180B — continuous-margin replication

Status: preregistered scientific shadow-detection experiment.

Scientific parent: sealed UP-179B 3add97cb27d37b2818a56312ca3b1807d5b3f59a.

## Question

Does continuous absolute margin again outperform the frozen rank-plus-temporal score on an independent held-out phase?

## Frozen scores

Evaluate unchanged on terminal phase 25:
- rank_plus_temporal = frozen probability from UP-175B/176B;
- margin_only = -abs(current decision margin);
- rank_plus_temporal_margin_tiebreak = frozen probability + 1e-6/(1+abs(margin)).

No coefficients, bins, or thresholds are fitted from phase 25.

## Measurements

For each score:
- AUROC;
- average precision;
- base crossing rate;
- actual crossings;
- total slots.

## Interpretation

Replication of the phase-24 ordering supports continuous margin as the primary B-lane ranking signal. Failure to replicate means the phase-24 advantage was context-specific.

## Bounds

Shadow only. No threshold tuning, intervention, model fitting, maintenance trigger, extra training, capacity change, result-informed retry, live activation, or production authority.
