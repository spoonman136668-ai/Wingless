# Wingless UP-LM1U — integrated fixed-mass language allocation

Status: preregistered scientific integration experiment.

Scientific parent: sealed UP-LM1T 2dd43d89cdd8480eec3c7be779b2de5f48596159.

## Question

UP-LM1T established a smooth fixed-budget stability/plasticity tradeoff at the byte-readout level. Does that same tradeoff survive the accepted fifth-family router + exact bounded-memory sequence task, or does integration introduce a new failure mode?

## Frozen starting point

Reuse the exact five-family pre-adaptation state, fifth-family corpus, router, and integrated evaluation from UP-LM1P/UP-LM1T:
- 64-D recurrent transition;
- recurrent parameters frozen;
- exact four-family rotating-palindromic training history;
- fifth-family alphabet extension;
- accepted fifth-family router training;
- exact recall cap 16;
- no attention;
- no future oracle.

## Fixed-mass arms

Carry all UP-LM1T arms unchanged:

1. equal_mass: prior lr 0.08 each, fifth lr 0.08.
2. fifth_1p5_mass: prior lr 0.07 each, fifth lr 0.12.
3. fifth_2x_mass: prior lr 0.06 each, fifth lr 0.16.

All arms:
- 4 adaptation epochs;
- one update per family per matched example;
- canonical base→paraphrase→third→fourth→fifth order;
- total learning-rate mass 0.40 per matched example;
- identical update count.

No arm is selected or dropped based on UP-LM1T results.

## Evaluation

For every arm:
- held-out byte accuracy/perplexity on all five families;
- minimum/mean/spread and fifth/prior mean summaries;
- exact accepted fifth-family router metrics;
- integrated fifth-family evaluation for all five orders:
  - block;
  - per_name;
  - paired_names;
  - stores_then_local_reports;
  - reverse_report_tail;
- stream depths 1 and 4 for every order.

The integrated metric is the exact UP-LM1P/UP-LM1M metric and includes its accepted dependent first-byte / query-memory / routing measures.

## Interpretation

If the byte-level mass reallocation trend persists while integrated memory/routing remains stable, the fixed-budget tradeoff is composable with sequence intelligence. If integrated metrics degrade disproportionately despite similar byte scores, the next boundary lies in integration rather than readout allocation.

No arm is declared a winner; exact measured tradeoffs are the result.

## Bounds

No extra learning-rate mass, no extra updates, no recurrent training, no router modification, no recall-cap change, no adaptive weighting, no sixth family, no attention, no future oracle, no result-informed retry, no live activation, no production authority.
