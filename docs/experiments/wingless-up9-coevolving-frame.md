# Wingless UP-9: co-evolving internal frame

Status: Windows-qualified internal-frame result; not activated, promoted, or connected to ckb-plane.

Parent research qualification: UP-8 Windows observer ablation sealed at `e357e8bf25b69e838aca9652e8454fe2e83cd97b`.

## Question

Can a compact internal reference frame remove the unseen-depth ambiguity demonstrated by UP-8 without supplying transport depth or explicitly inverting the latent state?

## Motivation from UP-8

After correcting UP-7's split and query-conditioning confounds, UP-8 established:

- fixed-frame magnitude-only held-out accuracy: 0.62109375;
- fixed-frame coherence-aware held-out accuracy: 1.0;
- unseen-depth coherence-aware held-out accuracy: 0.373046875.

That isolates the next problem: the memory information survives, but its observable representation rotates with the transport frame.

## Internal frame

UP-9 adds five reference/pilot states:

- one global anchor;
- one phase-coded pilot for each of the four entities.

The global anchor is uniform over all sixteen coordinates.

Each entity pilot occupies only that entity's four value coordinates and assigns phases:

`0, pi/2, pi, 3pi/2`.

The pilot bank contains 80 complex scalars total. It is intentionally much smaller than a full transported 16-vector basis and therefore does not encode the exact inverse.

## Co-evolution

For the co-evolving condition, memory and pilots pass through the same propagation operator for the same number of blocks.

For each queried entity the observer receives only two real features derived from:

`z = <phase_ref|memory> * conj(<anchor_ref|memory>)`

normalized to unit magnitude.

Under unitary co-evolution, both inner products are depth-invariant. The product also cancels an arbitrary global phase applied to the memory alone.

No depth integer, inverse state, prototype table, or full transported basis is provided.

## Global-phase nuisance

Every memory example receives a deterministic memory-only global phase before transport. The pilot bank does not receive that phase.

This forces the two-correlation construction to act as a gauge-invariant reference rather than relying on a fixed absolute complex phase convention.

## Ablations

UP-9 measures:

1. unitary static frame, fixed depth 128;
2. unitary co-evolving frame, fixed depth 128;
3. unitary static frame, unseen depths;
4. unitary co-evolving frame, unseen depths;
5. matched-control static frame, unseen depths;
6. matched-control co-evolving frame, unseen depths.

The balanced 64/64 table selection from UP-8 is reused.

Training depths:

`8, 24, 72, 216, 432, 648`

Held-out depths:

`32, 128, 512, 1024`

Memory perturbation amplitude is 0.05. Pilot states are not perturbed in this first frame experiment; pilot corruption is a later stress test if the concept works.

## Observer

Each entity has a separate four-class linear-softmax head.

Each head receives exactly two real features and trains for 200 deterministic full-batch steps at learning rate 1.0.

## Scientific gate

The internal-frame hypothesis is considered supported when:

- unitary co-evolving unseen-depth held-out accuracy >= 0.95;
- and co-evolution improves unitary unseen-depth accuracy over static references by at least 0.15.

Matched-control co-evolving accuracy is measured but is not required to fail.

## Interpretation boundary

A positive result would not mean that Wingless has learned a clock or discovered a reference frame autonomously. It would establish that a small co-evolving reference can make transported unitary memory depth-invariant and directly readable without exposing the inverse or an explicit depth value.

A positive UP-9 permits UP-10: make the frame trainable/compressible, reduce pilot overhead, and stress it with pilot noise/drift and mutable read/write programs.

A negative result would indicate that the proposed five-pilot frame is insufficient and should be redesigned before increasing observer complexity.

## Authority boundary

UP-9 is mathematical research only. It does not register an inference backend, invoke a language model, activate a Wingless worker/listener, execute model-selected tools, or alter ckb-plane queue, retry, workspace, acceptance, promotion, broker, credential, or production authority.


## Windows qualification

Authoritative operator proof on 2026-09-22 against source head `66f9678373831188ab00d26b9ea01ffc4d405600` passed the harness and complete Wingless regression.

Observed held-out accuracies:

- unitary static frame, fixed depth: 0.37109375;
- unitary co-evolving frame, fixed depth: 1.0;
- unitary static frame, unseen depths: 0.2607421875;
- unitary co-evolving frame, unseen depths: 1.0;
- matched-control static frame, unseen depths: 0.2783203125;
- matched-control co-evolving frame, unseen depths: 0.8076171875.

The unitary co-evolving frame improved unseen-depth held-out accuracy by 0.7392578125 over static references while using only five pilot states and two real observer features per entity. No depth value, explicit inverse, or runtime prototype lookup was supplied.

Interpretation: within the tested substrate, a compact co-evolving internal reference frame resolves the frame/depth ambiguity isolated by UP-8.
