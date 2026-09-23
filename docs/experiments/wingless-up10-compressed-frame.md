# Wingless UP-10: compressed and noisy internal frame

Status: research branch only; not activated, promoted, or connected to ckb-plane.

Parent research qualification: UP-9 Windows internal-frame qualification sealed at `8c4b2b9f108169fed4b01897c14dadb492541c78`.

## Question

Can the successful UP-9 co-evolving internal frame be made materially cheaper and still support unseen-depth readout, reference corruption, and mutable read/write programs?

## Compression

UP-9 used five 16-dimensional pilot states:

- one anchor;
- four entity-specific phase pilots.

That is 80 runtime complex scalars.

UP-10 multiplexes two entities into each phase pilot:

- one anchor;
- one phase-code pilot for entities 0/1;
- one phase-code pilot for entities 2/3.

The runtime frame therefore falls to three states / 48 complex scalars, a 40% reduction.

Each phase-code pilot uses four ordered levels:

`-3, -1, +1, +3`.

For one entity the levels occupy the real component; for the paired entity they occupy the imaginary component.

The observer divides the selected pilot correlation by the anchor correlation, cancelling memory-only global phase, and reads exactly one real scalar for the queried entity.

## Data and depths

A balanced deterministic 32/32 table selection is drawn from the UP-8/UP-9 combination split.

Training depths:

`8, 24, 72, 216, 432, 648`

Held-out depths:

`32, 128, 512, 1024`

Memory perturbation amplitude remains 0.05.

## Scientific cells

The experiment measures:

1. the original five-pilot unitary frame as an in-run baseline;
2. three-pilot unitary frame, clean;
3. three-pilot unitary frame with independent reference noise amplitude 0.01;
4. three-pilot unitary frame with independent reference noise amplitude 0.03;
5. three-pilot unitary frame with reference-transport scale mismatch 1.0001;
6. three-pilot unitary frame with reference-transport scale mismatch 1.001;
7. three-pilot matched non-unitary frame, clean.

The observer is trained only on clean references. Reference-noise and drift cells are out-of-distribution stress tests.

## Mutable integration

A separate unitary integration uses:

- 32 mutable-memory scenarios;
- 12 writes/overwrites per scenario;
- only held-out transport depths;
- memory noise amplitude 0.05;
- independent reference noise amplitude 0.01;
- no inverse;
- no depth value;
- no prototype lookup;
- the compressed three-pilot frame.

At every commit the compressed frame directly decodes the transported memory, the explicit irreversible write boundary updates one entity, and canonical memory/frame state begins the next segment.

The final query uses the learned relation head already established in UP-6.

## Scientific gates

**Compression pass**

- original five-pilot held-out accuracy >= 0.95;
- compressed clean held-out accuracy >= 0.95.

**Noisy-frame pass**

- compressed reference-noise-0.01 held-out accuracy >= 0.95.

**Mutable-program pass**

- commit decode accuracy >= 0.95;
- exact final-table accuracy >= 0.90;
- relational-query accuracy >= 0.90.

Reference-noise 0.03 and transport mismatch 1e-4 / 1e-3 are diagnostic measurements and do not control the primary scientific gates.

## Interpretation boundary

A positive result would show that the internal frame can be reduced from five co-evolving vectors to three, remain useful under bounded pilot corruption, and operate inside a mutable working-memory program.

The three pilot vectors still impose runtime state overhead. UP-10 does not claim that the frame is minimal or learned.

A positive UP-10 permits UP-11: train or synthesize the compressed frame parameters from task loss, quantify minimal pilot count, and test whether frame state can be folded into the recurrent latent substrate rather than maintained as explicit side vectors.

## Authority boundary

UP-10 is mathematical research only. It does not register an inference backend, invoke a language model, activate a Wingless worker/listener, execute model-selected tools, or alter ckb-plane queue, retry, workspace, acceptance, promotion, broker, credential, or production authority.
