package unitary

import "fmt"

type upLm9dRow struct {
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

func RunUpLm9dPermutationInvariance(permutations []string) interface{} {
	profiles := []string{"deferred_only", "layout_only", "hybrid_max"}
	rotations := []int{8, 21}
	rows := make([]upLm9dRow, 0, len(permutations)*len(profiles)*len(rotations)*len(uplm2xArms(8)))

	mismatchRows := 0
	maxPostTargetError := 0
	maxWriteError := 0

	for _, profile := range profiles {
		for _, rotation := range rotations {
			for _, permutation := range permutations {
				arms := uplm8cPermute(uplm2xArms(rotation), permutation)
				for armIndex := range arms {
					target := uplm9cTarget(arms[armIndex], profile)
					preCountdown := 0
					if _, countdown, ok := uplm2xFirstPending(&arms[armIndex]); ok {
						preCountdown = countdown
					}

					writes := 0
					for writeIndex := 0; writeIndex < 32; writeIndex++ {
						_, countdown, ok := uplm2xFirstPending(&arms[armIndex])
						if !ok || countdown <= target {
							break
						}
						arms[armIndex].r.write(fmt.Sprintf("9d-pre-%s-%d-%s-%d-%d", profile, rotation, permutation, armIndex, writeIndex), "x")
						writes++
					}

					postCountdown := 0
					if _, countdown, ok := uplm2xFirstPending(&arms[armIndex]); ok {
						postCountdown = countdown
					}

					deficit := preCountdown - target
					if deficit < 0 {
						deficit = 0
					}
					writeError := writes - deficit
					postTargetError := postCountdown - target

					rows = append(rows, upLm9dRow{
						Profile:         profile,
						Rotation:        rotation,
						Permutation:     permutation,
						Arm:             armIndex,
						Target:          target,
						PreCountdown:    preCountdown,
						Writes:          writes,
						Deficit:         deficit,
						WriteError:      writeError,
						PostCountdown:   postCountdown,
						PostTargetError: postTargetError,
					})

					if writeError != 0 || postTargetError != 0 {
						mismatchRows++
					}
					if absolute := writeError; absolute < 0 {
						if -absolute > maxWriteError {
							maxWriteError = -absolute
						}
					} else if absolute > maxWriteError {
						maxWriteError = absolute
					}
					if absolute := postTargetError; absolute < 0 {
						if -absolute > maxPostTargetError {
							maxPostTargetError = -absolute
						}
					} else if absolute > maxPostTargetError {
						maxPostTargetError = absolute
					}
				}
			}
		}
	}

	return map[string]interface{}{
		"rows": rows,
		"metrics": map[string]float64{
			"mismatch_rows":         float64(mismatchRows),
			"max_post_target_error": float64(maxPostTargetError),
			"max_write_error":       float64(maxWriteError),
		},
	}
}
