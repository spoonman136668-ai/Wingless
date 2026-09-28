package unitary

import "fmt"

const UPLM5XThreatPredicateSchema = "wingless.up-lm5x-threat-predicate-equivalence.v1"

type UPLM5XCellSummary struct {
	DeadlineProfile            string `json:"deadline_profile"`
	BudgetReduction            int    `json:"budget_reduction"`
	ThroughputReduction        int    `json:"throughput_reduction"`
	Rotation                   int    `json:"rotation"`
	Permutation                string `json:"permutation"`
	Conditions                 int    `json:"conditions"`
	DecisionPoints             int    `json:"decision_points"`
	PredicateDifferencePoints  int    `json:"predicate_difference_points"`
	ExposedPrehazardOnlyPoints int    `json:"exposed_prehazard_only_points"`
	HazardCandidates           int    `json:"hazard_candidates"`
	PrehazardOnlyCandidates    int    `json:"prehazard_only_candidates"`
}

type UPLM5XResult struct {
	Schema                          string              `json:"schema"`
	Experiment                      string              `json:"experiment"`
	SourceUPLM5WSeal                string              `json:"source_up_lm5w_seal"`
	DeadlineProfiles                []string            `json:"deadline_profiles"`
	PooledConditionCells            int                 `json:"pooled_condition_cells"`
	MatchedConditionsPerCell        int                 `json:"matched_conditions_per_cell"`
	TotalConditions                 int                 `json:"total_conditions"`
	TotalDecisionPoints             int                 `json:"total_decision_points"`
	TotalPredicateDifferencePoints  int                 `json:"total_predicate_difference_points"`
	TotalExposedPrehazardOnlyPoints int                 `json:"total_exposed_prehazard_only_points"`
	TotalHazardCandidates           int                 `json:"total_hazard_candidates"`
	TotalPrehazardOnlyCandidates    int                 `json:"total_prehazard_only_candidates"`
	ResourceReductionsUsed          bool                `json:"resource_reductions_used"`
	PrepressureUsed                 bool                `json:"prepressure_used"`
	FutureScheduleUsed              bool                `json:"future_schedule_used"`
	AdaptivePolicySelectionUsed     bool                `json:"adaptive_policy_selection_used"`
	CounterfactualOnly              bool                `json:"counterfactual_only"`
	LiveActivation                  bool                `json:"live_activation"`
	Summaries                       []UPLM5XCellSummary `json:"summaries"`
}

func uplm5xCandidateCounts(arms []uplm2xArm, used map[int]bool, rot int) (hazard, preOnly int) {
	original := map[string]bool{}
	for i := 0; i < 12; i++ {
		original[uplm2nName(i, rot)] = true
	}
	for i := range arms {
		if used[i] || len(arms[i].r.order) == 0 {
			continue
		}
		victim := arms[i].r.order[0]
		if !original[victim] || arms[i].reported[victim] {
			continue
		}
		if len(arms[i].r.order) >= 16 {
			hazard++
		}
		if len(arms[i].r.order) == 15 {
			preOnly++
		}
	}
	return
}

func uplm5xRun(rot int, perm, profile string, budget, start, tp int, w UPLM4WWindow, bred, tred int) (decisionPoints, predicateDiff, exposed, hazardCandidates, preOnlyCandidates int) {
	arms := uplm2yPermute(uplm2xArms(rot), perm)
	uplm3tPrepressure(arms, rot, profile)
	original := map[string]bool{}
	for i := 0; i < 12; i++ {
		original[uplm2nName(i, rot)] = true
	}
	actions := 0
	for round := 0; round < 6; round++ {
		if round >= start {
			cap := uplm4wCap(budget, round, w.BudgetCut, w.BudgetRestore, bred)
			eff := uplm4wTP(tp, round, w.ThroughputCut, w.ThroughputRestore, tred)
			used := map[int]bool{}
			for k := 0; k < eff && actions < cap; k++ {
				h, p := uplm5xCandidateCounts(arms, used, rot)
				decisionPoints++
				hazardCandidates += h
				preOnlyCandidates += p
				if p > 0 {
					predicateDiff++
				}
				if h == 0 && p > 0 {
					exposed++
				}
				i := uplm5uThreatArm(arms, used, rot, false)
				if i < 0 {
					break
				}
				if len(arms[i].r.order) == 0 {
					break
				}
				n := arms[i].r.order[0]
				if !original[n] || arms[i].reported[n] {
					if pn, _, ok := uplm2xFirstPending(&arms[i]); ok {
						n = pn
					} else {
						break
					}
				}
				arms[i].reported[n] = true
				actions++
				used[i] = true
			}
		}
		for i := range arms {
			arms[i].r.write(fmt.Sprintf("5x-%s-%s-%d-%d-%d-%d-%d-%d", profile, perm, bred, tred, budget, start, tp, round), "x")
		}
	}
	return
}

func RunUPLM5X() (UPLM5XResult, error) {
	profiles := []string{"deferred_only", "layout_only", "hybrid_min"}
	windows := []UPLM4WWindow{{1, 3, 2, 4}, {2, 4, 1, 3}, {1, 4, 2, 5}, {2, 5, 1, 4}, {1, 3, 3, 5}, {3, 5, 1, 3}}
	breds := []int{1, 2}
	treds := []int{1, 2}
	rots := []int{5, 13}
	perms := []string{"identity", "reverse", "rotate2"}
	budgets := []int{4, 5, 6, 7}
	starts := []int{2, 3, 4, 5}
	tps := []int{2, 3, 4, 5}
	res := UPLM5XResult{Schema: UPLM5XThreatPredicateSchema, Experiment: "UP-LM5X-threat-predicate-equivalence", SourceUPLM5WSeal: "ec0deabda198a24173a2f95091bf51b0ac817deb", DeadlineProfiles: profiles, PooledConditionCells: 72, MatchedConditionsPerCell: 384, ResourceReductionsUsed: true, PrepressureUsed: true, FutureScheduleUsed: false, AdaptivePolicySelectionUsed: false, CounterfactualOnly: true, LiveActivation: false}
	for _, profile := range profiles {
		for _, br := range breds {
			for _, tr := range treds {
				for _, rot := range rots {
					for _, perm := range perms {
						sm := UPLM5XCellSummary{DeadlineProfile: profile, BudgetReduction: br, ThroughputReduction: tr, Rotation: rot, Permutation: perm}
						for _, w := range windows {
							for _, budget := range budgets {
								for _, start := range starts {
									for _, tp := range tps {
										d, p, e, hc, pc := uplm5xRun(rot, perm, profile, budget, start, tp, w, br, tr)
										sm.Conditions++
										sm.DecisionPoints += d
										sm.PredicateDifferencePoints += p
										sm.ExposedPrehazardOnlyPoints += e
										sm.HazardCandidates += hc
										sm.PrehazardOnlyCandidates += pc
										res.TotalConditions++
										res.TotalDecisionPoints += d
										res.TotalPredicateDifferencePoints += p
										res.TotalExposedPrehazardOnlyPoints += e
										res.TotalHazardCandidates += hc
										res.TotalPrehazardOnlyCandidates += pc
									}
								}
							}
						}
						res.Summaries = append(res.Summaries, sm)
					}
				}
			}
		}
	}
	return res, nil
}
