# Wingless UP-LM2M — parametric fallback remap

Status: preregistered scientific fixed-capacity language experiment.

Scientific parent: sealed UP-LM2L c59dd2bce72ab784dec39ab356c00ccb1d379ee7.

## Question

UP-LM2L established a real memory miss at five deferred dependencies in every tested cell, yet 25/200 language cells still produced an exact report set by predicting the missing value without recall. Does that apparent fallback track a deterministic remapping of the required value, or is it an incidental learned/default guess?

## Frozen system

Reuse the accepted LM2L model, router, recall cap 16, two chunks of 12, five lexical families, two training schedules, identity rotations 0 and 7, and forward/reverse deferred-report order.

Freeze:
- allocation = equal_mass;
- deferred first-chunk dependencies = 5;
- value shifts = 0, 1, 2, 3.

Training is unchanged. Value shift is applied only to the evaluation stream's STORE/REPORT value association using the accepted value alphabet. OBSERVE values receive the same deterministic offset from their shifted entity value.

## Controls

- shift 0 must replicate the parent equal-mass arm structurally;
- memory pressure is identical across shifts;
- recall cap remains 16;
- no added rehearsal, attention, router change, or future oracle.

## Measurements

Per training schedule × family × identity rotation × report order × value shift:
- recall hit rate;
- dependent first-byte accuracy;
- report-set exactness;
- event/class routing;
- max recall entries;
- top-1 accuracy and perplexity.

Also count cases where recall is below 1 but report-set exactness is 1.

## Interpretation

If the same fallback cells remain exact across remapped required values, the language state is supplying structured information beyond explicit recall. If exactness collapses or moves with remap, the prior success was value-specific or incidental. Either result is valid.

## Bounds

No capacity increase, extra training, adaptive chunking, semantic priority, result-informed retry, live activation, or production authority.
