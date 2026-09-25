# Wingless UP-65B — 512-write depth-stratified margin calibration

Status: preregistered scientific calibration/scale experiment. Prepared while UP-64B runs; not result-informed by UP-64B.

Scientific parent: sealed UP-63B Windows evidence `e0b9d3ed351bdd20ee336655c0daf5dce0e9ea55`.

## Question

UP-63B found monotonic decoder-margin calibration in all 16 tested fixed-depth conditions at 256 writes. UP-65B asks whether that descriptive calibration survives a doubled 512-write closed-loop horizon at each propagation depth.

## Frozen design

Unchanged substrate and diagnostic semantics:

- frozen pair [0,5];
- existing observer bank and classifier/regressor training;
- training depth 0;
- held-out depths 32, 128, 512, 1024, evaluated separately;
- baseline closed-loop update only;
- no confidence-carry intervention;
- no oracle state;
- frozen margin bins [0.00,0.25), [0.25,0.50), [0.50,0.75), [0.75,1.00].

Schedules: 98M, 99M.
Memory noise: 0.085, 0.09.
Writes per scenario: exactly 512.

For every schedule/noise/depth condition, record event count and empirical commit correctness in every frozen margin bin.

No retry, retraining, adaptive binning, threshold selection, policy change, result-informed stopping, promotion, production authority, or activation is permitted.

## Interpretation

Persistence, flattening, inversion, sparse bins, or failure of the margin/correctness relationship are all valid results. The bins remain descriptive and are not candidate policy thresholds.
