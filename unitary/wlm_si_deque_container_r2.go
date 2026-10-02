package unitary

type wlmSiDequeContainerR2Row struct {
	Seed uint32
	ScheduleFamily string
	InitialOccupancy int
	RemovalBudget int
	Substrate string
	Step int
	RequestedPushes int
	AdmitFront bool
	RemoveFront bool
	AdmittedCount int
	RejectedCount int
	RemovedCount int
	DequeBefore int
	DequeAfter int
	Residual int
	CumulativeAdmitted int
	CumulativeRejected int
	CumulativeRemoved int
	DequeTokenIDs []int
	RemovedTokenIDs []int
}

type wlmSiDequeContainerR2Result struct {
	Rows []wlmSiDequeContainerR2Row
	Metrics map[string]float64
}

type wlmSiDequeContainerR2Op struct {
	Requested int
	AdmitFront bool
	RemoveFront bool
}

type wlmSiDequeContainerR2Deque interface {
	Len() int
	PushFront(int)
	PushBack(int)
	PopFront() int
	PopBack() int
	Snapshot() []int
}

type wlmSiDequeContainerR2SliceDeque struct {
	values []int
}

func (d *wlmSiDequeContainerR2SliceDeque) Len() int { return len(d.values) }
func (d *wlmSiDequeContainerR2SliceDeque) PushFront(v int) {
	d.values = append([]int{v}, d.values...)
}
func (d *wlmSiDequeContainerR2SliceDeque) PushBack(v int) {
	d.values = append(d.values, v)
}
func (d *wlmSiDequeContainerR2SliceDeque) PopFront() int {
	v := d.values[0]
	d.values = d.values[1:]
	return v
}
func (d *wlmSiDequeContainerR2SliceDeque) PopBack() int {
	i := len(d.values) - 1
	v := d.values[i]
	d.values = d.values[:i]
	return v
}
func (d *wlmSiDequeContainerR2SliceDeque) Snapshot() []int {
	return wlmSiDequeContainerR2Copy(d.values)
}

type wlmSiDequeContainerR2RingDeque struct {
	values [16]int
	head int
	length int
}

func (d *wlmSiDequeContainerR2RingDeque) Len() int { return d.length }
func (d *wlmSiDequeContainerR2RingDeque) PushFront(v int) {
	d.head = (d.head + 15) % 16
	d.values[d.head] = v
	d.length++
}
func (d *wlmSiDequeContainerR2RingDeque) PushBack(v int) {
	d.values[(d.head+d.length)%16] = v
	d.length++
}
func (d *wlmSiDequeContainerR2RingDeque) PopFront() int {
	v := d.values[d.head]
	d.head = (d.head + 1) % 16
	d.length--
	return v
}
func (d *wlmSiDequeContainerR2RingDeque) PopBack() int {
	i := (d.head + d.length - 1) % 16
	v := d.values[i]
	d.length--
	return v
}
func (d *wlmSiDequeContainerR2RingDeque) Snapshot() []int {
	out := make([]int, d.length)
	for i := range out {
		out[i] = d.values[(d.head+i)%16]
	}
	return out
}

