# Wingless UP-LM2H — combined locality stream-depth stress

Status: preregistered scientific five-family language integration experiment.

Scientific parent: sealed UP-LM2G be861515e9020fe1be8d8e3a1193a77582e2621d.

## Question

UP-LM2G showed that combined STORE+OBSERVE de-localization remains beneficial across multiple tail orders at stream depths 1 and 4. Does that integration rule survive substantially longer continuous language streams under the same exact recall cap?

## Frozen training

Reuse UP-LM2G unchanged:
- 64-D recurrent transition;
- recurrent parameters frozen;
- exact recall cap 16;
- accepted five-family router;
- five fixed-mass allocations;
- canonical_prior and balanced_prior schedules;
- fifth family always trained last;
- 4 adaptation epochs;
- total learning-rate mass exactly 0.40 per matched example;
- no training change.

## Frozen evaluation organizations

Use two representative combined-locality organizations already accepted in UP-LM2G:

1. forward_tail
   - all STORE clauses global;
   - all OBSERVE clauses global;
   - update+report groups forward.

2. reverse_report_tail
   - all STORE clauses global;
   - all OBSERVE clauses global;
   - all updates forward;
   - REPORT tail reversed.

Content is identical across organizations.

## Stream depths

Evaluate each held-out family at:
- 1;
- 4;
- 8;
- 16;
- 48 examples per continuous stream.

48 spans the complete frozen held-out set without resetting recurrent/recall state between examples.

The exact recall cap remains 16; no capacity is added for longer streams.

## Measurements

For every allocation × training schedule × organization × family × stream depth:
- top-1 accuracy;
- perplexity;
- dependent first-byte accuracy;
- query/update exactness;
- admission precision/recall;
- event/report routing accuracy;
- maximum recall entries.

Per allocation × schedule × organization × stream depth:
- fifth-family mean;
- prior-four-family mean;
- all-family mean;
- minimum family mean;
- structural exactness;
- maximum recall entries.

## Interpretation

Persistence of the balanced-prior integration benefit and exact structural behavior through longer depths supports a genuine fixed-capacity language integration mechanism rather than a short-stream artifact. Depth-specific degradation defines the next boundary.

No winner or numeric success threshold is introduced.

## Bounds

No capacity increase, no extra training, no recurrent training, no router modification, no attention, no future oracle, no adaptive evaluation, no result-informed retry, no live activation, no production authority.
