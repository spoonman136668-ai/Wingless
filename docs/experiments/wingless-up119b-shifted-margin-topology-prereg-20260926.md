# Wingless UP-119B — shifted stores three-class margin topology

Status: preregistered scientific lexical-classifier diagnostic.

Scientific parent: sealed UP-118B 3db20b638d492fb4e1d70f3765e8cf7f82ed31b5.

## Question

UP-118B moved the mean `stores` representation exactly onto the sibling STORE centroid and reduced STORE-vs-OBSERVE failure severity, but did not reduce failure count. Some shifted failures occur even when STORE still beats OBSERVE on mean margin. Are those remaining failures caused by REPORT-class competition rather than OBSERVE competition?

## Frozen training matrix

Repeat the exact UP-118B matrix with no training changes:
- representation arms:
  - baseline
  - stores_centroid_shift
- all six semantic-class acquisition orders;
- replay policies:
  - current_class_excluded
  - stage3_anchor2
- same subjects, verbs, 20 epochs per batch, learning rate 0.08, grounding count, and fixed replay budget 6.

The centroid shift is exactly the sealed UP-118B fixed offset. No geometry is changed or tuned.

## Diagnostics

At stage 0 and after every acquisition, for `stores` over all six subjects report:
- accuracy;
- mean and minimum STORE-minus-OBSERVE margin;
- mean and minimum STORE-minus-REPORT margin;
- mean STORE, OBSERVE, and REPORT logits;
- wrong-destination counts:
  - predicted OBSERVE
  - predicted REPORT.

Also report `keeps` accuracy as the sibling control.

## Interpretation

If shifted failures are dominated by REPORT while STORE-vs-OBSERVE remains positive, UP-118B repaired one geometric axis but exposed a second. If failures still point to OBSERVE, the mean-centroid shift is insufficient within-subject geometry. This diagnostic changes no learning behavior.

## Bounds

No training modification, no new replay, no new representation intervention, no threshold tuning, no adaptive diagnostics, no result-informed retry, no live activation, no production authority.
