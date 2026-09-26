# Wingless UP-LM1W — mild fixed-mass allocation densification

Status: preregistered scientific language stability/plasticity experiment.

Scientific parent: sealed UP-LM1V c9349f4251421bf9395d0bff3e954463655d73b1.

## Question

UP-LM1V showed that shifting fixed learning mass toward the fifth family improves fifth-family integrated language quality while gradually reducing prior-family quality, with sequence memory and routing remaining exact. What does the mild region between equal allocation and the existing 1.5× fifth allocation look like?

## Frozen starting point

Use the exact UP-LM1V starting state and evaluation path:
- 64-D recurrent transition;
- recurrent parameters frozen;
- exact four-family rotating-palindromic prehistory;
- fifth-family alphabet extension;
- accepted five-family router;
- exact recall cap 16;
- block held-out integrated evaluation;
- stream depths 1 and 4;
- no attention;
- no future oracle.

## Fixed training budget

Every arm:
- 4 adaptation epochs;
- one update per family per matched example;
- canonical base → paraphrase → third → fourth → fifth order;
- total learning-rate mass exactly 0.40 per matched example;
- identical update count.

## Preregistered allocation arms

1. equal_mass
   - prior family rate 0.0800 each;
   - fifth family rate 0.0800.

2. fifth_1p125_mass
   - prior family rate 0.0775 each;
   - fifth family rate 0.0900.

3. fifth_1p25_mass
   - prior family rate 0.0750 each;
   - fifth family rate 0.1000.

4. fifth_1p375_mass
   - prior family rate 0.0725 each;
   - fifth family rate 0.1100.

5. fifth_1p5_mass
   - prior family rate 0.0700 each;
   - fifth family rate 0.1200.

No point is selected adaptively and no wider allocation is introduced after results.

## Evaluation

For every arm:
- held-out byte accuracy/perplexity for all five families;
- accepted router metrics for all five held-out lexical families;
- integrated block evaluation for all five families at stream depths 1 and 4.

Per arm × family integrated summary:
- mean and minimum top-1 accuracy across the two stream depths;
- mean perplexity;
- whether dependent first-byte, query set, updates 0/1/2/4, admission precision/recall, and event/report routing remain exact.

## Interpretation

A smooth curve in this mild region would establish a continuously tunable fixed-budget language stability/plasticity frontier. A knee or discontinuity would locate a more useful boundary for subsequent mechanism work, without declaring an allocation winner.

No post-result threshold is introduced.

## Bounds

No extra learning-rate mass, no extra updates, no recurrent training, no router modification, no recall-cap change, no adaptive weighting, no sixth family, no attention, no future oracle, no result-informed retry, no live activation, no production authority.
