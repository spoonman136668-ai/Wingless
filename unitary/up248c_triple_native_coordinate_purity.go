package unitary

import (
	"fmt"
	"sort"
)

const UP248CTriplePuritySchema = "wingless.up248c-triple-native-coordinate-purity.v1"

type UP248CTripleSummary struct {
	CoordinateA        string `json:"coordinate_a"`
	CoordinateB        string `json:"coordinate_b"`
	CoordinateC        string `json:"coordinate_c"`
	Groups             int    `json:"groups"`
	PureGroups         int    `json:"pure_groups"`
	MixedGroups        int    `json:"mixed_groups"`
	MaxClassesPerGroup int    `json:"max_classes_per_group"`
}

type UP248CResult struct {
	Schema                       string                `json:"schema"`
	Experiment                   string                `json:"experiment"`
	SourceUP247CSeal             string                `json:"source_up247c_seal"`
	ConditionA                   string                `json:"condition_a"`
	ConditionB                   string                `json:"condition_b"`
	PairedInitialStates          int                   `json:"paired_initial_states"`
	Coordinates                  []string              `json:"coordinates"`
	CoordinateTriples            int                   `json:"coordinate_triples"`
	OutcomeClasses               []string              `json:"outcome_classes"`
	InterventionChanged          bool                  `json:"intervention_changed"`
	ThresholdFittingUsed         bool                  `json:"threshold_fitting_used"`
	ClassifierTrainingUsed       bool                  `json:"classifier_training_used"`
	AdaptiveFeatureSelectionUsed bool                  `json:"adaptive_feature_selection_used"`
	NewNativeFieldUsed           bool                  `json:"new_native_field_used"`
	LiveActivation               bool                  `json:"live_activation"`
	Summaries                    []UP248CTripleSummary `json:"summaries"`
}

func RunUP248C() (UP248CResult, error) {
	cohorts := [][]int{{0, 9, 10, 15}, {1, 4, 11, 14}, {2, 5, 8, 13}, {3, 6, 7, 12}}
	hands := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	coords := []string{"endangered_age", "predicted_is_endangered", "hand_distance", "age0", "age1", "age2", "age3", "adversarial_horizon", "no_query_horizon"}
	res := UP248CResult{
		Schema: UP248CTriplePuritySchema, Experiment: "UP-248C-triple-native-coordinate-purity", SourceUP247CSeal: "763c6a84d701cbd5f31bf64f2411b0ed644021c1",
		ConditionA: "schedule2|advance22", ConditionB: "schedule3|advance0", Coordinates: coords,
		OutcomeClasses: []string{"a_only_fail", "b_only_fail", "both_fail", "neither_fail"},
		InterventionChanged: false, ThresholdFittingUsed: false, ClassifierTrainingUsed: false, AdaptiveFeatureSelectionUsed: false, NewNativeFieldUsed: false, LiveActivation: false,
	}
	type sample struct {
		values map[string]string
		class  string
	}
	samples := []sample{}
	for _, cohort := range cohorts {
		for _, hand := range hands {
			x, e, ok := up161cTriggerState(hand, cohort)
			if !ok {
				continue
			}
			s := up223cSnapshot(x, e, 0)
			a, err := up237cRun(x, e, 22, "alternating_shield", "fixed_offset_refresh")
			if err != nil {
				return UP248CResult{}, err
			}
			b, err := up237cRun(x, e, 0, "fixed_offset_refresh", "alternating_shield")
			if err != nil {
				return UP248CResult{}, err
			}
			res.PairedInitialStates++
			samples = append(samples, sample{values: up246cValues(s), class: up246cClass(a, b)})
		}
	}
	for i := 0; i < len(coords); i++ {
		for j := i + 1; j < len(coords); j++ {
			for k := j + 1; k < len(coords); k++ {
				a, b, c := coords[i], coords[j], coords[k]
				groups := map[string]map[string]bool{}
				for _, s := range samples {
					key := fmt.Sprintf("%s=%s|%s=%s|%s=%s", a, s.values[a], b, s.values[b], c, s.values[c])
					if groups[key] == nil {
						groups[key] = map[string]bool{}
					}
					groups[key][s.class] = true
				}
				keys := make([]string, 0, len(groups))
				for key := range groups {
					keys = append(keys, key)
				}
				sort.Strings(keys)
				sm := UP248CTripleSummary{CoordinateA: a, CoordinateB: b, CoordinateC: c, Groups: len(groups)}
				for _, key := range keys {
					n := len(groups[key])
					if n == 1 {
						sm.PureGroups++
					} else {
						sm.MixedGroups++
					}
					if n > sm.MaxClassesPerGroup {
						sm.MaxClassesPerGroup = n
					}
				}
				res.Summaries = append(res.Summaries, sm)
			}
		}
	}
	res.CoordinateTriples = len(res.Summaries)
	return res, nil
}
