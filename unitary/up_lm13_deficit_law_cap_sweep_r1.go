package unitary

type upLm13DeficitLawCapSweepR1Row struct {
	Seed                int  `json:"seed"`
	Cap                 int  `json:"cap"`
	DemandOrdinal       int  `json:"demand_ordinal"`
	StateIndex          int  `json:"state_index"`
	Attempted           bool `json:"attempted"`
	Accepted            bool `json:"accepted"`
	AcceptedWritesSoFar int  `json:"accepted_writes_so_far"`
	DeficitWritesSoFar  int  `json:"deficit_writes_so_far"`
}

type upLm13DeficitLawCapSweepR1Metrics struct {
	SubstrateCaseRuns             float64 `json:"substrate_case_runs"`
	CapLevelCount                 float64 `json:"cap_level_count"`
	TotalEmittedRows              float64 `json:"total_emitted_rows"`
	AggregateDemandWrites         float64 `json:"aggregate_demand_writes"`
	AggregateAcceptedWrites       float64 `json:"aggregate_accepted_writes"`
	AggregateDeficit              float64 `json:"aggregate_deficit"`
	AcceptedWriteMismatchCases    float64 `json:"accepted_write_mismatch_cases"`
	DeficitMismatchCases          float64 `json:"deficit_mismatch_cases"`
	LawViolationCases             float64 `json:"law_violation_cases"`
	BudgetOverrunRows             float64 `json:"budget_overrun_rows"`
	InvalidStateSchedules         float64 `json:"invalid_state_schedules"`
	RowCountMismatchCases         float64 `json:"row_count_mismatch_cases"`
	StateCardinalityMismatchCases float64 `json:"state_cardinality_mismatch_cases"`
	MaxConservationError          float64 `json:"max_conservation_error"`
}

// RunUpLm13DeficitLawCapSweepR1 replays every preregistered seed-cap case from fresh state.
func RunUpLm13DeficitLawCapSweepR1() interface{} {
	seeds := [...]int{104729, 130363, 155921, 181081, 205759, 231481, 257053, 282571}
	caps := [...]int{64, 128, 192, 256, 320, 384, 448, 512}
	const demandWrites = 512

	rows := make([]upLm13DeficitLawCapSweepR1Row, 0, len(seeds)*len(caps)*demandWrites)
	caseRuns := 0
	aggregateDemand := 0
	aggregateAccepted := 0
	aggregateDeficit := 0
	acceptedMismatches := 0
	deficitMismatches := 0
	lawViolations := 0
	budgetOverruns := 0
	invalidSchedules := 0
	rowCountMismatches := 0
	stateCardinalityMismatches := 0
	maxConservationError := 0

	for _, seed := range seeds {
		schedule := upLm13DeficitLawCapSweepR1Schedule(seed)
		uniqueStates, scheduleValid := upLm13DeficitLawCapSweepR1ValidateSchedule(schedule)
		if !scheduleValid {
			invalidSchedules++
		}

		for _, cap := range caps {
			caseRuns++
			caseRows, acceptedWrites, deficitWrites, overrunRows, writtenStates := upLm13DeficitLawCapSweepR1Replay(seed, cap, schedule)
			rows = append(rows, caseRows...)
			aggregateDemand += len(caseRows)
			aggregateAccepted += acceptedWrites
			aggregateDeficit += deficitWrites
			budgetOverruns += overrunRows

			expectedAccepted := cap
			expectedDeficit := demandWrites - cap
			acceptedMismatch := acceptedWrites != expectedAccepted
			deficitMismatch := deficitWrites != expectedDeficit
			conservationError := upLm13DeficitLawCapSweepR1Abs(demandWrites - (acceptedWrites + deficitWrites))
			rowMismatch := len(caseRows) != demandWrites
			stateMismatch := !scheduleValid || uniqueStates != demandWrites || writtenStates > uniqueStates

			if acceptedMismatch {
				acceptedMismatches++
			}
			if deficitMismatch {
				deficitMismatches++
			}
			if rowMismatch {
				rowCountMismatches++
			}
			if stateMismatch {
				stateCardinalityMismatches++
			}
			if conservationError > maxConservationError {
				maxConservationError = conservationError
			}
			if acceptedMismatch || deficitMismatch || conservationError != 0 || overrunRows != 0 || rowMismatch || stateMismatch {
				lawViolations++
			}
		}
	}

	return map[string]interface{}{
		"rows": rows,
		"metrics": upLm13DeficitLawCapSweepR1Metrics{
			SubstrateCaseRuns:             float64(caseRuns),
			CapLevelCount:                 float64(len(caps)),
			TotalEmittedRows:              float64(len(rows)),
			AggregateDemandWrites:         float64(aggregateDemand),
			AggregateAcceptedWrites:       float64(aggregateAccepted),
			AggregateDeficit:              float64(aggregateDeficit),
			AcceptedWriteMismatchCases:    float64(acceptedMismatches),
			DeficitMismatchCases:          float64(deficitMismatches),
			LawViolationCases:             float64(lawViolations),
			BudgetOverrunRows:             float64(budgetOverruns),
			InvalidStateSchedules:         float64(invalidSchedules),
			RowCountMismatchCases:         float64(rowCountMismatches),
			StateCardinalityMismatchCases: float64(stateCardinalityMismatches),
			MaxConservationError:          float64(maxConservationError),
		},
	}
}

