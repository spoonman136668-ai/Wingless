package unitary

import "encoding/binary"

type wlmSiDenseSparseR2Row struct {
	Seed              uint64 `json:"seed"`
	StateCardinality  int    `json:"state_cardinality"`
	ResourceBudget    int    `json:"resource_budget"`
	ScheduleID        int    `json:"schedule_id"`
	Step              int    `json:"step"`
	LogicalStateIndex uint32 `json:"logical_state_index"`
	Substrate         string `json:"substrate"`
	ActiveCardinality int    `json:"active_cardinality"`
	Checksum          uint64 `json:"canonical_active_index_checksum"`
	Demand            int    `json:"demand"`
	Served            int    `json:"served"`
	Deficit           int    `json:"deficit"`
}

type wlmSiDenseSparseR2Result struct {
	Rows    []wlmSiDenseSparseR2Row `json:"rows"`
	Metrics map[string]float64       `json:"metrics"`
}

type wlmSiDenseSparseR2Dense struct{ words []uint64 }
type wlmSiDenseSparseR2Sparse map[uint32]struct{}

func RunWlmSiDenseSparseR2() interface{} {
	seeds := [...]uint64{11, 23, 47, 97, 193, 389, 769, 1543}
	dims := [...]int{8, 32, 128, 512}
	budgets := [...]int{0, 1, 4, 16}
	metrics := map[string]float64{
		"paired_case_count": 0, "completed_substrate_case_runs": 0, "total_emitted_rows": 0,
		"paired_observable_mismatch_cases": 0, "row_count_mismatch_cases": 0,
		"state_cardinality_mismatch_cases": 0, "invalid_state_schedules": 0,
		"law_violation_cases": 0, "max_conservation_error": 0,
		"max_deficit_law_error": 0, "budget_overrun_rows": 0,
	}
	rows := make([]wlmSiDenseSparseR2Row, 0, 65536)
	for _, seed := range seeds {
		for _, dimension := range dims {
			for scheduleID := 0; scheduleID < 4; scheduleID++ {
				trace, valid := wlmSiDenseSparseR2Trace(seed, dimension, scheduleID)
				if !valid { metrics["invalid_state_schedules"]++; continue }
				metrics["paired_case_count"]++
				for _, budget := range budgets {
					d := wlmSiDenseSparseR2Dense{words: make([]uint64, (dimension+63)/64)}
					s := make(wlmSiDenseSparseR2Sparse)
					a, al := wlmSiDenseSparseR2Run(seed, dimension, scheduleID, budget, trace, "dense-bitset-control", func(i uint32) { d.words[i/64] ^= uint64(1) << (i % 64) }, func() (int, uint64) { return wlmSiDenseSparseR2DenseObs(d, dimension) })
					b, bl := wlmSiDenseSparseR2Run(seed, dimension, scheduleID, budget, trace, "sparse-index-set-treatment", func(i uint32) { if _, ok := s[i]; ok { delete(s, i) } else { s[i] = struct{}{} } }, func() (int, uint64) { return wlmSiDenseSparseR2SparseObs(s, dimension) })
					rows = append(rows, a...); rows = append(rows, b...)
					metrics["completed_substrate_case_runs"] += 2
					if al || bl { metrics["law_violation_cases"]++ }
					if len(a) != len(b) { metrics["row_count_mismatch_cases"]++ }
					if !wlmSiDenseSparseR2Equal(a, b) { metrics["paired_observable_mismatch_cases"]++ }
					if !wlmSiDenseSparseR2CardinalityEqual(a, b) { metrics["state_cardinality_mismatch_cases"]++ }
					for _, r := range append(a, b...) { wlmSiDenseSparseR2Errors(metrics, r) }
				}
			}
		}
	}
	metrics["total_emitted_rows"] = float64(len(rows))
	return wlmSiDenseSparseR2Result{Rows: rows, Metrics: metrics}
}

