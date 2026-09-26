# Wingless UP-108B — stratified fixed-budget replay

Status: preregistered scientific continual lexical-learning experiment.

Scientific parent: sealed UP-107B b368ef8385e231d69ac2379e8f8b9a39d31bd4e5.

## Question

UP-107B showed that six replay examples per epoch can finish with all six acquired verbs intact, but the pair-major rotating schedule has intermediate losses. Is replay allocation coverage, rather than total replay count, the remaining limit?

## Frozen representation and learning

Identical to UP-107B:
- deterministic 64-D prefix-through-verb encoder;
- same base lexicon;
- same six-verb acquisition sequence;
- four grounding subjects for the current verb;
- 20 epochs per acquisition batch;
- learning rate 0.08;
- base replay remains one fixed subject (ada) for every base verb after every epoch.

## Frozen arms

1. full_two_per_prior
   - exact growing replay control from UP-106B.

2. pair_major_fixed6
   - exact UP-107B fixed_6 schedule.
   - candidate pairs are previous verbs in acquisition order, each with ada then ben.
   - deterministic rotating offset by epoch.
   - exactly 6 previous-acquired replay examples per epoch when previous verbs exist.

3. verb_stratified_fixed6
   - exactly 6 previous-acquired replay examples per epoch when previous verbs exist.
   - first allocate one replay example to each previous verb in acquisition order, up to the six-example budget.
   - subject for that first coverage pass alternates deterministically between ada and ben by epoch plus verb index.
   - if budget remains, allocate second examples across previous verbs in acquisition order using the other subject.
   - wrap deterministically only if fewer than three previous verbs require repeated examples to reach six.
   - no performance feedback or adaptive selection.

## Evaluation

For every arm, at stage 0 and after each of six acquisitions:
- base-lexicon accuracy;
- aggregate acquired-verb accuracy;
- per-acquired-verb accuracy.

## Interpretation

If stratified_fixed6 removes the intermediate losses seen in pair_major_fixed6, replay coverage—not replay quantity—is the mechanism. Failure would indicate six total examples are intrinsically insufficient for stable retention across every stage.

## Bounds

No encoder change, no extra state, no adaptive replay, no threshold search, no class weighting, no pretrained semantics, no result-informed retry, no live activation, no production authority.
