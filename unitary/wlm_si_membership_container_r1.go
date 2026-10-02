package unitary

type wlmMembershipQuery struct {
	Key int
	Present bool
}
type wlmMembershipRow struct {
	Seed uint32
	ScheduleFamily string
	InitialCardinality int
	OperationBudget int
	Substrate string
	Step int
	ExecutedOperations int
	AddedCount int
	AddExistingCount int
	DeletedCount int
	DeleteMissingCount int
	QueryHitCount int
	QueryMissCount int
	CardinalityBefore int
	CardinalityAfter int
	CumulativeAdded int
	CumulativeDeleted int
	Keys []int
	Queries []wlmMembershipQuery
}
type wlmMembershipResult struct { Rows []wlmMembershipRow; Metrics map[string]float64 }
type wlmMembershipOp struct { Kind byte; Key int }
type wlmMembership interface { Len() int; Add(int) bool; Delete(int) bool; Has(int) bool; Snapshot() []int }

type wlmMembershipDense struct{ present [64]bool; size int }
func(s *wlmMembershipDense)Len()int{return s.size}
func(s *wlmMembershipDense)Add(k int)bool{if s.present[k]{return false};s.present[k]=true;s.size++;return true}
func(s *wlmMembershipDense)Delete(k int)bool{if !s.present[k]{return false};s.present[k]=false;s.size--;return true}
func(s *wlmMembershipDense)Has(k int)bool{return s.present[k]}
func(s *wlmMembershipDense)Snapshot()[]int{out:=make([]int,0,s.size);for k:=0;k<64;k++{if s.present[k]{out=append(out,k)}};return out}

type wlmBitSet struct{ bits uint64 }
func(s *wlmBitSet)Len()int{n:=0;x:=s.bits;for x!=0{x&=x-1;n++};return n}
func(s *wlmBitSet)Add(k int)bool{mask:=uint64(1)<<uint(k);if s.bits&mask!=0{return false};s.bits|=mask;return true}
func(s *wlmBitSet)Delete(k int)bool{mask:=uint64(1)<<uint(k);if s.bits&mask==0{return false};s.bits&^=mask;return true}
func(s *wlmBitSet)Has(k int)bool{return s.bits&(uint64(1)<<uint(k))!=0}
func(s *wlmBitSet)Snapshot()[]int{out:=make([]int,0,s.Len());for k:=0;k<64;k++{if s.Has(k){out=append(out,k)}};return out}

