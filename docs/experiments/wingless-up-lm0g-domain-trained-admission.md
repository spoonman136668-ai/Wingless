# Wingless UP-LM0G — natural-domain learned admission

Status: preregistered scientific language/memory integration experiment.

Scientific parents:
- sealed UP-LM0F `c135cb029d072be3a78b8cdc3929094aaea6ebf6`;
- sealed UP-96B `87ca9afbb789733ca50cf39ee4ed730938cac9fd`.

## Question

UP-LM0F failed because a classifier trained on synthetic kXX/vXX events was applied to natural name/value clauses. UP-96B showed the same recurrent encoder/classifier can generalize perfectly within its own event-language domain. If the admission classifier is trained in the actual natural lexical domain, can it replace explicit store admission without losing mutable latest-value language recall?

## Frozen language corpus

Use the exact UP-LM0F paragraph generator:
- names: ada, ben, cy, dee, eli, fay;
- values: amber, cobalt, ivory, jade, mauve, silver;
- four initial store clauses;
- four observe clauses that must not alter memory;
- update counts 0, 1, 2, or 4;
- four report clauses;
- 144 paragraphs;
- structural train/held-out split unchanged from UP-LM0F.

## Frozen recurrent language model

Identical to UP-LM0F:
- 64 float64 recurrent state values;
- deterministic signed-permutation transport;
- deterministic bipolar byte embedding;
- same recurrent update;
- linear softmax output;
- 20 epochs;
- learning rate 0.08;
- no attention;
- no pretrained weights.

## Frozen natural-domain classifier corpus

Clauses:
- `<name> stores <value>.`
- `<name> observes <value>.`

All six names and all six values appear in both train and held-out classifier sets.

Classifier split:
- train iff `(name_index + 2*value_index + event_type_index) mod 3 != 2`;
- held out otherwise.

## Frozen classifier

- same 64-D deterministic recurrent event encoder family;
- binary linear logistic head;
- zero initialization;
- 20 epochs;
- online SGD;
- learning rate 0.08;
- threshold exactly 0.5;
- no validation tuning.

## Frozen arms

1. `explicit_store_admission`
   - upper-bound control;
   - admit only `stores` clauses.

2. `natural_domain_classifier`
   - classify each completed store/observe clause from its bytes only;
   - admit iff classifier probability >= 0.5;
   - no explicit event-type bit.

Both arms:
- exact recall cap exactly 16;
- updates overwrite existing name bindings in place;
- reports query exact recall;
- memory may override only the first reported value byte.

## Evaluation

Classifier:
- train accuracy/precision/recall;
- held-out name/value recombination accuracy/precision/recall.

Language:
- train, held-out, and held-out stream4;
- overall next-byte accuracy/perplexity;
- dependent first-byte accuracy;
- whole-paragraph query-set exact accuracy;
- exact accuracy by update count 0/1/2/4;
- admission precision/recall;
- maximum recall entries used.

## Bounds

No explicit event-type bit in learned arm, no attention, no future oracle, no recurrent-state expansion, no memory beyond 16 entries, no threshold search, no result-informed retry, no live activation, no production authority.
