# Wingless UP-LM2E — per-name locality ablation

Status: preregistered scientific five-family language integration experiment.

Scientific parent: sealed UP-LM2D 12374b0c3bfc3227b9e00b335ca4c5ca6af307fb.

## Question

UP-LM2D showed that balanced-prior training pays a small prior-family mean-accuracy cost under every permutation of per-name-local groups. Which component of per-name locality creates that interaction: local initial STORE placement, local OBSERVE placement, or immediate local REPORT placement?

## Frozen training

Reuse UP-LM2D unchanged:
- 64-D recurrent transition;
- recurrent parameters frozen;
- exact recall cap 16;
- accepted five-family router;
- fifth family always trained last;
- five fixed-mass allocations;
- canonical_prior and balanced_prior schedules as fixed interaction controls;
- 4 adaptation epochs;
- total learning-rate mass exactly 0.40 per matched example;
- identical update count.

No model is retrained on an evaluation variant.

## Frozen held-out content

Every variant uses the exact accepted held-out tuple set for every lexical family:
- same four names;
- same initial stores;
- same observations;
- same optional updates;
- same reports;
- same values;
- same update count;
- same query targets.

Only clause placement across names changes.

## Locality ablations

1. per_name
   - for each name: STORE → OBSERVE → UPDATE(if present) → REPORT.

2. stores_global
   - all four initial STORE clauses first;
   - then per name: OBSERVE → UPDATE → REPORT.
   - only initial STORE locality is removed.

3. observes_global
   - all four OBSERVE clauses first;
   - then per name: STORE → UPDATE → REPORT.
   - only OBSERVE locality is removed; observations do not modify exact recall state.

4. reports_global
   - per name: STORE → OBSERVE → UPDATE;
   - all four REPORT queries at the end.
   - only immediate REPORT locality is removed.

## Measurements

For every allocation × training schedule × locality variant × family × stream depth {1,4}:
- top-1 accuracy;
- perplexity;
- dependent first-byte accuracy;
- query exactness;
- update-count exactness;
- admission precision/recall;
- event/report routing accuracy;
- max recall entries.

Also report per allocation × schedule × variant:
- fifth-family integrated mean;
- prior-four-family integrated mean;
- all-family integrated mean;
- minimum family integrated mean.

## Interpretation

The balanced-minus-canonical delta for each ablation is the result. If removing one locality component eliminates or reverses the prior-family mean penalty while the other ablations retain it, that component explains the interaction. If all remain negative, the penalty belongs to deeper name-local composition rather than one clause boundary.

No post-result feature selection, threshold, or retraining is allowed.

## Bounds

No extra mass, no extra updates, no recurrent training, no router modification, no recall-cap change, no sixth family, no attention, no future oracle, no result-informed retry, no live activation, no production authority.
