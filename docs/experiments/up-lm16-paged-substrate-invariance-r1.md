# UP-LM16-PAGED-SUBSTRATE-INVARIANCE-R1

Controller-frozen preregistration:

```json
{
  "experiment_id": "UP-LM16-PAGED-SUBSTRATE-INVARIANCE-R1",
  "candidate_id": "C1-PAGED-SUBSTRATE-PAIR",
  "harness_id": "wingless-windows-research",
  "changed_paths": [
    ".wingless/qualification-request.json",
    "scripts/run-up-lm16-paged-substrate-invariance-r1.ps1",
    "docs/experiments/up-lm16-paged-substrate-invariance-r1.md",
    "unitary/up_lm16_paged_substrate_invariance_r1.go",
    "cmd/up-lm16-paged-substrate-invariance-r1/main.go"
  ],
  "controls": [
    "Verify Git HEAD equals the frozen baseline before generating any case or result.",
    "Generate each case manifest once and provide identical immutable inputs to both substrates.",
    "Change only state storage representation: dense-contiguous versus paged-sparse.",
    "Use identical logical state IDs, initial zero state, demand schedule, resource capacity, update law, tick horizon, and canonical row schema.",
    "Access logical states in ascending ID order on both substrates; never depend on map iteration order.",
    "Execute substrate runs serially and alternate first-run substrate by paired-case parity to expose order effects without changing inputs.",
    "Compare complete canonical rows field by field; hashes may be recorded but cannot replace direct equality.",
    "Evaluate completion, cardinality, row-count, schedule-validity, and conservation controls before interpreting substrate mismatches.",
    "Preserve partial, negative, mixed, and positive outputs exactly; perform no rerun with changed seeds, metrics, thresholds, budgets, or interpretation."
  ],
  "fixed_parameters": {
    "backlog_bound_per_state": "4095",
    "baseline_sha": "2dbc102971b7e8966515be390ab36e24bd8f9dbd",
    "candidate_substrate": "four sparse pages of four logical state IDs each; absent entries decode as zero",
    "canonical_integer_format": "unsigned base-10 integers with no omitted logical states",
    "canonical_row": "seed,case_index,substrate,tick,ordered_16_backlogs,total_demand,total_served,total_backlog,remaining_capacity,conservation_error",
    "completed_substrate_case_runs_derivation": "256 paired cases * 2 substrates = 512",
    "conservation_identity": "sum(next_backlog)=sum(prior_backlog)+sum(demand)-sum(served)",
    "demand_generator": "SHA-256 counter stream over UTF-8 experiment_id, decimal seed, decimal case_index, decimal tick, and decimal state_id with length-prefixed fields; demand is first digest byte modulo 4",
    "demand_range_per_state_tick": "0 through 3 inclusive",
    "evidence_id": "RSE-031736a38c704c5408953e6e8960a222",
    "evidence_sha256": "e9e00e81a9f849fbb7c45ca0b35793379dfbc40d6efa80f3795765a744842852",
    "execution_mode": "serial isolated research harness; no live runtime",
    "experimental_dimension": "state storage representation only",
    "initial_backlog_per_state": "0",
    "logical_state_cardinality": "16",
    "logical_state_ids": "integers 0 through 15 inclusive",
    "north_star_sha256": "efd6172772d76f03ac04bdab0af40941232311def8202d6f0f953d33b3d05dd8",
    "page_count": "4",
    "paired_case_count_derivation": "8 seeds * 32 schedules per seed = 256",
    "reference_substrate": "dense contiguous fixed array indexed by logical state ID",
    "resource_capacity_per_tick": "24 units",
    "result_policy": "single sealed run; preserve exact output; no post-result tuning",
    "schedules_per_seed": "32",
    "seed_count": "8",
    "service_order": "ascending logical state ID 0 through 15",
    "service_rule": "available_i=prior_backlog_i+demand_i; served_i=min(available_i,remaining_capacity); remaining_capacity decreases by served_i; next_backlog_i=available_i-served_i",
    "state_update_evaluations_derivation": "32768 rows * 16 logical states = 524288",
    "states_per_page": "4",
    "substrates_per_case": "2",
    "ticks_per_run": "64",
    "total_emitted_rows_derivation": "512 substrate-case runs * 64 ticks = 32768"
  },
  "metrics": [
    {
      "name": "paired_case_count",
      "comparator": "==",
      "threshold": 256
    },
    {
      "name": "completed_substrate_case_runs",
      "comparator": "==",
      "threshold": 512
    },
    {
      "name": "total_emitted_rows",
      "comparator": "==",
      "threshold": 32768
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
      "name": "paired_observable_mismatch_cases",
      "comparator": "==",
      "threshold": 0
    },
    {
      "name": "max_conservation_error",
      "comparator": "==",
      "threshold": 0
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
    }
  ],
  "seeds": [
    104729,
    130363,
    155921,
    181081,
    206369,
    231709,
    257053,
    282407
  ],
  "compute_seconds": 3600,
  "stop_conditions": [
    "Before execution, stop without running and record an invalid qualification outcome if Git HEAD is not 2dbc102971b7e8966515be390ab36e24bd8f9dbd.",
    "Stop before scientific interpretation if the generated manifest does not contain exactly 8 seeds, 32 cases per seed, 256 paired cases, or 512 substrate-case runs.",
    "Stop and preserve all partial artifacts if wall-clock compute reaches 3600 seconds.",
    "Stop and mark the run mixed or inconclusive if an input schedule differs between paired substrates, a logical state falls outside 0 through 15, or backlog would exceed 4095.",
    "Stop after exactly 512 substrate-case runs and refuse additional runs, replacement seeds, or selective reruns.",
    "Stop immediately if execution would require live activation, broker access, credentials, accepted-ref mutation, cross-lane authority, or a path outside the allowed roots."
  ],
  "positive_meaning": "All ten preregistered metric conditions hold, supporting exact dense-versus-paged substrate invariance only for the frozen 16-state, 64-tick, 256-pair domain; it does not establish universal invariance.",
  "negative_meaning": "If all completion and scientific-control gates pass but paired_observable_mismatch_cases is greater than zero, the exact negative result falsifies substrate invariance for the preregistered 256-case domain. Preserve every mismatching row without retuning or exclusion.",
  "mixed_meaning": "The result is mixed or inconclusive if any completion, schedule-validity, row-count, state-cardinality, conservation, law, or budget control fails, regardless of paired equality. Preserve the exact result and do not rerun with altered parameters."
}
```
