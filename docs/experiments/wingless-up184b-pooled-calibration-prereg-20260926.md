# Wingless UP-184B — pooled calibration

Status: preregistered scientific shadow-detection experiment.

Scientific parent: sealed UP-183B d4e715241c94f32b46e52f595bf6a84b7c675b1d.

## Question

Can a probability model calibrated on several already sealed prior phases stabilize numeric failure-risk estimates across new held-out phases better than the original phase-26-only calibration, while preserving discrimination?

## Frozen feature model

Reuse UP-182B/183B exactly:
- update class;
- fixed absolute-margin bin;
- temporal streak state;
- same 0.01 margin bins through .05, [.05,.10), [.10,+inf);
- Laplace cell estimate p=(crossings+1)/(slots+2);
- smoothed global fallback for unseen cells.

## Calibration arms

Before evaluation, construct two frozen models:

1. **single_phase_26** — phase 26 only, identical to UP-182B/183B.
2. **pooled_26_30** — aggregate phases 26, 27, 28, 29, and 30.

The pooled range is contiguous and fixed before looking at phases 31–33. No phase is omitted or weighted differently.

## Held-out phases

Evaluate both frozen models unchanged on:
- phase 31;
- phase 32;
- phase 33.

## Measurements

Per model × held-out phase:
- actual crossings;
- summed predicted crossings;
- predicted/actual ratio;
- absolute event-mass ratio error |ratio-1|;
- Brier score;
- ten-bin ECE;
- AUROC;
- average precision.

Aggregate each model's mean absolute event-mass ratio error across the three held-out phases.

## Interpretation

If pooled calibration reduces held-out event-mass error without materially degrading discrimination, cross-phase pooling is a better basis for numeric risk. If not, the remaining drift needs an additional native state variable rather than more averaging.

## Bounds

Shadow only. No evaluation fitting, adaptive bins, per-phase scaling, threshold tuning, intervention, maintenance action, capacity increase, extra model calls, live activation, or production authority.
