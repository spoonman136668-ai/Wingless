# WLM-SI-SPARSE-INDEXED-R1

Controller-frozen preregistration:

```json
{
  "experiment_id": "WLM-SI-SPARSE-INDEXED-R1",
  "candidate_id": "WLM-SI-SPARSE-INDEXED-R1",
  "harness_id": "wingless-windows-research",
  "changed_paths": [
    ".wingless/qualification-request.json",
    "scripts/run-wlm-si-sparse-indexed-r1.ps1",
    "docs/experiments/wlm-si-sparse-indexed-r1.md",
    "unitary/wlm_si_sparse_indexed_r1.go",
    "cmd/wlm-si-sparse-indexed-r1/main.go"
  ],
  "controls": [
    "Verify the baseline commit is exactly 7766b52158b3012c22e05b4b81268798f6d1c8dc before running.",
    "Verify evidence RSE-f5dd4302248b427ec1d92659e6a2888d against SHA-256 0656bd67d22ae044fc0885dac1e18da5a580ef6ee71d8059c3e299790daae420 before running.",
    "Change only the substrate representation between paired runs; use identical seeds, generated cases, resource schedules, row limits, and comparison logic.",
    "Use independent dense and sparse-indexed encoders and a shared, frozen canonical row comparator.",
    "Execute every planned case after scientific mismatches so exact negative-result cardinalities are preserved.",
    "Do not alter thresholds, seeds, metrics, budgets, case counts, ordering, or interpretation after any result is observed.",
    "Disable network, broker, credential, accepted-reference, live-runtime, and cross-lane access."
  ],
  "fixed_parameters": {
    "arithmetic": "exact integer arithmetic",
    "baseline_sha": "7766b52158b3012c22e05b4b81268798f6d1c8dc",
    "canonicalization": "ascending canonical state key; integer values encoded losslessly; no tolerance or row elision",
    "case_order": "ascending seed-list position, then case index 0 through 31, with dense immediately followed by sparse-indexed",
    "cases_per_seed": "32",
    "comparison_scope": "all 64 canonical emitted rows, state cardinality, conservation error, and deficit-law error for every paired case",
    "completed_substrate_case_runs": "512",
    "compute_limit_seconds": "3600",
    "continue_after_scientific_failure": "true",
    "evidence_id": "RSE-f5dd4302248b427ec1d92659e6a2888d",
    "harness": "wingless-windows-research",
    "paired_case_count": "256",
    "reference_substrate": "dense",
    "resource_schedule": "0,1,2,4,8,16,32,64 repeated in that order for 8 cycles per substrate-case run",
    "rows_per_substrate_case_run": "64",
    "seed_count": "8",
    "substrate_count": "2",
    "total_emitted_rows": "32768",
    "treatment_substrate": "sparse-indexed"
  },
  "metrics": [
    {
      "name": "completed_substrate_case_runs",
      "comparator": "==",
      "threshold": 512
    },
    {
      "name": "paired_case_count",
      "comparator": "==",
      "threshold": 256
    },
    {
      "name": "total_emitted_rows",
      "comparator": "==",
      "threshold": 32768
    },
    {
      "name": "paired_observable_mismatch_cases",
      "comparator": "==",
      "threshold": 0
    },
    {
      "name": "row_count_mismatch_cases",
      "comparator": "==",
      "threshold": 0
    },
    {
      "name": "state_cardinality_mismatch_cases",
      "comparator": "==",
      "threshold": 0
    },
    {
      "name": "law_violation_cases",
      "comparator": "==",
      "threshold": 0
    },
    {
      "name": "max_conservation_error",
      "comparator": "\u003c=",
      "threshold": 0
    },
    {
      "name": "max_deficit_law_error",
      "comparator": "\u003c=",
      "threshold": 0
    },
    {
      "name": "budget_overrun_rows",
      "comparator": "==",
      "threshold": 0
    },
    {
      "name": "invalid_state_schedules",
      "comparator": "==",
      "threshold": 0
    }
  ],
  "seeds": [
    101,
    211,
    307,
    401,
    503,
    601,
    701,
    809
  ],
  "compute_seconds": 3600,
  "stop_conditions": [
    "Stop before scientific execution if the baseline SHA or evidence SHA-256 verification fails.",
    "Stop if execution would require live activation, broker access, credential access, accepted-reference mutation, network access, or cross-lane authority.",
    "Stop if the sealed preregistration, qualification request, script, scientific source, or probe wrapper differs from its qualified content.",
    "Stop at 3600 compute seconds and classify incomplete execution as inconclusive while preserving all partial results.",
    "Do not stop early for a scientific mismatch; finish all planned cases unless an operational safety or integrity stop condition occurs."
  ],
  "positive_meaning": "The hypothesis is supported only if all 512 substrate-case runs complete, exactly 32768 rows are emitted, and every preregistered metric satisfies its threshold with no invalid schedules or budget overruns.",
  "negative_meaning": "Any nonzero paired observable mismatch, row-count mismatch, state-cardinality mismatch, law violation, conservation error, or deficit-law error falsifies the hypothesis; preserve the complete exact result without tuning or reinterpretation.",
  "mixed_meaning": "If all completed scientific comparisons show zero violations but any required run or row is missing, or an operational integrity check fails without producing a scientific mismatch, classify the experiment as inconclusive, preserve all emitted results, and do not rerun with changed parameters."
}
```
