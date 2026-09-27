package unitary

import "fmt"

const UPLM3EDeployabilitySchema = "wingless.up-lm3e-deployability-equivalence.v1"

type UPLM3EPoint struct {
	Profile string `json:"profile"`
	Budget int `json:"budget"`
	Onset int `json:"onset"`
	Throughput int `json:"throughput"`
	PredictedDeployablePerScenario int `json:"predicted_deployable_per_scenario"`
	PredictedDeployableAggregate int `json:"predicted_deployable_aggregate"`
	ActualActions int `json:"actual_actions"`
	BaselineFailed int `json:"baseline_failed"`
	Failed int `json:"failed"`
	Prevented int `json:"prevented"`
}

type UPLM3EActionGroup struct {
	Profile string `json:"profile"`
	ActualActions int `json:"actual_actions"`
	Configurations int `json:"configurations"`
	MinFailed int `json:"min_failed"`
	MaxFailed int `json:"max_failed"`
	FailureSpread int `json:"failure_spread"`
}

type UPLM3EResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM3DSeal string `json:"source_up_lm3d_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	Profiles []string `json:"profiles"`
	Budgets []int `json:"budgets"`
	Onsets []int `json:"onsets"`
	Throughputs []int `json:"throughputs"`
	GlobalRounds int `json:"global_rounds"`
	OneActionPerArmPerRound bool `json:"one_action_per_arm_per_round"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveResourceChoiceUsed bool `json:"adaptive_resource_choice_used"`
	FutureScheduleOracleUsed bool `json:"future_schedule_oracle_used"`
	Points []UPLM3EPoint `json:"points"`
	ActionGroups []UPLM3EActionGroup `json:"action_groups"`
}

func uplm3eMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func uplm3eRun(rot int, perm, profile string, budget, onset, throughput int, active bool) (completed, failed, actions int) {
	arms := uplm2yPermute(uplm2xArms(rot), perm)
	uplm3bPrepressure(arms, profile, rot)
	for round := 0; round < 6; round++ {
		if active && round >= onset {
			used := map[int]bool{}
			for k := 0; k < throughput && actions < budget; k++ {
				i := uplm3cChoose(arms, "earliest_deadline", used)
				if i < 0 {
					break
				}
				used[i] = true
				n, _, ok := uplm2xFirstPending(&arms[i])
				if ok {
					arms[i].reported[n] = true
					actions++
				}
			}
		}
		for i := range arms {
			key := fmt.Sprintf("lm3e-%s-%s-%d-%d-%d-%d-%d", profile, perm, budget, onset, throughput, rot, round)
			arms[i].r.write(key, "x")
		}
	}
	for i := range arms {
		c, f := uplm2vFinish(arms[i].r, arms[i].reported, rot)
		completed += c
		failed += f
	}
	return
}

func RunUPLM3E() (UPLM3EResult, error) {
	profiles := []string{"by_deferred_level", "by_layout"}
	budgets := []int{3, 4, 5, 6}
	onsets := []int{0, 1, 2, 3}
	throughputs := []int{1, 2, 3}
	rots := []int{3, 11}
	perms := []string{"identity", "reverse"}

	res := UPLM3EResult{
		Schema: UPLM3EDeployabilitySchema,
		Experiment: "UP-LM3E-deployability-equivalence",
		SourceUPLM3DSeal: "ec7a600cbea615c17edc3bd6dcf79961e9fa2c4c",
		ExactRecallCap: 16,
		Profiles: profiles,
		Budgets: budgets,
		Onsets: onsets,
		Throughputs: throughputs,
		GlobalRounds: 6,
		OneActionPerArmPerRound: true,
		CounterfactualOnly: true,
		LiveActivation: false,
		AdaptiveResourceChoiceUsed: false,
		FutureScheduleOracleUsed: false,
	}

	type groupKey struct {
		profile string
		actions int
	}
	type groupAcc struct {
		n int
		min int
		max int
	}
	groups := map[groupKey]*groupAcc{}

	for _, profile := range profiles {
		for _, budget := range budgets {
			for _, onset := range onsets {
				for _, tp := range throughputs {
					baselineFailed := 0
					failed := 0
					actions := 0
					for _, rot := range rots {
						for _, perm := range perms {
							_, bf, _ := uplm3eRun(rot, perm, profile, budget, onset, tp, false)
							_, f, a := uplm3eRun(rot, perm, profile, budget, onset, tp, true)
							baselineFailed += bf
							failed += f
							actions += a
						}
					}
					perScenario := uplm3eMin(budget, (6-onset)*tp)
					point := UPLM3EPoint{
						Profile: profile,
						Budget: budget,
						Onset: onset,
						Throughput: tp,
						PredictedDeployablePerScenario: perScenario,
						PredictedDeployableAggregate: perScenario * len(rots) * len(perms),
						ActualActions: actions,
						BaselineFailed: baselineFailed,
						Failed: failed,
						Prevented: baselineFailed - failed,
					}
					res.Points = append(res.Points, point)
					key := groupKey{profile: profile, actions: actions}
					g := groups[key]
					if g == nil {
						g = &groupAcc{min: failed, max: failed}
						groups[key] = g
					}
					g.n++
					if failed < g.min {
						g.min = failed
					}
					if failed > g.max {
						g.max = failed
					}
				}
			}
		}
	}
	for _, profile := range profiles {
		for actions := 0; actions <= 24; actions++ {
			g := groups[groupKey{profile: profile, actions: actions}]
			if g == nil {
				continue
			}
			res.ActionGroups = append(res.ActionGroups, UPLM3EActionGroup{
				Profile: profile,
				ActualActions: actions,
				Configurations: g.n,
				MinFailed: g.min,
				MaxFailed: g.max,
				FailureSpread: g.max - g.min,
			})
		}
	}
	return res, nil
}
