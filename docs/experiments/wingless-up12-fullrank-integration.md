# Wingless UP-12: full-rank decoder saturation and mutable integration

Status: research branch only; not activated, promoted, or connected to ckb-plane.

Parent research qualification: UP-11 Windows pilot-rank qualification sealed at `48db0c6fe4180a71a2f1a118eaea83691736c35f`.

## Question

UP-11 established that four independent reference-code directions are required to clear the current unseen-depth accuracy gate.

Its noisy rank-4 unitary observer reached 0.95703125 held-out accuracy using:

- 32 training tables;
- 32 held-out tables;
- one noisy sample per table/depth;
- 250 linear-softmax optimization steps.

Is the remaining ~4.3% error caused by sparse decoder training rather than missing representational rank?

## Frozen representation

UP-12 does not change the successful rank-4 reference representation.

It keeps:

- one common anchor;
- four equal-energy Fourier code directions;
- five total pilot states;
- 80 runtime complex pilot scalars;
- eight real pilot/anchor correlation features per queried entity;
- unitary co-evolution;
- no explicit depth;
- no inverse;
- no runtime prototype lookup.

This experiment therefore tests decoder/data saturation, not another frame architecture.

## Full balanced pools

All 256 possible four-entity memory tables are used.

The balanced combination rule from UP-8 partitions them into:

- 128 training tables;
- 128 disjoint held-out tables.

Every entity/value marginal occurs exactly 32 times in each pool.

## Depths and perturbations

Training depths:

`8, 24, 72, 216, 432, 648`

Held-out depths:

`32, 128, 512, 1024`

Memory perturbation amplitude remains 0.05.

Every memory also receives a deterministic global-phase nuisance.

For each table/depth:

- four independent deterministic noisy training trials are generated;
- two independent deterministic noisy held-out trials are generated.

## Decoder budget

Each entity keeps a separate linear-softmax decoder over the same eight full-rank features.

Training budget:

- 1,200 deterministic full-batch gradient steps;
- learning rate 1.0.

The matched non-unitary rank-4 transport receives the identical table pools, trials, feature count, decoder architecture, and optimization budget.

## Static scientific gate

The rank-4 unitary representation is considered decoder/data saturated when held-out accuracy reaches at least 0.99.

The experiment records the gain relative to the UP-11 noisy rank-4 result of 0.95703125.

If the saturated linear decoder remains below 0.99, UP-13 should test nonlinear readout capacity rather than further increasing pilot rank.

## Mutable integration

The trained unitary rank-4 heads are placed into a 48-scenario mutable working-memory workload.

Each scenario contains:

- 16 writes/overwrites;
- only held-out transport depths;
- memory noise amplitude 0.05;
- deterministic memory-only global-phase nuisance;
- direct co-evolving rank-4 frame readout at every commit;
- the explicit irreversible overwrite/re-encode boundary;
- final learned relational query.

This creates 768 sequential commit opportunities.

## Mutable scientific gate

Pass requires:

- commit decode accuracy >= 0.99;
- exact final-table accuracy >= 0.95;
- relational-query accuracy >= 0.95.

The mutable gate is intentionally stricter on commit accuracy because small per-step errors compound through subsequent writes.

## Interpretation boundary

If static saturation passes and mutable integration passes, the full-rank frame is sufficient under the current noise/depth regime and the next research step should make the frame basis learnable/foldable rather than altering its rank.

If static saturation fails, the next bottleneck is decoder capacity or feature nonlinearity.

If static saturation passes but mutable integration fails, the residual per-step error/margin is still too weak for recurrent working memory and should be attacked directly.

## Authority boundary

UP-12 is mathematical research only. It does not register an inference backend, invoke a language model, activate a Wingless worker/listener, execute model-selected tools, or alter ckb-plane queue, retry, workspace, acceptance, promotion, broker, credential, or production authority.
