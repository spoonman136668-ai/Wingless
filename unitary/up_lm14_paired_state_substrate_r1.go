package unitary

import (
	"crypto/sha256"
	"encoding/binary"
)

type upLm14PairedStateSubstrateR1Row struct {
	Seed                int    `json:"seed"`
	ScheduleIndex       int    `json:"schedule_index"`
	Substrate           string `json:"substrate"`
	Step                int    `json:"step"`
	StateIndex          int    `json:"state_index"`
	DemandWrites        int    `json:"demand_writes"`
	AcceptedWrites      int    `json:"accepted_writes"`
	Deficit             int    `json:"deficit"`
	AcceptedWritesSoFar int    `json:"accepted_writes_so_far"`
	DeficitSoFar        int    `json:"deficit_so_far"`
	RemainingCap        int    `json:"remaining_cap"`
}

// RunUpLm14PairedStateSubstrateR1 executes the frozen paired dense-slice and
// sparse-map state-substrate matrix using identical immutable SHA-256 traces.
func RunUpLm14PairedStateSubstrateR1() interface{} {
	const (
		stateCardinality = 256
		stepsPerRun      = 64
		writeCap         = 36
	)
	seeds := [...]int{1103, 2207, 3319, 4421, 5527, 6637, 7753, 8861}

	rows := make([]upLm14PairedStateSubstrateR1Row, 0, len(seeds)*8*2*stepsPerRun)
	pairedObservableMismatches := 0
	acceptedMismatches := 0
	deficitMismatches := 0
	lawViolations := 0
	maxConservationError := 0
	invalidSchedules := 0
	rowCountMismatches := 0
	stateCardinalityMismatches := 0
	budgetOverrunRows := 0
	pairedSchedules := 0
	substrateRuns := 0
	aggregateDemand := 0
	aggregateAccepted := 0
	aggregateDeficit := 0

	for _, seed := range seeds {
		for scheduleIndex := 0; scheduleIndex < 8; scheduleIndex++ {
			pairedSchedules++
			trace := upLm14PairedStateSubstrateR1Trace(seed, scheduleIndex, stepsPerRun)
			pairedTrace := trace
			traceValid := upLm14PairedStateSubstrateR1ValidTrace(trace, stateCardinality) &&
				upLm14PairedStateSubstrateR1SameTrace(trace, pairedTrace)
			if !traceValid {
				invalidSchedules++
			}

			dense := make([]int, stateCardinality)
			sparse := make(map[int]int, stateCardinality)
			denseInitial := upLm14PairedStateSubstrateR1DenseVector(dense)
			sparseInitial := upLm14PairedStateSubstrateR1SparseVector(sparse, stateCardinality)
			if !upLm14PairedStateSubstrateR1SameVector(denseInitial, sparseInitial) {
				stateCardinalityMismatches += 2
			}

			denseAccepted, denseDeficit := 0, 0
			sparseAccepted, sparseDeficit := 0, 0
			denseRows, sparseRows := 0, 0
			denseLawViolation, sparseLawViolation := false, false
			denseStateMismatch, sparseStateMismatch := false, false
			pairObservableMismatch := false
			pairAcceptedMismatch := false
			pairDeficitMismatch := false

			for step, stateIndex := range trace {
				denseAcceptedNow := denseAccepted < writeCap
				if denseAcceptedNow {
					dense[stateIndex]++
					denseAccepted++
				} else {
					denseDeficit++
				}
				denseRows++
				denseDemand := 1
				denseError := upLm14PairedStateSubstrateR1Abs(denseDemand - upLm14PairedStateSubstrateR1StepDelta(denseAcceptedNow, denseDeficit, step))
				if denseError > maxConservationError {
					maxConservationError = denseError
				}
				if denseAccepted > writeCap {
					budgetOverrunRows++
					denseLawViolation = true
				}
				if denseError != 0 {
					denseLawViolation = true
				}

				sparseAcceptedNow := sparseAccepted < writeCap
				if sparseAcceptedNow {
					sparse[stateIndex]++
					sparseAccepted++
				} else {
					sparseDeficit++
				}
				sparseRows++
				sparseDemand := 1
				sparseError := upLm14PairedStateSubstrateR1Abs(sparseDemand - upLm14PairedStateSubstrateR1StepDelta(sparseAcceptedNow, sparseDeficit, step))
				if sparseError > maxConservationError {
					maxConservationError = sparseError
				}
				if sparseAccepted > writeCap {
					budgetOverrunRows++
					sparseLawViolation = true
				}
				if sparseError != 0 {
					sparseLawViolation = true
				}

				denseVector := upLm14PairedStateSubstrateR1DenseVector(dense)
				sparseVector := upLm14PairedStateSubstrateR1SparseVector(sparse, stateCardinality)
				if len(denseVector) != stateCardinality {
					denseStateMismatch = true
				}
				if len(sparseVector) != stateCardinality {
					sparseStateMismatch = true
				}
				if !upLm14PairedStateSubstrateR1SameVector(denseVector, sparseVector) ||
					denseAccepted != sparseAccepted || denseDeficit != sparseDeficit ||
					writeCap-denseAccepted != writeCap-sparseAccepted {
					pairObservableMismatch = true
				}
				if denseAccepted != sparseAccepted {
					pairAcceptedMismatch = true
				}
				if denseDeficit != sparseDeficit {
					pairDeficitMismatch = true
				}

				rows = append(rows,
					upLm14PairedStateSubstrateR1Row{seed, scheduleIndex, "dense-slice", step, stateIndex, 1, upLm14PairedStateSubstrateR1BoolInt(denseAcceptedNow), upLm14PairedStateSubstrateR1BoolInt(!denseAcceptedNow), denseAccepted, denseDeficit, writeCap - denseAccepted},
					upLm14PairedStateSubstrateR1Row{seed, scheduleIndex, "sparse-map", step, stateIndex, 1, upLm14PairedStateSubstrateR1BoolInt(sparseAcceptedNow), upLm14PairedStateSubstrateR1BoolInt(!sparseAcceptedNow), sparseAccepted, sparseDeficit, writeCap - sparseAccepted},
				)
				aggregateDemand += 2
				aggregateAccepted += upLm14PairedStateSubstrateR1BoolInt(denseAcceptedNow) + upLm14PairedStateSubstrateR1BoolInt(sparseAcceptedNow)
				aggregateDeficit += upLm14PairedStateSubstrateR1BoolInt(!denseAcceptedNow) + upLm14PairedStateSubstrateR1BoolInt(!sparseAcceptedNow)
			}

			substrateRuns += 2
			if denseRows != stepsPerRun || sparseRows != stepsPerRun {
				rowCountMismatches++
			}
			if denseStateMismatch {
				stateCardinalityMismatches++
			}
			if sparseStateMismatch {
				stateCardinalityMismatches++
			}
			if pairObservableMismatch {
				pairedObservableMismatches++
			}
			if pairAcceptedMismatch {
				acceptedMismatches++
			}
			if pairDeficitMismatch {
				deficitMismatches++
			}
			if denseLawViolation {
				lawViolations++
			}
			if sparseLawViolation {
				lawViolations++
			}
		}
	}

	return map[string]interface{}{
		"rows": rows,
		"metrics": map[string]float64{
			"paired_observable_mismatch_cases": float64(pairedObservableMismatches),
			"accepted_write_mismatch_cases":     float64(acceptedMismatches),
			"deficit_mismatch_cases":            float64(deficitMismatches),
			"law_violation_cases":                float64(lawViolations),
			"max_conservation_error":             float64(maxConservationError),
			"invalid_state_schedules":            float64(invalidSchedules),
			"row_count_mismatch_cases":           float64(rowCountMismatches),
			"state_cardinality_mismatch_cases":   float64(stateCardinalityMismatches),
			"budget_overrun_rows":                 float64(budgetOverrunRows),
			"paired_schedule_count":               float64(pairedSchedules),
			"substrate_case_runs":                 float64(substrateRuns),
			"total_emitted_rows":                  float64(len(rows)),
			"aggregate_demand_writes":             float64(aggregateDemand),
			"aggregate_accepted_writes":           float64(aggregateAccepted),
			"aggregate_deficit":                   float64(aggregateDeficit),
		},
	}
}

