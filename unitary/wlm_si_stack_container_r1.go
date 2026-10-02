package unitary

type wlmSiStackContainerR1Row struct {
	Seed               uint32 `json:"seed"`
	ScheduleFamily     string `json:"schedule_family"`
	InitialOccupancy   int    `json:"initial_occupancy"`
	PopBudget          int    `json:"pop_budget"`
	Substrate          string `json:"substrate"`
	Step               int    `json:"step"`
	RequestedPushes    int    `json:"requested_pushes"`
	AdmittedCount      int    `json:"admitted_count"`
	RejectedCount      int    `json:"rejected_count"`
	PoppedCount        int    `json:"popped_count"`
	StackBefore        int    `json:"s_t"`
	StackAfter         int    `json:"s_t_plus_1"`
	Residual           int    `json:"residual"`
	CumulativeAdmitted int    `json:"cumulative_admitted_count"`
	CumulativeRejected int    `json:"cumulative_rejected_count"`
	CumulativePopped   int    `json:"cumulative_popped_count"`
	StackTokenIDs      []int  `json:"stack_token_ids"`
	PoppedTokenIDs     []int  `json:"popped_token_ids"`
}
type wlmSiStackContainerR1Result struct {
	Rows    []wlmSiStackContainerR1Row `json:"rows"`
	Metrics map[string]float64         `json:"metrics"`
}
type wlmSiStackContainerR1LIFO interface {
	Len() int
	Push(int)
	Pop() int
	Snapshot() []int
}
type wlmSiStackContainerR1SliceLIFO struct{ values []int }

func (s *wlmSiStackContainerR1SliceLIFO) Len() int   { return len(s.values) }
func (s *wlmSiStackContainerR1SliceLIFO) Push(v int) { s.values = append(s.values, v) }
func (s *wlmSiStackContainerR1SliceLIFO) Pop() int {
	i := len(s.values) - 1
	v := s.values[i]
	s.values = s.values[:i]
	return v
}
func (s *wlmSiStackContainerR1SliceLIFO) Snapshot() []int { return wlmSiStackContainerR1Copy(s.values) }

type wlmSiStackContainerR1ArrayLIFO struct {
	values [16]int
	length int
}

func (s *wlmSiStackContainerR1ArrayLIFO) Len() int   { return s.length }
func (s *wlmSiStackContainerR1ArrayLIFO) Push(v int) { s.values[s.length] = v; s.length++ }
func (s *wlmSiStackContainerR1ArrayLIFO) Pop() int   { s.length--; return s.values[s.length] }
func (s *wlmSiStackContainerR1ArrayLIFO) Snapshot() []int {
	out := make([]int, s.length)
	for i := range out {
		out[i] = s.values[i]
	}
	return out
}

