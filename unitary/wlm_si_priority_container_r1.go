package unitary

type wlmSiPriorityContainerR1Item struct {
	Priority int
	Token    int
}

type wlmSiPriorityContainerR1Row struct {
	Seed                 uint32
	ScheduleFamily       string
	InitialCardinality   int
	ExtractionBudget     int
	Substrate            string
	Step                 int
	RequestedPushes      int
	AdmittedCount        int
	RejectedCount        int
	ExtractedCount       int
	QueueBefore          int
	QueueAfter           int
	Residual             int
	CumulativeAdmitted   int
	CumulativeRejected   int
	CumulativeExtracted  int
	QueueItems           []wlmSiPriorityContainerR1Item
	ExtractedItems       []wlmSiPriorityContainerR1Item
}

type wlmSiPriorityContainerR1Result struct {
	Rows    []wlmSiPriorityContainerR1Row
	Metrics map[string]float64
}

type wlmSiPriorityContainerR1Op struct {
	Requested int
	Items     []wlmSiPriorityContainerR1Item
}

type wlmSiPriorityContainerR1Queue interface {
	Len() int
	Push(wlmSiPriorityContainerR1Item)
	Pop() wlmSiPriorityContainerR1Item
	Snapshot() []wlmSiPriorityContainerR1Item
}

type wlmSiPriorityContainerR1Sorted struct {
	items []wlmSiPriorityContainerR1Item
}

func (q *wlmSiPriorityContainerR1Sorted) Len() int { return len(q.items) }
func (q *wlmSiPriorityContainerR1Sorted) Push(v wlmSiPriorityContainerR1Item) {
	i := 0
	for i < len(q.items) && !wlmSiPriorityContainerR1Less(v, q.items[i]) {
		i++
	}
	q.items = append(q.items, wlmSiPriorityContainerR1Item{})
	copy(q.items[i+1:], q.items[i:])
	q.items[i] = v
}
func (q *wlmSiPriorityContainerR1Sorted) Pop() wlmSiPriorityContainerR1Item {
	v := q.items[0]
	q.items = q.items[1:]
	return v
}
func (q *wlmSiPriorityContainerR1Sorted) Snapshot() []wlmSiPriorityContainerR1Item {
	out := make([]wlmSiPriorityContainerR1Item, len(q.items))
	copy(out, q.items)
	return out
}

type wlmSiPriorityContainerR1Heap struct {
	items []wlmSiPriorityContainerR1Item
}

func (q *wlmSiPriorityContainerR1Heap) Len() int { return len(q.items) }
func (q *wlmSiPriorityContainerR1Heap) Push(v wlmSiPriorityContainerR1Item) {
	q.items = append(q.items, v)
	i := len(q.items) - 1
	for i > 0 {
		parent := (i - 1) / 2
		if !wlmSiPriorityContainerR1Less(q.items[i], q.items[parent]) {
			break
		}
		q.items[i], q.items[parent] = q.items[parent], q.items[i]
		i = parent
	}
}
func (q *wlmSiPriorityContainerR1Heap) Pop() wlmSiPriorityContainerR1Item {
	root := q.items[0]
	last := len(q.items) - 1
	q.items[0] = q.items[last]
	q.items = q.items[:last]
	i := 0
	for {
		left := i*2 + 1
		if left >= len(q.items) {
			break
		}
		right := left + 1
		child := left
		if right < len(q.items) && wlmSiPriorityContainerR1Less(q.items[right], q.items[left]) {
			child = right
		}
		if !wlmSiPriorityContainerR1Less(q.items[child], q.items[i]) {
			break
		}
		q.items[i], q.items[child] = q.items[child], q.items[i]
		i = child
	}
	return root
}
func (q *wlmSiPriorityContainerR1Heap) Snapshot() []wlmSiPriorityContainerR1Item {
	copyHeap := &wlmSiPriorityContainerR1Heap{items: append([]wlmSiPriorityContainerR1Item(nil), q.items...)}
	out := make([]wlmSiPriorityContainerR1Item, 0, len(q.items))
	for copyHeap.Len() > 0 {
		out = append(out, copyHeap.Pop())
	}
	return out
}

func wlmSiPriorityContainerR1Less(a, b wlmSiPriorityContainerR1Item) bool {
	if a.Priority != b.Priority {
		return a.Priority < b.Priority
	}
	return a.Token < b.Token
}

