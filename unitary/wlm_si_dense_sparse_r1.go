package unitary

type wlmSiDenseSparseR1Row struct {
	Seed               uint64 `json:"seed"`
	Capacity           int    `json:"capacity"`
	Schedule           string `json:"schedule"`
	Substrate          string `json:"substrate"`
	Step               int    `json:"step"`
	PreOccupied        int    `json:"pre_occupied"`
	PreFree            int    `json:"pre_free"`
	Released           int    `json:"released"`
	FreeBefore         int    `json:"free_before"`
	Demand             int    `json:"demand"`
	Accepted           int    `json:"accepted"`
	Deficit            int    `json:"deficit"`
	PostOccupied       int    `json:"post_occupied"`
	PostFree           int    `json:"post_free"`
	CumulativeAccepted int    `json:"cumulative_accepted"`
	CumulativeDeficit  int    `json:"cumulative_deficit"`
	StateDigest        uint64 `json:"state_digest"`
}

type wlmSiDenseSparseR1Dense struct{ slots []bool }
type wlmSiDenseSparseR1Sparse struct{ slots []uint16 }

// RunWlmSiDenseSparseR1 executes the frozen dense-versus-sparse substrate experiment.
func RunWlmSiDenseSparseR1() interface{} {
	const steps = 64
	seeds := [...]uint64{104729, 130363, 155921, 181081, 205019, 229447, 253907, 279967}
	capacities := [...]int{8, 16, 32, 64}
	schedules := [...]string{"zero", "half", "exact", "plus_one", "alternating", "ramp", "burst", "seeded"}
	metrics := map[string]float64{
		"completed_substrate_case_runs":      0,
		"paired_case_count":                  0,
		"total_emitted_rows":                 0,
		"paired_observable_mismatch_cases":   0,
		"state_cardinality_mismatch_cases":   0,
		"row_count_mismatch_cases":           0,
		"law_violation_cases":                0,
		"max_conservation_error":              0,
		"max_deficit_law_error":              0,
		"invalid_state_schedules":            0,
		"budget_overrun_rows":                0,
	}
	rows := make([]wlmSiDenseSparseR1Row, 0, len(seeds)*len(capacities)*len(schedules)*2*steps)

	for _, seed := range seeds {
		for _, capacity := range capacities {
			for _, schedule := range schedules {
				metrics["paired_case_count"]++
				initial := wlmSiDenseSparseR1Initial(seed, capacity)
				dense := wlmSiDenseSparseR1Dense{slots: append([]bool(nil), initial...)}
				sparse := wlmSiDenseSparseR1Sparse{slots: wlmSiDenseSparseR1SparseFromInitial(initial)}
				denseRows, denseInvalid, denseLaw, denseOverrun := wlmSiDenseSparseR1RunDense(seed, capacity, schedule, &dense)
				sparseRows, sparseInvalid, sparseLaw, sparseOverrun := wlmSiDenseSparseR1RunSparse(seed, capacity, schedule, &sparse)
				rows = append(rows, denseRows...)
				rows = append(rows, sparseRows...)
				metrics["completed_substrate_case_runs"] += 2
				metrics["budget_overrun_rows"] += float64(denseOverrun + sparseOverrun)
				if denseInvalid || sparseInvalid {
					metrics["invalid_state_schedules"]++
				}
				if len(denseRows) != steps || len(sparseRows) != steps {
					metrics["row_count_mismatch_cases"]++
				}
				if denseLaw || sparseLaw {
					metrics["law_violation_cases"]++
				}
				if !wlmSiDenseSparseR1SameCardinality(denseRows, sparseRows) {
					metrics["state_cardinality_mismatch_cases"]++
				}
				if !wlmSiDenseSparseR1SameObservables(denseRows, sparseRows) {
					metrics["paired_observable_mismatch_cases"]++
				}
				for _, row := range denseRows {
					wlmSiDenseSparseR1UpdateErrors(metrics, row)
				}
				for _, row := range sparseRows {
					wlmSiDenseSparseR1UpdateErrors(metrics, row)
				}
			}
		}
	}
	metrics["total_emitted_rows"] = float64(len(rows))
	if len(rows) > 32768 {
		metrics["budget_overrun_rows"] += float64(len(rows) - 32768)
	}
	return map[string]interface{}{"rows": rows, "metrics": metrics}
}

func wlmSiDenseSparseR1Initial(seed uint64, capacity int) []bool {
	slots := make([]bool, capacity)
	for slot := range slots {
		slots[slot] = wlmSiDenseSparseR1PRF(seed, 0, uint64(slot), 0x494E4954)%4 == 0
	}
	return slots
}

