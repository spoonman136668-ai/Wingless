# Wingless UP-162C — periodic recheck coverage

Status: preregistered scientific shadow-monitoring experiment.

Scientific parent: sealed UP-161C b3ae72f2192ef21892278bbba86afe45b6439f04.

## Question

UP-161C showed that later refreshes can sharply shorten a rescued item's survival horizon. How frequently must Wingless recompute its exact native one-write loss warning to catch that renewed danger before loss?

## Frozen starting states and pressure

Reuse the 16 UP-161C post-guided-refresh starting states.

Apply the same adversarial non-endangered shielding process:
- before each unique pressure write, predict the next replacement slot from current hand+ages;
- if it is not the endangered item, refresh that resident item;
- then issue one unique write;
- no refresh of the endangered item.

This pressure process is identical for every monitoring cadence.

## Shadow monitoring cadences

Evaluate independent copies with warning checks:
- every write (cadence 1);
- writes 1,3,5,... (cadence 2);
- writes 1,5,9,... (cadence 4).

At a scheduled check, after the adversarial refresh and immediately before the unique write, recompute the exact next replacement slot from current hand+ages. Warn iff that write will evict the endangered item.

Warnings do not change memory or the write stream.

## Measurements

Per cadence:
- 16 eventual losses;
- losses warned immediately before occurrence;
- missed losses;
- false warnings;
- arm-level warning recall.

## Interpretation

This measures the monitoring-frequency requirement for a dynamic warning. It does not authorize automatic rescue.

## Bounds

Shadow only. No corrective action, adaptive cadence, future schedule input, semantic priority, capacity increase, extra training, live activation, or result-informed retry.
