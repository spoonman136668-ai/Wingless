# UP-LM15-DEMAND-MULTIPLICITY-R1

Controller-frozen preregistration:

```json
{
  "experiment_id": "UP-LM15-DEMAND-MULTIPLICITY-R1",
  "candidate_id": "C1-DEMAND-MULTIPLICITY-LAW",
  "harness_id": "wingless-windows-research",
  "changed_paths": [
    ".wingless/qualification-request.json",
    "scripts/run-up-lm15-demand-multiplicity-r1.ps1",
    "docs/experiments/up-lm15-demand-multiplicity-r1.md",
    "unitary/up_lm15_demand_multiplicity_r1.go",
    "cmd/up-lm15-demand-multiplicity-r1/main.go"
  ],
  "controls": [
    "Use only baseline commit dbff168bf7a6c4a9356a5e7c850a4914b597f8f4 and evidence RSE-3b7536a64d307ad692206a2606eb2a20.",
    "Resolve and freeze the same two substrate implementations and the same resource-budget fixture used by the baseline experiment before any run.",
    "Generate exactly eight schedules per seed before applying demand multiplicity; reuse each generated schedule byte-for-byte for every factor and both substrates.",
    "Apply demand multiplicity only after schedule generation and apply it identically to both substrates.",
    "Keep row count, observation points, acceptance rule, deficit definition, state bounds, timeout, metrics, thresholds, seeds, and interpretation fixed.",
    "Randomize execution order once from the frozen seed list, record that order, and do not regenerate it after results are observed.",
    "Run every scientific case even after a threshold failure; threshold failures are recorded as negative results and do not trigger tuning.",
    "Perform no live activation, broker access, credential access, accepted-ref mutation, or cross-lane operation."
  ],
  "fixed_parameters": {
    "aggregate_demand_writes": "8 seeds * 8 schedules per seed * 2 substrates * 64 rows * (0+1+2+4) = 57344",
    "baseline_commit": "dbff168bf7a6c4a9356a5e7c850a4914b597f8f4",
    "demand_expansion_rule": "each baseline per-row demand event is replaced by factor identical ordered attempts; factor 0 emits the row with no demand attempt",
    "demand_multiplicity_factors": "0,1,2,4",
    "evidence_id": "RSE-3b7536a64d307ad692206a2606eb2a20",
    "frozen_schedule_instances": "64",
    "observable_vector": "accepted writes, deficit, state cardinality, and bounded-state digest at every emitted row",
    "paired_cases": "8 seeds * 8 schedules per seed * 4 factors = 256",
    "resource_budget": "unchanged baseline resource-budget fixture for every run",
    "rows_per_substrate_case": "64",
    "schedules_per_seed": "8",
    "scientific_failure_policy": "continue the complete frozen matrix and preserve exact output without retuning",
    "seed_count": "8",
    "substrate_case_runs": "8 seeds * 8 schedules per seed * 4 factors * 2 substrates = 512",
    "substrate_count": "2",
    "substrates": "the two substrate implementations resolved from the baseline experiment fixture at the baseline commit",
    "total_emitted_rows": "512 substrate-case runs * 64 rows = 32768"
  },
  "metrics": [
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
      "name": "paired_observable_mismatch_cases",
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
      "name": "aggregate_demand_writes",
      "comparator": "==",
      "threshold": 57344
    }
  ],
  "seeds": [
    11,
    29,
    47,
    71,
    101,
    131,
    173,
    211
  ],
  "compute_seconds": 7200,
  "stop_conditions": [
    "Stop before scientific execution if the baseline commit, evidence identifier, substrate fixture, resource-budget fixture, or generated matrix cannot be resolved exactly as preregistered.",
    "Stop on any attempted live activation, broker access, credential access, accepted-ref mutation, cross-lane access, or write outside changed_paths.",
    "Stop when wall-clock compute reaches 7200 seconds and report the package as incomplete and inconclusive with all partial outputs preserved.",
    "Do not stop for a scientific metric failure; finish the frozen matrix so the exact negative result is preserved."
  ],
  "positive_meaning": "All 512 substrate-case runs and 32768 rows complete, all cardinality checks equal their preregistered values, aggregate demand is exactly 57344, and every zero-threshold scientific metric remains zero. This supports only bounded generalization to factors 0, 1, 2, and 4 for the frozen schedules, substrates, and resource budget.",
  "negative_meaning": "A complete and integrity-valid matrix with any nonzero law violation, conservation error, budget overrun, or paired observable mismatch falsifies the corresponding hypothesis. Preserve all exact rows, aggregates, factor strata, and failing seeds without changing thresholds, seeds, factors, budget, metrics, or interpretation.",
  "mixed_meaning": "A complete, integrity-valid matrix in which the conservation metrics pass but paired observables fail supports deficit-law generalization while rejecting substrate invariance for at least one factor; the converse supports paired invariance while rejecting the generalized deficit law. Any incomplete or integrity-invalid matrix is inconclusive, not positive or negative."
}
```
