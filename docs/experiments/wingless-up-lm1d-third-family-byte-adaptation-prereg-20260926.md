# Wingless UP-LM1D — bounded third-family byte adaptation

Status: preregistered scientific continual language-model adaptation experiment.

Scientific parent: sealed UP-LM1C 4ac7668597fbe21ab7d2e7b17f65267cb3912e48.

## Question

LM1C showed perfect grounded routing and exact recall for a third lexical family while the frozen byte readout fell to about 49% top-1 accuracy. Can bounded readout adaptation recover the third family while preserving the already integrated base and paraphrase distributions?

## Frozen starting model

Start from the exact LM1B/LM1C warm20 byte model:
- recurrent transition fixed;
- 20 base-only epochs;
- 20 base/paraphrase joint-interleaved epochs;
- state dimension 64;
- learning rate 0.08;
- exact recall cap 16.

Router/memory:
- exact LM1C grounded router after its 20 third-family grounding epochs;
- router weights frozen during this experiment;
- no attention;
- no future oracle.

## Third-family adaptation corpus

Third-family verbs:
- STORE: archives
- OBSERVE: notices
- REPORT: recounts

Use the exact structural split rule from the LM lineage.

Adaptation training uses only block-order third-family examples.

## Frozen adaptation ladder

Additional adaptation epochs:
- 0
- 1
- 2
- 4

Each adaptation epoch uses matched deterministic example-level updates:
1. third-family training example i;
2. base training example i;
3. paraphrase training example i.

The three training corpora have the same deterministic example count/order. No arm receives extra examples or adaptive replay.

Each ladder arm begins from an identical copy of the same warm20 starting byte model.

## Evaluation

For every adaptation depth:
- base held-out block;
- paraphrase held-out block;
- third-family held-out block;
- third-family per_name;
- third-family paired_names;
- third-family stores_then_local_reports;
- third-family reverse_report_tail.

Third-family order families are evaluated at stream1 and stream4. Base/paraphrase controls are evaluated at stream1.

## Metrics

- byte top-1 accuracy;
- perplexity;
- dependent first-byte accuracy;
- whole-query-set exact accuracy;
- admission precision/recall;
- event-routing accuracy;
- report-routing accuracy;
- max recall entries.

## Interpretation

Third-family byte recovery with retained base/paraphrase quality would establish bounded continual lexical adaptation on top of the integrated LM1B procedure. Base/paraphrase degradation would expose three-family readout interference.

## Bounds

No router retraining, no recurrent training, no state expansion, no recall-cap increase, no attention, no adaptive replay, no threshold tuning, no result-informed retry, no live activation, no production authority.
