# Wingless UP-11: pilot-rank and memory-noise sweep

Status: Windows-qualified pilot-rank result; not activated, promoted, or connected to ckb-plane.

Parent research record: UP-10 Windows-qualified negative compression result sealed at `d125a6e69c95f0a887322afcf96a0870d8e599e9`.

## Why UP-11 exists

UP-10 reduced the successful UP-9 frame from five pilots to three and failed:

- full five-pilot held-out accuracy: 1.0;
- three-pilot compressed held-out accuracy: 0.85546875;
- reference noise 0.01: 0.853515625;
- reference transport mismatch 1e-3: 0.8515625.

Reference corruption and modest synchronization drift barely changed the already-degraded compressed baseline.

However, the UP-10 "clean" cell still used the standing memory perturbation amplitude 0.05. UP-10 therefore did not isolate:

1. loss from insufficient pilot rank;
2. loss from memory-noise sensitivity;
3. loss caused specifically by UP-10's one-scalar multiplexed readout.

UP-11 separates those effects.

## Reference construction

The four successful entity-phase pilots from UP-9 are an orthonormal rank-4 reference family because they occupy disjoint entity coordinate blocks.

For code rank `r = 1, 2, 3, 4`, UP-11 constructs the first `r` equal-energy four-point Fourier mixtures of those pilots.

A common anchor is retained.

Runtime pilot states therefore equal:

`1 + r`

and runtime frame storage is:

`16 * (1 + r)` complex scalars.

The complete complex pilot/anchor ratios are exposed to each entity-conditioned linear head, giving `2r` real features. Unlike UP-10, no component is discarded before the learned readout.

## Depth and table split

The balanced combination split from UP-8/UP-9 is retained.

A deterministic 32-table training selection and 32-table held-out selection are used.

Training depths:

`8, 24, 72, 216, 432, 648`

Held-out depths:

`32, 128, 512, 1024`

No explicit depth, inverse, or runtime prototype table is supplied.

Every memory also receives a deterministic global-phase nuisance.

## Noise sweep

Each unitary rank is trained and evaluated twice:

- memory noise = 0;
- memory noise = 0.05.

This produces eight unitary cells.

A matched non-unitary rank-4 control is also measured at memory noise 0.05.

## Diagnosis

The experiment reports:

- minimum pilot code rank reaching >=0.95 held-out accuracy without memory noise;
- minimum rank reaching >=0.95 at memory noise 0.05;
- rank-2, rank-3, and rank-4 accuracies with and without noise.

**Intrinsic rank loss supported** means the minimum noiseless passing rank is at least 3.

**Memory-noise sensitivity supported** means adding 0.05 noise drops rank-2 or rank-3 accuracy by at least 0.10.

**Four-pilot candidate supported** means code rank 3 — four total pilot states including the anchor — reaches >=0.95 under the noisy unseen-depth test.

## Interpretation boundary

UP-11 does not attempt mutable integration. Its purpose is to identify the smallest reference rank that actually retains enough observable information under the standing noise budget.

If rank 3 passes noisy evaluation, a four-total-pilot frame becomes the next integration candidate.

If only rank 4 passes, the UP-9 five-pilot frame is not redundant under the current linear/coherence interface and attempts to fold it should preserve four independent code directions.

If rank 2 passes noiselessly but fails with noise, the UP-10 failure is primarily robustness rather than information rank.

## Authority boundary

UP-11 is mathematical research only. It does not register an inference backend, invoke a language model, activate a Wingless worker/listener, execute model-selected tools, or alter ckb-plane queue, retry, workspace, acceptance, promotion, broker, credential, or production authority.


## Windows qualification

Authoritative operator proof on 2026-09-23 against source head `f74dce32c330969b599690a190419c8f16fa3fe6` passed the harness and complete Wingless regression.

Observed held-out accuracies:

- rank 1, no memory noise: 0.4296875;
- rank 1, memory noise 0.05: 0.421875;
- rank 2, no memory noise: 0.53125;
- rank 2, memory noise 0.05: 0.525390625;
- rank 3, no memory noise: 0.7734375;
- rank 3, memory noise 0.05: 0.70703125;
- rank 4, no memory noise: 0.96875;
- rank 4, memory noise 0.05: 0.95703125;
- matched non-unitary rank 4, memory noise 0.05: 0.732421875.

The minimum passing code rank is 4 both noiselessly and with memory noise 0.05. This supports intrinsic rank loss as the dominant explanation for UP-10's three-pilot failure under the current coherence/readout interface. It does not prove that five physical side vectors are globally minimal for every possible representation.
