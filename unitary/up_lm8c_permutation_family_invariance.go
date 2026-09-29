package unitary

import "fmt"

const UPLM8CSchema = "wingless.up-lm8c-permutation-family-invariance.v1"

type UPLM8CRow struct {
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

type UPLM8CResult struct {
	Schema               string         `json:"schema"`
	Experiment           string         `json:"experiment"`
	SourceUPLM7CSeal     string         `json:"source_up_lm7c_seal"`
	ParentClassification string         `json:"parent_classification"`
	Rotations            []int          `json:"rotations"`
	Profiles             []string       `json:"profiles"`
	Permutations         []string       `json:"permutations"`
	RandomFixed          []int          `json:"random_fixed"`
	Rows                 []UPLM8CRow    `json:"rows"`
	ExactLawRows         int            `json:"exact_law_rows"`
	MismatchRows         int            `json:"mismatch_rows"`
	MinWriteError        int            `json:"min_write_error"`
	MaxWriteError        int            `json:"max_write_error"`
	MinPostTargetError   int            `json:"min_post_target_error"`
	MaxPostTargetError   int            `json:"max_post_target_error"`
	MismatchByPermutation map[string]int `json:"mismatch_by_permutation"`
	MismatchByRotation   map[string]int `json:"mismatch_by_rotation"`
	MismatchByProfile    map[string]int `json:"mismatch_by_profile"`
	LiveActivation       bool           `json:"live_activation"`
	PolicyChanged        bool           `json:"policy_changed"`
	AdaptiveSelectionUsed bool          `json:"adaptive_selection_used"`
	Classification       string         `json:"classification"`
}

func uplm8cPermute(in []uplm2xArm, permutation string) []uplm2xArm {
	n := len(in)
	out := make([]uplm2xArm, n)
	switch permutation {
	case "rotate1":
		for i := 0; i < n; i++ {
			out[i] = in[(i+1)%n]
		}
	case "random_fixed":
		order := []int{3, 0, 5, 1, 4, 2}
		for i := range order {
			out[i] = in[order[i]]
		}
	default:
		copy(out, in)
	}
	return out
}

func RunUPLM8C() (UPLM8CResult, error) {
	parent, err := RunUPLM7C()
	if err != nil {
		return UPLM8CResult{}, err
	}

	rotations := []int{8, 21}
	profiles := []string{"deferred_only", "layout_only", "hybrid_min"}
	permutations := []string{"identity", "rotate1", "random_fixed"}
	randomFixed := []int{3, 0, 5, 1, 4, 2}
	result := UPLM8CResult{
		Schema:                UPLM8CSchema,
		Experiment:            "UP-LM8C-permutation-family-invariance",
		SourceUPLM7CSeal:      "8537c0a562d87d48cfcde366623c77a1c6d6941a",
		ParentClassification:  parent.Classification,
		Rotations:             append([]int(nil), rotations...),
		Profiles:              append([]string(nil), profiles...),
		Permutations:          append([]string(nil), permutations...),
		RandomFixed:           append([]int(nil), randomFixed...),
		MismatchByPermutation: map[string]int{},
		MismatchByRotation:    map[string]int{},
		MismatchByProfile:     map[string]int{},
	}

	first := true
	allExact := true
	for _, profile := range profiles {
		for _, rotation := range rotations {
			for _, permutation := range permutations {
				arms := uplm8cPermute(uplm2xArms(rotation), permutation)
				for i := range arms {
					target := uplm3tTarget(arms[i], profile)
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
						arms[i].r.write(fmt.Sprintf("8c-pre-%s-%d-%s-%d-%d", profile, rotation, permutation, i, n), "x")
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

					result.Rows = append(result.Rows, UPLM8CRow{
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
						result.MismatchByPermutation[permutation]++
						result.MismatchByRotation[fmt.Sprint(rotation)]++
						result.MismatchByProfile[profile]++
						allExact = false
					}
				}
			}
		}
	}

	if parent.Classification != "HELDOUT_ROTATION_EXACT_DEFICIT_LAW" {
		result.Classification = "ANCHOR_NOT_REPRODUCED"
	} else if allExact {
		result.Classification = "PERMUTATION_FAMILY_EXACT_DEFICIT_LAW"
	} else {
		result.Classification = "PERMUTATION_DEPENDENT_DEFICIT_RESIDUAL"
	}
	return result, nil
}
