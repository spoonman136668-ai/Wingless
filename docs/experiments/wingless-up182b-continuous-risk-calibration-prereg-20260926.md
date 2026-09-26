# Wingless UP-182B — continuous native-risk calibration

Status: preregistered scientific shadow-detection experiment.

Scientific parent: sealed UP-181B 9a6d2dccea33afd67bd93a10d3199b307f3203e7.

## Question

Can a fixed probability model built only from prior native state—absolute margin, temporal streak state, and update class—produce a calibrated and discriminative failure-risk estimate on a fresh held-out phase?

## Frozen train/evaluation split

- calibration phase: 26 only;
- held-out evaluation phase: 27 only;
- same six subjects;
- same two 12-step cleanup paths;
- same 120 old examples per state.

No phase-27 result is used to change bins, probabilities, thresholds, or features.

## Frozen features

For every example immediately before an update:
- update class: STORE / OBSERVE / REPORT;
- absolute margin bin:
  - [0,.01)
  - [.01,.02)
  - [.02,.03)
  - [.03,.04)
  - [.04,.05)
  - [.05,.10)
  - [.10,+inf)
- temporal state:
  - out_of_band;
  - current_only;
  - transition_to_persistent;
  - established_persistent.

The accepted class-specific vulnerability bands remain unchanged and define streak state only.

## Frozen probability estimator

For each phase-26 feature cell:
p = (crossings + 1) / (slots + 2)

This Laplace smoothing rule is fixed before evaluation. If a phase-27 cell was absent in phase 26, use the phase-26 global smoothed base rate.

## Measurements on phase 27

- Brier score;
- Expected Calibration Error using ten fixed probability bins;
- AUROC;
- average precision;
- actual crossings;
- summed predicted crossings;
- margin-only AUROC/AP comparator.

## Interpretation

Good discrimination plus reasonable held-out calibration supports a native scalar failure-risk estimate. Poor calibration or poor generalization is a valid negative and blocks intervention use.

## Bounds

Shadow only. No maintenance action, threshold tuning, adaptive bins, phase-27 fitting, extra model calls, training, capacity increase, live activation, or production authority.
