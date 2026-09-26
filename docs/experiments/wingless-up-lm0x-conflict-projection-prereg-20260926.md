# Wingless UP-LM0X — conflict-aware paired update projection

Status: preregistered scientific continual language-model adaptation experiment.

Scientific parent: sealed UP-LM0W 205791da63c3ce4a214d44006ab77564e93c305e.

## Question

LM0W showed that order-neutral paired deltas land between the base-first and paraphrase-first controls rather than improving both, supporting genuinely conflicting update directions. Can projecting conflicting paired deltas reduce that tradeoff without adding examples or replay?

## Frozen starting point

- deterministic 64-D recurrent byte language model;
- 20 base training epochs;
- 4 adaptation epochs;
- learning rate 0.08;
- matched base/paraphrase example pairs in deterministic index order;
- grounded paraphrase router frozen;
- exact recall cap 16;
- no attention;
- no future oracle;
- no anchoring.

For every pair, independently train clone P on the paraphrase example and clone B on the base example from the same pre-pair model, exactly as in LM0W. Let dP and dB be the resulting complete parameter deltas.

## Frozen arms

1. paired_delta_sum
   - exact LM0W control: apply dP + dB.

2. base_protected_projection
   - if dot(dP,dB) >= 0, apply dP + dB;
   - if dot(dP,dB) < 0:
     - project dP off dB:
       dP' = dP - dot(dP,dB)/||dB||^2 * dB
     - apply dB + dP'.

3. symmetric_projection
   - if dot(dP,dB) >= 0, apply dP + dB;
   - if dot(dP,dB) < 0:
     - dP' = dP - dot/||dB||^2 * dB
     - dB' = dB - dot/||dP||^2 * dP
     - apply dP' + dB'.

All parameter vectors include readout weights and biases.

## Diagnostics

For every arm:
- total pair updates;
- number and fraction of pre-projection pairs with negative delta dot product;
- mean pre-projection cosine similarity.

## Evaluation

For every arm:
1. base held-out top-1 accuracy and perplexity;
2. paraphrase block held-out top-1 accuracy and perplexity;
3. paraphrase block plus all four LM0N structural order families at stream1:
   - top-1 byte accuracy;
   - perplexity;
   - dependent first-byte accuracy;
   - whole-query-set exact accuracy;
   - event-routing accuracy;
   - report-routing accuracy;
   - maximum recall entries.

## Interpretation

Joint improvement under a projection arm would support gradient conflict as an actionable mechanism. No improvement would imply that the two distributions require a different capacity or representation strategy rather than local gradient surgery.

## Bounds

No router retraining, no extra examples, no anchoring, no adaptive arm selection, no state expansion, no recall-cap increase, no attention, no threshold tuning, no result-informed retry, no live activation, no production authority.