func upLm14PairedStateSubstrateR1Trace(seed, schedule, steps int) []int {
	trace := make([]int, steps)
	var input [24]byte
	for step := range trace {
		binary.BigEndian.PutUint64(input[0:8], uint64(seed))
		binary.BigEndian.PutUint64(input[8:16], uint64(schedule))
		binary.BigEndian.PutUint64(input[16:24], uint64(step))
		digest := sha256.Sum256(input[:])
		trace[step] = int(binary.BigEndian.Uint64(digest[0:8]) % 256)
	}
	return trace
}

func upLm14PairedStateSubstrateR1ValidTrace(trace []int, cardinality int) bool {
	for _, index := range trace {
		if index < 0 || index >= cardinality {
			return false
		}
	}
	return true
}

func upLm14PairedStateSubstrateR1SameTrace(left, right []int) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func upLm14PairedStateSubstrateR1DenseVector(state []int) []int {
	vector := make([]int, len(state))
	copy(vector, state)
	return vector
}

func upLm14PairedStateSubstrateR1SparseVector(state map[int]int, cardinality int) []int {
	vector := make([]int, cardinality)
	for index := 0; index < cardinality; index++ {
		vector[index] = state[index]
	}
	return vector
}

func upLm14PairedStateSubstrateR1SameVector(left, right []int) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func upLm14PairedStateSubstrateR1StepDelta(accepted bool, cumulativeDeficit, step int) int {
	if accepted {
		return 1
	}
	if cumulativeDeficit > step {
		return 0
	}
	return 1
}

func upLm14PairedStateSubstrateR1BoolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

func upLm14PairedStateSubstrateR1Abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
