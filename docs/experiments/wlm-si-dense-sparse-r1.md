# WLM-SI-DENSE-SPARSE-R1

Controller-frozen preregistration:

```json
{
  "experiment_id": "WLM-SI-DENSE-SPARSE-R1",
  "candidate_id": "WLM-SI-DENSE-SPARSE-R1",
  "harness_id": "wingless-windows-research",
  "changed_paths": [
    ".wingless/qualification-request.json",
    "scripts/run-wlm-si-dense-sparse-r1.ps1",
    "docs/experiments/wlm-si-dense-sparse-r1.md",
    "unitary/wlm_si_dense_sparse_r1.go",
    "cmd/wlm-si-dense-sparse-r1/main.go"
  ],
  "controls": [
    "Only substrate representation varies within each pair; seed, capacity, schedule, initial canonical occupied-slot set, transition order, and step count are identical.",
    "Dense and sparse transitions are implemented independently; neither implementation converts to the other for updates.",
    "Canonical observations are produced separately from each representation and include occupied cardinality, free cardinality, released count, accepted count, deficit, cumulative accepted, cumulative deficit, and canonical state digest.",
    "Sparse state must remain strictly increasing, unique, and in range; dense state length must equal capacity.",
    "Every scientific mismatch is recorded and the full frozen matrix continues; mismatches never trigger early stopping or parameter changes.",
    "An independent integer oracle checks accepted=min(demand,free-before), deficit=demand-accepted, post-occupied=pre-occupied-released+accepted, and post-occupied+post-free=capacity.",
    "Output ordering is fixed as seed, capacity, schedule, substrate, then step, using the preregistered list order.",
    "The result must retain exact zero and nonzero counts; failed hypotheses are preserved without threshold, seed, metric, budget, or interpretation changes."
  ],
  "fixed_parameters": {
    "allocation_rule": "After release, accept min(demand,free-before) units and occupy the lowest-numbered free slots.",
    "budget_overrun_rows_definition": "Number of emitted rows beyond the frozen total of 32768.",
    "capacity_count": "4",
    "capacity_values": "8,16,32,64",
    "conservation_laws": "post-occupied=pre-occupied-released+accepted and post-occupied+post-free=capacity",
    "deficit_law": "accepted=min(demand,free-before) and deficit=demand-accepted",
    "expected_rows_per_substrate_case": "64",
    "expected_total_emitted_rows": "32768 = 512 substrate case runs * 64 steps",
    "initial_state": "Slot i is occupied iff PRF(seed,0,i,0x494E4954) mod 4 equals 0; both substrates receive the resulting canonical occupied-slot set.",
    "integer_semantics": "Unsigned 64-bit PRF arithmetic wraps modulo 2^64; state counts and laws use exact nonnegative integer arithmetic.",
    "invalid_state_schedule_definition": "A substrate case run counts once if its dense length differs from capacity or its sparse indices are non-increasing, duplicated, or out of range.",
    "law_violation_case_definition": "A paired case counts once if either substrate violates any frozen deficit or conservation equality at any step.",
    "max_conservation_error_definition": "Maximum absolute integer error across both frozen conservation equalities.",
    "max_deficit_law_error_definition": "Maximum absolute integer error across accepted-min(demand,free-before) and deficit-(demand-accepted).",
    "paired_case_count": "256 = 8 seeds * 4 capacities * 8 schedules",
    "paired_observable_mismatch_case_definition": "A paired case counts once if any of its 64 steps differs between substrates in any canonical observable or state digest.",
    "prf": "SplitMix64 applied to uint64(seed) XOR domain XOR (uint64(step) shifted left 32) XOR uint64(slot), with standard fixed SplitMix64 constants.",
    "release_rule": "At each step, release each occupied slot i iff PRF(seed,step,i,0x52454C53) mod 5 equals 0.",
    "result_interpretation_rule": "Apply only the preregistered metric thresholds and outcome meanings; no post-result tuning.",
    "row_count_mismatch_case_definition": "A paired case counts once unless each substrate emits exactly 64 rows.",
    "schedule_alternating": "demand=0 on even steps and capacity+1 on odd steps",
    "schedule_burst": "demand=capacity+2 when step mod 8 equals 0; otherwise floor(capacity/2)",
    "schedule_count": "8",
    "schedule_exact": "demand=capacity",
    "schedule_half": "demand=floor(capacity/2)",
    "schedule_order": "zero,half,exact,plus_one,alternating,ramp,burst,seeded",
    "schedule_plus_one": "demand=capacity+1",
    "schedule_ramp": "demand=step mod (capacity+2)",
    "schedule_seeded": "demand=PRF(seed,step,0,0x444D4E44) mod (capacity+3)",
    "schedule_zero": "demand=0",
    "seed_count": "8",
    "state_cardinality_mismatch_case_definition": "A paired case counts once if pre-occupied, post-occupied, or post-free differs between substrates at any step.",
    "state_digest": "FNV-1a-64 over ascending occupied slot identifiers encoded as two little-endian bytes each, beginning from the standard FNV-1a-64 offset basis.",
    "steps_per_substrate_case": "64",
    "substrate_case_run_count": "512 = 256 paired cases * 2 substrates",
    "substrate_count": "2",
    "substrates": "dense-bool,sparse-sorted-u16",
    "transition_order": "observe pre-state; release; compute free-before; obtain demand; allocate; compute deficit; observe post-state; update cumulative counters"
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
      "name": "state_cardinality_mismatch_cases",
      "comparator": "==",
      "threshold": 0
    },
    {
      "name": "row_count_mismatch_cases",
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
      "name": "invalid_state_schedules",
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
    104729,
    130363,
    155921,
    181081,
    205019,
    229447,
    253907,
    279967
  ],
  "compute_seconds": 1800,
  "stop_conditions": [
    "Stop before scientific interpretation if baseline SHA, North Star SHA-256, evidence ID, harness ID, lane, or authority generation does not exactly match the sealed request.",
    "Stop and report an incomplete operational result if elapsed compute time reaches 1800 seconds.",
    "Stop and report an incomplete operational result if the harness cannot write the complete append-only result artifact or qualification metadata.",
    "Do not stop for scientific mismatches, law violations, or negative intermediate observations; execute the entire frozen matrix and preserve them.",
    "Do not retry with changed parameters, seeds, thresholds, metrics, budgets, code, or interpretations after any result is observed."
  ],
  "positive_meaning": "All 512 substrate case runs and 32768 rows complete with every preregistered metric meeting its threshold, supporting dense-versus-sparse invariance only for the frozen matrix.",
  "negative_meaning": "The run is complete and valid and at least one paired observable mismatch or frozen-law violation occurs. This falsifies the corresponding invariance claim; preserve all exact nonzero results without rerunning or tuning.",
  "mixed_meaning": "The run is complete and valid, but results split: either canonical substrate observables mismatch while all deficit and conservation laws remain exact, or substrate observables match while at least one frozen law fails. Preserve every exact count and classify neither full invariance nor full falsification beyond the affected claim."
}
```
