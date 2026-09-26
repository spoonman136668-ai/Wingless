# Wingless UP-LM1Z — fifth-family update position

Status: preregistered scientific five-family language adaptation experiment.

Scientific parent: sealed UP-LM1Y feeca16b52f60cb0dfdac58de46ac02ff061b0c5.

## Question

UP-LM1Y showed that fixed-budget stability/plasticity survives adaptation-family order, but the operating point moves with update recency: canonical order favors the newest fifth family, reverse order preserves more prior-family accuracy, and rotating order lies between. Is the fifth-family position itself sufficient to generate that recency effect?

## Frozen starting point

Reuse the exact UP-LM1Y starting state:
- 64-D recurrent transition;
- recurrent parameters frozen;
- accepted four-family rotating-palindromic prehistory;
- fifth-family alphabet extension;
- accepted five-family router;
- exact recall cap 16;
- no attention;
- no future oracle.

## Frozen allocation arms

All five accepted fixed-mass allocations remain unchanged:

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

## Fifth-family position arms

The relative order of the four prior families is frozen:
base → paraphrase → third → fourth.

Only the fifth-family update position changes:
- fifth_pos_1: fifth → base → paraphrase → third → fourth
- fifth_pos_2: base → fifth → paraphrase → third → fourth
- fifth_pos_3: base → paraphrase → fifth → third → fourth
- fifth_pos_4: base → paraphrase → third → fifth → fourth
- fifth_pos_5: base → paraphrase → third → fourth → fifth

No order is selected adaptively.

## Evaluation

Evaluate the exact block-heldout corpus only, at stream depths 1 and 4, for all five lexical families.

Per allocation × fifth position × family:
- mean/min integrated top-1 accuracy;
- mean perplexity;
- dependent first-byte accuracy;
- query-set exactness;
- update-count exactness;
- admission precision/recall;
- event/report routing accuracy;
- max recall entries.

Also report per allocation × fifth position:
- fifth-family integrated mean;
- prior-four-family integrated mean;
- all-family integrated mean;
- minimum family integrated mean.

## Interpretation

The exact five-position curve is the result. A systematic recency gradient would show that fifth-family update position is sufficient to shift the fixed-budget stability/plasticity operating point. A flat or irregular curve would indicate that the broader order patterns in UP-LM1Y depended on interactions among prior-family order as well.

No position or allocation is declared a winner. No numeric success threshold is introduced.

## Bounds

No extra learning-rate mass, no extra updates, no recurrent training, no router modification, no recall-cap change, no adaptive weighting, no sixth family, no attention, no future oracle, no result-informed retry, no live activation, no production authority.
