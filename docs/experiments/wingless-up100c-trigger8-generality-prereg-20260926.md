# Wingless UP-100C — trigger8 occupancy/order generality

Status: preregistered scientific bounded-memory generalization experiment.

Scientific parent: sealed UP-99C 816cc124614947ecdb2339f1f874b3c4cc88f355.

## Question

UP-99C showed that an online trigger after eight successful repeat admissions yields exact 16-hot retention through 12,288 one-shot churn writes. Is trigger8 robust to smaller hot-set occupancy and different deterministic initial key orders?

## Frozen memory mechanism

Identical to UP-99C trigger8:
- exact recall cap 16;
- two-bit-aging replacement;
- 1024-bit legacy two-hash seen-once filter;
- ordinary generation interval 32;
- one-shot special clear after the 8th successful repeat admission;
- post-trigger generation counter set to 8;
- no explicit phase label, query label, key identity oracle, or future information;
- 135 policy metadata bytes.

## Frozen workload grid

Hot-set sizes:
- 8
- 12
- 16

Hot keys are the first N even-numbered original keys.

Initial 32-key write orders:
1. ascending: 0..31
2. reverse: 31..0
3. evens_then_odds: all even keys ascending, then all odd keys ascending

After initial loading:
- each hot key receives exactly one second write in ascending key order;
- then 12,288 unique one-shot churn writes;
- query every hot key after every four second-hot writes and every four churn writes;
- value vocabulary 32;
- 64 episodes per grid cell;
- seeds 189M and 190M.

## Frozen arms

For each hot-size x initial-order cell:
1. continuous_control
2. trigger8

## Metrics

For every cell:
- hot query accuracy;
- hot-set exact accuracy;
- cold query accuracy;
- recall entries used;
- false-positive churn admissions;
- first false-positive churn index;
- trigger-fired episode count;
- mean filtered-write trigger index;
- policy metadata bytes;
- total bounded-memory bytes.

## Interpretation

Broad trigger8 success would support a genuine event-driven role-separation mechanism. Failures confined to cells where eight legitimate repeat admissions cannot occur would identify the trigger as occupancy-dependent rather than universally architectural.

## Bounds

No trigger-threshold tuning, no exact-memory expansion, no filter-width/hash/interval change, no adaptive initial-order logic, no query-derived labels, no future oracle, no result-informed retry, no live activation, no production authority.
