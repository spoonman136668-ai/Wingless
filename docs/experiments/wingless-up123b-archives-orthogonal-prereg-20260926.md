# Wingless UP-123B — archives class-local competitor-orthogonal correction

Status: preregistered scientific representation experiment.

Scientific parent: sealed UP-122B 3381d0138ca59c23d64196124d1d6936523206ec.

## Question

UP-122B localized every alternate-lexicon failure to archives, always as STORE→OBSERVE confusion. Can a frozen class-local competitor-orthogonal correction for archives remove that failure while preserving the already-successful stores correction and all other learned verbs?

## Frozen base system

Keep the exact UP-121B/UP-122B system:
- stores uses the frozen UP-120B competitor-orthogonal correction;
- 64-D prefix-through-verb encoder;
- learning rate 0.08;
- 20 epochs per acquisition batch;
- fixed replay budget 6;
- six class orders;
- current_class_excluded and stage3_anchor2 replay policies;
- original training subjects ada/ben/cy/dee;
- no adaptive geometry.

## Archives correction

Compute once before training from the frozen encoder:
- source = mean archives representation over the original six subjects;
- OBSERVE competitor centroid = mean of notices and spots mean representations;
- REPORT competitor centroid = mean of recounts and retells mean representations;
- project archives source out of the span of those two competitor centroids;
- rescale the residual to the original archives-source norm;
- correction delta = rescaled residual minus original archives source.

Apply this fixed delta only to archives at both training and inference.

## Frozen arms

1. stores_only
   - exact UP-122B representation setup.

2. stores_plus_archives
   - stores correction unchanged;
   - archives gets the frozen correction above.

## Evaluation

Across all six class orders, both replay policies, and every acquisition stage:
- archives accuracy on eli/fay and kia/leo;
- archives correct-class margins;
- all other alternate acquired-verb accuracies;
- stores/keeps original and unseen accuracy;
- base-lexicon original/unseen accuracy;
- aggregate acquired accuracy.

## Interpretation

Closing archives failures without degrading the other metrics would support reusable class-local representation deconfounding. New failures would indicate interacting local corrections are not compositional.

## Bounds

No correction recomputation after results, no adaptive geometry, no replay/training change, no threshold tuning, no state expansion, no result-informed retry, no live activation, no production authority.
