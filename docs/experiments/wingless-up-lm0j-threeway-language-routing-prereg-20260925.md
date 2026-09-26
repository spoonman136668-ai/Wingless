# Wingless UP-LM0J — learned three-way language routing

Status: preregistered scientific language/memory integration experiment.

Scientific parent: sealed UP-LM0I `792ea898d6c33ca42b6b4031db0d9d4a82044cd9`.

## Question

UP-LM0I fully restored binary STORE-vs-OBSERVE admission by classifying at the verb boundary, but REPORT routing is still hand-coded. Can one learned prefix-timed three-way classifier control STORE, OBSERVE, and REPORT behavior inside the language loop without explicit event-type labels?

## Frozen language model and memory

Identical to UP-LM0I:
- deterministic 64-D recurrent language state;
- 512 recurrent-state bytes;
- 20 language-model epochs;
- learning rate 0.08;
- exact recall cap 16;
- no attention;
- same paragraph corpus and structural train/held-out split.

## Frozen three-way event classifier

Inputs:
- `<name> stores`;
- `<name> observes`;
- `<name> reports`.

Training subjects:
- ada, ben, cy, dee.

Held-out subjects:
- eli, fay.

Classifier:
- same deterministic 64-D event encoder;
- three-class linear softmax head;
- zero initialization;
- 20 epochs;
- online SGD;
- learning rate 0.08;
- prediction at completion of the verb token;
- no value bytes in classifier input.

## Frozen arms

1. `explicit_event_routing`
   - explicit STORE/OBSERVE/REPORT control.

2. `learned_prefix_threeway`
   - learned class determines all three behaviors;
   - STORE: write name/value at clause completion;
   - OBSERVE: no exact-memory write;
   - REPORT: query name binding for first dependent value byte.

Generic lexical parsing may recover name/value fields, but may not supply the event class.

## Evaluation

Classifier:
- train-subject and held-out-subject three-way accuracy;
- per-class precision/recall.

Language:
- held-out single-paragraph and stream4;
- byte accuracy/perplexity;
- dependent first-byte accuracy;
- whole-query-set exact accuracy;
- admission precision/recall;
- report-query success;
- maximum recall entries.

## Interpretation

Matching the explicit control would remove the remaining hand-coded event-type routing from this LM0 lineage. Failure would identify a three-way integration boundary despite successful isolated prefix classification.

## Bounds

No explicit event type in learned inference, no value bytes at classifier inference, no attention, no future oracle, no state expansion, no recall-cap increase, no threshold search, no result-informed retry, no live activation, no production authority.
