# Wingless UP-97B — three-way natural event recognition

Status: preregistered scientific event-routing experiment.

Scientific parent: sealed UP-96B `87ca9afbb789733ca50cf39ee4ed730938cac9fd`.

## Question

UP-96B closed binary store-vs-observe classification across multiple verbs and held-out recombinations. Can the same fixed recurrent event encoder support a three-way local classifier that distinguishes STORE, OBSERVE, and REPORT events in natural lexical clauses?

## Frozen lexical domain

Names:
- ada
- ben
- cy
- dee
- eli
- fay

Values:
- amber
- cobalt
- ivory
- jade
- mauve
- silver

STORE verbs:
- stores
- keeps
- holds

OBSERVE verbs:
- observes
- sees
- notes

REPORT verbs:
- reports
- recalls
- tells

Surface forms:
- store/observe: `<name> <verb> <value>.`
- report: `<name> <verb> <value>.`

The value is present in all three classes during classifier qualification; the classifier is only deciding local event type.

## Frozen split

All names, values, and verbs appear in training.

Training iff:
`(name_index + 2*value_index + verb_global_index) mod 4 != 3`.

Held-out recombination iff the expression equals 3.

## Frozen encoder/classifier

- deterministic 64-D recurrent event encoder;
- no trainable recurrence;
- three-class linear softmax head;
- zero initialization;
- 20 epochs;
- online SGD;
- learning rate 0.08;
- deterministic lexical order;
- no threshold tuning.

## Evaluation

- train three-way accuracy;
- held-out recombination three-way accuracy;
- per-class precision/recall;
- per-verb held-out accuracy;
- confusion counts STORE/OBSERVE/REPORT.

## Downstream probe

On deterministic event sequences:
- STORE writes a binding;
- OBSERVE does not alter memory;
- REPORT queries the binding;
- exact recall cap 16;
- compare explicit event class vs learned three-way event class;
- target set sizes 4, 8, 16;
- distractor counts 32, 64, 128, 256;
- seeds 147M and 148M.

Report final query accuracy and whole-query-set exact accuracy.

## Bounds

No explicit event class in learned probe, no attention, no trainable recurrence, no capacity increase beyond 16 entries, no validation tuning, no result-informed retry, no production authority, no live activation.
