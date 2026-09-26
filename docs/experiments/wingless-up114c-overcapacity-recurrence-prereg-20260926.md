# Wingless UP-114C — exact probation over-capacity recurrence pressure

Status: preregistered scientific bounded-memory capacity experiment.

Scientific parent: sealed UP-113C f425d2410d413ee53086cf794efc397e66b3660f.

## Question

UP-113C kept a 12-hot + 4-recurring mixed working set exactly under long churn at the fixed 16-entry cap. What happens when legitimate recurring demand itself exceeds that fixed exact-memory capacity?

## Frozen mechanism

Exact UP-113C:
- exact recall cap 16;
- two-bit-aging replacement;
- exact32 probation table, 32 x uint32 = 128 bytes;
- generation interval 32;
- phase-free slot-reuse checkpoint;
- no query-derived admission evidence;
- no future oracle;
- policy metadata unchanged.

## Working-set pressure

Persistent hot set:
- 12 keys.

Legitimate recurring candidate counts:
- 4  => total target demand 16
- 5  => total target demand 17
- 8  => total target demand 20

Every recurring candidate:
- first appears once to probation;
- repeats after exactly 16 filtered writes;
- therefore qualifies for admission under the frozen exact32 recurrence rule.

## Churn

After all legitimate recurring candidates have been presented:
- 49,152 unique one-shot churn writes;
- query all 12 hot keys and every recurring candidate every four churn writes;
- 32 episodes per seed;
- seeds 215M and 216M.

## Metrics

- recurring admission rate;
- 12-hot accuracy;
- recurring-candidate accuracy;
- retained recurring count;
- total target-set accuracy;
- whole-target exact accuracy;
- one-shot false admissions;
- exact-memory occupancy;
- checkpoint count.

## Interpretation

The 17- and 20-demand cells intentionally exceed capacity. The purpose is not to force exactness; it is to map which information survives and whether replacement remains stable/selective without increasing the 16-entry cap.

## Bounds

No capacity increase, no adaptive admission, no priority labels, no future oracle, no phase labels, no threshold tuning, no result-informed retry, no live activation, no production authority.
