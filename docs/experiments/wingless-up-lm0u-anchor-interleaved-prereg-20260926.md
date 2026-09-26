# Wingless UP-LM0U — weak anchoring plus interleaved replay

Status: preregistered scientific continual language-model integration experiment.

Scientific parent: sealed UP-LM0T 549017ab748a2fffffb49807088c1ad624d3c8ef.

## Question

LM0T showed that weak readout anchoring improves base retention with only a small paraphrase cost, while LM0Q showed that example-interleaved base replay gives a strong joint balance. Do these mechanisms compose, or are they redundant forms of the same retention pressure?

## Frozen starting point

- deterministic 64-D recurrent byte language model;
- 20 base training epochs;
- 4 adaptation epochs;
- learning rate 0.08;
- grounded paraphrase router frozen;
- exact recall cap 16;
- no attention;
- no future oracle.

## Frozen 2x2 arms

Replay schedule:
- none
- example_interleaved: for each matched training index, paraphrase example then base example.

Anchor coefficient:
- lambda 0
- lambda 0.0001

This yields four arms:
1. no_replay_anchor0
2. no_replay_anchor0001
3. interleaved_anchor0
4. interleaved_anchor0001

For anchored arms, the quadratic readout penalty is applied on every training update in that arm relative to the same frozen pre-adaptation base readout.

The interleaved lambda-0 arm reproduces the LM0Q paraphrase-first paired control; the no-replay lambda arms reproduce LM0T controls.

## Evaluation

For every arm:
1. base held-out top-1 accuracy and perplexity;
2. paraphrase block held-out top-1 accuracy and perplexity;
3. paraphrase block plus all four LM0N order families at stream1:
   - top-1 byte accuracy;
   - perplexity;
   - dependent first-byte accuracy;
   - whole-query-set exact accuracy;
   - event-routing accuracy;
   - report-routing accuracy;
   - maximum recall entries.

## Interpretation

If weak anchoring improves the interleaved arm on base retention without materially reducing paraphrase performance, the mechanisms are complementary. Little or negative change would show that interleaved replay already supplies the useful retention constraint.

## Bounds

No router retraining, no adaptive lambda, no additional replay examples beyond the defined interleaved arm, no state expansion, no recall-cap increase, no attention, no result-informed retry, no live activation, no production authority.
