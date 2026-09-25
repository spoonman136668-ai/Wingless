# Wingless UP-LM0H — whole lexical-family admission holdout

Status: preregistered scientific language/memory integration experiment.

Scientific parent: sealed UP-LM0G `e2b6dd94977bfde44c80f91480a308aae5d943ee`.

## Question

UP-LM0G reached perfect admission on structurally held-out name/value recombinations when every name and value appeared during classifier training. Does the learned admission rule still generalize when entire subject and value families are absent from admission-classifier training?

## Frozen language system

The language model and paragraph corpus remain identical to UP-LM0G:
- deterministic 64-D recurrent state;
- 512 recurrent-state bytes;
- 20 epochs;
- learning rate 0.08;
- no attention;
- exact recall cap 16;
- same train/held-out paragraph generator.

Only the admission-classifier training domain changes.

## Frozen admission classifier

Same deterministic 64-D event encoder, binary linear logistic head, zero initialization, 20 epochs, learning rate 0.08, and threshold 0.5.

Training subjects:
- ada, ben, cy, dee.

Completely held-out subjects:
- eli, fay.

Training values:
- amber, cobalt, ivory, jade.

Completely held-out values:
- mauve, silver.

Both STORE and OBSERVE event types are represented for every classifier-training subject/value combination.

Classifier evaluation slices:
1. train-domain subjects x train-domain values;
2. held-out subjects x train-domain values;
3. train-domain subjects x held-out values;
4. held-out subjects x held-out values.

## Language integration

Evaluate explicit admission and learned admission on the unchanged UP-LM0G held-out paragraph corpus, both single-paragraph and stream4 modes.

Report overall byte accuracy/perplexity, dependent first-byte accuracy, whole-query-set exact accuracy, admission precision/recall, and max recall entries.

## Interpretation

If the classifier remains correct on entire held-out subject/value families and matches explicit admission in language integration, learned admission has moved beyond recombination memorization. Failure is sealed as an out-of-family generalization boundary, without expanding state or recall budgets.

## Bounds

No explicit event-type bit in the learned arm, no attention, no future oracle, no recurrent-state expansion, no memory beyond 16 entries, no threshold search, no result-informed retry, no live activation, no production authority.
