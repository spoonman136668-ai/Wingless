package unitary

import (
	"crypto/sha256"
	"encoding/binary"
)

type upLm15DemandMultiplicityR1Row struct {
	Seed                int    `json:"seed"`
	ScheduleIndex       int    `json:"schedule_index"`
	DemandFactor        int    `json:"demand_factor"`
	Substrate           string `json:"substrate"`
	Step                int    `json:"step"`
	StateIndex          int    `json:"state_index"`
	DemandWrites        int    `json:"demand_writes"`
	AcceptedWrites      int    `json:"accepted_writes"`
	Deficit             int    `json:"deficit"`
	AcceptedWritesSoFar int    `json:"accepted_writes_so_far"`
	DeficitSoFar        int    `json:"deficit_so_far"`
	RemainingCap        int    `json:"remaining_cap"`
	StateCardinality    int    `json:"state_cardinality"`
	StateDigest         string `json:"state_digest"`
}

// RunUpLm15DemandMultiplicityR1 replays the frozen schedule matrix with each
// preregistered per-row demand multiplicity on paired dense and sparse states.
func RunUpLm15DemandMultiplicityR1() interface{} {
	const (
		stateCardinality = 256
		rowsPerCase      = 64
		writeCap         = 36
	)
	seeds := [...]int{11, 29, 47, 71, 101, 131, 173, 211}
	factors := [...]int{0, 1, 2, 4}

	rows := make([]upLm15DemandMultiplicityR1Row, 0, len(seeds)*8*len(factors)*2*rowsPerCase)
	lawViolations := 0
	maxConservationError := 0
	pairedObservableMismatches := 0
	budgetOverrunRows := 0
	invalidSchedules := 0
	rowCountMismatches := 0
	stateCardinalityMismatches := 0
	completedRuns := 0
	pairedCases := 0
	aggregateDemand := 0

	for _, seed := range seeds {
		for scheduleIndex := 0; scheduleIndex < 8; scheduleIndex++ {
			trace := upLm15DemandMultiplicityR1Trace(seed, scheduleIndex, rowsPerCase)
			if !upLm15DemandMultiplicityR1ValidTrace(trace, stateCardinality) {
				invalidSchedules++
			}

			for _, factor := range factors {
				pairedCases++
				dense := make([]int, stateCardinality)
				sparse := make(map[int]int, stateCardinality)
				denseAccepted, denseDeficit, denseRowCount := 0, 0, 0
				sparseAccepted, sparseDeficit, sparseRowCount := 0, 0, 0
				denseLawViolation, sparseLawViolation := false, false
				denseCardinalityMismatch, sparseCardinalityMismatch := false, false
				pairedMismatch := false

				for step, stateIndex := range trace {
					denseAcceptedNow, denseDeficitNow := 0, 0
					sparseAcceptedNow, sparseDeficitNow := 0, 0
					for attempt := 0; attempt < factor; attempt++ {
						if denseAccepted < writeCap {
							dense[stateIndex]++
							denseAccepted++
							denseAcceptedNow++
						} else {
							denseDeficit++
							denseDeficitNow++
						}
						if sparseAccepted < writeCap {
							sparse[stateIndex]++
							sparseAccepted++
							sparseAcceptedNow++
						} else {
							sparseDeficit++
							sparseDeficitNow++
						}
					}
					denseRowCount++
					sparseRowCount++

					denseError := upLm15DemandMultiplicityR1Abs(factor - denseAcceptedNow - denseDeficitNow)
					sparseError := upLm15DemandMultiplicityR1Abs(factor - sparseAcceptedNow - sparseDeficitNow)
					if denseError > maxConservationError {
						maxConservationError = denseError
					}
					if sparseError > maxConservationError {
						maxConservationError = sparseError
					}
					if denseError != 0 || denseAccepted > writeCap {
						denseLawViolation = true
					}
					if sparseError != 0 || sparseAccepted > writeCap {
						sparseLawViolation = true
					}
					if denseAccepted > writeCap {
						budgetOverrunRows++
					}
					if sparseAccepted > writeCap {
						budgetOverrunRows++
					}

					denseVector := upLm15DemandMultiplicityR1DenseVector(dense)
					sparseVector := upLm15DemandMultiplicityR1SparseVector(sparse, stateCardinality)
					if len(denseVector) != stateCardinality {
						denseCardinalityMismatch = true
					}
					if len(sparseVector) != stateCardinality {
						sparseCardinalityMismatch = true
					}
					denseDigest := upLm15DemandMultiplicityR1Digest(denseVector)
					sparseDigest := upLm15DemandMultiplicityR1Digest(sparseVector)
					if !upLm15DemandMultiplicityR1SameVector(denseVector, sparseVector) ||
						denseDigest != sparseDigest ||
						denseAccepted != sparseAccepted || denseDeficit != sparseDeficit {
						pairedMismatch = true
					}

					rows = append(rows,
						upLm15DemandMultiplicityR1Row{seed, scheduleIndex, factor, "dense-slice", step, stateIndex, factor, denseAcceptedNow, denseDeficitNow, denseAccepted, denseDeficit, writeCap - denseAccepted, len(denseVector), denseDigest},
						upLm15DemandMultiplicityR1Row{seed, scheduleIndex, factor, "sparse-map", step, stateIndex, factor, sparseAcceptedNow, sparseDeficitNow, sparseAccepted, sparseDeficit, writeCap - sparseAccepted, len(sparseVector), sparseDigest},
					)
					aggregateDemand += factor * 2
				}

				completedRuns += 2
				if denseRowCount != rowsPerCase || sparseRowCount != rowsPerCase {
					rowCountMismatches++
				}
				if denseCardinalityMismatch {
					stateCardinalityMismatches++
				}
				if sparseCardinalityMismatch {
					stateCardinalityMismatches++
				}
				if pairedMismatch {
					pairedObservableMismatches++
				}
				if denseLawViolation {
					lawViolations++
				}
				if sparseLawViolation {
					lawViolations++
				}
			}
		}
	}

	return map[string]interface{}{
		"rows": rows,
		"metrics": map[string]float64{
			"law_violation_cases":              float64(lawViolations),
			"max_conservation_error":           float64(maxConservationError),
			"paired_observable_mismatch_cases": float64(pairedObservableMismatches),
			"budget_overrun_rows":               float64(budgetOverrunRows),
			"invalid_state_schedules":           float64(invalidSchedules),
			"row_count_mismatch_cases":          float64(rowCountMismatches),
			"state_cardinality_mismatch_cases":  float64(stateCardinalityMismatches),
			"completed_substrate_case_runs":     float64(completedRuns),
			"paired_case_count":                 float64(pairedCases),
			"total_emitted_rows":                float64(len(rows)),
			"aggregate_demand_writes":           float64(aggregateDemand),
		},
	}
}

func upLm15DemandMultiplicityR1Trace(seed, schedule, rows int) []int {
	trace := make([]int, rows)
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

func upLm15DemandMultiplicityR1ValidTrace(trace []int, cardinality int) bool {
	for _, index := range trace {
		if index < 0 || index >= cardinality {
			return false
		}
	}
	return true
}

func upLm15DemandMultiplicityR1DenseVector(state []int) []int {
	vector := make([]int, len(state))
	copy(vector, state)
	return vector
}

func upLm15DemandMultiplicityR1SparseVector(state map[int]int, cardinality int) []int {
	vector := make([]int, cardinality)
	for index := range vector {
		vector[index] = state[index]
	}
	return vector
}

func upLm15DemandMultiplicityR1SameVector(left, right []int) bool {
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

func upLm15DemandMultiplicityR1Digest(vector []int) string {
	encoded := make([]byte, len(vector)*8)
	for index, value := range vector {
		binary.BigEndian.PutUint64(encoded[index*8:(index+1)*8], uint64(value))
	}
	digest := sha256.Sum256(encoded)
	return string(digest[:])
}

func upLm15DemandMultiplicityR1Abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
