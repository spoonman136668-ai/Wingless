# MP-AI05: paired public-fixture Nemotron development result (2026-10-27)

## Execution evidence

- Run: https://github.com/spoonman136668-ai/Wingless/actions/runs/40495398737 (completed SUCCESS).
- Workflow commit: `d99ea4ce2a8fd5997c7799fcdab669df04be21ce`.
- Actual source adapter: public FORGE `tools/mpai05_retrieval_consumer.py` and `tools/mpai05_query_normalizer.py` pinned at `ca373bbab0ab297b79c9073de7e78d96ddf26903`.
- Runtime: GitHub-hosted Ubuntu via existing Wingless OpenRouter secret; requested `nvidia/nemotron-3-ultra-550b-a55b:free`.
- Scope: four actual model responses, two public **synthetic** development tasks, each answered under `off` and `eligible` retrieval modes, alternating treatment order.
- No private Mind-Palace checkout, accepted `.ice` data, production writes, or scientific holdouts transmitted.

## Raw logged receipts

| Task | Arm | Routing | Correct and source-grounded | Tokens | Latency (ms) |
|---|---|---|---|---:|---:|
| 1 Quartz build cutoff | off | disabled | false | 165 | 764 |
| 1 Quartz build cutoff | eligible | normalized | true | 202 | 903 |
| 2 Harbor test window | eligible | normalized | true | 207 | 1471 |
| 2 Harbor test window | off | disabled | false | 177 | 1057 |

Baseline: **0/2**, 342 tokens, 1,821 ms. Candidate: **2/2**, 409 tokens, 2,374 ms. Candidate added 67 tokens (+19.6%) and 553 ms (+30.4%). Both arms invoked the same free model, same model request settings, same number of retrieved records and a balanced task/arm order.

## Interpretation and nonclaims

This is an authentic **inference execution** and **mechanism-level synthetic fixture result**, *not* evidence of genuine unfamiliar-world transfer, real customer-task improvement, independent scientific acceptance, or general retrieval superiority.

**Fixture selection is intentionally favorable to normalization:** the fake snapshot retrieves a generic uninformative record for verbose queries and an answer-containing record for exactly normalized core queries. The observed accuracy gap thus largely tests whether the opt-in rewrite reaches a curated matching record. It cannot establish broad performance on natural mixed-query corpora. Also, tokens and latency regressed, so it does **not** meet the prospective conservative resource-neutral acceptance condition in Mind-Palace issue #96.

Follow-up: evaluate heterogeneous, non-leaked public development questions under a fair retrieval setup with no constructed exact-match advantage and equal source-access budgets; only after that seek independently custodied unfamiliar holdout evidence with two authenticated agent classes. Keep the tested result immutable. Do not promote based on this result.

Status: `PUBLIC_SYNTHETIC_MECHANISM_PILOT=EXECUTED`; `PILOT_ANSWER_SUCCESS=2/2_VS_0/2`; `RESOURCE_NONREGRESSION=FAIL`; `REAL_WORLD_IMPROVEMENT=NOT_ESTABLISHED`; `SCIENTIFIC_ACCEPTANCE=NOT_TESTED`.
