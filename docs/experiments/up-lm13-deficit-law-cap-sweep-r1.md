# UP-LM13-DEFICIT-LAW-CAP-SWEEP-R1

Controller-frozen preregistration:

```json
{
  "experiment_id": "UP-LM13-DEFICIT-LAW-CAP-SWEEP-R1",
  "candidate_id": "C1-CAP-SWEEP",
  "harness_id": "wingless-windows-research",
  "changed_paths": [
    ".wingless/qualification-request.json",
    "scripts/run-up-lm13-deficit-law-cap-sweep-r1.ps1",
    "docs/experiments/up-lm13-deficit-law-cap-sweep-r1.md",
    "unitary/up_lm13_deficit_law_cap_sweep_r1.go",
    "cmd/up-lm13-deficit-law-cap-sweep-r1/main.go"
  ],
  "controls": [
    "Verify the baseline SHA, North Star SHA-256, evidence ID, lane, and harness identity before any case runs.",
    "Use only the eight preregistered seeds and eight preregistered cap levels; do not add, remove, replace, or rerun selected cases after observing results.",
    "Change only the total acceptance cap between cases; freeze substrate implementation, demand count, state cardinality, operation generator, logical ordering rules, admission semantics, and metric definitions.",
    "Start every seed-cap case from fresh isolated state and execute cases in the fixed lexicographic order of seed then cap.",
    "Generate exactly 512 deterministic one-write demand events addressing 512 injective state identifiers in every case.",
    "Compute expected accepted_writes=cap and expected deficit=512-cap before reading case output; do not derive the oracle from observed counters.",
    "Use exact integer comparisons only; no statistical tolerance, threshold revision, metric substitution, or interpretation retuning is permitted.",
    "Write all case records, including failures and partial records, to sealed research output; negative results remain valid and must not be discarded.",
    "Do not access brokers, credentials, live runtimes, accepted refs, other lanes, or paths outside the authorized research package.",
    "Treat infrastructure failure, malformed output, or incomplete case coverage as invalid/inconclusive rather than as support for the hypothesis; preserve the associated records."
  ],
  "fixed_parameters": {
    "aggregate_demand_writes": "32768",
    "baseline_sha": "65fac4cd5d824c2bb6df92ce9dc58a08b4731067",
    "cap_level_count": "8",
    "cap_levels": "64,128,192,256,320,384,448,512",
    "case_order": "ascending seed, then ascending cap",
    "case_run_count": "64",
    "comparison_mode": "exact integer equality",
    "demand_writes_per_case": "512",
    "evidence_id": "RSE-41fb2d0299b8a9407a3f029732eafc22",
    "expected_aggregate_accepted_writes": "18432",
    "expected_aggregate_deficit": "14336",
    "harness": "wingless-windows-research",
    "lane": "LM",
    "north_star_sha256": "efd6172772d76f03ac04bdab0af40941232311def8202d6f0f953d33b3d05dd8",
    "oracle": "accepted_writes=cap; deficit=512-cap; accepted_writes+deficit=512",
    "result_policy": "preserve all positive, negative, mixed, partial, and invalid records without post-result tuning",
    "seed_count": "8",
    "substrate_count": "1",
    "substrate_selection": "frozen baseline scientific substrate",
    "sum_cap_levels_per_seed": "2304",
    "total_emitted_rows": "32768",
    "unique_states_per_case": "512",
    "writes_per_demand_event": "1"
  },
  "metrics": [
    {
      "name": "substrate_case_runs",
      "comparator": "==",
      "threshold": 64
    },
    {
      "name": "cap_level_count",
      "comparator": "==",
      "threshold": 8
    },
    {
      "name": "total_emitted_rows",
      "comparator": "==",
      "threshold": 32768
    },
    {
      "name": "aggregate_demand_writes",
      "comparator": "==",
      "threshold": 32768
    },
    {
      "name": "aggregate_accepted_writes",
      "comparator": "==",
      "threshold": 18432
    },
    {
      "name": "aggregate_deficit",
      "comparator": "==",
      "threshold": 14336
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
      "name": "max_conservation_error",
      "comparator": "==",
      "threshold": 0
    }
  ],
  "seeds": [
    104729,
    130363,
    155921,
    181081,
    205759,
    231481,
    257053,
    282571
  ],
  "compute_seconds": 3600,
  "stop_conditions": [
    "Stop before execution if baseline SHA, North Star SHA-256, evidence ID, lane, harness identity, or authorized paths do not match the preregistration.",
    "Stop and preserve all produced records if wall-clock compute reaches 3600 seconds; classify the run as invalid/inconclusive unless all 64 cases and required records already completed.",
    "Stop and preserve all produced records on harness crash, malformed or non-integer scientific output, duplicate or missing seed-cap identity, failure to create fresh isolated state, or unexpected filesystem mutation.",
    "Do not stop early for a scientific mismatch; complete all 64 cases so negative and mixed outcomes retain their preregistered meaning.",
    "Do not rerun, tune, extend, or substitute any case after results are observed."
  ],
  "positive_meaning": "All 64 preregistered seed-cap cases complete with exactly 512 rows and 512 unique states each; every case satisfies accepted_writes=cap, deficit=512-cap, exact conservation, and no budget overrun; all aggregate cardinalities and totals equal their preregistered values.",
  "negative_meaning": "With valid complete coverage, any accepted-write mismatch, deficit mismatch, conservation error, budget overrun, or workload-shape mismatch falsifies the universal cap-sweep hypothesis. The exact negative result is retained. An infrastructure, identity, timeout, malformed-output, or incomplete-coverage stop is invalid/inconclusive, not a scientific negative, and its records are also retained.",
  "mixed_meaning": "Integrity and coverage controls pass, but one or more exact law metrics fail only for a subset of preregistered caps or seeds. This rejects the universal hypothesis while providing boundary-localized evidence; the subset must be reported unchanged and must not trigger new thresholds, seeds, caps, or interpretations in this experiment."
}
```
