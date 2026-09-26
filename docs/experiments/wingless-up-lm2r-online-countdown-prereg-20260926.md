# Wingless UP-LM2R — online countdown

Status: preregistered scientific fixed-capacity shadow-forecast experiment.

Scientific parent: sealed UP-LM2Q 9a7ac068b8a685dbfd44b087771e3cf5ce46fd28.

## Question

Does the current-state countdown remain exact when evaluated after every event in a mixed stream of STORE pressure and dependency-closing REPORT events, and how often does a near-term warning legitimately retract after a later REPORT closes risk?

## Frozen arms

Use:
- exact recall cap 16;
- deferred levels D=4,5,6;
- identity rotations 0 and 7;
- value shifts 0,1,2,3.

Initialize the same 12 first-chunk memories and reported history as LM2Q.

## Frozen event stream

Apply 16 events:
1 STORE, 2 REPORT-close pending #1,
3 STORE, 4 STORE, 5 REPORT-close pending #2,
6 STORE, 7 STORE, 8 REPORT-close pending #3,
9 STORE, 10 STORE, 11 STORE, 12 REPORT-close pending #4,
13 STORE, 14 STORE, 15 STORE, 16 STORE.

REPORT-close targets are frozen from the initial pending set; no target is chosen adaptively.

## Online predictor

After every event, use only current recall occupancy/order and past REPORT history:
free slots + FIFO index of earliest currently-unreported resident + 1.
If no unreported resident remains, emit -1.

Ground truth is recomputed from a clone of the current state by applying unique STOREs only until first unreported eviction.

Freeze "near-term warning" at countdown 1..4 unique STOREs. A warning retraction is recorded when a later REPORT moves the countdown above 4 or to -1.

## Measurements

Across all arms and all 16 event snapshots:
- exact countdown prediction rate;
- near-term warning snapshots;
- warning retractions caused by REPORT closure;
- per-event predicted/actual countdown.

## Bounds

Shadow only. No intervention, no future schedule inspection by the predictor, no extra memory, no adaptive rule, no result-informed retry, no live activation, or production authority.
