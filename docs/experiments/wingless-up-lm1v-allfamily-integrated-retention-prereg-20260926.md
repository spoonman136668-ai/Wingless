# Wingless UP-LM1V — all-family integrated retention under fixed mass

Status: preregistered scientific five-family integration experiment.

Scientific parent: sealed UP-LM1U d47b0f9b03e98b819d31e574f49bcf22d63f2b66.

## Question

UP-LM1U showed that fixed-budget mass reallocation toward the fifth family improves fifth-family byte and integrated accuracy while the fifth-family router and exact memory remain perfect. Do the byte-level losses in the prior four families propagate into their full integrated sequence behavior, or does the accepted router-plus-memory path preserve those tasks despite readout drift?

## Frozen starting point

Reuse the exact UP-LM1U five-family starting state:
- 64-D recurrent transition;
- recurrent parameters frozen;
- exact four-family rotating-palindromic training history;
- fifth-family alphabet extension;
- accepted five-family router;
- exact recall cap 16;
- no attention;
- no future oracle.

## Fixed-mass arms

Carry the exact UP-LM1U arms unchanged:

1. equal_mass
   - prior lr 0.08 each;
   - fifth lr 0.08.

2. fifth_1p5_mass
   - prior lr 0.07 each;
   - fifth lr 0.12.

3. fifth_2x_mass
   - prior lr 0.06 each;
   - fifth lr 0.16.

Every arm:
- 4 adaptation epochs;
- one update per family per matched example;
- canonical base→paraphrase→third→fourth→fifth order;
- total learning-rate mass exactly 0.40 per matched example;
- identical update count.

## Integrated evaluation

Evaluate every family using its held-out block corpus at:
- stream depth 1;
- stream depth 4.

Families:
1. base;
2. paraphrase;
3. third;
4. fourth;
5. fifth.

For each arm × family × stream report the exact integrated metric:
- top-1 accuracy and perplexity;
- dependent first-byte accuracy;
- query-set exact accuracy;
- update-count exact accuracies;
- admission precision/recall;
- event-routing accuracy;
- report-routing accuracy;
- max recall entries.

Also carry the five-family byte metrics and fixed-mass summaries unchanged for alignment with UP-LM1U.

## Interpretation

If prior-family integrated top-1/perplexity degrade in the same direction as their byte readouts while dependent query, routing, and exact-memory metrics remain stable, the fixed-mass tradeoff is specifically a language-readout retention cost rather than a sequence-memory failure. If integrated prior tasks remain stable despite byte losses, the accepted memory/routing path is functionally insulating task behavior from moderate readout drift.

No arm is declared a winner; exact measured tradeoffs are the result.

## Bounds

No extra learning-rate mass, no extra updates, no recurrent training, no router modification, no recall-cap change, no adaptive weighting, no sixth family, no attention, no future oracle, no result-informed retry, no live activation, no production authority.
