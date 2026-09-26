# Wingless UP-154C — prewrite durable-loss warning

Status: preregistered scientific native-state warning experiment.

Scientific parent: sealed UP-153C da9e256704ddafc9cb20c990f347e9bb59cf2181.

## Question

Can the demonstrated sufficient replacement state (current hand + 16 two-bit ages), combined only with current slot occupancy and whether the incoming key is already present, warn before a write will erase an original durable fact?

## Frozen design

Reuse the accepted up81cAging memory unchanged:
- exact recall cap 16;
- same 16 initial durable facts;
- starting hands 0, 4, 8, 12.

Run 24 writes per hand. The incoming schedule is deterministic:
- steps divisible by 3 repeat the most recently introduced novel key;
- every other step introduces a new unique key;
- therefore each arm contains 16 novel admissions and 8 repeat updates.

Before steps 4, 8, 12, 16, 20, and 24 issue fixed refresh queries to original keys:
3, 7, 11, 15, 0, 4 respectively. If already absent, the query misses and no substitute is chosen.

## Frozen warning rule

Before each write:
1. if the incoming key already exists, predict no replacement loss;
2. otherwise compute the next replacement slot using only current hand + age vector;
3. predict durable loss iff that slot currently contains one of the original 16 durable entries.

No semantic class, value, future schedule, or outcome is used.

## Measurements

Per operation:
- incoming key and whether already present;
- refresh hit/miss;
- predicted replacement slot;
- predicted durable-loss warning;
- actual durable loss;
- original durable count before/after.

Aggregate TP/FP/FN/TN, precision, recall, and warning accuracy.

## Bounds

No intervention, refresh adaptation, semantic priority, capacity increase, policy change, future oracle, result-informed retry, live activation, or production authority.
