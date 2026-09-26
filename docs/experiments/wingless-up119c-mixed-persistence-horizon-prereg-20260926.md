# Wingless UP-119C — mixed cross-generation persistence horizon

Status: preregistered scientific bounded-memory persistence experiment.

Scientific parent: sealed UP-118C c04c533451f752810dce1cc718ae6428ced42eb3.

## Question

UP-118C showed that equal-memory age2 history admits candidates that recur after one skipped generation while rejecting one-generation bursts. Can the same mechanism handle adjacent and skipped persistence simultaneously, and where does its horizon fail?

## Frozen mechanisms

Compare:
1. previous_generation_control — exact UP-117C previous-generation crossgen_2x2.
2. age2_history — exact UP-118C age2 history.

Both:
- admission sidecar 128 bytes;
- exact recall cap 16;
- 12 established hot residents;
- two-bit-aging replacement;
- generation interval 32;
- phase-free slot-reuse checkpoint;
- no query-derived admission evidence;
- no semantic priority labels;
- no future oracle.

The age2 arm retains the preregistered 14-bit exact candidate-key encoding.

## Workload A — mixed_valid

Desired four recurring candidates:
- two adjacent candidates: 2 sightings in generation 1 and 2 sightings in generation 2;
- two skipped candidates: 2 sightings in generation 1, absent in generation 2, 2 sightings in generation 3.

Four burst distractors:
- 4 sightings in generation 1 only.

Target working set:
- 12 hot + 4 recurring = 16.

## Workload B — two_skip_boundary

Four target candidates:
- 2 sightings in generation 1;
- absent in generations 2 and 3;
- 2 sightings in generation 4.

Four burst distractors:
- 4 sightings in generation 1 only.

This cell intentionally exceeds the age2 temporal horizon without changing memory.

## Churn

After the final recurrence generation:
- 12,288 unique one-shot writes;
- query hot, target, and burst keys every four churn writes;
- 32 episodes per seed;
- seeds 225M and 226M.

## Metrics

Per arm x workload:
- target admission and retention;
- burst admission/retention;
- hot accuracy;
- target16 accuracy/exactness;
- one-shot false admissions;
- exact-memory occupancy;
- checkpoint count.

For mixed_valid also report adjacent-target and skipped-target accuracy separately.

## Interpretation

Success on mixed_valid would show age2 supports heterogeneous persistence timing in one working set. Failure on two_skip_boundary would establish the expected two-generation temporal horizon at fixed memory.

## Bounds

No admission-memory increase, no exact-memory capacity increase, no adaptive horizon, no semantic labels, no query-derived admission, no future oracle, no result-informed retry, no live activation, no production authority.
