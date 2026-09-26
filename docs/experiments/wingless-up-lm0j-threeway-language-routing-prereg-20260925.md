# Wingless UP-LM0J — learned three-way event routing in the language loop

Status: preregistered scientific language/memory integration experiment.

Scientific parent: sealed UP-LM0I `792ea898d6c33ca42b6b4031db0d9d4a82044cd9`.

Supporting evidence: sealed UP-99B established perfect prefix-through-verb three-way separation on its lexical routing domain.

## Question

UP-LM0I removed explicit STORE-vs-OBSERVE admission by making the decision at the verb boundary, but REPORT/query routing is still hand-coded. Can a learned prefix-timed three-way classifier control STORE, OBSERVE, and REPORT behavior inside the byte-level language loop without explicit event-type bits?

## Frozen language model and memory

Identical to UP-LM0I:
- deterministic 64-D recurrent language state;
- 512 recurrent-state bytes;
- 20 language-model epochs;
- learning rate 0.08;
- no attention;
- exact recall cap 16;
- unchanged UP-LM0F paragraph corpus and language split.

## Frozen three-way classifier

Classifier input is only `<name> <verb>`.

Classes:
- STORE: `stores`;
- OBSERVE: `observes`;
- REPORT: `reports`.

Training subjects:
- ada, ben, cy, dee.

Held-out subjects:
- eli, fay.

Classifier:
- deterministic 64-D event encoder;
- three-class linear softmax head;
- zero initialization;
- 20 epochs;
- online SGD;
- learning rate 0.08;
- no value bytes;
- no explicit class bit.

## Language routing

At the second-field boundary, predict one of STORE / OBSERVE / REPORT.

- STORE: at clause completion, write generic first-field -> third-field binding.
- OBSERVE: no exact-memory write.
- REPORT: query memory for the first field and allow exact memory to override only the first byte of the reported value.

The clause parser may identify field boundaries and copy lexical fields, but may not use the verb string to choose the semantic action in the learned arm.

Control arm uses explicit event class for the same actions.

## Metrics

Classifier:
- train-subject and held-out-subject overall accuracy;
- per-class precision/recall.

Language:
- explicit vs learned three-way routing;
- held-out single paragraph and stream4;
- overall byte accuracy/perplexity;
- dependent first-byte accuracy;
- whole-query-set exact accuracy;
- STORE-action precision/recall;
- max recall entries.

## Bounds

No explicit event-type bit in learned inference, no verb-string semantic branching in learned routing, no value bytes in classifier input, no attention, no future oracle, no state expansion, no recall-cap increase, no threshold search, no result-informed retry, no live activation, no production authority.
