package unitary

type wlmSparseEntry struct { Coord int; Value int }
type wlmSparseRead struct { Coord int; Value int }
type wlmSparseRow struct {
	Seed uint32
	ScheduleFamily string
	InitialNonzeroCount int
	OperationBudget int
	Substrate string
	Step int
	ExecutedOperations int
	AcceptedSetCount int
	RejectedSetCount int
	ReadCount int
	NonzeroBefore int
	NonzeroAfter int
	SumBefore int
	SumAfter int
	CumulativeAccepted int
	CumulativeRejected int
	Entries []wlmSparseEntry
	Reads []wlmSparseRead
}
type wlmSparseResult struct { Rows []wlmSparseRow; Metrics map[string]float64 }
type wlmSparseOp struct { Kind byte; Coord int; Value int }
type wlmSparseVector interface {
	Nonzero() int
	Get(int) int
	Set(int,int) bool
	Snapshot() []wlmSparseEntry
}

type wlmSparseDense struct { values [64]int; nonzero int }
func(v *wlmSparseDense)Nonzero()int{return v.nonzero}
func(v *wlmSparseDense)Get(i int)int{return v.values[i]}
func(v *wlmSparseDense)Set(i,x int)bool{
	old:=v.values[i]
	if old==0&&x!=0&&v.nonzero>=16{return false}
	if old==0&&x!=0{v.nonzero++}
	if old!=0&&x==0{v.nonzero--}
	v.values[i]=x
	return true
}
func(v *wlmSparseDense)Snapshot()[]wlmSparseEntry{out:=make([]wlmSparseEntry,0,v.nonzero);for i,x:=range v.values{if x!=0{out=append(out,wlmSparseEntry{Coord:i,Value:x})}};return out}

type wlmSparseList struct { entries []wlmSparseEntry }
func(v *wlmSparseList)Nonzero()int{return len(v.entries)}
func(v *wlmSparseList)find(i int)(int,bool){lo,hi:=0,len(v.entries);for lo<hi{m:=(lo+hi)/2;if v.entries[m].Coord<i{lo=m+1}else{hi=m}};return lo,lo<len(v.entries)&&v.entries[lo].Coord==i}
func(v *wlmSparseList)Get(i int)int{p,ok:=v.find(i);if !ok{return 0};return v.entries[p].Value}
func(v *wlmSparseList)Set(i,x int)bool{
	p,ok:=v.find(i)
	if ok{
		if x==0{copy(v.entries[p:],v.entries[p+1:]);v.entries=v.entries[:len(v.entries)-1]}else{v.entries[p].Value=x}
		return true
	}
	if x==0{return true}
	if len(v.entries)>=16{return false}
	v.entries=append(v.entries,wlmSparseEntry{})
	copy(v.entries[p+1:],v.entries[p:])
	v.entries[p]=wlmSparseEntry{Coord:i,Value:x}
	return true
}
func(v *wlmSparseList)Snapshot()[]wlmSparseEntry{out:=make([]wlmSparseEntry,len(v.entries));copy(out,v.entries);return out}

func RunWlmSiSparseVectorR1()interface{}{
	seeds:=[...]uint32{2531,2557,2591,2621,2657,2683,2711,2741}
	families:=[...]string{"sequential","hotset","zero-wave","lcg32"}
	initials:=[...]int{0,4,8,12}
	budgets:=[...]int{1,2,4,8}
	metrics:=map[string]float64{"paired_case_count":0,"completed_substrate_case_runs":0,"total_emitted_rows":0,"paired_observable_mismatch_cases":0,"row_count_mismatch_cases":0,"state_cardinality_mismatch_cases":0,"invalid_operation_schedules":0,"budget_overrun_rows":0,"law_violation_cases":0,"max_conservation_error":0,"max_residual_law_error":0}
	rows:=make([]wlmSparseRow,0,65536)
	for si,seed:=range seeds{for _,family:=range families{for _,initial:=range initials{for _,budget:=range budgets{
		metrics["paired_case_count"]++
		var a,b []wlmSparseRow;var badA,badB bool
		if si%2==0{a,badA=wlmSparseRun(seed,family,initial,budget,"dense-array",true,metrics);b,badB=wlmSparseRun(seed,family,initial,budget,"sorted-sparse-list",false,metrics)}else{b,badB=wlmSparseRun(seed,family,initial,budget,"sorted-sparse-list",false,metrics);a,badA=wlmSparseRun(seed,family,initial,budget,"dense-array",true,metrics)}
		rows=append(rows,a...);rows=append(rows,b...);metrics["completed_substrate_case_runs"]+=2
		if badA||badB{metrics["law_violation_cases"]++};if len(a)!=len(b){metrics["row_count_mismatch_cases"]++};if !wlmSparseRowsEqual(a,b){metrics["paired_observable_mismatch_cases"]++};if !wlmSparseCardinalitiesEqual(a,b){metrics["state_cardinality_mismatch_cases"]++}
	}}}}
	metrics["total_emitted_rows"]=float64(len(rows));return wlmSparseResult{Rows:rows,Metrics:metrics}
}

