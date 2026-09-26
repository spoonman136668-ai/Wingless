# Wingless UP-156C — bounded refresh counterfactual

Status: preregistered counterfactual-only response experiment.

Scientific parent: sealed UP-155C 4ba5db866befcbe662167e0d32cad10a2b317cf5.

## Question

When the exact prewrite warning identifies a protected cohort item that will be overwritten by the next write, can one bounded refresh of that endangered item prevent the immediate protected loss better than a same-cost sham refresh?

## Frozen cohorts and stream

Use the four accepted interleaved protected cohorts from UP-155C and starting hands 0,4,8,12.

Reuse the 24-step mixed stream:
- every third write repeats the most recently introduced novel key;
- all other writes are novel.

Reuse the UP-155C sparse refresh schedule unchanged.

## Trigger

For each hand × cohort arm, identify the first baseline step where the frozen hand+age predictor says the incoming novel write will replace a protected cohort member.

## Counterfactual arms

From the same initial state and same stream:

- baseline: no extra action;
- guided: at the trigger step, issue exactly one query/refresh to the endangered protected key immediately before the write;
- sham: at the same trigger step, issue exactly one query/refresh to the lowest-key currently present non-cohort item other than the endangered key.

Guided and sham have identical one-query action budgets. No later adaptive action is allowed.

## Measurements

Per hand × cohort:
- trigger step and endangered key;
- immediate protected loss in baseline/guided/sham;
- final protected cohort retention;
- guided and sham action targets.

Aggregate:
- baseline immediate harmful losses;
- guided prevented immediate losses;
- sham prevented immediate losses;
- final retention gain versus baseline;
- action counts.

## Bounds

Counterfactual only. No live activation, capacity increase, extra training, repeated intervention, semantic priority, future oracle, result-informed retry, or production authority.
