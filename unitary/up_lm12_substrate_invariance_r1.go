package unitary

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
)

type upLm12SubstrateInvarianceR1Row struct {
	Seed                int    `json:"seed"`
	Budget              int    `json:"budget"`
	Substrate           string `json:"substrate"`
	DemandOrdinal       int    `json:"demand_ordinal"`
	StateIndex          int    `json:"state_index"`
	Attempted           bool   `json:"attempted"`
	Accepted            bool   `json:"accepted"`
	AcceptedWritesSoFar int    `json:"accepted_writes_so_far"`
	DeficitWritesSoFar  int    `json:"deficit_writes_so_far"`
}

type upLm12SubstrateInvarianceR1Metrics struct {
	LogicalComparisonCases       float64 `json:"logical_comparison_cases"`
	SubstrateCaseRuns            float64 `json:"substrate_case_runs"`
	TotalEmittedRows             float64 `json:"total_emitted_rows"`
	AggregateDemandWrites        float64 `json:"aggregate_demand_writes"`
	AggregateAcceptedWrites      float64 `json:"aggregate_accepted_writes"`
	AggregateDeficit             float64 `json:"aggregate_deficit"`
	PrefixDigestComparisons      float64 `json:"prefix_digest_comparisons"`
	PrefixDigestMismatchCases    float64 `json:"prefix_digest_mismatch_cases"`
	AcceptedWriteMismatchCases   float64 `json:"accepted_write_mismatch_cases"`
	DeficitMismatchCases         float64 `json:"deficit_mismatch_cases"`
	LawViolationCases            float64 `json:"law_violation_cases"`
	MaxConservationError         float64 `json:"max_conservation_error"`
	BudgetOverrunRows            float64 `json:"budget_overrun_rows"`
	InvalidStateSchedules        float64 `json:"invalid_state_schedules"`
	MinUniqueStatesPerSeed       float64 `json:"min_unique_states_per_seed"`
	MaxUniqueStatesPerSeed       float64 `json:"max_unique_states_per_seed"`
}

type upLm12SubstrateInvarianceR1Dense struct {
	cells [512]bool
}

type upLm12SubstrateInvarianceR1Sparse struct {
	cells map[int]bool
}

func (s *upLm12SubstrateInvarianceR1Dense) set(index int) { s.cells[index] = true }
func (s *upLm12SubstrateInvarianceR1Dense) get(index int) bool { return s.cells[index] }
func (s *upLm12SubstrateInvarianceR1Sparse) set(index int) { s.cells[index] = true }
func (s *upLm12SubstrateInvarianceR1Sparse) get(index int) bool { return s.cells[index] }

// RunUpLm12SubstrateInvarianceR1 compares dense and sparse state substrates under identical schedules.
func RunUpLm12SubstrateInvarianceR1() interface{} {
	seeds := [...]int{101, 211, 307, 401, 503, 601, 701, 809}
	budgets := [...]int{0, 64, 128, 256, 512}
	rows := make([]upLm12SubstrateInvarianceR1Row, 0, len(seeds)*len(budgets)*2*512)

	logicalCases := 0
	substrateRuns := 0
	aggregateDemand := 0
	aggregateAccepted := 0
	aggregateDeficit := 0
	prefixComparisons := 0
	prefixMismatches := 0
	acceptedMismatches := 0
	deficitMismatches := 0
	lawViolations := 0
	maxConservationError := 0
	budgetOverruns := 0
	invalidSchedules := 0
	minUnique := 512
	maxUnique := 0

	for _, seed := range seeds {
		schedule := upLm12SubstrateInvarianceR1Schedule(seed)
		unique, valid := upLm12SubstrateInvarianceR1ValidateSchedule(schedule)
		if unique < minUnique {
			minUnique = unique
		}
		if unique > maxUnique {
			maxUnique = unique
		}
		if !valid {
			invalidSchedules++
			continue
		}

		for _, budget := range budgets {
			logicalCases++
			denseAccepted, denseDeficit, denseOverruns, denseDigest, denseRows := upLm12SubstrateInvarianceR1ReplayDense(seed, budget, schedule)
			sparseAccepted, sparseDeficit, sparseOverruns, sparseDigest, sparseRows := upLm12SubstrateInvarianceR1ReplaySparse(seed, budget, schedule)
			rows = append(rows, denseRows...)
			rows = append(rows, sparseRows...)
			substrateRuns += 2
			aggregateDemand += len(schedule) * 2
			aggregateAccepted += denseAccepted + sparseAccepted
			aggregateDeficit += denseDeficit + sparseDeficit
			budgetOverruns += denseOverruns + sparseOverruns

			prefixComparisons++
			if denseDigest != sparseDigest {
				prefixMismatches++
			}
			if denseAccepted != sparseAccepted {
				acceptedMismatches++
			}
			if denseDeficit != sparseDeficit {
				deficitMismatches++
			}

			denseError := upLm12SubstrateInvarianceR1Abs(len(schedule) - (denseAccepted + denseDeficit))
			sparseError := upLm12SubstrateInvarianceR1Abs(len(schedule) - (sparseAccepted + sparseDeficit))
			if denseError > maxConservationError {
				maxConservationError = denseError
			}
			if sparseError > maxConservationError {
				maxConservationError = sparseError
			}
			if denseError != 0 || sparseError != 0 || denseOverruns != 0 || sparseOverruns != 0 {
				lawViolations++
			}
		}
	}

	return map[string]interface{}{
		"rows": rows,
		"metrics": upLm12SubstrateInvarianceR1Metrics{
			LogicalComparisonCases:     float64(logicalCases),
			SubstrateCaseRuns:          float64(substrateRuns),
			TotalEmittedRows:           float64(len(rows)),
			AggregateDemandWrites:      float64(aggregateDemand),
			AggregateAcceptedWrites:    float64(aggregateAccepted),
			AggregateDeficit:           float64(aggregateDeficit),
			PrefixDigestComparisons:    float64(prefixComparisons),
			PrefixDigestMismatchCases:  float64(prefixMismatches),
			AcceptedWriteMismatchCases: float64(acceptedMismatches),
			DeficitMismatchCases:       float64(deficitMismatches),
			LawViolationCases:          float64(lawViolations),
			MaxConservationError:       float64(maxConservationError),
			BudgetOverrunRows:          float64(budgetOverruns),
			InvalidStateSchedules:      float64(invalidSchedules),
			MinUniqueStatesPerSeed:     float64(minUnique),
			MaxUniqueStatesPerSeed:     float64(maxUnique),
		},
	}
}

