# UP-LM14-PAIRED-STATE-SUBSTRATE-R1

Controller-frozen preregistration:

```json
{
  "experiment_id": "UP-LM14-PAIRED-STATE-SUBSTRATE-R1",
  "candidate_id": "C1-PAIRED-STATE-SUBSTRATE",
  "harness_id": "wingless-windows-research",
  "changed_paths": [
    ".wingless/qualification-request.json",
    "scripts/run-up-lm14-paired-state-substrate-r1.ps1",
    "docs/experiments/up-lm14-paired-state-substrate-r1.md",
    "unitary/up_lm14_paired_state_substrate_r1.go",
    "cmd/up-lm14-paired-state-substrate-r1/main.go"
  ],
  "controls": [
    "Replay the same precomputed immutable state-index trace in both substrate arms of every pair.",
    "Hold the write cap, demand per row, state cardinality, run length, seed set, and schedules per seed fixed across substrate arms.",
    "Generate each trace solely from its preregistered seed, schedule index, step index, and SHA-256 counter construction; substrate execution cannot consume randomness.",
    "Use direct indexed access during transitions and compare complete state vectors in canonical ascending-index order, preventing map traversal order from controlling evolution or comparison.",
    "Initialize both substrates from the same all-zero logical state and reject any run whose initial canonical vectors differ.",
    "Evaluate row-local demand equals accepted writes plus deficit and the cap bound independently in each substrate arm.",
    "Validate all frozen run and row cardinalities before interpreting scientific mismatches.",
    "Run the complete frozen matrix once; preserve valid negative and mixed results without changing seeds, thresholds, metrics, budgets, or interpretations."
  ],
  "fixed_parameters": {
    "accepted_write_derivation": "128 runs * min(64 demand writes,36 cap) = 4608 aggregate accepted writes",
    "baseline_sha": "35762f62d13404c252fb0cdf0954dc4510189598",
    "deficit_derivation": "8192 aggregate demand writes - 4608 aggregate accepted writes = 3584 aggregate deficit",
    "demand_derivation": "8192 emitted rows * 1 demand write per row = 8192 aggregate demand writes",
    "demand_writes_per_row": "1",
    "execution_count": "one complete sealed execution",
    "expected_aggregate_accepted_writes": "4608",
    "expected_aggregate_deficit": "3584",
    "expected_aggregate_demand_writes": "8192",
    "fixed_write_cap_per_run": "36",
    "initial_state": "all 256 counters equal zero",
    "observable_comparison": "exact equality of the complete 256-counter vector, cumulative accepted writes, cumulative deficit, and remaining cap after every step",
    "paired_schedule_count": "64",
    "result_policy": "preserve exact positive, negative, and mixed outputs; no post-result tuning",
    "row_count_derivation": "128 substrate case runs * 64 steps per run = 8192 emitted rows",
    "run_matrix_derivation": "8 seeds * 8 schedules per seed * 2 substrates = 128 substrate case runs",
    "schedule_generator": "SHA-256 counter construction over seed,schedule_index,step_index; state_index is the first 64-bit big-endian word modulo 256",
    "schedules_per_seed": "8",
    "seed_count": "8",
    "state_cardinality": "256",
    "state_substrates": "dense-slice,sparse-map",
    "steps_per_run": "64",
    "substrate_case_runs": "128",
    "substrate_variant_count": "2",
    "total_emitted_rows": "8192",
    "verified_evidence_id": "RSE-d525644a8db1666c030eb35c8d8a8f53"
  },
  "metrics": [
    {
      "name": "paired_observable_mismatch_cases",
      "comparator": "==",
      "threshold": 0
    },
    {
      "name": "accepted_write_mismatch_cases",
      "comparator": "==",
      "threshold": 0
    },
    {
      "name": "deficit_mismatch_cases",
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
      "name": "invalid_state_schedules",
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
      "name": "budget_overrun_rows",
      "comparator": "==",
      "threshold": 0
    },
    {
      "name": "paired_schedule_count",
      "comparator": "==",
      "threshold": 64
    },
    {
      "name": "substrate_case_runs",
      "comparator": "==",
      "threshold": 128
    },
    {
      "name": "total_emitted_rows",
      "comparator": "==",
      "threshold": 8192
    },
    {
      "name": "aggregate_demand_writes",
      "comparator": "==",
      "threshold": 8192
    },
    {
      "name": "aggregate_accepted_writes",
      "comparator": "==",
      "threshold": 4608
    },
    {
      "name": "aggregate_deficit",
      "comparator": "==",
      "threshold": 3584
    }
  ],
  "seeds": [
    1103,
    2207,
    3319,
    4421,
    5527,
    6637,
    7753,
    8861
  ],
  "compute_seconds": 3600,
  "stop_conditions": [
    "Do not start if the baseline SHA, evidence identifier, North-Star SHA-256, qualification request, preregistration, or changed-path manifest does not match the sealed request.",
    "Stop and classify as mixed if either substrate cannot initialize to the exact all-zero 256-counter logical state.",
    "Stop and classify as mixed if a generated state index is outside 0 through 255 or paired arms receive nonidentical immutable traces.",
    "Stop and classify as mixed if elapsed compute reaches 3600 seconds before the full frozen matrix completes.",
    "Stop after exactly 128 substrate case runs and 8192 emitted rows; reject any additional run or row.",
    "Do not rerun, extend, substitute seeds, modify thresholds, change metrics, increase budget, or reinterpret outcomes after observing any result."
  ],
  "positive_meaning": "All preregistered integrity and scientific metrics meet their exact thresholds in the single sealed execution. This supports, but does not prove beyond scope, substrate invariance and exact deficit conservation for the two tested state substrates at the fixed resource budget, dimensions, schedules, and seeds.",
  "negative_meaning": "With all qualification and cardinality controls valid, any nonzero paired observable, accepted-write, deficit, conservation-law, or budget-bound violation is a preserved negative result that falsifies substrate invariance for the tested dense-slice and sparse-map substrates under these frozen conditions.",
  "mixed_meaning": "The result is mixed and non-confirmatory if any qualification, initialization, trace, state-cardinality, run-cardinality, or row-cardinality control fails, or if execution stops before the complete frozen matrix is emitted. Preserve the exact output and do not alter or rerun the preregistration based on it."
}
```
