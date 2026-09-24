# Wingless UP-44 FIXA — seed-matched sqrt-probability causal control

Status: **SEALED — qualified causal-control result; weaker representation benefit replicated, stronger geometry-alignment claim not replicated.**

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

## Authoritative Windows result

GitHub Actions run `35975826227`, job `107555940483`, executed on dedicated runner `WINGLESS-LINKDEADKB` at research head `ba1314f784a2102c340273265e51cfabcf271c66`.

Qualification evidence:

- production-priority guard: PASS (`WINGLESS_HOST_GUARD_PASS`), BelowNormal priority, `GOMAXPROCS=4`;
- guarded automated qualification: PASS;
- focused UP-44 qualification and full Go regression: PASS;
- deterministic qualification classification: `qualified-scientific-result`;
- training/validation split remained 96/32;
- true held-out remained 128 tables and was not used for selection;
- selected checkpoint: step 48;
- initial objective: `0.6047706218899666`;
- selected objective: `0.6054161920939569`;
- selected objective gain: `+0.0006455702039903644`;
- first provisional fusion-basin entry: step 11;
- maximum provisional capacity: 8;
- selected provisional capacity: 6;
- selected nearest offset gap: `0.0038524422318966453`, outside the frozen `0.0025` provisional fusion tolerance;
- untouched selected true-heldout accuracy: `0.901123046875`;
- mutable commit accuracy: `0.3411458333333333`;
- exact final-table accuracy: `0.375`;
- relation accuracy: `0.5416666666666666`;
- capability gates without exact fusion: false.

## Sealed classification

This is the second preregistered interpretation case: **objective-positive, selected-geometry-negative**.

The seed-matched square-root representation again produced a positive training-selected objective gain, so evidence that the square-root representation can improve the smooth objective survives the causal control. However, the selected checkpoint was step 48 with provisional capacity 6 and nearest gap above the frozen symmetry-basin tolerance. Therefore the stronger claim that the representation change by itself aligns task selection with the symmetry basin is **not replicated** and must be downgraded.

The observed UP-45 trajectory facts remain unchanged. This control changes only the causal narrative around UP-43 -> UP-44.

## Plain speak

The new uncertainty representation still helped the score even after restoring the old deterministic noise sequence, so that part is real. But it did **not** make the task choose the near-fused geometry this time. That means the representation helps, while the stronger claim that it alone steers training into the useful symmetry basin did not survive the cleaner test.
