# Wingless UP-14: task-learned pilot support

Status: research branch only; not activated, promoted, or connected to ckb-plane.

Parent research qualification: UP-13 Windows learned-phase qualification sealed at `ea19f64307567fcab2853055fcabd6d18f1f5093`.

## Question

Can task loss discover **where the four rank-required pilot references should look** rather than relying on the hand-specified entity-local support used by UP-9 through UP-13?

UP-14 learns the pilot support/mixing directions while freezing:

- the rank-4 requirement established by UP-11;
- the five-pilot runtime budget;
- the global anchor;
- the four-value phase alphabet learned by UP-13.

This isolates support discovery.

## Starting support

The four entity-local phase templates from UP-13 are treated only as a spanning basis.

Each runtime pilot is initialized as an almost-uniform dense mixture of all four templates.

Before row normalization, every row contains:

- coefficient 1.04 for its nominal target entity;
- coefficient 1.00 for each of the other three entities.

Thus every pilot initially depends strongly on all four independently varying entities.

The nominal-target support purity starts below 0.30.

No target support matrix is supplied to the optimizer.

## Learned parameters

A real 4 × 4 support matrix mixes the four entity templates into four runtime pilots.

Rows are normalized after every update because overall pilot scale is unobservable after state normalization.

Each queried entity still receives only:

- its corresponding pilot/anchor correlation;
- two real features total.

A dense mixed row is therefore contaminated by the three nuisance entity values. To solve the task, support learning must make each row selective enough for its target entity.

## Task-loss optimization

The fixed value phases are the UP-13 learned values:

`[0, -1.7700278912430298, 1.1656136621055615, 2.5917033093902164]`

The complete balanced 128-table training pool is used with:

- memory noise amplitude 0.05;
- deterministic global-phase nuisance;
- two noisy trials per table/entity.

Training alternates:

1. ten full-batch linear-softmax updates for each entity;
2. central-difference gradients of that entity's classification loss with respect to its four support coefficients;
3. row normalization.

Parameters:

- 100 outer steps;
- head learning rate 1.0;
- support learning rate 0.35;
- finite-difference epsilon 1e-4.

## Fresh unseen-depth evaluation

After support learning, the alternating heads are discarded.

Fresh decoders train on:

- all 128 training tables;
- depths 8, 24, 72, 216, 432, 648;
- four noisy trials per table/depth;
- 1,200 softmax steps.

Evaluation uses:

- all 128 disjoint held-out tables;
- unseen depths 32, 128, 512, 1024;
- two noisy trials per table/depth.

The matched non-unitary path receives the identical learned support matrix, phase alphabet, pools, and decoder budget.

## Mutable integration

The learned unitary support frame then runs:

- 48 scenarios;
- 16 writes each;
- 768 commit opportunities;
- held-out transport depths only;
- memory noise 0.05;
- global-phase nuisance;
- no depth value;
- no inverse;
- no prototype lookup;
- explicit irreversible overwrite/re-encode boundary;
- learned final relation head.

## Scientific gates

**Support-learning pass**

- initial saturated capacity < 0.70;
- final unitary train accuracy >= 0.99;
- mean learned target support purity >= 0.85;
- minimum learned target support purity >= 0.75;
- final alternating-training loss < initial loss.

**Unseen-depth pass**

- learned-support unitary held-out accuracy >= 0.99.

**Mutable-integration pass**

- commit accuracy >= 0.99;
- exact final-table accuracy >= 0.95;
- relational-query accuracy >= 0.95.

## Interpretation boundary

A positive result would remove hand selection of both the phase alphabet and entity-local pilot support.

The global anchor, the existence of four pilot channels, and the rank-4/five-state frame topology would remain specified.

A positive UP-14 permits UP-15 to learn or eliminate the anchor and/or fold the five co-evolving states into one composite recurrent state so that the frame is no longer represented as explicit side vectors.

## Authority boundary

UP-14 is mathematical research only. It does not register an inference backend, invoke a language model, activate a Wingless worker/listener, execute model-selected tools, or alter ckb-plane queue, retry, workspace, acceptance, promotion, broker, credential, or production authority.
