# Wingless UP-LM0T — quadratic byte-readout anchoring

Status: preregistered scientific continual language-model adaptation experiment.

Scientific parent: sealed UP-LM0S 21c33f040b55e43fdd956c9c229d6afd4e08c419.

## Question

LM0S showed that simply masking adaptation loss to new verb targets damages both old and new byte-model performance. Can full-sentence adaptation preserve more of the original distribution by penalizing movement of the byte readout away from the pre-adaptation base model?

## Frozen starting point

- deterministic 64-D recurrent byte language model;
- 20 base training epochs;
- recurrent transition remains frozen;
- full-sentence paraphrase adaptation for exactly 4 epochs;
- learning rate 0.08;
- no base replay;
- grounded paraphrase router frozen;
- exact recall cap 16;
- no attention;
- no future oracle.

## Frozen anchoring ladder

Quadratic readout-anchor coefficient lambda:
- 0
- 0.0001
- 0.001
- 0.01

For every next-byte training step:
- compute the normal cross-entropy gradient for all output weights and biases;
- add lambda times the parameter difference from the frozen pre-adaptation base readout;
- update using the same learning rate 0.08.

The base reference is identical for every arm and never changes.

Lambda 0 is the exact full-sentence no-replay adaptation control.

## Evaluation

For every lambda:
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

A middle anchor strength that preserves more base performance while retaining paraphrase gains would establish parameter anchoring as a viable continual-language mechanism. Monotonic loss of adaptation would show that anchoring only trades one distribution for the other.

## Bounds

No base replay, no router retraining, no recurrent-state change, no recall-cap increase, no attention, no adaptive lambda, no choosing lambda after results, no result-informed retry, no live activation, no production authority.