func wlmSiDenseSparseR1SparseFromInitial(initial []bool) []uint16 {
	slots := make([]uint16, 0, len(initial))
	for i, occupied := range initial {
		if occupied {
			slots = append(slots, uint16(i))
		}
	}
	return slots
}

func wlmSiDenseSparseR1RunDense(seed uint64, capacity int, schedule string, state *wlmSiDenseSparseR1Dense) ([]wlmSiDenseSparseR1Row, bool, bool, int) {
	rows := make([]wlmSiDenseSparseR1Row, 0, 64)
	invalid, law, overrun := len(state.slots) != capacity, false, 0
	acceptedTotal, deficitTotal := 0, 0
	for step := 0; step < 64; step++ {
		pre := wlmSiDenseSparseR1DenseCount(state.slots)
		released := 0
		for slot, occupied := range state.slots {
			if occupied && wlmSiDenseSparseR1PRF(seed, uint64(step), uint64(slot), 0x52454C53)%5 == 0 {
				state.slots[slot] = false
				released++
			}
		}
		freeBefore := capacity - wlmSiDenseSparseR1DenseCount(state.slots)
		demand := wlmSiDenseSparseR1Demand(seed, step, capacity, schedule)
		accepted := demand
		if accepted > freeBefore { accepted = freeBefore }
		for slot := range state.slots {
			if accepted == 0 { break }
			if !state.slots[slot] { state.slots[slot] = true; accepted-- }
		}
		accepted = wlmSiDenseSparseR1DenseCount(state.slots) - (pre - released)
		deficit := demand - accepted
		acceptedTotal += accepted
		deficitTotal += deficit
		post := wlmSiDenseSparseR1DenseCount(state.slots)
		row := wlmSiDenseSparseR1Row{seed, capacity, schedule, "dense-bool", step, pre, capacity - pre, released, freeBefore, demand, accepted, deficit, post, capacity - post, acceptedTotal, deficitTotal, wlmSiDenseSparseR1DenseDigest(state.slots)}
		if !wlmSiDenseSparseR1ValidLaw(row) { law = true }
		rows = append(rows, row)
	}
	return rows, invalid, law, overrun
}

func wlmSiDenseSparseR1RunSparse(seed uint64, capacity int, schedule string, state *wlmSiDenseSparseR1Sparse) ([]wlmSiDenseSparseR1Row, bool, bool, int) {
	rows := make([]wlmSiDenseSparseR1Row, 0, 64)
	invalid, law, overrun := !wlmSiDenseSparseR1SparseValid(state.slots, capacity), false, 0
	acceptedTotal, deficitTotal := 0, 0
	for step := 0; step < 64; step++ {
		pre := len(state.slots)
		next := make([]uint16, 0, len(state.slots))
		released := 0
		for _, slot := range state.slots {
			if wlmSiDenseSparseR1PRF(seed, uint64(step), uint64(slot), 0x52454C53)%5 == 0 { released++ } else { next = append(next, slot) }
		}
		state.slots = next
		freeBefore := capacity - len(state.slots)
		demand := wlmSiDenseSparseR1Demand(seed, step, capacity, schedule)
		accepted := demand
		if accepted > freeBefore { accepted = freeBefore }
		need := accepted
		merged := make([]uint16, 0, len(state.slots)+accepted)
		pos := 0
		for slot := 0; slot < capacity; slot++ {
			if pos < len(state.slots) && int(state.slots[pos]) == slot { merged = append(merged, uint16(slot)); pos++; continue }
			if need > 0 { merged = append(merged, uint16(slot)); need-- }
		}
		state.slots = merged
		deficit := demand - accepted
		acceptedTotal += accepted
		deficitTotal += deficit
		post := len(state.slots)
		row := wlmSiDenseSparseR1Row{seed, capacity, schedule, "sparse-sorted-u16", step, pre, capacity - pre, released, freeBefore, demand, accepted, deficit, post, capacity - post, acceptedTotal, deficitTotal, wlmSiDenseSparseR1SparseDigest(state.slots)}
		if !wlmSiDenseSparseR1SparseValid(state.slots, capacity) { invalid = true }
		if !wlmSiDenseSparseR1ValidLaw(row) { law = true }
		rows = append(rows, row)
	}
	return rows, invalid, law, overrun
}

