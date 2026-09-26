# Wingless UP-LM1R — five-family readout-gradient conflict attribution

Status: preregistered scientific readout-geometry diagnostic.

Scientific parent: sealed UP-LM1Q 6c2111454a22feb60a36448effcc21997b4403e3.

## Question

UP-LM1Q showed that the fifth lexical family remains the weakest family under every cyclic update rotation, so simple update position does not explain its disadvantage. At the exact five-family pre-adaptation starting point, is the fifth family's output-readout gradient uniquely anti-aligned with the prior four families?

## Frozen starting model

Use the exact UP-LM1Q starting point before five-family adaptation:
- 64-D recurrent transition;
- recurrent parameters frozen;
- exact four-family training history ending with the four-family rotating-palindromic schedule;
- output alphabet extended for fifth-family bytes;
- no five-family adaptation updates.

Matched training corpora:
1. base
2. paraphrase
3. third
4. fourth
5. fifth

All five training corpora have identical example counts.

## Gradient scope

Measure gradients only for the byte-model output readout:
- output weight matrix;
- output bias vector.

Do not compute or update recurrent gradients.

For every family:
- sum the per-sentence readout gradients over the full matched training corpus;
- report aggregate gradient norm;
- report frozen starting held-out top-1 accuracy and perplexity.

For every unordered family pair (10 pairs):
- aggregate-gradient cosine;
- aggregate dot product;
- mean matched-example cosine;
- minimum and maximum matched-example cosine;
- fraction of matched examples with negative cosine;
- fraction with non-positive cosine.

Derived summaries:
- mean aggregate cosine for fifth-vs-prior pairs;
- minimum fifth-vs-prior aggregate cosine;
- mean aggregate cosine among prior-family pairs;
- fifth aggregate gradient norm relative to the prior-family mean norm.

## Interpretation

If fifth-vs-prior gradients are systematically more anti-aligned than prior-vs-prior gradients, five-family weakness is attributable to output-readout gradient conflict. If not, the evidence shifts toward representational difficulty or initialization/history effects rather than direct gradient opposition.

## Bounds

No parameter updates after constructing the frozen starting model, no router, no exact recall, no recurrent gradients, no adaptive weighting, no learning-rate change, no sixth family, no result-informed retry, no live activation, no production authority.
