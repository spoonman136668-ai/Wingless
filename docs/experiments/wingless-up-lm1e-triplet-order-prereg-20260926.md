# Wingless UP-LM1E — symmetric three-family update order

Status: preregistered scientific continual language-model adaptation experiment.

Scientific parent: sealed UP-LM1D 76c653cc9a20208d42041804b54400e946c3a52e.

## Question

LM1D recovered third-family byte quality but introduced moderate base/paraphrase interference. With update counts fixed exactly, can symmetric triplet ordering reduce this recency tradeoff?

## Frozen starting point

Use the exact LM1D extended warm20 starting model:
- 64-D recurrent state;
- recurrent transition frozen;
- deterministic output-alphabet extension;
- exact recall cap 16;
- router frozen after LM1C grounding;
- learning rate 0.08;
- four adaptation epochs.

Each adaptation epoch uses exactly one base, one paraphrase, and one third-family training example for every matched corpus index.

## Frozen schedule arms

1. third_base_para
   - third, then base, then paraphrase for every index.
   - exact LM1D adapt4 control.

2. base_para_third
   - base, then paraphrase, then third for every index.

3. cyclic_by_index
   - index mod 3 = 0: third, base, paraphrase
   - index mod 3 = 1: base, paraphrase, third
   - index mod 3 = 2: paraphrase, third, base

4. cyclic_by_epoch
   - epoch mod 3 = 0: third, base, paraphrase
   - epoch mod 3 = 1: base, paraphrase, third
   - epoch mod 3 = 2: paraphrase, third, base

Every arm receives identical examples and update counts.

## Evaluation

- base held-out block;
- paraphrase held-out block;
- third-family held-out block;
- all four third-family structural order families at stream1 and stream4.

Record byte quality and the full routing/memory metrics.

## Interpretation

A symmetric schedule improving base/paraphrase retention without sacrificing third-family adaptation would identify local update recency as the remaining three-family readout bottleneck.

## Bounds

No router retraining, no recurrent training, no extra examples, no state expansion, no recall-cap increase, no attention, no adaptive ordering, no threshold tuning, no result-informed retry, no live activation, no production authority.