func wlmSiDenseSparseR1Demand(seed uint64, step, capacity int, schedule string) int {
	switch schedule {
	case "zero": return 0
	case "half": return capacity / 2
	case "exact": return capacity
	case "plus_one": return capacity + 1
	case "alternating": if step%2 == 0 { return 0 }; return capacity + 1
	case "ramp": return step % (capacity + 2)
	case "burst": if step%8 == 0 { return capacity + 2 }; return capacity / 2
	default: return int(wlmSiDenseSparseR1PRF(seed, uint64(step), 0, 0x444D4E44) % uint64(capacity+3))
	}
}

func wlmSiDenseSparseR1PRF(seed, step, slot, domain uint64) uint64 {
	z := seed ^ domain ^ (step << 32) ^ slot
	z += 0x9e3779b97f4a7c15
	z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
	z = (z ^ (z >> 27)) * 0x94d049bb133111eb
	return z ^ (z >> 31)
}

func wlmSiDenseSparseR1DenseCount(slots []bool) int { n := 0; for _, occupied := range slots { if occupied { n++ } }; return n }
func wlmSiDenseSparseR1SparseValid(slots []uint16, capacity int) bool { for i, slot := range slots { if int(slot) >= capacity || (i > 0 && slots[i-1] >= slot) { return false } }; return true }
func wlmSiDenseSparseR1DenseDigest(slots []bool) uint64 { h := uint64(14695981039346656037); for i, occupied := range slots { if occupied { h ^= uint64(byte(i)); h *= 1099511628211; h ^= uint64(byte(i >> 8)); h *= 1099511628211 } }; return h }
func wlmSiDenseSparseR1SparseDigest(slots []uint16) uint64 { h := uint64(14695981039346656037); for _, slot := range slots { h ^= uint64(byte(slot)); h *= 1099511628211; h ^= uint64(byte(slot >> 8)); h *= 1099511628211 }; return h }
func wlmSiDenseSparseR1ValidLaw(r wlmSiDenseSparseR1Row) bool { expected := r.Demand; if expected > r.FreeBefore { expected = r.FreeBefore }; return r.Accepted == expected && r.Deficit == r.Demand-r.Accepted && r.PostOccupied == r.PreOccupied-r.Released+r.Accepted && r.PostOccupied+r.PostFree == r.Capacity }
func wlmSiDenseSparseR1SameCardinality(a, b []wlmSiDenseSparseR1Row) bool { if len(a) != len(b) { return false }; for i := range a { if a[i].PreOccupied != b[i].PreOccupied || a[i].PostOccupied != b[i].PostOccupied || a[i].PostFree != b[i].PostFree { return false } }; return true }
func wlmSiDenseSparseR1SameObservables(a, b []wlmSiDenseSparseR1Row) bool { if len(a) != len(b) { return false }; for i := range a { x, y := a[i], b[i]; if x.Seed != y.Seed || x.Capacity != y.Capacity || x.Schedule != y.Schedule || x.Step != y.Step || x.PreOccupied != y.PreOccupied || x.PreFree != y.PreFree || x.Released != y.Released || x.FreeBefore != y.FreeBefore || x.Demand != y.Demand || x.Accepted != y.Accepted || x.Deficit != y.Deficit || x.PostOccupied != y.PostOccupied || x.PostFree != y.PostFree || x.CumulativeAccepted != y.CumulativeAccepted || x.CumulativeDeficit != y.CumulativeDeficit || x.StateDigest != y.StateDigest { return false } }; return true }
func wlmSiDenseSparseR1UpdateErrors(metrics map[string]float64, r wlmSiDenseSparseR1Row) { conservation := wlmSiDenseSparseR1Abs(r.PostOccupied - (r.PreOccupied-r.Released+r.Accepted)); if x := wlmSiDenseSparseR1Abs(r.PostOccupied+r.PostFree-r.Capacity); x > conservation { conservation = x }; deficit := wlmSiDenseSparseR1Abs(r.Accepted-wlmSiDenseSparseR1Min(r.Demand, r.FreeBefore)); if x := wlmSiDenseSparseR1Abs(r.Deficit-(r.Demand-r.Accepted)); x > deficit { deficit = x }; if float64(conservation) > metrics["max_conservation_error"] { metrics["max_conservation_error"] = float64(conservation) }; if float64(deficit) > metrics["max_deficit_law_error"] { metrics["max_deficit_law_error"] = float64(deficit) } }
func wlmSiDenseSparseR1Abs(x int) int { if x < 0 { return -x }; return x }
func wlmSiDenseSparseR1Min(a, b int) int { if a < b { return a }; return b }
