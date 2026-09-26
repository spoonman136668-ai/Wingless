# Wingless UP-179B — within-cell margin tie-break

Status: preregistered scientific shadow-detection experiment.

Scientific parent: sealed UP-178B 65fa1a872cbaa372dbd72df9c7c7df9ab0269e30.

## Question

Can continuous absolute decision margin improve failure ranking *within* the frozen rank-plus-temporal risk cells, without fitting a new predictor?

## Frozen scores

Evaluate three scores on held-out terminal phase 24:

1. rank_plus_temporal:
   - exact frozen probability from UP-175B/176B.

2. rank_plus_temporal_margin_tiebreak:
   - frozen cell probability
   - plus `1e-6 / (1 + abs(current_margin))`.
   - The epsilon is deliberately too small to reorder examples from different frozen probability cells; it only breaks ties inside a cell.

3. margin_only:
   - `-abs(current_margin)`, so smaller absolute margin ranks as higher risk.

No coefficients are fitted from phase-24 outcomes.

## Measurements

For each score:
- AUROC;
- average precision;
- base crossing rate;
- predicted event count where meaningful;
- actual crossings.

## Bounds

Shadow only. No threshold tuning, intervention, model fitting, maintenance trigger, extra training, capacity change, result-informed retry, live activation, or production authority.
