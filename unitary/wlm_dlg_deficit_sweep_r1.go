package unitary

import (
	"crypto/sha256"
	"encoding/binary"
)

type wlmDlgDeficitSweepR1Row struct {
	Seed              uint64     `json:"seed"`
	ScheduleIndex     uint64     `json:"schedule_index"`
	DeficitLevel      uint64     `json:"deficit_level"`
	Substrate         string     `json:"substrate"`
	Step              uint64     `json:"step"`
	RequestedWork     uint64     `json:"requested_work"`
	Budget            uint64     `json:"budget"`
	CompletedWork     uint64     `json:"completed_work"`
	Deficit           uint64     `json:"deficit"`
	BoundedState      [16]uint64 `json:"bounded_state"`
	Overflow          uint64     `json:"overflow"`
	PriorState        uint64     `json:"prior_state"`
	AdmittedWork      uint64     `json:"admitted_work"`
	NextState         uint64     `json:"next_state"`
	ConservationError uint64     `json:"conservation_error"`
	DeficitLawError   uint64     `json:"deficit_law_error"`
}

type wlmDlgDeficitSweepR1Schedule struct {
	demand [64][16]uint64
	hash   [32]byte
}

type wlmDlgDeficitSweepR1State interface {
	get(int) uint64
	set(int, uint64)
	ordered() [16]uint64
}

type wlmDlgDeficitSweepR1Dense struct {
	values [16]uint64
}

func (s *wlmDlgDeficitSweepR1Dense) get(id int) uint64 { return s.values[id] }
func (s *wlmDlgDeficitSweepR1Dense) set(id int, value uint64) {
	s.values[id] = value
}
func (s *wlmDlgDeficitSweepR1Dense) ordered() [16]uint64 { return s.values }

type wlmDlgDeficitSweepR1Paged struct {
	pages [4]map[int]uint64
}

func (s *wlmDlgDeficitSweepR1Paged) get(id int) uint64 {
	return s.pages[id/4][id%4]
}

func (s *wlmDlgDeficitSweepR1Paged) set(id int, value uint64) {
	page, offset := id/4, id%4
	if value == 0 {
		delete(s.pages[page], offset)
		return
	}
	if s.pages[page] == nil {
		s.pages[page] = make(map[int]uint64, 4)
	}
	s.pages[page][offset] = value
}

func (s *wlmDlgDeficitSweepR1Paged) ordered() [16]uint64 {
	var values [16]uint64
	for id := range values {
		values[id] = s.get(id)
	}
	return values
}

