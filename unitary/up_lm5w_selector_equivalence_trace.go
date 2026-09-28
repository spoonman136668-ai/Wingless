package unitary

import (
	"fmt"
	"strings"
)

const UPLM5WSelectorEquivalenceSchema = "wingless.up-lm5w-selector-equivalence-trace.v1"

type UPLM5WCellSummary struct {
	DeadlineProfile     string `json:"deadline_profile"`
	BudgetReduction     int    `json:"budget_reduction"`
	ThroughputReduction int    `json:"throughput_reduction"`
	Rotation            int    `json:"rotation"`
	Permutation         string `json:"permutation"`
	ComparedConditions  int    `json:"compared_conditions"`
	SelectorDifferences int    `json:"selector_differences"`
	ActionDifferences   int    `json:"action_differences"`
	FailureDifferences  int    `json:"failure_differences"`
}

type UPLM5WResult struct {
	Schema                         string                 `json:"schema"`
	Experiment                     string                 `json:"experiment"`
	SourceUPLM5VSeal               string                 `json:"source_up_lm5v_seal"`
	ComparedPolicies               []string               `json:"compared_policies"`
	DeadlineProfiles               []string               `json:"deadline_profiles"`
	PooledConditionCells           int                    `json:"pooled_condition_cells"`
	MatchedConditionsPerCell       int                    `json:"matched_conditions_per_cell"`
	TotalComparedConditions        int                    `json:"total_compared_conditions"`
	TotalSelectorDifferences       int                    `json:"total_selector_differences"`
	TotalActionDifferences         int                    `json:"total_action_differences"`
	TotalFailureDifferences        int                    `json:"total_failure_differences"`
	ResourceReductionsUsed         bool                   `json:"resource_reductions_used"`
	PrepressureUsed                bool                   `json:"prepressure_used"`
	PrehazardRuleFixed             bool                   `json:"prehazard_rule_fixed"`
	FutureScheduleUsed             bool                   `json:"future_schedule_used"`
	AdaptivePolicySelectionUsed    bool                   `json:"adaptive_policy_selection_used"`
	CounterfactualOnly             bool                   `json:"counterfactual_only"`
	LiveActivation                 bool                   `json:"live_activation"`
	Summaries                      []UPLM5WCellSummary     `json:"summaries"`
}

func uplm5wRun(rot int, perm, profile, policy string, budget, start, tp int, w UPLM4WWindow, bred, tred int) (actions, failed int, trace string) {
	arms := uplm2yPermute(uplm2xArms(rot), perm)
	uplm3tPrepressure(arms, rot, profile)

	original := map[string]bool{}
	for i := 0; i < 12; i++ {
		original[uplm2nName(i, rot)] = true
	}

	events := []string{}
	for round := 0; round < 6; round++ {
		if round >= start {
			cap := uplm4wCap(budget, round, w.BudgetCut, w.BudgetRestore, bred)
			eff := uplm4wTP(tp, round, w.ThroughputCut, w.ThroughputRestore, tred)
			used := map[int]bool{}
			for k := 0; k < eff && actions < cap; k++ {
				i := uplm5uThreatArm(arms, used, rot, policy == "prehazard_triggered")
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
				events = append(events, fmt.Sprintf("r%d:a%d:%s", round, i, n))
				arms[i].reported[n] = true
				actions++
				used[i] = true
			}
		}
		for i := range arms {
			arms[i].r.write(fmt.Sprintf("5w-%s-%s-%s-%d-%d-%d-%d-%d-%d", profile, policy, perm, bred, tred, budget, start, tp, round), "x")
		}
	}
	for i := range arms {
		_, f := uplm2vFinish(arms[i].r, arms[i].reported, rot)
		failed += f
	}
	return actions, failed, strings.Join(events, "|")
}

func RunUPLM5W() (UPLM5WResult, error) {
	profiles := []string{"deferred_only", "layout_only", "hybrid_min"}
	windows := []UPLM4WWindow{{1, 3, 2, 4}, {2, 4, 1, 3}, {1, 4, 2, 5}, {2, 5, 1, 4}, {1, 3, 3, 5}, {3, 5, 1, 3}}
	breds := []int{1, 2}
	treds := []int{1, 2}
	rots := []int{5, 13}
	perms := []string{"identity", "reverse", "rotate2"}
	budgets := []int{4, 5, 6, 7}
	starts := []int{2, 3, 4, 5}
	tps := []int{2, 3, 4, 5}

	res := UPLM5WResult{
		Schema:                      UPLM5WSelectorEquivalenceSchema,
		Experiment:                  "UP-LM5W-selector-equivalence-trace",
		SourceUPLM5VSeal:            "1d7179582dfe06563c6b630cb6adb9c8090421d5",
		ComparedPolicies:            []string{"hazard_triggered", "prehazard_triggered"},
		DeadlineProfiles:            profiles,
		PooledConditionCells:        72,
		MatchedConditionsPerCell:    384,
		ResourceReductionsUsed:      true,
		PrepressureUsed:             true,
		PrehazardRuleFixed:          true,
		FutureScheduleUsed:          false,
		AdaptivePolicySelectionUsed: false,
		CounterfactualOnly:          true,
		LiveActivation:              false,
	}

	for _, profile := range profiles {
		for _, br := range breds {
			for _, tr := range treds {
				for _, rot := range rots {
					for _, perm := range perms {
						sm := UPLM5WCellSummary{
							DeadlineProfile:     profile,
							BudgetReduction:     br,
							ThroughputReduction: tr,
							Rotation:            rot,
							Permutation:         perm,
						}
						for _, w := range windows {
							for _, budget := range budgets {
								for _, start := range starts {
									for _, tp := range tps {
										ha, hf, ht := uplm5wRun(rot, perm, profile, "hazard_triggered", budget, start, tp, w, br, tr)
										pa, pf, pt := uplm5wRun(rot, perm, profile, "prehazard_triggered", budget, start, tp, w, br, tr)
										sm.ComparedConditions++
										res.TotalComparedConditions++
										if ht != pt {
											sm.SelectorDifferences++
											res.TotalSelectorDifferences++
										}
										if ha != pa {
											sm.ActionDifferences++
											res.TotalActionDifferences++
										}
										if hf != pf {
											sm.FailureDifferences++
											res.TotalFailureDifferences++
										}
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
