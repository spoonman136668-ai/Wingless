# Wingless UP-167C — warning confirmation tradeoff

Status: preregistered scientific shadow-monitoring experiment.

Scientific parent: sealed UP-166C 5e4879b1653490ebae62421f84620ce8e1c6a667.

## Question

Can repeated current-state warning evidence reduce trigger burden while preserving pre-loss recall?

## Frozen warning model

Reuse the UP-165C/166C warning unchanged:
- current-state adversarial-shield horizon;
- cadence 2 warns iff horizon <= 2;
- cadence 4 warns iff horizon <= 4.

## Frozen confirmation rules

A trigger-eligible interval requires:
- 1 consecutive warning interval;
- 2 consecutive warning intervals;
- 3 consecutive warning intervals.

The rule is evaluated independently for each arm. Confirmation never modifies memory.

## Frozen policies and cadences

Policies:
- no_refresh;
- alternating_shield;
- fixed_offset_refresh;
- hostile_shield.

Cadences:
- 2;
- 4.

## Measurements

Per policy × cadence × confirmation:
- eventual losses;
- losses occurring in trigger-eligible intervals;
- missed losses;
- trigger-eligible intervals;
- expired trigger-eligible intervals;
- trigger precision;
- trigger recall;
- first-trigger lead writes.

## Interpretation

Confirmation is useful only if it reduces expired trigger-eligible intervals without destroying recall. No rule is selected or activated in this experiment.

## Bounds

Shadow only. No corrective action, adaptive confirmation, threshold tuning, future policy schedule input, capacity change, live activation, or production authority.
