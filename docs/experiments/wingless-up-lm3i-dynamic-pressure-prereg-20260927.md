# Wingless UP-LM3I — dynamic-pressure reachability

Status: preregistered scientific counterfactual scheduling experiment.

Scientific parent: sealed UP-LM3H e2cff4b95aab789a872d53fcdc7b5012c804f0aa.

## Question

Does the reachable-budget law remain sufficient when hazard is nonstationary and pressure increases inside the intervention window?

## Frozen initial deadline profile

Reuse UP-LM3H hybrid_min prepressure unchanged.

## Frozen dynamic perturbation

At global round 3 (zero-based, the fourth pressure round), inject one additional ordinary write into every arm after that round's standard pressure write.

This burst is identical across all resource settings and is fixed before the run.

## Frozen topology

- rotations 5 and 13;
- permutations identity, reverse, rotate2.

## Frozen resource grid

Budgets:
- 4;
- 5;
- 6.

Action start rounds:
- 3;
- 4.

Throughput:
- 1;
- 2;
- 3.

Reachable budget remains:
min(total budget, throughput × remaining action rounds).

## Measurements

Per budget × onset × throughput:
- reachable budget;
- baseline failures;
- earliest-deadline failures;
- fixed-order failures;
- actions used;
- failures prevented.

Group earliest-deadline outcomes by reachable budget and report failure spread.

## Interpretation

Zero spread would extend the resource law to this fixed nonstationary hazard. Nonzero spread would show that timing relative to an internal pressure transition matters beyond aggregate reachable action capacity.

## Bounds

Counterfactual only. Fixed one-round burst. No adaptive burst timing or magnitude, no resource tuning, no future schedule oracle, no capacity change, no semantic priority, no extra training, live activation, or production authority.
