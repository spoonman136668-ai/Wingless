# Wingless UP-LM1J — fourth lexical-family boundary

Status: preregistered scientific language/memory breadth experiment.

Scientific parent: sealed UP-LM1I 7eedf2833867fa1b154a5113c07f08cd77ae170b.

## Question

LM1I established an integrated three-family adaptation procedure with exact routing/memory and improved worst-family byte quality. Where is the next breadth boundary when a fourth event-verb family is introduced without byte-model training?

## Frozen integrated model

Byte model:
- exact LM1I strang_split_base procedure;
- 64-D recurrent state;
- recurrent transition frozen;
- 20 base pretrain epochs;
- 20 base/paraphrase joint-interleaved epochs;
- 4 three-family strang_split_base adaptation epochs;
- no fourth-family byte training.

Router/memory:
- exact recall cap 16;
- no attention;
- no future oracle;
- existing base/paraphrase/third routing frozen.

## Fourth lexical family

STORE: retains
OBSERVE: inspects
REPORT: states

The fourth-family surfaces are never included in byte-model training.

Router grounding:
- linear 3-class grounded event router;
- 20 epochs at lr 0.08;
- fourth-family grounding subjects: ada, ben, cy, dee;
- one-subject rehearsal of each prior lexical family per epoch;
- held-out fourth-family subjects: eli, fay.

No explicit event class at inference.

## Evaluation

Controls:
- base block
- paraphrase block
- third-family block

Fourth-family orders:
- block
- per_name
- paired_names
- stores_then_local_reports
- reverse_report_tail

Each evaluated at stream1 and stream4.

## Metrics

- fourth-family router held-out accuracy;
- byte top-1 accuracy/perplexity;
- dependent first-byte accuracy;
- whole-query-set exact accuracy;
- admission precision/recall;
- event-routing accuracy;
- report-routing accuracy;
- max recall entries.

## Interpretation

Exact routing/memory with degraded byte quality would identify the next lexical-readout OOD boundary. Routing failure would instead show semantic grounding capacity is becoming the breadth limit.

## Bounds

No fourth-family byte training, no recurrent training, no recall-cap increase, no attention, no adaptive grounding, no threshold tuning, no result-informed retry, no live activation, no production authority.
