# Wingless UP-189B — native mass-correction replication

Status: preregistered scientific shadow-calibration experiment.

Scientific parent: sealed UP-188B ec79b82c146b0c75634f0abec190434d1f58cebd.

## Question

Does the UP-188B native-context multiplicative correction remain stable across a longer untouched phase range?

## Frozen calibration

Reuse UP-188B unchanged:
- pooled probability model trained on phases 26..30;
- context_A if native pre-cleanup correct count <= 621, context_B otherwise;
- context_A factor = 1.0103904235372925;
- context_B factor = 0.8983526543771945.

No factors are recomputed.

## Untouched evaluation

Evaluate phases:
- 37, 38, 39, 40, 41, 42.

Compare pooled versus corrected probabilities exactly as UP-188B.

## Measurements

Per phase:
- native context;
- actual crossing count;
- predicted crossing mass;
- mass ratio and absolute error;
- Brier;
- ECE;
- AUROC;
- average precision.

Aggregate mean and maximum absolute mass-ratio error across all six phases.

## Bounds

Shadow only. No held-out fitting, phase/parity input, factor update, intervention, maintenance action, capacity change, extra model calls, live activation, or production authority.
