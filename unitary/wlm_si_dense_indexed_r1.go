package unitary

import (
	"crypto/sha256"
	"encoding/binary"
	"sort"
)

type wlmSiDenseIndexedR1Row struct {
	Seed              uint64 `json:"seed"`
	CaseIndex         int    `json:"case_index"`
	Substrate         string `json:"substrate"`
	Step              int    `json:"step"`
	LogicalStateIndex int    `json:"logical_state_index"`
	ResourceBudget    int    `json:"resource_budget"`
	ActiveCount       int    `json:"active_count"`
	InactiveCount     int    `json:"inactive_count"`
	DemandUnits       int    `json:"demand_units"`
	DeficitUnits      int    `json:"deficit_units"`
	ConservationTotal int    `json:"conservation_total"`
}

type wlmSiDenseIndexedR1Observation struct {
	active, inactive, demand, deficit, conservation int
}

type wlmSiDenseIndexedR1DenseState struct {
	states [256]bool
}

type wlmSiDenseIndexedR1SparseState struct {
	active map[int]struct{}
}

// RunWlmSiDenseIndexedR1 replays every frozen sparse-indexed/dense-indexed pair.
func RunWlmSiDenseIndexedR1() interface{} {
	const (
		stateWidth    = 256
		casesPerSeed  = 32
		rowsPerRun    = 64
		substrateRuns = 2
	)

	seeds := [...]uint64{1103, 2207, 3313, 4421, 5531, 6653, 7757, 8861}
	budgets := [...]int{0, 1, 2, 4, 8, 16, 32, 64}
	metrics := map[string]float64{
		"budget_overrun_rows":               0,
		"completed_substrate_case_runs":     0,
		"invalid_state_schedules":           0,
		"law_violation_cases":               0,
		"max_conservation_error":            0,
		"max_deficit_law_error":             0,
		"paired_case_count":                 0,
		"paired_observable_mismatch_cases":  0,
		"row_count_mismatch_cases":          0,
		"state_cardinality_mismatch_cases":  0,
		"total_emitted_rows":                0,
	}
	rows := make([]wlmSiDenseIndexedR1Row, 0, len(seeds)*casesPerSeed*substrateRuns*rowsPerRun)

	for _, seed := range seeds {
		for caseIndex := 0; caseIndex < casesPerSeed; caseIndex++ {
			metrics["paired_case_count"]++
			schedule, canonical := wlmSiDenseIndexedR1Schedule(seed, caseIndex)
			if !wlmSiDenseIndexedR1ValidSchedule(schedule, canonical) {
				metrics["invalid_state_schedules"]++
				continue
			}

			initial := wlmSiDenseIndexedR1Initial(seed, caseIndex)
			oracle := wlmSiDenseIndexedR1Oracle(initial, schedule, budgets)
			sparseRows, sparseInvalid, sparseLaw := wlmSiDenseIndexedR1RunSparse(seed, caseIndex, schedule, budgets, initial, oracle)
			denseRows, denseInvalid, denseLaw := wlmSiDenseIndexedR1RunDense(seed, caseIndex, schedule, budgets, initial, oracle)
			rows = append(rows, sparseRows...)
			rows = append(rows, denseRows...)
			metrics["completed_substrate_case_runs"] += substrateRuns

			if sparseInvalid || denseInvalid {
				metrics["invalid_state_schedules"]++
			}
			if len(sparseRows) != rowsPerRun || len(denseRows) != rowsPerRun {
				metrics["row_count_mismatch_cases"]++
			}
			if sparseLaw || denseLaw {
				metrics["law_violation_cases"]++
			}
			if !wlmSiDenseIndexedR1SameCardinality(sparseRows, denseRows) {
				metrics["state_cardinality_mismatch_cases"]++
			}
			if !wlmSiDenseIndexedR1SameObservables(sparseRows, denseRows) {
				metrics["paired_observable_mismatch_cases"]++
			}
			for _, row := range sparseRows {
				wlmSiDenseIndexedR1UpdateErrors(metrics, row)
			}
			for _, row := range denseRows {
				wlmSiDenseIndexedR1UpdateErrors(metrics, row)
			}
		}
	}

	metrics["total_emitted_rows"] = float64(len(rows))
	frozenRowBudget := len(seeds) * casesPerSeed * substrateRuns * rowsPerRun
	if len(rows) > frozenRowBudget {
		metrics["budget_overrun_rows"] = float64(len(rows) - frozenRowBudget)
	}
	return map[string]interface{}{"rows": rows, "metrics": metrics}
}

