# Wingless UP-LM2S — bounded early closure

Status: preregistered counterfactual closed-loop simulation.

Scientific parent: sealed UP-LM2R 4e254762af5e380cbaf689002c884d5574512888.

## Question

Can the exact native warning prevent unreported dependency loss with a bounded sequence-only response, more often than it causes unnecessary intervention?

## Frozen arms

Use:
- exact recall cap 16;
- deferred levels D=4,5,6;
- identity rotations 0 and 7;
- value shifts 0,1,2,3;
- the same 12 first-chunk memories and 12 second-chunk unique STOREs.

## Baseline arm

Run the 12 second-chunk STOREs with no early closure. Count evictions of identities that have not yet been REPORTed.

## Counterfactual response arm

Immediately before each second-chunk STORE:
- if recall is full, inspect the current FIFO victim;
- if that victim has not yet been REPORTed, move only that victim's pending REPORT to immediately before the STORE by marking it closed now;
- then execute the same STORE.

No extra REPORT content is added. The moved REPORT is removed from its later pending position. Total dependency content is preserved.

## Measurements

Per D × rotation × shift:
- baseline unreported evictions;
- intervention count;
- counterfactual unreported evictions;
- prevented evictions;
- unnecessary interventions;
- moved versus remaining deferred REPORT count.

Aggregate:
- prevented fraction;
- intervention precision = prevented / interventions;
- intervention recall = prevented / baseline harmful evictions.

## Bounds

Counterfactual simulation only. No live activation, no memory increase, no extra training, no write suppression, no semantic priority, no adaptive threshold, no result-informed retry, or production authority.
