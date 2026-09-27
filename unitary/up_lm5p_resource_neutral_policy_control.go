package unitary

const UPLM5PResourceNeutralSchema = "wingless.up-lm5p-resource-neutral-policy-control.v1"

type UPLM5PSummary struct {
	DeadlineProfile string `json:"deadline_profile"`
	Rotation int `json:"rotation"`
	Permutation string `json:"permutation"`
	MatchedConditions int `json:"matched_conditions"`
	EarliestFixedDiffering int `json:"earliest_fixed_differing"`
	EarliestLatestDiffering int `json:"earliest_latest_differing"`
	EarliestLowerThanFixed int `json:"earliest_lower_than_fixed"`
	FixedLowerThanEarliest int `json:"fixed_lower_than_earliest"`
	EarliestFixedEqual int `json:"earliest_fixed_equal"`
	EarliestLowerThanLatest int `json:"earliest_lower_than_latest"`
	LatestLowerThanEarliest int `json:"latest_lower_than_earliest"`
	EarliestLatestEqual int `json:"earliest_latest_equal"`
}

type UPLM5PResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM5OSeal string `json:"source_up_lm5o_seal"`
	Policies []string `json:"policies"`
	ConditionCells int `json:"condition_cells"`
	MatchedConditionsPerCell int `json:"matched_conditions_per_cell"`
	ResourceReductionsUsed bool `json:"resource_reductions_used"`
	BackgroundProfilesPreserved bool `json:"background_profiles_preserved"`
	OnlyPolicyChanged bool `json:"only_policy_changed"`
	AdaptivePolicySelectionUsed bool `json:"adaptive_policy_selection_used"`
	AdaptiveCoordinateSearchUsed bool `json:"adaptive_coordinate_search_used"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	Summaries []UPLM5PSummary `json:"summaries"`
}

func RunUPLM5P() (UPLM5PResult, error) {
	profiles := []string{"deferred_only", "layout_only", "hybrid_min"}
	rots := []int{5, 13}
	perms := []string{"identity", "reverse", "rotate2"}
	budgets := []int{4, 5, 6, 7}
	starts := []int{2, 3, 4, 5}
	tps := []int{2, 3, 4, 5}
	neutral := UPLM4WWindow{}

	res := UPLM5PResult{
		Schema: UPLM5PResourceNeutralSchema,
		Experiment: "UP-LM5P-resource-neutral-policy-control",
		SourceUPLM5OSeal: "19079148a8405e8b4eb625f69206751316600476",
		Policies: []string{"earliest_deadline", "fixed_order", "latest_deadline"},
		ConditionCells: 18,
		MatchedConditionsPerCell: 64,
		ResourceReductionsUsed: false,
		BackgroundProfilesPreserved: true,
		OnlyPolicyChanged: true,
		AdaptivePolicySelectionUsed: false,
		AdaptiveCoordinateSearchUsed: false,
		CounterfactualOnly: true,
		LiveActivation: false,
	}

	for _, profile := range profiles {
		for _, rot := range rots {
			for _, perm := range perms {
				s := UPLM5PSummary{
					DeadlineProfile: profile,
					Rotation: rot,
					Permutation: perm,
				}
				for _, budget := range budgets {
					for _, start := range starts {
						for _, tp := range tps {
							earliest := uplm5oRun(rot, perm, profile, "earliest_deadline", budget, start, tp, neutral, 0, 0)
							fixed := uplm5oRun(rot, perm, profile, "fixed_order", budget, start, tp, neutral, 0, 0)
							latest := uplm5oRun(rot, perm, profile, "latest_deadline", budget, start, tp, neutral, 0, 0)
							s.MatchedConditions++
							if earliest != fixed {
								s.EarliestFixedDiffering++
							}
							if earliest != latest {
								s.EarliestLatestDiffering++
							}
							if earliest < fixed {
								s.EarliestLowerThanFixed++
							} else if fixed < earliest {
								s.FixedLowerThanEarliest++
							} else {
								s.EarliestFixedEqual++
							}
							if earliest < latest {
								s.EarliestLowerThanLatest++
							} else if latest < earliest {
								s.LatestLowerThanEarliest++
							} else {
								s.EarliestLatestEqual++
							}
						}
					}
				}
				res.Summaries = append(res.Summaries, s)
			}
		}
	}
	return res, nil
}
