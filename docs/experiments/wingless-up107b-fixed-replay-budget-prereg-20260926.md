# Wingless UP-107B — fixed total replay budget

Status: preregistered scientific continual lexical-learning experiment.

Scientific parent: sealed UP-106B 0918523a305c817e156a7326f0888445874d7f8d.

## Question

UP-106B preserved six sequentially acquired verbs with two replay subjects per previously acquired verb, but replay cost grows linearly with vocabulary size. Can a fixed total replay budget per epoch preserve the same six-batch sequence?

## Frozen representation, base lexicon, and acquisition

Identical to UP-106B:
- deterministic 64-D prefix-through-verb encoder;
- three-class linear softmax classifier;
- 20 epochs per acquisition batch;
- learning rate 0.08;
- four grounding subjects for the current verb: ada, ben, cy, dee;
- evaluation subjects: eli, fay;
- base lexicon and six-verb acquisition sequence unchanged;
- base replay remains one fixed subject (ada) for every base verb after every epoch.

## Frozen replay arms for previously acquired verbs

1. full_two_per_prior
   - exact UP-106B control: ada and ben for every previously acquired verb every epoch.

2. fixed_2
   - exactly 2 previous-acquired replay examples total per epoch.

3. fixed_4
   - exactly 4 previous-acquired replay examples total per epoch.

4. fixed_6
   - exactly 6 previous-acquired replay examples total per epoch.

For fixed-budget arms:
- candidate replay pairs are ordered by acquisition order, then subjects ada, ben;
- each epoch starts at a deterministic rotating offset based only on epoch number;
- take exactly the arm's budget, wrapping if needed;
- no accuracy feedback or adaptive selection.

If no previous verb exists, replay count is zero.

## Evaluation

For every arm, at stage 0 and after each of six acquisitions:
- base-lexicon accuracy across all six subjects;
- aggregate accuracy of acquired verbs on eli/fay;
- per-acquired-verb accuracy.

## Interpretation

A fixed-budget arm matching the full replay control would show that replay cost need not grow linearly over this sequence. Failure maps the minimum replay pressure needed as vocabulary accumulates.

## Bounds

No encoder change, no extra state, no adaptive replay, no threshold search, no class weighting, no pretrained semantics, no result-informed retry, no live activation, no production authority.
