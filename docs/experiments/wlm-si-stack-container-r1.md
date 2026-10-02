# WLM-SI-STACK-CONTAINER-R1

Controller-frozen preregistration:

```json
{
  "experiment_id": "WLM-SI-STACK-CONTAINER-R1",
  "candidate_id": "WLM-C1-STACK-CONTAINER",
  "harness_id": "wingless-windows-research",
  "question": "Under fixed bounded resources and identical push/pop traces, do complete ordered state and popped-token observables plus registered integer laws remain exactly invariant between an independently implemented slice-backed LIFO stack and fixed-capacity array LIFO stack?",
  "hypothesis": "For the complete frozen matrix, the two LIFO storage substrates emit exactly equal canonical row observables and satisfy the registered admission, rejection, pop, conservation, and residual-demand laws on every row.",
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
    "slice-backed-lifo",
    "fixed-capacity-array-lifo"
  ],
  "controls": [
    "Change only the LIFO storage substrate within each pair; semantic rule, logical capacity, input traces, initial ordered state, pop budget, seeds, step count, arithmetic, and observations are identical.",
    "Implement each substrate independently; neither transition implementation converts to or calls the other.",
    "Generate each complete push trace before either paired arm executes and reuse its exact integer sequence for both arms.",
    "Evaluate registered integer laws from scalar pre-state and frozen inputs independently of either substrate's snapshot method.",
    "Compare complete ordered bottom-to-top stack state and complete ordered popped-token sequences without hashes, tolerance, normalization, averaging, or omitted fields.",
    "Balance execution order deterministically by seed index, then sort observations into canonical seed, schedule, initial occupancy, pop budget, substrate, step order.",
    "Run the full valid matrix after scientific mismatches or law violations; preserve every failing stratum and exact negative-result cardinality.",
    "Use one process with GOMAXPROCS=1 and no parallel case execution, network, broker, credentials, accepted-reference mutation, live runtime, cross-lane input, scheduler, queue consumer, daemon, or retry loop.",
    "Do not change thresholds, seeds, inputs, metrics, budgets, capacities, code, or interpretation after observing any scientific output."
  ],
  "fixed_parameters": {
    "experimental_dimension": "LIFO storage representation only: dynamically sized slice versus fixed-capacity array with explicit length.",
    "logical_capacity": "16 tokens for every run and both arms.",
    "steps_per_run": "64 transitions indexed t=0..63.",
    "push_schedule_families": "[constant,alternating,pulse,lcg32], each producing integer requested pushes in [0,16].",
    "constant_schedule": "p_t=1+(seed mod 8) for every t.",
    "alternating_schedule": "p_t=16 when ((t+seed) mod 2)=0, otherwise 0.",
    "pulse_schedule": "p_t=16 when ((t+seed) mod 8)=1, otherwise 0.",
    "lcg32_schedule": "x_0=seed; x_(t+1)=(1664525*x_t+1013904223) mod 2^32; p_t=x_(t+1) mod 17.",
    "initial_occupancies": "[0,5,10,15].",
    "initial_token_order": "For initial occupancy s0, bottom-to-top token IDs are the ordered integers -s0 through -1; empty has no tokens.",
    "pop_budgets": "[1,2,4,8] tokens per step, constant within a run.",
    "push_token_ids": "At step t, requested push ordinal j in [0,p_t-1] has token ID t*16+j and requests are considered in ascending j order.",
    "transition_order": "Observe pre-state; admit requested pushes in ascending ordinal order until capacity; reject the remainder; pop from the LIFO top up to the fixed pop budget; observe post-state.",
    "registered_transition_laws": "admitted_t=min(p_t,16-s_t); rejected_t=p_t-admitted_t; popped_t=min(s_t+admitted_t,b); s_(t+1)=s_t+admitted_t-popped_t.",
    "registered_conservation_law": "s_(t+1)-(s_t+admitted_t-popped_t)=0 exactly on every row.",
    "registered_residual_law": "residual_t=max(0,s_t+admitted_t-b), and s_(t+1)=residual_t exactly on every row.",
    "row_observables": "seed, schedule family, initial occupancy, pop budget, substrate, step, requested pushes, admitted count, rejected count, popped count, s_t, s_(t+1), residual, cumulative admitted/rejected/popped counts, complete bottom-to-top token IDs, complete pop-order token IDs.",
    "pair_key": "(seed,push_schedule_family,initial_occupancy,pop_budget).",
    "paired_case_count": "8 seeds * 4 schedules * 4 initial occupancies * 4 pop budgets = 512.",
    "completed_substrate_case_runs": "512 pairs * 2 arms = 1024.",
    "total_emitted_rows": "1024 runs * 64 rows = 65536.",
    "execution_order": "For even seed indices run slice then array; for odd seed indices run array then slice; canonical result rows are ordered by seed, schedule list order, initial occupancy, pop budget, substrate list order, then step.",
    "process_resources": "One harness process, GOMAXPROCS=1, no parallel case execution, 3600-second wall-clock limit.",
    "comparison_rule": "Exact integer and exact ordered-sequence equality on every row; zero tolerance.",
    "result_policy": "One complete sealed result for the frozen matrix; no post-result tuning, case replacement, selective rerun, or interpretation change."
  },
  "seeds": [
    131,
    149,
    167,
    191,
    223,
    251,
    277,
    313
  ],
  "seed_control": "All eight seeds are fixed before implementation and disjoint from WLM-SI-QUEUE-CONTAINER-R1 seeds [11,23,37,53,71,89,107,127].",
  "compute_seconds": 3600,
  "metrics": [
    {"name":"paired_case_count","comparator":"==","threshold":512},
    {"name":"completed_substrate_case_runs","comparator":"==","threshold":1024},
    {"name":"total_emitted_rows","comparator":"==","threshold":65536},
    {"name":"paired_observable_mismatch_cases","comparator":"==","threshold":0},
    {"name":"row_count_mismatch_cases","comparator":"==","threshold":0},
    {"name":"state_cardinality_mismatch_cases","comparator":"==","threshold":0},
    {"name":"invalid_state_schedules","comparator":"==","threshold":0},
    {"name":"budget_overrun_rows","comparator":"==","threshold":0},
    {"name":"law_violation_cases","comparator":"==","threshold":0},
    {"name":"max_conservation_error","comparator":"==","threshold":0},
    {"name":"max_residual_law_error","comparator":"==","threshold":0}
  ],
  "validity_criteria": [
    "The checked-out history has exact parent 7287347f636c7c2aa7c1dc3b804d06a96dd28778 followed first by this preregistration-only commit.",
    "The North Star bytes hash to efd6172772d76f03ac04bdab0af40941232311def8202d6f0f953d33b3d05dd8.",
    "The compiled implementation and harness exactly encode every frozen arm, input, seed, budget, capacity, metric, law, ordering rule, cardinality, and stop condition.",
    "Every one of 512 pairs, 1024 arm runs, and 65536 rows completes and is retained; each input trace is generated before its pair executes and is identical across arms.",
    "Each emitted state is within [0,16], snapshot length equals state cardinality, and all token sequences are retained without hashing.",
    "Any implementation mismatch discovered after output is observed invalidates the implementation: restore every implementation/request change, retain only this preregistration commit, make no scientific claim, and do not repair or rerun under this preregistration."
  ],
  "classification_rules": {
    "supported": "The run is valid and complete and all eleven preregistered metrics exactly meet their thresholds; support is limited to the frozen matrix.",
    "negative": "The run is valid and complete and one or more invariance and/or registered-law metrics fail in every applicable preregistered stratum; preserve exact counts and make only the corresponding falsification claim.",
    "mixed": "The run is valid and complete but invariance or law outcomes differ across preregistered strata, or substrate equality and registered-law validity disagree; preserve every stratum and make no global claim.",
    "null": "The run is valid and complete but the registered observables do not distinguish the hypothesis from its frozen controls sufficiently to support or falsify it; preserve the exact result without reinterpretation.",
    "incomplete": "An environmental or compute stop prevents the frozen full cardinality; preserve diagnostics but make no scientific claim.",
    "invalid": "Any provenance, preregistration, implementation, input identity, ordering, cardinality, or artifact-integrity criterion fails; make no scientific claim and do not repair or rerun after output under this preregistration."
  },
  "stop_conditions": [
    "Before scientific execution, stop as invalid if parent history, branch, North Star hash, preregistration identity, implementation conformance, or qualification-request identity differs from the frozen values.",
    "Stop as incomplete at 3600 wall-clock seconds or if compiler, harness, process, or artifact writer failure prevents the full matrix from completing and being preserved.",
    "Stop immediately as invalid if any path requests network, broker, credential, accepted-reference, live-runtime, cross-lane, scheduler, queue-consumer, daemon, or retry-loop authority.",
    "Do not stop on scientific mismatch or law violation; complete the full valid matrix to preserve exact negative-result cardinality.",
    "After any scientific output is observed, do not change or repair code, inputs, seeds, thresholds, metrics, budgets, capacity, or interpretation under this preregistration."
  ],
  "no_post_result_tuning_rule": "No threshold, seed, schedule, metric, budget, capacity, arm, code path, case count, classification, or interpretation may be added, removed, weakened, repaired, or retuned after any scientific output is observed."
}
```
