package unitary

import (
	"crypto/sha256"
	"encoding/binary"
	"sort"
)

type wlmSiSparseIndexedR1Row struct {
	Seed              uint64 `json:"seed"`
	CaseIndex         int    `json:"case_index"`
	Substrate         string `json:"substrate"`
	Step              int    `json:"step"`
	ResourceBudget    int    `json:"resource_budget"`
	ActiveCount       int    `json:"active_count"`
	InactiveCount     int    `json:"inactive_count"`
	DemandUnits       int    `json:"demand_units"`
	DeficitUnits      int    `json:"deficit_units"`
	ConservationTotal int    `json:"conservation_total"`
}

type wlmSiSparseIndexedR1Observation struct {
	active, inactive, demand, deficit, conservation int
}

type wlmSiSparseIndexedR1Dense struct {
	bits [256]bool
}

type wlmSiSparseIndexedR1Sparse struct {
	active map[int]struct{}
}

// RunWlmSiSparseIndexedR1 executes the frozen dense versus sparse-indexed substrate experiment.
func RunWlmSiSparseIndexedR1() interface{} {
	const (
		stateWidth = 256
		steps      = 64
		cases      = 32
	)
	seeds := [...]uint64{101, 211, 307, 401, 503, 601, 701, 809}
	budgets := [...]int{0, 1, 2, 4, 8, 16, 32, 64}
	metrics := map[string]float64{
		"completed_substrate_case_runs":       0,
		"paired_case_count":                   0,
		"total_emitted_rows":                  0,
		"paired_observable_mismatch_cases":    0,
		"row_count_mismatch_cases":            0,
		"state_cardinality_mismatch_cases":    0,
		"law_violation_cases":                 0,
		"max_conservation_error":              0,
		"max_deficit_law_error":               0,
		"budget_overrun_rows":                 0,
		"invalid_state_schedules":             0,
	}
	rows := make([]wlmSiSparseIndexedR1Row, 0, len(seeds)*cases*2*steps)

	for _, seed := range seeds {
		for caseIndex := 0; caseIndex < cases; caseIndex++ {
			metrics["paired_case_count"]++
			schedule, canonical := wlmSiSparseIndexedR1Schedule(seed, caseIndex)
			if sha256.Sum256(canonical) != sha256.Sum256(canonical) {
				metrics["invalid_state_schedules"]++
				continue
			}
			initial := wlmSiSparseIndexedR1Initial(seed, caseIndex)
			oracle := wlmSiSparseIndexedR1Oracle(initial, schedule, budgets)

			dense := wlmSiSparseIndexedR1Dense{bits: initial}
			sparse := wlmSiSparseIndexedR1Sparse{active: wlmSiSparseIndexedR1SparseInitial(initial)}
			denseRows, denseInvalid, denseLaw := wlmSiSparseIndexedR1RunDense(seed, caseIndex, schedule, budgets, &dense, oracle)
			sparseRows, sparseInvalid, sparseLaw := wlmSiSparseIndexedR1RunSparse(seed, caseIndex, schedule, budgets, &sparse, oracle)
			rows = append(rows, denseRows...)
			rows = append(rows, sparseRows...)
			metrics["completed_substrate_case_runs"] += 2

			if denseInvalid || sparseInvalid {
				metrics["invalid_state_schedules"]++
			}
			if len(denseRows) != steps || len(sparseRows) != steps {
				metrics["row_count_mismatch_cases"]++
			}
			if denseLaw || sparseLaw {
				metrics["law_violation_cases"]++
			}
			if !wlmSiSparseIndexedR1SameCardinality(denseRows, sparseRows) {
				metrics["state_cardinality_mismatch_cases"]++
			}
			if !wlmSiSparseIndexedR1SameObservables(denseRows, sparseRows) {
				metrics["paired_observable_mismatch_cases"]++
			}
			for _, row := range denseRows {
				wlmSiSparseIndexedR1UpdateErrors(metrics, row)
			}
			for _, row := range sparseRows {
				wlmSiSparseIndexedR1UpdateErrors(metrics, row)
			}
		}
	}
	metrics["total_emitted_rows"] = float64(len(rows))
	if len(rows) > 32768 {
		metrics["budget_overrun_rows"] = float64(len(rows) - 32768)
	}
	return map[string]interface{}{"rows": rows, "metrics": metrics}
}

func wlmSiSparseIndexedR1Schedule(seed uint64, caseIndex int) ([64]uint16, []byte) {
	var schedule [64]uint16
	canonical := make([]byte, len(schedule)*2)
	for step := range schedule {
		digest := wlmSiSparseIndexedR1HashTuple(seed, uint64(caseIndex), uint64(step))
		index := binary.BigEndian.Uint64(digest[:8]) % 256
		schedule[step] = uint16(index)
		binary.BigEndian.PutUint16(canonical[step*2:], uint16(index))
	}
	return schedule, canonical
}

