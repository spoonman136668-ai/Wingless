package unitary

type wlmPrefixQuery struct {
	Kind string
	A int
	B int
	Value int
}
type wlmPrefixEntry struct {
	Index int
	Value int
}
type wlmPrefixRow struct {
	Seed uint32
	ScheduleFamily string
	InitialNonzeroCount int
	OperationBudget int
	Substrate string
	Step int
	ExecutedOperations int
	UpdateCount int
	QueryCount int
	NonzeroBefore int
	NonzeroAfter int
	TotalBefore int
	TotalAfter int
	CumulativeUpdateAmount int
	Entries []wlmPrefixEntry
	Queries []wlmPrefixQuery
}
type wlmPrefixResult struct {
	Rows []wlmPrefixRow
	Metrics map[string]float64
}
type wlmPrefixOp struct {
	Kind byte
	A int
	B int
	Delta int
}
type wlmPrefixIndex interface {
	Nonzero() int
	Total() int
	Add(int,int)
	Get(int) int
	Prefix(int) int
	Range(int,int) int
	Snapshot() []wlmPrefixEntry
}

type wlmPrefixDirect struct {
	values [64]int
	nonzero int
	total int
}
func (x *wlmPrefixDirect) Nonzero() int { return x.nonzero }
func (x *wlmPrefixDirect) Total() int { return x.total }
func (x *wlmPrefixDirect) Add(i,delta int) {
	old:=x.values[i]
	next:=old+delta
	if old==0 && next!=0 { x.nonzero++ }
	if old!=0 && next==0 { x.nonzero-- }
	x.values[i]=next
	x.total+=delta
}
func (x *wlmPrefixDirect) Get(i int) int { return x.values[i] }
func (x *wlmPrefixDirect) Prefix(end int) int {
	s:=0
	for i:=0;i<end;i++ { s+=x.values[i] }
	return s
}
func (x *wlmPrefixDirect) Range(a,b int) int {
	return x.Prefix(b)-x.Prefix(a)
}
func (x *wlmPrefixDirect) Snapshot() []wlmPrefixEntry {
	out:=make([]wlmPrefixEntry,0,x.nonzero)
	for i,v:=range x.values {
		if v!=0 { out=append(out,wlmPrefixEntry{Index:i,Value:v}) }
	}
	return out
}

type wlmPrefixFenwick struct {
	values [64]int
	tree [65]int
	nonzero int
	total int
}
func (x *wlmPrefixFenwick) Nonzero() int { return x.nonzero }
func (x *wlmPrefixFenwick) Total() int { return x.total }
func (x *wlmPrefixFenwick) Add(i,delta int) {
	old:=x.values[i]
	next:=old+delta
	if old==0 && next!=0 { x.nonzero++ }
	if old!=0 && next==0 { x.nonzero-- }
	x.values[i]=next
	x.total+=delta
	for p:=i+1;p<=64;p+=p&-p { x.tree[p]+=delta }
}
func (x *wlmPrefixFenwick) Get(i int) int { return x.values[i] }
func (x *wlmPrefixFenwick) Prefix(end int) int {
	s:=0
	for p:=end;p>0;p-=p&-p { s+=x.tree[p] }
	return s
}
func (x *wlmPrefixFenwick) Range(a,b int) int {
	return x.Prefix(b)-x.Prefix(a)
}
func (x *wlmPrefixFenwick) Snapshot() []wlmPrefixEntry {
	out:=make([]wlmPrefixEntry,0,x.nonzero)
	for i,v:=range x.values {
		if v!=0 { out=append(out,wlmPrefixEntry{Index:i,Value:v}) }
	}
	return out
}

