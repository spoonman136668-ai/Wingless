# Wingless UP-LM0I — verb-anchored learned admission

Status: preregistered scientific language/memory integration experiment.

Scientific parent: sealed UP-LM0H `d9b8b8db0042cd4b55eb05afc1fb4b52d1ef5c67`.

## Question

UP-LM0H produced perfect admission precision but poor recall when whole value families were unseen, suggesting later value bytes overwrite or interfere with event-type information. Can the same learned classifier generalize if the admission decision is made from the clause prefix through the verb, before any value bytes arrive?

## Frozen language model and memory

Identical to UP-LM0H:
- deterministic 64-D recurrent language state;
- 512 recurrent-state bytes;
- 20 language-model epochs;
- learning rate 0.08;
- exact recall cap 16;
- no attention;
- same UP-LM0F paragraph corpus and language train/held-out split.

## Frozen admission classifier

- deterministic 64-D event encoder;
- binary linear logistic head;
- zero initialization;
- 20 epochs;
- learning rate 0.08;
- threshold 0.5.

Classifier training subjects:
- ada, ben, cy, dee.

Held-out subjects:
- eli, fay.

Classifier input is only:
- `<name> stores`
- `<name> observes`

No value bytes or punctuation are supplied to the classifier.

## Inference timing

During paragraph processing:
- once the second lexical field (verb) is complete, classify the current `<name> <verb>` prefix;
- carry that one-bit learned admission decision until the clause terminator;
- if the clause is a store/observe clause, use the carried decision;
- report clauses never write memory;
- no explicit store/observe label is supplied to the classifier.

## Evaluation

Classifier:
- training-subject accuracy/precision/recall;
- held-out-subject accuracy/precision/recall.

Language:
- explicit admission vs verb-anchored learned admission;
- held-out single-paragraph and stream4;
- byte accuracy/perplexity;
- dependent first-byte accuracy;
- whole-query-set exact accuracy;
- admission precision/recall;
- maximum recall entries.

## Interpretation

Recovery would directly support suffix interference / readout timing as the LM0H failure mechanism. Failure would move the next experiment upstream to the event encoder itself.

## Bounds

No explicit event-type bit in learned inference, no value bytes in classifier input, no attention, no future oracle, no state expansion, no recall-cap increase, no threshold search, no result-informed retry, no live activation, no production authority.
