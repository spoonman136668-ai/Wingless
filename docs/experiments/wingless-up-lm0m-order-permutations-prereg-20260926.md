# Wingless UP-LM0M — clause-order permutation breadth

Status: preregistered scientific language/memory structural-generalization experiment.

Scientific parent: sealed UP-LM0L 2e75e78dc180e6b976572e63f6cb193552e4e6b4.

## Question

UP-LM0L remained exact on one interleaved clause order. Does learned prefix-timed routing remain stable across several deterministic dependency-valid orderings without retraining?

## Frozen models

Use the exact UP-LM0L model state:
- language model trained only on the original block-ordered corpus;
- three-way prefix classifier unchanged;
- recurrent state dimension 64;
- 512 recurrent-state bytes;
- exact recall cap 16;
- no attention;
- no future oracle;
- no training on the new orderings.

## Frozen held-out orderings

Use the same names, values, held-out split, and update-count schedule.

Evaluate four deterministic order families:

1. per_name
   - STORE, OBSERVE, optional update STORE, REPORT for each name.

2. paired_names
   - for names in pairs: STORE both, OBSERVE both, optional updates for the pair, REPORT both.

3. stores_then_local_reports
   - emit all initial STORE clauses first; then for each name emit OBSERVE, optional update STORE, REPORT.

4. reverse_report_tail
   - emit initial STORE, OBSERVE, and optional update STORE for all names; emit REPORT clauses in reverse name order.

All REPORT clauses occur after the latest STORE for that name.

## Frozen arms

1. explicit_event_routing
2. learned_prefix_threeway

Evaluate each order family at stream1 and stream4.

## Metrics

- top-1 byte accuracy;
- perplexity;
- dependent first-byte accuracy;
- whole-query-set exact accuracy;
- admission precision/recall;
- event-routing accuracy;
- report-routing accuracy;
- maximum recall entries.

## Interpretation

Learned routing matching explicit routing across all four orderings would support structural breadth of the routing mechanism. Shared degradation in generic byte metrics is recorded separately from memory-dependent correctness.

## Bounds

No retraining on these permutations, no state expansion, no memory-cap increase, no attention, no threshold tuning, no result-informed retry, no live activation, no production authority.
