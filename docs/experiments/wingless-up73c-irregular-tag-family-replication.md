# Wingless UP-73C — irregular tag-family replication

Status: preregistered scientific geometry closure experiment.

Scientific parent: sealed UP-72C Windows evidence `cdcd3c70b8b7c334df3fc8ee05168b63ababf052`.

## Question

UP-72C showed that the golden-rotation tag family remains viable at 10 banks through noise 0.002 while a uniform phase layout collapses through phase-code cancellation. UP-73C asks whether success is unique to the golden angle or generalizes to other fixed irrational rotations.

## Frozen design

Unchanged:

- state dimension 16;
- exactly 10 banks;
- four entities per bank;
- 256 deterministic scenarios;
- coherence decoder and magnitude control;
- reversible transport depth 64;
- value gate 0.99;
- exact-scenario gate 0.95.

Tag families:

- golden rotation: existing golden-angle tags;
- sqrt2 rotation: tag i = i * sqrt(2) radians;
- sqrt3 rotation: tag i = i * sqrt(3) radians.

Schedules:

- 117M;
- 118M.

Memory noise:

- 0.002;
- 0.003.

No tag optimization, phase search, state-dimension increase, decoder change, threshold change, adaptive stopping, family selection, promotion, or activation is permitted.

## Interpretation

This is a replication/control experiment, not a selector. If multiple irrational rotations preserve the same boundary, the mechanism is broader than one special tag family. Divergence identifies tag-geometry sensitivity and is sealed without tuning.
