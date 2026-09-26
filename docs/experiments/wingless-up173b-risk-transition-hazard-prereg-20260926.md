# Wingless UP-173B — risk-transition hazard

Status: preregistered scientific shadow-detection experiment.

Scientific parent: sealed UP-172B c41a28b9fe3c42f6e43aeaaeb0a6d334c96b672a.

## Question

UP-172B found that exactly-two-step persistence was the highest-risk streak bin, while 3+ step persistence was less risky. Does that temporal hazard shape reproduce in a fresh phase-18 context?

## Frozen design

Reuse the accepted class bands unchanged:
- STORE rank <= 10;
- OBSERVE rank <= 10;
- REPORT rank <= 4.

Use the same six subjects, two cleanup paths, and 12 chained updates as UP-171B/172B, but construct the frozen post-REPORT state at terminal phase 18.

Before every chained update, classify each old example by its consecutive in-band streak:
- out_of_band = 0;
- current_only = 1;
- transition_to_persistent = exactly 2;
- established_persistent = 3+.

No warning changes an update.

## Measurements

Per scope ALL/STORE/OBSERVE/REPORT:
- slots and crossings per streak state;
- crossing density per state;
- transition-to-persistent vs established-persistent density ratio.

## Interpretation

Reproduction of a higher crossing density at exactly streak 2 than at streak 3+ supports a transition-hazard signal rather than a monotonic persistence signal. Failure is a valid negative.

## Bounds

Shadow-only. No maintenance trigger, update suppression, adaptive band, schedule change, extra training, capacity change, result-informed retry, live activation, or production authority.
