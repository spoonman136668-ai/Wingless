# Wingless UP-44 FIXA — seed-matched sqrt-probability causal control

Status: preregistered causal replication.

## Why this exists

After UP-45 completed, audit of the UP-43 -> UP-44 comparison found that UP-44 used deterministic free-running noise seed prefix `44000000`, while UP-43 used `39000000`.

That means the original UP-44 result remains a valid observation, but the claim that the representation was the *only* changed scientific variable was not cleanly isolated.

This control repairs that causal comparison without changing any result already observed.

## Controlled comparison

Scientific comparison baseline: UP-43 seal `85f9ab7837bacda5deed635c5e7bf1e5ed5b4c05`.

Keep all UP-44 mechanics exactly as qualified, including:

- `0.5*sqrt(probability)` amplitude encoding;
- no global soft-memory renormalization;
- rolling full-rank directional memory;
- buffer 7;
- 48 updates;
- checkpoints 0/12/24/48;
- perturbation 0.002;
- learning rate 0.002;
- max coordinate update 0.004;
- same 96/32 training split;
- true held-out 128 tables;
- same harmonic phase/value/relation objective;
- resource price 0.02;
- soft-capacity tau 0.0025;
- no sticky projection;
- no target geometry;
- no partition menu;
- no hard merge search;
- no held-out checkpoint selection.

The only scientific correction relative to the original UP-44 implementation is:

- free-running deterministic noise prefix `44000000 -> 39000000`.

The value `39000000` is inherited exactly from UP-43; it is not chosen from any UP-44 or UP-45 result.

## Frozen interpretation

- If the seed-matched square-root control again produces positive selected objective gain and selects a provisional symmetry-basin checkpoint, the representation-mismatch diagnosis is restored with clean causal isolation.
- If the seed-matched control improves the objective but does not select a symmetry basin, the square-root representation helps but the stronger geometry-alignment claim is not replicated.
- If the seed-matched control fails to improve the objective, the original UP-44 causal attribution is rejected; treat the deterministic noise schedule as a material confound and reconsider the UP-44/45 mechanistic narrative.
- Capability gates remain unchanged and are reported separately.

No thresholds, optimizer settings, seeds, gates, or horizons may be changed after observing the result.

## Plain speak

We changed two things by accident in UP-44: the way uncertainty is represented and the exact deterministic noise sequence.

This rerun puts the old noise sequence back while keeping the new representation.

If the good behavior remains, we know the representation really caused it. If it disappears, we caught a confound before building more conclusions on top of it.
