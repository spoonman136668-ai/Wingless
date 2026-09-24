# Wingless UP-45 — reversible convex exact fusion result

Status: Windows-qualified scientific result. No activation, promotion, deployment, or external integration authorized.

Scientific parent: UP-44 seal `584373c3bd32fc305cffffeec16ad5e7ee07878c`.

## Question

Can the rolling task optimizer plus a reversible complete-graph convex clustering proximal operator produce exact spectral equality without sticky snapping, target groups, or a partition menu, and will the task objective select that exact-fusion state?

## Frozen mechanics

- square-root probability soft-memory representation;
- no global soft-memory renormalization;
- rolling seven-observation full-rank directional memory;
- 48 updates;
- checkpoints at 0/12/24/48;
- perturbation 0.002;
- learning rate 0.002;
- maximum coordinate update 0.004;
- resource price 0.02;
- soft-capacity tau 0.0025;
- proximal lambda derived as `0.002 * 0.02 = 0.00004`;
- no sticky state;
- no known target geometry;
- no partition menu;
- no hard merge search;
- no pair-affinity field;
- checkpoint selection by the unchanged training-side task objective only;
- exact fusion not used for checkpoint selection;
- true held-out evaluation only after selection.

## Authoritative Windows result

Workflow run: `35968534136`

Source head: `1c402ca58029941d1a5f9ebede2b9fd4fc1a5d0c`

Artifact: `10796101279`

Artifact digest: `sha256:ec079fc43e003544121184977071802e95af4e8f9906c41c8b91b0292bd394ee`

The production-priority guard, focused tests, deterministic double probe, and full repository regression passed.

### Task objective

- initial objective: `0.5267850187116467`;
- step 12 objective: `0.5256553271125561`;
- step 24 objective: `0.5581278352889602`;
- step 48 objective: `0.5281358881877997`;
- selected checkpoint: step `24`;
- selected objective gain: `0.031342816577313415`.

Selected step-24 smooth components:

- normalized phase score: `0.8723651983452807`;
- mean correct value probability: `0.5976892665979858`;
- mean correct relation probability: `0.3974934519527407`;
- harmonic task score: `0.5623049003135845`.

### Exact fusion behavior

The proximal operator did create exact equality:

- first exact fusion step: `29`;
- maximum exact capacity: `8`;
- exact fused group at step 29: channels `[0,3]`;
- step-29 exact offsets:
  `[-0.045664698831175514, -0.04501676812740249, 0.007218659931181362, -0.045664698831175514, 0.07894943527059196, 0.05017807058798021]`;
- the exact pair separated again on the next part of the trajectory;
- therefore the operator demonstrated exact fusion **and** reversibility.

The selected checkpoint itself was not exactly fused:

- selected exact capacity: `6`;
- selected exact fusion: `false`.

### Hard capability at selected step 24

- held-out accuracy: `0.907958984375`;
- mutable commit accuracy: `0.3385416666666667`;
- exact final-table accuracy: `0.3541666666666667`;
- relation accuracy: `0.5833333333333334`;
- held-out retention delta: `0.092041015625`;
- commit retention delta: `0.6614583333333333`;
- capability gates: `false`.

## Scientific classification

**Qualified Branch-C result: exact fusion is reachable and reversible, but the frozen task objective selected a non-fused checkpoint.**

This is not evidence that exact equality improves hard recurrent capability. It shows that a principled non-sticky operator can create exact symmetry, while the current smooth task objective still prefers another point on the trajectory.

The frozen post-UP-45 decision tree therefore points to an objective-decomposition audit, not stronger fusion pressure or lambda tuning.

## Causal-lineage caveat discovered after UP-45 completion

During extraction of the step-29 state, an earlier controlled-comparison defect was identified in the UP-44 lineage: UP-44's square-root free-running signal uses deterministic noise prefix `44000000`, whereas UP-43 used `39000000`.

UP-45 itself remains a valid experiment under its frozen UP-44-derived configuration, and its internal conclusions above are unchanged. However, the stronger earlier causal attribution that the UP-43 -> UP-44 change was due only to the square-root representation is confounded by the simultaneous deterministic-noise schedule change.

Before using that attribution as a foundation for the next mechanistic conclusion, a seed-matched UP-44 control must be run with the UP-43 `39000000` schedule while changing only the representation.

## Plain speak

The new fusion mechanism worked mechanically: it made two channels exactly equal and later let them separate again, so it is not sticky.

But the task score preferred step 24, before exact fusion happened. That puts us in the predeclared branch where we inspect why the objective prefers the near-symmetry state.

There is also an older A/B cleanliness issue: the square-root experiment changed the deterministic noise schedule at the same time. We caught it while auditing this result. The current observations still stand, but we need one clean seed-matched control before claiming the representation alone caused the earlier improvement.
