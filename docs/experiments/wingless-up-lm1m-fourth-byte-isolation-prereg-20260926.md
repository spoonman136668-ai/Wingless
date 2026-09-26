# Wingless UP-LM1M — corrected fourth-router byte-adaptation isolation

Status: preregistered scientific language/memory integration experiment.

Scientific parent: sealed UP-LM1L bab2384f16a154b0af448b96fd096cca3e897a2c.

## Question

LM1L made fourth-family semantic routing exact by correcting only the "states" REPORT representation. With routing no longer confounded, how much of the remaining fourth-family failure is purely byte-language out-of-distribution error, and how much can be recovered by bounded fourth-family byte adaptation?

## Frozen semantic system

Router:
- exact LM1L states-vs-STORE orthogonal correction;
- same 20-epoch fourth-family grounding;
- same prior-family rehearsal;
- correction frozen;
- no router changes between arms.

Memory:
- exact recall cap 16;
- same integrated evaluator as LM1J/LM1I;
- no attention;
- no future oracle.

## Frozen byte-model starting point

Use the exact integrated three-family model:
- 64-D recurrent transition frozen;
- 20 base-only epochs;
- 20 base/paraphrase joint-interleaved epochs;
- three-family Strang-split-base adaptation for 4 epochs;
- no fourth-family examples before the adaptation arm.

## Arms

1. no_fourth_byte_adaptation
   - evaluate the frozen three-family byte model with corrected fourth router/memory.

2. fourth_only_4ep
   - clone the same starting byte model;
   - train only the fourth-family byte corpus for 4 epochs at lr 0.08;
   - recurrent transition remains frozen;
   - router is not retrained.

3. four_family_interleaved_4ep
   - clone the same starting byte model;
   - for 4 epochs, train matched base / paraphrase / third / fourth examples in deterministic cyclic order;
   - one example from each family per index;
   - lr 0.08;
   - identical fourth-family example count to arm 2 plus rehearsal of prior families.

## Evaluation

For every arm:
- base / paraphrase / third / fourth held-out byte top-1 accuracy and perplexity;
- fourth-family dependent first-byte accuracy;
- fourth-family query-set exact accuracy;
- admission precision/recall;
- event-routing accuracy;
- report-routing accuracy;
- max recall entries.

Evaluate fourth family under:
- block
- per_name
- paired_names
- stores_then_local_reports
- reverse_report_tail
at stream1 and stream4.

## Interpretation

If routing/memory stays exact while fourth byte quality improves, LM1J’s remaining failure is isolated to byte adaptation. Comparing fourth-only with four-family rehearsal maps the adaptation-vs-retention tradeoff without changing recurrent state or memory capacity.

## Bounds

No recurrent training, no router retraining, no recall-cap increase, no attention, no adaptive schedule, no threshold tuning, no result-informed retry, no live activation, no production authority.
