# Wingless UP-165C — robust-warning transfer

Status: preregistered scientific shadow-monitoring experiment.

Scientific parent: sealed UP-164C e388b14467b099d0fce933570fc2d6c549d6d5c0.

## Question

Does the current-state robust horizon learned from the hostile shielding model remain useful when the *real* future pressure policy is different?

## Frozen warning model

At each sparse check, compute the UP-161C adversarial-shield horizon from current state only.

Warn for the next interval iff:
- horizon <= 2 for cadence 2;
- horizon <= 4 for cadence 4.

The warning model is identical across all real pressure policies.

## Frozen real pressure policies

1. **no_refresh** — unique writes only.
2. **alternating_shield** — before odd-numbered writes only, refresh the currently predicted non-endangered victim.
3. **fixed_offset_refresh** — before every write, refresh the resident at slot (hand+5) mod 16 if present and not endangered.
4. **hostile_shield** — original UP-161C policy: refresh the currently predicted non-endangered victim before every write.

All arms receive unique writes and no direct refresh of the endangered item.

## Monitoring cadences

- cadence 2;
- cadence 4.

Warnings remain active through the current interval and never modify memory.

## Measurements

Per real policy × cadence:
- arms;
- eventual losses;
- warned losses;
- missed losses;
- expired warning intervals;
- warning recall.

## Interpretation

High recall across policies supports transfer of the current-state robust horizon beyond the exact policy used to derive it. Misses identify where the horizon is not truly policy-robust.

## Bounds

Shadow only. No corrective action, no future real-policy schedule passed into the warning model, no adaptive cadence, no capacity change, no semantic priority, no live activation, or production authority.
