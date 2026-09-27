package unitary

import "strings"

const UP240CComponentOverlapSchema = "wingless.up240c-signature-component-overlap.v1"

type UP240CSummary struct {
	ScheduleA int `json:"schedule_a"`
	ScheduleB int `json:"schedule_b"`
	SignatureMode string `json:"signature_mode"`
	Component string `json:"component"`
	AFailureValues int `json:"a_failure_values"`
	BFailureValues int `json:"b_failure_values"`
	SharedFailureValues int `json:"shared_failure_values"`
	ASurvivorValues int `json:"a_survivor_values"`
	BSurvivorValues int `json:"b_survivor_values"`
	SharedSurvivorValues int `json:"shared_survivor_values"`
	AFailureBSurvivorConflicts int `json:"a_failure_b_survivor_conflicts"`
	ASurvivorBFailureConflicts int `json:"a_survivor_b_failure_conflicts"`
}

type UP240CResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP239CSeal string `json:"source_up239c_seal"`
	ArmsTotal int `json:"arms_total"`
	SchedulePairs int `json:"schedule_pairs"`
	CurrentComponents []string `json:"current_components"`
	TrajectoryComponents []string `json:"trajectory_components"`
	InterventionChanged bool `json:"intervention_changed"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	ClassifierTrainingUsed bool `json:"classifier_training_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	NewNativeFieldUsed bool `json:"new_native_field_used"`
	LiveActivation bool `json:"live_activation"`
	Summaries []UP240CSummary `json:"summaries"`
}

type up240cPoint struct {
	schedule int
	arm UP237CArm
}

func up240cFields(sig string) map[string]string {
	out := map[string]string{}
	for _, part := range strings.Split(sig, "|") {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			out[kv[0]] = kv[1]
		}
	}
	return out
}

func up240cIntersect(a, b map[string]bool) int {
	n := 0
	for k := range a {
		if b[k] {
			n++
		}
	}
	return n
}

func RunUP240C() (UP240CResult, error) {
	cohorts := [][]int{{0, 9, 10, 15}, {1, 4, 11, 14}, {2, 5, 8, 13}, {3, 6, 7, 12}}
	hands := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	type sched struct{ a, b string }
	schedules := []sched{
		{"no_refresh", "hostile_shield"},
		{"hostile_shield", "no_refresh"},
		{"alternating_shield", "fixed_offset_refresh"},
		{"fixed_offset_refresh", "alternating_shield"},
	}
	points := make([]up240cPoint, 0, 512)
	for si, s := range schedules {
		for _, adv := range []int{0, 22} {
			for _, c := range cohorts {
				for _, hand := range hands {
					x, e, ok := up161cTriggerState(hand, c)
					if !ok {
						continue
					}
					a, err := up237cRun(x, e, adv, s.a, s.b)
					if err != nil {
						return UP240CResult{}, err
					}
					points = append(points, up240cPoint{schedule: si, arm: a})
				}
			}
		}
	}

	current := []string{"age", "dist", "pred", "a0", "a1", "a2", "a3", "adv", "noq"}
	trajectory := []string{"age", "dist", "pred", "a0", "a1", "a2", "a3", "adv", "noq", "dage", "ddist", "dpred", "da0", "da1", "da2", "da3", "dadv", "dnoq"}
	res := UP240CResult{
		Schema: UP240CComponentOverlapSchema,
		Experiment: "UP-240C-signature-component-overlap",
		SourceUP239CSeal: "7c17964ae4b864f0ffdd01b6bca9cbdbd5aa4a01",
		ArmsTotal: len(points),
		SchedulePairs: 6,
		CurrentComponents: current,
		TrajectoryComponents: trajectory,
		InterventionChanged: false,
		ThresholdFittingUsed: false,
		ClassifierTrainingUsed: false,
		AdaptiveFeatureSelectionUsed: false,
		NewNativeFieldUsed: false,
		LiveActivation: false,
	}

	for a := 0; a < len(schedules); a++ {
		for b := a + 1; b < len(schedules); b++ {
			for _, mode := range []string{"current", "trajectory"} {
				components := current
				if mode == "trajectory" {
					components = trajectory
				}
				for _, component := range components {
					af := map[string]bool{}
					as := map[string]bool{}
					bf := map[string]bool{}
					bs := map[string]bool{}
					for _, p := range points {
						if !p.arm.EventPresent || (p.schedule != a && p.schedule != b) {
							continue
						}
						sig := p.arm.CurrentSignature
						if mode == "trajectory" {
							sig = p.arm.TrajectorySignature
						}
						if sig == "" {
							continue
						}
						value, ok := up240cFields(sig)[component]
						if !ok {
							continue
						}
						target := af
						if p.schedule == a {
							if p.arm.Outcome != "loss" {
								target = as
							}
						} else {
							target = bf
							if p.arm.Outcome != "loss" {
								target = bs
							}
						}
						target[value] = true
					}
					res.Summaries = append(res.Summaries, UP240CSummary{
						ScheduleA: a,
						ScheduleB: b,
						SignatureMode: mode,
						Component: component,
						AFailureValues: len(af),
						BFailureValues: len(bf),
						SharedFailureValues: up240cIntersect(af, bf),
						ASurvivorValues: len(as),
						BSurvivorValues: len(bs),
						SharedSurvivorValues: up240cIntersect(as, bs),
						AFailureBSurvivorConflicts: up240cIntersect(af, bs),
						ASurvivorBFailureConflicts: up240cIntersect(as, bf),
					})
				}
			}
		}
	}
	return res, nil
}
