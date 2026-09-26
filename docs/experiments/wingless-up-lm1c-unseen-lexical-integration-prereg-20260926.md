# Wingless UP-LM1C — unseen lexical-family integrated stress

Status: preregistered scientific language/memory generalization experiment.

Scientific parent: sealed UP-LM1B 33b5d83b477caa5104cdb67ce2ff3a9e1f9526fc.

## Question

LM1B established that the warm20 jointly optimized byte readout integrates cleanly with frozen routing and exact recall across base/paraphrase structural orderings. Does the integrated procedure generalize when event verbs come from a third lexical family not used in either base or paraphrase byte-model joint training?

## Frozen model

Byte model:
- exact LM1B warm20 procedure;
- 20 base-only epochs;
- 20 base/paraphrase joint-interleaved epochs;
- recurrent transition frozen;
- no further byte-model training.

Router/memory:
- exact LM1B grounded router and exact recall cap 16;
- no router retraining;
- no attention;
- no future oracle.

## Third lexical family

STORE: archives
OBSERVE: notices
REPORT: recounts

The third-family surfaces are never included in byte-model training.

Router grounding:
- fixed 20 epochs at lr 0.08 on ada, ben, cy, dee for archives/notices/recounts;
- same one-subject base/paraphrase rehearsal pattern used in the established lexical-grounding lineage;
- evaluation subjects eli/fay.

## Evaluation

Third-family language is evaluated under:
- block order
- per_name
- paired_names
- stores_then_local_reports
- reverse_report_tail

Each at stream1 and stream4.

Also re-evaluate LM1B base/paraphrase block controls.

## Metrics

- byte top-1 accuracy/perplexity;
- dependent first-byte accuracy;
- whole-query-set exact accuracy;
- admission precision/recall;
- event-routing accuracy;
- report-routing accuracy;
- max recall entries;
- router held-out third-family accuracy.

## Interpretation

Exact routing/memory with degraded byte quality would isolate remaining lexical language-model OOD limits. Failure in routing would indicate that grounded event semantics no longer transfer cleanly when integrated with a third surface family.

## Bounds

No byte-model training on the third family, no recurrent training, no recall-cap increase, no attention, no adaptive grounding, no threshold tuning, no result-informed retry, no live activation, no production authority.
