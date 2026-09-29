package unitary

import (
	"fmt"
	"sort"
	"strconv"
)

const UPLM5YOrderLengthSchema = "wingless.up-lm5y-eligible-order-length-spectrum.v1"

type UPLM5YCellSummary struct {
	DeadlineProfile                 string         `json:"deadline_profile"`
	BudgetReduction                 int            `json:"budget_reduction"`
	ThroughputReduction             int            `json:"throughput_reduction"`
	Rotation                        int            `json:"rotation"`
	Permutation                     string         `json:"permutation"`
	Conditions                      int            `json:"conditions"`
	DecisionPoints                  int            `json:"decision_points"`
	EligibleCandidates              int            `json:"eligible_candidates"`
	HazardCandidates                int            `json:"hazard_candidates"`
	SubhazardCandidates             int            `json:"subhazard_candidates"`
	ExposedSubhazardDecisionPoints  int            `json:"exposed_subhazard_decision_points"`
	LengthHistogram                 map[string]int `json:"length_histogram"`
}

type UPLM5YResult struct {
	Schema                           string               `json:"schema"`
	Experiment                       string               `json:"experiment"`
	SourceUPLM5XSeal                 string               `json:"source_up_lm5x_seal"`
	DeadlineProfiles                 []string             `json:"deadline_profiles"`
	PooledConditionCells             int                  `json:"pooled_condition_cells"`
	MatchedConditionsPerCell         int                  `json:"matched_conditions_per_cell"`
	TotalConditions                  int                  `json:"total_conditions"`
	TotalDecisionPoints              int                  `json:"total_decision_points"`
	TotalEligibleCandidates          int                  `json:"total_eligible_candidates"`
	TotalHazardCandidates            int                  `json:"total_hazard_candidates"`
	TotalSubhazardCandidates         int                  `json:"total_subhazard_candidates"`
	TotalExposedSubhazardPoints      int                  `json:"total_exposed_subhazard_decision_points"`
	ObservedLengths                  []int                `json:"observed_lengths"`
	LengthHistogram                  map[string]int       `json:"length_histogram"`
	ParentPrehazardOnlyCandidates    int                  `json:"parent_prehazard_only_candidates"`
	ParentPredicateDifferencePoints  int                  `json:"parent_predicate_difference_points"`
	CanonicalHazardSelectorUsed      bool                 `json:"canonical_hazard_selector_used"`
	PolicyChanged                    bool                 `json:"policy_changed"`
	FutureScheduleUsed               bool                 `json:"future_schedule_used"`
	AdaptivePolicySelectionUsed      bool                 `json:"adaptive_policy_selection_used"`
	CounterfactualOnly               bool                 `json:"counterfactual_only"`
	LiveActivation                   bool                 `json:"live_activation"`
	Summaries                        []UPLM5YCellSummary  `json:"summaries"`
}

func uplm5yEligibleLengths(arms []uplm2xArm, used map[int]bool, rot int) []int {
	original := map[string]bool{}
	for i := 0; i < 12; i++ {
		original[uplm2nName(i, rot)] = true
	}
	out := []int{}
	for i := range arms {
		if used[i] || len(arms[i].r.order) == 0 {
			continue
		}
		victim := arms[i].r.order[0]
		if !original[victim] || arms[i].reported[victim] {
			continue
		}
		out = append(out, len(arms[i].r.order))
	}
	return out
}

