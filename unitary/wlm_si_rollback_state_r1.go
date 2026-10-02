package unitary

type wlmRollbackEntry struct { Key int; Value int }
type wlmRollbackQuery struct { Key int; Present bool; Value int }
type wlmRollbackRow struct {
	Seed uint32
	ScheduleFamily string
	InitialCardinality int
	OperationBudget int
	Substrate string
	Step int
	Control string
	DepthBefore int
	DepthAfterControl int
	DepthAfter int
	CardinalityBefore int
	CardinalityAfterControl int
	CardinalityAfter int
	SetInsertCount int
	SetUpdateCount int
	DeleteHitCount int
	DeleteMissCount int
	QueryCount int
	CumulativeInsert int
	CumulativeUpdate int
	CumulativeDelete int
	Entries []wlmRollbackEntry
	Queries []wlmRollbackQuery
}
type wlmRollbackResult struct { Rows []wlmRollbackRow; Metrics map[string]float64 }
type wlmRollbackOp struct { Kind byte; Key int; Value int }

type wlmRollbackState interface {
	Depth() int
	Begin() bool
	Commit() bool
	Rollback() bool
	Set(int,int)(bool,bool)
	Delete(int) bool
	Get(int)(int,bool)
	Cardinality() int
	Snapshot() []wlmRollbackEntry
}

type wlmFullSnapshot struct {
	present [32]bool
	values [32]int
	size int
	frames []wlmFullFrame
}
type wlmFullFrame struct { present [32]bool; values [32]int; size int }
func(s *wlmFullSnapshot)Depth()int{return len(s.frames)}
func(s *wlmFullSnapshot)Begin()bool{if len(s.frames)>=2{return false};s.frames=append(s.frames,wlmFullFrame{present:s.present,values:s.values,size:s.size});return true}
func(s *wlmFullSnapshot)Commit()bool{if len(s.frames)==0{return false};s.frames=s.frames[:len(s.frames)-1];return true}
func(s *wlmFullSnapshot)Rollback()bool{if len(s.frames)==0{return false};f:=s.frames[len(s.frames)-1];s.frames=s.frames[:len(s.frames)-1];s.present=f.present;s.values=f.values;s.size=f.size;return true}
func(s *wlmFullSnapshot)Set(k,v int)(bool,bool){if s.present[k]{s.values[k]=v;return false,true};s.present[k]=true;s.values[k]=v;s.size++;return true,false}
func(s *wlmFullSnapshot)Delete(k int)bool{if !s.present[k]{return false};s.present[k]=false;s.values[k]=0;s.size--;return true}
func(s *wlmFullSnapshot)Get(k int)(int,bool){if !s.present[k]{return 0,false};return s.values[k],true}
func(s *wlmFullSnapshot)Cardinality()int{return s.size}
func(s *wlmFullSnapshot)Snapshot()[]wlmRollbackEntry{out:=make([]wlmRollbackEntry,0,s.size);for k:=0;k<32;k++{if s.present[k]{out=append(out,wlmRollbackEntry{Key:k,Value:s.values[k]})}};return out}

type wlmUndoEntry struct { key int; present bool; value int }
type wlmUndoLog struct {
	present [32]bool
	values [32]int
	size int
	frames [][]wlmUndoEntry
}
func(s *wlmUndoLog)Depth()int{return len(s.frames)}
func(s *wlmUndoLog)Begin()bool{if len(s.frames)>=2{return false};s.frames=append(s.frames,nil);return true}
func(s *wlmUndoLog)Commit()bool{
	if len(s.frames)==0{return false}
	last:=len(s.frames)-1
	child:=s.frames[last]
	s.frames=s.frames[:last]
	if len(s.frames)>0{s.frames[len(s.frames)-1]=append(s.frames[len(s.frames)-1],child...)}
	return true
}
func(s *wlmUndoLog)Rollback()bool{
	if len(s.frames)==0{return false}
	last:=len(s.frames)-1
	frame:=s.frames[last]
	s.frames=s.frames[:last]
	for i:=len(frame)-1;i>=0;i--{u:=frame[i];now:=s.present[u.key];if now&&!u.present{s.size--};if !now&&u.present{s.size++};s.present[u.key]=u.present;s.values[u.key]=u.value}
	return true
}
func(s *wlmUndoLog)record(k int){if len(s.frames)>0{s.frames[len(s.frames)-1]=append(s.frames[len(s.frames)-1],wlmUndoEntry{key:k,present:s.present[k],value:s.values[k]})}}
func(s *wlmUndoLog)Set(k,v int)(bool,bool){s.record(k);if s.present[k]{s.values[k]=v;return false,true};s.present[k]=true;s.values[k]=v;s.size++;return true,false}
func(s *wlmUndoLog)Delete(k int)bool{if !s.present[k]{return false};s.record(k);s.present[k]=false;s.values[k]=0;s.size--;return true}
func(s *wlmUndoLog)Get(k int)(int,bool){if !s.present[k]{return 0,false};return s.values[k],true}
func(s *wlmUndoLog)Cardinality()int{return s.size}
func(s *wlmUndoLog)Snapshot()[]wlmRollbackEntry{out:=make([]wlmRollbackEntry,0,s.size);for k:=0;k<32;k++{if s.present[k]{out=append(out,wlmRollbackEntry{Key:k,Value:s.values[k]})}};return out}