func wlmSiDenseIndexedR1Schedule(seed uint64, caseIndex int) ([64]uint16, []byte) {
	var schedule [64]uint16
	canonical := make([]byte, len(schedule)*2)
	for step := range schedule {
		digest := wlmSiDenseIndexedR1HashTuple(seed, uint64(caseIndex), uint64(step))
		index := uint16(binary.BigEndian.Uint64(digest[:8]) % 256)
		schedule[step] = index
		binary.BigEndian.PutUint16(canonical[step*2:], index)
	}
	return schedule, canonical
}

func wlmSiDenseIndexedR1ValidSchedule(schedule [64]uint16, canonical []byte) bool {
	if len(canonical) != len(schedule)*2 {
		return false
	}
	for step, index := range schedule {
		if index >= 256 || binary.BigEndian.Uint16(canonical[step*2:step*2+2]) != index {
			return false
		}
	}
	return true
}

func wlmSiDenseIndexedR1Initial(seed uint64, caseIndex int) [256]bool {
	type candidate struct {
		index int
		hash  [32]byte
	}
	candidates := make([]candidate, 256)
	for index := range candidates {
		candidates[index] = candidate{index: index, hash: wlmSiDenseIndexedR1HashTuple(seed, uint64(caseIndex), uint64(index))}
	}
	sort.Slice(candidates, func(i, j int) bool {
		for byteIndex := range candidates[i].hash {
			if candidates[i].hash[byteIndex] != candidates[j].hash[byteIndex] {
				return candidates[i].hash[byteIndex] < candidates[j].hash[byteIndex]
			}
		}
		return candidates[i].index < candidates[j].index
	})
	var initial [256]bool
	for _, candidate := range candidates[:128] {
		initial[candidate.index] = true
	}
	return initial
}

func wlmSiDenseIndexedR1HashTuple(seed, caseIndex, value uint64) [32]byte {
	var encoded [24]byte
	binary.BigEndian.PutUint64(encoded[0:8], seed)
	binary.BigEndian.PutUint64(encoded[8:16], caseIndex)
	binary.BigEndian.PutUint64(encoded[16:24], value)
	return sha256.Sum256(encoded[:])
}

func wlmSiDenseIndexedR1Oracle(initial [256]bool, schedule [64]uint16, budgets [8]int) [64]wlmSiDenseIndexedR1Observation {
	state := initial
	var observations [64]wlmSiDenseIndexedR1Observation
	for step, index := range schedule {
		state[index] = !state[index]
		active := 0
		for _, set := range state {
			if set {
				active++
			}
		}
		budget := budgets[step%len(budgets)]
		observations[step] = wlmSiDenseIndexedR1Observation{
			active:       active,
			inactive:     256 - active,
			demand:       active,
			deficit:      wlmSiDenseIndexedR1Deficit(active, budget),
			conservation: 256,
		}
	}
	return observations
}

func wlmSiDenseIndexedR1RunSparse(seed uint64, caseIndex int, schedule [64]uint16, budgets [8]int, initial [256]bool, oracle [64]wlmSiDenseIndexedR1Observation) ([]wlmSiDenseIndexedR1Row, bool, bool) {
	state := wlmSiDenseIndexedR1SparseState{active: make(map[int]struct{}, 128)}
	for index, set := range initial {
		if set {
			state.active[index] = struct{}{}
		}
	}
	rows := make([]wlmSiDenseIndexedR1Row, 0, len(schedule))
	invalid, lawViolation := false, false
	for step, index := range schedule {
		key := int(index)
		if _, present := state.active[key]; present {
			delete(state.active, key)
		} else {
			state.active[key] = struct{}{}
		}
		active := len(state.active)
		budget := budgets[step%len(budgets)]
		row := wlmSiDenseIndexedR1Row{seed, caseIndex, "sparse-indexed", step, key, budget, active, 256 - active, active, wlmSiDenseIndexedR1Deficit(active, budget), 256}
		if active < 0 || active > 256 {
			invalid = true
		}
		if !wlmSiDenseIndexedR1Valid(row, oracle[step]) {
			lawViolation = true
		}
		rows = append(rows, row)
	}
	return rows, invalid, lawViolation
}

