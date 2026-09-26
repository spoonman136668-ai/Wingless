# Wingless UP-LM2J — fixed working-set scheduling

Status: preregistered scientific fixed-capacity language scheduling experiment.

Scientific parent: sealed UP-LM2I.

## Question

UP-LM2I established a hard 16-entity exact-recall boundary: with 24 simultaneously stored entities, recall is exactly 16/24 while event routing remains perfect. Can fixed nonadaptive sequence organization keep the active working set at or below 16 and recover exact reporting without increasing memory, training, or model capacity?

## Frozen training and capacity

Reuse UP-LM2I unchanged:
- 64-D recurrent transition;
- exact recall cap 16;
- same five lexical families;
- same five fixed-mass allocations;
- canonical_prior and balanced_prior training schedules;
- fifth family always last;
- 4 adaptation epochs;
- total learning-rate mass 0.40 per matched example;
- no recurrent retraining, router modification, or memory increase.

## Frozen entities and event multiset

Use the exact same 24 deterministic entity names and value assignment as UP-LM2I.

Every arm contains exactly:
- 24 STORE clauses;
- 24 OBSERVE clauses;
- 24 REPORT clauses.

Only event order changes.

## Working-set schedules

1. global24
   - all 24 STORE;
   - all 24 OBSERVE;
   - all 24 REPORT.

2. chunk16_8
   - entities 1..16: STORE all, OBSERVE all, REPORT all;
   - then entities 17..24 likewise.

3. chunk12x2
   - two consecutive 12-entity STORE/OBSERVE/REPORT chunks.

4. chunk8x3
   - three consecutive 8-entity STORE/OBSERVE/REPORT chunks.

Reports are forward within each chunk.

## Preregistered capacity expectation

- global24 exact-recall hit rate = 16/24;
- every chunked schedule exact-recall hit rate = 1.0 because no active chunk exceeds 16.

## Measurements

For every allocation × training schedule × family × working-set schedule:
- top-1 accuracy;
- perplexity;
- dependent first-byte accuracy;
- report-set exact accuracy;
- exact-recall hit rate;
- preregistered expected hit rate;
- event/STORE/OBSERVE/REPORT routing accuracy;
- maximum recall entries.

Total cells: 5 × 2 × 5 × 4 = 200.

## Interpretation

Exact reporting under chunked schedules with global24 remaining at 16/24 demonstrates that sequence organization can avoid fixed-capacity overflow without adding memory. Failure despite recall-hit recovery would identify a recurrent/language sequencing limit.

## Bounds

No extra clauses, no capacity change, no extra training, no adaptive chunk size, no router modification, no attention, no future oracle, no result-informed retry, no live activation, no production authority.
