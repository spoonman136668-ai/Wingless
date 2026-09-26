# Wingless UP-121B — competitor-orthogonal lexical generalization

Status: preregistered scientific representation-generalization experiment.

Scientific parent: sealed UP-120B 17bb8323a0df4583e8b41807e19b3d473b3a1043.

## Question

UP-120B eliminated every observed stores failure by projecting the stores mean representation away from the OBSERVE/REPORT competitor span. Does that fixed correction generalize beyond the exact six training/evaluation subjects and acquisition verb set used to derive the failure matrix?

## Frozen representation intervention

Use the exact UP-120B competitor-orthogonal stores correction:
- correction is computed only from the original frozen encoder family means used by UP-120B;
- correction coefficients are frozen before this experiment;
- only the surface "stores" is corrected;
- no other verb encoding changes;
- no adaptive geometry.

## Frozen classifier/training

- 64-D prefix-through-verb encoder;
- same three-class linear softmax head;
- 20 epochs per acquisition batch;
- learning rate 0.08;
- fixed replay budget 6;
- current_class_excluded and stage3_anchor2 replay policies remain unchanged.

## Generalization grid

Subjects:
- original six: ada, ben, cy, dee, eli, fay
- unseen six: gia, hal, ivy, jon, kia, leo

Acquisition lexical set:
- original: holds / notes / tells / saves / watches / remembers
- alternate: keepsafe / notices / recounts / archives / spots / retells

Class orders:
- all six permutations of STORE / OBSERVE / REPORT.

Arms:
1. baseline stores encoding
2. frozen competitor-orthogonal stores encoding

Train on original four grounding subjects for each lexical set.
Evaluate base stores/keeps accuracy on:
- original six subjects
- unseen six subjects.

Evaluate acquired verbs on held-out subjects from the corresponding subject family.

## Metrics

For every arm x lexical set x class order x replay policy x stage:
- stores accuracy original subjects;
- stores accuracy unseen subjects;
- keeps accuracy original/unseen;
- base-lexicon accuracy original/unseen;
- acquired aggregate accuracy;
- STORE-vs-OBSERVE and STORE-vs-REPORT margins for stores.

## Interpretation

If the frozen correction remains robust on unseen subjects and alternate lexical surfaces, UP-120B reflects a reusable representation-level deconfounding mechanism rather than matrix-specific repair.

## Bounds

No correction recomputation, no adaptive geometry, no replay change, no state expansion, no threshold tuning, no result-informed retry, no live activation, no production authority.
