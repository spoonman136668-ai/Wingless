# Wingless UP-LM1B — warm20 joint-readout reintegration

Status: preregistered scientific language/memory integration experiment.

Scientific parent: sealed/corrected UP-LM1A 7764d7d75d1d4e9d2c200c68a0de26d425647d70.

## Question

LM1A showed that a base-pretrained readout followed by 20 joint interleaved epochs reaches strong held-out byte accuracy on both base and paraphrase distributions. Does that jointly optimized byte model remain compatible with the learned event router and exact-recall mechanism across block and reordered language streams?

## Frozen byte model training

- deterministic 64-D recurrent transition, never trained;
- zero-initialized readout;
- 20 base-only epochs at learning rate 0.08;
- then 20 joint-interleaved epochs:
  for each matched index, base example then paraphrase example;
- exact LM1A warm_20 procedure;
- no further byte-model updates during evaluation.

## Frozen router and memory

- use the grounded three-way classifier from LM0N unchanged;
- no router retraining during this experiment;
- exact recall cap 16;
- same exact-memory logic used by LM0N/LM0K;
- no attention;
- no future oracle.

## Evaluation corpus

Base distribution:
- original held-out block corpus at stream1 and stream4.

Paraphrase distribution:
- paraphrase block held-out corpus at stream1 and stream4;
- four structural order families from LM0N:
  1. per_name
  2. paired_names
  3. stores_then_local_reports
  4. reverse_report_tail
- each order family at stream1 and stream4.

## Metrics

Byte-only:
- base held-out top-1 accuracy/perplexity;
- paraphrase held-out top-1 accuracy/perplexity.

Integrated routing/memory for every evaluation cell:
- top-1 byte accuracy;
- perplexity;
- dependent first-byte accuracy;
- whole-query-set exact accuracy;
- admission precision/recall;
- event-routing accuracy;
- report-routing accuracy;
- maximum recall entries.

## Interpretation

If the warm20 jointly optimized readout preserves exact routing/memory behavior while improving generic byte quality on both distributions, the language lane has a viable integrated adaptation procedure rather than an isolated readout diagnostic.

## Bounds

No recurrent training, no router retraining, no recall-cap increase, no attention, no threshold tuning, no adaptive schedule, no result-informed retry, no live activation, no production authority.
