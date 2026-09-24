# Wingless UP-51B — frozen pair-[0,5] noise boundary

Status: preregistered scientific robustness experiment.

Scientific parent: UP-50B recorded result `d74bd3a457b59bbd21532dd72c40f46e36786023`.

## Question

UP-50B showed that the frozen pair-[0,5] structure passes the unchanged hard gate at memory noise 0.03 and fails at 0.08, while prior UP-49B evidence passed at the original 0.05 condition.

UP-51B asks where the hard-capability boundary lies between 0.05 and 0.08 without changing structure, thresholds, depths, or selection.

## Frozen design

- exact pair-[0,5] fused geometry from UP-49B / UP-50B;
- standard held depths 32, 128, 512, 1024;
- noise ladder 0.05, 0.055, 0.06, 0.065, 0.07, 0.075, 0.08;
- unchanged full-capacity control;
- unchanged UP-49B hard gate;
- no optimization or structure selection.

The first failing noise and last passing noise are reported directly.

## Interpretation

This is boundary mapping, not threshold tuning. A scientific negative is valid and must be sealed exactly as observed.

If a monotone crossing is observed, the next B-lane experiment should test the mechanism responsible for loss around that boundary rather than adjusting the gate. If behavior is non-monotone, the next experiment should reproduce the affected points with additional frozen deterministic noise schedules before changing architecture.
