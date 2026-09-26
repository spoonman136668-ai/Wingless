# Wingless UP-LM1N — four-family operator-balance scheduling

Status: preregistered scientific multi-family byte-readout optimization experiment.

Scientific parent: sealed UP-LM1M a87c788a0690caae6bb47fe84b9ba74f1935a61b.

## Question

LM1M showed that fourth-only adaptation learns the fourth lexical family strongly but catastrophically harms the prior three families, while simple four-family interleaving preserves prior families better at lower fourth-family quality. Can update ordering alone improve the four-family balance when examples, recurrent state, router, recall capacity, and nominal per-family learning-rate mass are held fixed?

## Frozen starting model

Exact LM1M three-family starting point:
- 64-D recurrent transition;
- recurrent parameters frozen;
- 20 base-only epochs;
- 20 base/paraphrase joint-interleaved epochs;
- three-family Strang-split-base adaptation for 4 epochs;
- output alphabet already extended for third/fourth byte surfaces;
- corrected LM1L fourth-family router is not involved in byte-readout training.

## Frozen adaptation budget

All arms:
- 4 adaptation epochs;
- matched base / paraphrase / third / fourth training examples;
- learning rate 0.08 for a full update;
- half step = 0.04;
- each matched index contributes nominal lr 0.08 to every family per epoch;
- gradients are recomputed after every sub-step;
- no extra examples.

## Arms

1. cyclic_control
   - base full step;
   - paraphrase full step;
   - third full step;
   - fourth full step.
   - exact LM1M four_family_interleaved_4ep control.

2. mirrored_by_index
   - even example index: base, paraphrase, third, fourth;
   - odd example index: fourth, third, paraphrase, base.
   - every update is a full 0.08 step.

3. rotating_palindromic_split
   - center family rotates by example index modulo 4;
   - the three non-center families each receive a 0.04 half step before center, in cyclic family order;
   - center receives one 0.08 full step;
   - the non-center half steps are then replayed in exact reverse order.
   - Thus every family receives total nominal lr 0.08 per matched example, while local last-update privilege rotates across all four families.

Family order is fixed as base, paraphrase, third, fourth.

## Evaluation

For each arm:
- base / paraphrase / third / fourth train and held-out top-1 byte accuracy;
- held-out perplexity per family;
- minimum held-out family accuracy;
- mean held-out family accuracy;
- held-out family spread.

No semantic router or exact-memory override is used in this probe; this experiment isolates byte-readout optimization.

## Interpretation

If mirrored or rotating-palindromic ordering improves minimum/mean held-out balance over cyclic control without changing update mass, the remaining four-family interference is at least partly an operator-order / recency effect. If not, the evidence shifts toward representational/readout capacity interference.

## Bounds

No recurrent training, no router retraining, no recall-cap change, no extra examples, no adaptive ordering, no learning-rate search, no result-informed retry, no live activation, no production authority.
