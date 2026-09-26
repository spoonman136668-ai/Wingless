# Wingless UP-LM0S — lexical-target loss masking

Status: preregistered scientific continual language-model adaptation experiment.

Scientific parent: sealed UP-LM0R d75ed3c7f6d3d0447b475477bab3496295cd5188.

## Question

LM0R showed that update order only modestly shifts the old-versus-new byte-model tradeoff. Can adaptation be made less destructive by applying gradient updates only when the target byte belongs to a new paraphrase verb, while still processing the full sentence through the frozen recurrent state?

## Frozen starting point

Identical to LM0R:
- deterministic 64-D recurrent byte language model;
- 20 base training epochs;
- learning rate 0.08;
- paraphrase corpus uses saves / sees / recalls;
- grounded paraphrase router frozen;
- exact recall cap 16;
- no attention;
- no future oracle.

## Frozen adaptation depth

- exactly 4 paraphrase adaptation epochs.
- no base replay in either arm.
- every arm sees the identical paraphrase training sentences in identical order.

## Frozen arms

1. full_sentence
   - exact full-sentence paraphrase adaptation used by the LM0P replay-0 control.
   - update output weights for every next-byte target.

2. verb_targets_only
   - process every byte of the full paraphrase sentence through the recurrent state exactly as normal.
   - update output weights only when the next-byte target position lies inside one of the verb tokens saves, sees, or recalls.
   - spaces and punctuation surrounding the verb are not update targets.
   - all non-verb target positions contribute context but no gradient update.

## Evaluation

For each arm:
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

Also report the exact number of gradient-updated target bytes per adaptation epoch for each arm.

## Interpretation

If verb-target masking improves paraphrase competence while preserving substantially more base performance, the current interference is caused by unnecessary full-sentence readout updates rather than the recurrent representation itself.

## Bounds

No router retraining, no base replay, no state expansion, no recall-cap increase, no attention, no threshold tuning, no adaptive mask, no result-informed retry, no live activation, no production authority.
