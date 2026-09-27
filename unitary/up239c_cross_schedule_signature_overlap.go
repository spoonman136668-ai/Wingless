package unitary

const UP239CSignatureOverlapSchema = "wingless.up239c-cross-schedule-signature-overlap.v1"

type UP239CSummary struct {
	ScheduleA int `json:"schedule_a"`
	ScheduleB int `json:"schedule_b"`
	SignatureMode string `json:"signature_mode"`
	AFailureSignatures int `json:"a_failure_signatures"`
	BFailureSignatures int `json:"b_failure_signatures"`
	SharedFailureSignatures int `json:"shared_failure_signatures"`
	ASurvivorSignatures int `json:"a_survivor_signatures"`
	BSurvivorSignatures int `json:"b_survivor_signatures"`
	SharedSurvivorSignatures int `json:"shared_survivor_signatures"`
	AFailureBSurvivorConflicts int `json:"a_failure_b_survivor_conflicts"`
	ASurvivorBFailureConflicts int `json:"a_survivor_b_failure_conflicts"`
}

type UP239CResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP238CSeal string `json:"source_up238c_seal"`
	ArmsTotal int `json:"arms_total"`
	SchedulePairs int `json:"schedule_pairs"`
	SignatureModes []string `json:"signature_modes"`
	InterventionChanged bool `json:"intervention_changed"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	ClassifierTrainingUsed bool `json:"classifier_training_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	NewNativeFieldUsed bool `json:"new_native_field_used"`
	LiveActivation bool `json:"live_activation"`
	Summaries []UP239CSummary `json:"summaries"`
}

type up239cPoint struct {
	schedule int
	arm UP237CArm
}

func up239cCountIntersection(a, b map[string]bool) int {
	n := 0
	for k := range a {
		if b[k] {
			n++
		}
	}
	return n
}

func RunUP239C() (UP239CResult, error) {
	cohorts := [][]int{{0, 9, 10, 15}, {1, 4, 11, 14}, {2, 5, 8, 13}, {3, 6, 7, 12}}
	hands := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	type sched struct{ a, b string }
	schedules := []sched{
		{"no_refresh", "hostile_shield"},
		{"hostile_shield", "no_refresh"},
		{"alternating_shield", "fixed_offset_refresh"},
		{"fixed_offset_refresh", "alternating_shield"},
	}
	points := make([]up239cPoint, 0, 512)
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
						return UP239CResult{}, err
					}
					points = append(points, up239cPoint{schedule: si, arm: a})
				}
			}
		}
	}

	res := UP239CResult{
		Schema: UP239CSignatureOverlapSchema,
		Experiment: "UP-239C-cross-schedule-signature-overlap",
		SourceUP238CSeal: "5de326f66c8de8f754c69ec6d12f63293d2e9168",
		ArmsTotal: len(points),
		SchedulePairs: 6,
		SignatureModes: []string{"current", "trajectory"},
		InterventionChanged: false,
		ThresholdFittingUsed: false,
		ClassifierTrainingUsed: false,
		AdaptiveFeatureSelectionUsed: false,
		NewNativeFieldUsed: false,
		LiveActivation: false,
	}

	for a := 0; a < len(schedules); a++ {
		for b := a + 1; b < len(schedules); b++ {
			for _, mode := range res.SignatureModes {
				af := map[string]bool{}
				as := map[string]bool{}
				bf := map[string]bool{}
				bs := map[string]bool{}
				for _, p := range points {
					if !p.arm.EventPresent || (p.schedule != a && p.schedule != b) {
						continue
					}
					sig := up238cSignature(p.arm, mode)
					if sig == "" {
						continue
					}
					if p.schedule == a {
						if p.arm.Outcome == "loss" {
							af[sig] = true
						} else {
							as[sig] = true
						}
					} else {
						if p.arm.Outcome == "loss" {
							bf[sig] = true
						} else {
							bs[sig] = true
						}
					}
				}
				res.Summaries = append(res.Summaries, UP239CSummary{
					ScheduleA: a,
					ScheduleB: b,
					SignatureMode: mode,
					AFailureSignatures: len(af),
					BFailureSignatures: len(bf),
					SharedFailureSignatures: up239cCountIntersection(af, bf),
					ASurvivorSignatures: len(as),
					BSurvivorSignatures: len(bs),
					SharedSurvivorSignatures: up239cCountIntersection(as, bs),
					AFailureBSurvivorConflicts: up239cCountIntersection(af, bs),
					ASurvivorBFailureConflicts: up239cCountIntersection(as, bf),
				})
			}
		}
	}
	return res, nil
}