// RunWlmSiPriorityContainerR1 evaluates the frozen paired priority-container matrix.
func RunWlmSiPriorityContainerR1() interface{} {
	seeds := [...]uint32{701, 719, 733, 751, 773, 797, 821, 853}
	schedules := [...]string{"monotone", "burst", "alternating", "lcg32"}
	initials := [...]int{0, 5, 10, 15}
	budgets := [...]int{1, 2, 4, 8}
	metrics := map[string]float64{
		"paired_case_count":                 0,
		"completed_substrate_case_runs":     0,
		"total_emitted_rows":                0,
		"paired_observable_mismatch_cases":  0,
		"row_count_mismatch_cases":          0,
		"state_cardinality_mismatch_cases":  0,
		"invalid_operation_schedules":       0,
		"budget_overrun_rows":               0,
		"law_violation_cases":               0,
		"max_conservation_error":            0,
		"max_residual_law_error":            0,
	}
	rows := make([]wlmSiPriorityContainerR1Row, 0, 65536)

	for seedIndex, seed := range seeds {
		for _, schedule := range schedules {
			ops, valid := wlmSiPriorityContainerR1Operations(seed, schedule)
			if !valid {
				metrics["invalid_operation_schedules"]++
				continue
			}
			for _, initial := range initials {
				for _, budget := range budgets {
					metrics["paired_case_count"]++
					var sortedRows, heapRows []wlmSiPriorityContainerR1Row
					var sortedBad, heapBad bool
					if seedIndex%2 == 0 {
						sortedRows, sortedBad = wlmSiPriorityContainerR1Run(seed, schedule, initial, budget, ops, "sorted-slice-priority", true, metrics)
						heapRows, heapBad = wlmSiPriorityContainerR1Run(seed, schedule, initial, budget, ops, "binary-heap-priority", false, metrics)
					} else {
						heapRows, heapBad = wlmSiPriorityContainerR1Run(seed, schedule, initial, budget, ops, "binary-heap-priority", false, metrics)
						sortedRows, sortedBad = wlmSiPriorityContainerR1Run(seed, schedule, initial, budget, ops, "sorted-slice-priority", true, metrics)
					}
					rows = append(rows, sortedRows...)
					rows = append(rows, heapRows...)
					metrics["completed_substrate_case_runs"] += 2
					if sortedBad || heapBad {
						metrics["law_violation_cases"]++
					}
					if len(sortedRows) != len(heapRows) {
						metrics["row_count_mismatch_cases"]++
					}
					if !wlmSiPriorityContainerR1EqualRows(sortedRows, heapRows) {
						metrics["paired_observable_mismatch_cases"]++
					}
					if !wlmSiPriorityContainerR1EqualCardinalities(sortedRows, heapRows) {
						metrics["state_cardinality_mismatch_cases"]++
					}
				}
			}
		}
	}
	metrics["total_emitted_rows"] = float64(len(rows))
	return wlmSiPriorityContainerR1Result{Rows: rows, Metrics: metrics}
}

func wlmSiPriorityContainerR1Operations(seed uint32, family string) ([]wlmSiPriorityContainerR1Op, bool) {
	ops := make([]wlmSiPriorityContainerR1Op, 64)
	state := seed
	for step := range ops {
		requested := 0
		switch family {
		case "monotone":
			requested = 1 + int(seed%4)
		case "burst":
			if (uint32(step)+seed)%4 < 2 {
				requested = 16
			}
		case "alternating":
			if (uint32(step)+seed)%2 == 0 {
				requested = 8
			} else {
				requested = 1
			}
		case "lcg32":
			state = state*1664525 + 1013904223
			requested = int(state % 17)
		default:
			return ops, false
		}
		if requested < 0 || requested > 16 {
			return ops, false
		}
		items := make([]wlmSiPriorityContainerR1Item, requested)
		for ordinal := 0; ordinal < requested; ordinal++ {
			token := step*32 + ordinal
			var priority int
			if family == "monotone" {
				priority = token
			} else if family == "lcg32" {
				state = state*1664525 + 1013904223
				priority = int(state % 257)
			} else {
				priority = int((seed*131 + uint32(step*17+ordinal*29)) % 257)
			}
			items[ordinal] = wlmSiPriorityContainerR1Item{Priority: priority, Token: token}
		}
		ops[step] = wlmSiPriorityContainerR1Op{Requested: requested, Items: items}
	}
	return ops, true
}

