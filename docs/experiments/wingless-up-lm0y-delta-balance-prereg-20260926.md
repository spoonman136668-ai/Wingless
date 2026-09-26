# Wingless UP-LM0Y — paired-delta norm balance and row projection

Status: preregistered scientific continual language-model mechanism experiment.

Scientific parent: sealed UP-LM0X 8b4bbbd93538f1f001f73f7e888e901769868a7e.

## Question

LM0X found mostly positive matched gradients but a roughly 3.55x paraphrase/base gradient-norm imbalance plus sparse negative conflict in four output rows. Can order-neutral paired adaptation improve joint retention by correcting scale imbalance, and does additionally removing only negative row conflict help?

## Frozen starting point

- exact LM0W base model and matched base/paraphrase training corpus;
- 20 base epochs;
- 4 adaptation epochs;
- nominal per-example learning rate 0.08;
- independent base and paraphrase one-example deltas are always computed from the same pre-pair model, preserving each example's internal token-level SGD;
- grounded paraphrase router frozen;
- exact recall cap 16;
- no attention;
- no anchoring;
- no future oracle.

## Frozen arms

1. paired_delta_sum
   - exact LM0W order-neutral control:
   - apply base delta + paraphrase delta.

2. norm_balanced_delta
   - compute full parameter norms of the independent base and paraphrase deltas;
   - multiply the paraphrase delta by base_norm / paraphrase_norm when paraphrase_norm > 0;
   - apply base delta + scaled paraphrase delta.

3. norm_balanced_row_projected
   - start from the same norm-balanced paraphrase delta;
   - for each output-byte row independently, if the row dot product with the base delta row is negative, project only the scaled paraphrase row onto the orthogonal complement of the base row;
   - leave nonnegative rows unchanged;
   - apply base delta + resulting paraphrase delta.

No row list is hard-coded from LM0X; projection is determined only by the current pair's row dot sign.

## Evaluation

For every arm:
- base held-out top-1 accuracy and perplexity;
- paraphrase block held-out top-1 accuracy and perplexity;
- paraphrase block plus all four structural order families at stream1:
  - top-1 byte accuracy;
  - perplexity;
  - dependent first-byte accuracy;
  - whole-query-set exact accuracy;
  - event-routing accuracy;
  - report-routing accuracy;
  - maximum recall entries.

Also report adaptation diagnostics:
- mean paraphrase/base raw delta norm ratio;
- number of pair-row projections applied.

## Interpretation

A gain from norm balancing would identify gradient-scale imbalance as causal. Additional gain from row projection would show that sparse negative rows contribute beyond scale. Failure of both would redirect the lane away from pairwise readout-gradient correction.

## Bounds

No router retraining, no extra examples, no anchoring, no adaptive arm selection, no state expansion, no recall-cap increase, no attention, no threshold tuning, no result-informed retry, no live activation, no production authority.
