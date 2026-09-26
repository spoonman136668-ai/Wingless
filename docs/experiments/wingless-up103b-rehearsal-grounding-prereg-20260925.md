# Wingless UP-103B — rehearsal-stabilized lexical grounding

Status: preregistered scientific continual lexical-learning experiment.

Scientific parent: sealed UP-102B 1fdb08115705ef9455e3c219fa5810e6f4f1b31e.

## Question

UP-102B reached perfect new-verb transfer with four labeled examples per verb, but seen-verb accuracy fell to 86.11%. Can a tiny fixed rehearsal set preserve old verb classes while retaining the new lexical grounding?

## Frozen base and grounding

Base classifier, encoder, seen verbs, held-out verbs, subjects, 20 epochs, and learning rate 0.08 are identical to UP-102B.

Grounding is fixed at four examples per held-out verb:
- subjects ada, ben, cy, dee.

Evaluation of held-out verbs remains on:
- eli, fay.

## Frozen arms

1. adaptation_only
   - exact UP-102B K=4 adaptation.

2. rehearsal_stabilized
   - same four held-out examples per verb;
   - after each grounding epoch, rehearse exactly one fixed anchor subject (ada) for each of the six seen verbs;
   - no additional examples;
   - no class weighting;
   - no threshold tuning.

Both arms start from the identical base classifier.

## Evaluation

For each arm:
- held-out-verb accuracy;
- held-out STORE/OBSERVE/REPORT precision and recall;
- seen-verb accuracy across all six subjects and six seen verbs.

## Interpretation

If the tiny rehearsal set restores seen-verb accuracy while retaining held-out transfer, the observed B102 forgetting is a continual-learning interference problem rather than a capacity requirement.

## Bounds

No encoder change, no extra state, no pretrained semantics, no class weighting, no threshold search, no result-informed retry, no live activation, no production authority.
