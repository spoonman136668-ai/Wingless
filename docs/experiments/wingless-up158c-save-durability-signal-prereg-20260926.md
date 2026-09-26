# Wingless UP-158C — save durability signal

Status: preregistered scientific counterfactual diagnostic.

Scientific parent: sealed UP-157C e528d4509aeb4de120b343328e4f90b75fb09ca5.

## Question

After a warning-guided refresh successfully prevents the immediate protected-item loss, can Wingless's *present native replacement state* distinguish saves that will persist from saves that are merely delayed?

## Frozen design

Reuse the exact UP-157C:
- 16 arms: four protected cohorts × four starting hands;
- 24-step mixed write stream;
- same sparse refresh schedule;
- same one-action guided refresh at the first predicted harmful write;
- exact recall cap 16;
- counterfactual only.

At the trigger, after applying the guided refresh but before the harmful write, compute one diagnostic from a copy of current memory state only:

**Unique-write eviction horizon** = number of hypothetical unique writes required to evict the endangered item if no further queries/refreshes occurred.

The diagnostic may read only:
- current replacement hand;
- current two-bit age state;
- current location of the already-known endangered item.

It does not inspect future stream events, future refresh schedule, semantic class, final outcome, or any oracle label.

Then discard the diagnostic copy and run the original frozen guided arm unchanged to completion.

## Measurements

For each arm:
- cohort;
- starting hand;
- trigger step;
- endangered key;
- present-state unique-write horizon;
- final outcome: permanent_save or delayed_loss;
- actual delayed-loss step when applicable.

Report:
- mean/min/max horizon for permanent saves;
- mean/min/max horizon for delayed losses;
- AUROC using horizon as a frozen scalar durability score.

## Interpretation

If larger present-state horizon ranks permanent saves above delayed losses, native state contains information about correction durability. If not, the permanence of a save depends on future pressure not recoverable from this local state alone.

## Bounds

No intervention policy change, no threshold fitting, no semantic priority, no future schedule input, no capacity increase, no extra training, no live activation, no result-informed retry.