func wlmSparseRun(seed uint32,family string,initial,budget int,substrate string,dense bool,metrics map[string]float64)([]wlmSparseRow,bool){
	var v wlmSparseVector;if dense{v=&wlmSparseDense{}}else{v=&wlmSparseList{}}
	for i:=0;i<initial;i++{coord:=(i*7+int(seed))&63;val:=1+((i+int(seed))%8);if !v.Set(coord,val){panic("initial sparse admission failed")}}
	rows:=make([]wlmSparseRow,0,64);ca,cr:=0,0;state:=seed;bad:=false
	for step:=0;step<64;step++{
		ops,ok,next:=wlmSparseOps(seed,state,family,step,budget);state=next;if !ok{metrics["invalid_operation_schedules"]++;return rows,true};if len(ops)>budget{metrics["budget_overrun_rows"]++;bad=true}
		before:=v.Nonzero();sumBefore:=wlmSparseSum(v.Snapshot());accepted,rejected,reads:=0,0,0;readRows:=make([]wlmSparseRead,0,budget)
		for _,op:=range ops{switch op.Kind{case 'S':if v.Set(op.Coord,op.Value){accepted++}else{rejected++};case 'G':reads++;readRows=append(readRows,wlmSparseRead{Coord:op.Coord,Value:v.Get(op.Coord)});default:metrics["invalid_operation_schedules"]++;bad=true}}
		after:=v.Nonzero();snapshot:=v.Snapshot();sumAfter:=wlmSparseSum(snapshot)
		residual:=wlmSparseAbs(len(snapshot)-after)
		conservation:=0
		for _,e:=range snapshot{if e.Value==0{conservation++}}
		wlmSparseMaxMetric(metrics,"max_conservation_error",conservation);wlmSparseMaxMetric(metrics,"max_residual_law_error",residual)
		if residual!=0||conservation!=0||after<0||after>16{bad=true}
		ca+=accepted;cr+=rejected
		rows=append(rows,wlmSparseRow{Seed:seed,ScheduleFamily:family,InitialNonzeroCount:initial,OperationBudget:budget,Substrate:substrate,Step:step,ExecutedOperations:len(ops),AcceptedSetCount:accepted,RejectedSetCount:rejected,ReadCount:reads,NonzeroBefore:before,NonzeroAfter:after,SumBefore:sumBefore,SumAfter:sumAfter,CumulativeAccepted:ca,CumulativeRejected:cr,Entries:snapshot,Reads:readRows})
	}
	return rows,bad
}
func wlmSparseOps(seed,state uint32,family string,step,budget int)([]wlmSparseOp,bool,uint32){
	ops:=make([]wlmSparseOp,0,budget)
	for ordinal:=0;ordinal<budget;ordinal++{var kind byte;var coord,val int
		switch family{
		case "sequential":
			coord=(step*budget+ordinal)&63;if (step+ordinal)%3==2{kind='G'}else{kind='S';val=1+((step+ordinal+int(seed))%8)}
		case "hotset":
			coord=(int(seed)+step+ordinal)&7;if (step+ordinal)%4==3{kind='G'}else{kind='S';val=-8+((step*3+ordinal+int(seed))%17);if val==0{val=1}}
		case "zero-wave":
			coord=(step+ordinal*5)&63;if step%2==0{kind='S';val=1+((step+ordinal)%8)}else{kind='S';val=0}
		case "lcg32":
			state=state*1664525+1013904223;coord=int(state&63);state=state*1664525+1013904223
			if state%4==0{kind='G'}else{kind='S';val=int((state>>8)%17)-8}
		default:return ops,false,state
		}
		if coord<0||coord>=64||val<-8||val>8{return ops,false,state}
		ops=append(ops,wlmSparseOp{Kind:kind,Coord:coord,Value:val})
	}
	return ops,true,state
}
func wlmSparseSum(entries []wlmSparseEntry)int{s:=0;for _,e:=range entries{s+=e.Value};return s}
func wlmSparseRowsEqual(a,b []wlmSparseRow)bool{if len(a)!=len(b){return false};for i:=range a{x,y:=a[i],b[i];if x.Seed!=y.Seed||x.ScheduleFamily!=y.ScheduleFamily||x.InitialNonzeroCount!=y.InitialNonzeroCount||x.OperationBudget!=y.OperationBudget||x.Step!=y.Step||x.ExecutedOperations!=y.ExecutedOperations||x.AcceptedSetCount!=y.AcceptedSetCount||x.RejectedSetCount!=y.RejectedSetCount||x.ReadCount!=y.ReadCount||x.NonzeroBefore!=y.NonzeroBefore||x.NonzeroAfter!=y.NonzeroAfter||x.SumBefore!=y.SumBefore||x.SumAfter!=y.SumAfter||x.CumulativeAccepted!=y.CumulativeAccepted||x.CumulativeRejected!=y.CumulativeRejected||!wlmSparseEntriesEqual(x.Entries,y.Entries)||!wlmSparseReadsEqual(x.Reads,y.Reads){return false}};return true}
func wlmSparseCardinalitiesEqual(a,b []wlmSparseRow)bool{if len(a)!=len(b){return false};for i:=range a{if a[i].NonzeroBefore!=b[i].NonzeroBefore||a[i].NonzeroAfter!=b[i].NonzeroAfter{return false}};return true}
func wlmSparseEntriesEqual(a,b []wlmSparseEntry)bool{if len(a)!=len(b){return false};for i:=range a{if a[i]!=b[i]{return false}};return true}
func wlmSparseReadsEqual(a,b []wlmSparseRead)bool{if len(a)!=len(b){return false};for i:=range a{if a[i]!=b[i]{return false}};return true}
func wlmSparseAbs(v int)int{if v<0{return -v};return v}
func wlmSparseMaxMetric(m map[string]float64,k string,v int){if float64(v)>m[k]{m[k]=float64(v)}}
