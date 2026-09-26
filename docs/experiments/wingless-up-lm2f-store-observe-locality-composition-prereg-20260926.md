# Wingless UP-LM2F — STORE/OBSERVE locality composition

Status: preregistered scientific five-family language integration experiment.

Scientific parent: sealed UP-LM2E a58f81e197cbd9f1928970a30fa59fb29a8e5b15.

## Question

UP-LM2E showed that breaking STORE locality removes the average balanced-prior mean penalty, breaking OBSERVE locality nearly removes it, and breaking REPORT locality alone does not. Are the STORE and OBSERVE effects additive, redundant, or nonlinear when both are de-localized together?

## Frozen training

Reuse UP-LM2E unchanged:
- 64-D recurrent transition;
- recurrent parameters frozen;
- exact recall cap 16;
- accepted five-family router;
- five fixed-mass allocation arms;
- canonical_prior and balanced_prior training schedules;
- fifth family always trained last;
- four adaptation epochs;
- total learning-rate mass exactly 0.40 per matched example;
- identical update count.

## Frozen evaluation content

Every variant uses the exact same held-out names, values, STORE/OBSERVE/update/REPORT clauses, update counts, and query targets. No evaluation retraining.

## Locality variants

1. per_name
   - accepted local store → observe → update → report grouping.

2. stores_global
   - accepted UP-LM2E STORE-only de-localization.

3. observes_global
   - accepted UP-LM2E OBSERVE-only de-localization.

4. stores_observes_global
   - emit all initial STORE clauses across names;
   - then all OBSERVE clauses across names;
   - then for each name emit its update (when present) and REPORT.
   - REPORT remains name-local to its update phase.

Only clause placement changes; lexical content and task semantics are frozen.

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
- fifth-family mean;
- prior-four-family mean;
- all-family mean;
- minimum-family mean.

## Interpretation

Compare balanced_prior minus canonical_prior for each variant. If combined STORE+OBSERVE de-localization improves prior mean beyond both single ablations, effects are additive/synergistic. If it matches one single ablation, the mechanisms are redundant. If it reintroduces a penalty, their interaction is nonlinear.

No winner or numeric threshold is introduced.

## Bounds

No extra mass, no extra updates, no recurrent training, no router modification, no recall-cap change, no attention, no future oracle, no result-informed retry, no live activation, no production authority.
