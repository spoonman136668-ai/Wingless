# UP-LM12-SUBSTRATE-INVARIANCE-R1

Controller-frozen preregistration:

```json
{
  "experiment_id": "UP-LM12-SUBSTRATE-INVARIANCE-R1",
  "candidate_id": "UP-LM12-SUBSTRATE-INVARIANCE-R1",
  "harness_id": "wingless-windows-research",
  "changed_paths": [
    ".wingless/qualification-request.json",
    "scripts/run-up-lm12-substrate-invariance-r1.ps1",
    "docs/experiments/up-lm12-substrate-invariance-r1.md",
    "unitary/up_lm12_substrate_invariance_r1.go",
    "cmd/up-lm12-substrate-invariance-r1/main.go"
  ],
  "controls": [
    "Run only against baseline commit 594259006bc4468f37870eae7cca93e7dcd5302c in the wingless-windows-research harness.",
    "Generate each state schedule once per seed and replay the identical schedule for every budget and both substrates.",
    "Change only the state representation: dense fixed array is the reference substrate and sparse keyed map is the test substrate.",
    "Use canonical ascending state-index traversal for every digest; never use native map iteration order as an observable.",
    "Use the same Go process, demand generator, deficit rule, digest function, budget vector, and logical resource limits for both substrates.",
    "Validate exactly 512 distinct state indices for every seed before any measured replay.",
    "Record every logical seed-budget case for both substrates, including failures and zero-budget cases.",
    "Do not change thresholds, seeds, budgets, metrics, compute budget, or interpretation after observing results.",
    "Preserve complete results and partial results from stopped or failed runs as valid negative evidence.",
    "Do not access brokers, credentials, accepted refs, live runtime paths, or cross-lane resources."
  ],
  "fixed_parameters": {
    "baseline_sha": "594259006bc4468f37870eae7cca93e7dcd5302c",
    "budget_count": "5",
    "budgets": "0,64,128,256,512",
    "compute_limit_seconds": "3600",
    "digest_algorithm": "sha256 over canonical length-prefixed observable records",
    "emitted_rows_per_substrate": "20480",
    "evidence_id": "RSE-859fac47358c95101e2c50b7ec4fd555",
    "expected_accepted_writes_per_substrate": "7680",
    "expected_deficit_per_substrate": "12800",
    "expected_total_accepted_writes": "15360",
    "expected_total_deficit": "25600",
    "logical_comparison_cases": "40",
    "observation_order": "ascending state index 0..511",
    "reference_substrate": "dense fixed array of 512 state cells",
    "rows_per_logical_case_per_substrate": "512",
    "seed_count": "8",
    "seeds": "101,211,307,401,503,601,701,809",
    "state_count_per_seed": "512",
    "substrate_case_runs": "80",
    "substrate_count": "2",
    "substrates": "dense-array,sparse-map",
    "test_substrate": "sparse map keyed by state index",
    "total_emitted_rows": "40960"
  },
  "metrics": [
    {
      "name": "logical_comparison_cases",
      "comparator": "==",
      "threshold": 40
    },
    {
      "name": "substrate_case_runs",
      "comparator": "==",
      "threshold": 80
    },
    {
      "name": "total_emitted_rows",
      "comparator": "==",
      "threshold": 40960
    },
    {
      "name": "aggregate_demand_writes",
      "comparator": "==",
      "threshold": 40960
    },
    {
      "name": "aggregate_accepted_writes",
      "comparator": "==",
      "threshold": 15360
    },
    {
      "name": "aggregate_deficit",
      "comparator": "==",
      "threshold": 25600
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
      "name": "budget_overrun_rows",
      "comparator": "==",
      "threshold": 0
    },
    {
      "name": "invalid_state_schedules",
      "comparator": "==",
      "threshold": 0
    },
    {
      "name": "min_unique_states_per_seed",
      "comparator": "==",
      "threshold": 512
    },
    {
      "name": "max_unique_states_per_seed",
      "comparator": "==",
      "threshold": 512
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
    "Stop after 3600 compute seconds and preserve all completed and partial records.",
    "Stop if baseline commit or harness identity differs from the fixed parameters.",
    "Stop if schedule validation finds anything other than 512 distinct state indices for any seed.",
    "Stop if the process attempts broker, credential, accepted-ref, live-runtime, or cross-lane access.",
    "Stop on compile failure, wrapper failure, nonzero scientific-process exit, or inability to write only the harness-designated local result output.",
    "Do not stop early because an outcome metric fails; finish all remaining safe cases so negative evidence remains complete."
  ],
  "positive_meaning": "The hypothesis is supported only if all 40 paired cases complete, both substrates emit all 40,960 rows, every paired digest and accepted/deficit total matches, both substrates satisfy exact conservation and budget bounds, and every cardinality metric equals its preregistered threshold.",
  "negative_meaning": "The hypothesis is falsified if any preregistered equality or control metric misses its threshold, the run reaches a stop condition, or the package cannot produce all fixed cardinalities within the compute limit. Preserve the exact result as negative evidence.",
  "mixed_meaning": "A mixed result occurs if substrate equality holds but a cardinality, budget, or conservation control fails, or if all law controls hold but at least one paired substrate observable differs. Preserve it without changing the claim or rerunning with altered parameters."
}
```
