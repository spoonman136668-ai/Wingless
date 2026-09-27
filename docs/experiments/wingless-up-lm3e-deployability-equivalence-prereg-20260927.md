# Wingless UP-LM3E — deployability equivalence

Status: preregistered scientific counterfactual scheduling experiment.

Scientific parent: sealed UP-LM3D delayed-onset result ec7a600cbea615c17edc3bd6dcf79961e9fa2c4c.

## Question

The LM3C/LM3D sequence shows that throughput matters when a short rescue window prevents the fixed budget from being deployed in time. Is the resulting failure boundary explained by the number of actions that can actually be deployed, or does the timing configuration itself matter even when the same number of actions is delivered?

## Frozen scheduler and scenarios

Reuse:
- exact recall cap 16;
- six arms;
- profiles by_deferred_level and by_layout;
- rotations 3 and 11;
- permutations identity and reverse;
- earliest_deadline only;
- one action per arm per round;
- six global pressure rounds.

No future schedule oracle is used.

## Frozen factorial

Total budgets:
- 3;
- 4;
- 5;
- 6.

Intervention onset, zero-based:
- 0;
- 1;
- 2;
- 3.

Throughput:
- 1;
- 2;
- 3.

For each configuration define predicted deployable capacity:

deployable_capacity = min(total_budget, (6 - onset) * throughput).

The policy still acts from current state; the formula is diagnostic only and does not select actions.

## Measurements

Per profile × budget × onset × throughput:
- predicted deployable capacity;
- actual actions used;
- failed originals;
- failures prevented versus no-action baseline.

Then group points by profile and actual actions used:
- number of configurations;
- minimum failures;
- maximum failures;
- failure spread.

## Interpretation

If equal deployed-action counts give equal failures across different budget/onset/throughput configurations, a simple deployability scalar captures this resource boundary.

If failure spread remains at equal action count, timing has an additional causal effect beyond the number of deployed actions.

## Bounds

Counterfactual only. No adaptive budget, onset, or throughput; no capacity change; no semantic priority; no extra training; no live activation; no production authority.
