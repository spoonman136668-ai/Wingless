package unitary

type wlmPriorityItem struct {
	Priority int
	Token    int
}

type wlmPriorityRow struct {
	Seed                uint32
	ScheduleFamily      string
	InitialCardinality  int
	ExtractionBudget    int
	Substrate           string
	Step                int
	RequestedInsertions int
	AdmittedCount       int
	RejectedCount       int
	ExtractedCount      int
	CardinalityBefore   int
	CardinalityAfter    int
	Residual            int
	CumulativeAdmitted  int
	CumulativeRejected  int
	CumulativeExtracted int
	OrderedItems        []wlmPriorityItem
	ExtractedItems      []wlmPriorityItem
}

type wlmPriorityResult struct {
	Rows    []wlmPriorityRow
	Metrics map[string]float64
}

type wlmPriorityOp struct {
	Requested  int
	Priorities [16]int
}

type wlmPriorityContainer interface {
	Len() int
	Push(wlmPriorityItem)
	PopMin() wlmPriorityItem
	Snapshot() []wlmPriorityItem
}

type wlmSortedPriority struct{ v []wlmPriorityItem }

func (q *wlmSortedPriority) Len() int { return len(q.v) }
func (q *wlmSortedPriority) Push(x wlmPriorityItem) {
	i := 0
	for i < len(q.v) && !wlmPriorityLess(x, q.v[i]) {
		i++
	}
	q.v = append(q.v, wlmPriorityItem{})
	copy(q.v[i+1:], q.v[i:])
	q.v[i] = x
}
func (q *wlmSortedPriority) PopMin() wlmPriorityItem {
	x := q.v[0]
	q.v = q.v[1:]
	return x
}
func (q *wlmSortedPriority) Snapshot() []wlmPriorityItem {
	out := make([]wlmPriorityItem, len(q.v))
	copy(out, q.v)
	return out
}

type wlmHeapPriority struct{ v []wlmPriorityItem }

func (q *wlmHeapPriority) Len() int { return len(q.v) }
func (q *wlmHeapPriority) Push(x wlmPriorityItem) {
	q.v = append(q.v, x)
	i := len(q.v) - 1
	for i > 0 {
		p := (i - 1) / 2
		if !wlmPriorityLess(q.v[i], q.v[p]) {
			break
		}
		q.v[i], q.v[p] = q.v[p], q.v[i]
		i = p
	}
}
func (q *wlmHeapPriority) PopMin() wlmPriorityItem {
	x := q.v[0]
	last := q.v[len(q.v)-1]
	q.v = q.v[:len(q.v)-1]
	if len(q.v) > 0 {
		q.v[0] = last
		for i := 0; ; {
			l := 2*i + 1
			r := l + 1
			if l >= len(q.v) {
				break
			}
			m := l
			if r < len(q.v) && wlmPriorityLess(q.v[r], q.v[l]) {
				m = r
			}
			if !wlmPriorityLess(q.v[m], q.v[i]) {
				break
			}
			q.v[i], q.v[m] = q.v[m], q.v[i]
			i = m
		}
	}
	return x
}
func (q *wlmHeapPriority) Snapshot() []wlmPriorityItem {
	tmp := &wlmHeapPriority{v: append([]wlmPriorityItem(nil), q.v...)}
	out := make([]wlmPriorityItem, 0, len(tmp.v))
	for tmp.Len() > 0 {
		out = append(out, tmp.PopMin())
	}
	return out
}

func wlmPriorityLess(a, b wlmPriorityItem) bool {
	if a.Priority != b.Priority {
		return a.Priority < b.Priority
	}
	return a.Token < b.Token
}

func RunWlmSiPriorityContainerR1() interface{} {
	seeds := [...]uint32{701, 719, 733, 751, 773, 797, 821, 853}
	schedules := [...]string{"monotone", "burst", "alternating", "lcg32"}
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
	rows := make([]wlmPriorityRow, 0, 65536)

	for si, seed := range seeds {
		for _, schedule := range schedules {
			ops, ok := wlmPriorityOps(seed, schedule)
			if !ok {
				metrics["invalid_operation_schedules"]++
				continue
			}
			for _, initial := range initials {
				for _, budget := range budgets {
					metrics["paired_case_count"]++
					var a, b []wlmPriorityRow
					var badA, badB bool
					if si%2 == 0 {
						a, badA = wlmPriorityRun(seed, schedule, initial, budget, ops, "sorted-slice-priority", true, metrics)
						b, badB = wlmPriorityRun(seed, schedule, initial, budget, ops, "binary-heap-priority", false, metrics)
					} else {
						b, badB = wlmPriorityRun(seed, schedule, initial, budget, ops, "binary-heap-priority", false, metrics)
						a, badA = wlmPriorityRun(seed, schedule, initial, budget, ops, "sorted-slice-priority", true, metrics)
					}
					rows = append(rows, a...)
					rows = append(rows, b...)
					metrics["completed_substrate_case_runs"] += 2
					if badA || badB {
						metrics["law_violation_cases"]++
					}
					if len(a) != len(b) {
						metrics["row_count_mismatch_cases"]++
					}
					if !wlmPriorityRowsEqual(a, b) {
						metrics["paired_observable_mismatch_cases"]++
					}
					if !wlmPriorityCardinalitiesEqual(a, b) {
						metrics["state_cardinality_mismatch_cases"]++
					}
				}
			}
		}
	}
	metrics["total_emitted_rows"] = float64(len(rows))
	return wlmPriorityResult{Rows: rows, Metrics: metrics}
}

