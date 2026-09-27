# Wingless UP-200C — hostile internal-trigger vs latest-safe oracle onset

Status: preregistered scientific counterfactual diagnostic.

Scientific parent: sealed UP-199C ca99fbda42a5e5b3e4b8d5bd9287c7329617387c.

## Question

How much action-efficiency headroom remains between Wingless's frozen internal critical-warning trigger and an external oracle that may choose the latest safe rescue onset, when both use the exact same correction schedule and budget?

## Frozen environment

Policy:
- hostile_shield

Heldout cohorts:
- contiguous quartets 0,1,2,3 / 4,5,6,7 / 8,9,10,11 / 12,13,14,15

Initial hand positions:
- 0..15

Monitoring cadence:
- 2 writes

Evaluation horizon:
- 56 writes

## Frozen correction schedule

After onset:
- immediate targeted refresh of endangered resident;
- then one targeted refresh every 4 monitoring intervals = every 8 writes;
- maximum 8 actions.

## Internal trigger

- onset at first monitoring boundary where adversarial_shield_horizon <= cadence.
- no future outcome information.

## Oracle onset

Offline comparator only.

For each arm:
- candidate onsets are all monitoring boundaries within the horizon;
- simulate the exact same spacing-4/cap-8 schedule from each candidate onset;
- among candidate onsets that preserve the endangered resident through write 56, choose the latest one;
- if none survive, record no safe oracle onset.

The oracle may inspect future outcomes only to choose onset. It receives no extra actions, no denser spacing, no different target, and no altered state.

## Measurements

Aggregate:
- valid arms
- baseline losses
- internal prevented losses
- internal actions
- oracle-safe arms
- oracle actions
- internal mean onset write
- oracle mean latest-safe onset write
- mean onset headroom = oracle onset - internal onset
- total action savings of oracle vs internal
- arms where internal onset equals latest-safe onset

Per-arm records:
- cohort index
- initial hand
- baseline loss step
- internal onset/action count/survival
- oracle latest-safe onset/action count
- onset headroom

## Interpretation

If the oracle starts materially later and uses fewer actions while preserving the same rescue rate, the internal trigger is conservative but effective. Little or no oracle headroom would show that the current warning onset is already close to the latest safe point and the remaining action burden is intrinsic to long-horizon coverage.

## Bounds

Diagnostic only. Oracle future access is comparator-only and may not feed any candidate rule. No threshold, cadence, spacing, cap, target, cohort, or horizon tuning; no live activation or production authority.
