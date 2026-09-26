# Wingless UP-LM0P — base replay ratio during lexical adaptation

Status: preregistered scientific continual language-model adaptation experiment.

Scientific parent: sealed UP-LM0O d2c6fd9ed91b607cbb86b4dc3615cbff5342b6c8.

## Question

LM0O showed rapid paraphrase adaptation but a modest drop on the original byte-language distribution. At fixed paraphrase adaptation depth, can additional deterministic base replay preserve the original distribution without sacrificing the new lexical surfaces?

## Frozen starting point

Identical to LM0O:
- deterministic 64-D recurrent language state;
- 20 base training epochs before adaptation;
- learning rate 0.08;
- exact recall cap 16;
- grounded paraphrase router frozen;
- no router retraining;
- no attention;
- no future oracle.

## Frozen adaptation depth

Paraphrase adaptation is fixed at:
- 4 epochs.

Every adaptation epoch performs exactly one pass over the LM0O paraphrase training corpus.

## Frozen base replay ladder

Immediately after each paraphrase pass, perform:
- 0 base replay passes
- 1 base replay pass
- 2 base replay passes
- 4 base replay passes

Each arm starts from an identical copy of the same 20-epoch base model.

No replay examples are selected adaptively; a replay pass is the complete original base training corpus in deterministic order.

## Evaluation

For every replay arm:
1. original base held-out corpus:
   - top-1 byte accuracy
   - perplexity
2. paraphrase block-order held-out corpus:
   - top-1 byte accuracy
   - perplexity
   - dependent first-byte accuracy
   - whole-query-set exact accuracy
3. four LM0N paraphrase order families at stream1:
   - learned grounded-paraphrase routing only
   - top-1 byte accuracy
   - perplexity
   - dependent first-byte accuracy
   - whole-query-set exact accuracy
   - event-routing accuracy
   - report-routing accuracy
   - maximum recall entries

## Interpretation

This is a preregistered retention/adaptation curve, not a post-hoc replay search. A replay level that improves base retention while preserving paraphrase gains would support cheap continual lexical adaptation of the byte model.

## Bounds

No router retraining, no adaptive replay, no state expansion, no recall-cap increase, no attention, no threshold tuning, no choosing a replay ratio after results, no result-informed retry, no live activation, no production authority.