func uplm5yRun(rot int, perm, profile string, budget, start, tp int, w UPLM4WWindow, bred, tred int) (decisionPoints, eligible, hazard, subhazard, exposed int, hist map[int]int) {
	arms := uplm2yPermute(uplm2xArms(rot), perm)
	uplm3tPrepressure(arms, rot, profile)
	original := map[string]bool{}
	for i := 0; i < 12; i++ {
		original[uplm2nName(i, rot)] = true
	}
	hist = map[int]int{}
	actions := 0
	for round := 0; round < 6; round++ {
		if round >= start {
			cap := uplm4wCap(budget, round, w.BudgetCut, w.BudgetRestore, bred)
			eff := uplm4wTP(tp, round, w.ThroughputCut, w.ThroughputRestore, tred)
			used := map[int]bool{}
			for k := 0; k < eff && actions < cap; k++ {
				lengths := uplm5yEligibleLengths(arms, used, rot)
				decisionPoints++
				h := 0
				s := 0
				for _, n := range lengths {
					hist[n]++
					eligible++
					if n >= 16 {
						h++
						hazard++
					} else {
						s++
						subhazard++
					}
				}
				if h == 0 && s > 0 {
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
			arms[i].r.write(fmt.Sprintf("5y-%s-%s-%d-%d-%d-%d-%d-%d", profile, perm, bred, tred, budget, start, tp, round), "x")
		}
	}
	return
}

func RunUPLM5Y() (UPLM5YResult, error) {
	parent, err := RunUPLM5X()
	if err != nil {
		return UPLM5YResult{}, err
	}
	profiles := []string{"deferred_only", "layout_only", "hybrid_min"}
	windows := []UPLM4WWindow{{1, 3, 2, 4}, {2, 4, 1, 3}, {1, 4, 2, 5}, {2, 5, 1, 4}, {1, 3, 3, 5}, {3, 5, 1, 3}}
	breds := []int{1, 2}
	treds := []int{1, 2}
	rots := []int{5, 13}
	perms := []string{"identity", "reverse", "rotate2"}
	budgets := []int{4, 5, 6, 7}
	starts := []int{2, 3, 4, 5}
	tps := []int{2, 3, 4, 5}
	res := UPLM5YResult{
		Schema: UPLM5YOrderLengthSchema,
		Experiment: "UP-LM5Y-eligible-order-length-spectrum",
		SourceUPLM5XSeal: "6dfb16f07df2281f5d0e636ae14536e01221111c",
		DeadlineProfiles: profiles,
		PooledConditionCells: 72,
		MatchedConditionsPerCell: 384,
		ParentPrehazardOnlyCandidates: parent.TotalPrehazardOnlyCandidates,
		ParentPredicateDifferencePoints: parent.TotalPredicateDifferencePoints,
		CanonicalHazardSelectorUsed: true,
		PolicyChanged: false,
		FutureScheduleUsed: false,
		AdaptivePolicySelectionUsed: false,
		CounterfactualOnly: true,
		LiveActivation: false,
		LengthHistogram: map[string]int{},
	}
	globalHist := map[int]int{}
	for _, profile := range profiles {
		for _, br := range breds {
			for _, tr := range treds {
				for _, rot := range rots {
					for _, perm := range perms {
						sm := UPLM5YCellSummary{DeadlineProfile: profile, BudgetReduction: br, ThroughputReduction: tr, Rotation: rot, Permutation: perm, LengthHistogram: map[string]int{}}
						cellHist := map[int]int{}
						for _, w := range windows {
							for _, budget := range budgets {
								for _, start := range starts {
									for _, tp := range tps {
										d, e, h, s, x, hist := uplm5yRun(rot, perm, profile, budget, start, tp, w, br, tr)
										sm.Conditions++
										sm.DecisionPoints += d
										sm.EligibleCandidates += e
										sm.HazardCandidates += h
										sm.SubhazardCandidates += s
										sm.ExposedSubhazardDecisionPoints += x
										res.TotalConditions++
										res.TotalDecisionPoints += d
										res.TotalEligibleCandidates += e
										res.TotalHazardCandidates += h
										res.TotalSubhazardCandidates += s
										res.TotalExposedSubhazardPoints += x
										for n, count := range hist {
											cellHist[n] += count
											globalHist[n] += count
										}
									}
								}
							}
						}
						keys := make([]int, 0, len(cellHist))
						for n := range cellHist {
							keys = append(keys, n)
						}
						sort.Ints(keys)
						for _, n := range keys {
							sm.LengthHistogram[strconv.Itoa(n)] = cellHist[n]
						}
						res.Summaries = append(res.Summaries, sm)
					}
				}
			}
		}
	}
	keys := make([]int, 0, len(globalHist))
	for n := range globalHist {
		keys = append(keys, n)
	}
	sort.Ints(keys)
	res.ObservedLengths = append(res.ObservedLengths, keys...)
	for _, n := range keys {
		res.LengthHistogram[strconv.Itoa(n)] = globalHist[n]
	}
	return res, nil
}