// RunWlmSiDequeContainerR2 evaluates the frozen paired FIFO deque matrix.
func RunWlmSiDequeContainerR2() interface{} {
	seeds := [...]uint32{10009, 10037, 10061, 10091, 10111, 10141, 10169, 10193}
	schedules := [...]string{"alternating-ends", "paired-bursts", "front-pulse", "lcg32"}
	initials := [...]int{0, 5, 10, 15}
	budgets := [...]int{1, 2, 4, 8}
	metrics := map[string]float64{
		"paired_case_count": 0,
		"completed_substrate_case_runs": 0,
		"total_emitted_rows": 0,
		"paired_observable_mismatch_cases": 0,
		"row_count_mismatch_cases": 0,
		"state_cardinality_mismatch_cases": 0,
		"invalid_operation_schedules": 0,
		"budget_overrun_rows": 0,
		"law_violation_cases": 0,
		"max_conservation_error": 0,
		"max_residual_law_error": 0,
	}
	rows := make([]wlmSiDequeContainerR2Row, 0, 65536)

	for seedIndex, seed := range seeds {
		for _, schedule := range schedules {
			ops, valid := wlmSiDequeContainerR2Operations(seed, schedule)
			if !valid {
				metrics["invalid_operation_schedules"]++
				continue
			}
			for _, initial := range initials {
				for _, budget := range budgets {
					metrics["paired_case_count"]++
					var sliceRows, ringRows []wlmSiDequeContainerR2Row
					var sliceBad, ringBad bool
					if seedIndex%2 == 0 {
						sliceRows, sliceBad = wlmSiDequeContainerR2Run(seed, schedule, initial, budget, ops, "slice-backed-deque", true, metrics)
						ringRows, ringBad = wlmSiDequeContainerR2Run(seed, schedule, initial, budget, ops, "fixed-capacity-ring-deque", false, metrics)
					} else {
						ringRows, ringBad = wlmSiDequeContainerR2Run(seed, schedule, initial, budget, ops, "fixed-capacity-ring-deque", false, metrics)
						sliceRows, sliceBad = wlmSiDequeContainerR2Run(seed, schedule, initial, budget, ops, "slice-backed-deque", true, metrics)
					}
					rows = append(rows, sliceRows...)
					rows = append(rows, ringRows...)
					metrics["completed_substrate_case_runs"] += 2
					if sliceBad || ringBad {
						metrics["law_violation_cases"]++
					}
					if len(sliceRows) != len(ringRows) {
						metrics["row_count_mismatch_cases"]++
					}
					if !wlmSiDequeContainerR2EqualRows(sliceRows, ringRows) {
						metrics["paired_observable_mismatch_cases"]++
					}
					if !wlmSiDequeContainerR2EqualCardinalities(sliceRows, ringRows) {
						metrics["state_cardinality_mismatch_cases"]++
					}
				}
			}
		}
	}
	metrics["total_emitted_rows"] = float64(len(rows))
	return wlmSiDequeContainerR2Result{Rows: rows, Metrics: metrics}
}

func wlmSiDequeContainerR2Operations(seed uint32, family string) ([]wlmSiDequeContainerR2Op, bool) {
	ops := make([]wlmSiDequeContainerR2Op, 64)
	state := seed
	for step := range ops {
		var op wlmSiDequeContainerR2Op
		switch family {
		case "alternating-ends":
			op.Requested = 1 + int(seed%8)
			op.AdmitFront = (uint32(step)+seed)%2 == 0
			op.RemoveFront = !op.AdmitFront
		case "paired-bursts":
			p := (uint32(step) + seed) % 4
			if p < 2 {
				op.Requested = 16
			}
			op.AdmitFront = p < 2
			op.RemoveFront = (uint32(step)+seed)%2 == 0
		case "front-pulse":
			if (uint32(step)+seed)%8 == 2 {
				op.Requested = 16
			}
			op.AdmitFront = true
			op.RemoveFront = (uint32(step)+seed)%3 != 0
		case "lcg32":
			state = state*1664525 + 1013904223
			op.Requested = int(state % 17)
			state = state*1664525 + 1013904223
			op.AdmitFront = state%2 == 0
			state = state*1664525 + 1013904223
			op.RemoveFront = state%2 == 0
		default:
			return ops, false
		}
		if op.Requested < 0 || op.Requested > 16 {
			return ops, false
		}
		ops[step] = op
	}
	return ops, true
}

