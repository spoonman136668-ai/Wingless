# Wingless UP-50B — pair [0,5] robustness/generalization

Status: preregistered scientific generalization experiment.

Scientific parent: UP-49B seal `3b69129d94c97eccc90af320f99dcac7c2450053`.

## Question

UP-49B selected exact pair [0,5] using training-side smooth objective only, and that same pair passed the full frozen hard gate.

UP-50B freezes that geometry and asks whether the gate survives distributional stress rather than selecting another structure.

## Frozen cases

No selection or optimization is run.

The exact step-42 pair-[0,5] fused offsets are evaluated under four combinations:

- memory noise 0.03, held depths 32/128/512/1024;
- memory noise 0.08, held depths 32/128/512/1024;
- memory noise 0.03, held depths 64/256/1024/2048;
- memory noise 0.08, held depths 64/256/1024/2048.

Each case is compared with the unchanged full-capacity control and uses the exact existing hard gate.

## Interpretation

If all four pass, the learned exact structure generalizes beyond the distribution on which it was selected and can move to persistence/restart testing.

If only high-noise or shifted-depth cases fail, that condition becomes the next bounded robustness target. No fusion pressure or thresholds are changed.

## Plain speak

The score found a structure that passed the real test.

Now we stop choosing anything and simply make the world noisier and the propagation depths different to see whether that structure was actually robust.
