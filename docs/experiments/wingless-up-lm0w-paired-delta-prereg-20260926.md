# Wingless UP-LM0W — order-neutral paired delta composition

Status: preregistered scientific continual language-model adaptation experiment.

Scientific parent: sealed UP-LM0V 7d8d743c05ef86806713f26e1f41146618ce46bc.

## Question

LM0Q/LM0R showed that base-versus-paraphrase recency matters, while LM0V showed that readout anchoring does not improve the interleaved tradeoff. Can we remove between-example recency directly by combining independently computed base and paraphrase parameter deltas from the same pre-pair model?

## Frozen starting point

- deterministic 64-D recurrent byte language model;
- 20 base training epochs;
- 4 adaptation epochs;
- learning rate 0.08;
- one matched paraphrase example and one matched base example per index per adaptation epoch;
- grounded paraphrase router frozen;
- exact recall cap 16;
- no attention;
- no future oracle;
- no anchoring.

## Frozen arms

1. paraphrase_then_base
   - for each matched index: normal paraphrase trainSentence update, then normal base trainSentence update.

2. base_then_paraphrase
   - base trainSentence, then paraphrase trainSentence.

3. paired_delta_sum
   - clone the current model twice from the same pre-pair state;
   - train the paraphrase example normally on clone P at lr 0.08;
   - train the base example normally on clone B at lr 0.08;
   - compute both complete parameter deltas relative to the shared pre-pair model;
   - apply the sum of those two deltas once to the current model;
   - proceed to the next pair.

This preserves each example's internal token-level SGD while eliminating order between the two examples.

Every arm sees identical examples and nominal per-example learning rate.

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

If paired_delta_sum improves the joint base/paraphrase balance relative to both sequential orders, the remaining tradeoff is partly a between-example recency artifact. Failure would indicate genuinely conflicting gradients rather than ordering alone.

## Bounds

No router retraining, no additional examples, no anchoring, no adaptive mixing, no state expansion, no recall-cap increase, no attention, no threshold tuning, no result-informed retry, no live activation, no production authority.
