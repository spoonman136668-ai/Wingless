# UP-LM9D-SUBSTRATE-R1

Controller-frozen preregistration:

```json
{
  "experiment_id": "UP-LM9D-SUBSTRATE-R1",
  "candidate_id": "UP-LM9D-SUBSTRATE-R1",
  "harness_id": "wingless-windows-research",
  "changed_paths": [
    ".wingless/qualification-request.json",
    "scripts/run-up-lm9d-substrate-r1.ps1",
    "docs/experiments/up-lm9d-substrate-r1.md",
    "unitary/up_lm9d_substrate_r1.go",
    "cmd/up-lm9d-substrate-r1/main.go"
  ],
  "controls": [
    "Verify the checked-out baseline SHA equals ccc800672f456ba56bc89d00d8332641e60d8544 before starting.",
    "Use only the wingless-windows-research harness and the five preregistered changed paths.",
    "Hold dimension, complete state domain, observable definitions, permutations, execution process, toolchain selection, and resource limits constant; change only the state representation substrate.",
    "Enumerate every integer state from 0 through 511 exactly once per seed and feed identical state and permutation values to both substrate paths.",
    "Generate each permutation once from its frozen seed, validate that it is a bijection on indices 0 through 8, then reuse it unchanged for both substrates.",
    "Keep the array and packed-bitmask encode, mutation, and decode implementations independent; share only immutable inputs, row schema, and final comparator.",
    "Compare integer observables exactly with no tolerance, rounding, filtering, exclusions, or discarded rows.",
    "Emit results in deterministic seed-then-state-then-substrate order and preserve the complete result even when a scientific threshold fails.",
    "Do not alter thresholds, seeds, metrics, dimensions, budgets, interpretation, or code after observing any result; any revision requires a new experiment ID and preregistration.",
    "Do not activate live runtime behavior, access brokers or credentials, mutate accepted references, or cross lanes."
  ],
  "fixed_parameters": {
    "baseline_sha": "ccc800672f456ba56bc89d00d8332641e60d8544",
    "comparison_rule": "exact integer equality after canonical decoding",
    "dimension": "9",
    "experimental_factor": "state representation substrate only",
    "go_toolchain": "repository-pinned toolchain resolved at the baseline SHA",
    "gomaxprocs": "1",
    "harness": "wingless-windows-research",
    "memory_budget_mib": "512",
    "observable_names": "post_target,write",
    "paired_comparisons_per_seed": "512",
    "permutations_per_seed": "1",
    "result_retention": "preserve complete positive, mixed, and negative results without overwrite",
    "rows_per_pair": "2",
    "seed_count": "8",
    "state_cardinality": "2^9=512",
    "state_enumeration_order": "ascending 0..511",
    "state_values": "all integers 0..511 inclusive",
    "substrate_a": "[9]uint8 with each cell restricted to 0 or 1",
    "substrate_b": "uint16 with bits 0..8 used and bits 9..15 required to remain zero",
    "total_emitted_rows": "2*8*512=8192",
    "total_paired_comparisons": "8*512=4096",
    "verified_evidence_id": "RSE-9824b96ed93f501bbd85e3740f13f0f4",
    "verified_evidence_scope": "UP-LM9D-permutation-invariance: max_post_target_error=0; max_write_error=0; mismatch_rows=0",
    "wall_clock_budget_seconds": "1800"
  },
  "metrics": [
    {
      "name": "mismatch_rows",
      "comparator": "==",
      "threshold": 0
    },
    {
      "name": "max_post_target_error",
      "comparator": "\u003c=",
      "threshold": 0
    },
    {
      "name": "max_write_error",
      "comparator": "\u003c=",
      "threshold": 0
    },
    {
      "name": "paired_comparisons",
      "comparator": "==",
      "threshold": 4096
    },
    {
      "name": "emitted_rows",
      "comparator": "==",
      "threshold": 8192
    },
    {
      "name": "unique_states_per_seed_per_substrate",
      "comparator": "==",
      "threshold": 512
    },
    {
      "name": "invalid_permutations",
      "comparator": "==",
      "threshold": 0
    },
    {
      "name": "packed_high_bit_violations",
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
    232003,
    257591,
    283163
  ],
  "compute_seconds": 1800,
  "stop_conditions": [
    "Before execution, refuse to start if the baseline SHA, harness ID, lane, authority generation, changed-path set, or preregistration digest differs from this proposal.",
    "Before measurement, refuse to start if any frozen seed fails to produce a bijection on indices 0 through 8.",
    "Stop and mark the run invalid if the 1800-second compute cap or 512 MiB memory budget is exceeded; preserve partial output and do not rerun under this experiment ID.",
    "Stop and mark the run invalid on probe-wrapper failure, malformed row schema, duplicate or missing state identity, arithmetic overflow, or any packed high-bit violation; preserve diagnostics and partial output.",
    "Do not stop early for a scientific mismatch; complete all 4096 comparisons so an exact negative result is preserved.",
    "Immediately stop without activation if any step requests live runtime launch, broker access, credential access, accepted-reference mutation, cross-lane authority, or a path outside the preregistered set.",
    "After any result is observable, perform no tuning or replacement run under this experiment ID; a changed design requires a new sealed preregistration."
  ],
  "positive_meaning": "All preregistered metrics pass simultaneously across all 4096 paired comparisons, supporting exact LM9D invariance between the two specified representation substrates under the frozen controls; no claim beyond these substrates, observables, seeds, dimension, harness, or resources is authorized.",
  "negative_meaning": "A run-valid result with any post-target or write mismatch falsifies substrate invariance. Preserve all rows and summary metrics exactly; do not change code, seeds, thresholds, resources, metrics, or interpretation after observation.",
  "mixed_meaning": "Run-valid evidence in which the cardinality, permutation, and substrate-integrity controls pass but only one of the two scientific observable families remains exact. It is not a positive result and must be preserved without tuning or rerun under this experiment ID."
}
```