func RunWlmSiRollbackStateR1()interface{}{
	seeds:=[...]uint32{2521,2551,2579,2609,2633,2659,2683,2707}
	families:=[...]string{"commit-heavy","rollback-heavy","nested-alternating","lcg32"}
	initials:=[...]int{0,8,16,24}
	budgets:=[...]int{1,2,4,8}
	metrics:=map[string]float64{"paired_case_count":0,"completed_substrate_case_runs":0,"total_emitted_rows":0,"paired_observable_mismatch_cases":0,"row_count_mismatch_cases":0,"state_cardinality_mismatch_cases":0,"invalid_operation_schedules":0,"budget_overrun_rows":0,"law_violation_cases":0,"max_conservation_error":0,"max_residual_law_error":0}
	rows:=make([]wlmRollbackRow,0,65536)
	for si,seed:=range seeds{for _,family:=range families{for _,initial:=range initials{for _,budget:=range budgets{
		metrics["paired_case_count"]++
		ops:=wlmRollbackOps(seed,family,budget)
		var a,b []wlmRollbackRow;var badA,badB bool
		if si%2==0{a,badA=wlmRollbackRun(seed,family,initial,budget,ops,"full-copy-snapshot",true,metrics);b,badB=wlmRollbackRun(seed,family,initial,budget,ops,"nested-undo-log",false,metrics)}else{b,badB=wlmRollbackRun(seed,family,initial,budget,ops,"nested-undo-log",false,metrics);a,badA=wlmRollbackRun(seed,family,initial,budget,ops,"full-copy-snapshot",true,metrics)}
		rows=append(rows,a...);rows=append(rows,b...);metrics["completed_substrate_case_runs"]+=2
		if badA||badB{metrics["law_violation_cases"]++};if len(a)!=len(b){metrics["row_count_mismatch_cases"]++};if !wlmRollbackRowsEqual(a,b){metrics["paired_observable_mismatch_cases"]++};if !wlmRollbackCardsEqual(a,b){metrics["state_cardinality_mismatch_cases"]++}
	}}}}
	metrics["total_emitted_rows"]=float64(len(rows));return wlmRollbackResult{Rows:rows,Metrics:metrics}
}
func wlmRollbackOps(seed uint32,family string,budget int)[][]wlmRollbackOp{
	out:=make([][]wlmRollbackOp,64);state:=seed
	for step:=0;step<64;step++{row:=make([]wlmRollbackOp,0,budget);for ordinal:=0;ordinal<budget;ordinal++{var kind byte;var key,value int
		switch family{
		case "commit-heavy":key=(step*budget+ordinal)&31;if (step+ordinal)%4==3{kind='D'}else if (step+ordinal)%4==2{kind='Q'}else{kind='S'};value=int(seed)+step*97+ordinal
		case "rollback-heavy":key=(step+ordinal*7+int(seed))&31;if (step+ordinal)%3==0{kind='D'}else{kind='S'};value=int(seed)+step*71+ordinal*5
		case "nested-alternating":key=(step*3+ordinal*11+int(seed))&31;switch(step+ordinal)%3{case 0:kind='S';case 1:kind='Q';default:kind='D'};value=int(seed)+step*53+ordinal*13
		case "lcg32":state=state*1664525+1013904223;key=int(state&31);state=state*1664525+1013904223;switch state%3{case 0:kind='S';case 1:kind='D';default:kind='Q'};state=state*1664525+1013904223;value=int(state&0x7fffffff)
		default:panic("unknown family")}
		row=append(row,wlmRollbackOp{Kind:kind,Key:key,Value:value})
	};out[step]=row};return out
}
func wlmRollbackControl(family string,step int,depth int,seed uint32)byte{
	switch family{
	case "commit-heavy":switch step%4{case 0:return 'B';case 1:if depth<2{return 'B'};return 'C';default:return 'C'}
	case "rollback-heavy":switch step%4{case 0:return 'B';case 1:if depth<2{return 'B'};return 'R';default:return 'R'}
	case "nested-alternating":switch step%4{case 0:return 'B';case 1:if depth<2{return 'B'};return 'C';case 2:return 'C';default:return 'R'}
	case "lcg32":
		x:=uint32(step)*1103515245+seed*12345
		if depth==0{return 'B'}
		if depth==2{if x&1==0{return 'C'};return 'R'}
		switch x%3{case 0:return 'B';case 1:return 'C';default:return 'R'}
	}
	return 'N'
}
func wlmRollbackRun(seed uint32,family string,initial,budget int,ops [][]wlmRollbackOp,substrate string,full bool,metrics map[string]float64)([]wlmRollbackRow,bool){
	var s wlmRollbackState;if full{s=&wlmFullSnapshot{}}else{s=&wlmUndoLog{}}
	for k:=0;k<initial;k++{s.Set(k,int(seed)+k*17)}
	rows:=make([]wlmRollbackRow,0,64);ci,cu,cd:=0,0,0;bad:=false
	for step:=0;step<64;step++{
		before:=s.Cardinality();depthBefore:=s.Depth();ctl:=wlmRollbackControl(family,step,depthBefore,seed);valid:=true
		switch ctl{case 'B':valid=s.Begin();case 'C':valid=s.Commit();case 'R':valid=s.Rollback();case 'N':default:valid=false}
		if !valid{metrics["invalid_operation_schedules"]++;bad=true}
		afterCtl:=s.Cardinality();depthCtl:=s.Depth();inserted,updated,deleted,missing,qcount:=0,0,0,0,0;queries:=make([]wlmRollbackQuery,0,budget)
		if len(ops[step])>budget{metrics["budget_overrun_rows"]++;bad=true}
		for _,op:=range ops[step]{switch op.Kind{case 'S':ins,upd:=s.Set(op.Key,op.Value);if ins{inserted++};if upd{updated++};case 'D':if s.Delete(op.Key){deleted++}else{missing++};case 'Q':v,ok:=s.Get(op.Key);queries=append(queries,wlmRollbackQuery{Key:op.Key,Present:ok,Value:v});qcount++;default:metrics["invalid_operation_schedules"]++;bad=true}}
		after:=s.Cardinality();conservation:=wlmRollbackAbs(after-(afterCtl+inserted-deleted));snapshot:=s.Snapshot();residual:=wlmRollbackAbs(len(snapshot)-after);wlmRollbackMaxMetric(metrics,"max_conservation_error",conservation);wlmRollbackMaxMetric(metrics,"max_residual_law_error",residual);if conservation!=0||residual!=0||s.Depth()<0||s.Depth()>2{bad=true}
		ci+=inserted;cu+=updated;cd+=deleted
		rows=append(rows,wlmRollbackRow{Seed:seed,ScheduleFamily:family,InitialCardinality:initial,OperationBudget:budget,Substrate:substrate,Step:step,Control:string([]byte{ctl}),DepthBefore:depthBefore,DepthAfterControl:depthCtl,DepthAfter:s.Depth(),CardinalityBefore:before,CardinalityAfterControl:afterCtl,CardinalityAfter:after,SetInsertCount:inserted,SetUpdateCount:updated,DeleteHitCount:deleted,DeleteMissCount:missing,QueryCount:qcount,CumulativeInsert:ci,CumulativeUpdate:cu,CumulativeDelete:cd,Entries:snapshot,Queries:queries})
	}
	return rows,bad
}
func wlmRollbackRowsEqual(a,b []wlmRollbackRow)bool{if len(a)!=len(b){return false};for i:=range a{x,y:=a[i],b[i];if x.Seed!=y.Seed||x.ScheduleFamily!=y.ScheduleFamily||x.InitialCardinality!=y.InitialCardinality||x.OperationBudget!=y.OperationBudget||x.Step!=y.Step||x.Control!=y.Control||x.DepthBefore!=y.DepthBefore||x.DepthAfterControl!=y.DepthAfterControl||x.DepthAfter!=y.DepthAfter||x.CardinalityBefore!=y.CardinalityBefore||x.CardinalityAfterControl!=y.CardinalityAfterControl||x.CardinalityAfter!=y.CardinalityAfter||x.SetInsertCount!=y.SetInsertCount||x.SetUpdateCount!=y.SetUpdateCount||x.DeleteHitCount!=y.DeleteHitCount||x.DeleteMissCount!=y.DeleteMissCount||x.QueryCount!=y.QueryCount||x.CumulativeInsert!=y.CumulativeInsert||x.CumulativeUpdate!=y.CumulativeUpdate||x.CumulativeDelete!=y.CumulativeDelete||!wlmRollbackEntriesEqual(x.Entries,y.Entries)||!wlmRollbackQueriesEqual(x.Queries,y.Queries){return false}};return true}
func wlmRollbackCardsEqual(a,b []wlmRollbackRow)bool{if len(a)!=len(b){return false};for i:=range a{if a[i].CardinalityBefore!=b[i].CardinalityBefore||a[i].CardinalityAfterControl!=b[i].CardinalityAfterControl||a[i].CardinalityAfter!=b[i].CardinalityAfter{return false}};return true}
func wlmRollbackEntriesEqual(a,b []wlmRollbackEntry)bool{if len(a)!=len(b){return false};for i:=range a{if a[i]!=b[i]{return false}};return true}
func wlmRollbackQueriesEqual(a,b []wlmRollbackQuery)bool{if len(a)!=len(b){return false};for i:=range a{if a[i]!=b[i]{return false}};return true}
func wlmRollbackAbs(v int)int{if v<0{return -v};return v}
func wlmRollbackMaxMetric(m map[string]float64,k string,v int){if float64(v)>m[k]{m[k]=float64(v)}}
