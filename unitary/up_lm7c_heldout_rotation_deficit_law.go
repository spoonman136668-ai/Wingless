package unitary

import "fmt"

const UPLM7CSchema = "wingless.up-lm7c-heldout-rotation-deficit-law.v1"

type UPLM7CRow struct {
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

type UPLM7CResult struct {
	Schema               string         `json:"schema"`
	Experiment           string         `json:"experiment"`
	SourceUPLM6CSeal     string         `json:"source_up_lm6c_seal"`
	ParentClassification string         `json:"parent_classification"`
	HeldoutRotations     []int          `json:"heldout_rotations"`
	Rows                 []UPLM7CRow    `json:"rows"`
	ExactLawRows         int            `json:"exact_law_rows"`
	MismatchRows         int            `json:"mismatch_rows"`
	MinWriteError        int            `json:"min_write_error"`
	MaxWriteError        int            `json:"max_write_error"`
	MinPostTargetError   int            `json:"min_post_target_error"`
	MaxPostTargetError   int            `json:"max_post_target_error"`
	MismatchByRotation   map[string]int `json:"mismatch_by_rotation"`
	MismatchByProfile    map[string]int `json:"mismatch_by_profile"`
	LiveActivation       bool           `json:"live_activation"`
	PolicyChanged        bool           `json:"policy_changed"`
	Classification       string         `json:"classification"`
}

func RunUPLM7C() (UPLM7CResult, error) {
	parent, err := RunUPLM6C()
	if err != nil {
		return UPLM7CResult{}, err
	}

	profiles := []string{"deferred_only", "layout_only", "hybrid_min"}
	rotations := []int{8, 21}
	permutations := []string{"identity", "reverse", "rotate2"}
	result := UPLM7CResult{
		Schema:               UPLM7CSchema,
		Experiment:           "UP-LM7C-heldout-rotation-deficit-law",
		SourceUPLM6CSeal:     "0801e441f75298d08ab17f772f8595f34e649bc9",
		ParentClassification: parent.Classification,
		HeldoutRotations:     append([]int(nil), rotations...),
		MismatchByRotation:   map[string]int{},
		MismatchByProfile:    map[string]int{},
	}

	first := true
	allExact := true
	for _, profile := range profiles {
		for _, rotation := range rotations {
			for _, permutation := range permutations {
				arms := uplm2yPermute(uplm2xArms(rotation), permutation)
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
						arms[i].r.write(fmt.Sprintf("3t-pre-%s-%d-%d-%d", profile, rotation, i, n), "x")
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

					result.Rows = append(result.Rows, UPLM7CRow{
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
						result.MismatchByRotation[fmt.Sprint(rotation)]++
						result.MismatchByProfile[profile]++
						allExact = false
					}
				}
			}
		}
	}

	if parent.Classification != "EXACT_DEFICIT_LAW" {
		result.Classification = "ANCHOR_NOT_REPRODUCED"
	} else if allExact {
		result.Classification = "HELDOUT_ROTATION_EXACT_DEFICIT_LAW"
	} else {
		result.Classification = "ROTATION_DEPENDENT_DEFICIT_RESIDUAL"
	}
	return result, nil
}
