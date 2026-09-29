UP-LM5Z PREREGISTRATION — RAW ORDER-LENGTH ORIGIN DECOMPOSITION

Parent UP-LM5Y qualified scientific result.
Observed parent anchor: every eligible candidate across the frozen LM5X/LM5Y manifold had pending-order length exactly 16; subhazard eligible candidates=0.

Question: are shorter pending-order states genuinely absent from the underlying arm dynamics, or are they present but removed by the frozen eligibility predicate before selection?

Freeze exact LM5Y substrate and canonical hazard-selector trajectory:
- same 27,648 matched conditions;
- same deadline profiles, resource reductions, rotations, permutations, windows, budgets, starts, throughputs;
- same prepressure;
- same canonical hazard selector and action policy;
- no future-schedule information;
- no adaptive policy selection;
- counterfactual only;
- no live activation.

Immediately before every selector decision, inspect every arm and assign it to exactly one mutually exclusive category in this order:
1. USED_THIS_DECISION — arm already selected in the current selector decision.
2. EMPTY — pending order empty.
3. NON_ORIGINAL_VICTIM — first pending victim is not one of the original 12 victims.
4. ALREADY_REPORTED — first pending victim is original but already reported.
5. ELIGIBLE — exact LM5Y eligible predicate.

For every non-empty arm, record its exact pending-order length regardless of category.
Report:
- global raw non-empty order-length histogram;
- category-by-length histograms;
- total counts per category;
- raw sub-16 candidate count;
- eligible sub-16 candidate count;
- same summaries per frozen condition cell.

Parent anchors must reproduce exactly:
- observed eligible lengths=[16];
- eligible subhazard candidates=0;
- LM5Y decision point and eligible candidate totals.

Interpretation:
- raw sub-16=0 means the hard length-16 support comes from dynamics/state construction itself;
- raw sub-16>0 with eligible sub-16=0 localizes the boundary to eligibility filtering and identifies the exclusion category;
- eligible sub-16>0 is parent-anchor drift.

Observational only. No policy change. Scientific negatives are valid. No post-result tuning.
