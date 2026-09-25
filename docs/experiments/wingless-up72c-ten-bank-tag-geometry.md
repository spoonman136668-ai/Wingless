# Wingless UP-72C — ten-bank tag-geometry control

Status: preregistered scientific representation-control experiment. Prepared while UP-71C runs; not result-informed by UP-71C.

Scientific parent: sealed UP-70C Windows evidence `52af265b43b4074e62bdc3836ab68c7d90e58452`.

## Question

UP-70C showed that 10 banks in the fixed 16-dimensional state passed at zero noise but failed at noise 0.004 and 0.008 under the existing golden-rotation tag family. UP-72C asks whether that noisy 10-bank failure is specific to tag geometry or persists under a deterministic equally spaced phase control.

## Frozen design

Unchanged:

- state dimension 16;
- exactly 10 banks;
- four entities per bank;
- 256 deterministic scenarios;
- coherence decoder and magnitude control;
- reversible transport depth 64;
- value-accuracy gate 0.99;
- exact-scenario gate 0.95.

Tag families:

- `golden_rotation`: existing `up57cGoldenTags(10)`;
- `uniform_phase`: tag i = 2*pi*i/10 for bank i.

Schedules: 115M, 116M.
Memory noise: 0, 0.001, 0.002, 0.004.

No tag optimization, phase search, state-dimension increase, decoder change, threshold change, adaptive stopping, schedule search, family selection, promotion, or activation is permitted.

## Interpretation

Each tag-family/schedule/noise point is sealed exactly as observed. The experiment is a control, not a selector: superiority of either family does not authorize replacing the existing tags. Similar failure across families supports a broader capacity/noise limitation; divergent behavior supports tag-geometry sensitivity.
