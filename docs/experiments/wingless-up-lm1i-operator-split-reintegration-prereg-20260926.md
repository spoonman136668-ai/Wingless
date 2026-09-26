# Wingless UP-LM1I — operator-split reintegration

Status: preregistered scientific language/memory integration experiment.

Scientific parent: sealed UP-LM1H 5fd53050638d4d636132392647b56933778d82cf.

## Question

LM1H showed that nonlinear symmetric operator splitting can improve worst-family balance, with strang_split_base slightly improving both mean and minimum held-out byte accuracy over cyclic control. Does that optimization improvement remain compatible with the frozen grounded router and exact-memory subsystem?

## Frozen model and training

Start exactly from the LM1H warm three-family state:
- 64-D recurrent transition;
- recurrent parameters frozen;
- 20 base pretrain epochs;
- 20 base/paraphrase joint-interleaved epochs;
- output alphabet extended exactly as LM1D;
- 4 three-family adaptation epochs;
- learning rate 0.08, half-step 0.04.

Arms:
1. cyclic_by_index control
2. strang_split_base

Router/memory:
- exact frozen LM1C third-family grounded router;
- exact recall cap 16;
- no router retraining in this experiment;
- no attention;
- no future oracle.

## Evaluation

Byte quality:
- base heldout
- paraphrase heldout
- third-family heldout

Integrated routing/memory:
- base block stream1 and stream4
- paraphrase block stream1 and stream4
- third-family:
  - block
  - per_name
  - paired_names
  - stores_then_local_reports
  - reverse_report_tail
- each third-family order at stream1 and stream4

## Metrics

- top-1 byte accuracy;
- perplexity;
- dependent first-byte accuracy;
- whole-query-set exact accuracy;
- admission precision/recall;
- event-routing accuracy;
- report-routing accuracy;
- max recall entries.

## Interpretation

If strang_split_base retains exact routing/memory and its byte-balance advantage, it becomes the integrated three-family adaptation procedure for the next lexical-breadth stage.

## Bounds

No recurrent training, no router retraining, no recall-cap increase, no attention, no adaptive schedule, no threshold tuning, no result-informed retry, no live activation, no production authority.
