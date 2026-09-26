# Wingless UP-116B — stores-versus-keeps margin diagnostic

Status: preregistered scientific lexical-interference diagnostic.

Scientific parent: sealed UP-115B 96e7968598df235267e07d099f0d77171a183f89.

## Question

UP-115B localized every base-lexicon error to the single STORE verb `stores`, while `keeps` and every OBSERVE/REPORT base verb remained perfect. Is `stores` failing through a consistent low-margin decision boundary and a consistent wrong semantic class?

## Frozen training

Repeat the exact UP-115B / UP-114B matrix without any learning change:
- six acquisition orders: SOR, SRO, OSR, ORS, RSO, ROS;
- policies: current_class_excluded and stage3_anchor2;
- same base/acquired lexicon;
- same 64-D encoder;
- four grounding subjects per current verb;
- 20 epochs per acquisition batch;
- learning rate 0.08;
- same fixed replay budget and policies.

## Frozen diagnostic measurements

At stage 0 and after every acquisition, for both base STORE verbs `stores` and `keeps`, evaluate all six subjects.

For each verb report:
- accuracy;
- mean target-class margin;
- minimum target-class margin;
- predicted STORE count;
- predicted OBSERVE count;
- predicted REPORT count.

Margin is:
target-class probability minus the maximum non-target class probability.

Also preserve aggregate base and acquired accuracy for cross-checking the parent result.

No metric affects training.

## Interpretation

A consistently negative or collapsing `stores` margin with stable `keeps` margins would confirm lexical-surface geometry as the remaining B-lane failure. A consistent wrong-class destination would identify the competing decision direction to test next.

## Bounds

No training-rule change, no replay-budget change, no adaptive intervention, no threshold tuning, no class weighting, no state change, no result-informed retry, no live activation, no production authority.
