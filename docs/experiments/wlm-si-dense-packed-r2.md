# WLM-SI-DENSE-PACKED-R2

Controller-frozen preregistration:

```json
{
  "experiment_id": "WLM-SI-DENSE-PACKED-R2",
  "candidate_id": "WLM-C1-DENSE-PACKED",
  "harness_id": "wingless-windows-research",
  "changed_paths": [
    ".wingless/qualification-request.json",
    "scripts/run-wlm-si-dense-packed-r2.ps1",
    "docs/experiments/wlm-si-dense-packed-r2.md",
    "unitary/wlm_si_dense_packed_r2.go",
    "cmd/wlm-si-dense-packed-r2/main.go"
  ],
  "controls": [
    "The scientific source is test-only and must not be imported by production or live-runtime code.",
    "The only experimental change within each pair is state representation: dense Boolean array versus packed four-word bitset.",
    "Both substrates receive the same immutable schedule bytes and fixed resource budget; state is freshly initialized before every run.",
    "Schedules are generated before either substrate executes, serialized canonically, hashed, and rejected as an invalid run if their digest changes.",
    "Substrate order is deterministically balanced: dense first for even schedule indices and packed first for odd schedule indices.",
    "The oracle replays toggle indices directly from the immutable schedule and computes expected active count independently of either substrate's enumeration method.",
    "All observables and laws use exact integer arithmetic; no floating-point tolerance is permitted.",
    "Scientific failures do not terminate the case matrix early; every valid case is run and all mismatches are retained.",
    "Thresholds, seeds, schedules, metrics, budgets, and interpretations are frozen before execution and must not be changed after results are observed.",
    "The harness must remain offline and must not access brokers, credentials, accepted refs, live activation, or cross-lane resources."
  ],
  "fixed_parameters": {
    "baseline_sha": "12fdfa153c54a90f063170fb8f4b17be1ae34f65",
    "changed_dimension": "state substrate representation only",
    "comparison_rule": "exact equality at every row; no tolerance",
    "conservation_law": "active_count + inactive_count = 256",
    "deficit_law": "deficit_units = max(0,demand_units-resource_budget_units)",
    "demand_law": "demand_units = active_count",
    "evidence_id": "RSE-ad939416a5eaf9931409a8505cbf5e97",
    "execution_order": "dense then packed for even schedule_index; packed then dense for odd schedule_index",
    "initial_active_count": "128",
    "initial_state_generator": "Sort indices 0..255 by SHA-256 of canonical big-endian seed,schedule_index,index tuples with index as deterministic tie-breaker; activate the first 128 indices",
    "observable_vector": "active_count,inactive_count,demand_units,deficit_units,conservation_total",
    "observation_steps_per_run": "64",
    "paired_case_count": "8 seeds * 32 schedules = 256",
    "resource_budget_units": "128",
    "result_policy": "emit all valid rows and preserve every negative or mixed result without retuning",
    "rows_per_run": "64",
    "schedule_generator": "SHA-256 counter stream over canonical big-endian seed,schedule_index,step tuples; each digest-derived unsigned index is reduced modulo 256",
    "schedules_per_seed": "32",
    "seed_count": "8",
    "state_width": "256 bits",
    "substrate_case_run_count": "256 paired cases * 2 substrates = 512",
    "substrate_count": "2",
    "substrates": "dense-bool-array,packed-4xuint64-bitset",
    "total_emitted_rows": "512 runs * 64 rows = 32768",
    "transition_rule": "toggle exactly one indexed state bit per observation step"
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
    "Stop before scientific execution if the checked-out commit is not exactly 12fdfa153c54a90f063170fb8f4b17be1ae34f65.",
    "Stop if the sealed qualification request or preregistration digest does not match the qualified artifact.",
    "Stop and preserve the invalid attempt if schedule serialization or its pre-execution digest changes after generation.",
    "Stop and preserve the invalid attempt on compilation failure, harness failure, or nonzero probe exit before the case matrix begins.",
    "Stop and preserve the partial result if elapsed compute reaches 3600 seconds.",
    "Stop immediately if any code path requests network, broker, credential, accepted-ref, live-runtime, or cross-lane access.",
    "Do not stop for a scientific mismatch; finish the valid case matrix so negative-result cardinality is preserved."
  ],
  "positive_meaning": "Support the bounded hypothesis only if all 256 pairs and 512 substrate runs complete, exactly 32768 rows are emitted, and every preregistered zero-threshold metric remains exactly zero.",
  "negative_meaning": "If the run is complete and valid but any paired observable mismatch, conservation error, deficit-law error, state-cardinality mismatch, row-count mismatch, or budget overrun is nonzero, reject exact dense-to-packed substrate invariance and preserve the exact failing cases. An integrity stop produces an invalid sealed attempt, not support for the hypothesis, and must also be preserved.",
  "mixed_meaning": "If integrity and cardinality metrics are complete but only a subset of paired cases or observables disagree, conclude bounded substrate dependence localized to the recorded seeds, schedules, steps, and observable fields; preserve the complete result and do not alter the preregistration."
}
```
