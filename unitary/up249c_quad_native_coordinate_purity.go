package unitary

import (
	"fmt"
	"sort"
)

const UP249CQuadPuritySchema = "wingless.up249c-quad-native-coordinate-purity.v1"

type UP249CQuadSummary struct {
	CoordinateA        string `json:"coordinate_a"`
	CoordinateB        string `json:"coordinate_b"`
	CoordinateC        string `json:"coordinate_c"`
	CoordinateD        string `json:"coordinate_d"`
	Groups             int    `json:"groups"`
	PureGroups         int    `json:"pure_groups"`
	MixedGroups        int    `json:"mixed_groups"`
	MaxClassesPerGroup int    `json:"max_classes_per_group"`
}

type UP249CResult struct {
	Schema                       string              `json:"schema"`
	Experiment                   string              `json:"experiment"`
	SourceUP248CSeal             string              `json:"source_up248c_seal"`
	ConditionA                   string              `json:"condition_a"`
	ConditionB                   string              `json:"condition_b"`
	PairedInitialStates          int                 `json:"paired_initial_states"`
	Coordinates                  []string            `json:"coordinates"`
	ParentCoordinateTriples      int                 `json:"parent_coordinate_triples"`
	ParentFullyPureTriples       int                 `json:"parent_fully_pure_triples"`
	CoordinateQuadruples         int                 `json:"coordinate_quadruples"`
	FullyPureQuadruples          int                 `json:"fully_pure_quadruples"`
	OutcomeClasses               []string            `json:"outcome_classes"`
	InterventionChanged          bool                `json:"intervention_changed"`
	ThresholdFittingUsed         bool                `json:"threshold_fitting_used"`
	ClassifierTrainingUsed       bool                `json:"classifier_training_used"`
	AdaptiveFeatureSelectionUsed bool                `json:"adaptive_feature_selection_used"`
	NewNativeFieldUsed           bool                `json:"new_native_field_used"`
	LiveActivation               bool                `json:"live_activation"`
	Summaries                    []UP249CQuadSummary `json:"summaries"`
}

func RunUP249C() (UP249CResult, error) {
	parent, err := RunUP248C()
	if err != nil {
		return UP249CResult{}, err
	}
	parentPure := 0
	for _, s := range parent.Summaries {
		if s.MixedGroups == 0 {
			parentPure++
		}
	}
	cohorts := [][]int{{0, 9, 10, 15}, {1, 4, 11, 14}, {2, 5, 8, 13}, {3, 6, 7, 12}}
	hands := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}
	coords := []string{"endangered_age", "predicted_is_endangered", "hand_distance", "age0", "age1", "age2", "age3", "adversarial_horizon", "no_query_horizon"}
	res := UP249CResult{
		Schema: UP249CQuadPuritySchema,
		Experiment: "UP-249C-quad-native-coordinate-purity",
		SourceUP248CSeal: "f55d8b01a26f4a90edf4e380225cb54ae3a29acb",
		ConditionA: "schedule2|advance22",
		ConditionB: "schedule3|advance0",
		Coordinates: coords,
		ParentCoordinateTriples: parent.CoordinateTriples,
		ParentFullyPureTriples: parentPure,
		OutcomeClasses: []string{"a_only_fail", "b_only_fail", "both_fail", "neither_fail"},
		InterventionChanged: false,
		ThresholdFittingUsed: false,
		ClassifierTrainingUsed: false,
		AdaptiveFeatureSelectionUsed: false,
		NewNativeFieldUsed: false,
		LiveActivation: false,
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
				return UP249CResult{}, err
			}
			b, err := up237cRun(x, e, 0, "fixed_offset_refresh", "alternating_shield")
			if err != nil {
				return UP249CResult{}, err
			}
			res.PairedInitialStates++
			samples = append(samples, sample{values: up246cValues(s), class: up246cClass(a, b)})
		}
	}
	for i := 0; i < len(coords); i++ {
		for j := i + 1; j < len(coords); j++ {
			for k := j + 1; k < len(coords); k++ {
				for l := k + 1; l < len(coords); l++ {
					a, b, c, d := coords[i], coords[j], coords[k], coords[l]
					groups := map[string]map[string]bool{}
					for _, s := range samples {
						key := fmt.Sprintf("%s=%s|%s=%s|%s=%s|%s=%s", a, s.values[a], b, s.values[b], c, s.values[c], d, s.values[d])
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
					sm := UP249CQuadSummary{CoordinateA: a, CoordinateB: b, CoordinateC: c, CoordinateD: d, Groups: len(groups)}
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
					if sm.MixedGroups == 0 {
						res.FullyPureQuadruples++
					}
					res.Summaries = append(res.Summaries, sm)
				}
			}
		}
	}
	res.CoordinateQuadruples = len(res.Summaries)
	return res, nil
}
