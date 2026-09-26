# Wingless UP-155C — protected-cohort prewrite warning

Status: preregistered scientific native-state warning experiment.

Scientific parent: sealed UP-154C 6dcfd9ecaf4d58f989943517a14cb545f409b086.

## Question

Can the exact prewrite replacement forecast identify loss of a specific protected cohort, rather than merely any original durable fact, without semantic priority?

## Frozen protected cohorts

Use four interleaved cohorts over the 16 original durable keys:
- cohort_A = {0,4,8,12}
- cohort_B = {1,5,9,13}
- cohort_C = {2,6,10,14}
- cohort_D = {3,7,11,15}

Each cohort spans multiple lexical classes and positions.

## Frozen stream

Reuse exact recall cap 16 and starting hands 0,4,8,12.

Run 24 writes per arm:
- steps divisible by 3 repeat the most recently introduced novel key;
- all other steps introduce a new unique key.

Use a disjoint sparse-refresh schedule:
- step 5 -> key 2
- step 10 -> key 6
- step 15 -> key 10
- step 20 -> key 14

If a refresh target is absent, the query misses with no substitute.

## Frozen warning

For each cohort independently, before every write:
1. if incoming key already exists, warn false;
2. otherwise predict next replacement slot from hand + age vector;
3. warn iff the predicted slot currently contains a member of that cohort.

No lexical class or future outcome is used.

## Measurements

Per cohort and globally:
- TP/FP/FN/TN;
- precision;
- recall;
- accuracy.

## Bounds

No intervention, semantic priority, adaptive refresh, capacity increase, policy change, future oracle, result-informed retry, live activation, or production authority.
