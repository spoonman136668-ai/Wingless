# Wingless UP-140C — unresolved candidate identity density

Status: preregistered scientific bounded-memory temporal-consolidation experiment.

Scientific parent: sealed UP-139C d0b747eada76bb3ca82b140b4888e033620e3a6e.

## Question

UP-139C established that the consolidation clock advances with unresolved candidate events, not generic process calls. With the unresolved-event count fixed, does history survival also depend on how those events are distributed across candidate identities and within-window recurrence?

## Frozen mechanism

Reuse the accepted age-eviction machine unchanged:
- exact recall cap 16;
- current/history tables 32/32;
- 128-byte admission sidecar;
- maxAge = 2;
- lowest-index tie break;
- no semantic/query/future priority.

Each arm uses one recurrence target only.

## Frozen setup

As in UP-139C:
- initialize exact recall and bounded history;
- allow original target history to expire;
- rebuild fresh target history with a G3 same-window pair;
- finish the G3 boundary so target history begins the manipulation at age 0.

## Fixed unresolved-event budget

Every arm then executes exactly 64 unresolved candidate process calls: 32 calls in each of two complete candidate-counter windows.

No manipulation key is present in history before its window, and window 2 always uses a fresh disjoint key set, so none can complete two-stage durable admission during the manipulation.

Identity-density arms, per 32-call window:

1. unique_32
   - 32 distinct keys × 1 sighting.

2. paired_16
   - 16 distinct keys × 2 sightings each.

3. quartet_8
   - 8 distinct keys × 4 sightings each.

4. repeat_1
   - 1 key × 32 sightings.

Thus every arm has exactly 64 unresolved events and exactly two candidate boundaries; only identity count and current-window qualification density differ.

Then present the recurrence target twice.

Two target identities:
- 100
- 103.

Seeds:
- 267000000;
- 268000000.

32 episodes per seed.

## Measurements

Per target × arm:
- unresolved calls;
- unique manipulation identities;
- target-history presence/age after boundary 1;
- target-history presence/age after boundary 2;
- history occupancy after each boundary;
- age0/age1/age2 evictions during the manipulation;
- unexpected manipulation admissions;
- target admission rate and final accuracy;
- hot accuracy;
- max current/history occupancy;
- recall entries used;
- panic rate.

## Interpretation

If all four identity patterns leave the target at the same age and admission outcome, clock pressure depends on unresolved event count alone. If denser within-window recurrence changes target survival despite identical event count, history insertion/replacement pressure adds an identity-topology effect on top of the event clock.

No post-result threshold is introduced.

## Bounds

No memory increase, no replacement-rule change, no semantic labels, no query priority, no future oracle, no adaptive horizon, no result-informed retry, no live activation, no production authority.
