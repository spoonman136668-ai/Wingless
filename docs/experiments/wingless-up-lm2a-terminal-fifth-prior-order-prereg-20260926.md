# Wingless UP-LM2A — terminal fifth with prior-order variation

Status: preregistered scientific five-family language adaptation experiment.

Scientific parent: sealed UP-LM1Z 89c18dbebdc5c2d26cf763235e97dfb77c6cfcf4.

## Question

UP-LM1Z showed a clear terminal-position bonus for the fifth family, while prior-family means changed little when only fifth position moved. Does that terminal privilege survive when the relative order of the four prior families changes, and can prior-order variation explain the larger prior-family shifts seen in UP-LM1Y?

## Frozen starting point

Reuse the exact UP-LM1Z starting state:
- 64-D recurrent transition;
- recurrent parameters frozen;
- accepted four-family rotating-palindromic prehistory;
- fifth-family alphabet extension;
- accepted five-family router;
- exact recall cap 16;
- no attention;
- no future oracle.

## Frozen allocation arms

All five accepted fixed-mass allocations remain unchanged:
- equal_mass: prior 0.0800 each, fifth 0.0800;
- fifth_1p125_mass: prior 0.0775 each, fifth 0.0900;
- fifth_1p25_mass: prior 0.0750 each, fifth 0.1000;
- fifth_1p375_mass: prior 0.0725 each, fifth 0.1100;
- fifth_1p5_mass: prior 0.0700 each, fifth 0.1200.

Every arm:
- 4 adaptation epochs;
- one update per family per matched example;
- total learning-rate mass exactly 0.40 per matched example;
- identical update count.

## Prior-order arms

Fifth family is always position 5.

1. canonical_prior
   - base → paraphrase → third → fourth → fifth.

2. reverse_prior
   - fourth → third → paraphrase → base → fifth.

3. rotating_prior
   - deterministic cyclic start among the four prior families using (epoch + matched-example index) mod 4;
   - fifth always last.

No order is selected adaptively.

## Evaluation

Evaluate exact block-heldout corpora only, stream depths 1 and 4, for all five lexical families.

Per allocation × prior order × family:
- mean/min integrated top-1 accuracy;
- mean perplexity;
- dependent first-byte accuracy;
- query-set exactness;
- update-count exactness;
- admission precision/recall;
- event/report routing accuracy;
- max recall entries.

Also report per allocation × prior order:
- fifth-family integrated mean;
- prior-four-family integrated mean;
- all-family integrated mean;
- minimum family integrated mean.

## Interpretation

If fifth-family performance remains close across prior-order arms while prior-family means move, terminal fifth privilege is robust and prior-order interaction explains the older-family shift. If fifth performance itself changes materially, terminal privilege depends on the preceding prior-family sequence.

No arm is declared a winner and no numeric success threshold is introduced.

## Bounds

No extra learning-rate mass, no extra updates, no recurrent training, no router modification, no recall-cap change, no adaptive weighting, no sixth family, no attention, no future oracle, no result-informed retry, no live activation, no production authority.
