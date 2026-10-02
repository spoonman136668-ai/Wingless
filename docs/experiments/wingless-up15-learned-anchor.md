# Wingless UP-15: task-learned global anchor

Status: Windows-qualified task-learned anchor result; not activated, promoted, or connected to ckb-plane.

Parent research qualification: UP-14 Windows learned-support qualification sealed at `eb8affea150f4564502145bc50607ef863b2f425`.

## Question

Can task loss learn the global reference anchor itself rather than relying on the uniform hand-designed anchor used by UP-9 through UP-14?

UP-15 freezes the two structures already learned by prior experiments:

- UP-13 four-value phase alphabet;
- UP-14 pilot support matrix.

It learns only the global anchor.

## Starting anchor

The anchor is a full dense 16-dimensional complex state.

It starts from a deterministic irregular amplitude/phase pattern rather than the uniform real anchor used by the successful earlier frame.

The initial anchor is normalized and gauge-fixed so coordinate zero is real-positive.

That gauge removes irrelevant scale/global-phase freedom. The remaining anchor geometry is trainable.

## Learned parameters

All sixteen complex anchor coordinates are optimized.

Equivalent parameter count:

- 16 real components;
- 16 imaginary components;
- normalization removes scale;
- deterministic gauge removes one global phase.

No target anchor coordinates are supplied.

## Task-loss learning

Training uses the full balanced 128-table training pool with:

- memory noise amplitude 0.05;
- deterministic memory-only global phase;
- two noisy trials per table/entity.

The frozen entity pilots use:

- UP-13 learned phase alphabet;
- UP-14 learned support matrix.

Training alternates:

1. ten full-batch softmax updates for each entity;
2. central-difference gradients of mean classification cross-entropy over all 32 raw anchor parameters;
3. anchor normalization and deterministic gauge fixing.

Parameters:

- 80 outer steps;
- head learning rate 1.0;
- anchor learning rate 0.18;
- finite-difference epsilon 1e-4.

The learned anchor is also evaluated by its table-phase concentration: the circular concentration of anchor/memory correlation phase across all 256 clean memory tables. This is diagnostic, not a target supplied to learning.

## Fresh unseen-depth evaluation

After anchor learning, all alternating-training heads are discarded.

Fresh decoders train on:

- all 128 training tables;
- depths 8, 24, 72, 216, 432, 648;
- four noisy trials per table/depth;
- 1,200 softmax steps.

Evaluation uses:

- all 128 disjoint held-out tables;
- unseen depths 32, 128, 512, 1024;
- two noisy trials per table/depth.

The matched non-unitary path receives the identical learned anchor, phase alphabet, support matrix, pools, and decoder budget.

## Mutable integration

The learned unitary anchor/frame then runs:

- 48 scenarios;
- 16 writes each;
- 768 commit opportunities;
- held-out transport depths only;
- memory noise 0.05;
- global-phase nuisance;
- no depth value;
- no inverse;
- no runtime prototype lookup;
- explicit overwrite/re-encode boundary;
- learned final relation head.

## Scientific gates

**Anchor-learning pass**

- initial saturated capacity < 0.70;
- final unitary train accuracy >= 0.99;
- final alternating-training loss < initial loss;
- learned anchor state moves at least 0.20 L2 from its initialization.

No similarity to the historical uniform anchor is required.

**Unseen-depth pass**

- learned-anchor unitary held-out accuracy >= 0.99.

**Mutable-integration pass**

- commit accuracy >= 0.99;
- exact final-table accuracy >= 0.95;
- relational-query accuracy >= 0.95.

## Interpretation boundary

A positive result would remove manual specification of:

- the value-phase alphabet;
- the pilot support matrix;
- the global anchor geometry.

The experiment would still retain the explicit five-state frame topology: one anchor channel plus four rank-required pilot channels.

That topology becomes the UP-16 boundary: determine whether those five co-evolving states can be packed into one composite recurrent latent state without losing direct readout, rather than existing as explicit side vectors.

## Authority boundary

UP-15 is mathematical research only. It does not register an inference backend, invoke a language model, activate a Wingless worker/listener, execute model-selected tools, or alter ckb-plane queue, retry, workspace, acceptance, promotion, broker, credential, or production authority.


## Windows qualification

Authoritative operator proof on 2026-09-23 against source head `062eabadceea220e17398911378789228f1d702d` passed the harness and complete Wingless regression.

The irregular starting anchor had saturated capacity 0.447265625 and table-phase concentration 0.19313507842984184. Task-loss optimization moved the anchor by L2 distance 1.2810354629307779, raised table-phase concentration to 0.9893853411379732, and reduced average alternating-training loss from 1.2161035415312518 to 0.023262457632572883.

Fresh unitary decoders reached 1.0 train and 1.0 unseen-depth held-out accuracy. The matched non-unitary path reached 0.763427734375 held-out accuracy.

Mutable integration used 48 scenarios × 16 writes and achieved 1.0 commit decode, 1.0 exact final-table, and 1.0 relational-query accuracy. Minimum value margin was 0.17830499668449967 and maximum norm drift was 6.750155989720952e-14.

Interpretation: UP-15 removes manual specification of the anchor geometry. The remaining explicit design is the multi-state frame topology itself.
