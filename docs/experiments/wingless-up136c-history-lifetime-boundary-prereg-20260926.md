# Wingless UP-136C — rebuilt history lifetime boundary

Status: preregistered scientific bounded-memory temporal-consolidation experiment.

Scientific parent: sealed UP-135C 6dad48c4b5522a2991a600b33757522dc9bd2ea1.

## Question

UP-135C showed that after history expiry, a same-generation pair rebuilds qualified history and a later pair can complete a fresh two-stage admission. How many full generations can separate the rebuilding pair from the admitting pair before that rebuilt history trace expires?

## Frozen mechanism

Reuse UP-135C unchanged:
- exact recall cap 16;
- current/history tables 32/32;
- 128-byte admission sidecar;
- maximum-history-age replacement with maxAge=2;
- lowest-index tie break;
- no semantic/query/future priority.

Each arm contains exactly one recurrence target so every scheduled position advances the admission counter exactly once until admission.

## Frozen seven-generation schedules

Generation windows are 32 writes each. Every non-target position is a unique one-shot write.

All requalification arms first force original-history expiry with isolated sightings in G1 and G2, then rebuild fresh qualified history with a pair in G3 at positions 74,75.

Second-pair timing:
- next_g4: positions 106,107;
- after_one_blank_g5: positions 138,139;
- after_two_blank_g6: positions 170,171;
- after_three_blank_g7: positions 202,203.

Controls:
- rebuild_only: G1/G2 isolated sightings + G3 rebuilding pair, no second pair;
- live_pair_g1: positions 15,16 while original history remains live.

Two target identities are tested independently: 100 and 103.

Seeds:
- 259000000;
- 260000000.

32 episodes per seed.

## Interpretation

The exact delay response defines the lifetime of rebuilt qualified history. Admission at G4/G5/G6 but failure at G7 would match maxAge=2 semantics: age-0 history survives two complete intervening generations and is cleared at the next boundary.

No post-result threshold is introduced.

## Bounds

No memory increase, no replacement-rule change, no semantic labels, no query priority, no future oracle, no adaptive horizon, no result-informed retry, no live activation, no production authority.
