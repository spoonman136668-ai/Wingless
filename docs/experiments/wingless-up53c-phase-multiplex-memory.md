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


## Authoritative Windows result

Workflow run: `36031972244`

Runner: `WINGLESS-UP-C`

Source head: `fcfb8c3a76aede0d83c06f9d5a6daf6be5ed8dca`

Artifact: `10823715628`

Artifact digest: `sha256:904f222d1411edfdddf1e49637b3460cacc0c8f9d11551ddab2ad33d8c3091f5`

The guarded qualification, focused tests, deterministic double probe, and full repository regression passed.

Coherence-aware decoding:

- 1 bank: value `1.0`, exact scenario `1.0`;
- 2 banks: value `1.0`, exact scenario `1.0`;
- 3 banks: value `1.0`, exact scenario `1.0`;
- 4 banks: value `0.962646484375`, exact scenario `0.6875`;
- 5 banks: value `0.8068359375`, exact scenario `0`;
- 6 banks: value `0.7478841145833334`, exact scenario `0`.

The minimum coherence margin collapsed from `0.023861231972672314` at three banks to `0.00008703271519971967` at four banks.

Magnitude-only value accuracy was already only `0.7177734375` at two banks and `0.4501953125` at three banks.

## Scientific classification

At fixed dimension 16 and memory noise 0.05, phase/coherence carries substantially more simultaneous ordered-memory information than magnitude alone. The frozen phase code supports three banks perfectly but crosses a sharp robustness boundary at four banks.

The next C-lane experiment freezes the four-bank phase geometry and maps its noise ladder down to zero. This distinguishes a finite-noise separation problem from a representational collision/information-limit problem.
