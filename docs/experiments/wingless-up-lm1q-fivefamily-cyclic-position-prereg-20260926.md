# Wingless UP-LM1Q — five-family cyclic position geometry

Status: preregistered scientific multi-family byte-readout optimization experiment.

Scientific parent: sealed UP-LM1P b34310666d604423be4a98a46cfc0a8d49fec26e.

## Question

LM1P showed that five-family cyclic updates produce a stronger worst-family floor and smaller spread than the five-family rotating-palindromic schedule, while the palindromic schedule has a slightly higher mean. Is the cyclic result driven by which family occupies each recency position within the per-example update sequence?

## Frozen starting model and data

Use the exact LM1P five-family starting point before fifth-family adaptation:
- 64-D recurrent transition frozen;
- exact four-family rotating-palindromic training history;
- base / paraphrase / third / fourth / fifth matched training corpora;
- output alphabet already extended for fifth-family bytes;
- no router or exact-memory mechanism in this optimization probe.

## Frozen adaptation budget

All arms:
- 4 adaptation epochs;
- one full update per family per matched example;
- learning rate 0.08;
- identical examples and update count;
- no half steps;
- no adaptive order.

Canonical family order:
0. base
1. paraphrase
2. third
3. fourth
4. fifth

## Arms

Five cyclic rotations of the canonical order:

1. rotate_0: base, paraphrase, third, fourth, fifth
   - exact LM1P five_family_cyclic control.

2. rotate_1: paraphrase, third, fourth, fifth, base

3. rotate_2: third, fourth, fifth, base, paraphrase

4. rotate_3: fourth, fifth, base, paraphrase, third

5. rotate_4: fifth, base, paraphrase, third, fourth

Each arm uses a fixed order for all examples and epochs.

## Evaluation

For each arm:
- held-out top-1 accuracy and perplexity for all five families;
- minimum held-out family accuracy;
- mean held-out family accuracy;
- family spread;
- identity of minimum and maximum families.

## Interpretation

A systematic relationship between final update position and family accuracy would identify recency position as the dominant five-family interference mechanism. Similar results across rotations would instead point toward family-specific representational/readout interference.

## Bounds

No recurrent training, no router use, no exact recall, no extra examples, no adaptive order, no learning-rate search, no result-informed retry, no sixth family, no live activation, no production authority.
