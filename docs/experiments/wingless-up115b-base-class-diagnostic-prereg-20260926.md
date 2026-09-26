# Wingless UP-115B — base-class interference localization

Status: preregistered scientific diagnostic experiment.

Scientific parent: sealed UP-114B 3c64686896dd22342a1d5f230c1292f6f9c91755.

## Question

UP-114B showed that acquired-verb retention is perfect across all six semantic-class acquisition orders, while base-lexicon errors move among early stages and the stage-3 anchor can help, hurt, or be irrelevant. Which base semantic class and base verbs are being displaced at each order/stage?

## Frozen training

Repeat the exact UP-114B training matrix without any learning change:
- all six class orders: SOR, SRO, OSR, ORS, RSO, ROS;
- both policies: current_class_excluded and stage3_anchor2;
- same deterministic 64-D prefix-through-verb encoder;
- same base and acquired verbs;
- four grounding subjects per current verb;
- 20 epochs per acquisition batch;
- learning rate 0.08;
- fixed six-example additional replay budget.

## Additional diagnostic measurements

At stage 0 and after every acquisition, report:
- total base-lexicon accuracy;
- STORE base-class accuracy across stores + keeps and all six subjects;
- OBSERVE base-class accuracy across observes + sees and all six subjects;
- REPORT base-class accuracy across reports + recalls and all six subjects;
- per-base-verb accuracy for all six base verbs;
- aggregate acquired-verb accuracy and per-acquired-verb accuracy exactly as in UP-114B.

No diagnostic metric affects training.

## Interpretation

The experiment localizes the order-dependent base interference before another replay intervention is introduced. A consistent displaced class relative to the current acquisition class would support a class-competition mechanism; heterogeneous failures would argue for lexical or trajectory-specific interference.

## Bounds

No training-rule change, no replay-budget change, no adaptive replay, no threshold tuning, no class weighting, no extra state, no result-informed retry, no live activation, no production authority.
