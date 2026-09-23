# Wingless UP-37: adaptive task-weighted continuous symmetry fusion

Status: research branch only; not activated, promoted, or connected to ckb-plane.

Parent research result: UP-36 qualified scientific negative sealed at `fde6ad40fdd9548e24374c16a2459d42f9c076d3`.

## Question

UP-36 proved that continuous task-weighted attraction can create exact repeated-spectrum symmetry from six independent offsets without evaluating hard merge candidates.

It failed because the attraction field was estimated once at the singleton starting geometry and then frozen. The best checkpoint fused only one pair, stopping at commutant capacity 8 and missing the preregistered held-out and mutable gates.

UP-37 asks:

**Does continuous symmetry formation become sufficient when task-derived attraction is recomputed after every flow step from the current geometry?**

## Controlled change from UP-36

UP-37 preserves all of UP-36's frozen mechanics:

- six independently parameterized spectral offsets;
- no finished partition menu;
- no hard merge candidate evaluation;
- 50% continuous pair-probe fraction;
- simultaneous flow rate 0.45;
- sticky fusion tolerance 0.0025 radians;
- at most 10 flow steps;
- capacity price 0.02;
- same balanced 96/32 training-only inner split;
- same transport-discovered observer and phase-code learner;
- same untouched 128-table true held-out pool.

The only scientific change is:

**the pairwise attraction field is rebuilt from the current geometry before every flow step.**

## Group-preserving adaptive probes

Once two directions fuse, they remain an exact group.

Subsequent probes operate on current fused groups rather than individual members:

- choose two distinct current groups;
- move both group centroids 50% toward their midpoint;
- keep all members of each group exactly equal;
- re-center and restore RMS offset magnitude 0.05;
- evaluate only on training-side validation.

A probe never intentionally hard-merges the two groups.

Its positive score gain over the current geometry becomes the attraction weight for that step.

## Adaptive flow

At every step:

1. evaluate the current geometry;
2. continuously probe every remaining group pair;
3. rebuild the task-attraction matrix from those current-geometry probes;
4. move all group centroids simultaneously;
5. preserve all previous fusions exactly;
6. fuse any groups whose centroid gap crosses 0.0025;
7. evaluate the new checkpoint;
8. repeat from the new geometry.

Thus guidance changes as symmetry forms.

## Model selection

Every flow checkpoint is scored using:

`score = min(validation held, validation commit) - 0.02 * capacity / 36`.

The best checkpoint is frozen before true held-out evaluation.

No true held-out information participates in affinity estimation, checkpoint selection, or stopping.

## Scientific gates

Preflight:

- balanced inner split;
- no true held-out selection.

Adaptive requirement:

- at least two affinity-field recomputations;
- at least one positive task-derived affinity;
- selected flow step >= 2;
- selected commutant capacity > 8, demonstrating additional exact symmetry beyond the UP-36 capacity-8 endpoint;
- selected commutant capacity < 36.

Final capability:

- held-out accuracy >= 0.985;
- mutable commit accuracy >= 0.95;
- exact final table accuracy >= 0.90;
- relational query accuracy >= 0.95.

Retention versus full capacity:

- held-out loss <= 0.02;
- mutable-commit loss <= 0.05.

## Interpretation boundary

A positive result would isolate the UP-36 failure to stale task guidance: continuous symmetry formation can build sufficient internal stable structure when the task-attraction field adapts as the geometry reorganizes.

A negative result would indicate that merely refreshing the local attraction field is not enough. The next experiment should then separate the continuous flow rule from the affinity estimator—for example by comparing adaptive local flow with a direct differentiable optimization of the spectral parameters.

UP-37 still contains scaffolding:

- the six spectral directions are predefined;
- the flow equation is hand-designed;
- the sticky-fusion tolerance is hand-designed;
- affinities are estimated by finite continuous probes.

A positive successor should remove the hand-designed flow itself and learn continuous transport parameters directly under task/resource pressure.

UP-29 already established exact real-orthogonal equivalence; no result here supports uniquely complex or quantum-computation claims.

## Plain speak

UP-36 looked once at the starting geometry, decided which directions should attract each other, and then kept following that old map even after the landscape changed.

UP-37 redraws the map after every movement.

As internal symmetry forms, the system asks again: "given where I am now, which remaining directions should move together?"

If that lets the system grow enough symmetry to recover robust mutable memory, it means the missing ingredient in UP-36 was adaptive self-evaluation—not the idea of continuous self-organization itself.

## Authority boundary

UP-37 is mathematical research only. It does not invoke a language model, activate Wingless in production, alter ckb-plane/KTRADE authority, access brokers or credentials, modify accepted refs, or perform promotion/deployment.
