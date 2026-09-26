# Wingless UP-LM2I — entity-concurrency capacity boundary

Status: preregistered scientific fixed-capacity language integration experiment.

Scientific parent: sealed UP-LM2H 0e79f05ed3bc3ac46f290de9bdfc09ff0a0f0239.

## Question

UP-LM2H showed combined STORE+OBSERVE de-localized language remains structurally exact through a 48-example continuous stream, but recall usage never exceeded six entries. What happens when simultaneous entity concurrency approaches and crosses the frozen 16-entry exact-recall cap?

## Frozen training

Reuse UP-LM2H unchanged:
- 64-D recurrent transition;
- exact recall cap 16;
- accepted five-family router;
- five fixed-mass allocations;
- canonical_prior and balanced_prior schedules;
- fifth family always trained last;
- 4 adaptation epochs;
- total learning-rate mass exactly 0.40 per matched example;
- no recurrent retraining or capacity change.

## Fixed entity-name set

Use 24 deterministic names composed only of characters already present in the accepted six-name vocabulary:
- ada, ben, cy, dee, eli, fay;
- adaben, adacy, adadee, adaeli, adafay;
- benada, bency, bendee, beneli, benfay;
- cyada, cyben, cydee, cyeli, cyfay;
- deeada, deeben, deecy.

Novel concatenated names test entity generalization, so event-routing accuracy is reported explicitly as a confound control.

## Concurrency levels

- 8;
- 12;
- 16;
- 20;
- 24 simultaneous entities.

For each level:
1. all entity STORE clauses occur globally;
2. all entity OBSERVE clauses occur globally;
3. REPORT every stored entity.

No updates occur in this capacity-isolation probe.

## Report-order variants

1. forward_report — report entities in store order.
2. reverse_report — report entities in reverse store order.

Because exact recall is FIFO-16, the preregistered capacity-only expected recall-hit fraction is:
- 1.0 for N <= 16;
- 16/N for N > 16,
provided STORE routing succeeds.

## Families

Evaluate the exact same structure independently with all five accepted lexical families:
- base;
- paraphrase;
- third;
- fourth;
- fifth.

## Measurements

For every allocation × training schedule × family × report order × concurrency:
- byte top-1 accuracy;
- perplexity;
- report dependent first-byte accuracy;
- report-set exact accuracy;
- exact-recall hit rate at REPORT;
- preregistered FIFO capacity-only expected hit rate;
- overall event-routing accuracy;
- STORE/OBSERVE/REPORT routing accuracy;
- maximum recall entries.

## Interpretation

A recall-hit boundary at 16 with routing remaining intact establishes the actual fixed-capacity language-memory boundary. Earlier degradation with routing failure identifies entity-name generalization as the limiting factor instead. Persistence above 16 would falsify the assumed FIFO-16 behavior.

No capacity-dependent training or adaptive naming is allowed.

## Bounds

Exact recall cap remains 16. No capacity increase, no extra training, no router modification, no attention, no future oracle, no adaptive evaluation, no result-informed retry, no live activation, no production authority.
