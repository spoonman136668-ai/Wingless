package unitary

import "fmt"

const UPLM4NRevokeRestoreSchema="wingless.up-lm4n-budget-revoke-restore-transfer.v1"

type UPLM4NSchedule struct{
	CutRound int `json:"cut_round"`
	RestoreRound int `json:"restore_round"`
}
type UPLM4NSummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	CutRound int `json:"cut_round"`
	RestoreRound int `json:"restore_round"`
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
type UPLM4NResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM4MSeal string `json:"source_up_lm4m_seal"`
	DeadlineProfiles []string `json:"deadline_profiles"`
	Schedules []UPLM4NSchedule `json:"schedules"`
	RevokeAmounts []int `json:"revoke_amounts"`
	ResourceConfigurationsPerCell int `json:"resource_configurations_per_cell"`
	TopologyConditionCells int `json:"topology_condition_cells"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveScheduleUsed bool `json:"adaptive_schedule_used"`
	AdaptiveCoordinateUsed bool `json:"adaptive_coordinate_used"`
	Summaries []UPLM4NSummary `json:"summaries"`
}
func uplm4nCap(budget,round,cut,restore,revoke int)int{
	if round<cut||round>=restore{return budget}
	x:=budget-revoke;if x<0{return 0};return x
}
func uplm4nDynamicSchedule(budget,start,tp,cut,restore,revoke int)(reachable,end int){
	actions:=0;end=start-1
	for round:=0;round<6;round++{
		if round<start{continue}
		cap:=uplm4nCap(budget,round,cut,restore,revoke);before:=actions
		for k:=0;k<tp&&actions<cap;k++{actions++}
		if actions>before{end=round}
	}
	return actions,end
}
func uplm4nRun(rot int,perm,profile string,budget,start,tp,cut,restore,revoke int)int{
	arms:=uplm2yPermute(uplm2xArms(rot),perm);uplm3tPrepressure(arms,rot,profile)
	actions:=0
	for round:=0;round<6;round++{
		if round>=start{
			cap:=uplm4nCap(budget,round,cut,restore,revoke);used:=map[int]bool{}
			for k:=0;k<tp&&actions<cap;k++{
				i:=uplm3cChoose(arms,"earliest_deadline",used);if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++;used[i]=true}else{break}
			}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("4n-%s-%d-%d-%d-%s-%d-%d-%d-%d",profile,cut,restore,revoke,perm,budget,start,tp,round),"x")}
	}
	failed:=0
	for i:=range arms{_,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);failed+=f}
	return failed
}
func uplm4nKey(coord string,rs,tp,es,rd,ed int)string{
	if coord=="static_r+tp+end"{return fmt.Sprintf("r=%d|tp=%d|end=%d",rs,tp,es)}
	return fmt.Sprintf("r=%d|tp=%d|end=%d",rd,tp,ed)
}
func RunUPLM4N()(UPLM4NResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	schedules:=[]UPLM4NSchedule{{2,4},{2,5},{3,5}};revokes:=[]int{1,2}
	rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{1,2,3,4}
	coords:=[]string{"static_r+tp+end","dynamic_r+tp+end"}
	res:=UPLM4NResult{Schema:UPLM4NRevokeRestoreSchema,Experiment:"UP-LM4N-budget-revoke-restore-transfer",SourceUPLM4MSeal:"45ffd2f3cec82ef3f034020798434845610ad1c1",DeadlineProfiles:profiles,Schedules:schedules,RevokeAmounts:revokes,ResourceConfigurationsPerCell:len(budgets)*len(starts)*len(tps),TopologyConditionCells:len(profiles)*len(schedules)*len(revokes)*len(rots)*len(perms),CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveScheduleUsed:false,AdaptiveCoordinateUsed:false}
	type point struct{rs,tp,es,rd,ed,out int}
	type acc struct{min,max,n int}
	for _,profile:=range profiles{for _,sch:=range schedules{for _,revoke:=range revokes{for _,rot:=range rots{for _,perm:=range perms{
		points:=make([]point,0,res.ResourceConfigurationsPerCell)
		for _,budget:=range budgets{for _,start:=range starts{for _,tp:=range tps{
			rd,ed:=uplm4nDynamicSchedule(budget,start,tp,sch.CutRound,sch.RestoreRound,revoke)
			points=append(points,point{rs:uplm3fReachable(budget,start,tp),tp:tp,es:uplm4cCoverageEnd(budget,start,tp),rd:rd,ed:ed,out:uplm4nRun(rot,perm,profile,budget,start,tp,sch.CutRound,sch.RestoreRound,revoke)})
		}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range points{
				k:=uplm4nKey(coord,p.rs,p.tp,p.es,p.rd,p.ed);a:=g[k]
				if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a}
				a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
			}
			sm:=UPLM4NSummary{DeadlineProfile:profile,CutRound:sch.CutRound,RestoreRound:sch.RestoreRound,RevokeAmount:revoke,Rotation:rot,Permutation:perm,Coordinate:coord,Groups:len(g)};sum:=0
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
