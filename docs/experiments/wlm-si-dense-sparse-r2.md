# WLM-SI-DENSE-SPARSE-R2

Controller-frozen preregistration:

```json
{
  "experiment_id": "WLM-SI-DENSE-SPARSE-R2",
  "candidate_id": "WLM-C1-DENSE-SPARSE-PAIR",
  "harness_id": "wingless-windows-research",
  "changed_paths": [
    ".wingless/qualification-request.json",
    "scripts/run-wlm-si-dense-sparse-r2.ps1",
    "docs/experiments/wlm-si-dense-sparse-r2.md",
    "unitary/wlm_si_dense_sparse_r2.go",
    "cmd/wlm-si-dense-sparse-r2/main.go"
  ],
  "controls": [
    "Change exactly one experimental dimension: state storage is dense bitset for the control and sparse index set for the treatment; schedules, seeds, state cardinalities, budgets, timesteps, observation procedure, and law are identical.",
    "Generate each canonical toggle trace once and replay the same immutable trace into both substrates.",
    "Initialize both substrates empty and apply exactly one valid in-range toggle before each timestep observation.",
    "Compute active cardinality and active-index checksum for both substrates by scanning indices in ascending order; never use sparse-map iteration order.",
    "Use integer arithmetic only for demand, served, deficit, cardinality, and conservation checks.",
    "Define demand as active cardinality, served as min(demand,budget), and deficit as max(0,demand-budget); independently check demand=served+deficit.",
    "Domain-separate deterministic schedule generation by experiment ID, seed, state cardinality, and schedule ID; budget is excluded so every budget sees the same state traces.",
    "Run all frozen cases without data-dependent parameter changes; preserve partial counts and the exact failing rows if a stop condition occurs.",
    "Do not alter thresholds, seeds, metrics, budgets, dimensions, interpretation, or code after observing a result; any later change requires a new experiment ID and preregistration."
  ],
  "fixed_parameters": {
    "arithmetic": "exact nonnegative integer arithmetic; checksum uses specified uint64 wraparound only",
    "baseline_commit": "f5a14070fe3c09c4639b700997b4fcaeefff91a8",
    "checksum_definition": "unsigned 64-bit FNV-1a over each active uint32 index encoded little-endian in ascending order",
    "completed_substrate_case_runs": "1024",
    "compute_limit_seconds": "7200",
    "conservation_law": "demand=served+deficit",
    "deficit_law": "max(0,demand-resource_budget)",
    "demand_definition": "active_cardinality",
    "dense_layout": "fixed []uint64 bitset of ceil(state_cardinality/64) words",
    "harness": "wingless-windows-research",
    "initial_state": "empty",
    "observables": "active_cardinality,canonical_active_index_checksum,demand,served,deficit",
    "observation_order": "ascending state index",
    "paired_case_count": "512",
    "paired_case_formula": "8 seeds * 4 state cardinalities * 4 resource budgets * 4 schedule IDs = 512",
    "resource_budget_count": "4",
    "resource_budgets": "0,1,4,16",
    "result_policy": "preserve exact positive, negative, mixed, invalid, timeout, and partial results without retuning",
    "row_formula": "1024 substrate-case runs * 64 timesteps = 65536",
    "schedule_domain": "experiment_id|seed|state_cardinality|schedule_id; resource budget excluded",
    "schedule_generator": "SplitMix64 with unsigned 64-bit wraparound and rejection sampling for unbiased in-range indices",
    "schedule_ids": "0,1,2,3",
    "schedules_per_seed_cardinality": "4",
    "served_law": "min(demand,resource_budget)",
    "sparse_layout": "map[uint32]struct{} containing exactly the active indices",
    "state_cardinalities": "8,32,128,512",
    "state_cardinality_count": "4",
    "substrate_case_run_formula": "512 paired cases * 2 substrates = 1024",
    "substrate_count": "2",
    "substrates": "dense-bitset-control,sparse-index-set-treatment",
    "timesteps_per_run": "64",
    "total_emitted_rows": "65536",
    "transition_per_timestep": "toggle exactly one in-range state index before observation",
    "verified_evidence_id": "RSE-e9db560866685437d6eac19ea5ab2bac"
  },
  "metrics": [
    {
      "name": "paired_case_count",
      "comparator": "==",
      "threshold": 512
    },
    {
      "name": "completed_substrate_case_runs",
      "comparator": "==",
      "threshold": 1024
    },
    {
      "name": "total_emitted_rows",
      "comparator": "==",
      "threshold": 65536
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
      "name": "budget_overrun_rows",
      "comparator": "==",
      "threshold": 0
    }
  ],
  "seeds": [
    11,
    23,
    47,
    97,
    193,
    389,
    769,
    1543
  ],
  "compute_seconds": 7200,
  "stop_conditions": [
    "Stop when all 1024 substrate-case runs and 65536 rows have been emitted and evaluated.",
    "Stop at 7200 elapsed seconds; preserve the timeout and all partial rows and counts as a negative incomplete result.",
    "Stop on harness process failure, unreadable frozen artifact, or inability to emit the result schema; preserve the error and any partial output without changing or rerunning parameters.",
    "Stop if a generated index is out of range or a canonical trace cannot be replayed identically; increment invalid_state_schedules, preserve the trace and failure, and do not replace its seed or schedule.",
    "Do not stop merely because a scientific metric fails; continue the frozen matrix when the harness remains valid so the exact negative-result extent is preserved."
  ],
  "positive_meaning": "All 512 paired cases and 1024 substrate-case runs complete, emit exactly 65536 rows, and satisfy every frozen zero-error and zero-mismatch metric, supporting invariance for only the preregistered dimensions and bounds.",
  "negative_meaning": "Any incomplete cardinality, invalid schedule, timeout, row-count mismatch, state-cardinality mismatch, paired observable mismatch, budget overrun, conservation error, or deficit-law error falsifies the full hypothesis or invalidates completion as identified by the corresponding frozen metric; exact failures and partial output are preserved.",
  "mixed_meaning": "If both substrates independently have zero conservation and deficit-law error but paired_observable_mismatch_cases is greater than zero, the shared law survives while substrate-invariant observables are falsified. If paired observables match but either law error is nonzero, substrate agreement survives while the specified resource-response law is falsified. Either pattern is preserved as a non-positive result without retuning."
}
```
