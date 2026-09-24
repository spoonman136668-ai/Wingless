# Wingless UP-52C — memory load interference

Status: preregistered scientific scale experiment.

Scientific parent: UP-51C seal `c4fb28b3852e916d419d77cc85a014b9f95946b9`.

## Question

UP-51C located the dimension-128 floating-point depth boundary while prior evidence showed task classification remained perfect even at one million propagation steps.

The more consequential scaling question is whether a persistent mutable memory remains reliable as the number of sequential write/commit operations grows.

UP-52C therefore asks how repeated memory updates affect the already-established hard memory gate.

## Frozen design

- full-capacity unitary memory geometry;
- memory noise fixed at 0.05;
- held propagation depths fixed at 32, 128, 512, 1024;
- the same training tables, held-out tables, observable discovery, decoder training, relation head, and hard thresholds used by the existing breadth/multiplicity memory harness;
- write-count ladder: 16, 32, 64, 128, 256 writes per scenario;
- 48 deterministic scenarios at every write count;
- no threshold change, precision change, renormalization, reunitarization, or result-informed parameter adjustment.

The decoder/observable model is trained once and reused for all five write loads.

## Frozen gate

A load point passes only if:

- static held-out accuracy >= 0.99;
- commit decode accuracy >= 0.99;
- exact final-table accuracy >= 0.95;
- relational query accuracy >= 0.95;
- maximum norm drift <= 1e-10.

## Interpretation

The first failing write count is a memory-continuity/load boundary. If all tested counts pass, the next C-lane experiment should increase state/load diversity rather than simply lengthening the same write chain. If a boundary appears, the next experiment should separate decoder error accumulation from state-transport error without changing the frozen thresholds.
