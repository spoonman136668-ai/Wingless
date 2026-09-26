# Wingless UP-LM1Y — adaptation-family order robustness

Status: preregistered scientific five-family language adaptation experiment.

Scientific parent: sealed UP-LM1X 4205c6d519cd6f07a346d71d36e4739474a6c345.

## Question

UP-LM1X showed that the five-point fixed-budget stability/plasticity frontier is robust across five evaluation sequence orders. Is that same frontier robust to the order in which lexical families receive their adaptation updates?

## Frozen starting point

Reuse the exact UP-LM1X starting state:
- 64-D recurrent transition;
- recurrent parameters frozen;
- accepted four-family rotating-palindromic prehistory;
- fifth-family alphabet extension;
- accepted five-family router;
- exact recall cap 16;
- no attention;
- no future oracle.

## Frozen allocation arms

All five UP-LM1X allocations remain unchanged:

1. equal_mass: prior 0.0800 each, fifth 0.0800.
2. fifth_1p125_mass: prior 0.0775 each, fifth 0.0900.
3. fifth_1p25_mass: prior 0.0750 each, fifth 0.1000.
4. fifth_1p375_mass: prior 0.0725 each, fifth 0.1100.
5. fifth_1p5_mass: prior 0.0700 each, fifth 0.1200.

Every arm:
- 4 adaptation epochs;
- exactly one update per family per matched example;
- total learning-rate mass exactly 0.40 per matched example;
- identical update count.

## Adaptation-order arms

1. canonical
   - base → paraphrase → third → fourth → fifth.

2. reverse
   - fifth → fourth → third → paraphrase → base.

3. rotating
   - deterministic cyclic start = (epoch + matched-example index) mod 5;
   - each matched example still updates all five families exactly once.

No order is selected adaptively.

## Evaluation

To isolate adaptation order, evaluate the exact block-heldout corpus only, at stream depths 1 and 4, for all five lexical families.

Per allocation × adaptation order × family:
- mean/min integrated top-1 accuracy;
- mean perplexity;
- dependent first-byte accuracy;
- query-set exactness;
- update-count exactness;
- admission precision/recall;
- event/report routing accuracy;
- max recall entries.

Also report per allocation × adaptation order:
- fifth-family integrated mean;
- prior-four-family integrated mean;
- all-family integrated mean;
- minimum family integrated mean.

## Interpretation

If the fixed-mass curves retain the same direction under canonical, reverse, and rotating update order while structural sequence metrics remain exact, the stability/plasticity mechanism is adaptation-order robust. If an order changes the curve materially or causes structural failure, training-order sensitivity becomes the next language boundary.

No allocation or order is declared a winner. No numeric success threshold is introduced.

## Bounds

No extra learning-rate mass, no extra updates, no recurrent training, no router modification, no recall-cap change, no adaptive weighting, no sixth family, no attention, no future oracle, no result-informed retry, no live activation, no production authority.
