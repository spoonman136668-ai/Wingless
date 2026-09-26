# Wingless UP-144C — live age-stack replacement order

Status: preregistered scientific bounded-memory replacement closeout.

Scientific parent: sealed UP-143C 35572f35815a5acbafd86b02fba9b5cae930e383.

## Question

UP-143C proved that one older live qualified trace is evicted before a younger trace under full-table pressure. Does the rule remain ordered across successive replacement events when multiple age cohorts coexist?

## Frozen mechanism

Reuse the accepted age-eviction machine unchanged:
- exact recall cap 16;
- current/history tables 32/32;
- 128-byte admission sidecar;
- maxAge = 2;
- maximum-history-age replacement;
- lowest-index tie break;
- no semantic/query/future priority.

## Live target cohorts

Two mirrored assignments:

A:
- oldest1=100;
- oldest2=101;
- younger1=102;
- younger2=103.

B:
- oldest1=102;
- oldest2=103;
- younger1=100;
- younger2=101.

Start from UP-140C setup using oldest1, leaving oldest1 as live history age 0.

### Window 1 — 32 unresolved calls

- oldest2 twice;
- 15 fresh distractors twice each.

At boundary:
- oldest1 becomes age 1;
- oldest2 + 15 distractors enter at age 0;
- occupancy = 17.

### Window 2 — 32 unresolved calls

- younger1 twice;
- younger2 twice;
- 14 fresh distractors twice each.

At boundary:
- oldest1 becomes age 2;
- oldest2 cohort becomes age 1;
- younger cohort enters age 0;
- table requires exactly one replacement;
- oldest1 is expected to be evicted.

### Window 3 — 32 unresolved calls

- one fresh candidate twice;
- 30 unique singletons.

At boundary:
- oldest2 cohort becomes age 2;
- younger cohort becomes age 1;
- one new qualified history forces exactly one replacement;
- oldest2 was inserted first in its cohort and is expected to be the lowest-index age-2 victim.

After boundary 3, clone the machine independently and test two sightings of all four targets.

Seeds:
- 273000000;
- 274000000.

32 episodes per seed.

## Measurements

Per assignment:
- target history presence/age after every boundary;
- history occupancy;
- age0/age1/age2 eviction counts;
- admission rate for each target after boundary 3;
- hot exact-memory accuracy;
- unexpected manipulation admissions;
- panic rate.

## Interpretation

Sequential loss of oldest1 then oldest2, with both younger targets surviving and admitting, demonstrates ordered age-stack replacement in a live stream. Failure defines the closeout boundary.

## Bounds

No direct table fixture, no memory increase, no rule change, no semantic/query priority, no future oracle, no adaptive selection, no result-informed retry, no live activation, no production authority.
