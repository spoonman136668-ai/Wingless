# WLM-DLG-DEFICIT-SWEEP-R1

Controller-frozen preregistration:

```json
{
  "experiment_id": "WLM-DLG-DEFICIT-SWEEP-R1",
  "candidate_id": "WLM-DLG-DEFICIT-SWEEP-R1",
  "harness_id": "wingless-windows-research",
  "changed_paths": [
    ".wingless/qualification-request.json",
    "scripts/run-wlm-dlg-deficit-sweep-r1.ps1",
    "docs/experiments/wlm-dlg-deficit-sweep-r1.md",
    "unitary/wlm_dlg_deficit_sweep_r1.go",
    "cmd/wlm-dlg-deficit-sweep-r1/main.go"
  ],
  "controls": [
    "Use only baseline commit 10b30d53eb39e6e847ecf7f43aafad9ccffc0fe3 and verified evidence RSE-4a1d1ee42c9e362e89c0ef4415c63daa.",
    "Resolve the same two-member substrate pair used by the referenced baseline experiment; fail qualification if it cannot be resolved unambiguously.",
    "Change only the resource-deficit level within the experimental sweep; keep substrate pair, seeds, schedules, step count, observable schema, arithmetic, and comparison procedure frozen.",
    "Generate each seed's eight schedules once, hash their canonical encodings, and reuse identical encodings for every deficit level and both substrates.",
    "Use exact integer arithmetic; no floating-point tolerance, regression fitting, threshold adjustment, seed replacement, rerun selection, or post-result interpretation change.",
    "Keep bounded-state overflow as an explicit observable and require stepwise conservation so saturation cannot hide work.",
    "Run cases in a deterministic preregistered order and retain every emitted row, including failing and incomplete rows.",
    "Preserve and report negative, mixed, timeout, and qualification-failure outcomes without tuning or suppressing cases."
  ],
  "fixed_parameters": {
    "baseline_sha": "10b30d53eb39e6e847ecf7f43aafad9ccffc0fe3",
    "budget_definition": "max(0, requested_work - deficit_level)",
    "comparison_mode": "canonical row-for-row exact integer equality",
    "completed_substrate_case_runs": "256 paired_cases * 2 substrates_per_case = 512",
    "completed_work_law": "min(requested_work, budget)",
    "conservation_law": "prior_state + admitted_work = next_state + completed_work + overflow",
    "deficit_law": "requested_work - completed_work = max(0, requested_work - budget)",
    "deficit_levels": "[0,1,2,4]",
    "evidence_id": "RSE-4a1d1ee42c9e362e89c0ef4415c63daa",
    "evidence_result_sha256": "bf662c5e938d2fbf171a3b6dead61dd2baf920031d08b8cfbd9134eeaf54770c",
    "paired_case_count": "8 seeds * 8 schedules_per_seed * 4 deficit_levels = 256",
    "post_result_tuning": "forbidden",
    "requested_work_domain": "nonnegative integers generated solely by the frozen schedule generator",
    "result_tolerance": "0",
    "run_order": "seed ascending, schedule index ascending, deficit level ascending, substrate index ascending",
    "schedules_per_seed": "8",
    "seed_count": "8",
    "steps_per_run": "64",
    "substrate_pair": "exact two-member pair resolved from experiment UP-LM16-PAGED-SUBSTRATE-INVARIANCE-R1; qualification fails if resolution is not unique",
    "substrates_per_case": "2",
    "total_emitted_rows": "512 runs * 64 steps_per_run = 32768"
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
    "Stop before scientific execution if the baseline commit, North Star hash, evidence hash, harness identity, or uniquely resolved evidence-linked substrate pair does not match the preregistration; preserve the qualification failure.",
    "Stop and preserve partial artifacts if elapsed compute reaches 3600 seconds; classify as incomplete, not positive or negative.",
    "Stop and preserve partial artifacts on harness crash, non-integer observable, malformed row, duplicate case identity, or inability to write the sealed local result artifact.",
    "Do not stop early because a scientific metric fails; complete all remaining cases unless a preceding integrity or compute-limit condition applies."
  ],
  "positive_meaning": "Qualification succeeds; exactly 256 paired cases, 512 substrate-case runs, and 32768 rows complete; every preregistered metric equals its threshold; and the exact deficit and conservation laws and paired bounded-state observables hold at all four deficit levels.",
  "negative_meaning": "Qualification succeeds and the run completes, but one or more preregistered exact acceptance metrics fail. This is valid evidence against the joint generalization claim and must be preserved exactly.",
  "mixed_meaning": "All qualification and cardinality checks complete, but at least one preregistered exact metric fails while another remains satisfied; or the deficit law holds on both substrates while paired observables differ, or vice versa. Record the complete result as mixed without changing any parameter or interpretation."
}
```
