# Wingless UP-LM2B — prior-family cyclic position map

Status: preregistered scientific five-family language adaptation experiment.

Scientific parent: sealed UP-LM2A ccdbc414cc87ccf331d7504a370368a7c550cf17.

## Question

UP-LM2A showed that fifth-family terminal privilege survives changes in prior-family order, while prior-family mean and minimum accuracy move with prior order. Does each prior family benefit systematically from occupying later preterminal positions?

## Frozen starting point

Reuse the exact UP-LM2A starting state:
- 64-D recurrent transition;
- recurrent parameters frozen;
- accepted four-family rotating-palindromic prehistory;
- fifth-family alphabet extension;
- accepted five-family router;
- exact recall cap 16;
- no attention;
- no future oracle.

## Frozen allocation arms

All five fixed-mass allocations remain unchanged:
- equal_mass: prior 0.0800 each, fifth 0.0800;
- fifth_1p125_mass: prior 0.0775 each, fifth 0.0900;
- fifth_1p25_mass: prior 0.0750 each, fifth 0.1000;
- fifth_1p375_mass: prior 0.0725 each, fifth 0.1100;
- fifth_1p5_mass: prior 0.0700 each, fifth 0.1200.

Every arm:
- 4 adaptation epochs;
- one update per family per matched example;
- total learning-rate mass exactly 0.40 per matched example;
- identical update count;
- fifth family always position 5.

## Prior-family cyclic rotations

Prior family indices:
0 = base
1 = paraphrase
2 = third
3 = fourth

Four fixed rotations:
- rotation_0: 0 → 1 → 2 → 3 → fifth
- rotation_1: 1 → 2 → 3 → 0 → fifth
- rotation_2: 2 → 3 → 0 → 1 → fifth
- rotation_3: 3 → 0 → 1 → 2 → fifth

Across the four rotations, every prior family occupies preterminal positions 1,2,3,4 exactly once.

## Evaluation

Evaluate exact block-heldout corpora at stream depths 1 and 4 for all five families.

Per allocation × rotation × family:
- integrated mean/min top-1 accuracy;
- mean perplexity;
- dependent first-byte accuracy;
- query-set exactness;
- update-count exactness;
- admission precision/recall;
- event/report routing accuracy;
- max recall entries.

Also report per allocation × rotation:
- fifth-family integrated mean;
- prior-four-family integrated mean;
- all-family integrated mean;
- minimum family integrated mean.

## Interpretation

The exact per-family position response is the result. If a prior family improves when moved later in the four-prior sequence, that establishes a prior-family recency transfer distinct from fifth terminal privilege. If family responses differ idiosyncratically, prior-order effects are interaction-specific rather than a generic recency rule.

No rotation or allocation is selected as a winner. No numeric success threshold is introduced.

## Bounds

No extra learning-rate mass, no extra updates, no recurrent training, no router modification, no recall-cap change, no adaptive weighting, no sixth family, no attention, no future oracle, no result-informed retry, no live activation, no production authority.
