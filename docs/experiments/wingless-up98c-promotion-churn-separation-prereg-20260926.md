# Wingless UP-98C — promotion/churn filter-role separation

Status: preregistered scientific bounded-memory diagnostic experiment.

Scientific parent: sealed UP-97C a52f15992f6e0476f2f36eb88d222822694a9772.

## Question

UP-97C showed that the single seen-once filter is serving two competing roles: preserving legitimate second-write promotion evidence and rejecting later one-shot churn. If promotion is completed first and the filter is then cleared experimentally, can churn-generation phase be changed without sacrificing the promoted hot set?

## Frozen memory and workload

- exact recall cap 16;
- two-bit-aging replacement;
- 1024-bit legacy two-hash seen-once filter;
- generation interval 32 filtered unseen writes;
- 32 original keys;
- 16 hot even keys;
- each hot key receives one second write before churn;
- 12288 unique one-shot churn writes;
- hot queries after every four hot second writes and every four churn writes;
- value vocabulary 32;
- 64 episodes per arm;
- seeds 185M and 186M.

## Frozen arms

1. continuous_control
   - exact UP-97C phase-offset-0 behavior with no phase clear.

2. clear_after_promotion_phase0
3. clear_after_promotion_phase8
4. clear_after_promotion_phase16
5. clear_after_promotion_phase24

For every clear_after_promotion arm:
- process initial writes and all hot second writes using the ordinary mechanism;
- immediately after the complete hot second-write phase, clear only the seen-once filter;
- leave exact memory and two-bit-aging state untouched;
- set the generation counter to the named phase offset;
- then process churn normally with 32-write generations.

The post-promotion clear is an explicit diagnostic phase label and is not proposed as a general online policy.

## Metrics

For every arm:
- hot query accuracy;
- hot-set exact accuracy;
- cold query accuracy;
- recall entries used;
- rejected one-shot churn writes;
- false-positive churn admissions;
- first false-positive churn index;
- episodes with any false admission;
- filter reset count;
- policy metadata bytes;
- total bounded-memory bytes.

## Interpretation

If a cleared/filter-phase arm preserves the promoted hot set and removes false admissions, the C97 tradeoff is specifically caused by sharing one filter state across promotion and churn roles. If not, within-churn hash geometry remains limiting.

## Bounds

No exact-memory expansion, no filter-width change, no hash-count change, no generation-interval tuning, no adaptive phase, no query-derived labels, no future oracle, no result-informed retry, no live activation, no production authority.
