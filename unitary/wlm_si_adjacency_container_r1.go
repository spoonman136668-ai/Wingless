package unitary

type wlmAdjEdge struct { A int; B int }
type wlmAdjQuery struct { A int; B int; Present bool }
type wlmAdjRow struct {
	Seed uint32
	ScheduleFamily string
	InitialEdgeCount int
	OperationBudget int
	Substrate string
	Step int
	ExecutedOperations int
	AddedCount int
	RejectedAddCount int
	DeletedCount int
	DeleteMissingCount int
	QueryHitCount int
	QueryMissCount int
	EdgeCountBefore int
	EdgeCountAfter int
	CumulativeAdded int
	CumulativeRejected int
	CumulativeDeleted int
	Degrees []int
	Edges []wlmAdjEdge
	Queries []wlmAdjQuery
}
type wlmAdjResult struct { Rows []wlmAdjRow; Metrics map[string]float64 }
type wlmAdjOp struct { Kind byte; A int; B int }
type wlmAdjacency interface {
	EdgeCount() int
	Add(int,int) bool
	Delete(int,int) bool
	Has(int,int) bool
	Degrees() []int
	Snapshot() []wlmAdjEdge
}

type wlmAdjMatrix struct { edge [32][32]bool; count int }
func(g *wlmAdjMatrix)EdgeCount()int{return g.count}
func(g *wlmAdjMatrix)Add(a,b int)bool{if a==b||g.edge[a][b]||g.count>=96{return false};g.edge[a][b]=true;g.edge[b][a]=true;g.count++;return true}
func(g *wlmAdjMatrix)Delete(a,b int)bool{if a==b||!g.edge[a][b]{return false};g.edge[a][b]=false;g.edge[b][a]=false;g.count--;return true}
func(g *wlmAdjMatrix)Has(a,b int)bool{return a!=b&&g.edge[a][b]}
func(g *wlmAdjMatrix)Degrees()[]int{d:=make([]int,32);for a:=0;a<32;a++{for b:=0;b<32;b++{if g.edge[a][b]{d[a]++}}};return d}
func(g *wlmAdjMatrix)Snapshot()[]wlmAdjEdge{out:=make([]wlmAdjEdge,0,g.count);for a:=0;a<32;a++{for b:=a+1;b<32;b++{if g.edge[a][b]{out=append(out,wlmAdjEdge{A:a,B:b})}}};return out}

type wlmAdjLists struct { neighbors [32][]int; count int }
func(g *wlmAdjLists)EdgeCount()int{return g.count}
func wlmAdjFind(xs []int,v int)(int,bool){lo,hi:=0,len(xs);for lo<hi{m:=(lo+hi)/2;if xs[m]<v{lo=m+1}else{hi=m}};return lo,lo<len(xs)&&xs[lo]==v}
func wlmAdjInsert(xs []int,i,v int)[]int{xs=append(xs,0);copy(xs[i+1:],xs[i:]);xs[i]=v;return xs}
func wlmAdjRemove(xs []int,i int)[]int{copy(xs[i:],xs[i+1:]);return xs[:len(xs)-1]}
func(g *wlmAdjLists)Add(a,b int)bool{
	if a==b||g.count>=96{return false}
	ia,ok:=wlmAdjFind(g.neighbors[a],b);if ok{return false}
	ib,_:=wlmAdjFind(g.neighbors[b],a)
	g.neighbors[a]=wlmAdjInsert(g.neighbors[a],ia,b);g.neighbors[b]=wlmAdjInsert(g.neighbors[b],ib,a);g.count++;return true
}
func(g *wlmAdjLists)Delete(a,b int)bool{
	if a==b{return false}
	ia,ok:=wlmAdjFind(g.neighbors[a],b);if !ok{return false}
	ib,ok2:=wlmAdjFind(g.neighbors[b],a);if !ok2{panic("asymmetric adjacency")}
	g.neighbors[a]=wlmAdjRemove(g.neighbors[a],ia);g.neighbors[b]=wlmAdjRemove(g.neighbors[b],ib);g.count--;return true
}
func(g *wlmAdjLists)Has(a,b int)bool{if a==b{return false};_,ok:=wlmAdjFind(g.neighbors[a],b);return ok}
func(g *wlmAdjLists)Degrees()[]int{d:=make([]int,32);for i:=0;i<32;i++{d[i]=len(g.neighbors[i])};return d}
func(g *wlmAdjLists)Snapshot()[]wlmAdjEdge{out:=make([]wlmAdjEdge,0,g.count);for a:=0;a<32;a++{for _,b:=range g.neighbors[a]{if a<b{out=append(out,wlmAdjEdge{A:a,B:b})}}};return out}

