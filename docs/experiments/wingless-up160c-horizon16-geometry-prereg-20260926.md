# Wingless UP-160C — horizon-16 native geometry

Status: preregistered scientific diagnostic.

Scientific parent: sealed UP-159C 3f621ebede6fc7b4aa09bcb0eeb5b06bde295943.

## Question

UP-159C found one durable save and several delayed losses sharing the same scalar unique-write horizon of 16. Does additional present native replacement geometry distinguish those outcomes without future schedule information?

## Frozen design

Reuse UP-159C exactly:
- disjoint mixed-write schedule;
- sparse refresh schedule;
- four protected cohorts × four starting hands;
- same warning-guided refresh;
- exact recall cap 16.

Only arms whose post-refresh unique-write horizon equals 16 are included in the primary diagnostic.

Immediately after guided refresh and before the threatened write, record only present native state:
- current hand;
- endangered slot;
- cyclic slot distance from hand to endangered slot;
- endangered two-bit age;
- count of zero-age slots;
- count of nonzero-age slots;
- sum of all age counters;
- maximum age counter;
- count of slots with age >= endangered age.

No future stream event, future refresh, semantic class, final outcome, or oracle label is read by the diagnostic.

Then run the frozen arm unchanged to completion and label actual durability.

## Interpretation

If the durable exception differs from delayed-loss horizon-16 cases on a preregistered geometry field, that field becomes a candidate durability signal. If the recorded native geometry is identical or nonseparating, the exception likely depends on future pressure structure outside this present local summary.

## Bounds

Diagnostic only. No intervention policy change, no threshold fitting, no capacity increase, no extra training, no live activation, no future schedule input.
