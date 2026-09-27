package unitary

import "fmt"

const UPLM4OPathSchema="wingless.up-lm4o-revoke-restore-path-checkpoints.v1"

type UPLM4OSummary struct{
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
type UPLM4OResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM4NBranch string `json:"source_up_lm4n_branch"`
	DeadlineProfiles []string `json:"deadline_profiles"`
	Schedules []UPLM4NSchedule `json:"schedules"`
	RevokeAmounts []int `json:"revoke_amounts"`
	ResourceConfigurationsPerCell int `json:"resource_configurations_per_cell"`
	TopologyConditionCells int `json:"topology_condition_cells"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveCheckpointSearchUsed bool `json:"adaptive_checkpoint_search_used"`
	AdaptiveCoordinateUsed bool `json:"adaptive_coordinate_used"`
	Summaries []UPLM4OSummary `json:"summaries"`
}
func uplm4oSchedule(budget,start,tp,cut,restore,revoke int)(reachable,end,preCut,preRestore int){
	actions:=0;end=start-1;cutSet,restoreSet:=false,false
	for round:=0;round<6;round++{
		if round==cut{preCut=actions;cutSet=true}
		if round==restore{preRestore=actions;restoreSet=true}
		if round<start{continue}
		cap:=uplm4nCap(budget,round,cut,restore,revoke);before:=actions
		for k:=0;k<tp&&actions<cap;k++{actions++}
		if actions>before{end=round}
	}
	if !cutSet{preCut=actions}
	if !restoreSet{preRestore=actions}
	return actions,end,preCut,preRestore
}
func uplm4oKey(coord string,r,tp,end,preCut,preRestore int)string{
	if coord=="dynamic_r+tp+end"{return fmt.Sprintf("r=%d|tp=%d|end=%d",r,tp,end)}
	return fmt.Sprintf("r=%d|tp=%d|end=%d|precut=%d|prerestore=%d",r,tp,end,preCut,preRestore)
}
func RunUPLM4O()(UPLM4OResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	schedules:=[]UPLM4NSchedule{{2,4},{2,5},{3,5}};revokes:=[]int{1,2}
	rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{1,2,3,4}
	coords:=[]string{"dynamic_r+tp+end","path_r+tp+end+precut+prerestore"}
	res:=UPLM4OResult{Schema:UPLM4OPathSchema,Experiment:"UP-LM4O-revoke-restore-path-checkpoints",SourceUPLM4NBranch:"research/wingless-up-lm4n-budget-revoke-restore-transfer-r1",DeadlineProfiles:profiles,Schedules:schedules,RevokeAmounts:revokes,ResourceConfigurationsPerCell:len(budgets)*len(starts)*len(tps),TopologyConditionCells:len(profiles)*len(schedules)*len(revokes)*len(rots)*len(perms),CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveCheckpointSearchUsed:false,AdaptiveCoordinateUsed:false}
	type point struct{r,tp,end,preCut,preRestore,out int}
	type acc struct{min,max,n int}
	for _,profile:=range profiles{for _,sch:=range schedules{for _,revoke:=range revokes{for _,rot:=range rots{for _,perm:=range perms{
		points:=make([]point,0,res.ResourceConfigurationsPerCell)
		for _,budget:=range budgets{for _,start:=range starts{for _,tp:=range tps{
			r,end,preCut,preRestore:=uplm4oSchedule(budget,start,tp,sch.CutRound,sch.RestoreRound,revoke)
			points=append(points,point{r:r,tp:tp,end:end,preCut:preCut,preRestore:preRestore,out:uplm4nRun(rot,perm,profile,budget,start,tp,sch.CutRound,sch.RestoreRound,revoke)})
		}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range points{
				k:=uplm4oKey(coord,p.r,p.tp,p.end,p.preCut,p.preRestore);a:=g[k]
				if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a}
				a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
			}
			sm:=UPLM4OSummary{DeadlineProfile:profile,CutRound:sch.CutRound,RestoreRound:sch.RestoreRound,RevokeAmount:revoke,Rotation:rot,Permutation:perm,Coordinate:coord,Groups:len(g)};sum:=0
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