func RunWlmSiPrefixIndexR1() interface{} {
	seeds:=[...]uint32{2777,2801,2833,2861,2897,2927,2953,2999}
	families:=[...]string{"sequential","signed-wave","range-hot","lcg32"}
	initials:=[...]int{0,16,32,48}
	budgets:=[...]int{1,2,4,8}
	metrics:=map[string]float64{
		"paired_case_count":0,
		"completed_substrate_case_runs":0,
		"total_emitted_rows":0,
		"paired_observable_mismatch_cases":0,
		"row_count_mismatch_cases":0,
		"state_cardinality_mismatch_cases":0,
		"invalid_operation_schedules":0,
		"budget_overrun_rows":0,
		"law_violation_cases":0,
		"max_conservation_error":0,
		"max_residual_law_error":0,
	}
	rows:=make([]wlmPrefixRow,0,65536)
	for si,seed:=range seeds {
		for _,family:=range families {
			for _,initial:=range initials {
				for _,budget:=range budgets {
					metrics["paired_case_count"]++
					var a,b []wlmPrefixRow
					var badA,badB bool
					if si%2==0 {
						a,badA=wlmPrefixRun(seed,family,initial,budget,"direct-prefix-scan",true,metrics)
						b,badB=wlmPrefixRun(seed,family,initial,budget,"fenwick-tree",false,metrics)
					} else {
						b,badB=wlmPrefixRun(seed,family,initial,budget,"fenwick-tree",false,metrics)
						a,badA=wlmPrefixRun(seed,family,initial,budget,"direct-prefix-scan",true,metrics)
					}
					rows=append(rows,a...)
					rows=append(rows,b...)
					metrics["completed_substrate_case_runs"]+=2
					if badA||badB { metrics["law_violation_cases"]++ }
					if len(a)!=len(b) { metrics["row_count_mismatch_cases"]++ }
					if !wlmPrefixRowsEqual(a,b) { metrics["paired_observable_mismatch_cases"]++ }
					if !wlmPrefixCardinalitiesEqual(a,b) { metrics["state_cardinality_mismatch_cases"]++ }
				}
			}
		}
	}
	metrics["total_emitted_rows"]=float64(len(rows))
	return wlmPrefixResult{Rows:rows,Metrics:metrics}
}

func wlmPrefixRun(seed uint32,family string,initial,budget int,substrate string,direct bool,metrics map[string]float64)([]wlmPrefixRow,bool) {
	var idx wlmPrefixIndex
	if direct { idx=&wlmPrefixDirect{} } else { idx=&wlmPrefixFenwick{} }
	for i:=0;i<initial;i++ {
		slot:=(i*5+int(seed))&63
		delta:=1+((i+int(seed))%7)
		idx.Add(slot,delta)
	}
	rows:=make([]wlmPrefixRow,0,64)
	state:=seed
	cumulative:=0
	bad:=false
	for step:=0;step<64;step++ {
		ops,ok,next:=wlmPrefixOps(seed,state,family,step,budget)
		state=next
		if !ok { metrics["invalid_operation_schedules"]++;return rows,true }
		if len(ops)>budget { metrics["budget_overrun_rows"]++;bad=true }
		beforeN:=idx.Nonzero()
		beforeTotal:=idx.Total()
		updates,queries:=0,0
		queryRows:=make([]wlmPrefixQuery,0,budget)
		for _,op:=range ops {
			switch op.Kind {
			case 'U':
				idx.Add(op.A,op.Delta)
				updates++
				cumulative+=op.Delta
			case 'P':
				queryRows=append(queryRows,wlmPrefixQuery{Kind:"prefix",A:op.A,Value:idx.Prefix(op.A)})
				queries++
			case 'R':
				queryRows=append(queryRows,wlmPrefixQuery{Kind:"range",A:op.A,B:op.B,Value:idx.Range(op.A,op.B)})
				queries++
			default:
				metrics["invalid_operation_schedules"]++
				bad=true
			}
		}
		afterN:=idx.Nonzero()
		afterTotal:=idx.Total()
		snapshot:=idx.Snapshot()
		sum:=0
		for _,e:=range snapshot { sum+=e.Value }
		conservation:=wlmPrefixAbs((afterTotal-beforeTotal)-wlmPrefixAppliedDelta(ops))
		residual:=wlmPrefixAbs(sum-afterTotal)
		wlmPrefixMaxMetric(metrics,"max_conservation_error",conservation)
		wlmPrefixMaxMetric(metrics,"max_residual_law_error",residual)
		if conservation!=0||residual!=0||afterN<0||afterN>64||len(snapshot)!=afterN { bad=true }
		rows=append(rows,wlmPrefixRow{
			Seed:seed,ScheduleFamily:family,InitialNonzeroCount:initial,OperationBudget:budget,
			Substrate:substrate,Step:step,ExecutedOperations:len(ops),UpdateCount:updates,QueryCount:queries,
			NonzeroBefore:beforeN,NonzeroAfter:afterN,TotalBefore:beforeTotal,TotalAfter:afterTotal,
			CumulativeUpdateAmount:cumulative,Entries:snapshot,Queries:queryRows,
		})
	}
	return rows,bad
}

