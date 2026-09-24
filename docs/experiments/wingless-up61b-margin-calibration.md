# Wingless UP-61B — decoder-margin calibration diagnostic

Status: preregistered scientific self-evaluation/calibration experiment.

Scientific parent: UP-58B seal `850e010027508a9c463656b9f4ce5e73542048dd`.

This experiment is defined independently of UP-59B and UP-60B outcomes.

## Question

UP-58B showed that all three frozen confidence-carry thresholds improved commit accuracy over the closed-loop baseline across its untouched conditions.

UP-61B asks whether the decoder margin used by that policy is actually informative about commit correctness across a broader noise range, without changing policy behavior or selecting a threshold.

## Frozen design

The decoder/training substrate remains unchanged:

- frozen pair [0,5];
- existing observer bank and classifier/regressor training;
- training depth 0;
- held-out depths 32, 128, 512, 1024;
- 48 deterministic scenarios;
- 64 writes per scenario;
- baseline closed-loop state update only;
- no confidence-carry intervention.

Untouched deterministic schedule bases:

- 86M;
- 87M;
- 88M.

Memory-noise levels:

- 0.06;
- 0.07;
- 0.08;
- 0.09.

The commit margin is binned into fixed descriptive ranges before execution:

- [0.00, 0.25);
- [0.25, 0.50);
- [0.50, 0.75);
- [0.75, 1.00].

For every schedule/noise condition, record event count and empirical commit correctness in each bin.

No oracle state, carry policy, retry, retraining, adaptive binning, threshold selection, result-informed stopping, production authority, or activation is permitted.

## Interpretation

This is a calibration diagnostic only. Monotonic, non-monotonic, sparse, or inverted relationships between margin and correctness are all valid scientific results. The bins are not candidate policy thresholds and no policy is promoted by this experiment.
