# Wingless UP-61C — six-bank boundary refinement

Status: preregistered scientific robustness/boundary experiment.

Scientific parent: UP-60C seal `a157a3580438a28627ba885c13435f23bb53bd8a`.

## Question

UP-60C bounded the golden-rotation six-bank failure transition between noise 0.01 and 0.015, while fixed-irregular already failed at 0.01 after having passed 0.0075 in UP-59C.

UP-61C asks where those two frozen families cross the existing qualification gate on untouched schedules.

## Frozen design

Untouched deterministic schedule bases:

- 77M;
- 78M.

Golden-rotation noise ladder:

- 0.010;
- 0.011;
- 0.012;
- 0.013;
- 0.014;
- 0.015.

Fixed-irregular noise ladder:

- 0.0075;
- 0.0080;
- 0.0085;
- 0.0090;
- 0.0095;
- 0.0100.

All other conditions are unchanged: six banks, state dimension 16, 256 scenarios, depth 64, coherence decoder, reversible transport, and the existing 0.99 value / 0.95 exact-scenario gate.

No geometry modification, decoder change, threshold change, result-informed stopping, family promotion, production authority, or activation is permitted.

## Interpretation

This is a fixed boundary-refinement experiment only. Any failure is a valid scientific result. The family-specific transition intervals will be sealed exactly as observed and will not themselves authorize promotion.
