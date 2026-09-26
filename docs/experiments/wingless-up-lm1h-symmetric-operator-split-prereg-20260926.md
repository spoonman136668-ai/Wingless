# Wingless UP-LM1H — symmetric operator splitting

Status: preregistered scientific three-family readout optimization experiment.

Scientific parent: sealed UP-LM1G 4d7ad1c003f23c67749c229ac766c2a7353ac41c.

## Question

LM1G showed that simple symmetric gradient scaling recovers part of the third-family loss but still underperforms sequential cyclic updates. Can symmetric operator splitting preserve nonlinear sequential optimization while reducing triplet recency bias?

## Frozen model

- same 64-D recurrent transition;
- recurrent parameters frozen;
- same output alphabet;
- same base/paraphrase/third training corpora;
- same base warm start;
- no router or recall involvement in this probe;
- learning rate 0.08.

## Frozen arms

1. cyclic_by_index
   - exact LM1E/LM1G control.

2. stranged_split_base
   - for each matched triplet: base half-step, paraphrase half-step, third full-step, paraphrase half-step, base half-step.
   - gradients recomputed after every sub-step.

3. stranged_split_para
   - paraphrase half-step, third half-step, base full-step, third half-step, paraphrase half-step.

4. stranged_split_third
   - third half-step, base half-step, paraphrase full-step, base half-step, third half-step.

Half-step uses lr 0.04; full-step uses lr 0.08.
All arms run 20 epochs and consume identical matched example triplets.

## Evaluation

Base/paraphrase/third train and held-out:
- top-1 accuracy;
- perplexity.

Report:
- minimum held-out accuracy across the three families;
- mean held-out accuracy;
- family spread.

## Interpretation

A symmetric split that improves worst-family and mean balance over cyclic control would support nonlinear operator ordering as the remaining mechanism. Failure would point toward true readout/representation capacity interference.

## Bounds

No recurrent training, no extra examples, no adaptive schedule, no threshold tuning, no result-informed retry, no live activation, no production authority.
