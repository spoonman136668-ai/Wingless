# Wingless UP-143C — live mixed-age replacement

Status: preregistered scientific bounded-memory consolidation experiment.

Scientific parent: sealed UP-142C 33d3928831081f204fecfdc28c58d2d213b85b87.

## Question

UP-142C proved age-first replacement in a controlled full-table fixture. Does the same rule hold in a live admission stream when two real qualified histories of different ages coexist and a single replacement event is forced?

## Frozen mechanism

Reuse the accepted age-eviction machine unchanged:
- exact recall cap 16;
- current/history tables 32/32;
- 128-byte admission sidecar;
- maxAge = 2;
- maximum-history-age replacement;
- lowest-index tie break;
- no semantic/query/future priority.

## Live schedule

Start from the accepted UP-140C setup, which leaves an `older` target as qualified history age 0.

Window 1 — 32 unresolved calls:
- present the `younger` target twice;
- present 15 fresh distractor identities twice each.

At boundary 1:
- older target becomes age 1;
- younger + 15 distractors enter history at age 0;
- history occupancy = 17.

Window 2 — 32 unresolved calls:
- present 16 new distractor identities twice each.

At boundary 2:
- older target becomes age 2;
- younger and first-window distractors become age 1;
- 16 new qualified histories attempt insertion;
- exactly one replacement is required.

Two mirrored target assignments:
- older=100, younger=103;
- older=103, younger=100.

Seeds:
- 271000000;
- 272000000.

32 episodes per seed.

After boundary 2, clone the machine state:
- test two sightings of the older target on one clone;
- test two sightings of the younger target on the other clone.

## Measurements

Per assignment:
- older/younger history presence and age after boundary 1;
- older/younger history presence and age after boundary 2;
- history occupancy after each boundary;
- age0/age1/age2 eviction counts;
- older re-admission rate;
- younger admission rate;
- hot exact-memory accuracy;
- unexpected manipulation admissions;
- panic rate.

## Interpretation

If the older age-2 live trace is uniquely evicted while the younger age-1 trace survives and still admits, the controlled age-rank rule generalizes to a live bounded stream.

## Bounds

No direct history-table fixture, no memory increase, no replacement-rule change, no semantic priority, no query priority, no future oracle, no adaptive selection, no result-informed retry, no live activation, no production authority.