func wlmSiDenseIndexedR1RunDense(seed uint64, caseIndex int, schedule [64]uint16, budgets [8]int, initial [256]bool, oracle [64]wlmSiDenseIndexedR1Observation) ([]wlmSiDenseIndexedR1Row, bool, bool) {
	state := wlmSiDenseIndexedR1DenseState{states: initial}
	rows := make([]wlmSiDenseIndexedR1Row, 0, len(schedule))
	invalid, lawViolation := false, false
	for step, index := range schedule {
		state.states[index] = !state.states[index]
		active := 0
		for _, set := range state.states {
			if set {
				active++
			}
		}
		budget := budgets[step%len(budgets)]
		row := wlmSiDenseIndexedR1Row{seed, caseIndex, "dense-indexed", step, int(index), budget, active, 256 - active, active, wlmSiDenseIndexedR1Deficit(active, budget), 256}
		if active < 0 || active > 256 {
			invalid = true
		}
		if !wlmSiDenseIndexedR1Valid(row, oracle[step]) {
			lawViolation = true
		}
		rows = append(rows, row)
	}
	return rows, invalid, lawViolation
}

func wlmSiDenseIndexedR1Deficit(demand, budget int) int {
	if demand > budget {
		return demand - budget
	}
	return 0
}

func wlmSiDenseIndexedR1Valid(row wlmSiDenseIndexedR1Row, expected wlmSiDenseIndexedR1Observation) bool {
	return row.ActiveCount == expected.active &&
		row.InactiveCount == expected.inactive &&
		row.DemandUnits == expected.demand &&
		row.DeficitUnits == expected.deficit &&
		row.ConservationTotal == expected.conservation &&
		row.ConservationTotal == row.ActiveCount+row.InactiveCount &&
		row.DeficitUnits == wlmSiDenseIndexedR1Deficit(row.DemandUnits, row.ResourceBudget)
}

func wlmSiDenseIndexedR1SameCardinality(a, b []wlmSiDenseIndexedR1Row) bool {
	if len(a) != len(b) {
		return false
	}
	for index := range a {
		if a[index].ActiveCount != b[index].ActiveCount || a[index].InactiveCount != b[index].InactiveCount {
			return false
		}
	}
	return true
}

func wlmSiDenseIndexedR1SameObservables(a, b []wlmSiDenseIndexedR1Row) bool {
	if len(a) != len(b) {
		return false
	}
	for index := range a {
		left, right := a[index], b[index]
		if left.Seed != right.Seed ||
			left.CaseIndex != right.CaseIndex ||
			left.Step != right.Step ||
			left.LogicalStateIndex != right.LogicalStateIndex ||
			left.ResourceBudget != right.ResourceBudget ||
			left.ActiveCount != right.ActiveCount ||
			left.InactiveCount != right.InactiveCount ||
			left.DemandUnits != right.DemandUnits ||
			left.DeficitUnits != right.DeficitUnits ||
			left.ConservationTotal != right.ConservationTotal {
			return false
		}
	}
	return true
}

func wlmSiDenseIndexedR1UpdateErrors(metrics map[string]float64, row wlmSiDenseIndexedR1Row) {
	conservationError := wlmSiDenseIndexedR1Abs(row.ConservationTotal - (row.ActiveCount + row.InactiveCount))
	deficitLawError := wlmSiDenseIndexedR1Abs(row.DeficitUnits - wlmSiDenseIndexedR1Deficit(row.DemandUnits, row.ResourceBudget))
	if float64(conservationError) > metrics["max_conservation_error"] {
		metrics["max_conservation_error"] = float64(conservationError)
	}
	if float64(deficitLawError) > metrics["max_deficit_law_error"] {
		metrics["max_deficit_law_error"] = float64(deficitLawError)
	}
}

func wlmSiDenseIndexedR1Abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
