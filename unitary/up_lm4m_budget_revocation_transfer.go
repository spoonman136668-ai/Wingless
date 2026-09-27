package unitary

import "fmt"

const UPLM4MBudgetRevokeSchema="wingless.up-lm4m-budget-revocation-transfer.v1"

type UPLM4MSummary struct{
	DeadlineProfile string `json:"deadline_profile"`
	CutRound int `json:"cut_round"`
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
type UPLM4MResult struct{
	Schema string `json:"schema"`
	Experiment string `json:"experiment"`
	SourceUPLM4LSeal string `json:"source_up_lm4l_seal"`
	DeadlineProfiles []string `json:"deadline_profiles"`
	CutRounds []int `json:"cut_rounds"`
	RevokeAmounts []int `json:"revoke_amounts"`
	Rotations []int `json:"rotations"`
	Permutations []string `json:"permutations"`
	ResourceConfigurationsPerCell int `json:"resource_configurations_per_cell"`
	TopologyConditionCells int `json:"topology_condition_cells"`
	CandidateCoordinates []string `json:"candidate_coordinates"`
	CounterfactualOnly bool `json:"counterfactual_only"`
	LiveActivation bool `json:"live_activation"`
	AdaptiveRevocationUsed bool `json:"adaptive_revocation_used"`
	AdaptiveCoordinateUsed bool `json:"adaptive_coordinate_used"`
	Summaries []UPLM4MSummary `json:"summaries"`
}
func uplm4mCap(budget,round,cut,revoke int)int{
	if round<cut{return budget}
	x:=budget-revoke;if x<0{return 0};return x
}
func uplm4mDynamicSchedule(budget,start,tp,cut,revoke int)(reachable,end int){
	actions:=0;end=start-1
	for round:=0;round<6;round++{
		if round<start{continue}
		cap:=uplm4mCap(budget,round,cut,revoke)
		before:=actions
		for k:=0;k<tp&&actions<cap;k++{actions++}
		if actions>before{end=round}
	}
	return actions,end
}
func uplm4mRun(rot int,perm,profile string,budget,start,tp,cut,revoke int)int{
	arms:=uplm2yPermute(uplm2xArms(rot),perm);uplm3tPrepressure(arms,rot,profile)
	actions:=0
	for round:=0;round<6;round++{
		if round>=start{
			cap:=uplm4mCap(budget,round,cut,revoke);used:=map[int]bool{}
			for k:=0;k<tp&&actions<cap;k++{
				i:=uplm3cChoose(arms,"earliest_deadline",used);if i<0{break}
				if n,_,ok:=uplm2xFirstPending(&arms[i]);ok{arms[i].reported[n]=true;actions++;used[i]=true}else{break}
			}
		}
		for i:=range arms{arms[i].r.write(fmt.Sprintf("4m-%s-%d-%d-%s-%d-%d-%d-%d",profile,cut,revoke,perm,budget,start,tp,round),"x")}
	}
	failed:=0
	for i:=range arms{_,f:=uplm2vFinish(arms[i].r,arms[i].reported,rot);failed+=f}
	return failed
}
func uplm4mKey(coord string,rs,tp,es,rd,ed int)string{
	if coord=="static_r+tp+end"{return fmt.Sprintf("r=%d|tp=%d|end=%d",rs,tp,es)}
	return fmt.Sprintf("r=%d|tp=%d|end=%d",rd,tp,ed)
}
func RunUPLM4M()(UPLM4MResult,error){
	profiles:=[]string{"deferred_only","layout_only","hybrid_min"}
	cuts:=[]int{2,3,4};revokes:=[]int{1,2};rots:=[]int{5,13};perms:=[]string{"identity","reverse","rotate2"}
	budgets:=[]int{4,5,6,7};starts:=[]int{2,3,4,5};tps:=[]int{1,2,3,4}
	coords:=[]string{"static_r+tp+end","dynamic_r+tp+end"}
	res:=UPLM4MResult{Schema:UPLM4MBudgetRevokeSchema,Experiment:"UP-LM4M-budget-revocation-transfer",SourceUPLM4LSeal:"f598b8acf5c912d22c6f30ce53dbbe3714095157",DeadlineProfiles:profiles,CutRounds:cuts,RevokeAmounts:revokes,Rotations:rots,Permutations:perms,ResourceConfigurationsPerCell:len(budgets)*len(starts)*len(tps),TopologyConditionCells:len(profiles)*len(cuts)*len(revokes)*len(rots)*len(perms),CandidateCoordinates:coords,CounterfactualOnly:true,LiveActivation:false,AdaptiveRevocationUsed:false,AdaptiveCoordinateUsed:false}
	type point struct{rs,tp,es,rd,ed,out int}
	type acc struct{min,max,n int}
	for _,profile:=range profiles{for _,cut:=range cuts{for _,revoke:=range revokes{for _,rot:=range rots{for _,perm:=range perms{
		points:=make([]point,0,res.ResourceConfigurationsPerCell)
		for _,budget:=range budgets{for _,start:=range starts{for _,tp:=range tps{
			rd,ed:=uplm4mDynamicSchedule(budget,start,tp,cut,revoke)
			points=append(points,point{rs:uplm3fReachable(budget,start,tp),tp:tp,es:uplm4cCoverageEnd(budget,start,tp),rd:rd,ed:ed,out:uplm4mRun(rot,perm,profile,budget,start,tp,cut,revoke)})
		}}}
		for _,coord:=range coords{
			g:=map[string]*acc{}
			for _,p:=range points{
				k:=uplm4mKey(coord,p.rs,p.tp,p.es,p.rd,p.ed);a:=g[k]
				if a==nil{a=&acc{min:p.out,max:p.out};g[k]=a}
				a.n++;if p.out<a.min{a.min=p.out};if p.out>a.max{a.max=p.out}
			}
			sm:=UPLM4MSummary{DeadlineProfile:profile,CutRound:cut,RevokeAmount:revoke,Rotation:rot,Permutation:perm,Coordinate:coord,Groups:len(g)};sum:=0
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