func RunWlmSiMembershipContainerR1()interface{}{
	seeds:=[...]uint32{1411,1451,1481,1511,1543,1571,1601,1637}
	families:=[...]string{"sequential","toggle-wave","hot-cold","lcg32"}
	initials:=[...]int{0,16,32,48}
	budgets:=[...]int{1,2,4,8}
	metrics:=map[string]float64{"paired_case_count":0,"completed_substrate_case_runs":0,"total_emitted_rows":0,"paired_observable_mismatch_cases":0,"row_count_mismatch_cases":0,"state_cardinality_mismatch_cases":0,"invalid_operation_schedules":0,"budget_overrun_rows":0,"law_violation_cases":0,"max_conservation_error":0,"max_residual_law_error":0}
	rows:=make([]wlmMembershipRow,0,65536)
	for si,seed:=range seeds{for _,family:=range families{for _,initial:=range initials{for _,budget:=range budgets{
		metrics["paired_case_count"]++
		var a,b []wlmMembershipRow;var badA,badB bool
		if si%2==0{a,badA=wlmMembershipRun(seed,family,initial,budget,"dense-bool-membership-table",true,metrics);b,badB=wlmMembershipRun(seed,family,initial,budget,"uint64-bitset",false,metrics)}else{b,badB=wlmMembershipRun(seed,family,initial,budget,"uint64-bitset",false,metrics);a,badA=wlmMembershipRun(seed,family,initial,budget,"dense-bool-membership-table",true,metrics)}
		rows=append(rows,a...);rows=append(rows,b...);metrics["completed_substrate_case_runs"]+=2
		if badA||badB{metrics["law_violation_cases"]++};if len(a)!=len(b){metrics["row_count_mismatch_cases"]++};if !wlmMembershipRowsEqual(a,b){metrics["paired_observable_mismatch_cases"]++};if !wlmMembershipCardinalitiesEqual(a,b){metrics["state_cardinality_mismatch_cases"]++}
	}}}}
	metrics["total_emitted_rows"]=float64(len(rows));return wlmMembershipResult{Rows:rows,Metrics:metrics}
}
func wlmMembershipRun(seed uint32,family string,initial,budget int,substrate string,sorted bool,metrics map[string]float64)([]wlmMembershipRow,bool){
	var set wlmMembership;if sorted{set=&wlmMembershipDense{}}else{set=&wlmBitSet{}}
	for k:=0;k<initial;k++{set.Add(k)}
	rows:=make([]wlmMembershipRow,0,64);ca,cd:=0,0;state:=seed;bad:=false
	for step:=0;step<64;step++{
		ops,ok,next:=wlmMembershipOps(seed,state,family,step,budget);state=next
		if !ok{metrics["invalid_operation_schedules"]++;return rows,true}
		if len(ops)>budget{metrics["budget_overrun_rows"]++;bad=true}
		before:=set.Len();add,existing,del,missing,qhit,qmiss:=0,0,0,0,0,0
		queries:=make([]wlmMembershipQuery,0,budget)
		for _,op:=range ops{switch op.Kind{
		case 'A':if set.Add(op.Key){add++}else{existing++}
		case 'D':if set.Delete(op.Key){del++}else{missing++}
		case 'Q':hit:=set.Has(op.Key);if hit{qhit++}else{qmiss++};queries=append(queries,wlmMembershipQuery{Key:op.Key,Present:hit})
		default:metrics["invalid_operation_schedules"]++;bad=true}}
		after:=set.Len();conservation:=wlmMembershipAbs(after-(before+add-del));residual:=wlmMembershipAbs(len(set.Snapshot())-after)
		wlmMembershipMaxMetric(metrics,"max_conservation_error",conservation);wlmMembershipMaxMetric(metrics,"max_residual_law_error",residual)
		if conservation!=0||residual!=0||after<0||after>64{bad=true};ca+=add;cd+=del
		rows=append(rows,wlmMembershipRow{Seed:seed,ScheduleFamily:family,InitialCardinality:initial,OperationBudget:budget,Substrate:substrate,Step:step,ExecutedOperations:len(ops),AddedCount:add,AddExistingCount:existing,DeletedCount:del,DeleteMissingCount:missing,QueryHitCount:qhit,QueryMissCount:qmiss,CardinalityBefore:before,CardinalityAfter:after,CumulativeAdded:ca,CumulativeDeleted:cd,Keys:set.Snapshot(),Queries:queries})
	};return rows,bad
}
func wlmMembershipOps(seed,state uint32,family string,step,budget int)([]wlmMembershipOp,bool,uint32){
	ops:=make([]wlmMembershipOp,0,budget)
	for ordinal:=0;ordinal<budget;ordinal++{var kind byte;var key int
		switch family{
		case "sequential":key=(step*budget+ordinal)&63;switch step%3{case 0:kind='A';case 1:kind='Q';default:kind='D'}
		case "toggle-wave":key=(step+ordinal*7)&63;if step%2==0{kind='A'}else{kind='D'}
		case "hot-cold":if (step+ordinal)%2==0{key=(int(seed)+step+ordinal)&15}else{key=48+((int(seed)+step+ordinal)&15)};switch (step+ordinal)%3{case 0:kind='A';case 1:kind='Q';default:kind='D'}
		case "lcg32":state=state*1664525+1013904223;key=int(state&63);state=state*1664525+1013904223;switch state%3{case 0:kind='A';case 1:kind='D';default:kind='Q'}
		default:return ops,false,state}
		if key<0||key>=64{return ops,false,state};ops=append(ops,wlmMembershipOp{Kind:kind,Key:key})
	};return ops,true,state
}
func wlmMembershipRowsEqual(a,b []wlmMembershipRow)bool{if len(a)!=len(b){return false};for i:=range a{x,y:=a[i],b[i];if x.Seed!=y.Seed||x.ScheduleFamily!=y.ScheduleFamily||x.InitialCardinality!=y.InitialCardinality||x.OperationBudget!=y.OperationBudget||x.Step!=y.Step||x.ExecutedOperations!=y.ExecutedOperations||x.AddedCount!=y.AddedCount||x.AddExistingCount!=y.AddExistingCount||x.DeletedCount!=y.DeletedCount||x.DeleteMissingCount!=y.DeleteMissingCount||x.QueryHitCount!=y.QueryHitCount||x.QueryMissCount!=y.QueryMissCount||x.CardinalityBefore!=y.CardinalityBefore||x.CardinalityAfter!=y.CardinalityAfter||x.CumulativeAdded!=y.CumulativeAdded||x.CumulativeDeleted!=y.CumulativeDeleted||!wlmMembershipIntsEqual(x.Keys,y.Keys)||!wlmMembershipQueriesEqual(x.Queries,y.Queries){return false}};return true}
func wlmMembershipCardinalitiesEqual(a,b []wlmMembershipRow)bool{if len(a)!=len(b){return false};for i:=range a{if a[i].CardinalityBefore!=b[i].CardinalityBefore||a[i].CardinalityAfter!=b[i].CardinalityAfter||len(a[i].Keys)!=a[i].CardinalityAfter||len(b[i].Keys)!=b[i].CardinalityAfter{return false}};return true}
func wlmMembershipIntsEqual(a,b []int)bool{if len(a)!=len(b){return false};for i:=range a{if a[i]!=b[i]{return false}};return true}
func wlmMembershipQueriesEqual(a,b []wlmMembershipQuery)bool{if len(a)!=len(b){return false};for i:=range a{if a[i]!=b[i]{return false}};return true}
func wlmMembershipAbs(v int)int{if v<0{return -v};return v}
func wlmMembershipMaxMetric(m map[string]float64,k string,v int){if float64(v)>m[k]{m[k]=float64(v)}}
