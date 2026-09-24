# Wingless UP-50B — pair [0,5] robustness/generalization

Status: Windows-qualified scientific result.

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


## Authoritative Windows result

Workflow run: `35988173601`

Runner: `WINGLESS-UP-B`

Artifact: `10803036693`

Artifact digest: `sha256:52e00829b5694d247cc0548732532e47414633d2fe006bff612849dd7d39d5f5`

Results:

- noise 0.03 / standard depths: gate PASS; held `0.998291015625`, commit `0.9947916666666666`, final `1.0`, relation `1.0`
- noise 0.08 / standard depths: gate FAIL; held `0.98583984375`, commit `0.9361979166666666`, final `0.8541666666666666`, relation `1.0`
- noise 0.03 / shifted depths: gate PASS; held `0.998291015625`, commit `0.9947916666666666`, final `1.0`, relation `1.0`
- noise 0.08 / shifted depths: gate FAIL; held `0.98583984375`, commit `0.93359375`, final `0.8333333333333334`, relation `1.0`

## Scientific classification

The frozen pair-[0,5] structure generalizes across the tested depth shift at low noise, so depth distribution is not the limiting variable in this range.

At noise 0.08 the mutable-memory gate fails under both depth schedules, with the main loss in commit and exact final-table accuracy while relational accuracy remains perfect.

Because UP-49B already passed at the original noise 0.05, the next B-lane experiment is a frozen noise ladder between 0.05 and 0.08 using the standard depth schedule. No structure or threshold changes are permitted.

## Plain speak

The structure is robust to changing how deep the memory travels.

What breaks it is enough noise. At 0.03 it is excellent; at 0.08 its relation logic still works, but repeated memory writes start getting corrupted.

The next run finds where that noise limit actually sits.
