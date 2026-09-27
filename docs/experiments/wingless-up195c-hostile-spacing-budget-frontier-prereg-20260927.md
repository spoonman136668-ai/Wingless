# Wingless UP-195C — hostile spacing × budget frontier

Status: preregistered scientific counterfactual intervention diagnostic.

Scientific parent: sealed UP-194C 5a25ed974a31ca3d34a4f2b8a55133db800ff361.

## Question

What is the smallest bounded targeted-refresh budget that preserves hostile-pressure prevention when corrective capacity is spread over time?

## Frozen policy
- hostile_shield

## Frozen monitoring cadence
- 2 writes

## Frozen first trigger
- first native critical warning: adversarial_shield_horizon <= 2

## Frozen post-warning spacings
- every 2 monitoring intervals = every 4 writes
- every 3 monitoring intervals = every 6 writes
- every 4 monitoring intervals = every 8 writes

## Frozen action caps
- 4
- 6
- 8
- 10
- 12
- 16

At the first critical warning, one targeted refresh is used. Later targeted refreshes follow the frozen spacing until the cap is exhausted.

## Measurements

Per spacing × cap:
- arms
- baseline losses
- prevented losses
- actions used
- accelerated losses
- mean and maximum treated-minus-baseline loss-step extension among failures

## Interpretation

The smallest cap producing prevention identifies the onset of sufficient temporal coverage. The smallest cap producing 64/64 prevention identifies the full-rescue boundary for that spacing.

## Bounds

Counterfactual only. No adaptive spacing, adaptive cap, threshold tuning, policy change, new action type, future-policy input to warning, live activation, or production authority.
