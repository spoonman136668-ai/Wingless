# Wingless UP-159C — durability-horizon replication

Status: preregistered scientific counterfactual replication.

Scientific parent: sealed UP-158C 60d470b7d11aae27c98a3d73d62bea2984b8e5e6.

## Question

Does the present-state unique-write horizon continue to distinguish durable warning-guided saves from delayed-only saves under a disjoint mixed-write and refresh schedule?

## Frozen disjoint stream

Reuse:
- exact recall cap 16;
- same four protected cohorts;
- starting hands 0,4,8,12;
- one guided refresh of the first endangered cohort item;
- counterfactual-only execution.

Change the workload before results:
- 24 steps;
- every 4th write repeats the previous novel key; all other writes are novel;
- sparse refreshes:
  - step 2 -> key 5
  - step 6 -> key 10
  - step 11 -> key 15
  - step 15 -> key 4
  - step 20 -> key 9

At the first predicted harmful novel write for each cohort/hand, refresh the endangered item and compute the same present-state unique-write eviction horizon from a state copy.

## Frozen durability rule

C158 separated delayed horizon=16 from permanent horizons>=25. Before this run, freeze:

- predict durable save if horizon >= 21;
- predict delayed-only otherwise.

No threshold changes after seeing this run.

## Measurements

- permanent vs delayed outcomes;
- horizon AUROC;
- frozen-rule TP/FP/FN/TN;
- precision, recall, accuracy;
- horizon ranges by outcome.

## Bounds

No future schedule is read by the horizon predictor. No semantic class, adaptive threshold, capacity increase, extra training, result-informed retry, live activation, or production authority.
