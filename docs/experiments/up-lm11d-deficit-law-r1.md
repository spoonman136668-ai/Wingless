# UP-LM11D-DEFICIT-LAW-R1

Controller-frozen preregistration:

```json
{
  "experiment_id": "UP-LM11D-DEFICIT-LAW-R1",
  "candidate_id": "C2-BOUNDED-WRITE-DEFICIT-LAW",
  "harness_id": "wingless-windows-research",
  "changed_paths": [
    ".wingless/qualification-request.json",
    "scripts/run-up-lm11d-deficit-law-r1.ps1",
    "docs/experiments/up-lm11d-deficit-law-r1.md",
    "unitary/up_lm11d_deficit_law_r1.go",
    "cmd/up-lm11d-deficit-law-r1/main.go"
  ],
  "controls": [
    "Use baseline commit 98b9aecf97719b7168d05264c6481a7c3f80c684 and the wingless-windows-research harness only.",
    "Change only the write-budget level; keep source, substrate, state schema, demand count, seed, permutation, ordering, and comparison logic fixed within each seed.",
    "Generate exactly one deterministic permutation of 512 unique states per seed and reuse that permutation unchanged across all five budgets.",
    "Run cases serially with no concurrency, wall-clock metric, adaptive retry, or early stopping based on scientific results.",
    "Record attempted writes independently before budget admission; derive accepted and deficit totals from separately recorded admission outcomes.",
    "Construct each expected prefix state independently from the same-seed full-demand permutation, not from the bounded run's reported counters.",
    "Emit every preregistered row and case summary even after a scientific threshold has failed.",
    "Use exact integer and canonical-byte comparisons only; no tolerances, discarded seeds, threshold changes, or post-result interpretation changes."
  ],
  "fixed_parameters": {
    "admission_law": "accepted=min(write_budget,512); deficit=512-accepted",
    "budget_cases": "8 seeds * 5 budgets = 40",
    "budget_level_count": "5",
    "budget_levels": "[0,64,128,256,512]",
    "conservation_law": "accepted+deficit=512 for every seed-budget case",
    "demand_writes_per_seed_budget_case": "512",
    "execution_order": "ascending seed list, then budgets [0,64,128,256,512]",
    "expected_aggregate_accepted_writes": "8 * 960 = 7680",
    "expected_aggregate_deficit": "20480 - 7680 = 12800",
    "expected_aggregate_demand_writes": "8 * 5 * 512 = 20480",
    "expected_emitted_rows": "8 seeds * 5 budgets * 512 demand writes = 20480",
    "expected_prefix_digest_comparisons": "8 seeds * 5 budgets = 40",
    "experimental_dimension": "write_budget_only",
    "parallelism": "1",
    "result_retention": "retain complete output for positive, mixed, negative, and operationally incomplete outcomes",
    "seed_count": "8",
    "state_oracle": "canonical state after the first accepted entries of the same-seed full-demand permutation",
    "sum_budget_per_seed": "0 + 64 + 128 + 256 + 512 = 960",
    "unique_states_per_seed": "512"
  },
  "metrics": [
    {
      "name": "emitted_rows",
      "comparator": "==",
      "threshold": 20480
    },
    {
      "name": "budget_cases",
      "comparator": "==",
      "threshold": 40
    },
    {
      "name": "invalid_permutations",
      "comparator": "==",
      "threshold": 0
    },
    {
      "name": "unique_states_per_seed",
      "comparator": "==",
      "threshold": 512
    },
    {
      "name": "aggregate_demand_writes",
      "comparator": "==",
      "threshold": 20480
    },
    {
      "name": "aggregate_accepted_writes",
      "comparator": "==",
      "threshold": 7680
    },
    {
      "name": "aggregate_deficit",
      "comparator": "==",
      "threshold": 12800
    },
    {
      "name": "law_violation_cases",
      "comparator": "==",
      "threshold": 0
    },
    {
      "name": "budget_overrun_rows",
      "comparator": "==",
      "threshold": 0
    },
    {
      "name": "max_conservation_error",
      "comparator": "==",
      "threshold": 0
    },
    {
      "name": "prefix_digest_comparisons",
      "comparator": "==",
      "threshold": 40
    },
    {
      "name": "prefix_digest_mismatch_cases",
      "comparator": "==",
      "threshold": 0
    }
  ],
  "seeds": [
    104729,
    130363,
    155921,
    181081,
    206369,
    231701,
    257053,
    282407
  ],
  "compute_seconds": 1200,
  "stop_conditions": [
    "Before execution, stop as operationally invalid if the baseline SHA, evidence reference, harness identity, preregistration, changed-path set, seed list, budget list, or source hashes do not match the sealed qualification request.",
    "Stop and retain partial output on process crash, nonzero probe exit, output-schema failure, or compute time exceeding 1200 seconds; classify these as operationally incomplete rather than scientific negatives.",
    "Do not stop early for mismatches or threshold failures; complete every seed-budget case so negative results remain exact and unbiased.",
    "Do not retry, alter seeds, budgets, metrics, thresholds, ordering, source, or interpretation after any result is observed."
  ],
  "positive_meaning": "All 40 cases and 20480 demand rows are present; every exact cardinality and aggregate matches; no permutation, admission-law, conservation, budget-overrun, or prefix-digest violation occurs.",
  "negative_meaning": "A complete valid run with any preregistered scientific metric outside its threshold falsifies the exact deficit-law hypothesis; preserve all rows, summaries, seeds, and failures unchanged.",
  "mixed_meaning": "All runs and exact cardinalities complete, but at least one conservation, admission, budget-overrun, or prefix-equivalence metric fails while another passes; preserve the complete result and classify the hypothesis as not supported without retuning."
}
```