func RunWlmSiAdjacencyContainerR1()interface{}{
	seeds:=[...]uint32{2081,2111,2137,2161,2203,2237,2267,2293}
	families:=[...]string{"chain","star-wave","rewire","lcg32"}
	initials:=[...]int{0,24,48,72}
	budgets:=[...]int{1,2,4,8}
	metrics:=map[string]float64{"paired_case_count":0,"completed_substrate_case_runs":0,"total_emitted_rows":0,"paired_observable_mismatch_cases":0,"row_count_mismatch_cases":0,"state_cardinality_mismatch_cases":0,"invalid_operation_schedules":0,"budget_overrun_rows":0,"law_violation_cases":0,"max_conservation_error":0,"max_residual_law_error":0}
	rows:=make([]wlmAdjRow,0,65536)
	for si,seed:=range seeds{for _,family:=range families{for _,initial:=range initials{for _,budget:=range budgets{
		metrics["paired_case_count"]++
		var a,b []wlmAdjRow;var badA,badB bool
		if si%2==0{a,badA=wlmAdjRun(seed,family,initial,budget,"adjacency-matrix",true,metrics);b,badB=wlmAdjRun(seed,family,initial,budget,"sorted-adjacency-lists",false,metrics)}else{b,badB=wlmAdjRun(seed,family,initial,budget,"sorted-adjacency-lists",false,metrics);a,badA=wlmAdjRun(seed,family,initial,budget,"adjacency-matrix",true,metrics)}
		rows=append(rows,a...);rows=append(rows,b...);metrics["completed_substrate_case_runs"]+=2
		if badA||badB{metrics["law_violation_cases"]++};if len(a)!=len(b){metrics["row_count_mismatch_cases"]++};if !wlmAdjRowsEqual(a,b){metrics["paired_observable_mismatch_cases"]++};if !wlmAdjCardinalitiesEqual(a,b){metrics["state_cardinality_mismatch_cases"]++}
	}}}}
	metrics["total_emitted_rows"]=float64(len(rows));return wlmAdjResult{Rows:rows,Metrics:metrics}
}
func wlmAdjInitialEdges(seed uint32,count int)[]wlmAdjEdge{
	all:=make([]wlmAdjEdge,0,496);for a:=0;a<32;a++{for b:=a+1;b<32;b++{all=append(all,wlmAdjEdge{A:a,B:b})}}
	offset:=int(seed%uint32(len(all)));out:=make([]wlmAdjEdge,0,count);for i:=0;i<count;i++{out=append(out,all[(offset+i)%len(all)])};return out
}
func wlmAdjRun(seed uint32,family string,initial,budget int,substrate string,matrix bool,metrics map[string]float64)([]wlmAdjRow,bool){
	var g wlmAdjacency;if matrix{g=&wlmAdjMatrix{}}else{g=&wlmAdjLists{}}
	for _,e:=range wlmAdjInitialEdges(seed,initial){if !g.Add(e.A,e.B){panic("initial edge admission failed")}}
	rows:=make([]wlmAdjRow,0,64);ca,cr,cd:=0,0,0;state:=seed;bad:=false
	for step:=0;step<64;step++{
		ops,ok,next:=wlmAdjOps(seed,state,family,step,budget);state=next;if !ok{metrics["invalid_operation_schedules"]++;return rows,true};if len(ops)>budget{metrics["budget_overrun_rows"]++;bad=true}
		before:=g.EdgeCount();added,rejected,deleted,missing,qhit,qmiss:=0,0,0,0,0,0;queries:=make([]wlmAdjQuery,0,budget)
		for _,op:=range ops{switch op.Kind{case 'A':if g.Add(op.A,op.B){added++}else{rejected++};case 'D':if g.Delete(op.A,op.B){deleted++}else{missing++};case 'Q':hit:=g.Has(op.A,op.B);if hit{qhit++}else{qmiss++};queries=append(queries,wlmAdjQuery{A:op.A,B:op.B,Present:hit});default:metrics["invalid_operation_schedules"]++;bad=true}}
		after:=g.EdgeCount();conservation:=wlmAdjAbs(after-(before+added-deleted));edges:=g.Snapshot();degrees:=g.Degrees();sumDegree:=0;for _,d:=range degrees{sumDegree+=d};residual:=wlmAdjAbs(sumDegree-2*after)
		wlmAdjMaxMetric(metrics,"max_conservation_error",conservation);wlmAdjMaxMetric(metrics,"max_residual_law_error",residual)
		if conservation!=0||residual!=0||after<0||after>96||len(edges)!=after{bad=true}
		ca+=added;cr+=rejected;cd+=deleted
		rows=append(rows,wlmAdjRow{Seed:seed,ScheduleFamily:family,InitialEdgeCount:initial,OperationBudget:budget,Substrate:substrate,Step:step,ExecutedOperations:len(ops),AddedCount:added,RejectedAddCount:rejected,DeletedCount:deleted,DeleteMissingCount:missing,QueryHitCount:qhit,QueryMissCount:qmiss,EdgeCountBefore:before,EdgeCountAfter:after,CumulativeAdded:ca,CumulativeRejected:cr,CumulativeDeleted:cd,Degrees:degrees,Edges:edges,Queries:queries})
	}
	return rows,bad
}
func wlmAdjOps(seed,state uint32,family string,step,budget int)([]wlmAdjOp,bool,uint32){
	ops:=make([]wlmAdjOp,0,budget)
	for ordinal:=0;ordinal<budget;ordinal++{var kind byte;var a,b int
		switch family{
		case "chain":a=(step+ordinal)&31;b=(a+1)&31;switch step%3{case 0:kind='A';case 1:kind='Q';default:kind='D'}
		case "star-wave":a=int(seed&31);b=(step+ordinal+1)&31;if b==a{b=(b+1)&31};switch(step+ordinal)%3{case 0:kind='A';case 1:kind='Q';default:kind='D'}
		case "rewire":a=(step+ordinal*3)&31;b=(a+7+step)&31;if b==a{b=(b+1)&31};if step%2==0{kind='D'}else{kind='A'}
		case "lcg32":state=state*1664525+1013904223;a=int(state&31);state=state*1664525+1013904223;b=int(state&31);if b==a{b=(b+1)&31};state=state*1664525+1013904223;switch state%3{case 0:kind='A';case 1:kind='D';default:kind='Q'}
		default:return ops,false,state
		}
		if a<0||a>=32||b<0||b>=32||a==b{return ops,false,state};if a>b{a,b=b,a};ops=append(ops,wlmAdjOp{Kind:kind,A:a,B:b})
	}
	return ops,true,state
}
func wlmAdjRowsEqual(a,b []wlmAdjRow)bool{if len(a)!=len(b){return false};for i:=range a{x,y:=a[i],b[i];if x.Seed!=y.Seed||x.ScheduleFamily!=y.ScheduleFamily||x.InitialEdgeCount!=y.InitialEdgeCount||x.OperationBudget!=y.OperationBudget||x.Step!=y.Step||x.ExecutedOperations!=y.ExecutedOperations||x.AddedCount!=y.AddedCount||x.RejectedAddCount!=y.RejectedAddCount||x.DeletedCount!=y.DeletedCount||x.DeleteMissingCount!=y.DeleteMissingCount||x.QueryHitCount!=y.QueryHitCount||x.QueryMissCount!=y.QueryMissCount||x.EdgeCountBefore!=y.EdgeCountBefore||x.EdgeCountAfter!=y.EdgeCountAfter||x.CumulativeAdded!=y.CumulativeAdded||x.CumulativeRejected!=y.CumulativeRejected||x.CumulativeDeleted!=y.CumulativeDeleted||!wlmAdjIntsEqual(x.Degrees,y.Degrees)||!wlmAdjEdgesEqual(x.Edges,y.Edges)||!wlmAdjQueriesEqual(x.Queries,y.Queries){return false}};return true}
func wlmAdjCardinalitiesEqual(a,b []wlmAdjRow)bool{if len(a)!=len(b){return false};for i:=range a{if a[i].EdgeCountBefore!=b[i].EdgeCountBefore||a[i].EdgeCountAfter!=b[i].EdgeCountAfter{return false}};return true}
func wlmAdjIntsEqual(a,b []int)bool{if len(a)!=len(b){return false};for i:=range a{if a[i]!=b[i]{return false}};return true}
func wlmAdjEdgesEqual(a,b []wlmAdjEdge)bool{if len(a)!=len(b){return false};for i:=range a{if a[i]!=b[i]{return false}};return true}
func wlmAdjQueriesEqual(a,b []wlmAdjQuery)bool{if len(a)!=len(b){return false};for i:=range a{if a[i]!=b[i]{return false}};return true}
func wlmAdjAbs(v int)int{if v<0{return -v};return v}
func wlmAdjMaxMetric(m map[string]float64,k string,v int){if float64(v)>m[k]{m[k]=float64(v)}}