func wlmSiDenseSparseR2Trace(seed uint64, dimension, scheduleID int) ([]uint32, bool) {
	trace := make([]uint32, 64)
	state := wlmSiDenseSparseR2Seed(seed, dimension, scheduleID)
	for i := range trace { state = wlmSiDenseSparseR2Next(state); trace[i] = uint32(state % uint64(dimension)); if int(trace[i]) >= dimension { return trace, false } }
	return trace, true
}

func wlmSiDenseSparseR2Seed(seed uint64, dimension, scheduleID int) uint64 {
	v := seed ^ 0x574c4d2d53492d44
	v ^= uint64(dimension) * 0x9e3779b97f4a7c15
	v ^= uint64(scheduleID+1) * 0xbf58476d1ce4e5b9
	return v
}
func wlmSiDenseSparseR2Next(v uint64) uint64 { v += 0x9e3779b97f4a7c15; z := v; z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9; z = (z ^ (z >> 27)) * 0x94d049bb133111eb; return z ^ (z >> 31) }

func wlmSiDenseSparseR2Run(seed uint64, dimension, scheduleID, budget int, trace []uint32, substrate string, toggle func(uint32), observe func() (int, uint64)) ([]wlmSiDenseSparseR2Row, bool) {
	rows := make([]wlmSiDenseSparseR2Row, 0, len(trace)); invalid := false
	for step, index := range trace {
		toggle(index); active, checksum := observe(); demand := active; served := demand; if served > budget { served = budget }; deficit := demand - served
		row := wlmSiDenseSparseR2Row{seed, dimension, budget, scheduleID, step, index, substrate, active, checksum, demand, served, deficit}
		if index >= uint32(dimension) || demand != served+deficit || deficit != wlmSiDenseSparseR2Max(0, demand-budget) { invalid = true }
		rows = append(rows, row)
	}
	return rows, invalid
}
func wlmSiDenseSparseR2DenseObs(d wlmSiDenseSparseR2Dense, n int) (int, uint64) { active := 0; h := uint64(14695981039346656037); for i := 0; i < n; i++ { if d.words[i/64]&(uint64(1)<<(uint(i)%64)) != 0 { active++; var b [4]byte; binary.LittleEndian.PutUint32(b[:], uint32(i)); for _, x := range b { h ^= uint64(x); h *= 1099511628211 } } }; return active, h }
func wlmSiDenseSparseR2SparseObs(s wlmSiDenseSparseR2Sparse, n int) (int, uint64) { active := 0; h := uint64(14695981039346656037); for i := 0; i < n; i++ { if _, ok := s[uint32(i)]; ok { active++; var b [4]byte; binary.LittleEndian.PutUint32(b[:], uint32(i)); for _, x := range b { h ^= uint64(x); h *= 1099511628211 } } }; return active, h }
func wlmSiDenseSparseR2Equal(a, b []wlmSiDenseSparseR2Row) bool { if len(a) != len(b) { return false }; for i := range a { if a[i].ActiveCardinality != b[i].ActiveCardinality || a[i].Checksum != b[i].Checksum || a[i].Demand != b[i].Demand || a[i].Served != b[i].Served || a[i].Deficit != b[i].Deficit { return false } }; return true }
func wlmSiDenseSparseR2CardinalityEqual(a, b []wlmSiDenseSparseR2Row) bool { if len(a) != len(b) { return false }; for i := range a { if a[i].ActiveCardinality != b[i].ActiveCardinality { return false } }; return true }
func wlmSiDenseSparseR2Errors(m map[string]float64, r wlmSiDenseSparseR2Row) { c := wlmSiDenseSparseR2Abs(r.Demand - r.Served - r.Deficit); d := wlmSiDenseSparseR2Abs(r.Deficit - wlmSiDenseSparseR2Max(0, r.Demand-r.ResourceBudget)); if float64(c) > m["max_conservation_error"] { m["max_conservation_error"] = float64(c) }; if float64(d) > m["max_deficit_law_error"] { m["max_deficit_law_error"] = float64(d) } }
func wlmSiDenseSparseR2Max(a, b int) int { if a > b { return a }; return b }
func wlmSiDenseSparseR2Abs(v int) int { if v < 0 { return -v }; return v }
