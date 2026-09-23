# Wingless UP-13: task-learned phase frame

Status: research branch only; not activated, promoted, or connected to ckb-plane.

Parent research qualification: UP-12 Windows full-rank mutable-memory qualification sealed at `20ae7845042af43f60634926c96347d80e214265`.

## Question

Can the successful rank-4 internal frame begin to be learned from task loss rather than supplied as a hand-selected four-value phase alphabet?

UP-13 keeps the frame rank and entity-local support structure fixed. It learns only the four-value phase alphabet.

This is an incremental learning test, not yet full frame discovery.

## Starting point

The successful UP-9 entity pilots used manually supplied quarter-turn phases.

UP-13 deliberately starts from the poor clustered alphabet:

`[0, 0.12, 0.24, 0.36]`

The first phase is held at zero to remove the irrelevant common phase degree of freedom. The remaining three relative phases are trainable.

Before phase learning, a saturated linear readout is fitted to this clustered alphabet to measure its actual representational capacity under memory noise 0.05.

## Task-loss learning

The phase alphabet is optimized on the complete balanced 128-table training pool.

Every training state receives:

- memory perturbation amplitude 0.05;
- deterministic global-phase nuisance.

Training alternates:

1. ten deterministic full-batch softmax-head gradient steps for each entity;
2. central-difference gradients of the actual average classification cross-entropy with respect to the three trainable phases;
3. one phase update.

Parameters:

- 80 outer steps;
- head learning rate 1.0;
- phase learning rate 0.25;
- phase finite-difference epsilon 1e-4;
- two noisy training trials per table/entity.

No target phase angles are supplied.

## Fresh post-learning decoders

After the phase alphabet is learned, the alternating-training heads are discarded.

Fresh entity decoders are trained using:

- the learned phase frame;
- all 128 training tables;
- training depths 8, 24, 72, 216, 432, 648;
- four noisy trials per table/depth;
- 1,200 softmax steps.

Evaluation uses:

- all 128 disjoint held-out tables;
- unseen depths 32, 128, 512, 1024;
- two noisy trials per table/depth.

A matched non-unitary transport receives the same learned phase alphabet, full pools, decoder architecture, and optimizer budget.

## Mutable integration

The learned unitary frame and fresh decoders then run the UP-12-class mutable workload:

- 48 scenarios;
- 16 writes each;
- 768 intermediate commit opportunities;
- only held-out transport depths;
- memory noise 0.05;
- global-phase nuisance;
- explicit irreversible overwrite/re-encode boundary;
- learned final relation head.

No explicit depth, inverse transport, or runtime prototype lookup is supplied.

## Scientific gates

**Phase-learning pass**

- clustered-alphabet saturated capacity < 0.70;
- final unitary train accuracy >= 0.99;
- learned minimum phase separation >= 0.50 radians;
- final alternating-training loss < initial loss.

**Unseen-depth pass**

- learned-frame unitary held-out accuracy >= 0.99.

**Mutable-integration pass**

- commit decode accuracy >= 0.99;
- exact final-table accuracy >= 0.95;
- relational-query accuracy >= 0.95.

## Interpretation boundary

A positive result would show that task loss can transform a deliberately poor phase code into a useful internal reference alphabet that survives unseen depths and mutable integration.

The following structure would still be hand specified:

- four independent entity-local pilot supports;
- one fixed global anchor;
- rank 4 / five total pilot states.

Therefore a UP-13 pass would **not** mean the complete frame was learned from scratch.

A positive UP-13 permits UP-14: learn the pilot support/mixing directions themselves from a generic full-dimensional initialization, or fold the learned frame channels into one composite recurrent state.

## Authority boundary

UP-13 is mathematical research only. It does not register an inference backend, invoke a language model, activate a Wingless worker/listener, execute model-selected tools, or alter ckb-plane queue, retry, workspace, acceptance, promotion, broker, credential, or production authority.
