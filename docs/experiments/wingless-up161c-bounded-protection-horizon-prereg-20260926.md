# Wingless UP-161C — bounded protection horizon

Status: preregistered scientific counterfactual stress test.

Scientific parent: sealed UP-160C c70c1f53358d88dd984080724562242920c19ff1.

## Question

After a warning-guided refresh, how many unique writes can the endangered item survive under a present-state-only adversarial refresh stress, without assuming the actual future stream?

## Frozen starting states

Reuse the 16 C159 trigger states:
- four protected cohorts × four starting hands;
- same first predicted harmful write;
- same guided refresh of the endangered item;
- exact recall cap 16.

After the guided refresh, discard the original future stream and fork diagnostic copies.

## Frozen stress arms

1. **No-query horizon:** unique writes only, identical to the accepted UP-158C horizon.
2. **Adversarial non-endangered shielding:** before every unique write, predict the next replacement slot from current hand+ages. If that slot does not contain the endangered item, query that resident item to protect it; then issue the unique write. If the predicted slot is the endangered item, do not refresh the endangered item and allow the write.

The stress policy reads only current state. It does not know the original future schedule or final outcome.

## Measurements

Per starting state:
- no-query horizon;
- adversarial-shield horizon;
- tested robust horizon = min(two horizons).

Aggregate min/mean/max robust horizon.

## Interpretation

This does not claim a universal mathematical guarantee over every possible future query sequence. It establishes a bounded tested protection horizon under a deliberately hostile native-state policy.

## Bounds

Counterfactual diagnostic only. No live activation, semantic priority, future schedule input, capacity increase, extra training, or adaptive intervention.
