# Wingless UP-126B — learned STORE gate topology diagnostic

Status: preregistered scientific representation-routing diagnostic.

Scientific parent: sealed UP-125B 2e86cdcd20f88a1d6480ff647349b3baf90a8758.

## Question

UP-125B produced perfect STORE-gate precision but only 0.20 held-out and 0.267 unseen recall. Is the low recall broad across STORE surfaces, or localized to specific lexical surfaces / subject shifts?

## Frozen gate

Use the exact UP-125B gate:
- raw 64-D prefix-through-verb encoder;
- zero-initialized logistic gate;
- training surfaces and subjects unchanged;
- 20 epochs;
- learning rate 0.08;
- threshold 0.5;
- no retraining, threshold search, or class rebalancing.

## Evaluation surface grid

STORE:
- stores
- keeps
- holds
- saves
- archives

Non-STORE:
- observes
- sees
- notes
- watches
- notices
- reports
- recalls
- tells
- remembers
- recounts

Subject splits:
1. train: ada, ben, cy, dee
2. heldout: eli, fay
3. unseen: gia, hal, ivy, jon, kia, leo

## Metrics

For every split x surface:
- mean gate probability;
- minimum gate probability;
- maximum gate probability;
- positive prediction rate;
- target class.

For every split:
- accuracy;
- precision;
- recall.

Also report:
- STORE surface mean-probability ordering;
- minimum STORE margin to threshold;
- maximum non-STORE probability.

## Interpretation

If only a subset of STORE surfaces fall below threshold, the next experiment should target their representation geometry. If all STORE surfaces are uniformly depressed, the next experiment should test a preregistered class-balance/training-objective intervention rather than threshold tuning.

## Bounds

Diagnostic only. No gate retraining, no threshold adjustment, no projector change, no downstream classifier training, no adaptive geometry, no result-informed retry, no live activation, no production authority.
