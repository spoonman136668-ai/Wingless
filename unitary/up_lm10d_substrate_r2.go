package unitary

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"math/bits"
)

const upLm10dMask uint16 = 0x01ff

type upLm10dRow struct {
	Seed               int    `json:"seed"`
	InitialState       int    `json:"initial_state"`
	Substrate          string `json:"substrate"`
	FinalState         int    `json:"final_state"`
	FinalPopcount      int    `json:"final_popcount"`
	TrajectoryDigest   string `json:"trajectory_digest"`
	MaxPostTargetError int    `json:"max_post_target_error"`
	MaxWriteError      int    `json:"max_write_error"`
}

type upLm10dMetrics struct {
	EmittedRows                     float64 `json:"emitted_rows"`
	PairedComparisons               float64 `json:"paired_comparisons"`
	UniqueStatesPerSeedPerSubstrate float64 `json:"unique_states_per_seed_per_substrate"`
	InvalidPermutations             float64 `json:"invalid_permutations"`
	MismatchRows                    float64 `json:"mismatch_rows"`
	MaxPostTargetError              float64 `json:"max_post_target_error"`
	MaxWriteError                   float64 `json:"max_write_error"`
	PackedHighBitViolations         float64 `json:"packed_high_bit_violations"`
}

type upLm10dOperation struct {
	kind int
	i    int
	j    int
}

type upLm10dOutcome struct {
	finalState int
	popcount   int
	digest     string
	postError  int
	writeError int
	highBits   int
}

