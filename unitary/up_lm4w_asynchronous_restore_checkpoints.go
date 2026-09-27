package unitary

import "fmt"

const UPLM4WAsyncSchema="wingless.up-lm4w-asynchronous-restore-checkpoints.v1"

type UPLM4WWindow struct{
	BudgetCut int `json:"budget_cut"`
	BudgetRestore int `json:"budget_restore"`
	ThroughputCut int `json:"throughput_cut"`
	ThroughputRestore int `json:"throughput_restore"`
}
type UPLM4WSummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	BudgetCut int `json:"budget_cut"`
	BudgetRestore int `json:"budget_restore"`
	ThroughputCut int `json:"throughput_cut"`
	ThroughputRestore int `json:"throughput_restore"`
	BudgetReduction int `json:"budget_reduction"`
	ThroughputReduction int `json:"throughput_reduction"`
	Rotation int `json:"rotation"`
	Permutation string `json:"permutation"`
	Coordinate string `json:"coordinate"`
	Groups int `json:"groups"`
	MultiMemberGroups int `json:"multi_member_groups"`
	NonzeroSpreadGroups int `json:"nonzero_spread_groups"`
	MaxFailureSpread int `json:"max_failure_spread"`
	MeanFailureSpread float64 `json:"mean_failure_spread"`
}
type UPLM4WResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM4VSeal string `json:"source_up_lm4v_seal"`
	Windows []UPLM4WWindow `json:"windows"`
	TopologyConditionCells int `json:"topology_condition_cells"`
	ResourceConfigurationsPerCell int `json:"resource_configurations_per_cell"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveSelectionUsed bool `json:"adaptive_selection_used"`
	Summaries []UPLM4WSummary `json:"summaries"`
}
func uplm4wCap(budget,round,cut,restore,reduce int)int{
	if round<cut||round>=restore{return budget}
	x:=budget-reduce;if x<0{return 0};return x
}
func uplm4wTP(tp,round,cut,restore,reduce int)int{
	if round<cut||round>=restore{return tp}
	x:=tp-reduce;if x<0{return 0};return x
}
func uplm4wSchedule(budget,start,tp int,w UPLM4WWindow,bred,tred int)(reach,end,preBudget,preThroughput,preFull int){
	actions:=0;end=start-1
	fullRestore:=w.BudgetRestore;if w.ThroughputRestore>fullRestore{fullRestore=w.ThroughputRestore}
	for round:=0;round<6;round++{
		if round==w.BudgetRestore{preBudget=actions}
		if round==w.ThroughputRestore{preThroughput=actions}
		if round==fullRestore{preFull=actions}
		if round<start{continue}
		cap:=uplm4wCap(budget,round,w.BudgetCut,w.BudgetRestore,bred)
		eff:=uplm4wTP(tp,round,w.ThroughputCut,w.ThroughputRestore,tred)
		before:=actions
		for k:=0;k<eff&&actions<cap;k++{actions++}
		if actions>before{end=round}
	}
	return actions,end,preBudget,preThroughput,preFull
}
func uplm4wRun(rot int,perm,profile string,budget,start,tp int,w UPLM4WWindow,bred,tred int)int{
	arms:=uplm2yPermute(uplm2xArms(rot),perm);uplm3tPrepressure(arms,rot,profile)
	actions:=0
	for round:=0;round<6;round++{
		if round>=start{
			cap:=uplm4wCap(budget,round,w.BudgetCut,w.BudgetRestore,bred)
			eff:=uplm4wTP(tp,round,w.ThroughputCut,w.ThroughputRestore,tred)
			used:=map[int]bool{}
			for k:=0;k<eff&&actions<cap;k++{
				i:=uplm3cChoose(arms,"earliest_deadline",used);if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++;used[i]=true}else{break}
			}
		}
		for i:=range arms{
			arms[i].r.write(fmt.Sprintf("4w-%s-b%d-%d-t%d-%d-br%d-tr%d-%s-%d-%d-%d-%d",profile,w.BudgetCut,w.BudgetRestore,w.ThroughputCut,w.ThroughputRestore,bred,tred,perm,budget,start,tp,round),"x")
		}
	}
	failed:=0;for i:=range arms{_,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);failed+=f}
	return failed
}
func uplm4wKey(coord string,r,tp,end,preBudget,preThroughput,preFull int)string{
	switch coord{
	case "final":
		return fmt.Sprintf("r=%d|tp=%d|end=%d",r,tp,end)
	case "final+pre-full-restore":
		return fmt.Sprintf("r=%d|tp=%d|end=%d|pre=%d",r,tp,end,preFull)
	default:
		return fmt.Sprintf("r=%d|tp=%d|end=%d|pb=%d|pt=%d",r,tp,end,preBudget,preThroughput)
	}
}
func RunUPLM4W()(UPLM4WResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	windows:=[]UPLM4WWindow{{1,3,2,4},{2,4,1,3},{1,4,2,5},{2,5,1,4},{1,3,3,5},{3,5,1,3}}
	breds:=[]int{1,2};treds:=[]int{1,2}
	rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{2,3,4,5}
	coords:=[]string{"final","final+pre-full-restore","final+both-restores"}
	res:=UPLM4WResult{Schema:UPLM4WAsyncSchema,Experiment:"UP-LM4W-asynchronous-restore-checkpoints",SourceUPLM4VSeal:"7a0bf7dc75c59be1d2c83cda189ba17d263a80c5",Windows:windows,TopologyConditionCells:432,ResourceConfigurationsPerCell:64,CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveSelectionUsed:false}
	type point struct{r,tp,end,pb,pt,pf,out int};type acc struct{min,max,n int}
	for _,profile:=range profiles{for _,w:=range windows{for _,br:=range breds{for _,tr:=range treds{for _,rot:=range rots{for _,perm:=range perms{
		var pts []point
		for _,b:=range budgets{for _,st:=range starts{for _,tp:=range tps{
			r,end,pb,pt,pf:=uplm4wSchedule(b,st,tp,w,br,tr)
			pts=append(pts,point{r,tp,end,pb,pt,pf,uplm4wRun(rot,perm,profile,b,st,tp,w,br,tr)})
		}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range pts{
				k:=uplm4wKey(coord,p.r,p.tp,p.end,p.pb,p.pt,p.pf);a:=g[k]
				if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a}
				a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
			}
			sm:=UPLM4WSummary{DeadlineProfile:profile,BudgetCut:w.BudgetCut,BudgetRestore:w.BudgetRestore,ThroughputCut:w.ThroughputCut,ThroughputRestore:w.ThroughputRestore,BudgetReduction:br,ThroughputReduction:tr,Rotation:rot,Permutation:perm,Coordinate:coord,Groups:len(g)};sum:=0
			for _,a:=range g{if a.n>1{sm.MultiMemberGroups++};d:=a.max-a.min;sum+=d;if d>0{sm.NonzeroSpreadGroups++};if d>sm.MaxFailureSpread{sm.MaxFailureSpread=d}}
			if len(g)>0{sm.MeanFailureSpread=float64(sum)/float64(len(g))}
			res.Summaries=append(res.Summaries,sm)
		}
	}}}}}}
	return res,nil
}
