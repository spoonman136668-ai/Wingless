package unitary

import "fmt"

const UPLM4RDoubleSchema="wingless.up-lm4r-double-revoke-restore-transfer.v1"

type UPLM4RWindowSchedule struct{
	Cut1 int `json:"cut_1"`
	Restore1 int `json:"restore_1"`
	Cut2 int `json:"cut_2"`
	Restore2 int `json:"restore_2"`
}
type UPLM4RSummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	ScheduleIndex int `json:"schedule_index"`
	RevokeAmount int `json:"revoke_amount"`
	Rotation int `json:"rotation"`
	Permutation string `json:"permutation"`
	Coordinate string `json:"coordinate"`
	Groups int `json:"groups"`
	MultiMemberGroups int `json:"multi_member_groups"`
	NonzeroSpreadGroups int `json:"nonzero_spread_groups"`
	MaxFailureSpread int `json:"max_failure_spread"`
	MeanFailureSpread float64 `json:"mean_failure_spread"`
}
type UPLM4RResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM4QSeal string `json:"source_up_lm4q_seal"`
	Schedules []UPLM4RWindowSchedule `json:"schedules"`
	RevokeAmounts []int `json:"revoke_amounts"`
	TopologyConditionCells int `json:"topology_condition_cells"`
	ResourceConfigurationsPerCell int `json:"resource_configurations_per_cell"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveCheckpointSearchUsed bool `json:"adaptive_checkpoint_search_used"`
	Summaries []UPLM4RSummary `json:"summaries"`
}
func uplm4rCap(budget,round,revoke int,s UPLM4RWindowSchedule)int{
	revoked:=(round>=s.Cut1&&round<s.Restore1)||(round>=s.Cut2&&round<s.Restore2)
	if !revoked{return budget}
	x:=budget-revoke;if x<0{return 0};return x
}
func uplm4rSchedule(budget,start,tp,revoke int,s UPLM4RWindowSchedule)(reachable,end,preR1,preR2 int){
	actions:=0;end=start-1
	for round:=0;round<6;round++{
		if round==s.Restore1{preR1=actions}
		if round==s.Restore2{preR2=actions}
		if round<start{continue}
		cap:=uplm4rCap(budget,round,revoke,s);before:=actions
		for k:=0;k<tp&&actions<cap;k++{actions++}
		if actions>before{end=round}
	}
	return actions,end,preR1,preR2
}
func uplm4rRun(rot int,perm,profile string,budget,start,tp,revoke int,s UPLM4RWindowSchedule)int{
	arms:=uplm2yPermute(uplm2xArms(rot),perm);uplm3tPrepressure(arms,rot,profile)
	actions:=0
	for round:=0;round<6;round++{
		if round>=start{
			cap:=uplm4rCap(budget,round,revoke,s);used:=map[int]bool{}
			for k:=0;k<tp&&actions<cap;k++{
				i:=uplm3cChoose(arms,"earliest_deadline",used);if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++;used[i]=true}else{break}
			}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("4r-%s-%d-%d-%d-%d-%d-%s-%d-%d-%d-%d",profile,s.Cut1,s.Restore1,s.Cut2,s.Restore2,revoke,perm,budget,start,tp,round),"x")}
	}
	failed:=0
	for i:=range arms{_,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);failed+=f}
	return failed
}
func uplm4rKey(coord string,r,tp,end,preR1,preR2 int)string{
	switch coord{
	case "final":
		return fmt.Sprintf("r=%d|tp=%d|end=%d",r,tp,end)
	case "final+lastrestore":
		return fmt.Sprintf("r=%d|tp=%d|end=%d|r2=%d",r,tp,end,preR2)
	default:
		return fmt.Sprintf("r=%d|tp=%d|end=%d|r1=%d|r2=%d",r,tp,end,preR1,preR2)
	}
}
func RunUPLM4R()(UPLM4RResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	schedules:=[]UPLM4RWindowSchedule{{1,3,4,5},{2,3,4,5},{1,2,3,5}}
	revokes:=[]int{1,2};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{1,2,3,4}
	coords:=[]string{"final","final+lastrestore","final+bothrestores"}
	res:=UPLM4RResult{Schema:UPLM4RDoubleSchema,Experiment:"UP-LM4R-double-revoke-restore-transfer",SourceUPLM4QSeal:"8b3dfbc3253add2c2688f68dbaa9bc36ed945bfd",Schedules:schedules,RevokeAmounts:revokes,TopologyConditionCells:len(profiles)*len(schedules)*len(revokes)*len(rots)*len(perms),ResourceConfigurationsPerCell:len(budgets)*len(starts)*len(tps),CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveCheckpointSearchUsed:false}
	type point struct{r,tp,end,preR1,preR2,out int}
	type acc struct{min,max,n int}
	for _,profile:=range profiles{for si,s:=range schedules{for _,revoke:=range revokes{for _,rot:=range rots{for _,perm:=range perms{
		points:=[]point{}
		for _,budget:=range budgets{for _,start:=range starts{for _,tp:=range tps{
			r,end,preR1,preR2:=uplm4rSchedule(budget,start,tp,revoke,s)
			points=append(points,point{r:r,tp:tp,end:end,preR1:preR1,preR2:preR2,out:uplm4rRun(rot,perm,profile,budget,start,tp,revoke,s)})
		}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range points{
				k:=uplm4rKey(coord,p.r,p.tp,p.end,p.preR1,p.preR2);a:=g[k]
				if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a}
				a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
			}
			sm:=UPLM4RSummary{DeadlineProfile:profile,ScheduleIndex:si,RevokeAmount:revoke,Rotation:rot,Permutation:perm,Coordinate:coord,Groups:len(g)};sum:=0
			for _,a:=range g{
				if a.n>1{sm.MultiMemberGroups++}
				d:=a.max-a.min;sum+=d
				if d>0{sm.NonzeroSpreadGroups++};if d>sm.MaxFailureSpread{sm.MaxFailureSpread=d}
			}
			if len(g)>0{sm.MeanFailureSpread=float64(sum)/float64(len(g))}
			res.Summaries=append(res.Summaries,sm)
		}
	}}}}}
	return res,nil
}
