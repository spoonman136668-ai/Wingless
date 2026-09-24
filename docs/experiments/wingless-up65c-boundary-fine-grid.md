# Wingless UP-65C — six-bank boundary fine grid

Status: preregistered scientific robustness/boundary experiment.

Scientific parent: UP-63C seal `66e705ac348294947625093728d4c6b792874dbb`.

This experiment is defined independently of UP-64C outcomes.

## Question

UP-63C showed schedule-stable fixed-irregular failure at 0.009, while the golden transition near 0.011 exhibited schedule sensitivity.

UP-65C asks for a finer untouched-schedule map of both frozen boundary intervals.

## Frozen design

Untouched deterministic schedule bases:

- 93M;
- 94M;
- 95M;
- 96M;
- 97M;
- 98M.

Golden-rotation noise grid:

- 0.0100;
- 0.0102;
- 0.0104;
- 0.0106;
- 0.0108;
- 0.0110.

Fixed-irregular noise grid:

- 0.0085;
- 0.0086;
- 0.0087;
- 0.0088;
- 0.0089;
- 0.0090.

All other conditions are unchanged: six banks, state dimension 16, 256 scenarios, depth 64, coherence decoder, reversible transport, and the existing 0.99 value / 0.95 exact-scenario gate.

No geometry modification, decoder change, threshold change, family selection, adaptive stopping, production authority, or activation is permitted.

## Interpretation

This is a fixed boundary-density experiment. The output is a schedule-by-noise qualification map, not a tuned threshold or promoted mechanism. Every passing and failing point is retained exactly as observed.
