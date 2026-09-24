# Wingless UP-52B — boundary write-chain decomposition

Status: preregistered scientific robustness experiment.

Scientific parent: UP-51B seal `3452be7012cb385dfc32941f89a5d1d8c35b4517`.

## Question

UP-51B found the frozen pair-[0,5] hard-gate boundary between noise 0.065 (PASS) and 0.07 (FAIL), while relation accuracy remained perfect and mutable commit/final-table accuracy degraded.

UP-52B asks whether that boundary is caused by immediate decode corruption or by cumulative error propagation through repeated write/commit cycles.

## Frozen design

Two memory-noise levels are tested:

- 0.065, the last passing UP-51B point;
- 0.07, the first failing UP-51B point.

For each noise level, the exact pair-[0,5] geometry and the unchanged full-capacity control are independently trained/evaluated exactly as in UP-51B.

Mutable integration is then evaluated at frozen write-chain lengths:

`1, 2, 4, 8, 16, 32, 64`.

Held depths remain `32, 128, 512, 1024`. The UP-49B hard gate is reused unchanged.

The 16-write result must reproduce the already-sealed UP-51B gate classification: 0.065 PASS and 0.07 FAIL. Failure to reproduce is an infrastructure/implementation inconsistency, not a scientific negative.

## Interpretation

If 0.07 fails even at very short chains, the boundary is primarily immediate decode/noise sensitivity.

If short chains pass and longer chains fail, cumulative commit-chain error is the dominant mechanism.

If the full-capacity control degrades similarly, the boundary is broader memory-decoder robustness rather than a pair-[0,5] structural cost.

No geometry, thresholds, optimizer settings, or noise values may be changed after observing results.