func wlmSiDequeContainerR2Run(seed uint32, schedule string, initial, budget int, ops []wlmSiDequeContainerR2Op, substrate string, sliceBacked bool, metrics map[string]float64) ([]wlmSiDequeContainerR2Row, bool) {
	var deque wlmSiDequeContainerR2Deque
	if sliceBacked {
		deque = &wlmSiDequeContainerR2SliceDeque{}
	} else {
		deque = &wlmSiDequeContainerR2RingDeque{}
	}
	for token := -initial; token < 0; token++ {
		deque.PushBack(token)
	}

	rows := make([]wlmSiDequeContainerR2Row, 0, len(ops))
	cumulativeAdmitted, cumulativeRejected, cumulativeRemoved := 0, 0, 0
	caseBad := false

	for step, op := range ops {
		before := deque.Len()
		admitted := wlmSiDequeContainerR2Min(op.Requested, 16-before)
		rejected := op.Requested - admitted
		for ordinal := 0; ordinal < admitted; ordinal++ {
			token := step*16 + ordinal
			if op.AdmitFront {
				deque.PushFront(token)
			} else {
				deque.PushBack(token)
			}
		}

		removedCount := wlmSiDequeContainerR2Min(deque.Len(), budget)
		removed := make([]int, removedCount)
		for i := range removed {
			if op.RemoveFront {
				removed[i] = deque.PopFront()
			} else {
				removed[i] = deque.PopBack()
			}
		}

		after := deque.Len()
		residual := wlmSiDequeContainerR2Max(0, before+admitted-budget)
		cumulativeAdmitted += admitted
		cumulativeRejected += rejected
		cumulativeRemoved += removedCount
		row := wlmSiDequeContainerR2Row{
			Seed: seed,
			ScheduleFamily: schedule,
			InitialOccupancy: initial,
			RemovalBudget: budget,
			Substrate: substrate,
			Step: step,
			RequestedPushes: op.Requested,
			AdmitFront: op.AdmitFront,
			RemoveFront: op.RemoveFront,
			AdmittedCount: admitted,
			RejectedCount: rejected,
			RemovedCount: removedCount,
			DequeBefore: before,
			DequeAfter: after,
			Residual: residual,
			CumulativeAdmitted: cumulativeAdmitted,
			CumulativeRejected: cumulativeRejected,
			CumulativeRemoved: cumulativeRemoved,
			DequeTokenIDs: deque.Snapshot(),
			RemovedTokenIDs: removed,
		}

		conservationError := wlmSiDequeContainerR2Abs(after - (before + admitted - removedCount))
		residualError := wlmSiDequeContainerR2Abs(after - residual)
		wlmSiDequeContainerR2MaxMetric(metrics, "max_conservation_error", conservationError)
		wlmSiDequeContainerR2MaxMetric(metrics, "max_residual_law_error", residualError)
		if removedCount > budget {
			metrics["budget_overrun_rows"]++
		}
		if conservationError != 0 || residualError != 0 || len(row.DequeTokenIDs) != after || after < 0 || after > 16 {
			caseBad = true
		}
		rows = append(rows, row)
	}
	return rows, caseBad
}

func wlmSiDequeContainerR2EqualRows(a, b []wlmSiDequeContainerR2Row) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		x, y := a[i], b[i]
		if x.Seed != y.Seed || x.ScheduleFamily != y.ScheduleFamily || x.InitialOccupancy != y.InitialOccupancy || x.RemovalBudget != y.RemovalBudget || x.Step != y.Step || x.RequestedPushes != y.RequestedPushes || x.AdmitFront != y.AdmitFront || x.RemoveFront != y.RemoveFront || x.AdmittedCount != y.AdmittedCount || x.RejectedCount != y.RejectedCount || x.RemovedCount != y.RemovedCount || x.DequeBefore != y.DequeBefore || x.DequeAfter != y.DequeAfter || x.Residual != y.Residual || x.CumulativeAdmitted != y.CumulativeAdmitted || x.CumulativeRejected != y.CumulativeRejected || x.CumulativeRemoved != y.CumulativeRemoved || !wlmSiDequeContainerR2IntsEqual(x.DequeTokenIDs, y.DequeTokenIDs) || !wlmSiDequeContainerR2IntsEqual(x.RemovedTokenIDs, y.RemovedTokenIDs) {
			return false
		}
	}
	return true
}

func wlmSiDequeContainerR2EqualCardinalities(a, b []wlmSiDequeContainerR2Row) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].DequeBefore != b[i].DequeBefore || a[i].DequeAfter != b[i].DequeAfter || len(a[i].DequeTokenIDs) != a[i].DequeAfter || len(b[i].DequeTokenIDs) != b[i].DequeAfter {
			return false
		}
	}
	return true
}

func wlmSiDequeContainerR2Copy(in []int) []int {
	out := make([]int, len(in))
	copy(out, in)
	return out
}
func wlmSiDequeContainerR2IntsEqual(a, b []int) bool {
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
func wlmSiDequeContainerR2Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func wlmSiDequeContainerR2Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func wlmSiDequeContainerR2Abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}
func wlmSiDequeContainerR2MaxMetric(metrics map[string]float64, name string, value int) {
	if float64(value) > metrics[name] {
		metrics[name] = float64(value)
	}
}
