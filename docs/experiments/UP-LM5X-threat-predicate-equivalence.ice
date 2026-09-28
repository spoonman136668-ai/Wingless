UP-LM5X PREREGISTRATION — THREAT-PREDICATE EQUIVALENCE
Parent UP-LM5W qualified scientific result.
Parent anchor: hazard_triggered and prehazard_triggered produced zero selector, action, or failure differences across 27,648 matched conditions.

Question: are the hazard and prehazard predicates themselves equivalent on the evaluated state manifold, or does prehazard expose additional len=15 candidates that are masked because a hazard-priority candidate is always available?

Freeze exact UP-LM5W transfer substrate:
- deadline profiles deferred_only, layout_only, hybrid_min;
- budget reductions {1,2};
- throughput reductions {1,2};
- rotations {5,13};
- permutations identity, reverse, rotate2;
- six frozen reduction windows;
- budgets {4,5,6,7};
- starts {2,3,4,5};
- throughputs {2,3,4,5};
- prepressure on;
- no future-schedule information;
- no adaptive policy selection;
- counterfactual only;
- no live activation.

Use the already-proven hazard-equivalent trajectory as the canonical evolution. At every selector decision point BEFORE selection, count:
- hazard candidates: eligible original unreported victims with order length >=16;
- prehazard-only candidates: eligible original unreported victims with order length exactly 15.

Report:
- total decision points;
- decision points with any prehazard-only candidate;
- decision points with hazard=0 and prehazard-only>0;
- total hazard candidates;
- total prehazard-only candidates;
- the same statistics per frozen condition cell.

Interpretation:
- prehazard-only count=0 implies predicate equivalence on this manifold;
- prehazard-only>0 but exposed=0 implies priority masking;
- exposed>0 identifies states where the extended predicate could actually alter selection.

Scientific negatives are valid. No post-result tuning.