# Wingless UP-162B — cleanup geometry and order-sign attribution

Status: preregistered scientific lexical-stability mechanism diagnostic.

Scientific parent: sealed UP-161B 4678afd1a0338aa849a4185d7d02d90166cd7b10.
Mechanism parent: sealed UP-160B d4a12a5ecd2d5f24671d1af3a2a7267732f10fdd.

## Question

UP-160B proved cleanup updates are noncommutative and that the STORE↔OBSERVE retention sign changes by subject. UP-161B rejected simple fixed cycling as a universal solution. Can frozen post-REPORT state geometry explain the already-sealed subject-specific STORE↔OBSERVE sign without adapting the schedule?

## Frozen subject states

For each accepted new-family subject:
- use its frozen matched pair from UP-159B;
- reconstruct the accepted 15-epoch common prefix;
- execute terminal epoch 16 through the fixed target-subject REPORT interference block exactly as UP-160B;
- freeze the resulting post-REPORT gate.

## Frozen geometric diagnostics

From the identical frozen gate for each subject:

1. Aggregate one-step cleanup vectors
   - STORE: old rehearsal indices 0..4
   - OBSERVE: indices 5..9
   - REPORT: indices 13..14
   - each individual one-step vector is measured from a fresh clone of the same frozen gate, then summed; no chaining.

2. Old-rehearsal reference vector
   - all 15 accepted old-rehearsal items, each measured as a one-step vector from the same frozen gate and summed.

3. Pairwise block geometry
   - vector norms;
   - STORE/OBSERVE, STORE/REPORT, OBSERVE/REPORT cosines;
   - cosine of each block with the old-rehearsal reference.

4. STORE↔OBSERVE commutator geometry
   - execute full STORE_OBSERVE_REPORT and OBSERVE_STORE_REPORT cleanup orders on separate clones;
   - compute final-state difference vector SOR minus OSR;
   - measure its norm and cosine with the frozen old-rehearsal reference.

The sealed UP-160B retention delta SOR minus OSR is carried as a fixed parent label for each subject. It does not control the harness and is not used to choose or modify an order.

## Interpretation

The exact six-subject feature table is the result. Consistent relation between geometric direction/projection and the fixed parent retention sign would identify a state-space explanation for context-sensitive order preference. Failure to separate signs would show that these first/second-order diagnostics are insufficient.

No post-result classifier, threshold, feature selection, or schedule adaptation is allowed.

## Bounds

Diagnostic only. No retained diagnostic updates, no adaptive ordering, no extra training, no memory change, no projector recomputation, no architecture change, no result-informed retry, no live activation, no production authority.
