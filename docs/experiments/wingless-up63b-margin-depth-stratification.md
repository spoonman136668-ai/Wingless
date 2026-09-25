# Wingless UP-63B — decoder-margin calibration by propagation depth

Status: preregistered scientific calibration diagnostic.

Scientific parent: sealed UP-62B Windows evidence `9fe4b71e267e34af346fe771da0025c2542d5555`.

## Question

UP-62B preserved a strong descriptive relationship between decoder margin and commit correctness through 128/256-write closed-loop horizons: the highest frozen margin bin beat the lowest in all 12 conditions, with 11/12 conditions monotonic across all four bins. UP-63B asks whether that signal is stable within each propagation depth rather than being an artifact of aggregating depth regimes.

## Frozen design

Unchanged substrate and diagnostic semantics:

- frozen pair [0,5];
- existing observer bank and classifier/regressor training;
- training depth 0;
- no confidence-carry intervention;
- no oracle state;
- frozen margin bins [0.00,0.25), [0.25,0.50), [0.50,0.75), [0.75,1.00].

Held-out propagation depths are stratified separately:

- 32;
- 128;
- 512;
- 1024.

Untouched deterministic schedule bases:

- 94M;
- 95M.

Memory-noise levels:

- 0.085;
- 0.09.

Writes per scenario:

- 256.

For every schedule/noise/depth condition, all 48 deterministic scenarios use only that propagation depth for their write gaps, and event count plus empirical commit correctness is recorded in each frozen margin bin.

No retry, retraining, adaptive binning, threshold selection, policy change, result-informed stopping, production authority, or activation is permitted.

## Interpretation

Stable calibration, depth-specific flattening/inversion, sparse bins, or failure are all valid results. The bins remain descriptive; this experiment does not select a control threshold.