// RunWlmSiStackContainerR1 evaluates the frozen paired LIFO matrix.
func RunWlmSiStackContainerR1() interface{} {
	seeds := [...]uint32{131, 149, 167, 191, 223, 251, 277, 313}
	schedules := [...]string{"constant", "alternating", "pulse", "lcg32"}
	initials := [...]int{0, 5, 10, 15}
	budgets := [...]int{1, 2, 4, 8}
	metrics := map[string]float64{"paired_case_count": 0, "completed_substrate_case_runs": 0, "total_emitted_rows": 0, "paired_observable_mismatch_cases": 0, "row_count_mismatch_cases": 0, "state_cardinality_mismatch_cases": 0, "invalid_state_schedules": 0, "budget_overrun_rows": 0, "law_violation_cases": 0, "max_conservation_error": 0, "max_residual_law_error": 0}
	rows := make([]wlmSiStackContainerR1Row, 0, 65536)
	for seedIndex, seed := range seeds {
		for _, schedule := range schedules {
			pushes, valid := wlmSiStackContainerR1Pushes(seed, schedule)
			if !valid {
				metrics["invalid_state_schedules"]++
				continue
			}
			for _, initial := range initials {
				for _, budget := range budgets {
					metrics["paired_case_count"]++
					var sliceRows, arrayRows []wlmSiStackContainerR1Row
					var sliceBad, arrayBad bool
					if seedIndex%2 == 0 {
						sliceRows, sliceBad = wlmSiStackContainerR1Run(seed, schedule, initial, budget, pushes, "slice-backed-lifo", true, metrics)
						arrayRows, arrayBad = wlmSiStackContainerR1Run(seed, schedule, initial, budget, pushes, "fixed-capacity-array-lifo", false, metrics)
					} else {
						arrayRows, arrayBad = wlmSiStackContainerR1Run(seed, schedule, initial, budget, pushes, "fixed-capacity-array-lifo", false, metrics)
						sliceRows, sliceBad = wlmSiStackContainerR1Run(seed, schedule, initial, budget, pushes, "slice-backed-lifo", true, metrics)
					}
					rows = append(rows, sliceRows...)
					rows = append(rows, arrayRows...)
					metrics["completed_substrate_case_runs"] += 2
					if sliceBad || arrayBad {
						metrics["law_violation_cases"]++
					}
					if len(sliceRows) != len(arrayRows) {
						metrics["row_count_mismatch_cases"]++
					}
					if !wlmSiStackContainerR1EqualRows(sliceRows, arrayRows) {
						metrics["paired_observable_mismatch_cases"]++
					}
					if !wlmSiStackContainerR1EqualCardinalities(sliceRows, arrayRows) {
						metrics["state_cardinality_mismatch_cases"]++
					}
				}
			}
		}
	}
	metrics["total_emitted_rows"] = float64(len(rows))
	return wlmSiStackContainerR1Result{Rows: rows, Metrics: metrics}
}
func wlmSiStackContainerR1Pushes(seed uint32, family string) ([]int, bool) {
	trace := make([]int, 64)
	state := seed
	for step := range trace {
		switch family {
		case "constant":
			trace[step] = 1 + int(seed%8)
		case "alternating":
			if (uint32(step)+seed)%2 == 0 {
				trace[step] = 16
			}
		case "pulse":
			if (uint32(step)+seed)%8 == 1 {
				trace[step] = 16
			}
		case "lcg32":
			state = state*1664525 + 1013904223
			trace[step] = int(state % 17)
		default:
			return trace, false
		}
		if trace[step] < 0 || trace[step] > 16 {
			return trace, false
		}
	}
	return trace, true
}
func wlmSiStackContainerR1Run(seed uint32, schedule string, initial, budget int, pushes []int, substrate string, sliceBacked bool, metrics map[string]float64) ([]wlmSiStackContainerR1Row, bool) {
	var stack wlmSiStackContainerR1LIFO
	if sliceBacked {
		stack = &wlmSiStackContainerR1SliceLIFO{}
	} else {
		stack = &wlmSiStackContainerR1ArrayLIFO{}
	}
	for token := -initial; token < 0; token++ {
		stack.Push(token)
	}
	rows := make([]wlmSiStackContainerR1Row, 0, len(pushes))
	cumulativeAdmitted, cumulativeRejected, cumulativePopped := 0, 0, 0
	caseBad := false
	for step, requested := range pushes {
		before := stack.Len()
		admitted := wlmSiStackContainerR1Min(requested, 16-before)
		rejected := requested - admitted
		for ordinal := 0; ordinal < admitted; ordinal++ {
			stack.Push(step*16 + ordinal)
		}
		poppedCount := wlmSiStackContainerR1Min(stack.Len(), budget)
		popped := make([]int, poppedCount)
		for i := range popped {
			popped[i] = stack.Pop()
		}
		after := stack.Len()
		residual := wlmSiStackContainerR1Max(0, before+admitted-budget)
		cumulativeAdmitted += admitted
		cumulativeRejected += rejected
		cumulativePopped += poppedCount
		row := wlmSiStackContainerR1Row{Seed: seed, ScheduleFamily: schedule, InitialOccupancy: initial, PopBudget: budget, Substrate: substrate, Step: step, RequestedPushes: requested, AdmittedCount: admitted, RejectedCount: rejected, PoppedCount: poppedCount, StackBefore: before, StackAfter: after, Residual: residual, CumulativeAdmitted: cumulativeAdmitted, CumulativeRejected: cumulativeRejected, CumulativePopped: cumulativePopped, StackTokenIDs: stack.Snapshot(), PoppedTokenIDs: popped}
		conservationError := wlmSiStackContainerR1Abs(after - (before + admitted - poppedCount))
		residualError := wlmSiStackContainerR1Abs(after - residual)
		wlmSiStackContainerR1MaxMetric(metrics, "max_conservation_error", conservationError)
		wlmSiStackContainerR1MaxMetric(metrics, "max_residual_law_error", residualError)
		if poppedCount > budget {
			metrics["budget_overrun_rows"]++
		}
		if conservationError != 0 || residualError != 0 || len(row.StackTokenIDs) != after || after < 0 || after > 16 {
			caseBad = true
		}
		rows = append(rows, row)
	}
	return rows, caseBad
}
func wlmSiStackContainerR1EqualRows(a, b []wlmSiStackContainerR1Row) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		x, y := a[i], b[i]
		if x.Seed != y.Seed || x.ScheduleFamily != y.ScheduleFamily || x.InitialOccupancy != y.InitialOccupancy || x.PopBudget != y.PopBudget || x.Step != y.Step || x.RequestedPushes != y.RequestedPushes || x.AdmittedCount != y.AdmittedCount || x.RejectedCount != y.RejectedCount || x.PoppedCount != y.PoppedCount || x.StackBefore != y.StackBefore || x.StackAfter != y.StackAfter || x.Residual != y.Residual || x.CumulativeAdmitted != y.CumulativeAdmitted || x.CumulativeRejected != y.CumulativeRejected || x.CumulativePopped != y.CumulativePopped || !wlmSiStackContainerR1IntsEqual(x.StackTokenIDs, y.StackTokenIDs) || !wlmSiStackContainerR1IntsEqual(x.PoppedTokenIDs, y.PoppedTokenIDs) {
			return false
		}
	}
	return true
}
func wlmSiStackContainerR1EqualCardinalities(a, b []wlmSiStackContainerR1Row) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].StackBefore != b[i].StackBefore || a[i].StackAfter != b[i].StackAfter || len(a[i].StackTokenIDs) != a[i].StackAfter || len(b[i].StackTokenIDs) != b[i].StackAfter {
			return false
		}
	}
	return true
}
func wlmSiStackContainerR1Copy(in []int) []int {
	out := make([]int, len(in))
	copy(out, in)
	return out
}
func wlmSiStackContainerR1IntsEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
func wlmSiStackContainerR1Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func wlmSiStackContainerR1Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func wlmSiStackContainerR1Abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
func wlmSiStackContainerR1MaxMetric(metrics map[string]float64, name string, value int) {
	if float64(value) > metrics[name] {
		metrics[name] = float64(value)
	}
}
