# Wingless UP-LM0R — paired update-order symmetry

Status: preregistered scientific continual language-model adaptation experiment.

Scientific parent: sealed UP-LM0Q 530f5fd1e6f86578525f0d5ae105cc518c6af1d3.

## Question

LM0Q showed that example-level interleaving balances base retention and paraphrase adaptation better than either full-corpus block order. Within each matched base/paraphrase example pair, does pair-level recency still control the remaining tradeoff?

## Frozen starting point

Identical to LM0Q:
- deterministic 64-D recurrent byte language model;
- 20 base training epochs;
- 4 adaptation epochs;
- learning rate 0.08;
- one complete base corpus and one complete paraphrase corpus worth of updates per adaptation epoch;
- grounded paraphrase router frozen;
- exact recall cap 16;
- no attention;
- no future oracle.

Base and paraphrase training corpora have equal counts and remain in deterministic corresponding index order.

## Frozen schedule arms

1. paraphrase_first_pairs
   - for every index i: paraphrase example i, then base example i.
   - exact LM0Q example_interleaved control.

2. base_first_pairs
   - for every index i: base example i, then paraphrase example i.

3. alternating_pair_order
   - for even i: paraphrase example i, then base example i;
   - for odd i: base example i, then paraphrase example i.
   - pair order is fixed by index and does not depend on results or epoch.

Every arm receives exactly the same examples and update count.

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

Improvement from alternating pair order would demonstrate that even example-level replay retains a local recency bias that can be reduced without adding updates.

## Bounds

No router retraining, no extra examples, no adaptive ordering, no state expansion, no recall-cap increase, no attention, no threshold tuning, no result-informed retry, no live activation, no production authority.
