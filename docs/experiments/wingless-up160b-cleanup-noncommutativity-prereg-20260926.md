# Wingless UP-160B — cleanup-order noncommutativity

Status: preregistered scientific lexical-stability mechanism diagnostic.

Scientific parent: sealed UP-159B c1e458477a74c553edbe5f28f8281b9575beced7.

## Question

UP-159B established that the same fixed cleanup composition prefers different class orders across subjects. Does that context sensitivity arise from genuinely noncommutative chained updates, where swapping the same cleanup blocks produces different gate states from the identical post-REPORT state?

## Frozen state construction

For each accepted new-family subject:
- use its frozen matched two-subject corpus from UP-159B;
- reconstruct the accepted 15-epoch common prefix;
- execute terminal epoch 16 only up to and including the fixed four-example target-subject REPORT block;
- use exactly 20 non-tail new-family updates and old rehearsal indices 10,11,12 before REPORT;
- freeze the resulting post-REPORT gate.

No cleanup update is retained between diagnostic arms.

## Frozen cleanup composition

The same 12 old-rehearsal cleanup items are used in every arm:
- STORE indices 0..4;
- OBSERVE indices 5..9;
- REPORT indices 13..14.

Three matched order-swap comparisons:

1. STORE_OBSERVE_REPORT vs OBSERVE_STORE_REPORT.
2. STORE_REPORT_OBSERVE vs REPORT_STORE_OBSERVE.
3. OBSERVE_REPORT_STORE vs REPORT_OBSERVE_STORE.

Each sequence uses all 12 cleanup items exactly once.

## Measurements

Per subject × swap pair:
- old retention after sequence A;
- old retention after sequence B;
- signed retention difference A-B;
- gate-parameter Euclidean distance between the two final states;
- gate displacement norm from the common post-REPORT state for each sequence.

## Interpretation

A nonzero final-state distance proves chained cleanup updates are noncommutative at that subject's frozen state. Subject-specific sign changes in retention difference would provide a direct mechanism for the context-dependent order preferences seen in UP-159B.

No threshold is tuned after observation.

## Bounds

Diagnostic only. No adaptive scheduling, no extra training epochs, no memory change, no projector recomputation, no architecture change, no result-informed retry, no live activation, no production authority.
