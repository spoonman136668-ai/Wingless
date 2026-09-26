# Wingless UP-LM2G — combined locality order robustness

Status: preregistered scientific five-family language integration experiment.

Scientific parent: sealed UP-LM2F d7c5cd06e94cb91764380ae9b4a7b89fe3177a65.

## Question

UP-LM2F showed that de-localizing STORE and OBSERVE together removes the balanced-prior mean penalty across all tested allocations. Does that benefit survive changes to the remaining name-group/update/report order, including a reverse REPORT tail?

## Frozen training

Reuse UP-LM2F unchanged:
- 64-D recurrent transition;
- recurrent parameters frozen;
- exact recall cap 16;
- accepted five-family router;
- five fixed-mass allocation arms;
- canonical_prior and balanced_prior training schedules;
- fifth family always trained last;
- 4 adaptation epochs;
- total learning-rate mass exactly 0.40 per matched example;
- identical update count.

## Frozen evaluation content

Every variant uses the identical held-out tuples and clause content.

In every variant:
1. all four STORE-initial clauses are global;
2. all four OBSERVE clauses are global.

Only the remaining update/report organization changes.

## Tail variants

1. forward_tail: update+report groups 0→1→2→3.
2. rotate_1_tail: 1→2→3→0.
3. rotate_2_tail: 2→3→0→1.
4. rotate_3_tail: 3→0→1→2.
5. reverse_tail: 3→2→1→0.
6. reverse_report_tail:
   - all updates 0→1→2→3;
   - reports 3→2→1→0.

No evaluation variant retrains the model.

## Measurements

For every allocation × training schedule × locality-tail variant × family × stream depth {1,4}:
- top-1 accuracy;
- perplexity;
- dependent first-byte accuracy;
- query/update exactness;
- admission precision/recall;
- event/report routing accuracy;
- max recall entries.

Also report per allocation × training schedule × variant:
- fifth-family mean;
- prior-four-family mean;
- all-family mean;
- minimum family mean.

## Interpretation

If balanced_prior remains neutral/positive in prior mean and improves minimum-family accuracy across these tail variants, combined STORE+OBSERVE de-localization is sequence-order robust. A variant-specific reappearance of the penalty defines the remaining boundary.

No winner, adaptive variant selection, or numeric success threshold is introduced.

## Bounds

No training change, no extra mass, no extra updates, no recurrent training, no router modification, no recall-cap change, no attention, no future oracle, no result-informed retry, no live activation, no production authority.