func wlmPriorityOps(seed uint32, family string) ([]wlmPriorityOp, bool) {
	ops := make([]wlmPriorityOp, 64)
	state := seed
	for step := range ops {
		var op wlmPriorityOp
		switch family {
		case "monotone":
			op.Requested = 1 + int(seed%4)
			for i := 0; i < op.Requested; i++ {
				op.Priorities[i] = step*16 + i
			}
		case "burst":
			if (uint32(step)+seed)%4 < 2 {
				op.Requested = 16
			}
			for i := 0; i < op.Requested; i++ {
				op.Priorities[i] = 4096 - step*16 - i
			}
		case "alternating":
			op.Requested = int((uint32(step)+seed)%9)
			for i := 0; i < op.Requested; i++ {
				if step%2 == 0 {
					op.Priorities[i] = i
				} else {
					op.Priorities[i] = 31 - i
				}
			}
		case "lcg32":
			state = state*1664525 + 1013904223
			op.Requested = int(state % 17)
			for i := 0; i < op.Requested; i++ {
				state = state*1664525 + 1013904223
				op.Priorities[i] = int(state % 257)
			}
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

func wlmPriorityRun(seed uint32, schedule string, initial, budget int, ops []wlmPriorityOp, substrate string, sorted bool, metrics map[string]float64) ([]wlmPriorityRow, bool) {
	var q wlmPriorityContainer
	if sorted {
		q = &wlmSortedPriority{}
	} else {
		q = &wlmHeapPriority{}
	}
	for i := 0; i < initial; i++ {
		q.Push(wlmPriorityItem{Priority: initial - i, Token: -initial + i})
	}
	rows := make([]wlmPriorityRow, 0, 64)
	ca, cr, ce := 0, 0, 0
	bad := false
	for step, op := range ops {
		before := q.Len()
		admitted := wlmPriorityMin(op.Requested, 16-before)
		rejected := op.Requested - admitted
		for i := 0; i < admitted; i++ {
			q.Push(wlmPriorityItem{Priority: op.Priorities[i], Token: step*16 + i})
		}
		extractedCount := wlmPriorityMin(q.Len(), budget)
		extracted := make([]wlmPriorityItem, extractedCount)
		for i := range extracted {
			extracted[i] = q.PopMin()
		}
		after := q.Len()
		residual := wlmPriorityMax(0, before+admitted-budget)
		ca += admitted
		cr += rejected
		ce += extractedCount
		row := wlmPriorityRow{
			Seed: seed, ScheduleFamily: schedule, InitialCardinality: initial, ExtractionBudget: budget,
			Substrate: substrate, Step: step, RequestedInsertions: op.Requested, AdmittedCount: admitted,
			RejectedCount: rejected, ExtractedCount: extractedCount, CardinalityBefore: before,
			CardinalityAfter: after, Residual: residual, CumulativeAdmitted: ca, CumulativeRejected: cr,
			CumulativeExtracted: ce, OrderedItems: q.Snapshot(), ExtractedItems: extracted,
		}
		conservation := wlmPriorityAbs(after - (before + admitted - extractedCount))
		residualErr := wlmPriorityAbs(after - residual)
		wlmPriorityMaxMetric(metrics, "max_conservation_error", conservation)
		wlmPriorityMaxMetric(metrics, "max_residual_law_error", residualErr)
		if extractedCount > budget {
			metrics["budget_overrun_rows"]++
		}
		if conservation != 0 || residualErr != 0 || len(row.OrderedItems) != after || after < 0 || after > 16 {
			bad = true
		}
		rows = append(rows, row)
	}
	return rows, bad
}

func wlmPriorityRowsEqual(a, b []wlmPriorityRow) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		x, y := a[i], b[i]
		if x.Seed != y.Seed || x.ScheduleFamily != y.ScheduleFamily || x.InitialCardinality != y.InitialCardinality ||
			x.ExtractionBudget != y.ExtractionBudget || x.Step != y.Step || x.RequestedInsertions != y.RequestedInsertions ||
			x.AdmittedCount != y.AdmittedCount || x.RejectedCount != y.RejectedCount || x.ExtractedCount != y.ExtractedCount ||
			x.CardinalityBefore != y.CardinalityBefore || x.CardinalityAfter != y.CardinalityAfter || x.Residual != y.Residual ||
			x.CumulativeAdmitted != y.CumulativeAdmitted || x.CumulativeRejected != y.CumulativeRejected ||
			x.CumulativeExtracted != y.CumulativeExtracted || !wlmPriorityItemsEqual(x.OrderedItems, y.OrderedItems) ||
			!wlmPriorityItemsEqual(x.ExtractedItems, y.ExtractedItems) {
			return false
		}
	}
	return true
}

func wlmPriorityCardinalitiesEqual(a, b []wlmPriorityRow) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].CardinalityBefore != b[i].CardinalityBefore || a[i].CardinalityAfter != b[i].CardinalityAfter ||
			len(a[i].OrderedItems) != a[i].CardinalityAfter || len(b[i].OrderedItems) != b[i].CardinalityAfter {
			return false
		}
	}
	return true
}

func wlmPriorityItemsEqual(a, b []wlmPriorityItem) bool {
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
func wlmPriorityMin(a, b int) int { if a < b { return a }; return b }
func wlmPriorityMax(a, b int) int { if a > b { return a }; return b }
func wlmPriorityAbs(v int) int { if v < 0 { return -v }; return v }
func wlmPriorityMaxMetric(m map[string]float64, k string, v int) { if float64(v) > m[k] { m[k] = float64(v) } }
