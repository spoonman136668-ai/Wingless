package unitary

import (
	"crypto/sha256"
	"encoding/binary"
)

type upLm16PagedSubstrateInvarianceR1Row struct {
	Seed              uint64     `json:"seed"`
	CaseIndex         uint64     `json:"case_index"`
	Substrate         string     `json:"substrate"`
	Tick              uint64     `json:"tick"`
	Ordered16Backlogs [16]uint64 `json:"ordered_16_backlogs"`
	TotalDemand       uint64     `json:"total_demand"`
	TotalServed       uint64     `json:"total_served"`
	TotalBacklog      uint64     `json:"total_backlog"`
	RemainingCapacity uint64     `json:"remaining_capacity"`
	ConservationError uint64     `json:"conservation_error"`
}

type upLm16PagedSubstrateInvarianceR1State interface {
	get(int) uint64
	set(int, uint64)
	ordered() [16]uint64
}

type upLm16PagedSubstrateInvarianceR1Dense struct {
	values [16]uint64
}

func (s *upLm16PagedSubstrateInvarianceR1Dense) get(id int) uint64 { return s.values[id] }
func (s *upLm16PagedSubstrateInvarianceR1Dense) set(id int, value uint64) {
	s.values[id] = value
}
func (s *upLm16PagedSubstrateInvarianceR1Dense) ordered() [16]uint64 { return s.values }

type upLm16PagedSubstrateInvarianceR1Paged struct {
	pages [4]map[int]uint64
}

func (s *upLm16PagedSubstrateInvarianceR1Paged) get(id int) uint64 {
	return s.pages[id/4][id%4]
}

