package unitary

type upLm9dSubstrateRow struct {
	Seed       int    `json:"seed"`
	State      int    `json:"state"`
	Substrate  string `json:"substrate"`
	PostTarget int    `json:"post_target"`
	Write      int    `json:"write"`
}

type upLm9dSubstrateMetrics struct {
	MismatchRows                  float64 `json:"mismatch_rows"`
	MaxPostTargetError            float64 `json:"max_post_target_error"`
	MaxWriteError                 float64 `json:"max_write_error"`
	PairedComparisons             float64 `json:"paired_comparisons"`
	EmittedRows                   float64 `json:"emitted_rows"`
	UniqueStatesPerSeedPerSubstrate float64 `json:"unique_states_per_seed_per_substrate"`
	InvalidPermutations           float64 `json:"invalid_permutations"`
	PackedHighBitViolations       float64 `json:"packed_high_bit_violations"`
}

func RunUpLm9dSubstrateR1() interface{} {
	seeds := [...]int{104729, 130363, 155921, 181081, 206369, 232003, 257591, 283163}
	rows := make([]upLm9dSubstrateRow, 0, len(seeds)*512*2)

	mismatchRows := 0
	maxPostTargetError := 0
	maxWriteError := 0
	pairedComparisons := 0
	invalidPermutations := 0
	packedHighBitViolations := 0
	minimumUniqueStates := 512

	for _, seed := range seeds {
		permutation := upLm9dPermutation(seed)
		if !upLm9dValidPermutation(permutation) {
			invalidPermutations++
			continue
		}

		arrayStates := make(map[int]struct{}, 512)
		packedStates := make(map[int]struct{}, 512)
		for state := 0; state < 512; state++ {
			arrayPostTarget, arrayWrite := upLm9dArrayPath(state, permutation)
			packedPostTarget, packedWrite, highBitsSet := upLm9dPackedPath(state, permutation)
			if highBitsSet {
				packedHighBitViolations++
			}

			arrayStates[state] = struct{}{}
			packedStates[state] = struct{}{}
			rows = append(rows,
				upLm9dSubstrateRow{Seed: seed, State: state, Substrate: "array", PostTarget: arrayPostTarget, Write: arrayWrite},
				upLm9dSubstrateRow{Seed: seed, State: state, Substrate: "packed", PostTarget: packedPostTarget, Write: packedWrite},
			)
			pairedComparisons++

			postTargetError := upLm9dAbsolute(arrayPostTarget - packedPostTarget)
			writeError := upLm9dAbsolute(arrayWrite - packedWrite)
			if postTargetError != 0 || writeError != 0 {
				mismatchRows += 2
			}
			if postTargetError > maxPostTargetError {
				maxPostTargetError = postTargetError
			}
			if writeError > maxWriteError {
				maxWriteError = writeError
			}
		}

		if len(arrayStates) < minimumUniqueStates {
			minimumUniqueStates = len(arrayStates)
		}
		if len(packedStates) < minimumUniqueStates {
			minimumUniqueStates = len(packedStates)
		}
	}

	return map[string]interface{}{
		"rows": rows,
		"metrics": upLm9dSubstrateMetrics{
			MismatchRows:                    float64(mismatchRows),
			MaxPostTargetError:              float64(maxPostTargetError),
			MaxWriteError:                   float64(maxWriteError),
			PairedComparisons:               float64(pairedComparisons),
			EmittedRows:                     float64(len(rows)),
			UniqueStatesPerSeedPerSubstrate: float64(minimumUniqueStates),
			InvalidPermutations:             float64(invalidPermutations),
			PackedHighBitViolations:         float64(packedHighBitViolations),
		},
	}
}

func upLm9dPermutation(seed int) [9]int {
	permutation := [9]int{0, 1, 2, 3, 4, 5, 6, 7, 8}
	state := uint32(seed)
	for i := len(permutation) - 1; i > 0; i-- {
		state = state*1664525 + 1013904223
		j := int(state % uint32(i+1))
		permutation[i], permutation[j] = permutation[j], permutation[i]
	}
	return permutation
}

func upLm9dValidPermutation(permutation [9]int) bool {
	seen := [9]bool{}
	for _, index := range permutation {
		if index < 0 || index >= len(seen) || seen[index] {
			return false
		}
		seen[index] = true
	}
	return true
}

func upLm9dArrayPath(state int, permutation [9]int) (int, int) {
	var cells [9]uint8
	for index := 0; index < len(cells); index++ {
		cells[index] = uint8((state >> index) & 1)
	}

	var permuted [9]uint8
	for destination, source := range permutation {
		permuted[destination] = cells[source]
	}

	postTarget := 0
	write := 0
	for index, cell := range permuted {
		postTarget |= int(cell) << index
		if cell != 0 {
			write++
		}
	}
	return postTarget, write
}

func upLm9dPackedPath(state int, permutation [9]int) (int, int, bool) {
	packed := uint16(state)
	if packed&^uint16(0x01ff) != 0 {
		return 0, 0, true
	}

	var permuted uint16
	for destination, source := range permutation {
		bit := (packed >> uint(source)) & 1
		permuted |= bit << uint(destination)
	}
	if permuted&^uint16(0x01ff) != 0 {
		return 0, 0, true
	}

	postTarget := int(permuted)
	write := 0
	for index := 0; index < 9; index++ {
		write += int((permuted >> uint(index)) & 1)
	}
	return postTarget, write, false
}

func upLm9dAbsolute(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
