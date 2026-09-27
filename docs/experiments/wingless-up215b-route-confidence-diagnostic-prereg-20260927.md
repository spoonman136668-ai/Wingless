# Wingless UP-215B — native route-confidence diagnostic

Status: preregistered scientific shadow-state diagnostic.

Scientific parent: sealed UP-214B 3664d805bbb1a9acf70228edb8286a8df9336838.

## Question

Are the persistent mixed4 route errors detectable from the frozen router's own native confidence geometry?

## Frozen regimes
- mixed4: 0,5,1,6
- observe4: 5,6,7,8
- store4: 0,1,2,3
- cross3: 0,5,13

## Frozen native features
Canonical history only:
- native_correct_count
- mean_absolute_margin
- near_zero_margin_count
- min_absolute_margin

## Frozen training
- phases 58 through 72
- pooled standardization
- one centroid per regime
- nearest Euclidean centroid

No retraining from the UP-212B / UP-213B router design.

## Untouched evaluation phases
- 160 through 183 inclusive

24 decisions per regime; 96 total.

## Frozen confidence observables
For each route decision:
- nearest-centroid distance
- second-nearest-centroid distance
- confidence margin = second-nearest distance - nearest distance
- predicted regime
- true regime

No confidence threshold is chosen and no route is altered.

## Measurements
Per true regime:
- correct / total
- mean nearest distance for correct decisions
- mean nearest distance for incorrect decisions
- mean confidence margin for correct decisions
- mean confidence margin for incorrect decisions
- minimum and maximum confidence margin

Overall:
- route accuracy
- confusion counts for all 4×4 true/predicted pairs

## Interpretation

Lower confidence margin on misroutes would show the ambiguity is visible to Wingless's own native routing state and could support a separately preregistered bounded fallback rule. Similar or higher confidence on misroutes would mean the current geometry cannot safely self-identify routing uncertainty.

## Bounds

Shadow diagnostic only. No evaluation-label fitting, retraining, adaptive feature selection, confidence threshold, route fallback, phase/parity input, nonlinear classifier, maintenance action, live activation, or production authority.
