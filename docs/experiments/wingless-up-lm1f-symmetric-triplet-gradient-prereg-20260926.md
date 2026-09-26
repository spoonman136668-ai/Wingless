# Wingless UP-LM1F — symmetric three-family gradient update

Status: preregistered scientific continual language-model adaptation experiment.

Scientific parent: sealed UP-LM1E 1968a301cad7ad318713143037ede9113fa43a0a.

## Question

LM1E showed that sequential order shifts base/paraphrase/third-family byte quality and no tested ordering improves all three simultaneously. Does removing within-triplet recency by averaging the three readout gradients improve the joint balance?

## Frozen starting point

Use the exact LM1E extended warm20 starting model:
- 64-D recurrent state;
- recurrent transition frozen;
- deterministic output alphabet extension;
- exact recall cap 16;
- grounded router frozen;
- learning rate 0.08;
- four adaptation epochs.

Matched training corpora:
- base;
- paraphrase;
- third lexical family.

## Frozen arms

1. cyclic_by_index
   - exact LM1E sequential control.

2. symmetric_triplet_average
   - for matched corpus index i, compute the readout gradient for base[i], paraphrase[i], and third[i] from the same pre-update model parameters;
   - arithmetic-average the three gradients;
   - apply exactly one readout update at learning rate 0.08;
   - proceed to the next index.

3. symmetric_triplet_sum
   - same simultaneous gradient calculation;
   - sum the three gradients and apply with learning rate 0.08 / 3, making the effective average scale identical while preserving a separate implementation check.

Arms receive identical data and effective gradient scale.

## Evaluation

- base held-out block;
- paraphrase held-out block;
- third held-out block;
- all third-family structural order families at stream1 and stream4.

Record byte quality and full routing/memory metrics.

## Interpretation

If symmetric gradient updates improve all three distributions or reduce the spread without harming exact routing/memory, the remaining tradeoff is primarily sequential-gradient interference rather than readout capacity.

## Bounds

No recurrent/router training, no extra examples, no state/recall expansion, no attention, no adaptive weighting, no result-informed retry, no live activation, no production authority.