// RunWlmDlgDeficitSweepR1 executes the preregistered paired deficit sweep.
func RunWlmDlgDeficitSweepR1() interface{} {
	const (
		steps        = 64
		stateBound   = uint64(4095)
		schedules    = 8
	)
	seeds := [...]uint64{101, 211, 307, 401, 503, 601, 701, 809}
	deficitLevels := [...]uint64{0, 1, 2, 4}
	rows := make([]wlmDlgDeficitSweepR1Row, 0, len(seeds)*schedules*len(deficitLevels)*2*steps)
	metrics := map[string]float64{
		"budget_overrun_rows":              0,
		"completed_substrate_case_runs":    0,
		"invalid_state_schedules":          0,
		"law_violation_cases":              0,
		"max_conservation_error":           0,
		"max_deficit_law_error":            0,
		"paired_case_count":                0,
		"paired_observable_mismatch_cases": 0,
		"row_count_mismatch_cases":         0,
		"state_cardinality_mismatch_cases": 0,
		"total_emitted_rows":               0,
	}

	for _, seed := range seeds {
		for scheduleIndex := 0; scheduleIndex < schedules; scheduleIndex++ {
			schedule := wlmDlgDeficitSweepR1MakeSchedule(seed, uint64(scheduleIndex))
			if !wlmDlgDeficitSweepR1ValidSchedule(schedule) {
				metrics["invalid_state_schedules"]++
			}
			for _, deficitLevel := range deficitLevels {
				metrics["paired_case_count"]++
				denseRows, denseViolation, denseOverruns := wlmDlgDeficitSweepR1Run(seed, uint64(scheduleIndex), deficitLevel, "dense-contiguous", &wlmDlgDeficitSweepR1Dense{}, schedule, stateBound)
				pagedRows, pagedViolation, pagedOverruns := wlmDlgDeficitSweepR1Run(seed, uint64(scheduleIndex), deficitLevel, "paged-sparse", &wlmDlgDeficitSweepR1Paged{}, schedule, stateBound)
				rows = append(rows, denseRows...)
				rows = append(rows, pagedRows...)
				metrics["completed_substrate_case_runs"] += 2
				metrics["budget_overrun_rows"] += float64(denseOverruns + pagedOverruns)

				if len(denseRows) != steps || len(pagedRows) != steps {
					metrics["row_count_mismatch_cases"]++
				}
				if !wlmDlgDeficitSweepR1RowsHaveCardinality(denseRows) || !wlmDlgDeficitSweepR1RowsHaveCardinality(pagedRows) {
					metrics["state_cardinality_mismatch_cases"]++
				}
				if !wlmDlgDeficitSweepR1SameObservables(denseRows, pagedRows) {
					metrics["paired_observable_mismatch_cases"]++
				}
				if denseViolation || pagedViolation {
					metrics["law_violation_cases"]++
				}
				for _, row := range denseRows {
					wlmDlgDeficitSweepR1UpdateMaxima(metrics, row)
				}
				for _, row := range pagedRows {
					wlmDlgDeficitSweepR1UpdateMaxima(metrics, row)
				}
			}
		}
	}
	metrics["total_emitted_rows"] = float64(len(rows))
	return map[string]interface{}{"rows": rows, "metrics": metrics}
}

func wlmDlgDeficitSweepR1MakeSchedule(seed, scheduleIndex uint64) wlmDlgDeficitSweepR1Schedule {
	var schedule wlmDlgDeficitSweepR1Schedule
	for step := range schedule.demand {
		for stateID := range schedule.demand[step] {
			schedule.demand[step][stateID] = uint64(wlmDlgDeficitSweepR1Demand(seed, scheduleIndex, uint64(step), uint64(stateID)))
		}
	}
	schedule.hash = sha256.Sum256(wlmDlgDeficitSweepR1ScheduleBytes(schedule.demand))
	return schedule
}

func wlmDlgDeficitSweepR1Demand(seed, scheduleIndex, step, stateID uint64) byte {
	fields := [][]byte{
		[]byte("WLM-DLG-DEFICIT-SWEEP-R1"),
		[]byte(wlmDlgDeficitSweepR1Decimal(seed)),
		[]byte(wlmDlgDeficitSweepR1Decimal(scheduleIndex)),
		[]byte(wlmDlgDeficitSweepR1Decimal(step)),
		[]byte(wlmDlgDeficitSweepR1Decimal(stateID)),
	}
	input := make([]byte, 0, 96)
	for _, field := range fields {
		var length [8]byte
		binary.BigEndian.PutUint64(length[:], uint64(len(field)))
		input = append(input, length[:]...)
		input = append(input, field...)
	}
	digest := sha256.Sum256(input)
	return digest[0] % 4
}

func wlmDlgDeficitSweepR1ScheduleBytes(demand [64][16]uint64) []byte {
	encoded := make([]byte, 0, 64*16)
	for _, step := range demand {
		for _, value := range step {
			encoded = append(encoded, byte(value))
		}
	}
	return encoded
}

func wlmDlgDeficitSweepR1ValidSchedule(schedule wlmDlgDeficitSweepR1Schedule) bool {
	if schedule.hash != sha256.Sum256(wlmDlgDeficitSweepR1ScheduleBytes(schedule.demand)) {
		return false
	}
	for _, step := range schedule.demand {
		for _, value := range step {
			if value > 3 {
				return false
			}
		}
	}
	return true
}

