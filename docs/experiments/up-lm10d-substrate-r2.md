# UP-LM10D-SUBSTRATE-R2

Controller-frozen preregistration:

```json
{
  "experiment_id": "UP-LM10D-SUBSTRATE-R2",
  "candidate_id": "C1-array-packed-trajectory-invariance",
  "harness_id": "wingless-windows-research",
  "changed_paths": [
    ".wingless/qualification-request.json",
    "scripts/run-up-lm10d-substrate-r2.ps1",
    "docs/experiments/up-lm10d-substrate-r2.md",
    "unitary/up_lm10d_substrate_r2.go",
    "cmd/up-lm10d-substrate-r2/main.go"
  ],
  "controls": [
    "Use exactly the eight preregistered seeds; do not add, remove, replace, or rerun seeds based on results.",
    "For each seed, enumerate each initial state 0 through 511 exactly once under each substrate; seed-dependent ordering is a control and does not alter membership.",
    "Generate one 64-transition schedule per seed and apply that identical schedule to both substrates and all 512 initial states in that seed block.",
    "Make storage substrate the sole experimental factor: substrate A is a nine-element Boolean array and substrate B is a uint16 restricted to bits 0 through 8.",
    "Run one process with GOMAXPROCS=1 and the frozen memory and wall-clock limits.",
    "Compute the canonical reference transition from the logical pre-state and operation before invoking either substrate implementation.",
    "Canonicalize every observed state to a nine-bit integer after every transition and hash the initial state plus all 64 post-transition canonical states.",
    "Stratify and preserve exact counts by seed and substrate; do not use pooled success to override a failing seed.",
    "Do not stop on a scientific mismatch and do not modify thresholds, seeds, schedules, metrics, resources, or interpretation after execution begins.",
    "Preserve complete, partial, negative, mixed, and invalid results verbatim with their qualification and integrity status."
  ],
  "fixed_parameters": {
    "baseline_sha": "2bd5dfb4f05ff4b5f90aa339c9c7668548a7bdf6",
    "cpu_parallelism": "1 process; GOMAXPROCS=1",
    "declared_write_sets": "toggle={i}; copy={j}; swap={i,j}",
    "emitted_row_definition": "one row per seed, initial state, and substrate containing final observables and per-case error maxima",
    "enumeration_order": "Fisher-Yates permutation of 0..511 using a separately initialized SplitMix64 stream and index selection value mod current_length",
    "evidence_reference": "RSE-0f75f4c81c65629752ba912c9899c790",
    "expected_emitted_rows": "8 seeds * 512 states * 2 substrates = 8192",
    "expected_paired_comparisons": "8 seeds * 512 states = 4096",
    "expected_total_transition_applications": "8 seeds * 512 states * 2 substrates * 64 transitions = 524288",
    "expected_unique_states_per_seed_per_substrate": "512",
    "final_observables": "canonical nine-bit final state, final popcount, and trajectory digest",
    "harness": "wingless-windows-research",
    "initial_state_domain": "all integers 0..511 inclusive",
    "initial_states_per_seed": "512",
    "invalid_permutation_definition": "a seed order that is not a bijection over integers 0..511",
    "memory_limit_mib": "256",
    "mismatch_row_definition": "a seed-and-initial-state pair whose substrates differ in final canonical state, final popcount, or trajectory digest",
    "operation_set": "toggle(i); copy(i,j); swap(i,j)",
    "order_stream_initialization": "uint64(seed) XOR 0xD1B54A32D192ED03",
    "packed_high_bit_violation_definition": "any packed-substrate post-state with a nonzero bit in positions 9..15",
    "post_target_error_definition": "Hamming distance between the substrate's canonical post-state and the independently computed logical reference post-state",
    "reference_semantics": "toggle flips cell i; copy assigns cell j the pre-operation value of cell i; swap exchanges the pre-operation values of cells i and j",
    "result_policy": "sealed exact result; no post-result tuning or selective reruns",
    "schedule_generation": "For each seed, a separately initialized SplitMix64 stream emits three uint64 values per transition: i=value1 mod 9, j=value2 mod 9, operation=value3 mod 3 where 0=toggle, 1=copy, 2=swap.",
    "schedule_stream_initialization": "uint64(seed) XOR 0x9E3779B97F4A7C15",
    "seed_count": "8",
    "state_width_bits": "9",
    "substrate_a": "[9]bool with index 0 as canonical bit 0",
    "substrate_b": "uint16 with canonical cells in bits 0..8 and bits 9..15 required zero",
    "substrate_count": "2",
    "toolchain_control": "Use only the exact Go toolchain pinned by successful harness qualification; if qualification does not identify an exact version, do not execute.",
    "trajectory_digest": "SHA-256 over 65 canonical states in order, each encoded as exactly two little-endian bytes with bits 9..15 zero: initial state followed by every post-transition state",
    "transitions_per_seed_state_substrate": "64",
    "wall_clock_limit_seconds": "1800",
    "write_error_definition": "count of cells outside the operation's declared write set whose value changed during that operation"
  },
  "metrics": [
    {
      "name": "emitted_rows",
      "comparator": "==",
      "threshold": 8192
    },
    {
      "name": "paired_comparisons",
      "comparator": "==",
      "threshold": 4096
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
      "name": "mismatch_rows",
      "comparator": "==",
      "threshold": 0
    },
    {
      "name": "max_post_target_error",
      "comparator": "==",
      "threshold": 0
    },
    {
      "name": "max_write_error",
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
    231701,
    257053,
    282407
  ],
  "compute_seconds": 1800,
  "stop_conditions": [
    "Do not begin unless baseline SHA, envelope ID, lane, harness, allowed paths, and pinned toolchain all match the qualification request.",
    "After execution begins, continue through all preregistered cases despite observed scientific mismatches.",
    "Stop successfully only after exactly 8192 rows and 4096 paired comparisons have been emitted and integrity checks have completed.",
    "Stop at 1800 wall-clock seconds, preserve the exact partial output, and classify it as incomplete rather than positive, negative, or mixed.",
    "Stop and classify the run as invalid if the process exceeds 256 MiB, qualification changes, a non-preregistered seed or parameter is detected, or output persistence fails; preserve all output produced before the stop."
  ],
  "positive_meaning": "Qualification succeeds, all 8192 rows and 4096 pairs are present with exactly 512 unique states per seed per substrate, and every preregistered metric meets its exact threshold in every seed block. This supports invariance only for the two specified substrates, state domain, schedules, observables, and resource bounds.",
  "negative_meaning": "All integrity cardinalities pass, and every seed block contains at least one paired mismatch, post-target error, unintended-write error, or packed high-bit violation. This uniformly falsifies exact substrate invariance; the exact failing rows and per-seed counts are retained.",
  "mixed_meaning": "All integrity cardinalities pass, and scientific failures occur in a strict nonempty subset of the eight seed blocks while at least one seed block satisfies every scientific threshold. This schedule-dependent outcome still falsifies the universal hypothesis and must be preserved without pooled reinterpretation."
}
```
