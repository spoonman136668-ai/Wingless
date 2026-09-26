# Wingless UP-LM1W — mild fixed-mass all-family frontier

Status: preregistered scientific five-family allocation experiment.

Scientific parent: sealed UP-LM1V c9349f4251421bf9395d0bff3e954463655d73b1.

## Question

UP-LM1V showed that the fixed-mass stability/plasticity tradeoff propagates into all-family integrated language scores while exact sequence memory and semantic routing remain perfect. Is there a lower-cost plasticity region between equal allocation and the previous 1.5x fifth-family allocation?

## Frozen starting point

Reuse the exact UP-LM1V five-family starting state:
- 64-D recurrent transition;
- recurrent parameters frozen;
- exact four-family rotating-palindromic training history;
- fifth-family alphabet extension;
- accepted five-family router;
- exact recall cap 16;
- no attention;
- no future oracle.

## Fixed-mass arms

All arms use:
- 4 adaptation epochs;
- one update per family per matched example;
- canonical base→paraphrase→third→fourth→fifth order;
- total learning-rate mass exactly 0.40 per matched example;
- identical update count.

Arms:

1. equal_mass
   - prior lr 0.0800 each
   - fifth lr 0.0800

2. fifth_1p125_mass
   - prior lr 0.0775 each
   - fifth lr 0.0900

3. fifth_1p25_mass
   - prior lr 0.0750 each
   - fifth lr 0.1000

4. fifth_1p375_mass
   - prior lr 0.0725 each
   - fifth lr 0.1100

5. fifth_1p5_mass
   - prior lr 0.0700 each
   - fifth lr 0.1200

These points were fixed before execution and do not depend on intermediate results.

## Evaluation

For every arm:
- held-out byte top-1 accuracy/perplexity on all five families;
- minimum / mean / spread;
- fifth-family accuracy/perplexity;
- prior-family mean accuracy.

Evaluate every family on its held-out block corpus at stream depths 1 and 4 using the exact integrated metric from UP-LM1V:
- top-1 accuracy/perplexity;
- dependent first-byte accuracy;
- query-set exact accuracy;
- update-count exact accuracies;
- admission precision/recall;
- event-routing accuracy;
- report-routing accuracy;
- max recall entries.

Carry the accepted five-family router evaluation unchanged.

## Interpretation

The scientific result is the full frozen tradeoff curve. A smooth region where fifth-family gains initially cost little prior-family or minimum accuracy would identify a practical stability/plasticity operating regime under fixed budget. If every fifth-family gain immediately lowers the prior-family floor proportionally, the tradeoff has no mild shoulder in this interval.

No arm is declared a winner and no numeric success threshold is introduced.

## Bounds

No extra learning-rate mass, no extra updates, no recurrent training, no router modification, no recall-cap change, no adaptive weighting, no sixth family, no attention, no future oracle, no result-informed retry, no live activation, no production authority.
