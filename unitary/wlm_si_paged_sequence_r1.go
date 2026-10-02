package unitary

type wlmPagedRow struct {
	Seed uint32
	ScheduleFamily string
	InitialCardinality int
	OperationBudget int
	Substrate string
	Step int
	RequestedOperations int
	ExecutedOperations int
	AppendedCount int
	RejectedAppendCount int
	ReplacedCount int
	ReplaceMissCount int
	RemovedCount int
	RemoveMissCount int
	CardinalityBefore int
	CardinalityAfter int
	CumulativeAppended int
	CumulativeRejected int
	CumulativeReplaced int
	CumulativeRemoved int
	Values []int
	RemovedValues []int
}
type wlmPagedResult struct { Rows []wlmPagedRow; Metrics map[string]float64 }
type wlmPagedOp struct { Kind byte; Index int; Value int }
type wlmSequence interface { Len() int; Append(int) bool; Replace(int,int) bool; Remove(int)(int,bool); Snapshot() []int }

type wlmContiguousSequence struct{ v []int }
func(s *wlmContiguousSequence)Len()int{return len(s.v)}
func(s *wlmContiguousSequence)Append(v int)bool{if len(s.v)>=32{return false};s.v=append(s.v,v);return true}
func(s *wlmContiguousSequence)Replace(i,v int)bool{if i<0||i>=len(s.v){return false};s.v[i]=v;return true}
func(s *wlmContiguousSequence)Remove(i int)(int,bool){if i<0||i>=len(s.v){return 0,false};x:=s.v[i];copy(s.v[i:],s.v[i+1:]);s.v=s.v[:len(s.v)-1];return x,true}
func(s *wlmContiguousSequence)Snapshot()[]int{out:=make([]int,len(s.v));copy(out,s.v);return out}

type wlmFixedPagedSequence struct{ pages [4][8]int; length int }
func(s *wlmFixedPagedSequence)Len()int{return s.length}
func(s *wlmFixedPagedSequence)get(i int)int{return s.pages[i/8][i%8]}
func(s *wlmFixedPagedSequence)set(i,v int){s.pages[i/8][i%8]=v}
func(s *wlmFixedPagedSequence)Append(v int)bool{if s.length>=32{return false};s.set(s.length,v);s.length++;return true}
func(s *wlmFixedPagedSequence)Replace(i,v int)bool{if i<0||i>=s.length{return false};s.set(i,v);return true}
func(s *wlmFixedPagedSequence)Remove(i int)(int,bool){if i<0||i>=s.length{return 0,false};x:=s.get(i);for j:=i;j<s.length-1;j++{s.set(j,s.get(j+1))};s.length--;s.set(s.length,0);return x,true}
func(s *wlmFixedPagedSequence)Snapshot()[]int{out:=make([]int,s.length);for i:=range out{out[i]=s.get(i)};return out}

