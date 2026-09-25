# Wingless UP-96B — local event classifier surface generalization

Status: preregistered scientific routing experiment.

Scientific parent: sealed UP-95B `f15ff278753d3af1786b38f9958e41ded3c00fc5`.

## Question

UP-95B perfectly separated the exact phrases `stores` and `observes`. Does the same fixed recurrent event encoder plus linear classifier generalize across multiple lexical surface forms and structurally held-out key/value/verb recombinations?

## Frozen event language

Keys: `k00` through `k31`.
Values: `v00` through `v31`.

Store-class verbs:
- stores
- keeps
- holds

Observe-class verbs:
- observes
- sees
- notes

Surface form:
`<key> <verb> <value>.`

## Frozen split

All six verbs, all 32 keys, and all 32 values appear in training.

Training examples are combinations satisfying:
`(key + 2*value + verb_index) mod 4 != 3`.

Held-out recombination examples satisfy:
`(key + 2*value + verb_index) mod 4 == 3`.

No lexical item is held out; only combinations are.

## Frozen encoder/classifier

- same deterministic 64-D recurrent event encoder as UP-95B;
- binary linear logistic head only;
- zero initialization;
- 20 epochs;
- learning rate 0.08;
- deterministic order;
- threshold exactly 0.5;
- no validation tuning.

## Frozen routing task

Sparse target-change stream:
- target keys: 4, 8, 16;
- total writes: 32, 64, 128, 256;
- store events choose store-class verbs deterministically;
- distractor observation events choose observe-class verbs deterministically;
- final target queries unchanged;
- value vocabulary 32;
- 64 episodes per setting;
- seeds 143M and 144M.

Arms:
1. `explicit_event_class`
2. `learned_surface_classifier`

## Metrics

Classifier:
- train accuracy/precision/recall;
- held-out recombination accuracy/precision/recall;
- per-verb accuracy.

Routing:
- final query accuracy;
- whole-target-set exact accuracy;
- admission precision/recall;
- false-positive admissions.

## Bounds

No explicit event-class bit in learned routing arm, no attention, no trainable recurrence, no threshold search, no capacity increase beyond 16 entries, no result-informed retry, no production authority, no live activation.
