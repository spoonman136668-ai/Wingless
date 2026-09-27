package unitary

import (
	"fmt"
	"math"
)

const UPLM5ODeadlineOrderSchema = "wingless.up-lm5o-deadline-order-ablation.v1"

type UPLM5OSummary struct {
	DeadlineProfile string `json:"deadline_profile"`
	BudgetReduction int `json:"budget_reduction"`
	ThroughputReduction int `json:"throughput_reduction"`
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
	MaxAbsEarliestFixedDelta int `json:"max_abs_earliest_fixed_delta"`
	MaxAbsEarliestLatestDelta int `json:"max_abs_earliest_latest_delta"`
}

type UPLM5OResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM5NSeal string `json:"source_up_lm5n_seal"`
	Policies []string `json:"policies"`
	PooledConditionCells int `json:"pooled_condition_cells"`
	MatchedConditionsPerCell int `json:"matched_conditions_per_cell"`
	OnlyPolicyChanged bool `json:"only_policy_changed"`
	AdaptivePolicySelectionUsed bool `json:"adaptive_policy_selection_used"`
	AdaptiveCoordinateSearchUsed bool `json:"adaptive_coordinate_search_used"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	Summaries []UPLM5OSummary `json:"summaries"`
}

func uplm5oChooseLatest(arms []uplm2xArm, used map[int]bool) int {
	best := -1
	bestCd := -1
	for i := range arms {
		if used[i] {
			continue
		}
		_, cd, ok := uplm2xFirstPending(&arms[i])
		if ok && (cd > bestCd || (cd == bestCd && (best < 0 || i < best))) {
			best = i
			bestCd = cd
		}
	}
	return best
}

func uplm5oRun(rot int, perm, profile, policy string, budget, start, tp int, w UPLM4WWindow, bred, tred int) int {
	arms := uplm2yPermute(uplm2xArms(rot), perm)
	uplm3tPrepressure(arms, rot, profile)
	actions := 0
	for round := 0; round < 6; round++ {
		if round >= start {
			cap := uplm4wCap(budget, round, w.BudgetCut, w.BudgetRestore, bred)
			eff := uplm4wTP(tp, round, w.ThroughputCut, w.ThroughputRestore, tred)
			used := map[int]bool{}
			for k := 0; k < eff && actions < cap; k++ {
				i := -1
				if policy == "latest_deadline" {
					i = uplm5oChooseLatest(arms, used)
				} else {
					i = uplm3cChoose(arms, policy, used)
				}
				if i < 0 {
					break
				}
				if n, _, ok := uplm2xFirstPending(&arms[i]); ok {
					arms[i].reported[n] = true
					actions++
					used[i] = true
				} else {
					break
				}
			}
		}
		for i := range arms {
			arms[i].r.write(fmt.Sprintf("5o-%s-%s-%s-%d-%d-%d-%d", profile, policy, perm, budget, start, tp, round), "x")
		}
	}
	failed := 0
	for i := range arms {
		_, f := uplm2vFinish(arms[i].r, arms[i].reported, rot)
		failed += f
	}
	return failed
}

func RunUPLM5O() (UPLM5OResult, error) {
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

	res := UPLM5OResult{
		Schema: UPLM5ODeadlineOrderSchema,
		Experiment: "UP-LM5O-deadline-order-ablation",
		SourceUPLM5NSeal: "df80cefe725c82f7af0ac86c0abf7313b68ddf79",
		Policies: []string{"earliest_deadline", "fixed_order", "latest_deadline"},
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
						s := UPLM5OSummary{
							DeadlineProfile: profile,
							BudgetReduction: br,
							ThroughputReduction: tr,
							Rotation: rot,
							Permutation: perm,
						}
						for _, w := range windows {
							for _, b := range budgets {
								for _, st := range starts {
									for _, tp := range tps {
										earliest := uplm5oRun(rot, perm, profile, "earliest_deadline", b, st, tp, w, br, tr)
										fixed := uplm5oRun(rot, perm, profile, "fixed_order", b, st, tp, w, br, tr)
										latest := uplm5oRun(rot, perm, profile, "latest_deadline", b, st, tp, w, br, tr)
										s.MatchedConditions++
										df := earliest - fixed
										dl := earliest - latest
										if df != 0 {
											s.EarliestFixedDiffering++
										}
										if dl != 0 {
											s.EarliestLatestDiffering++
										}
										if df < 0 {
											s.EarliestLowerThanFixed++
										} else if df > 0 {
											s.FixedLowerThanEarliest++
										} else {
											s.EarliestFixedEqual++
										}
										if dl < 0 {
											s.EarliestLowerThanLatest++
										} else if dl > 0 {
											s.LatestLowerThanEarliest++
										} else {
											s.EarliestLatestEqual++
										}
										adf := int(math.Abs(float64(df)))
										adl := int(math.Abs(float64(dl)))
										if adf > s.MaxAbsEarliestFixedDelta {
											s.MaxAbsEarliestFixedDelta = adf
										}
										if adl > s.MaxAbsEarliestLatestDelta {
											s.MaxAbsEarliestLatestDelta = adl
										}
									}
								}
							}
						}
						res.Summaries = append(res.Summaries, s)
					}
				}
			}
		}
	}
	return res, nil
}