func RunWlmSiPagedSequenceR1()interface{}{
	seeds:=[...]uint32{1117,1151,1181,1213,1249,1283,1307,1361}
	families:=[...]string{"append-tail","replace-wave","remove-stride","lcg32"}
	initials:=[...]int{0,8,16,24}
	budgets:=[...]int{1,2,4,8}
	metrics:=map[string]float64{"paired_case_count":0,"completed_substrate_case_runs":0,"total_emitted_rows":0,"paired_observable_mismatch_cases":0,"row_count_mismatch_cases":0,"state_cardinality_mismatch_cases":0,"invalid_operation_schedules":0,"budget_overrun_rows":0,"law_violation_cases":0,"max_conservation_error":0,"max_residual_law_error":0}
	rows:=make([]wlmPagedRow,0,65536)
	for si,seed:=range seeds{for _,family:=range families{for _,initial:=range initials{for _,budget:=range budgets{
		metrics["paired_case_count"]++
		var a,b []wlmPagedRow;var badA,badB bool
		if si%2==0{a,badA=wlmPagedRun(seed,family,initial,budget,"contiguous-bounded-sequence",true,metrics);b,badB=wlmPagedRun(seed,family,initial,budget,"fixed-page-sequence",false,metrics)}else{b,badB=wlmPagedRun(seed,family,initial,budget,"fixed-page-sequence",false,metrics);a,badA=wlmPagedRun(seed,family,initial,budget,"contiguous-bounded-sequence",true,metrics)}
		rows=append(rows,a...);rows=append(rows,b...);metrics["completed_substrate_case_runs"]+=2
		if badA||badB{metrics["law_violation_cases"]++};if len(a)!=len(b){metrics["row_count_mismatch_cases"]++};if !wlmPagedRowsEqual(a,b){metrics["paired_observable_mismatch_cases"]++};if !wlmPagedCardinalitiesEqual(a,b){metrics["state_cardinality_mismatch_cases"]++}
	}}}}
	metrics["total_emitted_rows"]=float64(len(rows));return wlmPagedResult{Rows:rows,Metrics:metrics}
}
func wlmPagedRun(seed uint32,family string,initial,budget int,substrate string,contiguous bool,metrics map[string]float64)([]wlmPagedRow,bool){
	var seq wlmSequence;if contiguous{seq=&wlmContiguousSequence{}}else{seq=&wlmFixedPagedSequence{}}
	for i:=0;i<initial;i++{seq.Append(-initial+i)}
	rows:=make([]wlmPagedRow,0,64);ca,cr,cp,cd:=0,0,0,0;state:=seed;bad:=false
	for step:=0;step<64;step++{ops,ok,next:=wlmPagedOps(seed,state,family,step,budget);state=next;if !ok{metrics["invalid_operation_schedules"]++;return rows,true};before:=seq.Len();app,rej,rep,repMiss,rem,remMiss:=0,0,0,0,0,0;removed:=make([]int,0,budget);if len(ops)>budget{metrics["budget_overrun_rows"]++;bad=true}
		for _,op:=range ops{switch op.Kind{case 'A':if seq.Append(op.Value){app++}else{rej++};case 'R':if seq.Replace(op.Index,op.Value){rep++}else{repMiss++};case 'D':if v,ok:=seq.Remove(op.Index);ok{rem++;removed=append(removed,v)}else{remMiss++};default:metrics["invalid_operation_schedules"]++;bad=true}}
		after:=seq.Len();expected:=before+app-rem;conservation:=wlmPagedAbs(after-expected);residual:=wlmPagedAbs(len(seq.Snapshot())-after);wlmPagedMaxMetric(metrics,"max_conservation_error",conservation);wlmPagedMaxMetric(metrics,"max_residual_law_error",residual);if conservation!=0||residual!=0||after<0||after>32{bad=true};ca+=app;cr+=rej;cp+=rep;cd+=rem
		rows=append(rows,wlmPagedRow{Seed:seed,ScheduleFamily:family,InitialCardinality:initial,OperationBudget:budget,Substrate:substrate,Step:step,RequestedOperations:budget,ExecutedOperations:len(ops),AppendedCount:app,RejectedAppendCount:rej,ReplacedCount:rep,ReplaceMissCount:repMiss,RemovedCount:rem,RemoveMissCount:remMiss,CardinalityBefore:before,CardinalityAfter:after,CumulativeAppended:ca,CumulativeRejected:cr,CumulativeReplaced:cp,CumulativeRemoved:cd,Values:seq.Snapshot(),RemovedValues:removed})
	};return rows,bad
}
func wlmPagedOps(seed,state uint32,family string,step,budget int)([]wlmPagedOp,bool,uint32){ops:=make([]wlmPagedOp,0,budget);for ordinal:=0;ordinal<budget;ordinal++{var kind byte;idx:=0;value:=int(seed)+step*101+ordinal*7;switch family{case "append-tail":kind='A';case "replace-wave":if (step+ordinal)%5==0{kind='A'}else{kind='R'};idx=(step*3+ordinal*5)&31;case "remove-stride":if (step+ordinal)%3==0{kind='A'}else{kind='D'};idx=(step*5+ordinal*7)&31;case "lcg32":state=state*1664525+1013904223;switch state%3{case 0:kind='A';case 1:kind='R';default:kind='D'};state=state*1664525+1013904223;idx=int(state&31);state=state*1664525+1013904223;value=int(state&0x7fffffff);default:return ops,false,state};ops=append(ops,wlmPagedOp{Kind:kind,Index:idx,Value:value})};return ops,true,state}
func wlmPagedRowsEqual(a,b []wlmPagedRow)bool{if len(a)!=len(b){return false};for i:=range a{x,y:=a[i],b[i];if x.Seed!=y.Seed||x.ScheduleFamily!=y.ScheduleFamily||x.InitialCardinality!=y.InitialCardinality||x.OperationBudget!=y.OperationBudget||x.Step!=y.Step||x.RequestedOperations!=y.RequestedOperations||x.ExecutedOperations!=y.ExecutedOperations||x.AppendedCount!=y.AppendedCount||x.RejectedAppendCount!=y.RejectedAppendCount||x.ReplacedCount!=y.ReplacedCount||x.ReplaceMissCount!=y.ReplaceMissCount||x.RemovedCount!=y.RemovedCount||x.RemoveMissCount!=y.RemoveMissCount||x.CardinalityBefore!=y.CardinalityBefore||x.CardinalityAfter!=y.CardinalityAfter||x.CumulativeAppended!=y.CumulativeAppended||x.CumulativeRejected!=y.CumulativeRejected||x.CumulativeReplaced!=y.CumulativeReplaced||x.CumulativeRemoved!=y.CumulativeRemoved||!wlmPagedIntsEqual(x.Values,y.Values)||!wlmPagedIntsEqual(x.RemovedValues,y.RemovedValues){return false}};return true}
func wlmPagedCardinalitiesEqual(a,b []wlmPagedRow)bool{if len(a)!=len(b){return false};for i:=range a{if a[i].CardinalityBefore!=b[i].CardinalityBefore||a[i].CardinalityAfter!=b[i].CardinalityAfter||len(a[i].Values)!=a[i].CardinalityAfter||len(b[i].Values)!=b[i].CardinalityAfter{return false}};return true}
func wlmPagedIntsEqual(a,b []int)bool{if len(a)!=len(b){return false};for i:=range a{if a[i]!=b[i]{return false}};return true}
func wlmPagedAbs(v int)int{if v<0{return -v};return v}
func wlmPagedMaxMetric(m map[string]float64,k string,v int){if float64(v)>m[k]{m[k]=float64(v)}}