func upLm13DeficitLawCapSweepR1Replay(seed, cap int, schedule [512]int) ([]upLm13DeficitLawCapSweepR1Row, int, int, int, int) {
	state := make(map[int]bool, len(schedule))
	rows := make([]upLm13DeficitLawCapSweepR1Row, 0, len(schedule))
	acceptedWrites := 0
	deficitWrites := 0
	overrunRows := 0

	for ordinal, stateIndex := range schedule {
		accepted := acceptedWrites < cap
		if accepted {
			state[stateIndex] = true
			acceptedWrites++
		} else {
			deficitWrites++
		}
		if acceptedWrites > cap {
			overrunRows++
		}
		rows = append(rows, upLm13DeficitLawCapSweepR1Row{
			Seed:                seed,
			Cap:                 cap,
			DemandOrdinal:       ordinal,
			StateIndex:          stateIndex,
			Attempted:           true,
			Accepted:            accepted,
			AcceptedWritesSoFar: acceptedWrites,
			DeficitWritesSoFar:  deficitWrites,
		})
	}
	return rows, acceptedWrites, deficitWrites, overrunRows, len(state)
}

func upLm13DeficitLawCapSweepR1Schedule(seed int) [512]int {
	var schedule [512]int
	for index := range schedule {
		schedule[index] = index
	}
	stream := uint64(seed) ^ 0xD1B54A32D192ED03
	for index := len(schedule) - 1; index > 0; index-- {
		stream = upLm13DeficitLawCapSweepR1SplitMix64(stream)
		other := int(stream % uint64(index+1))
		schedule[index], schedule[other] = schedule[other], schedule[index]
	}
	return schedule
}

func upLm13DeficitLawCapSweepR1ValidateSchedule(schedule [512]int) (int, bool) {
	var seen [512]bool
	unique := 0
	valid := true
	for _, stateIndex := range schedule {
		if stateIndex < 0 || stateIndex >= len(seen) || seen[stateIndex] {
			valid = false
			continue
		}
		seen[stateIndex] = true
		unique++
	}
	return unique, valid && unique == len(seen)
}

func upLm13DeficitLawCapSweepR1SplitMix64(state uint64) uint64 {
	state += 0x9E3779B97F4A7C15
	z := state
	z = (z ^ (z >> 30)) * 0xBF58476D1CE4E5B9
	z = (z ^ (z >> 27)) * 0x94D049BB133111EB
	return z ^ (z >> 31)
}

func upLm13DeficitLawCapSweepR1Abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
