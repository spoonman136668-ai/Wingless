package unitary

type wlmSiQueueContainerR1Row struct {
	Seed                uint32 `json:"seed"`
	ScheduleFamily      string `json:"schedule_family"`
	InitialOccupancy    int    `json:"initial_occupancy"`
	ServiceBudget       int    `json:"service_budget"`
	Substrate           string `json:"substrate"`
	Step                int    `json:"step"`
	Arrivals            int    `json:"arrivals"`
	AdmittedCount       int    `json:"admitted_count"`
	RejectedCount       int    `json:"rejected_count"`
	ServedCount         int    `json:"served_count"`
	QueueBefore         int    `json:"q_t"`
	QueueAfter          int    `json:"q_t_plus_1"`
	Deficit             int    `json:"deficit"`
	CumulativeAdmitted  int    `json:"cumulative_admitted_count"`
	CumulativeRejected  int    `json:"cumulative_rejected_count"`
	CumulativeServed    int    `json:"cumulative_served_count"`
	QueuedTokenIDs      []int  `json:"queued_token_ids"`
	ServicedTokenIDs    []int  `json:"serviced_token_ids"`
}

type wlmSiQueueContainerR1Result struct {
	Rows    []wlmSiQueueContainerR1Row `json:"rows"`
	Metrics map[string]float64         `json:"metrics"`
}

type wlmSiQueueContainerR1FIFO interface {
	Len() int
	Push(int)
	Pop() int
	Snapshot() []int
}

type wlmSiQueueContainerR1SliceFIFO struct{ values []int }

func (q *wlmSiQueueContainerR1SliceFIFO) Len() int { return len(q.values) }
func (q *wlmSiQueueContainerR1SliceFIFO) Push(v int) { q.values = append(q.values, v) }
func (q *wlmSiQueueContainerR1SliceFIFO) Pop() int {
	v := q.values[0]
	copy(q.values, q.values[1:])
	q.values = q.values[:len(q.values)-1]
	return v
}
func (q *wlmSiQueueContainerR1SliceFIFO) Snapshot() []int {
	return wlmSiQueueContainerR1Copy(q.values)
}

type wlmSiQueueContainerR1RingFIFO struct {
	values [16]int
	head   int
	length int
}

func (q *wlmSiQueueContainerR1RingFIFO) Len() int { return q.length }
func (q *wlmSiQueueContainerR1RingFIFO) Push(v int) {
	q.values[(q.head+q.length)%len(q.values)] = v
	q.length++
}
func (q *wlmSiQueueContainerR1RingFIFO) Pop() int {
	v := q.values[q.head]
	q.head = (q.head + 1) % len(q.values)
	q.length--
	return v
}
func (q *wlmSiQueueContainerR1RingFIFO) Snapshot() []int {
	out := make([]int, q.length)
	for i := range out {
		out[i] = q.values[(q.head+i)%len(q.values)]
	}
	return out
}

