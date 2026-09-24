# Wingless UP-53C — fixed-state phase-multiplex memory capacity

Status: preregistered scientific scale experiment.

Scientific parent: UP-52C seal `9a1df4dd9a9556905e7030db633d3064068663a8`.

## Question

UP-52C showed that a single mutable table survives at least 256 sequential write/commit cycles with perfect hard metrics. The next scale question is not longer time; it is how many independent memories can coexist in the same fixed-size state.

UP-53C phase-codes multiple ordered memory banks into the same 16-dimensional four-entity/four-value state and asks where recovery degrades under fixed noise.

## Frozen design

- state dimension fixed at 16 for every condition;
- 4 entities and 4 values per entity;
- simultaneous bank counts: 1, 2, 3, 4, 5, 6;
- each bank has a frozen deterministic phase tag `0.41*b + 0.173*b^2`;
- 256 deterministic scenarios per bank count;
- additive memory noise amplitude fixed at 0.05;
- global phase nuisance applied after normalization;
- 64 steps of the same local reversible unitary transport;
- ordered bank values are decoded using global-phase-invariant coherence/fidelity against the complete per-entity prototype set;
- a magnitude-only prototype decoder is reported as a matched information-ablation control.

## Frozen gate

A bank count passes when coherence decoding achieves:

- value accuracy >= 0.99;
- exact all-bank/all-entity scenario accuracy >= 0.95;
- maximum norm drift <= 1e-9.

Magnitude-only performance does not control acceptance.

Scientific negatives are valid. No phase tag, noise level, threshold, transport depth, or decoder rule may be changed after execution.
