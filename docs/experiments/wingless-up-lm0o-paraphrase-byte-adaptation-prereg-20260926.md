# Wingless UP-LM0O — bounded paraphrase byte-model adaptation

Status: preregistered scientific language-model adaptation experiment.

Scientific parent: sealed UP-LM0N 7ca4aefc2d11c6e2bc94267d8d11b90e3770c8e0.

## Question

LM0N showed that grounded event routing and exact-memory behavior remain perfect under paraphrase plus clause-order variation, while the frozen byte language model falls to roughly 59–61% top-1 accuracy. How quickly can the byte model adapt to the new lexical surfaces without losing the original language distribution?

## Frozen starting point

Start from the exact LM0N system:
- deterministic 64-D recurrent language state;
- language model first trained for the original 20 base epochs;
- exact recall cap 16;
- grounded paraphrase router unchanged from LM0N;
- router weights frozen throughout this experiment;
- no attention;
- no future oracle.

## Paraphrase adaptation corpus

Use the same names, values, split rule, and update-count schedule as the LM0F/LM0N lineage.

Replace only event verbs:
- stores -> saves
- observes -> sees
- reports -> recalls

Adaptation examples use the original block-ordered clause structure, not the new order permutations.

Training split:
- exact existing structural train split.

Held-out split:
- exact existing structural held-out split.

## Frozen adaptation ladder

Additional paraphrase adaptation epochs:
- 0
- 1
- 2
- 4

For each ladder point:
- begin from an identical copy of the 20-epoch base byte model;
- each adaptation epoch performs one pass over the paraphrase training corpus at learning rate 0.08;
- immediately after that pass, perform exactly one replay pass over the original base training corpus at learning rate 0.08;
- no other training.

Thus the ladder measures bounded adaptation depth while keeping replay fixed.

## Evaluation

For every adaptation depth:
1. original base held-out corpus:
   - top-1 byte accuracy
   - perplexity
2. paraphrase held-out block-order corpus:
   - top-1 byte accuracy
   - perplexity
   - dependent first-byte accuracy
   - whole-query-set exact accuracy
3. paraphrase structural breadth:
   - the four LM0N order families at stream1
   - learned grounded-paraphrase routing only
   - top-1 byte accuracy
   - perplexity
   - dependent first-byte accuracy
   - whole-query-set exact accuracy
   - event-routing accuracy
   - report-routing accuracy
   - maximum recall entries

## Interpretation

Improved paraphrase byte accuracy with retained base accuracy would show that the current language bottleneck is quickly adaptable lexical modeling rather than routing or memory. Base degradation would expose continual-language interference analogous to the B lane.

## Bounds

No router retraining, no state expansion, no recall-cap increase, no attention, no threshold tuning, no adaptive replay, no choosing an epoch after seeing results, no result-informed retry, no live activation, no production authority.
