# Wingless UP-157C — protected-loss redistribution

Status: preregistered counterfactual-only cost-trace experiment.

Scientific parent: sealed UP-156C 527fe0829e1224b02d47e0827d07368ead8ce730.

## Question

When one warning-guided refresh prevents the predicted immediate protected loss, is that save permanent, merely delayed, or paid for by displacing a different original durable item?

## Frozen design

Reuse UP-156C exactly:
- four protected cohorts;
- starting hands 0,4,8,12;
- same 24-step mixed write stream;
- same sparse refresh schedule;
- same first predicted harmful write as trigger;
- baseline, one guided refresh, one matched sham refresh.

## Measurements

Per hand × cohort:
- endangered key;
- baseline/guided/sham trigger-step loss;
- whether endangered key remains at end;
- first later step where endangered key is lost, if any;
- final protected cohort retention;
- final total original-durable retention;
- final non-cohort original retention.

Classify guided outcome:
- permanent_save: endangered key survives to end;
- delayed_loss: immediate loss prevented but endangered key is lost later.

Also measure whether guided final total-original retention differs from baseline. This distinguishes true durable preservation from redistribution within the fixed 16-slot budget.

## Bounds

Counterfactual only. One guided action or one sham action per arm. No live activation, capacity increase, extra training, repeated intervention, future oracle, semantic priority, or result-informed retry.