// RunUpLm10dSubstrateR2 exhaustively compares array and packed nine-cell substrates.
func RunUpLm10dSubstrateR2() interface{} {
	seeds := [...]int{104729, 130363, 155921, 181081, 206369, 231701, 257053, 282407}
	rows := make([]upLm10dRow, 0, len(seeds)*512*2)
	invalidPermutations := 0
	mismatchRows := 0
	pairedComparisons := 0
	packedHighBitViolations := 0
	maxPostTargetError := 0
	maxWriteError := 0
	minimumUniqueStates := 512

	for _, seed := range seeds {
		order := upLm10dOrder(seed)
		if !upLm10dValidOrder(order) {
			invalidPermutations++
			continue
		}
		schedule := upLm10dSchedule(seed)
		arrayStates := make(map[int]struct{}, 512)
		packedStates := make(map[int]struct{}, 512)

		for _, initial := range order {
			arrayOutcome := upLm10dRunArray(initial, schedule)
			packedOutcome := upLm10dRunPacked(initial, schedule)
			arrayStates[initial] = struct{}{}
			packedStates[initial] = struct{}{}
			packedHighBitViolations += packedOutcome.highBits
			if arrayOutcome.postError > maxPostTargetError {
				maxPostTargetError = arrayOutcome.postError
			}
			if packedOutcome.postError > maxPostTargetError {
				maxPostTargetError = packedOutcome.postError
			}
			if arrayOutcome.writeError > maxWriteError {
				maxWriteError = arrayOutcome.writeError
			}
			if packedOutcome.writeError > maxWriteError {
				maxWriteError = packedOutcome.writeError
			}
			if arrayOutcome.finalState != packedOutcome.finalState || arrayOutcome.popcount != packedOutcome.popcount || arrayOutcome.digest != packedOutcome.digest {
				mismatchRows++
			}
			rows = append(rows,
				upLm10dRow{Seed: seed, InitialState: initial, Substrate: "array", FinalState: arrayOutcome.finalState, FinalPopcount: arrayOutcome.popcount, TrajectoryDigest: arrayOutcome.digest, MaxPostTargetError: arrayOutcome.postError, MaxWriteError: arrayOutcome.writeError},
				upLm10dRow{Seed: seed, InitialState: initial, Substrate: "packed", FinalState: packedOutcome.finalState, FinalPopcount: packedOutcome.popcount, TrajectoryDigest: packedOutcome.digest, MaxPostTargetError: packedOutcome.postError, MaxWriteError: packedOutcome.writeError},
			)
			pairedComparisons++
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
		"metrics": upLm10dMetrics{
			EmittedRows:                     float64(len(rows)),
			PairedComparisons:               float64(pairedComparisons),
			UniqueStatesPerSeedPerSubstrate: float64(minimumUniqueStates),
			InvalidPermutations:             float64(invalidPermutations),
			MismatchRows:                    float64(mismatchRows),
			MaxPostTargetError:              float64(maxPostTargetError),
			MaxWriteError:                   float64(maxWriteError),
			PackedHighBitViolations:         float64(packedHighBitViolations),
		},
	}
}

func upLm10dOrder(seed int) [512]int {
	var order [512]int
	for i := range order {
		order[i] = i
	}
	stream := uint64(seed) ^ 0xD1B54A32D192ED03
	for i := len(order) - 1; i > 0; i-- {
		stream = upLm10dSplitMix64(stream)
		j := int(stream % uint64(i+1))
		order[i], order[j] = order[j], order[i]
	}
	return order
}

func upLm10dValidOrder(order [512]int) bool {
	var seen [512]bool
	for _, state := range order {
		if state < 0 || state >= len(seen) || seen[state] {
			return false
		}
		seen[state] = true
	}
	return true
}

func upLm10dSchedule(seed int) [64]upLm10dOperation {
	var schedule [64]upLm10dOperation
	stream := uint64(seed) ^ 0x9E3779B97F4A7C15
	for n := range schedule {
		stream = upLm10dSplitMix64(stream)
		i := int(stream % 9)
		stream = upLm10dSplitMix64(stream)
		j := int(stream % 9)
		stream = upLm10dSplitMix64(stream)
		schedule[n] = upLm10dOperation{kind: int(stream % 3), i: i, j: j}
	}
	return schedule
}

func upLm10dSplitMix64(state uint64) uint64 {
	state += 0x9E3779B97F4A7C15
	z := state
	z = (z ^ (z >> 30)) * 0xBF58476D1CE4E5B9
	z = (z ^ (z >> 27)) * 0x94D049BB133111EB
	return z ^ (z >> 31)
}

func upLm10dRunArray(initial int, schedule [64]upLm10dOperation) upLm10dOutcome {
	var cells [9]bool
	for i := range cells {
		cells[i] = initial&(1<<uint(i)) != 0
	}
	trace := make([]byte, 0, 130)
	trace = upLm10dAppendState(trace, initial)
	maximumPostError, maximumWriteError := 0, 0
	for _, operation := range schedule {
		before := upLm10dArrayCanonical(cells)
		reference := upLm10dReference(before, operation)
		prior := cells
		switch operation.kind {
		case 0:
			cells[operation.i] = !cells[operation.i]
		case 1:
			cells[operation.j] = cells[operation.i]
		case 2:
			cells[operation.i], cells[operation.j] = cells[operation.j], cells[operation.i]
		}
		observed := upLm10dArrayCanonical(cells)
		maximumPostError = upLm10dMax(maximumPostError, bits.OnesCount16(uint16(observed^reference)))
		maximumWriteError = upLm10dMax(maximumWriteError, upLm10dArrayWriteError(prior, cells, operation))
		trace = upLm10dAppendState(trace, observed)
	}
	finalState := upLm10dArrayCanonical(cells)
	digest := sha256.Sum256(trace)
	return upLm10dOutcome{finalState: finalState, popcount: bits.OnesCount16(uint16(finalState)), digest: hex.EncodeToString(digest[:]), postError: maximumPostError, writeError: maximumWriteError}
}

func upLm10dRunPacked(initial int, schedule [64]upLm10dOperation) upLm10dOutcome {
	packed := uint16(initial) & upLm10dMask
	trace := make([]byte, 0, 130)
	trace = upLm10dAppendState(trace, int(packed))
	maximumPostError, maximumWriteError, highBits := 0, 0, 0
	for _, operation := range schedule {
		before := int(packed & upLm10dMask)
		reference := upLm10dReference(before, operation)
		prior := packed
		switch operation.kind {
		case 0:
			packed ^= uint16(1 << uint(operation.i))
		case 1:
			source := (packed >> uint(operation.i)) & 1
			mask := uint16(1 << uint(operation.j))
			packed = (packed &^ mask) | (source << uint(operation.j))
		case 2:
			left := (packed >> uint(operation.i)) & 1
			right := (packed >> uint(operation.j)) & 1
			if left != right {
				packed ^= uint16(1<<uint(operation.i)) | uint16(1<<uint(operation.j))
			}
		}
		if packed&^upLm10dMask != 0 {
			highBits++
		}
		observed := int(packed & upLm10dMask)
		maximumPostError = upLm10dMax(maximumPostError, bits.OnesCount16(uint16(observed^reference)))
		maximumWriteError = upLm10dMax(maximumWriteError, upLm10dPackedWriteError(prior, packed, operation))
		trace = upLm10dAppendState(trace, observed)
	}
	finalState := int(packed & upLm10dMask)
	digest := sha256.Sum256(trace)
	return upLm10dOutcome{finalState: finalState, popcount: bits.OnesCount16(uint16(finalState)), digest: hex.EncodeToString(digest[:]), postError: maximumPostError, writeError: maximumWriteError, highBits: highBits}
}

func upLm10dReference(state int, operation upLm10dOperation) int {
	state &= int(upLm10dMask)
	switch operation.kind {
	case 0:
		return state ^ (1 << uint(operation.i))
	case 1:
		mask := 1 << uint(operation.j)
		return (state &^ mask) | (((state >> uint(operation.i)) & 1) << uint(operation.j))
	default:
		left := (state >> uint(operation.i)) & 1
		right := (state >> uint(operation.j)) & 1
		if left != right {
			state ^= (1 << uint(operation.i)) | (1 << uint(operation.j))
		}
		return state
	}
}

func upLm10dArrayCanonical(cells [9]bool) int {
	state := 0
	for i, cell := range cells {
		if cell {
			state |= 1 << uint(i)
		}
	}
	return state
}

func upLm10dArrayWriteError(before, after [9]bool, operation upLm10dOperation) int {
	errors := 0
	for i := range before {
		allowed := i == operation.i
		if operation.kind == 1 {
			allowed = i == operation.j
		} else if operation.kind == 2 {
			allowed = i == operation.i || i == operation.j
		}
		if !allowed && before[i] != after[i] {
			errors++
		}
	}
	return errors
}

func upLm10dPackedWriteError(before, after uint16, operation upLm10dOperation) int {
	errors := 0
	for i := 0; i < 9; i++ {
		allowed := i == operation.i
		if operation.kind == 1 {
			allowed = i == operation.j
		} else if operation.kind == 2 {
			allowed = i == operation.i || i == operation.j
		}
		if !allowed && ((before>>uint(i))&1) != ((after>>uint(i))&1) {
			errors++
		}
	}
	return errors
}

func upLm10dAppendState(trace []byte, state int) []byte {
	var encoded [2]byte
	binary.LittleEndian.PutUint16(encoded[:], uint16(state)&upLm10dMask)
	return append(trace, encoded[:]...)
}

func upLm10dMax(left, right int) int {
	if right > left {
		return right
	}
	return left
}