func wlmPrefixOps(seed,state uint32,family string,step,budget int)([]wlmPrefixOp,bool,uint32) {
	ops:=make([]wlmPrefixOp,0,budget)
	for ordinal:=0;ordinal<budget;ordinal++ {
		var op wlmPrefixOp
		switch family {
		case "sequential":
			if (step+ordinal)%3==0 {
				op=wlmPrefixOp{Kind:'P',A:1+((step*budget+ordinal)&63)}
			} else {
				op=wlmPrefixOp{Kind:'U',A:(step*budget+ordinal)&63,Delta:1+((step+ordinal)%5)}
			}
		case "signed-wave":
			slot:=(step+ordinal*3)&63
			if (step+ordinal)%4==3 {
				op=wlmPrefixOp{Kind:'P',A:1+slot}
			} else {
				delta:=1+((step+ordinal)%5)
				if step%2==1 { delta=-delta }
				op=wlmPrefixOp{Kind:'U',A:slot,Delta:delta}
			}
		case "range-hot":
			if (step+ordinal)%2==0 {
				a:=(int(seed)+step+ordinal)&15
				b:=a+1+((step+ordinal)%16)
				if b>64 { b=64 }
				op=wlmPrefixOp{Kind:'R',A:a,B:b}
			} else {
				op=wlmPrefixOp{Kind:'U',A:(int(seed)+step+ordinal)&15,Delta:1+((step+ordinal)%3)}
			}
		case "lcg32":
			state=state*1664525+1013904223
			a:=int(state&63)
			state=state*1664525+1013904223
			switch state%3 {
			case 0:
				op=wlmPrefixOp{Kind:'U',A:a,Delta:int((state>>8)%11)-5}
			case 1:
				op=wlmPrefixOp{Kind:'P',A:a+1}
			default:
				state=state*1664525+1013904223
				b:=int(state&63)+1
				if b<=a { a,b=b-1,a+1 }
				op=wlmPrefixOp{Kind:'R',A:a,B:b}
			}
		default:
			return ops,false,state
		}
		if op.A<0||op.A>64||op.B<0||op.B>64||op.Delta < -8||op.Delta > 8 { return ops,false,state }
		if op.Kind=='U' && op.A>=64 { return ops,false,state }
		if op.Kind=='R' && (op.A>=op.B) { return ops,false,state }
		ops=append(ops,op)
	}
	return ops,true,state
}

func wlmPrefixAppliedDelta(ops []wlmPrefixOp) int {
	s:=0
	for _,op:=range ops { if op.Kind=='U' { s+=op.Delta } }
	return s
}
func wlmPrefixRowsEqual(a,b []wlmPrefixRow) bool {
	if len(a)!=len(b) { return false }
	for i:=range a {
		x,y:=a[i],b[i]
		if x.Seed!=y.Seed||x.ScheduleFamily!=y.ScheduleFamily||x.InitialNonzeroCount!=y.InitialNonzeroCount||
			x.OperationBudget!=y.OperationBudget||x.Step!=y.Step||x.ExecutedOperations!=y.ExecutedOperations||
			x.UpdateCount!=y.UpdateCount||x.QueryCount!=y.QueryCount||x.NonzeroBefore!=y.NonzeroBefore||
			x.NonzeroAfter!=y.NonzeroAfter||x.TotalBefore!=y.TotalBefore||x.TotalAfter!=y.TotalAfter||
			x.CumulativeUpdateAmount!=y.CumulativeUpdateAmount||!wlmPrefixEntriesEqual(x.Entries,y.Entries)||
			!wlmPrefixQueriesEqual(x.Queries,y.Queries) { return false }
	}
	return true
}
func wlmPrefixCardinalitiesEqual(a,b []wlmPrefixRow) bool {
	if len(a)!=len(b) { return false }
	for i:=range a { if a[i].NonzeroBefore!=b[i].NonzeroBefore||a[i].NonzeroAfter!=b[i].NonzeroAfter{return false} }
	return true
}
func wlmPrefixEntriesEqual(a,b []wlmPrefixEntry) bool {
	if len(a)!=len(b) { return false }
	for i:=range a { if a[i]!=b[i] { return false } }
	return true
}
func wlmPrefixQueriesEqual(a,b []wlmPrefixQuery) bool {
	if len(a)!=len(b) { return false }
	for i:=range a { if a[i]!=b[i] { return false } }
	return true
}
func wlmPrefixAbs(v int) int { if v<0{return -v};return v }
func wlmPrefixMaxMetric(m map[string]float64,k string,v int) { if float64(v)>m[k]{m[k]=float64(v)} }
