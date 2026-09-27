package unitary

import "fmt"

const UPLM3DDelayedOnsetSchema = "wingless.up-lm3d-delayed-onset-throughput.v1"

type UPLM3DResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM3CSeal string `json:"source_up_lm3c_seal"`
	ExactRecallCap int `json:"exact_recall_cap"`
	Profiles []string `json:"profiles"`
	IdentityRotations []int `json:"identity_rotations"`
	Permutations []string `json:"permutations"`
	Throughputs []int `json:"throughputs"`
	TotalBudget int `json:"total_budget"`
	GlobalRounds int `json:"global_rounds"`
	InterventionStartRound int `json:"intervention_start_round"`
	OneActionPerArmPerRound bool `json:"one_action_per_arm_per_round"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveOnsetUsed bool `json:"adaptive_onset_used"`
	AdaptiveBudgetUsed bool `json:"adaptive_budget_used"`
	AdaptiveThroughputUsed bool `json:"adaptive_throughput_used"`
	Points []UPLM3BPoint `json:"points"`
	Metrics []UPLM3BMetric `json:"metrics"`
}

func uplm3dRun(rot int, perm, profile, policy string, throughput int) (completed, failed, actions int) {
	arms := uplm2yPermute(uplm2xArms(rot), perm)
	uplm3bPrepressure(arms, profile, rot)
	const budget = 6
	for round := 0; round < 6; round++ {
		if policy != "baseline" && round >= 3 {
			used := map[int]bool{}
			for k := 0; k < throughput && actions < budget; k++ {
				i := uplm3cChoose(arms, policy, used)
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
			key := fmt.Sprintf("lm3d-%s-%s-%s-%d-%d-%d-%d", profile, policy, perm, throughput, rot, round, i)
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

func RunUPLM3D() (UPLM3DResult, error) {
	profiles := []string{"by_deferred_level", "by_layout"}
	rots := []int{3, 11}
	perms := []string{"identity", "reverse"}
	tps := []int{1, 2, 3}
	res := UPLM3DResult{
		Schema: UPLM3DDelayedOnsetSchema,
		Experiment: "UP-LM3D-delayed-onset-throughput",
		SourceUPLM3CSeal: "bb29b01bf50424f94cc40412fd395059d3b0cc48",
		ExactRecallCap: 16,
		Profiles: profiles,
		IdentityRotations: rots,
		Permutations: perms,
		Throughputs: tps,
		TotalBudget: 6,
		GlobalRounds: 6,
		InterventionStartRound: 4,
		OneActionPerArmPerRound: true,
		CounterfactualOnly: true,
		LiveActivation: false,
		AdaptiveOnsetUsed: false,
		AdaptiveBudgetUsed: false,
		AdaptiveThroughputUsed: false,
	}
	for _, profile := range profiles {
		for _, tp := range tps {
			m := UPLM3BMetric{Profile: profile, Throughput: tp}
			for _, rot := range rots {
				for _, perm := range perms {
					bc, bf, _ := uplm3dRun(rot, perm, profile, "baseline", tp)
					ec, ef, ea := uplm3dRun(rot, perm, profile, "earliest_deadline", tp)
					fc, ff, fa := uplm3dRun(rot, perm, profile, "fixed_order", tp)
					res.Points = append(res.Points,
						UPLM3BPoint{Profile: profile, IdentityRotation: rot, Permutation: perm, Throughput: tp, Policy: "baseline", Completed: bc, Failed: bf},
						UPLM3BPoint{Profile: profile, IdentityRotation: rot, Permutation: perm, Throughput: tp, Policy: "earliest_deadline", Actions: ea, Completed: ec, Failed: ef},
						UPLM3BPoint{Profile: profile, IdentityRotation: rot, Permutation: perm, Throughput: tp, Policy: "fixed_order", Actions: fa, Completed: fc, Failed: ff},
					)
					m.BaselineFailed += bf
					m.EarliestFailed += ef
					m.FixedFailed += ff
					m.EarliestActions += ea
					m.FixedActions += fa
				}
			}
			m.EarliestPrevented = m.BaselineFailed - m.EarliestFailed
			m.FixedPrevented = m.BaselineFailed - m.FixedFailed
			m.EarliestAdvantage = m.FixedFailed - m.EarliestFailed
			res.Metrics = append(res.Metrics, m)
		}
	}
	return res, nil
}
