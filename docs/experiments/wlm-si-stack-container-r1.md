# WLM-SI-STACK-CONTAINER-R1

Controller-frozen preregistration:

```json
{
  "experiment_id": "WLM-SI-STACK-CONTAINER-R1",
  "candidate_id": "WLM-C1-STACK-CONTAINER",
  "harness_id": "wingless-windows-research",
  "question": "Under fixed capacity, inputs, initial state, pop budget, and transition order, are exact bounded-state observables and resource-response laws invariant between a slice-backed LIFO stack and an independently implemented linked-node LIFO stack?",
  "hypothesis": "Every preregistered paired row and aggregate is exactly equal across the two stack-storage substrates, and both substrates satisfy the registered conservation and deficit laws on the complete frozen matrix.",
  "exact_parent_sha": "7287347f636c7c2aa7c1dc3b804d06a96dd28778",
  "north_star_path": "research/bootstrap/wingless-lm-north-star.txt",
  "north_star_sha256": "efd6172772d76f03ac04bdab0af40941232311def8202d6f0f953d33b3d05dd8",
  "changed_paths": [
    ".wingless/qualification-request.json",
    "scripts/run-wlm-si-stack-container-r1.ps1",
    "docs/experiments/wlm-si-stack-container-r1.md",
    "unitary/wlm_si_stack_container_r1.go",
    "cmd/wlm-si-stack-container-r1/main.go"
  ],
  "arms": [
    "control: slice-backed LIFO stack with append and truncation",
    "treatment: singly linked LIFO stack with independently allocated nodes"
  ],
  "controls": [
    "Change only the LIFO storage substrate within each pair; capacity, immutable push trace, initial ordered token state, pop budget, transition order, arithmetic, and observations are identical.",
    "Generate each 64-step push-count trace once before either arm executes and replay that same immutable integer sequence in both arms.",
    "Implement transitions independently; neither arm converts to or invokes the other arm for push, pop, length, or snapshot operations.",
    "Evaluate an independent scalar oracle from the registered recurrences; neither substrate supplies oracle values.",
    "Compare complete ordered bottom-to-top remaining-token sequences and complete ordered pop sequences directly, without hashing, normalization, averaging, or tolerance.",
    "Use exact nonnegative integer arithmetic for all registered laws and counters.",
    "Balance execution order deterministically: slice first for even seed indices and linked first for odd seed indices; sort emitted rows canonically before paired comparison.",
    "Continue the complete valid matrix after any scientific mismatch or law violation so exact negative-result cardinality is retained.",
    "Use no network, broker, credentials, live runtime, accepted-reference mutation, cross-lane input, scheduler, queue consumer, daemon, retry loop, or path outside the authorized research checkout.",
    "Thresholds, seeds, schedules, budgets, metrics, capacity, code, and interpretations are frozen before output and may not be tuned after any output is observed."
  ],
  "fixed_parameters": {
    "experimental_dimension": "LIFO stack storage substrate only: slice-backed versus linked-node",
    "logical_capacity": 16,
    "steps_per_run": 64,
    "transition_order": "At step t, admit push tokens in ascending ordinal order up to remaining capacity, reject the remainder, then pop from the LIFO top up to the fixed pop budget.",
    "push_token_ids": "At step t, push ordinal j in [0,p_t-1] has token ID t*16+j.",
    "initial_occupancies": [0, 5, 10, 15],
    "initial_token_order": "For initial occupancy s0, bottom-to-top token IDs are -s0 through -1; the empty case has no initial tokens.",
    "pop_budgets": [1, 2, 4, 8],
    "push_schedule_order": ["constant", "alternating", "pulse", "lcg32"],
    "constant_schedule": "p_t=1+(seed mod 8) for every t.",
    "alternating_schedule": "p_t=16 when ((t+seed) mod 2)=1, otherwise 0.",
    "pulse_schedule": "p_t=16 when ((t+seed) mod 8)=0, otherwise 0.",
    "lcg32_schedule": "x_0=seed; x_(t+1)=(1664525*x_t+1013904223) mod 2^32; p_t=x_(t+1) mod 17.",
    "registered_transition_laws": "admitted_t=min(p_t,16-s_t); rejected_t=p_t-admitted_t; popped_t=min(s_t+admitted_t,b); s_(t+1)=s_t+admitted_t-popped_t.",
    "registered_conservation_law": "s_(t+1)-(s_t+admitted_t-popped_t)=0 exactly on every row.",
    "registered_deficit_law": "deficit_t=max(0,s_t+admitted_t-b), and s_(t+1)=deficit_t exactly on every row.",
    "row_observables": "seed,schedule,initial_occupancy,pop_budget,substrate,step,pushes,admitted,rejected,popped,s_t,s_(t+1),deficit,cumulative_admitted,cumulative_rejected,cumulative_popped,complete bottom-to-top remaining-token IDs,complete pop-order token IDs",
    "pair_key": "(seed,push_schedule,initial_occupancy,pop_budget)",
    "paired_case_count": "8 seeds * 4 schedules * 4 initial occupancies * 4 pop budgets = 512",
    "completed_substrate_case_runs": "512 pairs * 2 substrates = 1024",
    "total_emitted_rows": "1024 runs * 64 steps = 65536",
    "process_resources": "one harness process, GOMAXPROCS=1, no parallel case execution",
    "comparison_rule": "exact integer and exact ordered-sequence equality; tolerance 0",
    "run_order": "seed order, schedule order, initial occupancy order, pop budget order; arm order balanced by seed index",
    "compute_limit_seconds": 3600,
    "post_result_tuning": "forbidden"
  },
  "seeds": [131, 151, 173, 197, 223, 251, 281, 313],
  "budgets": {
    "logical_capacity_tokens": 16,
    "pop_budgets_tokens_per_step": [1, 2, 4, 8],
    "compute_seconds": 3600,
    "max_emitted_rows": 65536
  },
  "metrics": [
    {"name": "paired_case_count", "comparator": "==", "threshold": 512},
    {"name": "completed_substrate_case_runs", "comparator": "==", "threshold": 1024},
    {"name": "total_emitted_rows", "comparator": "==", "threshold": 65536},
    {"name": "paired_observable_mismatch_cases", "comparator": "==", "threshold": 0},
    {"name": "row_count_mismatch_cases", "comparator": "==", "threshold": 0},
    {"name": "state_cardinality_mismatch_cases", "comparator": "==", "threshold": 0},
    {"name": "invalid_state_schedules", "comparator": "==", "threshold": 0},
    {"name": "budget_overrun_rows", "comparator": "==", "threshold": 0},
    {"name": "law_violation_cases", "comparator": "==", "threshold": 0},
    {"name": "max_conservation_error", "comparator": "==", "threshold": 0},
    {"name": "max_deficit_law_error", "comparator": "==", "threshold": 0}
  ],
  "validity_criteria": [
    "The checked-out parent before the preregistration commit is exactly 7287347f636c7c2aa7c1dc3b804d06a96dd28778 and the North Star digest matches the frozen SHA-256.",
    "All four push traces for every seed contain exactly 64 integers in [0,16] and are byte-for-byte identical between paired arms.",
    "Every initial occupancy and every intermediate stack length is in [0,16], every snapshot length equals the observed stack length, and every pop count is at most its frozen budget.",
    "The full matrix completes within the compute cap and emits exactly the registered pair, run, and row cardinalities.",
    "The implementation, command, harness, and request remain confined to the five declared changed paths and the authorized research-only scope."
  ],
  "classification_rules": {
    "supported": "The run is valid and complete and all eleven metric thresholds pass; support is limited to this frozen matrix.",
    "negative": "The run is valid and complete and at least one paired-observable or registered-law threshold fails without the split pattern defined as mixed.",
    "mixed": "The run is valid and complete and substrate equality passes while a registered law fails, or both laws pass while paired substrate observables differ.",
    "null": "The run is valid and complete, all integrity and cardinality metrics pass, but neither the supported, negative, nor mixed rule can be applied because the registered scientific contrasts contain no evaluable paired rows; this classification is expected to be unreachable under the frozen cardinalities and must still be preserved if encountered.",
    "incomplete": "A preregistered compute, environment, artifact, or integrity stop occurs before a valid complete matrix; no scientific conclusion is made.",
    "invalid_implementation": "After any output is observed, implementation behavior is found not to match this preregistration; restore every implementation and request change, retain only the preregistration commit, do not rerun, and make no scientific claim."
  },
  "stop_conditions": [
    "Before scientific execution, stop if parent identity, North Star digest, branch, harness identity, declared paths, seeds, dimensions, or frozen formulas do not match this preregistration.",
    "Stop and preserve an incomplete result if elapsed compute reaches 3600 seconds, compilation or harness execution fails, or the complete result artifact cannot be emitted.",
    "Stop and preserve an integrity failure if a generated trace is malformed, differs between arms, or any state is out of bounds; do not replace its seed or schedule.",
    "Do not stop early for a scientific mismatch or law violation; complete all remaining valid cases.",
    "After any output is observed, do not change or rerun the implementation under this preregistration; apply the invalid-implementation disposition exactly if a defect is discovered."
  ],
  "no_post_result_tuning_rule": "No threshold, metric, seed, schedule, budget, capacity, arm, code path, validity criterion, classification rule, stop condition, or interpretation may be changed after any probe or test output is observed. Any later scientific attempt requires a new experiment ID and a new preregistration."
}
```
