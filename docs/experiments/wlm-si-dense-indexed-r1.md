# WLM-SI-DENSE-INDEXED-R1

Controller-frozen preregistration:

```json
{
  "experiment_id": "WLM-SI-DENSE-INDEXED-R1",
  "candidate_id": "C1-SPARSE-DENSE-PAIRED",
  "harness_id": "wingless-windows-research",
  "changed_paths": [
    ".wingless/qualification-request.json",
    "scripts/run-wlm-si-dense-indexed-r1.ps1",
    "docs/experiments/wlm-si-dense-indexed-r1.md",
    "unitary/wlm_si_dense_indexed_r1.go",
    "cmd/wlm-si-dense-indexed-r1/main.go"
  ],
  "controls": [
    "Use the sparse-indexed implementation as the control and the dense-indexed implementation as the sole treatment.",
    "Replay each pair with the identical seed, initial state, state schedule, resource schedule, deficit schedule, arithmetic, and canonical row ordering.",
    "Run in the wingless-windows-research harness only, with no network, broker, credentials, live runtime, accepted-ref mutation, or cross-lane access.",
    "Evaluate all 256 preregistered pairs and preserve every row and aggregate, including exact negative results.",
    "Do not alter thresholds, seeds, metrics, dimensions, budgets, or interpretations after any result is observed."
  ],
  "fixed_parameters": {
    "baseline_sha": "9bfcde8419332d1e0fedff49e45c9cdd11d98cc8",
    "cases_per_seed": "32",
    "changed_dimension": "state-storage substrate only",
    "comparison_granularity": "exact per-row and per-pair comparison",
    "completed_substrate_case_runs": "256*2=512",
    "control_substrate": "sparse-indexed",
    "early_stopping": "forbidden except preregistered integrity or compute-cap conditions",
    "evidence_id": "RSE-e24f879e92e84fa04c5b6e5f6b58d87a",
    "harness": "wingless-windows-research",
    "output_order": "canonical ascending logical state index",
    "paired_case_count": "8*32=256",
    "post_result_tuning": "forbidden",
    "resource_policy": "identical fixed resource bounds for control and treatment",
    "rows_per_substrate_case_run": "64",
    "seed_count": "8",
    "state_schedule_source": "deterministically generated once per seed and case, then shared by both substrates",
    "substrates_per_pair": "2",
    "total_emitted_rows": "512*64=32768",
    "treatment_substrate": "dense-indexed"
  },
  "metrics": [
    {
      "name": "budget_overrun_rows",
      "comparator": "==",
      "threshold": 0
    },
    {
      "name": "completed_substrate_case_runs",
      "comparator": "==",
      "threshold": 512
    },
    {
      "name": "invalid_state_schedules",
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
      "comparator": "==",
      "threshold": 0
    },
    {
      "name": "max_deficit_law_error",
      "comparator": "==",
      "threshold": 0
    },
    {
      "name": "paired_case_count",
      "comparator": "==",
      "threshold": 256
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
      "name": "total_emitted_rows",
      "comparator": "==",
      "threshold": 32768
    }
  ],
  "seeds": [
    1103,
    2207,
    3313,
    4421,
    5531,
    6653,
    7757,
    8861
  ],
  "compute_seconds": 7200,
  "stop_conditions": [
    "Stop and preserve the partial result if elapsed compute reaches 7200 seconds; classify the experiment as incomplete and non-positive.",
    "Stop and preserve the partial result if the harness identity, baseline SHA, seed list, frozen dimensions, or substrate identities do not match the preregistration.",
    "Stop and preserve the partial result if any schedule cannot be reproduced identically for both members of its pair.",
    "Do not stop early because of a scientific mismatch; continue all remaining preregistered cases so negative-result extent is measured.",
    "Do not rerun with altered parameters after observing any result."
  ],
  "positive_meaning": "All 256 pairs and 512 substrate-case runs complete within the fixed compute cap, emit exactly 32768 rows, and satisfy every preregistered zero-error and zero-mismatch threshold; this supports bounded sparse-to-dense substrate invariance only within the frozen dimensions.",
  "negative_meaning": "The hypothesis is falsified if any preregistered metric threshold fails, including a single paired mismatch or nonzero exact law error. An integrity or compute-cap stop is an incomplete negative result, not support for invariance, and must be preserved.",
  "mixed_meaning": "The run is mixed if all completion, budget, schedule-validity, cardinality, and row-count thresholds pass but one invariant family passes while another fails: bounded observable equality, conservation, or the deficit law. Preserve the complete result without changing the interpretation."
}
```
