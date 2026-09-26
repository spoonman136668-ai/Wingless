# Wingless UP-109C — slot reuse confirmation detector

Status: preregistered scientific bounded-memory phase-free checkpoint experiment.

Scientific parent: sealed UP-108C 05e006e769ccb16b6f8ff78eac02edd29ae8dc77.

## Question

UP-108C showed that a resident-write mask can trigger the correct phase-free checkpoint when all final residents are refreshed in place, but misses orders where final residents enter by successful repeat admission. Does marking a slot reuse-confirmed on either in-place write or successful seen-before admission close that detector gap?

## Frozen memory/checkpoint action

Identical to UP-108C:
- exact recall cap 16;
- two-bit-aging exact memory;
- 1024-bit two-hash seen-once filter;
- generation interval 32;
- checkpoint action = clear filter and set counter to 8;
- no explicit phase label;
- no future oracle.

## New detector evidence

Maintain one reuse-confirmed bit per exact-memory slot.

Set slot bit when:
1. an existing resident key is written again; or
2. an unseen key is successfully admitted because seen-before evidence allows it.

Clear slot bit whenever that slot is replaced by another key.

When all currently occupied slots are reuse-confirmed, emit the checkpoint action once and clear the reuse-confirmed mask.

No reads contribute detector evidence.

## Frozen workload

Orders:
- ascending
- evens_then_odds
- odds_then_evens
- rotate8

Promotion:
- exact two-pass promotion sequence from UP-107C/UP-108C.

Churn:
- 12,288 unique one-shot writes;
- ordinary hot-query cadence unchanged;
- 64 episodes per cell;
- seeds 199M and 200M.

## Controls

1. explicit_post_promotion_checkpoint
2. resident_write_mask control from UP-108C
3. slot_reuse_confirmation

## Metrics

- checkpoint count and first checkpoint index;
- resident hot before churn;
- false-positive admissions;
- final hot accuracy;
- final hot-set exact accuracy;
- metadata bytes.

## Interpretation

If slot reuse confirmation matches explicit checkpoint behavior across all four orders, phase-local checkpointing has been converted into a purely local, phase-free mechanism.

## Bounds

No query-derived detector evidence, no future labels, no adaptive threshold, no capacity increase, no checkpoint-action change, no result-informed retry, no live activation, no production authority.
