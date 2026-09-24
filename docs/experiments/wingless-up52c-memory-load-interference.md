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


## Authoritative Windows result

Workflow run: `36030618415`

Runner: `WINGLESS-UP-C`

Source head: `b859614e83e7662cb0a1cc2fa88c450aff7bb440`

Artifact: `10821505974`

Artifact digest: `sha256:6d77207ae7c7fc67771ce8c5c6ed7796eb689560c138362cdeb936ef8f80b667`

The guarded qualification, focused tests, deterministic double probe, and full repository regression passed.

Results:

- 16 writes: PASS; commit/final/relation all `1.0`;
- 32 writes: PASS; commit/final/relation all `1.0`;
- 64 writes: PASS; commit/final/relation all `1.0`;
- 128 writes: PASS; commit/final/relation all `1.0`;
- 256 writes: PASS; commit/final/relation all `1.0`.

Maximum norm drift remained approximately `1.3e-12` to `1.44e-12` throughout.

## Scientific classification

Repeated mutable write-chain length is not the limiting scale variable in this frozen full-capacity memory configuration through 256 sequential writes.

The next C-lane experiment must increase **simultaneous memory/state diversity or concurrency**, rather than merely extending the same single-table write chain. No existing threshold is changed.
