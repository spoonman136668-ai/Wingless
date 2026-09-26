# Wingless UP-LM1T — fixed-total-mass stability/plasticity map

Status: preregistered scientific five-family language-capacity experiment.

Scientific parent: sealed UP-LM1S 4f572229717f454cbebe0f369c4fe8b5c6f6080e.

## Question

UP-LM1S showed that the fifth lexical family is highly learnable in isolation, but isolated adaptation damages prior-family retention. With total readout learning-rate mass fixed, can redistribution toward the fifth family move the stability/plasticity tradeoff smoothly, or does fifth-family improvement require a disproportionate collapse of prior-family accuracy?

## Frozen starting point

Use the exact five-family pre-adaptation starting point from UP-LM1S:
- 64-D recurrent transition;
- recurrent parameters frozen;
- exact four-family rotating-palindromic training history;
- output alphabet extended for fifth-family bytes;
- matched base / paraphrase / third / fourth / fifth corpora.

## Frozen update budget

All arms:
- 4 adaptation epochs;
- exactly one readout update per family per matched example;
- canonical update order: base, paraphrase, third, fourth, fifth;
- exactly the same examples and update count;
- total learning-rate mass per matched example = 0.40;
- no recurrent updates.

## Preregistered arms

1. equal_mass
   - prior-family learning rate: 0.08 each;
   - fifth-family learning rate: 0.08;
   - total = 4×0.08 + 0.08 = 0.40.
   - exact five-family cyclic control.

2. fifth_1p5_mass
   - prior-family learning rate: 0.07 each;
   - fifth-family learning rate: 0.12;
   - total = 4×0.07 + 0.12 = 0.40.

3. fifth_2x_mass
   - prior-family learning rate: 0.06 each;
   - fifth-family learning rate: 0.16;
   - total = 4×0.06 + 0.16 = 0.40.

These allocations are frozen before execution and are not selected adaptively.

## Evaluation

For every arm:
- held-out top-1 accuracy and perplexity for all five families;
- fifth-family held-out accuracy/perplexity;
- mean prior-family held-out accuracy;
- minimum and mean five-family accuracy;
- five-family spread.

Derived comparisons relative to equal_mass:
- fifth-family accuracy change;
- prior-family mean accuracy change;
- minimum-family accuracy change.

## Interpretation

A graded increase in fifth-family accuracy with a graded prior-family cost establishes a controllable fixed-budget stability/plasticity frontier. Little fifth gain despite reallocation would point toward a deeper representation/history boundary. A sharp prior-family collapse for modest fifth gain would show that the shared readout has little remaining balanced capacity at five families.

No arm is declared a winner and no post-result numeric success threshold is introduced; the measured tradeoff curve is the result.

## Bounds

No extra update count, no increase in total learning-rate mass, no recurrent training, no router, no exact recall, no adaptive weighting, no learning-rate search beyond these preregistered arms, no sixth family, no result-informed retry, no live activation, no production authority.
