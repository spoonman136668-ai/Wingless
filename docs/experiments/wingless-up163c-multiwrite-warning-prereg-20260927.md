# Wingless UP-163C — multi-write warning horizon

Status: preregistered scientific shadow-monitoring experiment.

Scientific parent: sealed UP-162C 181b0bb3dace65aaad51d3926627c70142439bd1.

## Question

Can sparse monitoring recover high loss-warning recall if each check forecasts across the whole interval until the next check rather than asking only whether the immediately next write is harmful?

## Frozen starting states and pressure

Reuse the same 16 UP-161C post-guided-refresh states and the same adversarial non-endangered shielding policy used by UP-161C/162C.

The pressure policy is deterministic from current state:
1. predict next replacement slot;
2. if that slot is not the endangered item, refresh that resident item;
3. issue one unique write.

## Frozen monitoring arms

- cadence 2 with a 2-write lookahead;
- cadence 4 with a 4-write lookahead.

Checks occur on writes 1, 1+cadence, 1+2*cadence, ...

At each check, copy the current memory state and simulate the already-frozen pressure policy for at most the lookahead length. Warn iff the endangered item would be lost within that simulated interval.

The copy is discarded. Warnings do not alter real memory or pressure.

## Measurements

Per cadence:
- eventual losses;
- losses with a warning active before occurrence;
- missed losses;
- false warnings;
- warning recall.

## Interpretation

This tests whether a bounded current-state forecast can trade computation for lower sampling frequency in a known pressure regime. It does not establish robustness to unknown future policies.

## Bounds

Shadow only. No corrective action, adaptive cadence, future external schedule, semantic priority, capacity increase, extra training, live activation, or result-informed retry.
