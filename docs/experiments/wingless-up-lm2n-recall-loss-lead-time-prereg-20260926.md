# Wingless UP-LM2N — recall-loss lead time

Status: preregistered scientific fixed-capacity shadow-observability experiment.

Scientific parent: sealed UP-LM2M e4482856eb3f80456c88997ac6a80496d99986a0.

## Question

At the established D=5 dependency boundary, how long before the deferred REPORT request does Wingless's native recall state reveal that one still-pending dependency has already been evicted?

## Frozen stream

Reuse the LM2M evaluation structure:
- 24 entities in two chunks of 12;
- exact recall cap 16;
- deferred first-chunk reports = 5;
- identity rotations 0 and 7;
- deferred report order forward and reverse;
- deterministic value shifts 0,1,2,3.

This experiment observes the recall mechanism only. It does not train or modify the language model.

## Pending dependency definition

First-chunk positions 7..11 are the five deferred REPORT obligations. A dependency is pending after its STORE and until its deferred REPORT event.

## Shadow detection

After every event, inspect whether any still-pending deferred key that was previously present in recall has disappeared.

Record:
- the event ordinal of that internal recall loss;
- the later REPORT event ordinal for that same dependency;
- lead_events = report_event - loss_event.

No intervention occurs.

## Interpretation

A positive lead establishes that Wingless's native memory state can expose dependency failure before the stream later asks for the missing value. Report order may change the available lead window.

## Bounds

No extra memory, no model training, no intervention, no semantic priority, no future oracle, no adaptive schedule, no result-informed retry, no live activation, or production authority.
