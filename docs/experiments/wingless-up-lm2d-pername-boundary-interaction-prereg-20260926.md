# Wingless UP-LM2D — per-name boundary interaction

Status: preregistered scientific five-family language integration experiment.

Scientific parent: sealed UP-LM2C b1b4dec076ca63765b31b12e99a53ffbf5f89921.

## Question

UP-LM2C showed that position-balanced prior training improves weakest-family accuracy overall, but the per_name evaluation order uniquely reduces prior-family mean accuracy. Is that penalty a generic consequence of per-name locality, or does it depend on which name-group occupies each within-example position?

## Frozen training

Reuse UP-LM2C unchanged:
- 64-D recurrent transition;
- recurrent parameters frozen;
- exact recall cap 16;
- accepted five-family router;
- fifth family always trained last;
- five fixed-mass allocation arms;
- two training schedules: canonical_prior and balanced_prior;
- 4 adaptation epochs;
- total learning-rate mass exactly 0.40 per matched example;
- identical update count.

## Frozen evaluation content

Use the exact held-out tuple set for every lexical family. Every example contains the same four names, initial stores, observations, updates, reports, values, and update count as accepted per_name evaluation.

Only the order of the four complete name-local groups changes.

## Per-name group-order variants

1. forward: 0 → 1 → 2 → 3.
2. rotate_1: 1 → 2 → 3 → 0.
3. rotate_2: 2 → 3 → 0 → 1.
4. rotate_3: 3 → 0 → 1 → 2.
5. reverse: 3 → 2 → 1 → 0.

Within each name-group, the accepted per_name clause order remains:
store → observe → update(if present) → report.

## Measurements

For every allocation × training schedule × family × group-order variant × stream depth {1,4}:
- top-1 accuracy;
- perplexity;
- dependent first-byte accuracy;
- query exactness;
- update-count exactness;
- admission precision/recall;
- event/report routing accuracy;
- max recall entries.

Also report per allocation × training schedule × group order:
- fifth-family integrated mean;
- prior-four-family integrated mean;
- all-family integrated mean;
- minimum family integrated mean.

## Interpretation

If balanced_prior remains worse in prior-family mean across all five name-group permutations, the LM2C penalty is a generic per-name-locality interaction. If the sign or magnitude changes substantially with group position, the mechanism depends on name-boundary placement/recency.

No variant is selected adaptively and no numeric success threshold is introduced.

## Bounds

No retraining on evaluation variants, no extra mass, no extra updates, no recurrent training, no router modification, no recall-cap change, no sixth family, no attention, no future oracle, no result-informed retry, no live activation, no production authority.
