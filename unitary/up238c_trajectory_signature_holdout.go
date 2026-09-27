package unitary

const UP238CTrajectoryHoldoutSchema = "wingless.up238c-trajectory-signature-holdout.v1"

type UP238CSummary struct {
	HoldoutSchedule int `json:"holdout_schedule"`
	SignatureMode string `json:"signature_mode"`
	TrainFailureSignatures int `json:"train_failure_signatures"`
	TrainSurvivorSignatures int `json:"train_survivor_signatures"`
	TrainSharedSignatures int `json:"train_shared_signatures"`
	TestEventArms int `json:"test_event_arms"`
	TestFailures int `json:"test_failures"`
	TestSurvivors int `json:"test_survivors"`
	TruePositive int `json:"true_positive"`
	FalseNegative int `json:"false_negative"`
	FalsePositive int `json:"false_positive"`
	TrueNegative int `json:"true_negative"`
	HoldoutSilentFailures int `json:"holdout_silent_failures"`
	Sensitivity float64 `json:"sensitivity"`
	Specificity float64 `json:"specificity"`
}

type UP238CResult struct {
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUP237CSeal string `json:"source_up237c_seal"`
	ArmsTotal int `json:"arms_total"`
	HoldoutFolds int `json:"holdout_folds"`
	SignatureModes []string `json:"signature_modes"`
	InterventionChanged bool `json:"intervention_changed"`
	ThresholdFittingUsed bool `json:"threshold_fitting_used"`
	ClassifierTrainingUsed bool `json:"classifier_training_used"`
	SignatureLookupUsed bool `json:"signature_lookup_used"`
	AdaptiveFeatureSelectionUsed bool `json:"adaptive_feature_selection_used"`
	HoldoutScheduleExcludedFromTraining bool `json:"holdout_schedule_excluded_from_training"`
	NewNativeFieldUsed bool `json:"new_native_field_used"`
	LiveActivation bool `json:"live_activation"`
	Summaries []UP238CSummary `json:"summaries"`
}

type up238cPoint struct {
	schedule int
	arm UP237CArm
}

func up238cSignature(a UP237CArm, mode string) string {
	if mode == "trajectory" {
		return a.TrajectorySignature
	}
	return a.CurrentSignature
}

func RunUP238C() (UP238CResult, error) {
	cohorts := [][]int{{0, 9, 10, 15}, {1, 4, 11, 14}, {2, 5, 8, 13}, {3, 6, 7, 12}}
	hands := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	type sched struct{ a, b string }
	schedules := []sched{
		{"no_refresh", "hostile_shield"},
		{"hostile_shield", "no_refresh"},
		{"alternating_shield", "fixed_offset_refresh"},
		{"fixed_offset_refresh", "alternating_shield"},
	}
	points := make([]up238cPoint, 0, 512)
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
						return UP238CResult{}, err
					}
					points = append(points, up238cPoint{schedule: si, arm: a})
				}
			}
		}
	}
	r := UP238CResult{
		Schema: UP238CTrajectoryHoldoutSchema,
		Experiment: "UP-238C-trajectory-signature-holdout",
		SourceUP237CSeal: "cf01d1e9193c77e286b01578649d39e252ab518e",
		ArmsTotal: len(points),
		HoldoutFolds: len(schedules),
		SignatureModes: []string{"current", "trajectory"},
		InterventionChanged: false,
		ThresholdFittingUsed: false,
		ClassifierTrainingUsed: false,
		SignatureLookupUsed: true,
		AdaptiveFeatureSelectionUsed: false,
		HoldoutScheduleExcludedFromTraining: true,
		NewNativeFieldUsed: false,
		LiveActivation: false,
	}
	for hold := range schedules {
		for _, mode := range r.SignatureModes {
			failSet := map[string]bool{}
			surviveSet := map[string]bool{}
			for _, p := range points {
				if p.schedule == hold || !p.arm.EventPresent {
					continue
				}
				sig := up238cSignature(p.arm, mode)
				if sig == "" {
					continue
				}
				if p.arm.Outcome == "loss" {
					failSet[sig] = true
				} else {
					surviveSet[sig] = true
				}
			}
			s := UP238CSummary{
				HoldoutSchedule: hold,
				SignatureMode: mode,
				TrainFailureSignatures: len(failSet),
				TrainSurvivorSignatures: len(surviveSet),
			}
			for sig := range failSet {
				if surviveSet[sig] {
					s.TrainSharedSignatures++
				}
			}
			for _, p := range points {
				if p.schedule != hold {
					continue
				}
				if p.arm.Outcome == "loss" && !p.arm.EventPresent {
					s.HoldoutSilentFailures++
				}
				if !p.arm.EventPresent {
					continue
				}
				sig := up238cSignature(p.arm, mode)
				if sig == "" {
					continue
				}
				pred := failSet[sig] && !surviveSet[sig]
				s.TestEventArms++
				if p.arm.Outcome == "loss" {
					s.TestFailures++
					if pred {
						s.TruePositive++
					} else {
						s.FalseNegative++
					}
				} else {
					s.TestSurvivors++
					if pred {
						s.FalsePositive++
					} else {
						s.TrueNegative++
					}
				}
			}
			if s.TestFailures > 0 {
				s.Sensitivity = float64(s.TruePositive) / float64(s.TestFailures)
			}
			if s.TestSurvivors > 0 {
				s.Specificity = float64(s.TrueNegative) / float64(s.TestSurvivors)
			}
			r.Summaries = append(r.Summaries, s)
		}
	}
	return r, nil
}
