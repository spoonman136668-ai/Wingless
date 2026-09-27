package unitary

import "fmt"

const UPLM4QHeldoutSchema="wingless.up-lm4q-heldout-revoke-restore-transfer.v1"

type UPLM4QSummary struct{
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
type UPLM4QResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM4PSeal string `json:"source_up_lm4p_seal"`
	HeldoutSchedules []UPLM4NSchedule `json:"heldout_schedules"`
	RevokeAmounts []int `json:"revoke_amounts"`
	TopologyConditionCells int `json:"topology_condition_cells"`
	ResourceConfigurationsPerCell int `json:"resource_configurations_per_cell"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveSelectionUsed bool `json:"adaptive_selection_used"`
	Summaries []UPLM4QSummary `json:"summaries"`
}
func uplm4qKey(coord string,r,tp,end,preRestore int)string{
	if coord=="final"{return fmt.Sprintf("r=%d|tp=%d|end=%d",r,tp,end)}
	return fmt.Sprintf("r=%d|tp=%d|end=%d|prerestore=%d",r,tp,end,preRestore)
}
func RunUPLM4Q()(UPLM4QResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	schedules:=[]UPLM4NSchedule{{1,3},{1,4},{1,5},{2,3},{3,4},{4,5}}
	revokes:=[]int{1,2};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{1,2,3,4}
	coords:=[]string{"final","final+prerestore"}
	res:=UPLM4QResult{Schema:UPLM4QHeldoutSchema,Experiment:"UP-LM4Q-heldout-revoke-restore-transfer",SourceUPLM4PSeal:"32880aedcc636afebbb7d305b50ee5d37e717efb",HeldoutSchedules:schedules,RevokeAmounts:revokes,TopologyConditionCells:len(profiles)*len(schedules)*len(revokes)*len(rots)*len(perms),ResourceConfigurationsPerCell:len(budgets)*len(starts)*len(tps),CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveSelectionUsed:false}
	type point struct{r,tp,end,preRestore,out int}
	type acc struct{min,max,n int}
	for _,profile:=range profiles{for _,sch:=range schedules{for _,revoke:=range revokes{for _,rot:=range rots{for _,perm:=range perms{
		points:=[]point{}
		for _,budget:=range budgets{for _,start:=range starts{for _,tp:=range tps{
			r,end,_,preRestore:=uplm4oSchedule(budget,start,tp,sch.CutRound,sch.RestoreRound,revoke)
			points=append(points,point{r:r,tp:tp,end:end,preRestore:preRestore,out:uplm4nRun(rot,perm,profile,budget,start,tp,sch.CutRound,sch.RestoreRound,revoke)})
		}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range points{
				k:=uplm4qKey(coord,p.r,p.tp,p.end,p.preRestore);a:=g[k]
				if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a}
				a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
			}
			sm:=UPLM4QSummary{DeadlineProfile:profile,CutRound:sch.CutRound,RestoreRound:sch.RestoreRound,RevokeAmount:revoke,Rotation:rot,Permutation:perm,Coordinate:coord,Groups:len(g)};sum:=0
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
