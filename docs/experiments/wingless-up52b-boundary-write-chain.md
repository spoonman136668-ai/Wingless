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


## Authoritative Windows result

Workflow run: `36031119720`

Runner: `WINGLESS-UP-B`

Source head: `1a14f0ec4fc0a768a76b6e17529dc604ec2754cc`

Artifact: `10821941255`

Artifact digest: `sha256:b5fa16163e33afc3b148c0f422364c34472aeff03d12be9d37e0e50e835ce69f`

The guarded qualification, focused tests, deterministic double probe, and full repository regression passed. The required 16-write reproduction succeeded: noise 0.065 PASS and noise 0.07 FAIL.

At noise 0.065:

- writes 1/2/4/8/16: PASS;
- writes 32/64: FAIL;
- commit accuracy declines from `0.9791666666666666` at one write to `0.923828125` at 64 writes;
- relation accuracy remains at or above the unchanged gate.

At noise 0.07:

- writes 1/2/4/8: PASS;
- writes 16/32/64: FAIL;
- commit accuracy declines from `0.9791666666666666` at one write to `0.90625` at 64 writes;
- relation accuracy remains at or above the unchanged gate.

The full-capacity control remains at commit/final/relation accuracy `1.0` throughout the reported long-chain points.

## Scientific classification

The pair-[0,5] robustness boundary is **cumulative commit-chain error**, not immediate decode failure. Higher noise moves the chain-length failure earlier: from between 16 and 32 writes at noise 0.065 to between 8 and 16 writes at noise 0.07.

The next B-lane experiment should separate closed-loop error propagation from per-step decoder noise by comparing the unchanged recurrent commit path against a teacher-forced/oracle-reset path at the same frozen geometries, noise levels, and write counts.