func (s *upLm16PagedSubstrateInvarianceR1Paged) set(id int, value uint64) {
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

func (s *upLm16PagedSubstrateInvarianceR1Paged) ordered() [16]uint64 {
	var values [16]uint64
	for id := 0; id < len(values); id++ {
		values[id] = s.get(id)
	}
	return values
}

// RunUpLm16PagedSubstrateInvarianceR1 executes the frozen dense-versus-paged
// bounded-state matrix and returns its complete canonical observations.
func RunUpLm16PagedSubstrateInvarianceR1() interface{} {
	const (
		ticks        = 64
		casesPerSeed = 32
		capacity     = uint64(24)
		backlogBound = uint64(4095)
	)
	seeds := [...]uint64{104729, 130363, 155921, 181081, 206369, 231709, 257053, 282407}
	rows := make([]upLm16PagedSubstrateInvarianceR1Row, 0, len(seeds)*casesPerSeed*2*ticks)
	metrics := map[string]float64{
		"paired_case_count":                 0,
		"completed_substrate_case_runs":     0,
		"total_emitted_rows":                0,
		"invalid_state_schedules":           0,
		"row_count_mismatch_cases":          0,
		"state_cardinality_mismatch_cases":  0,
		"paired_observable_mismatch_cases":  0,
		"max_conservation_error":            0,
		"law_violation_cases":               0,
		"budget_overrun_rows":               0,
	}

	for _, seed := range seeds {
		for caseIndex := 0; caseIndex < casesPerSeed; caseIndex++ {
			schedule := upLm16PagedSubstrateInvarianceR1Schedule(seed, uint64(caseIndex))
			metrics["paired_case_count"]++
			if !upLm16PagedSubstrateInvarianceR1ValidSchedule(schedule) {
				metrics["invalid_state_schedules"]++
			}

			var denseRows, pagedRows []upLm16PagedSubstrateInvarianceR1Row
			var denseLaw, pagedLaw bool
			var denseBudget, pagedBudget uint64
			if (caseIndex%2) == 0 {
				denseRows, denseLaw, denseBudget = upLm16PagedSubstrateInvarianceR1Run(seed, uint64(caseIndex), "dense-contiguous", &upLm16PagedSubstrateInvarianceR1Dense{}, schedule, capacity, backlogBound)
				pagedRows, pagedLaw, pagedBudget = upLm16PagedSubstrateInvarianceR1Run(seed, uint64(caseIndex), "paged-sparse", &upLm16PagedSubstrateInvarianceR1Paged{}, schedule, capacity, backlogBound)
			} else {
				pagedRows, pagedLaw, pagedBudget = upLm16PagedSubstrateInvarianceR1Run(seed, uint64(caseIndex), "paged-sparse", &upLm16PagedSubstrateInvarianceR1Paged{}, schedule, capacity, backlogBound)
				denseRows, denseLaw, denseBudget = upLm16PagedSubstrateInvarianceR1Run(seed, uint64(caseIndex), "dense-contiguous", &upLm16PagedSubstrateInvarianceR1Dense{}, schedule, capacity, backlogBound)
			}
			rows = append(rows, denseRows...)
			rows = append(rows, pagedRows...)
			metrics["completed_substrate_case_runs"] += 2
			metrics["budget_overrun_rows"] += float64(denseBudget + pagedBudget)

			if len(denseRows) != ticks || len(pagedRows) != ticks {
				metrics["row_count_mismatch_cases"]++
			}
			if !upLm16PagedSubstrateInvarianceR1RowsHaveCardinality(denseRows) || !upLm16PagedSubstrateInvarianceR1RowsHaveCardinality(pagedRows) {
				metrics["state_cardinality_mismatch_cases"]++
			}
			if !upLm16PagedSubstrateInvarianceR1SameObservables(denseRows, pagedRows) {
				metrics["paired_observable_mismatch_cases"]++
			}
			if denseLaw || pagedLaw {
				metrics["law_violation_cases"]++
			}
			for _, row := range append(denseRows, pagedRows...) {
				if float64(row.ConservationError) > metrics["max_conservation_error"] {
					metrics["max_conservation_error"] = float64(row.ConservationError)
				}
			}
		}
	}
	metrics["total_emitted_rows"] = float64(len(rows))
	return map[string]interface{}{"rows": rows, "metrics": metrics}
}

func upLm16PagedSubstrateInvarianceR1Schedule(seed, caseIndex uint64) [64][16]uint64 {
	var schedule [64][16]uint64
	for tick := range schedule {
		for stateID := range schedule[tick] {
			schedule[tick][stateID] = uint64(upLm16PagedSubstrateInvarianceR1Demand(seed, caseIndex, uint64(tick), uint64(stateID)))
		}
	}
	return schedule
}

func upLm16PagedSubstrateInvarianceR1Demand(seed, caseIndex, tick, stateID uint64) byte {
	fields := [][]byte{
		[]byte("UP-LM16-PAGED-SUBSTRATE-INVARIANCE-R1"),
		[]byte(upLm16PagedSubstrateInvarianceR1Decimal(seed)),
		[]byte(upLm16PagedSubstrateInvarianceR1Decimal(caseIndex)),
		[]byte(upLm16PagedSubstrateInvarianceR1Decimal(tick)),
		[]byte(upLm16PagedSubstrateInvarianceR1Decimal(stateID)),
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

func upLm16PagedSubstrateInvarianceR1Decimal(value uint64) string {
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

func upLm16PagedSubstrateInvarianceR1ValidSchedule(schedule [64][16]uint64) bool {
	for _, tick := range schedule {
		for _, demand := range tick {
			if demand > 3 {
				return false
			}
		}
	}
	return true
}

func upLm16PagedSubstrateInvarianceR1Run(seed, caseIndex uint64, substrate string, state upLm16PagedSubstrateInvarianceR1State, schedule [64][16]uint64, capacity, backlogBound uint64) ([]upLm16PagedSubstrateInvarianceR1Row, bool, uint64) {
	rows := make([]upLm16PagedSubstrateInvarianceR1Row, 0, len(schedule))
	lawViolation := false
	var budgetOverruns uint64
	for tick, demand := range schedule {
		var priorTotal, totalDemand, totalServed uint64
		for id := 0; id < 16; id++ {
			priorTotal += state.get(id)
			totalDemand += demand[id]
		}
		remaining := capacity
		for id := 0; id < 16; id++ {
			available := state.get(id) + demand[id]
			served := available
			if served > remaining {
				served = remaining
			}
			remaining -= served
			next := available - served
			if next > backlogBound {
				lawViolation = true
			}
			state.set(id, next)
			totalServed += served
		}
		ordered := state.ordered()
		var totalBacklog uint64
		for _, backlog := range ordered {
			totalBacklog += backlog
		}
		expected := priorTotal + totalDemand - totalServed
		conservationError := upLm16PagedSubstrateInvarianceR1Difference(totalBacklog, expected)
		if conservationError != 0 {
			lawViolation = true
		}
		if totalServed > capacity || remaining+totalServed != capacity {
			budgetOverruns++
			lawViolation = true
		}
		rows = append(rows, upLm16PagedSubstrateInvarianceR1Row{seed, caseIndex, substrate, uint64(tick), ordered, totalDemand, totalServed, totalBacklog, remaining, conservationError})
	}
	return rows, lawViolation, budgetOverruns
}

func upLm16PagedSubstrateInvarianceR1Difference(left, right uint64) uint64 {
	if left >= right {
		return left - right
	}
	return right - left
}

func upLm16PagedSubstrateInvarianceR1RowsHaveCardinality(rows []upLm16PagedSubstrateInvarianceR1Row) bool {
	for _, row := range rows {
		if len(row.Ordered16Backlogs) != 16 {
			return false
		}
	}
	return true
}

func upLm16PagedSubstrateInvarianceR1SameObservables(left, right []upLm16PagedSubstrateInvarianceR1Row) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		a, b := left[i], right[i]
		if a.Seed != b.Seed || a.CaseIndex != b.CaseIndex || a.Tick != b.Tick || a.TotalDemand != b.TotalDemand || a.TotalServed != b.TotalServed || a.TotalBacklog != b.TotalBacklog || a.RemainingCapacity != b.RemainingCapacity || a.ConservationError != b.ConservationError || a.Ordered16Backlogs != b.Ordered16Backlogs {
			return false
		}
	}
	return true
}
