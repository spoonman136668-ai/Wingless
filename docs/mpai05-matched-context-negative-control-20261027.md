# MP-AI05 matched-context negative control — 2026-10-27

**Result:** Real Nemotron inference completed, but no answer-quality improvement attributable to query normalization when baseline and candidate received identical records.

- Live Actions receipt: https://github.com/spoonman136668-ai/Wingless/actions/runs/40497116262
- Source workflow commit: `0f4c901fb99636ec6ff8c6ac8e24a45b7f0ff330`
- Shared FORGE adapter pinned at `ca373bbab0ab297b79c9073de7e78d96ddf26903`
- Requested model: `nvidia/nemotron-3-ultra-550b-a55b:free` via pre-existing Wingless GitHub secret.
- No private Mind-Palace records transmitted; no confidential first-use scientific holdouts accessed.
- Four **public synthetic** questions: two known normalization patterns and two unrecognized forms. Both arms received the exact same six-record corpus, including current and superseded documents. Order alternated between the arms.
- All eight actual model calls completed. Scoring: required answer-string match, grounded cited source, and syntactically valid response; the evaluator's simple lexical rule is not a substitute for independent human or scientific adjudication.

| Task | Baseline correct | Candidate correct | Baseline tokens | Candidate tokens | Baseline latency ms | Candidate latency ms |
|---|---:|---:|---:|---:|---:|---:|
| Quartz final cutoff | 1 | 1 | 490 | 494 | 2656 | 2599 |
| Harbor final window | 1 | 1 | 470 | 471 | 1763 | 16391 |
| Maple deployment prerequisites | 1 | 1 | 411 | 415 | 976 | 3015 |
| Maple rotation negative | 1 | 1 | 420 | 416 | 9218 | 1234 |
| **Totals** | **4/4** | **4/4** | **1791** | **1796** | **14613** | **23239** |

The candidate used 5 additional tokens (+0.28%) and 8626 ms additional aggregate latency (+59.0%), dominated by provider variation on task 2. The small sample cannot reliably attribute latency differences to the normalizer; the service is an uncontrolled external endpoint.

**Scientific interpretation:** The earlier favorable synthetic setup yielded 2/2 candidate vs 0/2 baseline by guaranteeing that only normalized exact-match queries retrieved useful content. Under this matched-corpus counterexample, both arms answered all questions correctly. Thus the earlier positive result is **not** evidence of independent cognition improvement; it is consistent with fixing a contrived retrieval miss. The null result is important and retained.

**Decision:** `MATCHED_CONTEXT_INFERENCE=PASS`; `ANSWER_QUALITY_GAIN=NOT_OBSERVED`; `RESOURCE_NONREGRESSION=NOT_DEMONSTRATED`; `SCIENTIFIC_ACCEPTANCE=NOT_TESTED`. Do not promote, enable global retrieval normalization, or claim unfamiliar-task transfer. Any next experiment should use a frozen, heterogeneous, realistic corpus and a retrieval evaluator independent of the normalizer; preregister tasks and limit free-inference consumption.
