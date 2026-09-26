# Wingless UP-LM1O — balanced four-family schedule reintegration

Status: preregistered scientific language/memory integration experiment.

Scientific parent: sealed UP-LM1N cf646c9a3da3da131ea245ca757f6ab3236a9014.

## Question

UP-LM1N showed that mirrored and rotating-palindromic four-family update order improve held-out family balance over simple cyclic training at identical nominal update mass. Do those byte-readout improvements remain compatible with the corrected fourth-family semantic router and bounded exact memory?

## Frozen starting model

Exact LM1M/LM1N starting point:
- 64-D recurrent transition frozen;
- 20 base-only epochs;
- 20 base/paraphrase joint-interleaved epochs;
- three-family Strang-split-base adaptation for 4 epochs;
- output alphabet extended for third/fourth surfaces;
- exact recall cap 16.

## Byte adaptation arms

1. cyclic_control
2. mirrored_by_index
3. rotating_palindromic_split

Use the exact UP-LM1N schedules for 4 epochs:
- full step lr 0.08;
- half step lr 0.04 where applicable;
- identical examples and nominal per-family learning-rate mass.

## Frozen semantic system

- exact LM1L corrected fourth-family router;
- states-vs-STORE correction frozen;
- no router retraining;
- no threshold change;
- no attention;
- no future oracle.

## Evaluation

Byte metrics for each arm:
- base held-out;
- paraphrase held-out;
- third held-out;
- fourth held-out.

Integrated routing/memory metrics:
- base block stream1 and stream4;
- paraphrase block stream1 and stream4;
- third block stream1 and stream4;
- fourth family under block, per_name, paired_names, stores_then_local_reports, reverse_report_tail;
- each fourth order at stream1 and stream4.

Metrics:
- byte top-1 accuracy/perplexity;
- dependent first-byte accuracy;
- whole-query-set exact accuracy;
- admission precision/recall;
- event-routing accuracy;
- report-routing accuracy;
- max recall entries.

## Interpretation

If balanced byte schedules retain exact semantic routing and memory behavior, the optimization improvement is integration-safe and becomes the new four-family language baseline.

## Bounds

No recurrent training, no router retraining, no recall-cap increase, no extra examples, no adaptive ordering, no attention, no result-informed retry, no live activation, no production authority.