func wlmDlgDeficitSweepR1Run(seed, scheduleIndex, deficitLevel uint64, substrate string, state wlmDlgDeficitSweepR1State, schedule wlmDlgDeficitSweepR1Schedule, stateBound uint64) ([]wlmDlgDeficitSweepR1Row, bool, uint64) {
	rows := make([]wlmDlgDeficitSweepR1Row, 0, len(schedule.demand))
	lawViolation := false
	var budgetOverruns uint64
	for step, demand := range schedule.demand {
		var priorState, requestedWork uint64
		for id := range demand {
			priorState += state.get(id)
			requestedWork += demand[id]
		}
		budget := requestedWork
		if budget > deficitLevel {
			budget -= deficitLevel
		} else {
			budget = 0
		}
		remaining := budget
		var completedWork, overflow uint64
		for id := range demand {
			available := state.get(id) + demand[id]
			completed := available
			if completed > remaining {
				completed = remaining
			}
			remaining -= completed
			completedWork += completed
			next := available - completed
			if next > stateBound {
				overflow += next - stateBound
				next = stateBound
			}
			state.set(id, next)
		}
		ordered := state.ordered()
		var nextState uint64
		for _, value := range ordered {
			nextState += value
		}
		deficit := requestedWork - completedWork
		expectedDeficit := requestedWork - budget
		deficitLawError := wlmDlgDeficitSweepR1Difference(deficit, expectedDeficit)
		conservationError := wlmDlgDeficitSweepR1Difference(priorState+requestedWork, nextState+completedWork+overflow)
		if completedWork != budget || remaining != 0 {
			budgetOverruns++
			lawViolation = true
		}
		if deficitLawError != 0 || conservationError != 0 {
			lawViolation = true
		}
		rows = append(rows, wlmDlgDeficitSweepR1Row{
			Seed: seed, ScheduleIndex: scheduleIndex, DeficitLevel: deficitLevel, Substrate: substrate, Step: uint64(step),
			RequestedWork: requestedWork, Budget: budget, CompletedWork: completedWork, Deficit: deficit,
			BoundedState: ordered, Overflow: overflow, PriorState: priorState, AdmittedWork: requestedWork,
			NextState: nextState, ConservationError: conservationError, DeficitLawError: deficitLawError,
		})
	}
	return rows, lawViolation, budgetOverruns
}

func wlmDlgDeficitSweepR1UpdateMaxima(metrics map[string]float64, row wlmDlgDeficitSweepR1Row) {
	if float64(row.ConservationError) > metrics["max_conservation_error"] {
		metrics["max_conservation_error"] = float64(row.ConservationError)
	}
	if float64(row.DeficitLawError) > metrics["max_deficit_law_error"] {
		metrics["max_deficit_law_error"] = float64(row.DeficitLawError)
	}
}

func wlmDlgDeficitSweepR1RowsHaveCardinality(rows []wlmDlgDeficitSweepR1Row) bool {
	for _, row := range rows {
		if len(row.BoundedState) != 16 {
			return false
		}
	}
	return true
}

func wlmDlgDeficitSweepR1SameObservables(left, right []wlmDlgDeficitSweepR1Row) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		a, b := left[i], right[i]
		if a.Seed != b.Seed || a.ScheduleIndex != b.ScheduleIndex || a.DeficitLevel != b.DeficitLevel || a.Step != b.Step || a.RequestedWork != b.RequestedWork || a.Budget != b.Budget || a.CompletedWork != b.CompletedWork || a.Deficit != b.Deficit || a.BoundedState != b.BoundedState || a.Overflow != b.Overflow || a.PriorState != b.PriorState || a.AdmittedWork != b.AdmittedWork || a.NextState != b.NextState || a.ConservationError != b.ConservationError || a.DeficitLawError != b.DeficitLawError {
			return false
		}
	}
	return true
}

func wlmDlgDeficitSweepR1Difference(left, right uint64) uint64 {
	if left >= right {
		return left - right
	}
	return right - left
}

func wlmDlgDeficitSweepR1Decimal(value uint64) string {
	if value == 0 {
		return "0"
	}
	var digits [20]byte
	end := len(digits)
	for value > 0 {
		end--
		digits[end] = byte('0' + value%10)
		value /= 10
	}
	return string(digits[end:])
}
