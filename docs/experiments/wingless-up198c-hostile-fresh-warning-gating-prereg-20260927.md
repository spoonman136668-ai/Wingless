# Wingless UP-198C — hostile fresh-warning gating ablation

Status: preregistered scientific counterfactual intervention experiment.

Scientific parent: sealed UP-197C d5dcbca3fde66f0079b73ecd9831267b18ec6d1e.

## Question

Can later hostile-rescue actions be reduced by requiring a fresh native critical warning at each scheduled refresh opportunity without losing protection?

## Frozen environment

Policy:
- hostile_shield

Heldout cohorts:
- contiguous quartets 0,1,2,3 / 4,5,6,7 / 8,9,10,11 / 12,13,14,15

Initial hand positions:
- 0..15

Monitoring cadence:
- 2 writes

## Frozen base schedule

- first native critical warning triggers targeted refresh of endangered resident;
- later refresh opportunities every 4 monitoring intervals = every 8 writes;
- maximum 8 actions.

## Frozen arms

Comparator: committed schedule
- every later opportunity refreshes, as in UP-197C.

Primary: fresh-warning-gated schedule
- at each later scheduled opportunity, recompute adversarial_shield_horizon from current native state;
- refresh only if horizon <= cadence;
- if not critical, skip that opportunity;
- next opportunity remains on the original 8-write grid;
- cap remains 8.

No schedule compression after a skipped action.

## Frozen horizons

- 16,24,32,40,48,56 writes.

Each horizon starts from the same initial state.

## Measurements

Per horizon × schedule:
- baseline failures;
- treated failures;
- prevented failures;
- total actions;
- arms with at least one action;
- accelerated failures;
- skipped scheduled opportunities.

Also report action reduction of gated vs committed.

## Interpretation

Equal prevention with fewer actions supports fresh native warning as an efficiency gate. Any loss of protection falsifies the gating rule under the frozen timing.

## Bounds

Counterfactual only. No threshold, cadence, spacing, cap, target, cohort, or horizon tuning after results; no live activation or production authority.
