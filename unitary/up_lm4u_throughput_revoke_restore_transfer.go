package unitary

import "fmt"

const UPLM4UThroughputSchema="wingless.up-lm4u-throughput-revoke-restore-transfer.v1"

type UPLM4USummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	CutRound int `json:"cut_round"`
	RestoreRound int `json:"restore_round"`
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
type UPLM4UResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM4TSeal string `json:"source_up_lm4t_seal"`
	TopologyConditionCells int `json:"topology_condition_cells"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveSelectionUsed bool `json:"adaptive_selection_used"`
	Summaries []UPLM4USummary `json:"summaries"`
}
func uplm4uEffTP(tp,round,cut,restore,reduce int)int{
	if round<cut||round>=restore{return tp}
	x:=tp-reduce;if x<0{return 0};return x
}
func uplm4uSchedule(budget,start,tp,cut,restore,reduce int)(reach,end,preRestore int){
	actions:=0;end=start-1
	for round:=0;round<6;round++{
		if round==restore{preRestore=actions}
		if round<start{continue}
		eff:=uplm4uEffTP(tp,round,cut,restore,reduce);before:=actions
		for k:=0;k<eff&&actions<budget;k++{actions++}
		if actions>before{end=round}
	}
	return actions,end,preRestore
}
func uplm4uRun(rot int,perm,profile string,budget,start,tp,cut,restore,reduce int)int{
	arms:=uplm2yPermute(uplm2xArms(rot),perm);uplm3tPrepressure(arms,rot,profile);actions:=0
	for round:=0;round<6;round++{
		if round>=start{
			eff:=uplm4uEffTP(tp,round,cut,restore,reduce);used:=map[int]bool{}
			for k:=0;k<eff&&actions<budget;k++{
				i:=uplm3cChoose(arms,"earliest_deadline",used);if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++;used[i]=true}else{break}
			}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("4u-%s-%d-%d-%d-%s-%d-%d-%d-%d",profile,cut,restore,reduce,perm,budget,start,tp,round),"x")}
	}
	failed:=0;for i:=range arms{_,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);failed+=f};return failed
}
func uplm4uKey(coord string,r,tp,end,preRestore int)string{
	if coord=="final"{return fmt.Sprintf("r=%d|tp=%d|end=%d",r,tp,end)}
	return fmt.Sprintf("r=%d|tp=%d|end=%d|prerestore=%d",r,tp,end,preRestore)
}
func RunUPLM4U()(UPLM4UResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	schedules:=[]UPLM4NSchedule{{1,3},{2,4},{3,5}};reductions:=[]int{1,2}
	rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{2,3,4,5}
	coords:=[]string{"final","final+prerestore"}
	res:=UPLM4UResult{Schema:UPLM4UThroughputSchema,Experiment:"UP-LM4U-throughput-revoke-restore-transfer",SourceUPLM4TSeal:"3b6332e7c7e9cfb3e447ac24c899fd9b1f474ef6",TopologyConditionCells:108,CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveSelectionUsed:false}
	type point struct{r,tp,end,pre,out int};type acc struct{min,max,n int}
	for _,profile:=range profiles{for _,sch:=range schedules{for _,red:=range reductions{for _,rot:=range rots{for _,perm:=range perms{
		var pts []point
		for _,b:=range budgets{for _,st:=range starts{for _,tp:=range tps{
			r,end,pre:=uplm4uSchedule(b,st,tp,sch.CutRound,sch.RestoreRound,red)
			pts=append(pts,point{r,tp,end,pre,uplm4uRun(rot,perm,profile,b,st,tp,sch.CutRound,sch.RestoreRound,red)})
		}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range pts{k:=uplm4uKey(coord,p.r,p.tp,p.end,p.pre);a:=g[k];if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a};a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}}
			sm:=UPLM4USummary{DeadlineProfile:profile,CutRound:sch.CutRound,RestoreRound:sch.RestoreRound,ThroughputReduction:red,Rotation:rot,Permutation:perm,Coordinate:coord,Groups:len(g)};sum:=0
			for _,a:=range g{if a.n>1{sm.MultiMemberGroups++};d:=a.max-a.min;sum+=d;if d>0{sm.NonzeroSpreadGroups++};if d>sm.MaxFailureSpread{sm.MaxFailureSpread=d}}
			if len(g)>0{sm.MeanFailureSpread=float64(sum)/float64(len(g))};res.Summaries=append(res.Summaries,sm)
		}
	}}}}}
	return res,nil
}
