# Wingless UP-LM3S — pressure saturation diagnostic

Status: preregistered scientific counterfactual scheduling diagnostic.

Scientific parent: sealed UP-LM3R 4f546af41f8e85af2267c13d0ac8c4851efab798.

## Question

Does the apparent high-pressure deterministic collapse remain informative because resource settings still change outcomes, or has the system saturated so all resource configurations fail the same way?

## Frozen scenario

Reuse:
- hybrid_min deadline profile
- rotations 5 and 13
- permutations identity, reverse, rotate2
- budgets 4,5,6
- action start rounds 3,4
- throughput 1,2,3
- six global rounds

## Frozen pressure levels

- 0 through 16 inclusive

## Diagnostic

For each pressure level:
- evaluate all 18 resource configurations;
- aggregate earliest-failed count across the six heldout topology cases exactly as prior LM3 experiments;
- report global minimum outcome;
- global maximum outcome;
- global outcome range;
- number of distinct outcomes across the 18 resource configurations;
- number of configurations at the minimum and maximum.

No coordinate fitting is performed.

## Interpretation

A nonzero range and multiple distinct outcomes mean resource configuration still matters, so zero within-coordinate spread is scientifically informative. Range zero with one distinct outcome identifies pressure saturation.

## Bounds

Diagnostic only. No resource tuning, coordinate search, pressure selection after results, topology search, capacity change, semantic priority, extra training, live activation, or production authority.