// RunWlmSiQueueContainerR1 evaluates the frozen slice-backed versus ring-backed
// FIFO matrix and returns its complete row-level observations and metrics.
func RunWlmSiQueueContainerR1() interface{} {
	seeds := [...]uint32{11, 23, 37, 53, 71, 89, 107, 127}
	schedules := [...]string{"constant", "alternating", "pulse", "lcg32"}
	initialOccupancies := [...]int{0, 5, 10, 15}
	budgets := [...]int{1, 2, 4, 8}
	metrics := map[string]float64{
		"paired_case_count":                 0,
		"completed_substrate_case_runs":     0,
		"total_emitted_rows":                0,
		"paired_observable_mismatch_cases":  0,
		"row_count_mismatch_cases":          0,
		"state_cardinality_mismatch_cases":  0,
		"invalid_state_schedules":           0,
		"budget_overrun_rows":                0,
		"law_violation_cases":                0,
		"max_conservation_error":             0,
		"max_deficit_law_error":              0,
	}
	rows := make([]wlmSiQueueContainerR1Row, 0, 65536)

	for seedIndex, seed := range seeds {
		for _, schedule := range schedules {
			arrivals, valid := wlmSiQueueContainerR1Arrivals(seed, schedule)
			if !valid {
				metrics["invalid_state_schedules"]++
				continue
			}
			for _, initial := range initialOccupancies {
				for _, budget := range budgets {
					metrics["paired_case_count"]++
					var sliceRows, ringRows []wlmSiQueueContainerR1Row
					var sliceBad, ringBad bool
					if seedIndex%2 == 0 {
						sliceRows, sliceBad = wlmSiQueueContainerR1Run(seed, schedule, initial, budget, arrivals, "slice-backed-fifo", true, metrics)
						ringRows, ringBad = wlmSiQueueContainerR1Run(seed, schedule, initial, budget, arrivals, "fixed-capacity-ring-buffer-fifo", false, metrics)
					} else {
						ringRows, ringBad = wlmSiQueueContainerR1Run(seed, schedule, initial, budget, arrivals, "fixed-capacity-ring-buffer-fifo", false, metrics)
						sliceRows, sliceBad = wlmSiQueueContainerR1Run(seed, schedule, initial, budget, arrivals, "slice-backed-fifo", true, metrics)
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
					if !wlmSiQueueContainerR1EqualRows(sliceRows, ringRows) {
						metrics["paired_observable_mismatch_cases"]++
					}
					if !wlmSiQueueContainerR1EqualCardinalities(sliceRows, ringRows) {
						metrics["state_cardinality_mismatch_cases"]++
					}
				}
			}
		}
	}
	metrics["total_emitted_rows"] = float64(len(rows))
	return wlmSiQueueContainerR1Result{Rows: rows, Metrics: metrics}
}

func wlmSiQueueContainerR1Arrivals(seed uint32, family string) ([]int, bool) {
	trace := make([]int, 64)
	state := seed
	for step := range trace {
		switch family {
		case "constant":
			trace[step] = 1 + int(seed%8)
		case "alternating":
			if (uint32(step)+seed)%2 == 1 {
				trace[step] = 16
			}
		case "pulse":
			if (uint32(step)+seed)%8 == 0 {
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

func wlmSiQueueContainerR1Run(seed uint32, schedule string, initial, budget int, arrivals []int, substrate string, sliceBacked bool, metrics map[string]float64) ([]wlmSiQueueContainerR1Row, bool) {
	var queue wlmSiQueueContainerR1FIFO
	if sliceBacked {
		queue = &wlmSiQueueContainerR1SliceFIFO{}
	} else {
		queue = &wlmSiQueueContainerR1RingFIFO{}
	}
	for token := -initial; token < 0; token++ {
		queue.Push(token)
	}
	rows := make([]wlmSiQueueContainerR1Row, 0, len(arrivals))
	cumulativeAdmitted, cumulativeRejected, cumulativeServed := 0, 0, 0
	caseBad := false
	for step, arriving := range arrivals {
		before := queue.Len()
		admitted := wlmSiQueueContainerR1Min(arriving, 16-before)
		rejected := arriving - admitted
		for ordinal := 0; ordinal < admitted; ordinal++ {
			queue.Push(step*16 + ordinal)
		}
		served := wlmSiQueueContainerR1Min(queue.Len(), budget)
		serviced := make([]int, served)
		for i := range serviced {
			serviced[i] = queue.Pop()
		}
		after := queue.Len()
		deficit := wlmSiQueueContainerR1Max(0, before+admitted-budget)
		cumulativeAdmitted += admitted
		cumulativeRejected += rejected
		cumulativeServed += served
		row := wlmSiQueueContainerR1Row{
			Seed: seed, ScheduleFamily: schedule, InitialOccupancy: initial, ServiceBudget: budget,
			Substrate: substrate, Step: step, Arrivals: arriving, AdmittedCount: admitted,
			RejectedCount: rejected, ServedCount: served, QueueBefore: before, QueueAfter: after,
			Deficit: deficit, CumulativeAdmitted: cumulativeAdmitted, CumulativeRejected: cumulativeRejected,
			CumulativeServed: cumulativeServed, QueuedTokenIDs: queue.Snapshot(), ServicedTokenIDs: serviced,
		}
		conservationError := wlmSiQueueContainerR1Abs(after - (before + admitted - served))
		deficitError := wlmSiQueueContainerR1Abs(after - deficit)
		wlmSiQueueContainerR1MaxMetric(metrics, "max_conservation_error", conservationError)
		wlmSiQueueContainerR1MaxMetric(metrics, "max_deficit_law_error", deficitError)
		if served > budget {
			metrics["budget_overrun_rows"]++
		}
		if conservationError != 0 || deficitError != 0 || len(row.QueuedTokenIDs) != after || after < 0 || after > 16 {
			caseBad = true
		}
		rows = append(rows, row)
	}
	return rows, caseBad
}

func wlmSiQueueContainerR1EqualRows(a, b []wlmSiQueueContainerR1Row) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		x, y := a[i], b[i]
		if x.Seed != y.Seed || x.ScheduleFamily != y.ScheduleFamily || x.InitialOccupancy != y.InitialOccupancy || x.ServiceBudget != y.ServiceBudget || x.Step != y.Step || x.Arrivals != y.Arrivals || x.AdmittedCount != y.AdmittedCount || x.RejectedCount != y.RejectedCount || x.ServedCount != y.ServedCount || x.QueueBefore != y.QueueBefore || x.QueueAfter != y.QueueAfter || x.Deficit != y.Deficit || x.CumulativeAdmitted != y.CumulativeAdmitted || x.CumulativeRejected != y.CumulativeRejected || x.CumulativeServed != y.CumulativeServed || !wlmSiQueueContainerR1IntsEqual(x.QueuedTokenIDs, y.QueuedTokenIDs) || !wlmSiQueueContainerR1IntsEqual(x.ServicedTokenIDs, y.ServicedTokenIDs) {
			return false
		}
	}
	return true
}

func wlmSiQueueContainerR1EqualCardinalities(a, b []wlmSiQueueContainerR1Row) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].QueueBefore != b[i].QueueBefore || a[i].QueueAfter != b[i].QueueAfter || len(a[i].QueuedTokenIDs) != a[i].QueueAfter || len(b[i].QueuedTokenIDs) != b[i].QueueAfter {
			return false
		}
	}
	return true
}

func wlmSiQueueContainerR1Copy(in []int) []int { out := make([]int, len(in)); copy(out, in); return out }
func wlmSiQueueContainerR1IntsEqual(a, b []int) bool { if len(a) != len(b) { return false }; for i := range a { if a[i] != b[i] { return false } }; return true }
func wlmSiQueueContainerR1Min(a, b int) int { if a < b { return a }; return b }
func wlmSiQueueContainerR1Max(a, b int) int { if a > b { return a }; return b }
func wlmSiQueueContainerR1Abs(v int) int { if v < 0 { return -v }; return v }
func wlmSiQueueContainerR1MaxMetric(metrics map[string]float64, name string, value int) { if float64(value) > metrics[name] { metrics[name] = float64(value) } }