func upLm12SubstrateInvarianceR1ReplayDense(seed, budget int, schedule [512]int) (int, int, int, string, []upLm12SubstrateInvarianceR1Row) {
	state := &upLm12SubstrateInvarianceR1Dense{}
	return upLm12SubstrateInvarianceR1Replay(seed, budget, schedule, "dense-array", state.set, state.get)
}

func upLm12SubstrateInvarianceR1ReplaySparse(seed, budget int, schedule [512]int) (int, int, int, string, []upLm12SubstrateInvarianceR1Row) {
	state := &upLm12SubstrateInvarianceR1Sparse{cells: make(map[int]bool, 512)}
	return upLm12SubstrateInvarianceR1Replay(seed, budget, schedule, "sparse-map", state.set, state.get)
}

func upLm12SubstrateInvarianceR1Replay(seed, budget int, schedule [512]int, substrate string, set func(int), get func(int) bool) (int, int, int, string, []upLm12SubstrateInvarianceR1Row) {
	rows := make([]upLm12SubstrateInvarianceR1Row, 0, len(schedule))
	acceptedWrites := 0
	deficitWrites := 0
	overrunRows := 0
	for ordinal, index := range schedule {
		accepted := acceptedWrites < budget
		if accepted {
			set(index)
			acceptedWrites++
		} else {
			deficitWrites++
		}
		if acceptedWrites > budget {
			overrunRows++
		}
		rows = append(rows, upLm12SubstrateInvarianceR1Row{seed, budget, substrate, ordinal, index, true, accepted, acceptedWrites, deficitWrites})
	}
	return acceptedWrites, deficitWrites, overrunRows, upLm12SubstrateInvarianceR1Digest(get), rows
}

func upLm12SubstrateInvarianceR1Schedule(seed int) [512]int {
	var schedule [512]int
	for i := range schedule {
		schedule[i] = i
	}
	stream := uint64(seed) ^ 0xD1B54A32D192ED03
	for i := len(schedule) - 1; i > 0; i-- {
		stream = upLm12SubstrateInvarianceR1SplitMix64(stream)
		j := int(stream % uint64(i+1))
		schedule[i], schedule[j] = schedule[j], schedule[i]
	}
	return schedule
}

func upLm12SubstrateInvarianceR1ValidateSchedule(schedule [512]int) (int, bool) {
	var seen [512]bool
	unique := 0
	valid := true
	for _, index := range schedule {
		if index < 0 || index >= len(seen) || seen[index] {
			valid = false
			continue
		}
		seen[index] = true
		unique++
	}
	return unique, valid && unique == len(seen)
}

func upLm12SubstrateInvarianceR1Digest(get func(int) bool) string {
	hash := sha256.New()
	var record [5]byte
	for index := 0; index < 512; index++ {
		binary.BigEndian.PutUint32(record[:4], uint32(index))
		if get(index) {
			record[4] = 1
		} else {
			record[4] = 0
		}
		var length [4]byte
		binary.BigEndian.PutUint32(length[:], uint32(len(record)))
		_, _ = hash.Write(length[:])
		_, _ = hash.Write(record[:])
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func upLm12SubstrateInvarianceR1SplitMix64(state uint64) uint64 {
	state += 0x9E3779B97F4A7C15
	z := state
	z = (z ^ (z >> 30)) * 0xBF58476D1CE4E5B9
	z = (z ^ (z >> 27)) * 0x94D049BB133111EB
	return z ^ (z >> 31)
}

func upLm12SubstrateInvarianceR1Abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
