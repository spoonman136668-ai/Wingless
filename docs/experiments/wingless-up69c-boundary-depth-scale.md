# Wingless UP-69C — six-bank boundary depth scale

Status: preregistered scientific transport-scale experiment.

Scientific parent: sealed UP-68C Windows evidence `f336e130724d0a61dc279b51fdaad55c6ae3eb8a`.

## Question

UP-68C showed identical schedule/noise outcomes across depths 256, 512, and 1024. UP-69C asks whether the same boundary-local behavior remains stable at materially longer reversible transport without spending runner time on already-easy interior noise points.

## Frozen design

Unchanged substrate:

- six banks;
- state dimension 16;
- 256 deterministic scenarios;
- the same coherence decoder;
- the same reversible transport;
- prototype and scenario transport at the same tested depth;
- value-accuracy gate 0.99;
- exact-scenario gate 0.95.

Untouched deterministic schedule bases:

- 109M;
- 110M.

Preregistered boundary-local noise subset:

- golden rotation: 0.0108 and 0.0110;
- fixed irregular: 0.0087.

Preregistered depth ladder:

- 1024;
- 2048;
- 4096.

The 1024 point is retained as an internal bridge to the UP-68C depth range. No geometry, decoder, gate, family definition, adaptive stopping, family selection, promotion, or activation change is permitted.

## Interpretation

Depth invariance, numerical degradation, schedule sensitivity, or complete failure are all valid results. No threshold or noise level may be changed after execution.