func wlmSiPriorityContainerR1Run(seed uint32, schedule string, initial, budget int, ops []wlmSiPriorityContainerR1Op, substrate string, sorted bool, metrics map[string]float64) ([]wlmSiPriorityContainerR1Row, bool) {
	var queue wlmSiPriorityContainerR1Queue
	if sorted {
		queue = &wlmSiPriorityContainerR1Sorted{}
	} else {
		queue = &wlmSiPriorityContainerR1Heap{}
	}
	for i := 0; i < initial; i++ {
		queue.Push(wlmSiPriorityContainerR1Item{
			Priority: int((seed*73 + uint32(i*19)) % 257),
			Token:    -initial + i,
		})
	}

	rows := make([]wlmSiPriorityContainerR1Row, 0, len(ops))
	cumulativeAdmitted, cumulativeRejected, cumulativeExtracted := 0, 0, 0
	caseBad := false

	for step, op := range ops {
		before := queue.Len()
		admitted := wlmSiPriorityContainerR1Min(op.Requested, 16-before)
		rejected := op.Requested - admitted
		for i := 0; i < admitted; i++ {
			queue.Push(op.Items[i])
		}

		extractedCount := wlmSiPriorityContainerR1Min(queue.Len(), budget)
		extracted := make([]wlmSiPriorityContainerR1Item, extractedCount)
		for i := range extracted {
			extracted[i] = queue.Pop()
		}

		after := queue.Len()
		residual := wlmSiPriorityContainerR1Max(0, before+admitted-budget)
		cumulativeAdmitted += admitted
		cumulativeRejected += rejected
		cumulativeExtracted += extractedCount

		row := wlmSiPriorityContainerR1Row{
			Seed: seed, ScheduleFamily: schedule, InitialCardinality: initial,
			ExtractionBudget: budget, Substrate: substrate, Step: step,
			RequestedPushes: op.Requested, AdmittedCount: admitted, RejectedCount: rejected,
			ExtractedCount: extractedCount, QueueBefore: before, QueueAfter: after,
			Residual: residual, CumulativeAdmitted: cumulativeAdmitted,
			CumulativeRejected: cumulativeRejected, CumulativeExtracted: cumulativeExtracted,
			QueueItems: queue.Snapshot(), ExtractedItems: extracted,
		}

		conservationError := wlmSiPriorityContainerR1Abs(after - (before + admitted - extractedCount))
		residualError := wlmSiPriorityContainerR1Abs(after - residual)
		wlmSiPriorityContainerR1MaxMetric(metrics, "max_conservation_error", conservationError)
		wlmSiPriorityContainerR1MaxMetric(metrics, "max_residual_law_error", residualError)
		if extractedCount > budget {
			metrics["budget_overrun_rows"]++
		}
		if conservationError != 0 || residualError != 0 || len(row.QueueItems) != after || after < 0 || after > 16 {
			caseBad = true
		}
		rows = append(rows, row)
	}
	return rows, caseBad
}

func wlmSiPriorityContainerR1EqualRows(a, b []wlmSiPriorityContainerR1Row) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		x, y := a[i], b[i]
		if x.Seed != y.Seed || x.ScheduleFamily != y.ScheduleFamily ||
			x.InitialCardinality != y.InitialCardinality || x.ExtractionBudget != y.ExtractionBudget ||
			x.Step != y.Step || x.RequestedPushes != y.RequestedPushes ||
			x.AdmittedCount != y.AdmittedCount || x.RejectedCount != y.RejectedCount ||
			x.ExtractedCount != y.ExtractedCount || x.QueueBefore != y.QueueBefore ||
			x.QueueAfter != y.QueueAfter || x.Residual != y.Residual ||
			x.CumulativeAdmitted != y.CumulativeAdmitted ||
			x.CumulativeRejected != y.CumulativeRejected ||
			x.CumulativeExtracted != y.CumulativeExtracted ||
			!wlmSiPriorityContainerR1ItemsEqual(x.QueueItems, y.QueueItems) ||
			!wlmSiPriorityContainerR1ItemsEqual(x.ExtractedItems, y.ExtractedItems) {
			return false
		}
	}
	return true
}

func wlmSiPriorityContainerR1EqualCardinalities(a, b []wlmSiPriorityContainerR1Row) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].QueueBefore != b[i].QueueBefore || a[i].QueueAfter != b[i].QueueAfter ||
			len(a[i].QueueItems) != a[i].QueueAfter || len(b[i].QueueItems) != b[i].QueueAfter {
			return false
		}
	}
	return true
}

func wlmSiPriorityContainerR1ItemsEqual(a, b []wlmSiPriorityContainerR1Item) bool {
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
func wlmSiPriorityContainerR1Min(a, b int) int { if a < b { return a }; return b }
func wlmSiPriorityContainerR1Max(a, b int) int { if a > b { return a }; return b }
func wlmSiPriorityContainerR1Abs(v int) int { if v < 0 { return -v }; return v }
func wlmSiPriorityContainerR1MaxMetric(metrics map[string]float64, name string, value int) {
	if float64(value) > metrics[name] { metrics[name] = float64(value) }
}
