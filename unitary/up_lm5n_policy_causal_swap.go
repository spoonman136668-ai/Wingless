package unitary

import "math"

const UPLM5NPolicyCausalSchema = "wingless.up-lm5n-policy-causal-swap.v1"

type UPLM5NSummary struct {
	DeadlineProfile string  `json:"deadline_profile"`
	BudgetReduction int     `json:"budget_reduction"`
	ThroughputReduction int `json:"throughput_reduction"`
	Rotation int            `json:"rotation"`
	Permutation string      `json:"permutation"`
	MatchedPairs int        `json:"matched_pairs"`
	DifferingPairs int      `json:"differing_pairs"`
	MaxAbsFailureDelta int  `json:"max_abs_failure_delta"`
	MeanAbsFailureDelta float64 `json:"mean_abs_failure_delta"`
	EarliestLowerFailures int `json:"earliest_lower_failures"`
	FixedLowerFailures int    `json:"fixed_lower_failures"`
	EqualFailures int         `json:"equal_failures"`
}

type UPLM5NResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM5MSeal string `json:"source_up_lm5m_seal"`
	Policies []string `json:"policies"`
	PooledConditionCells int `json:"pooled_condition_cells"`
	MatchedConditionsPerCell int `json:"matched_conditions_per_cell"`
	OnlyPolicyChanged bool `json:"only_policy_changed"`
	AdaptivePolicySelectionUsed bool `json:"adaptive_policy_selection_used"`
	AdaptiveCoordinateSearchUsed bool `json:"adaptive_coordinate_search_used"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	Summaries []UPLM5NSummary `json:"summaries"`
}

func RunUPLM5N() (UPLM5NResult, error) {
	profiles := []string{"deferred_only", "layout_only", "hybrid_min"}
	windows := []UPLM4WWindow{
		{BudgetCut: 0, BudgetRestore: 1, ThroughputCut: 1, ThroughputRestore: 3},
		{BudgetCut: 1, BudgetRestore: 3, ThroughputCut: 0, ThroughputRestore: 1},
		{BudgetCut: 0, BudgetRestore: 1, ThroughputCut: 2, ThroughputRestore: 4},
		{BudgetCut: 2, BudgetRestore: 4, ThroughputCut: 0, ThroughputRestore: 1},
		{BudgetCut: 0, BudgetRestore: 1, ThroughputCut: 3, ThroughputRestore: 5},
		{BudgetCut: 3, BudgetRestore: 5, ThroughputCut: 0, ThroughputRestore: 1},
	}
	breds := []int{1, 2}
	treds := []int{1, 2}
	rots := []int{5, 13}
	perms := []string{"identity", "reverse", "rotate2"}
	budgets := []int{4, 5, 6, 7}
	starts := []int{2, 3, 4, 5}
	tps := []int{2, 3, 4, 5}

	r := UPLM5NResult{
		Schema: UPLM5NPolicyCausalSchema,
		Experiment: "UP-LM5N-policy-causal-swap",
		SourceUPLM5MSeal: "3c18dc54cce22bbc133852193838977c32646de3",
		Policies: []string{"earliest_deadline", "fixed_order"},
		PooledConditionCells: 72,
		MatchedConditionsPerCell: 384,
		OnlyPolicyChanged: true,
		AdaptivePolicySelectionUsed: false,
		AdaptiveCoordinateSearchUsed: false,
		CounterfactualOnly: true,
		LiveActivation: false,
	}

	for _, profile := range profiles {
		for _, br := range breds {
			for _, tr := range treds {
				for _, rot := range rots {
					for _, perm := range perms {
						s := UPLM5NSummary{
							DeadlineProfile: profile,
							BudgetReduction: br,
							ThroughputReduction: tr,
							Rotation: rot,
							Permutation: perm,
						}
						totalAbs := 0
						for _, w := range windows {
							for _, b := range budgets {
								for _, st := range starts {
									for _, tp := range tps {
										earliest := uplm5mRun(rot, perm, profile, "earliest_deadline", b, st, tp, w, br, tr)
										fixed := uplm5mRun(rot, perm, profile, "fixed_order", b, st, tp, w, br, tr)
										s.MatchedPairs++
										d := earliest - fixed
										ad := int(math.Abs(float64(d)))
										totalAbs += ad
										if ad > s.MaxAbsFailureDelta {
											s.MaxAbsFailureDelta = ad
										}
										if d < 0 {
											s.EarliestLowerFailures++
										} else if d > 0 {
											s.FixedLowerFailures++
										} else {
											s.EqualFailures++
										}
										if d != 0 {
											s.DifferingPairs++
										}
									}
								}
							}
						}
						if s.MatchedPairs > 0 {
							s.MeanAbsFailureDelta = float64(totalAbs) / float64(s.MatchedPairs)
						}
						r.Summaries = append(r.Summaries, s)
					}
				}
			}
		}
	}
	return r, nil
}
