# MP-AI05 public source-backed answer quality: matched full vs selected context

Date: 2026-10-27. Actual real-model run: https://github.com/spoonman136668-ai/Wingless/actions/runs/40509363334 (PASS). Workflow source commit `98a6d8c96db3e9f008348d7c9707b91f9c24ef1e`.

## Exact scope

This experiment used **public** FORGE documentation `docs/MPAI05-ON-DEMAND-CONSUMER.md` pinned at FORGE commit `ca373bbab0ab297b79c9073de7e78d96ddf26903`, with the pre-existing Wingless GitHub OpenRouter credential and requested `nvidia/nemotron-3-ultra-550b-a55b:free`. Both arms answered the same four source-grounded questions; order alternated. The baseline received the complete public document, and the compact arm received the top three paragraphs selected by a **simple deterministic lexical overlap script**. The compact selector is **NOT Mind-Palace's actual tiered context pack**, and the selection cases are not an independent hidden holdout.

No private Mind-Palace records or accepted `.ice` records were sent to the model. No paid API keys, new credentials or controllers were used.

## Observed receipts

| Task | Full answer PASS | Selected answer PASS | Full tokens | Selected tokens | Full source bytes | Selected source bytes |
|---|---:|---:|---:|---:|---:|---:|
| 1 original vs shadow mode | true | true | 949 | 740 | 2,974 | 1,881 |
| 2 global activation after merge | true | true | 982 | 724 | 2,974 | 1,701 |
| 3 pinned source SHA | true | true | 993 | 784 | 2,974 | 1,881 |
| 4 eligible match count | true | true | 950 | 753 | 2,974 | 1,935 |
| **TOTAL** | **4/4** | **4/4** | **3,874** | **3,001** | **11,896** | **7,398** |

The compact arm used **873 fewer tokens (22.53% reduction)** and **4,498 fewer source bytes (37.81% reduction)**, retaining perfect rubric-based correct-source answers in these four cases. The token counts include full model request and response usage. Small-sample provider latency varied; no speedup inference is warranted. Source-citation check required the exact public documentation URL; answer scoring used predeclared answer strings, not independent human adjudication.

## Interpretation

This is evidence that **bounded public source selection can preserve answer correctness while reducing inference tokens** in an actual Nemotron invocation. It is **not** causal evidence that MP's tiered retrieval itself improves model answers: the live test uses a surrogate lexical paragraph selector and one short public document, and the real-corpus comparison in FORGE run #40506908983 used no model calls. There is no genuine private-corpus combined end-to-end inference or cross-agent/holdout first-use claim.

The pair of experiments together motivates further real-source, privately contained answer validation but **does not qualify scientific acceptance**. Preserve the equal-context null, resource regressions, and identity boundaries. Do not enable global deployment or export private memory to third-party inference.

Status: `PUBLIC_DOC_LIVE_INFERENCE=PASS`; `SELECTED_CORRECT=4/4`; `FULL_CORRECT=4/4`; `REQUEST_TOKENS=3001_VS_3874`; `TRUE_MP_TIERED_QA=NOT_TESTED`; `SCIENTIFIC_ACCEPTANCE=NOT_TESTED`.
