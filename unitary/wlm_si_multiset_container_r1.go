package unitary

type wlmMultisetEntry struct { Key int; Count int }
type wlmMultisetQuery struct { Key int; Count int }
type wlmMultisetRow struct {
	Seed uint32
	ScheduleFamily string
	InitialTotalMultiplicity int
	OperationBudget int
	Substrate string
	Step int
	ExecutedOperations int
	AddedCount int
	RejectedAddCount int
	RemovedCount int
	RemoveMissingCount int
	QueryCount int
	TotalBefore int
	TotalAfter int
	CumulativeAdded int
	CumulativeRejected int
	CumulativeRemoved int
	Entries []wlmMultisetEntry
	Queries []wlmMultisetQuery
}
type wlmMultisetResult struct { Rows []wlmMultisetRow; Metrics map[string]float64 }
type wlmMultisetOp struct { Kind byte; Key int }
type wlmMultiset interface {
	Total() int
	Add(int) bool
	Remove(int) bool
	Count(int) int
	Snapshot() []wlmMultisetEntry
}

type wlmDenseMultiset struct { counts [32]int; total int }
func(s *wlmDenseMultiset)Total()int{return s.total}
func(s *wlmDenseMultiset)Add(k int)bool{if s.total>=64{return false};s.counts[k]++;s.total++;return true}
func(s *wlmDenseMultiset)Remove(k int)bool{if s.counts[k]==0{return false};s.counts[k]--;s.total--;return true}
func(s *wlmDenseMultiset)Count(k int)int{return s.counts[k]}
func(s *wlmDenseMultiset)Snapshot()[]wlmMultisetEntry{out:=make([]wlmMultisetEntry,0,32);for k:=0;k<32;k++{if s.counts[k]>0{out=append(out,wlmMultisetEntry{Key:k,Count:s.counts[k]})}};return out}

type wlmRunMultiset struct { entries []wlmMultisetEntry; total int }
func(s *wlmRunMultiset)Total()int{return s.total}
func(s *wlmRunMultiset)find(k int)(int,bool){lo,hi:=0,len(s.entries);for lo<hi{m:=(lo+hi)/2;if s.entries[m].Key<k{lo=m+1}else{hi=m}};return lo,lo<len(s.entries)&&s.entries[lo].Key==k}
func(s *wlmRunMultiset)Add(k int)bool{
	if s.total>=64{return false}
	i,ok:=s.find(k)
	if ok{s.entries[i].Count++}else{s.entries=append(s.entries,wlmMultisetEntry{});copy(s.entries[i+1:],s.entries[i:]);s.entries[i]=wlmMultisetEntry{Key:k,Count:1}}
	s.total++;return true
}
func(s *wlmRunMultiset)Remove(k int)bool{
	i,ok:=s.find(k);if !ok{return false}
	s.entries[i].Count--;s.total--
	if s.entries[i].Count==0{copy(s.entries[i:],s.entries[i+1:]);s.entries=s.entries[:len(s.entries)-1]}
	return true
}
func(s *wlmRunMultiset)Count(k int)int{i,ok:=s.find(k);if !ok{return 0};return s.entries[i].Count}
func(s *wlmRunMultiset)Snapshot()[]wlmMultisetEntry{out:=make([]wlmMultisetEntry,len(s.entries));copy(out,s.entries);return out}

func RunWlmSiMultisetContainerR1()interface{}{
	seeds:=[...]uint32{1901,1931,1951,1973,1999,2017,2039,2063}
	families:=[...]string{"sequential","repeat-hot","remove-wave","lcg32"}
	initials:=[...]int{0,16,32,48}
	budgets:=[...]int{1,2,4,8}
	metrics:=map[string]float64{
		"paired_case_count":0,"completed_substrate_case_runs":0,"total_emitted_rows":0,
		"paired_observable_mismatch_cases":0,"row_count_mismatch_cases":0,"state_cardinality_mismatch_cases":0,
		"invalid_operation_schedules":0,"budget_overrun_rows":0,"law_violation_cases":0,
		"max_conservation_error":0,"max_residual_law_error":0,
	}
	rows:=make([]wlmMultisetRow,0,65536)
	for si,seed:=range seeds{for _,family:=range families{for _,initial:=range initials{for _,budget:=range budgets{
		metrics["paired_case_count"]++
		var a,b []wlmMultisetRow;var badA,badB bool
		if si%2==0{
			a,badA=wlmMultisetRun(seed,family,initial,budget,"dense-count-table",true,metrics)
			b,badB=wlmMultisetRun(seed,family,initial,budget,"sorted-run-length",false,metrics)
		}else{
			b,badB=wlmMultisetRun(seed,family,initial,budget,"sorted-run-length",false,metrics)
			a,badA=wlmMultisetRun(seed,family,initial,budget,"dense-count-table",true,metrics)
		}
		rows=append(rows,a...);rows=append(rows,b...);metrics["completed_substrate_case_runs"]+=2
		if badA||badB{metrics["law_violation_cases"]++}
		if len(a)!=len(b){metrics["row_count_mismatch_cases"]++}
		if !wlmMultisetRowsEqual(a,b){metrics["paired_observable_mismatch_cases"]++}
		if !wlmMultisetTotalsEqual(a,b){metrics["state_cardinality_mismatch_cases"]++}
	}}}}
	metrics["total_emitted_rows"]=float64(len(rows))
	return wlmMultisetResult{Rows:rows,Metrics:metrics}
}

