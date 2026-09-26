# Wingless UP-LM1X — mild fixed-mass sequence-order robustness

Status: preregistered scientific five-family language integration experiment.

Scientific parent: sealed UP-LM1W 8564bb800df3d0373efd3d361475037eb731aaf2.

## Question

UP-LM1W found a smooth five-point fixed-budget stability/plasticity frontier in block-heldout evaluation, with exact sequence memory and routing throughout. Does that frontier persist across the accepted sequence-order transformations rather than being specific to block order?

## Frozen starting point and training budget

Reuse the exact UP-LM1W starting state, classifier, five allocation arms, and 4-epoch adaptation:
- total learning-rate mass exactly 0.40 per matched example;
- identical update counts;
- recurrent transition frozen;
- router unchanged;
- exact recall cap 16;
- no adaptive weighting;
- no attention;
- no future oracle.

Allocation arms remain:
- equal_mass: prior 0.0800 each, fifth 0.0800;
- fifth_1p125_mass: prior 0.0775 each, fifth 0.0900;
- fifth_1p25_mass: prior 0.0750 each, fifth 0.1000;
- fifth_1p375_mass: prior 0.0725 each, fifth 0.1100;
- fifth_1p5_mass: prior 0.0700 each, fifth 0.1200.

## Frozen sequence orders

Evaluate all five lexical families under:
1. block;
2. per_name;
3. paired_names;
4. stores_then_local_reports;
5. reverse_report_tail.

For every arm × family × order, evaluate stream depths 1 and 4 with the accepted integrated metric.

## Measurements

Per arm × family × order:
- mean/min top-1 accuracy across stream depths;
- mean perplexity;
- dependent first-byte accuracy;
- query-set exactness;
- update-count exactness;
- admission precision/recall;
- event/report routing accuracy;
- max recall entries.

Per arm × order:
- fifth-family integrated mean;
- prior-four-family integrated mean;
- all-family integrated mean;
- minimum family integrated mean.

## Interpretation

If the mild allocation curves remain smooth with the same direction across sequence orders while structural memory/routing metrics remain exact, the stability/plasticity frontier is order-robust. If specific order transformations introduce discontinuities or structural failures, the next boundary lies in sequence organization rather than allocation alone.

No allocation is selected as a winner and no post-result threshold is introduced.

## Bounds

No extra learning-rate mass, no extra updates, no recurrent training, no router modification, no recall-cap change, no adaptive weighting, no sixth family, no attention, no future oracle, no result-informed retry, no live activation, no production authority.
