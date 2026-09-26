# Wingless UP-LM1A — initialization-history versus joint-training budget

Status: preregistered scientific continual language optimization diagnostic.

Scientific parent: sealed/corrected UP-LM0Z 761ed2d9c275d6707602d5b658bf2e1d7e46d253.

## Question

LM0Z showed that a fresh jointly trained readout can represent base and paraphrase distributions simultaneously, so the prior tradeoff is not a hard frozen-representation capacity limit. Is the remaining gap caused mainly by base-pretrained initialization history, or by too few joint-training epochs during continual adaptation?

## Frozen representation and corpora

- exact deterministic 64-D recurrent transition;
- recurrent parameters never train;
- exact LM0O base and paraphrase train/held-out corpora;
- same byte alphabet;
- output readout only;
- learning rate 0.08;
- matched joint training order:
  for every training index, base example then paraphrase example;
- no router, exact recall, attention, anchoring, projection, or adaptive scheduling.

## Frozen arms

1. fresh_4
   - zero-initialized readout;
   - 4 joint-interleaved epochs.

2. warm_4
   - first train 20 base-only epochs;
   - then 4 joint-interleaved epochs.

3. fresh_20
   - zero-initialized readout;
   - 20 joint-interleaved epochs.
   - exact LM0Z joint_interleaved training control.

4. warm_20
   - first train 20 base-only epochs;
   - then 20 joint-interleaved epochs.

The fresh-vs-warm comparison is made at identical joint-training budgets (4 vs 4 and 20 vs 20). Warm arms intentionally contain prior base history; that is the experimental factor.

## Evaluation

Every arm is evaluated on:
- base train and held-out top-1 accuracy and perplexity;
- paraphrase train and held-out top-1 accuracy and perplexity.

Also report:
- readout parameter count;
- recurrent state dimension;
- recurrent parameters trained (must be false).

## Interpretation

If warm and fresh converge at 20 joint epochs, the earlier continual gap is primarily insufficient joint optimization time. If warm remains systematically base-biased at equal joint budget, initialization/history has a persistent effect.

## Bounds

No recurrent training, no adaptive epoch count, no router, no exact recall, no attention, no anchoring, no projection, no result-informed retry, no live activation, no production authority.