func wlmMultisetRun(seed uint32,family string,initial,budget int,substrate string,dense bool,metrics map[string]float64)([]wlmMultisetRow,bool){
	var set wlmMultiset
	if dense{set=&wlmDenseMultiset{}}else{set=&wlmRunMultiset{}}
	for i:=0;i<initial;i++{if !set.Add((i+int(seed))&31){panic("initial multiset admission failed")}}
	rows:=make([]wlmMultisetRow,0,64);ca,cr,crm:=0,0,0;state:=seed;bad:=false
	for step:=0;step<64;step++{
		ops,ok,next:=wlmMultisetOps(seed,state,family,step,budget);state=next
		if !ok{metrics["invalid_operation_schedules"]++;return rows,true}
		if len(ops)>budget{metrics["budget_overrun_rows"]++;bad=true}
		before:=set.Total();added,rejected,removed,missing,qcount:=0,0,0,0,0
		queries:=make([]wlmMultisetQuery,0,budget)
		for _,op:=range ops{switch op.Kind{
		case 'A':if set.Add(op.Key){added++}else{rejected++}
		case 'R':if set.Remove(op.Key){removed++}else{missing++}
		case 'Q':qcount++;queries=append(queries,wlmMultisetQuery{Key:op.Key,Count:set.Count(op.Key)})
		default:metrics["invalid_operation_schedules"]++;bad=true
		}}
		after:=set.Total()
		conservation:=wlmMultisetAbs(after-(before+added-removed))
		snapshot:=set.Snapshot();sum:=0;for _,e:=range snapshot{sum+=e.Count}
		residual:=wlmMultisetAbs(sum-after)
		wlmMultisetMaxMetric(metrics,"max_conservation_error",conservation)
		wlmMultisetMaxMetric(metrics,"max_residual_law_error",residual)
		if conservation!=0||residual!=0||after<0||after>64{bad=true}
		ca+=added;cr+=rejected;crm+=removed
		rows=append(rows,wlmMultisetRow{Seed:seed,ScheduleFamily:family,InitialTotalMultiplicity:initial,OperationBudget:budget,Substrate:substrate,Step:step,ExecutedOperations:len(ops),AddedCount:added,RejectedAddCount:rejected,RemovedCount:removed,RemoveMissingCount:missing,QueryCount:qcount,TotalBefore:before,TotalAfter:after,CumulativeAdded:ca,CumulativeRejected:cr,CumulativeRemoved:crm,Entries:snapshot,Queries:queries})
	}
	return rows,bad
}

func wlmMultisetOps(seed,state uint32,family string,step,budget int)([]wlmMultisetOp,bool,uint32){
	ops:=make([]wlmMultisetOp,0,budget)
	for ordinal:=0;ordinal<budget;ordinal++{var kind byte;var key int
		switch family{
		case "sequential":
			key=(step*budget+ordinal)&31;switch step%3{case 0:kind='A';case 1:kind='Q';default:kind='R'}
		case "repeat-hot":
			key=(int(seed)+step+ordinal)&7;switch(step+ordinal)%4{case 0,1:kind='A';case 2:kind='Q';default:kind='R'}
		case "remove-wave":
			key=(step+ordinal*5)&31;if step%3==0{kind='A'}else{kind='R'}
		case "lcg32":
			state=state*1664525+1013904223;key=int(state&31)
			state=state*1664525+1013904223;switch state%3{case 0:kind='A';case 1:kind='R';default:kind='Q'}
		default:return ops,false,state
		}
		if key<0||key>=32{return ops,false,state}
		ops=append(ops,wlmMultisetOp{Kind:kind,Key:key})
	}
	return ops,true,state
}
func wlmMultisetRowsEqual(a,b []wlmMultisetRow)bool{
	if len(a)!=len(b){return false}
	for i:=range a{x,y:=a[i],b[i];if x.Seed!=y.Seed||x.ScheduleFamily!=y.ScheduleFamily||x.InitialTotalMultiplicity!=y.InitialTotalMultiplicity||x.OperationBudget!=y.OperationBudget||x.Step!=y.Step||x.ExecutedOperations!=y.ExecutedOperations||x.AddedCount!=y.AddedCount||x.RejectedAddCount!=y.RejectedAddCount||x.RemovedCount!=y.RemovedCount||x.RemoveMissingCount!=y.RemoveMissingCount||x.QueryCount!=y.QueryCount||x.TotalBefore!=y.TotalBefore||x.TotalAfter!=y.TotalAfter||x.CumulativeAdded!=y.CumulativeAdded||x.CumulativeRejected!=y.CumulativeRejected||x.CumulativeRemoved!=y.CumulativeRemoved||!wlmMultisetEntriesEqual(x.Entries,y.Entries)||!wlmMultisetQueriesEqual(x.Queries,y.Queries){return false}}
	return true
}
func wlmMultisetTotalsEqual(a,b []wlmMultisetRow)bool{if len(a)!=len(b){return false};for i:=range a{if a[i].TotalBefore!=b[i].TotalBefore||a[i].TotalAfter!=b[i].TotalAfter{return false}};return true}
func wlmMultisetEntriesEqual(a,b []wlmMultisetEntry)bool{if len(a)!=len(b){return false};for i:=range a{if a[i]!=b[i]{return false}};return true}
func wlmMultisetQueriesEqual(a,b []wlmMultisetQuery)bool{if len(a)!=len(b){return false};for i:=range a{if a[i]!=b[i]{return false}};return true}
func wlmMultisetAbs(v int)int{if v<0{return -v};return v}
func wlmMultisetMaxMetric(m map[string]float64,k string,v int){if float64(v)>m[k]{m[k]=float64(v)}}
