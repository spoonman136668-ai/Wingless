# Wingless UP-142C — age-rank replacement

Status: preregistered scientific bounded-history replacement diagnostic.

Scientific parent: sealed UP-141C 349da133458274f0d25cbde00135928ce4b79494.

## Question

UP-141C found the exact tested occupancy boundary: the age-2 target survives at 31/32 history entries and is evicted when the next qualified insertion forces replacement. Under identical full-table pressure, is survival determined by history age rank, with lowest table index used only as a tie-break?

## Frozen replacement mechanism

Use the accepted UP-125C history table unchanged:
- 32 history entries;
- maximum-history-age replacement;
- lowest-index tie break;
- maxAge = 2;
- no semantic/query/future priority;
- no memory increase.

## Controlled full-table fixture

For each arm:
- history table is exactly 32/32 occupied;
- one target key is present;
- 31 distractor keys are present;
- all distractors have age 0;
- perform exactly one `insertHistory(newKey, age=0)` call.

Target factors:
- target identity: 100 or 103;
- target age: 0, 1, or 2;
- target table index: 0 or 31.

This yields 12 deterministic arms.

The target-index factor explicitly calibrates the documented lowest-index tie-break:
- at target age 0, every entry has equal age;
- at target index 0, the target should be the tie victim;
- at target index 31, index-0 distractor should be the tie victim.

At target age 1 or 2, target is strictly older than all distractors, so age rank rather than index should determine replacement.

## Measurements

Per arm:
- target present before/after;
- incoming key present after;
- target age;
- target index;
- evicted key;
- evicted table index;
- age0/age1/age2 eviction counters;
- final history occupancy.

## Interpretation

Eviction of age-1/age-2 targets at either index, combined with the age-0 index-dependent tie result, demonstrates age-first replacement with index-only tie-breaking. Any deviation falsifies that description.

## Bounds

Diagnostic state fixture only. No replacement-rule change, no memory increase, no semantic priority, no query priority, no future oracle, no adaptive selection, no live activation, no production authority.
