package unitary

import "fmt"

const UPLM9CSchema = "wingless.up-lm9c-profile-family-invariance.v1"

type UPLM9CRow struct {
	Profile         string `json:"profile"`
	Rotation        int    `json:"rotation"`
	Permutation     string `json:"permutation"`
	Arm             int    `json:"arm"`
	Target          int    `json:"target"`
	PreCountdown    int    `json:"pre_countdown"`
	Writes          int    `json:"writes"`
	Deficit         int    `json:"deficit"`
	WriteError      int    `json:"write_error"`
	PostCountdown   int    `json:"post_countdown"`
	PostTargetError int    `json:"post_target_error"`
}

type UPLM9CResult struct {
	Schema                string         `json:"schema"`
	Experiment            string         `json:"experiment"`
	SourceUPLM8CSeal      string         `json:"source_up_lm8c_seal"`
	ParentClassification  string         `json:"parent_classification"`
	Rotations             []int          `json:"rotations"`
	Profiles              []string       `json:"profiles"`
	Permutations          []string       `json:"permutations"`
	RandomFixed           []int          `json:"random_fixed"`
	Rows                  []UPLM9CRow    `json:"rows"`
	ExactLawRows          int            `json:"exact_law_rows"`
	MismatchRows          int            `json:"mismatch_rows"`
	MinWriteError         int            `json:"min_write_error"`
	MaxWriteError         int            `json:"max_write_error"`
	MinPostTargetError    int            `json:"min_post_target_error"`
	MaxPostTargetError    int            `json:"max_post_target_error"`
	MismatchByProfile     map[string]int `json:"mismatch_by_profile"`
	MismatchByPermutation map[string]int `json:"mismatch_by_permutation"`
	MismatchByRotation    map[string]int `json:"mismatch_by_rotation"`
	LiveActivation        bool           `json:"live_activation"`
	PolicyChanged         bool           `json:"policy_changed"`
	AdaptiveSelectionUsed bool           `json:"adaptive_selection_used"`
	Classification        string         `json:"classification"`
}

func uplm9cTarget(a uplm2xArm, profile string) int {
	switch profile {
	case "deferred_only":
		return uplm3hDeferredTarget(a)
	case "layout_only":
		return uplm3hLayoutTarget(a)
	case "hybrid_max":
		d := uplm3hDeferredTarget(a)
		l := uplm3hLayoutTarget(a)
		if d > l {
			return d
		}
		return l
	default:
		panic("UPLM9C_PROFILE_INVALID")
	}
}

func RunUPLM9C() (UPLM9CResult, error) {
	parent, err := RunUPLM8C()
	if err != nil {
		return UPLM9CResult{}, err
	}

	rotations := []int{8, 21}
	profiles := []string{"deferred_only", "layout_only", "hybrid_max"}
	permutations := []string{"identity", "rotate1", "random_fixed"}
	randomFixed := []int{3, 0, 5, 1, 4, 2}
	result := UPLM9CResult{
		Schema:                UPLM9CSchema,
		Experiment:            "UP-LM9C-profile-family-invariance",
		SourceUPLM8CSeal:      "f7874c5262f34906261b67d28c9d1d3f19c62665",
		ParentClassification:  parent.Classification,
		Rotations:             append([]int(nil), rotations...),
		Profiles:              append([]string(nil), profiles...),
		Permutations:          append([]string(nil), permutations...),
		RandomFixed:           append([]int(nil), randomFixed...),
		MismatchByProfile:     map[string]int{},
		MismatchByPermutation: map[string]int{},
		MismatchByRotation:    map[string]int{},
	}

	first := true
	allExact := true
	for _, profile := range profiles {
		for _, rotation := range rotations {
			for _, permutation := range permutations {
				arms := uplm8cPermute(uplm2xArms(rotation), permutation)
				for i := range arms {
					target := uplm9cTarget(arms[i], profile)
					preCountdown := 0
					if _, countdown, ok := uplm2xFirstPending(&arms[i]); ok {
						preCountdown = countdown
					}

					writes := 0
					for n := 0; n < 32; n++ {
						_, countdown, ok := uplm2xFirstPending(&arms[i])
						if !ok || countdown <= target {
							break
						}
						arms[i].r.write(fmt.Sprintf("9c-pre-%s-%d-%s-%d-%d", profile, rotation, permutation, i, n), "x")
						writes++
					}

					postCountdown := 0
					if _, countdown, ok := uplm2xFirstPending(&arms[i]); ok {
						postCountdown = countdown
					}
					deficit := preCountdown - target
					if deficit < 0 {
						deficit = 0
					}
					writeError := writes - deficit
					postTargetError := postCountdown - target

					result.Rows = append(result.Rows, UPLM9CRow{
						Profile: profile, Rotation: rotation, Permutation: permutation,
						Arm: i, Target: target, PreCountdown: preCountdown, Writes: writes,
						Deficit: deficit, WriteError: writeError,
						PostCountdown: postCountdown, PostTargetError: postTargetError,
					})

					if first {
						result.MinWriteError, result.MaxWriteError = writeError, writeError
						result.MinPostTargetError, result.MaxPostTargetError = postTargetError, postTargetError
						first = false
					} else {
						if writeError < result.MinWriteError {
							result.MinWriteError = writeError
						}
						if writeError > result.MaxWriteError {
							result.MaxWriteError = writeError
						}
						if postTargetError < result.MinPostTargetError {
							result.MinPostTargetError = postTargetError
						}
						if postTargetError > result.MaxPostTargetError {
							result.MaxPostTargetError = postTargetError
						}
					}

					if writeError == 0 && postTargetError == 0 {
						result.ExactLawRows++
					} else {
						result.MismatchRows++
						result.MismatchByProfile[profile]++
						result.MismatchByPermutation[permutation]++
						result.MismatchByRotation[fmt.Sprint(rotation)]++
						allExact = false
					}
				}
			}
		}
	}

	if parent.Classification != "PERMUTATION_FAMILY_EXACT_DEFICIT_LAW" {
		result.Classification = "ANCHOR_NOT_REPRODUCED"
	} else if allExact {
		result.Classification = "PROFILE_FAMILY_EXACT_DEFICIT_LAW"
	} else {
		result.Classification = "PROFILE_DEPENDENT_DEFICIT_RESIDUAL"
	}
	return result, nil
}
