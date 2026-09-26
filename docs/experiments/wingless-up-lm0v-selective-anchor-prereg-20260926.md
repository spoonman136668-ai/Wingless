# Wingless UP-LM0V — selective byte-readout anchoring

Status: preregistered scientific continual language-model adaptation experiment.

Scientific parent: sealed UP-LM0U 3939e1f232c159bad6b0f90297d60ef0b3201f39.

## Question

LM0U showed that full-readout weak anchoring adds little base retention to example-interleaved replay while materially reducing paraphrase accuracy. Can anchoring only byte-output rows specific to the original event-verb surfaces preserve useful base structure without constraining the new paraphrase surfaces?

## Frozen starting point

Identical to the LM0U interleaved arms:
- deterministic 64-D recurrent byte language model;
- 20 base training epochs;
- 4 adaptation epochs;
- learning rate 0.08;
- example-interleaved replay: paraphrase example i, then base example i;
- grounded paraphrase router frozen;
- exact recall cap 16;
- no attention;
- no future oracle.

## Frozen anchor coefficient

- lambda = 0.0001 for anchored arms.

## Frozen row mask

Original event verbs:
- stores
- observes
- reports

Paraphrase event verbs:
- saves
- sees
- recalls

Construct the anchor byte set before training as:
- bytes that occur in at least one original event verb;
- excluding any byte that occurs in any paraphrase event verb.

For the frozen strings above this yields:
- t
- o
- b
- p

Selective anchoring applies the quadratic penalty only to output rows for those bytes.
All other output rows receive normal data gradients only.

## Frozen arms

1. interleaved_anchor0
   - exact LM0U interleaved lambda-0 control.

2. interleaved_full_anchor0001
   - exact LM0U full-readout lambda-0.0001 control.

3. interleaved_selective_anchor0001
   - same interleaved updates and lambda;
   - anchoring only on output rows t, o, b, p.

Every arm receives identical examples and update counts.

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

Selective anchoring that improves base retention without the full-anchor paraphrase penalty would show that retention pressure can be localized to old lexical readout structure. No benefit would argue against readout anchoring as the next mechanism.

## Bounds

No router retraining, no adaptive row selection, no additional examples, no state expansion, no recall-cap increase, no attention, no lambda tuning after results, no result-informed retry, no live activation, no production authority.
