# Wingless UP-LM0Z — frozen-representation joint readout capacity

Status: preregistered scientific representation-capacity diagnostic.

Scientific parent: sealed UP-LM0Y fc8d640b19fdd89fcdac2b79e423df2fbdcce6a6.

## Question

LM0X/LM0Y indicate mostly aligned base/paraphrase gradients, sparse row conflict, and no joint gain from pairwise gradient correction. Can the frozen recurrent representation support both distributions simultaneously when a fresh readout is trained jointly from scratch, or is the tradeoff a representational/readout-capacity boundary?

## Frozen representation and corpora

- use the exact deterministic 64-D recurrent transition from the LM lineage;
- recurrent transition is never trained;
- exact LM0O base and paraphrase train/held-out corpora;
- same byte alphabet;
- no exact recall, router, or memory mechanism is involved in this diagnostic;
- all readouts start from zero weights and zero biases.

## Frozen readout arms

1. base_only
   - train fresh readout for 20 epochs on the base training corpus using the lineage SGD rule at lr 0.08.

2. paraphrase_only
   - train fresh readout for 20 epochs on the paraphrase training corpus at lr 0.08.

3. joint_interleaved
   - train fresh readout for 20 epochs;
   - for every matched training index: base example then paraphrase example;
   - lr 0.08 per example.

4. joint_fullbatch
   - train fresh readout for 100 deterministic full-batch epochs;
   - at each epoch compute the average token-level cross-entropy gradient over the complete union of base + paraphrase training corpora at unchanged epoch-start weights;
   - apply one update with learning rate 0.08.
   - no example-order effect exists inside the epoch.

## Evaluation

Every arm is evaluated on:
- base held-out top-1 byte accuracy and perplexity;
- paraphrase held-out top-1 byte accuracy and perplexity.

Also report:
- train-set top-1 accuracy and perplexity on each distribution;
- readout parameter count;
- recurrent state dimension;
- whether recurrent parameters were trained (must be false).

## Interpretation

A joint arm performing strongly on both distributions would show that the frozen recurrent representation/readout capacity is adequate and the prior tradeoff is primarily continual optimization. A persistent joint tradeoff despite direct joint training would implicate representation/readout capacity and justify moving beyond readout-only adaptation.

## Bounds

No recurrent training, no router, no exact recall, no attention, no anchoring, no gradient projection, no adaptive epoch count, no result-informed retry, no live activation, no production authority.
