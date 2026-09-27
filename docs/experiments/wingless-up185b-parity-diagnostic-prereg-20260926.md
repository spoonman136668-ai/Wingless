# Wingless UP-185B — parity diagnostic

Status: preregistered scientific calibration diagnostic.

Scientific parent: sealed UP-184B d0f944790268afd4e7a6d15380710e123e818587.

## Question

UP-184B reduced calibration error through multi-phase pooling, but phases 31 and 33 repeated while phase 32 differed. Does the remaining drift track a simple two-state alternating context class?

## Frozen diagnostic models

Use the same UP-184B feature table and Laplace estimator.

Construct before evaluation:
- pooled_26_30: phases 26,27,28,29,30 together;
- parity_diagnostic:
  - even calibration: phases 26,28,30;
  - odd calibration: phases 27,29.

For evaluation only, phase parity chooses the corresponding already-frozen table.

## Held-out phases

31, 32, 33.

## Measurements

For pooled and parity-diagnostic models:
- predicted/actual event-mass ratio;
- absolute event-mass ratio error;
- Brier;
- ECE;
- AUROC;
- average precision;
- mean absolute mass-ratio error across held-out phases.

## Interpretation

Parity is an **external diagnostic tag**, not an acceptable native deployment feature. If parity conditioning removes most residual drift, the missing information is consistent with a repeating two-state context and the next experiment must search for an internal-state correlate. If not, residual calibration drift is not explained by this simple alternation.

## Bounds

No evaluation fitting, adaptive scaling, intervention, maintenance action, capacity change, extra model calls, live activation, or production authority.