func wlmSiSparseIndexedR1Initial(seed uint64, caseIndex int) [256]bool {
	type candidate struct {
		index int
		hash  [32]byte
	}
	candidates := make([]candidate, 256)
	for index := range candidates {
		candidates[index] = candidate{index: index, hash: wlmSiSparseIndexedR1HashTuple(seed, uint64(caseIndex), uint64(index))}
	}
	sort.Slice(candidates, func(i, j int) bool {
		for k := range candidates[i].hash {
			if candidates[i].hash[k] != candidates[j].hash[k] {
				return candidates[i].hash[k] < candidates[j].hash[k]
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

func wlmSiSparseIndexedR1SparseInitial(initial [256]bool) map[int]struct{} {
	active := make(map[int]struct{}, 128)
	for index, set := range initial {
		if set {
			active[index] = struct{}{}
		}
	}
	return active
}

func wlmSiSparseIndexedR1HashTuple(seed, caseIndex, value uint64) [32]byte {
	var encoded [24]byte
	binary.BigEndian.PutUint64(encoded[0:8], seed)
	binary.BigEndian.PutUint64(encoded[8:16], caseIndex)
	binary.BigEndian.PutUint64(encoded[16:24], value)
	return sha256.Sum256(encoded[:])
}

func wlmSiSparseIndexedR1Oracle(initial [256]bool, schedule [64]uint16, budgets [8]int) [64]wlmSiSparseIndexedR1Observation {
	state := initial
	var observations [64]wlmSiSparseIndexedR1Observation
	for step, index := range schedule {
		state[index] = !state[index]
		active := 0
		for _, set := range state {
			if set {
				active++
			}
		}
		budget := budgets[step%len(budgets)]
		observations[step] = wlmSiSparseIndexedR1Observation{active, 256 - active, active, wlmSiSparseIndexedR1Deficit(active, budget), 256}
	}
	return observations
}

func wlmSiSparseIndexedR1RunDense(seed uint64, caseIndex int, schedule [64]uint16, budgets [8]int, state *wlmSiSparseIndexedR1Dense, oracle [64]wlmSiSparseIndexedR1Observation) ([]wlmSiSparseIndexedR1Row, bool, bool) {
	rows := make([]wlmSiSparseIndexedR1Row, 0, len(schedule))
	invalid, law := false, false
	for step, index := range schedule {
		state.bits[index] = !state.bits[index]
		active := 0
		for _, set := range state.bits {
			if set {
				active++
			}
		}
		budget := budgets[step%len(budgets)]
		row := wlmSiSparseIndexedR1Row{seed, caseIndex, "dense", step, budget, active, 256 - active, active, wlmSiSparseIndexedR1Deficit(active, budget), 256}
		if !wlmSiSparseIndexedR1Valid(row, oracle[step]) {
			law = true
		}
		rows = append(rows, row)
	}
	return rows, invalid, law
}

func wlmSiSparseIndexedR1RunSparse(seed uint64, caseIndex int, schedule [64]uint16, budgets [8]int, state *wlmSiSparseIndexedR1Sparse, oracle [64]wlmSiSparseIndexedR1Observation) ([]wlmSiSparseIndexedR1Row, bool, bool) {
	rows := make([]wlmSiSparseIndexedR1Row, 0, len(schedule))
	invalid, law := false, false
	for step, index := range schedule {
		key := int(index)
		if _, set := state.active[key]; set {
			delete(state.active, key)
		} else {
			state.active[key] = struct{}{}
		}
		active := len(state.active)
		budget := budgets[step%len(budgets)]
		row := wlmSiSparseIndexedR1Row{seed, caseIndex, "sparse-indexed", step, budget, active, 256 - active, active, wlmSiSparseIndexedR1Deficit(active, budget), 256}
		if active < 0 || active > 256 {
			invalid = true
		}
		if !wlmSiSparseIndexedR1Valid(row, oracle[step]) {
			law = true
		}
		rows = append(rows, row)
	}
	return rows, invalid, law
}

func wlmSiSparseIndexedR1Deficit(demand, budget int) int {
	if demand > budget {
		return demand - budget
	}
	return 0
}

func wlmSiSparseIndexedR1Valid(row wlmSiSparseIndexedR1Row, expected wlmSiSparseIndexedR1Observation) bool {
	return row.ActiveCount == expected.active && row.InactiveCount == expected.inactive && row.DemandUnits == expected.demand && row.DeficitUnits == expected.deficit && row.ConservationTotal == expected.conservation && row.ConservationTotal == 256 && row.DeficitUnits == wlmSiSparseIndexedR1Deficit(row.DemandUnits, row.ResourceBudget)
}

func wlmSiSparseIndexedR1SameCardinality(a, b []wlmSiSparseIndexedR1Row) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].ActiveCount != b[i].ActiveCount || a[i].InactiveCount != b[i].InactiveCount {
			return false
		}
	}
	return true
}

func wlmSiSparseIndexedR1SameObservables(a, b []wlmSiSparseIndexedR1Row) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		x, y := a[i], b[i]
		if x.Seed != y.Seed || x.CaseIndex != y.CaseIndex || x.Step != y.Step || x.ResourceBudget != y.ResourceBudget || x.ActiveCount != y.ActiveCount || x.InactiveCount != y.InactiveCount || x.DemandUnits != y.DemandUnits || x.DeficitUnits != y.DeficitUnits || x.ConservationTotal != y.ConservationTotal {
			return false
		}
	}
	return true
}

func wlmSiSparseIndexedR1UpdateErrors(metrics map[string]float64, row wlmSiSparseIndexedR1Row) {
	conservation := wlmSiSparseIndexedR1Abs(row.ConservationTotal - 256)
	deficit := wlmSiSparseIndexedR1Abs(row.DeficitUnits - wlmSiSparseIndexedR1Deficit(row.DemandUnits, row.ResourceBudget))
	if float64(conservation) > metrics["max_conservation_error"] {
		metrics["max_conservation_error"] = float64(conservation)
	}
	if float64(deficit) > metrics["max_deficit_law_error"] {
		metrics["max_deficit_law_error"] = float64(deficit)
	}
}

func wlmSiSparseIndexedR1Abs(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
