# WLM-SI-QUEUE-CONTAINER-R1

Controller-frozen preregistration:

```json
{
  "experiment_id": "WLM-SI-QUEUE-CONTAINER-R1",
  "candidate_id": "WLM-C1-QUEUE-CONTAINER",
  "harness_id": "wingless-windows-research",
  "changed_paths": [
    ".wingless/qualification-request.json",
    "scripts/run-wlm-si-queue-container-r1.ps1",
    "docs/experiments/wlm-si-queue-container-r1.md",
    "unitary/wlm_si_queue_container_r1.go",
    "cmd/wlm-si-queue-container-r1/main.go"
  ],
  "controls": [
    "Change only the FIFO storage substrate within each pair: slice-backed FIFO versus fixed-capacity ring-buffer FIFO.",
    "Use one compiled scientific test binary for both substrates and the same harness process, logical capacity, inputs, initial state, and service budget within each pair.",
    "Generate every arrival trace before either member of its substrate pair is evaluated and reuse it byte-for-byte.",
    "Evaluate an independent scalar oracle from the registered integer recurrences; neither substrate adapter supplies oracle values.",
    "Compare complete ordered queued-token and serviced-token sequences without hashes, normalization, or tolerance.",
    "Emit rows in canonical order by seed, schedule, initial occupancy, service budget, substrate, and step; do not use execution timing as a metric.",
    "Run the entire valid matrix after a scientific mismatch; do not early-stop, alter thresholds, replace seeds, or add cases.",
    "Preserve the exact result, including every failing stratum and zero-valued metric, for positive, mixed, and negative outcomes.",
    "Use no network, broker, credentials, live runtime, accepted-reference mutation, cross-lane inputs, or paths outside the authorized roots."
  ],
  "fixed_parameters": {
    "alternating_schedule": "a_t=16 when ((t+seed) mod 2)=1, otherwise 0.",
    "arrival_schedule_families": "[constant,alternating,pulse,lcg32], each producing integer arrivals in [0,16].",
    "arrival_token_ids": "At step t, arrival ordinal j in [0,a_t-1] has token ID t*16+j; token order is ascending j.",
    "baseline_commit": "93a3da439a2ec6f8b5eef4b6193b26a11351294b",
    "comparison_rule": "Exact integer and exact ordered-sequence equality only; no tolerance, hashing, averaging, normalization, or dropped fields.",
    "constant_schedule": "a_t=1+(seed mod 8) for every t.",
    "evidence_scope": "RSE-68705d58285e3f62c4963deadb9dbe59 reports 128 paired cases, 1024 completed substrate case runs, 65536 rows, and zero reported mismatches, invalid schedules, budget overruns, law violations, conservation error, and deficit-law error for WLM-SI-DENSE-SPARSE-R2; no broader conclusion is assumed.",
    "execution_order": "Canonical case order is seed, schedule family, initial occupancy, then service budget; evaluate slice first for even seed indices and ring first for odd seed indices, then sort emitted rows canonically before comparison.",
    "experimental_dimension": "FIFO storage substrate only: slice-backed FIFO versus fixed-capacity ring-buffer FIFO.",
    "frozen_dimensions": "8 seeds * 4 arrival schedules * 4 initial occupancies * 4 service budgets = 512 paired cases; 512 pairs * 2 substrates = 1024 substrate case runs; 1024 runs * 64 rows = 65536 emitted rows.",
    "initial_occupancies": "[0,5,10,15].",
    "initial_token_order": "For initial occupancy q0, queue token IDs are the ordered integers -q0 through -1; the empty case has no initial tokens.",
    "lcg32_schedule": "x_0=seed; x_(t+1)=(1664525*x_t+1013904223) mod 2^32; a_t=x_(t+1) mod 17.",
    "logical_capacity": "16 tokens for every run and both substrates.",
    "pair_key": "(seed,arrival_schedule_family,initial_occupancy,service_budget); the two substrates form one pair.",
    "process_resources": "One harness process with GOMAXPROCS=1; no parallel case execution; logical resources are identical within each pair.",
    "pulse_schedule": "a_t=16 when ((t+seed) mod 8)=0, otherwise 0.",
    "registered_conservation_law": "q_(t+1)-(q_t+admitted_t-served_t)=0 exactly on every row.",
    "registered_deficit_law": "deficit_t=max(0,q_t+admitted_t-b), and q_(t+1) must equal deficit_t exactly.",
    "registered_transition_laws": "admitted_t=min(a_t,16-q_t); rejected_t=a_t-admitted_t; served_t=min(q_t+admitted_t,b); q_(t+1)=q_t+admitted_t-served_t.",
    "result_policy": "Write one exact sealed result covering the full valid matrix; retain all negative and mixed strata and perform no post-result tuning or rerun with altered parameters.",
    "row_observables": "seed, schedule family, initial occupancy, service budget, substrate, step, arrivals, admitted count, rejected count, served count, q_t, q_(t+1), deficit, cumulative admitted/rejected/served counts, complete ordered queued-token IDs, and complete ordered serviced-token IDs.",
    "service_budgets": "[1,2,4,8] tokens per step, constant within a run.",
    "steps_per_run": "64 transitions, indexed t=0..63.",
    "transition_order": "At each step, admit arrivals in token order up to remaining capacity, reject the remainder, then service from the FIFO head up to the fixed service budget."
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
      "name": "budget_overrun_rows",
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
    }
  ],
  "seeds": [
    11,
    23,
    37,
    53,
    71,
    89,
    107,
    127
  ],
  "compute_seconds": 3600,
  "stop_conditions": [
    "Before scientific evaluation, stop and record an integrity failure if the checked-out commit is not 93a3da439a2ec6f8b5eef4b6193b26a11351294b or any sealed package input differs from its qualification record.",
    "Stop and record an incomplete environmental run if the harness, compiler, test process, or result writer fails such that the full registered matrix cannot be evaluated or preserved.",
    "Stop at 3600 wall-clock seconds and preserve all emitted diagnostics as an incomplete run; do not interpret partial rows as a scientific result.",
    "Do not stop on an observable mismatch or law violation; complete all remaining valid cases so negative and mixed results retain the preregistered cardinalities.",
    "Do not rerun with changed seeds, dimensions, thresholds, metrics, budgets, schedules, code, or interpretations after observing any result."
  ],
  "positive_meaning": "All eleven registered metric thresholds pass on the complete 512-pair, 1024-run, 65536-row matrix. This supports only the preregistered slice-versus-ring invariance claim under the frozen bounds and does not establish universal substrate invariance.",
  "negative_meaning": "For a complete valid matrix, any registered metric misses its threshold, including any paired observable mismatch or nonzero conservation or deficit-law error. Preserve the exact negative result without tuning. An integrity or environmental abort is recorded separately as an incomplete run and carries no scientific conclusion.",
  "mixed_meaning": "The full valid matrix completes, but invariance or law criteria pass in some preregistered strata and fail in others. Preserve and report every stratum exactly; make no global invariance claim and do not change the matrix, thresholds, metrics, budget, seeds, or interpretation."
}
```
