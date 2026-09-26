# Wingless UP-104B — sequential lexical acquisition

Status: preregistered scientific continual lexical-learning experiment.

Scientific parent: sealed UP-103B dfb9e9378ed50e3075dc2ee3dd4697a941f66a1b.

## Question

UP-103B showed that one fixed rehearsal anchor per known verb eliminates forgetting while grounding one new STORE/OBSERVE/REPORT verb set. Does the same rule preserve prior lexical knowledge when the three new verbs are acquired sequentially rather than together?

## Frozen base

Use the exact UP-103B / UP-102B prefix-through-verb representation and classifier:
- state dimension 64;
- base seen verbs: stores, keeps, observes, sees, reports, recalls;
- base training unchanged;
- learning rate 0.08;
- 20 epochs per acquisition batch;
- no threshold tuning.

## Sequential acquisition batches

Batch 1:
- holds -> STORE

Batch 2:
- notes -> OBSERVE

Batch 3:
- tells -> REPORT

For each acquired verb:
- four grounding subjects: ada, ben, cy, dee;
- evaluation subjects: eli, fay.

After every epoch of a batch, rehearse exactly one anchor subject (ada) for:
- every base seen verb;
- every previously acquired held-out verb.

The current batch receives its normal four grounding examples and is not separately rehearsed in the same epoch.

## Evaluation

After the base state and after each acquisition batch:
- accuracy on all six base seen verbs across all six subjects;
- accuracy on each acquired verb over eli/fay;
- aggregate acquired-verb accuracy.

## Interpretation

Stable base and earlier-acquired accuracy across all three batches would support bounded sequential lexical acquisition with tiny rehearsal. Degradation identifies where interference reappears as learned vocabulary accumulates.

## Bounds

No encoder change, no extra hidden state, no pretrained semantics, no class weighting, no threshold search, no result-informed retry, no live activation, no production authority.
